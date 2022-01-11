// fake uses local disk as storage backend so as to test upload/download
package fake

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"

	"github.com/magicleap/tfmodules/pkg/backends/types"
)

var (
	// returned module metadata
	CodeRepo = "https://whatever.com/wherever.git"
	Metadata = []types.Metadata{
		types.Metadata{Version: "0.0.1", Source: CodeRepo},
		types.Metadata{Version: "1.2.3", Source: CodeRepo},
		types.Metadata{Version: "0.1.2", Source: CodeRepo},
	}
	// test file to download
	DownloadFile = "testModule.tar.gz"
	// test file uploaded
	UploadedFile string

	// getting root path of the package to access test files
	_, b, _, _  = runtime.Caller(0)
	basepath    = filepath.Dir(b)
	StoragePath = path.Join(basepath, "fake_storage")
)

type Fake struct{}

func (f *Fake) ListMetadata(namespace, name, system string) ([]types.Metadata, error) {
	failedTest := "unexisting"
	if namespace == failedTest || name == failedTest || system == failedTest {
		return []types.Metadata{}, fmt.Errorf("failed to list %s", failedTest)
	}
	return Metadata, nil
}

func (f *Fake) Write(namespace string, name string, system string, version string, source string, body io.Reader) error {
	UploadedFile = fmt.Sprintf("testfile-%s.tar.gz", version)
	tf, err := os.Create(path.Join(StoragePath, UploadedFile))
	if err != nil {
		return err
	}
	defer tf.Close()

	if _, err := io.Copy(tf, body); err != nil {
		return err
	}
	return nil
}

func (f *Fake) Read(namespace string, name string, system string, version string, response io.Writer) error {

	tf, err := os.Open(path.Join(StoragePath, DownloadFile))
	if err != nil {
		return err
	}
	defer tf.Close()

	if _, err := io.Copy(response, tf); err != nil {
		return err
	}

	return nil
}
