package app

import (
	"context"
	"fmt"
	"io"
	"net/http"

	"github.com/AleksandrVishniakov/tgbots-observer/internal/configs"
	"github.com/AleksandrVishniakov/tgbots-observer/internal/handlers"
	apiv1 "github.com/AleksandrVishniakov/tgbots-observer/internal/handlers/v1"
	"github.com/AleksandrVishniakov/tgbots-observer/internal/usecases/mqadapters"
	"github.com/AleksandrVishniakov/tgbots-observer/internal/usecases/observer"
	tgbotstgapi "github.com/AleksandrVishniakov/tgbots-tgapi"
	"github.com/AleksandrVishniakov/tgbots-util/http/server"
	"github.com/AleksandrVishniakov/tgbots-util/logger"
	"github.com/nats-io/nats.go"
)

type App struct {
	cfg *configs.Configs
	stdout io.Writer
}

func New(cfg *configs.Configs, stdout io.Writer) *App {
	return &App{
		cfg: cfg,
		stdout: stdout,
	}
}

func (app *App) Run(ctx context.Context) error {
	const src = "App.Run"

	log := logger.New(app.stdout, app.cfg.Debug)

	nc, err := initNATSConnection(app.cfg.NATS.URL)
	if err != nil {
		return fmt.Errorf("%s connect NATS: %w", src, err)
	}
	defer nc.Drain()

	mq, err := mqadapters.NewNATSAdapter(ctx, log, nc)
	if err != nil {
		return fmt.Errorf("%s create NATS adapter: %w", src, err)
	}

	tgapi := tgbotstgapi.New()
	observers := observer.NewObserversManager(log, tgapi, mq)

	apiv1 := apiv1.New(observers)
	handler := handlers.New(log)
	srv := server.New(server.Configs{
		Port: app.cfg.HTTP.Port,
		Host: app.cfg.HTTP.Host,
	}, handler.InitRoutes(map[string]http.Handler{
		"/api/v1": apiv1.InitRoutes(),
	}))
	defer srv.Shutdown(ctx)

	log.Info(fmt.Sprintf("server started on http://%s:%d", app.cfg.HTTP.Host, app.cfg.HTTP.Port))

	return srv.Run()
}

func initNATSConnection(url string) (*nats.Conn, error) {
	if url == "" {
		url = nats.DefaultURL
	}

	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return nc, nil
}
