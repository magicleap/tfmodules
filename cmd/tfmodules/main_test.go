package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/deepmap/oapi-codegen/pkg/testutil"
	"github.com/go-chi/chi/v5"
	"github.com/magicleap/tfmodules/pkg/backends/fake"
	"github.com/magicleap/tfmodules/pkg/modules"
	"github.com/magicleap/tfmodules/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// related to cmd
	receivedFilePath = "../../test/receivedModule.tar.gz"
	sentFilePath     = path.Join(fake.StoragePath, fake.DownloadFile)
)

// doGet helper to GET requests
func doGet(t *testing.T, mux *chi.Mux, url string) *httptest.ResponseRecorder {
	return testutil.NewRequest().Get(url).WithAcceptJson().GoWithHTTPHandler(t, mux).Recorder
}

// compare the size of the files
func assertSameFiles(t *testing.T, file1 string, file2 string) {
	t.Logf("file1 %s", file1)
	s1, err := utils.GetFileSize(file1)
	require.NoError(t, err)
	t.Logf("file2 %s", file2)
	s2, err := utils.GetFileSize(file2)
	require.NoError(t, err)
	assert.Equal(t, s1, s2)
}

// check that the content type of the reposnse is JSON
func assertJsonContentType(t *testing.T, rr *httptest.ResponseRecorder) {
	header := rr.Result().Header.Get("Content-Type")
	assert.Equal(t, "application/json", header)
}

func TestModuleRegistry(t *testing.T) {
	var err error

	// get handler for fake backend
	r, err := getHandler("fake")
	require.NoError(t, err)

	modFile, err := ioutil.ReadFile(sentFilePath)
	require.NoError(t, err)

	t.Run("Test discovery", func(t *testing.T) {
		restore := utils.EnvSetAndReset(map[string]string{"MODULE_PATH": "/whatever"}, false)

		rr := doGet(t, r, "/.well-known/terraform.json")
		assert.Equal(t, http.StatusOK, rr.Code)
		t.Logf("message : %v", rr.Body)

		discovery := make(map[string]string)
		err = json.NewDecoder(rr.Body).Decode(&discovery)
		assert.NoError(t, err)
		expected := map[string]string{"modules.v1": "/whatever"}
		assert.Equal(t, discovery, expected)
		assertJsonContentType(t, rr)
		t.Cleanup(func() { restore() })
	})

	t.Run("List Successfully Versions", func(t *testing.T) {
		rr := doGet(t, r, "/namespace/module/sys/versions")
		assert.Equal(t, http.StatusOK, rr.Code)

		var moduleRegistry modules.ModuleRegistry
		err = json.NewDecoder(rr.Body).Decode(&moduleRegistry)
		assert.NoError(t, err)
		// analyse response from fake server that sends 3 versions
		module := (*moduleRegistry.Modules)[0]
		// As we use pointers, ensure that the items are different
		versions := *module.Versions
		assert.Equal(t, 3, len(versions))
		assert.NotEqual(t, *versions[0].Version, *versions[1].Version)
		assertJsonContentType(t, rr)
	})

	t.Run("Failed to List Versions", func(t *testing.T) {
		rr := doGet(t, r, "/namespace/unexisting/sys/versions")
		assert.Equal(t, http.StatusTeapot, rr.Code)
	})

	t.Run("Get latest Version", func(t *testing.T) {
		rr := doGet(t, r, "/namespace/module/sys")
		assert.Equal(t, http.StatusOK, rr.Code)
		var moduleDetails modules.ModuleDetails
		err = json.NewDecoder(rr.Body).Decode(&moduleDetails)
		assert.NoError(t, err)
		// analyse response from fake server : expected highest version
		assert.Equal(t, "1.2.3", *moduleDetails.Version)
		assert.Equal(t, "https://whatever.com/wherever.git", *moduleDetails.Source)
		assert.Equal(t, 3, len(*moduleDetails.Versions))
		t.Logf("Versions %v\n", *moduleDetails.Versions)

		assert.Equal(t, "module", *moduleDetails.Name)
		assert.Equal(t, "namespace", *moduleDetails.Namespace)
		assert.Equal(t, "sys", *moduleDetails.Provider)

	})

	t.Run("Get Download link", func(t *testing.T) {
		rr := doGet(t, r, "/namespace/module/sys/0.0.1/download")
		assert.Equal(t, http.StatusNoContent, rr.Code)
		header := rr.Result().Header.Get("X-Terraform-Get")
		assert.Equal(t, header, "/namespace/module/sys/0.0.1/archive.tgz")
		assertJsonContentType(t, rr)
	})

	t.Run("Download file", func(t *testing.T) {
		// dowload the file locally
		rr := doGet(t, r, "/namespace/module/sys/0.0.1/archive.tgz")
		assert.Equal(t, http.StatusOK, rr.Code)
		out, err := os.Create(receivedFilePath)
		require.NoError(t, err)
		_, err = io.Copy(out, rr.Body) // equivalent to curl -o
		require.NoError(t, err)
		// compare the size of the files (received and orginal in fake storage)
		assertSameFiles(t, receivedFilePath, sentFilePath)
		err = os.Remove(receivedFilePath)
		require.NoError(t, err)
	})

	t.Run("Upload file without module source", func(t *testing.T) {
		restore := utils.EnvSetAndReset(map[string]string{"OVERWRITE": "1"}, false)
		// upload the file
		rr := testutil.NewRequest().Post("/namespace/module/sys/1.2.3").WithBody(modFile).GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusCreated, rr.Code)
		uploadedFile := path.Join(fake.StoragePath, fake.UploadedFile)
		assertSameFiles(t, sentFilePath, uploadedFile)
		// finally delete the file
		err = os.Remove(uploadedFile)
		require.NoError(t, err)
		t.Cleanup(func() { restore() })
	})

	t.Run("Upload file with module source", func(t *testing.T) {
		restore := utils.EnvSetAndReset(map[string]string{"OVERWRITE": "1"}, false)
		// upload the file
		rr := testutil.NewRequest().Post("/namespace/module/sys/1.2.3").WithBody(modFile).WithHeader("module-source", "https://whatever.com/wherever.git").GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusCreated, rr.Code)
		uploadedFile := path.Join(fake.StoragePath, fake.UploadedFile)
		assertSameFiles(t, sentFilePath, uploadedFile)
		// finally delete the file
		err = os.Remove(uploadedFile)
		require.NoError(t, err)
		t.Cleanup(func() { restore() })
	})

	t.Run("Upload file with existing version", func(t *testing.T) {
		restore := utils.EnvSetAndReset(map[string]string{"OVERWRITE": "0"}, false)
		// upload the file
		rr := testutil.NewRequest().Post("/namespace/module/sys/1.2.3").WithBody(modFile).WithHeader("module-source", "https://whatever.com/wherever.git").GoWithHTTPHandler(t, r).Recorder
		assert.Equal(t, http.StatusForbidden, rr.Code)
		t.Cleanup(func() { restore() })
	})
}
