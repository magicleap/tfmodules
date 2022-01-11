package backends

import (
	"errors"
	"fmt"
	"testing"

	"github.com/magicleap/tfmodules/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestBackendConfigurations test the envars requirements to initiate backends
func TestBackendConfigurations(t *testing.T) {

	t.Run("Wrong Backend Name", func(t *testing.T) {
		_, err := GetBackend("unexisting")
		require.Error(t, err)
	})

	// TODO: loop over envar for next backends
	e := "GOOGLE_BUCKET"
	var configError *utils.ConfigurationError

	t.Run(fmt.Sprintf("GCS Missing %s", e), func(t *testing.T) {
		restore := utils.EnvSetAndReset(map[string]string{e: ""}, true)
		_, err := GetBackend("gcs")
		assert.ErrorAs(t, err, &configError)
		t.Cleanup(func() { restore() })
	})

	t.Run(fmt.Sprintf("GCS Correct %s", e), func(t *testing.T) {
		restore := utils.EnvSetAndReset(map[string]string{e: "my-bucket"}, false)
		_, err := GetBackend("gcs")
		// We do not want to receive a configError, but we may receive an error from actual gcs backend
		// if not connected
		if errors.As(err, &configError) {
			t.Failed()
		}
		t.Cleanup(func() { restore() })
	})
}
