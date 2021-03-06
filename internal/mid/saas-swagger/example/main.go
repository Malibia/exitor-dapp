package main

import (
	"context"
	"geeks-accelerator/oss/saas-starter-kit/internal/platform/web/webcontext"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"geeks-accelerator/oss/saas-starter-kit/internal/mid"
	saasSwagger "geeks-accelerator/oss/saas-starter-kit/internal/mid/saas-swagger"
	_ "geeks-accelerator/oss/saas-starter-kit/internal/mid/saas-swagger/example/docs" // docs is generated by Swag CLI, you have to import it.
	"geeks-accelerator/oss/saas-starter-kit/internal/platform/flag"
	"geeks-accelerator/oss/saas-starter-kit/internal/platform/web"
	"github.com/kelseyhightower/envconfig"
)

// build is the git version of this program. It is set using build flags in the makefile.
var build = "develop"

// service is the name of the program used for logging, tracing and the
// the prefix used for loading env variables
// ie: export WEB_API_ENV=dev
var service = "EXAMPLE_API"

// @title SaaS Example API
// @version 1.0
// @description This is a sample server celler server.
// @termsOfService http://geeksinthewoods.com/terms

// @contact.name API Support
// @contact.email support@geeksinthewoods.com
// @contact.url https://gitlab.com/geeks-accelerator/oss/saas-starter-kit

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host example-api.saas.geeksinthewoods.com
// @BasePath /v1

func main() {

	// =========================================================================
	// Logging

	log := log.New(os.Stdout, service+" : ", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

	// =========================================================================
	// Configuration
	var cfg struct {
		Env  string `default:"dev" envconfig:"ENV"`
		HTTP struct {
			Host         string        `default:"0.0.0.0:1323" envconfig:"HOST"`
			ReadTimeout  time.Duration `default:"10s" envconfig:"READ_TIMEOUT"`
			WriteTimeout time.Duration `default:"10s" envconfig:"WRITE_TIMEOUT"`
		}
		App struct {
			ShutdownTimeout time.Duration `default:"5s" envconfig:"SHUTDOWN_TIMEOUT"`
		}
	}

	// For additional details refer to https://github.com/kelseyhightower/envconfig
	if err := envconfig.Process(service, &cfg); err != nil {
		log.Fatalf("main : Parsing Config : %v", err)
	}

	if err := flag.Process(&cfg); err != nil {
		if err != flag.ErrHelp {
			log.Fatalf("main : Parsing Command Line : %v", err)
		}
		return // We displayed help.
	}

	// =========================================================================
	// Start API Service

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:           cfg.HTTP.Host,
		Handler:        API(shutdown, log),
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Printf("main : API Listening %s", cfg.HTTP.Host)
		serverErrors <- api.ListenAndServe()
	}()

	// =========================================================================
	// Shutdown

	// Blocking main and waiting for shutdown.
	select {
	case err := <-serverErrors:
		log.Fatalf("main : Error starting server: %v", err)

	case sig := <-shutdown:
		log.Printf("main : %v : Start shutdown..", sig)

		// Create context for Shutdown call.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.App.ShutdownTimeout)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("main : Graceful shutdown did not complete in %v : %v", cfg.App.ShutdownTimeout, err)
			err = api.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			log.Fatal("main : Integrity issue caused shutdown")
		case err != nil:
			log.Fatalf("main : Could not stop server gracefully : %v", err)
		}
	}
}

// API returns a handler for a set of routes.
func API(shutdown chan os.Signal, log *log.Logger) http.Handler {

	// Construct the web.App which holds all routes as well as common Middleware.
	app := web.NewApp(shutdown, log, webcontext.Env_Dev, mid.Logger(log))

	app.Handle("GET", "/swagger/", saasSwagger.WrapHandler)
	app.Handle("GET", "/swagger/*", saasSwagger.WrapHandler)

	/*
		Or can use SaasWrapHandler func with configurations.
		url := saasSwagger.URL("http://localhost:1323/swagger/doc.json") //The url pointing to API definition
		e.GET("/swagger/*", saasSwagger.SaasWrapHandler(url))
	*/

	return app
}
