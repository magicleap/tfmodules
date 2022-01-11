package backends

import (
	"fmt"
	"io"

	"github.com/magicleap/tfmodules/pkg/backends/fake"
	"github.com/magicleap/tfmodules/pkg/backends/gcs"
	"github.com/magicleap/tfmodules/pkg/backends/types"
	"github.com/magicleap/tfmodules/pkg/utils"
)

// ServerInterface represents all server handlers.
type BackendInterface interface {
	// Returns metadata of each stored files (versions, sources, ...)
	ListMetadata(namespace string, name string, system string) ([]types.Metadata, error)
	// Upload tar.gz file into storage
	Write(namespace string, name string, system string, version string, source string, body io.Reader) error
	// Download tar.gz file from storage
	Read(namespace string, name string, system string, version string, response io.Writer) error
}

// GetBackend ensures that the envar are correctly set and return the expected implementation
func GetBackend(backend string) (BackendInterface, error) {
	switch backend {

	case "gcs":
		if err := utils.RequireEnvar("GOOGLE_BUCKET"); err != nil {
			return nil, err
		}
		return gcs.New()

	case "fake":
		return &fake.Fake{}, nil

	default:
		return nil, fmt.Errorf("%s is not a supported backend", backend)
	}
}
