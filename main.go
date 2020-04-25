package main

import (
	"context"
	"github.com/ice1n36/kurapika/clients"
	"github.com/ice1n36/kurapika/handlers"
	"go.uber.org/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

func NewMux(lc fx.Lifecycle, logger *zap.SugaredLogger) *http.ServeMux {
	logger.Infow("Executing NewMux.")
	// First, we construct the mux and server. We don't want to start the server
	// until all handlers are registered.
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8081",
		Handler: mux,
	}
	// If NewMux is called, we know that another function is using the mux. In
	// that case, we'll use the Lifecycle type to register a Hook that starts
	// and stops our HTTP server.
	//
	// Hooks are executed in dependency order. At startup, NewLogger's hooks run
	// before NewMux's. On shutdown, the order is reversed.
	//
	// Returning an error from OnStart hooks interrupts application startup. Fx
	// immediately runs the OnStop portions of any successfully-executed OnStart
	// hooks (so that types which started cleanly can also shut down cleanly),
	// then exits.
	//
	// Returning an error from OnStop hooks logs a warning, but Fx continues to
	// run the remaining hooks.
	lc.Append(fx.Hook{
		// To mitigate the impact of deadlocks in application startup and
		// shutdown, Fx imposes a time limit on OnStart and OnStop hooks. By
		// default, hooks have a total of 30 seconds to complete. Timeouts are
		// passed via Go's usual context.Context.
		OnStart: func(context.Context) error {
			logger.Infow("Starting HTTP server.")
			// In production, we'd want to separate the Listen and Serve phases for
			// better error-handling.
			go server.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Infow("Stopping HTTP server.")
			return server.Shutdown(ctx)
		},
	})

	return mux
}

func NewHandler(logger *zap.SugaredLogger) (http.Handler, error) {
	logger.Infow("Executing NewHandler.")
	return http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		logger.Infow("Got a request.")
	}), nil
}

func NewLogger() *zap.SugaredLogger {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	sugar.Infow("Executing NewLogger.")
	return sugar
}

func NewConfig() config.Provider {
	configDir := os.Getenv("CONFIGDIR")

	if configDir == "" {
		cwd, _ := os.Getwd()
		configDir = filepath.Join(cwd, "config")
	}

	basePath := filepath.Join(configDir, "base.yaml")
	secretsPath := filepath.Join(configDir, "secrets.yaml")
	provider, err := config.NewYAML(config.File(basePath), config.File(secretsPath))
	if err != nil {
		// configs are important, if this fails, the whole app should fail
		panic(err)
	}

	return provider
}

func Register(mux *http.ServeMux,
	h http.Handler,
	sah *handlers.NewAppHandler) {
	mux.Handle("/", h)
	mux.Handle("/new_app", sah)
}

func main() {
	fx.New(opts()).Run()
}

func opts() fx.Option {
	return fx.Options(
		fx.Provide(
			NewMux,
			NewHandler,
			handlers.NewNewAppHandler,
			NewLogger,
			clients.NewMobSFHTTPClient,
			NewConfig,
		),
		fx.Invoke(Register),
	)
}
