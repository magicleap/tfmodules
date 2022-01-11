package modules

import (
	"fmt"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// GetDownloadLink returns an empty content but a header `X-Terraform-Get`` with the URL to actually download the file
func (m *ModuleServer) GetDownloadLink(w http.ResponseWriter, r *http.Request, namespace Namespace, name Name, system System, version Version) {
	w.Header().Set("X-Terraform-Get", fmt.Sprintf("/%s/%s/%s/%s/archive.tgz", namespace, name, system, version))
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
}

// Download actually downloads module source
// ie it returns the content of the file to the requester
func (m *ModuleServer) Download(w http.ResponseWriter, r *http.Request, namespace Namespace, name Name, system System, version Version) {

	if err := m.Backend.Read(string(namespace), string(name), string(system), string(version), w); err != nil {
		logrus.WithError(err).Error("failed downloading file")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func containsVersion(versions []string, version Version) bool {
	for _, v := range versions {
		if v == string(version) {
			return true
		}
	}
	return false
}

// Upload uploads the tarball into the storage backend
// It can also optionnaly receive the URL of the module source as header
// TODO: check that module file is actually tgz
func (m *ModuleServer) Upload(w http.ResponseWriter, r *http.Request, namespace Namespace, name Name, system System, version Version, params UploadParams) {
	logrus.Debugf("received a request to upload a file %v/%v/%v/%v", namespace, name, system, version)

	// If we preventing overwriting a version, we check first that it does not exist yet
	if os.Getenv("OVERWRITE") == "0" {
		_, versions, err := m.fetchMetadata(namespace, name, system)
		if err != nil {
			logrus.WithError(err).Error("failed fetching versions")
			w.WriteHeader(http.StatusTeapot)
			return
		}

		if containsVersion(versions, version) {
			logrus.Warning("version already exists in storage")
			w.WriteHeader(http.StatusForbidden)
			return
		}
	}

	var source string
	if params.ModuleSource != nil {
		source = string(*params.ModuleSource)
		logrus.Debugf("module source is %s", source)
	}

	if err := m.Backend.Write(string(namespace), string(name), string(system), string(version), source, r.Body); err != nil {
		logrus.WithError(err).Error("failed uploading file")
		w.WriteHeader(http.StatusTeapot)
		return
	}

	w.WriteHeader(http.StatusCreated)

}
