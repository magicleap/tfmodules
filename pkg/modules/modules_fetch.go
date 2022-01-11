package modules

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/magicleap/tfmodules/pkg/backends/types"
	"github.com/sirupsen/logrus"
)

// fetchMetadata retrieves the metadata of the objects from backend
// the metadata is sorted by version descending
func (m *ModuleServer) fetchMetadata(namespace Namespace, name Name, system System) (latest types.Metadata, versions []string, err error) {

	metadata, err := m.Backend.ListMetadata(string(namespace), string(name), string(system))
	if err != nil {
		return latest, versions, err
	}

	// sort the metadata by version descending, so that the first element is highest version
	sort.Slice(metadata, func(i, j int) bool { return metadata[i].Version > metadata[j].Version })
	// get a list of the available versions form the metadata
	for m := range metadata {
		versions = append(versions, metadata[m].Version)
	}

	return metadata[0], versions, nil
}

// GetLatestVersion returns a list of 1 element contain the latest version of a module
func (m *ModuleServer) GetLatestVersion(w http.ResponseWriter, r *http.Request, namespace Namespace, name Name, system System) {

	latest, versions, err := m.fetchMetadata(namespace, name, system)
	if err != nil {
		// TODO: return as message to API response
		w.WriteHeader(http.StatusTeapot)
		return
	}

	name_str := string(name)
	namespace_str := string(namespace)
	system_str := string(system)
	moduleDetails := &ModuleDetails{
		Versions:  &versions,
		Name:      &name_str,
		Namespace: &namespace_str,
		Provider:  &system_str,
		Version:   &latest.Version,
		Source:    &latest.Source,
	}

	// FIXME: would have thought it was managed by openAPI
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(moduleDetails); err != nil {
		// TODO: return as message to API response
		logrus.WithError(err).Error("failed to encode JSON")
		w.WriteHeader(http.StatusTeapot)
		return
	}

}

// ListVersions lists available versions for a given module
func (m *ModuleServer) ListVersions(w http.ResponseWriter, r *http.Request, namespace Namespace, name Name, system System) {

	_, versions, err := m.fetchMetadata(namespace, name, system)
	if err != nil {
		// TODO: return as message to API response
		w.WriteHeader(http.StatusTeapot)
		return
	}

	moduleVersions := []ModuleVersion{}
	for i := range versions {
		moduleVersions = append(moduleVersions, ModuleVersion{Version: &versions[i]})
	}
	module := &Module{Versions: &moduleVersions}
	modules := []Module{*module}
	mr := &ModuleRegistry{Modules: &modules}

	// FIXME: would have thought it was managed by openAPI
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(mr); err != nil {
		// TODO: return as message to API response
		logrus.WithError(err).Error("failed to encode JSON")
		w.WriteHeader(http.StatusTeapot)
		return
	}

}
