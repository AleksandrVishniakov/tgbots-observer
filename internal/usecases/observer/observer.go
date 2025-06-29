package observer

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/AleksandrVishniakov/tgbots-tgapi/dto"
	"github.com/AleksandrVishniakov/tgbots-tgapi/models"
	"github.com/AleksandrVishniakov/tgbots-util/logger"
)

type UpdatesProvider interface {
	GetUpdates(
		ctx context.Context,
		token string,
		req *dto.GetUpdatesRequest,
	) (*dto.GetUpdatesResponse, error)
}

type observer struct {
	log             *slog.Logger
	updatesProvider UpdatesProvider
	id              int64
	token           string
	done            chan struct{}
	UpdateHandler   func(context.Context, *models.Update)
}

func newObserver(
	log *slog.Logger,
	updatesProvider UpdatesProvider,
	id int64,
	token string,
) *observer {
	return &observer{
		log:             log.With(slog.Int64("worker_id", id)),
		updatesProvider: updatesProvider,
		id:              id,
		token:           token,
		done:            make(chan struct{}),
		UpdateHandler:   func(ctx context.Context, u *models.Update) {},
	}
}

func (w *observer) Run(
	ctx context.Context,
	pollingInterval time.Duration,
) {
	const src = "Worker.Run"
	log := w.log.With("src", src)

	go func() {
		log.DebugContext(ctx, "new worker")

		ticker := time.NewTicker(pollingInterval)
		defer ticker.Stop()

		req := &dto.GetUpdatesRequest{
			Offset:  0,
			Timeout: int(pollingInterval.Seconds()),
		}

		for {
			select {
			case <-w.done:
				return
			case <-ticker.C:
				reqCtx, cancel := context.WithTimeout(ctx, pollingInterval)
				resp, err := w.updatesProvider.GetUpdates(reqCtx, w.token, req)
				cancel()
				if err != nil {
					if errors.Is(err, context.DeadlineExceeded) {
						log.WarnContext(ctx, "getUpdates: deadline")
					} else {
						log.ErrorContext(ctx, "getUpdates", logger.Err(err))
					}
					continue
				}

				if len(resp.Result) > 0 {
					req.Offset = resp.Result[len(resp.Result)-1].UpdateID + 1

					for _, u := range resp.Result {
						w.UpdateHandler(ctx, &u)
					}
				}
			}
		}
	}()
}

func (w *observer) Close() {
	close(w.done)
}
