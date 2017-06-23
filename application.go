package main

import (
	"go/build"
	"log"
	"os"
	"time"

	graceful "gopkg.in/tylerb/graceful.v1"

	"net/http/pprof"

	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/iAmPlus/TrackCards/handlers"
	negronilogrus "github.com/iAmPlus/negroni-logrus"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	newrelic "github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	newrelic_negroni "github.com/yadvendar/negroni-newrelic-go-agent"
)

func attachProfiler(router *mux.Router) {
	router.HandleFunc("/debug/pprof/", pprof.Index)
	router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
}

func attachHandlers(router *mux.Router) {
	router.HandleFunc("/health", handlers.Health)
}

func main() {
	router, n, recovery, logger := mux.NewRouter().StrictSlash(true), negroni.New(), negroni.NewRecovery(), logrus.New()
	w := logger.WriterLevel(logrus.ErrorLevel)
	defer w.Close()

	logger.Level = logrus.InfoLevel
	logger.Formatter = &logrus.TextFormatter{ForceColors: true}
	recovery.PrintStack = false
	recovery.Logger = log.New(w, "[negroni] ", 0)

	if os.Getenv("GO_ENV") != "production" {
		err := godotenv.Load(".env")
		if err != nil {
			importPath := "github.com/iAmPlus/TrackCards" // modify to match import path of main
			p, err := build.Default.Import(importPath, "", build.FindOnly)
			certsDir := filepath.Join(p.Dir, ".env")
			err = godotenv.Load(certsDir)
			if err != nil {
				logger.Error("Error loading .env file: %v", err)
			}
		}
	} else {
		logger.Formatter = logrustash.DefaultFormatter(logrus.Fields{"type": "TrackCards"})
		config := newrelic.NewConfig(os.Getenv("NEWRELIC_APPLICATION_NAME"), os.Getenv("NEWRELIC_LICENSE_KEY"))
		config.Enabled = true
		newRelicMiddleware, err := newrelic_negroni.New(config)
		if err != nil {

		}

		n.Use(newRelicMiddleware)
	}

	attachProfiler(router)
	attachHandlers(router)

	n.Use(recovery)
	n.Use(negronilogrus.NewMiddlewareFromLogger(logger, "web"))
	n.UseHandler(router)

	port := ":" + os.Getenv("PORT")
	if port == ":" {
		port = ":8080"
	}
	logger.Printf("Running on port: %v", port)
	graceful.Run(port, 10*time.Second, n)
}
