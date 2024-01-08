package akumu

import (
	"io"
	"log/slog"
	netHttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/studiolambda/akumu/event"
)

type Config struct {
	address string
	events  *event.System
	logger  *slog.Logger
}

type Akumu struct {
	Config
	server netHttp.Server
}

type Builder func(*Config)

const (
	EventStart event.Event = "akumu.start"
	EventStop  event.Event = "akumu.stop"
)

func WithLogger(logger *slog.Logger) Builder {
	return func(config *Config) {
		config.logger = logger
	}
}

func WithAddress(address string) Builder {
	return func(config *Config) {
		config.address = address
	}
}

func WithEvents(events *event.System) Builder {
	return func(config *Config) {
		config.events = events
	}
}

func New(builders ...Builder) *Akumu {
	config := Config{
		address: ":3000",
		logger: slog.New(
			slog.NewTextHandler(io.Discard, nil),
		),
	}

	for _, builder := range builders {
		builder(&config)
	}

	return &Akumu{
		Config: config,
		server: netHttp.Server{
			Addr:                         config.address,
			Handler:                      chi.NewRouter(),
			DisableGeneralOptionsHandler: false,
			ReadTimeout:                  0,
			ReadHeaderTimeout:            0,
			WriteTimeout:                 0,
			IdleTimeout:                  0,
			MaxHeaderBytes:               0,
		},
	}
}

func (akumu *Akumu) Start() {
	akumu.logger.Info("starting server", "address", akumu.address)

	if akumu.events != nil {
		akumu.events.Emit(EventStart, nil)
	}

	if err := akumu.server.ListenAndServe(); err != nil && err != netHttp.ErrServerClosed {
		akumu.logger.Error("error while starting server", "error", err)
	}
}

func (akumu *Akumu) Stop() {
	akumu.logger.Info("stopping server")

	if akumu.events != nil {
		akumu.events.Emit(EventStop, nil)
	}

	if err := akumu.server.Close(); err != nil {
		akumu.logger.Error("error while stopping server", "error", err)
	}
}
