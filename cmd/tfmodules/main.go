package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"

	"github.com/magicleap/tfmodules/pkg/backends"
	"github.com/magicleap/tfmodules/pkg/modules"
	"github.com/magicleap/tfmodules/pkg/utils"
)

var (
	version = "undefined"
)

// getHandler is an helper to initiate the modules handler
func getHandler(backend string) (*chi.Mux, error) {

	// initiate the storage backend
	b, err := backends.GetBackend(backend)
	if err != nil {
		return nil, fmt.Errorf("failed to initiate backend : %w", err)
	}

	// register modules as the handler
	module := &modules.ModuleServer{Backend: b}

	r := chi.NewRouter()
	r.Mount("/", modules.Handler(module))

	return r, nil
}

func main() {

	logrus.Infof("Running version %s", version)
	if os.Getenv("VERBOSE") == "1" {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.SetReportCaller(true)
		logrus.Debug("Let's talk !")
	}

	utils.SetDefaultValue("PORT", "8080")
	utils.SetDefaultValue("LISTEN", "0.0.0.0")
	utils.SetDefaultValue("BACKEND", "gcs")
	utils.SetDefaultValue("OVERWRITE", "0")
	utils.SetDefaultValue("MODULE_PATH", "/")

	addr := fmt.Sprintf("%s:%s", os.Getenv("LISTEN"), os.Getenv("PORT"))
	logrus.Infof("Listening to %s", addr)

	handler, err := getHandler(os.Getenv("BACKEND"))
	if err != nil {
		logrus.WithError(err).Fatal("error initiating http server")
	}

	srv := &http.Server{
		Handler: handler,
		Addr:    addr,
	}
	logrus.Fatal(srv.ListenAndServe())
}
