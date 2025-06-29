package observer

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/AleksandrVishniakov/tgbots-tgapi/models"
	"github.com/AleksandrVishniakov/tgbots-util/logger"
)

type MQ interface {
	PublishUpdate(
		ctx context.Context,
		subjectID int64,
		update *models.Update,
	) error
}

type ObserversManager struct {
	log             *slog.Logger
	mu              sync.Mutex
	observers       map[int64]*observer
	updatesProvider UpdatesProvider
	mq              MQ
}

func NewObserversManager(
	log *slog.Logger,
	updatesProvider UpdatesProvider,
	mq MQ,
) *ObserversManager {
	return &ObserversManager{
		log:             log,
		mu:              sync.Mutex{},
		updatesProvider: updatesProvider,
		observers:       make(map[int64]*observer),
		mq:              mq,
	}
}

func (w *ObserversManager) Observe(
	ctx context.Context,
	id int64,
	token string,
	pollingInterval time.Duration,
) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if _, ok := w.observers[id]; ok {
		return nil
	}

	observer := newObserver(w.log, w.updatesProvider, id, token)
	observer.UpdateHandler = func(ctx context.Context, u *models.Update) {
		err := w.mq.PublishUpdate(ctx, id, u)
		if err != nil {
			w.log.ErrorContext(ctx, "update publishing", logger.Err(err))
		}
	}
	w.observers[id] = observer

	observer.Run(ctx, pollingInterval)

	return nil
}

func (w *ObserversManager) CloseObserver(id int64) {
	w.mu.Lock()
	defer w.mu.Lock()

	if observer, ok := w.observers[id]; ok {
		observer.Close()
		delete(w.observers, id)
	}
}

func (w *ObserversManager) Close() {
	for id := range w.observers {
		w.CloseObserver(id)
	}
}
