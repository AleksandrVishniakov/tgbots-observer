package v1

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/AleksandrVishniakov/tgbots-util/http/e"
	"github.com/AleksandrVishniakov/tgbots-util/http/json"
	"github.com/AleksandrVishniakov/tgbots-util/http/middlewares"
)

type Observer interface {
	Observe(
		ctx context.Context,
		id int64,
		token string,
		pollingInterval time.Duration,
	) error

	CloseObserver(id int64)
}

type API struct {
	observer Observer
}

func New(obsevrer Observer) *API {
	return &API{
		observer: obsevrer,
	}
}

func (api *API) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /observe/{id}", middlewares.Error(api.StartObserve))
	mux.Handle("POST /observe/{id}/stop", middlewares.Error(api.StopObserve))

	return mux
}

type StartObserveRequest struct {
	Token           string `json:"token"`
	PollingInterval time.Duration `json:"pollingInterval"`
}

func (api *API) StartObserve(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	dto, err := json.Decode[StartObserveRequest](r.Body)
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	api.observer.Observe(
		context.WithoutCancel(r.Context()),
		int64(id),
		dto.Token,
		dto.PollingInterval * time.Second,
	)
	return nil
}

func (api *API) StopObserve(w http.ResponseWriter, r *http.Request) error {
	defer r.Body.Close()

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		return e.BadRequest(e.WithError(err))
	}

	api.observer.CloseObserver(int64(id))
	return nil
}
