package utils

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

// EnvSetAndReset allows to set envars values, then to reset them to original values after the tests
// if unset boolean is set to true, the envar is unset instead of being updated
func EnvSetAndReset(envs map[string]string, unset bool) (reset func()) {

	// record original values
	originalEnvs := map[string]string{}
	for name, value := range envs {
		if originalValue, ok := os.LookupEnv(name); ok {
			originalEnvs[name] = originalValue
		}
		if unset {
			// we unset the envar from the os
			_ = os.Unsetenv(name)
		} else {
			// we set the value of the envar
			_ = os.Setenv(name, value)
		}

	}

	// helper restoring the initial values
	return func() {
		for name := range envs {
			origValue, has := originalEnvs[name]
			if has {
				_ = os.Setenv(name, origValue)
			} else {
				_ = os.Unsetenv(name)
			}
		}
	}
}

// SetDefaultValues checks if an envar exists and set a value if not
func SetDefaultValue(envar string, value string) {
	// if not set, we set the default value
	if len(os.Getenv(envar)) == 0 {
		os.Setenv(envar, value)
		logrus.Infof("Defaulting %s to %s", envar, os.Getenv((envar)))
	}
}

// ConfigurationError is an error sent when an envar is missing
type ConfigurationError struct {
	env string
}

func (ce *ConfigurationError) Error() string {
	return fmt.Sprintf("required envar %s", ce.env)
}

// helper to ensure that mandatory envar is set
func RequireEnvar(key string) error {
	if len(os.Getenv(key)) == 0 {
		return &ConfigurationError{env: key}
	}
	return nil
}
