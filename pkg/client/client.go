package client

import (
	"context"
	"fmt"
	"net/http"

	"github.com/AleksandrVishniakov/tgbots-observer/pkg/dto"
)

type Client struct {
	httpClient *http.Client
	addr string
}

func NewClient(addr string) *Client {
	return &Client{
		addr: addr,
		httpClient: &http.Client{},
	}
}

func (cl *Client) Ping(ctx context.Context) error {
	const src = "Client.Ping"

	url := fmt.Sprintf("%s/ping", cl.addr)

	_, err := request[any, any](ctx, cl.httpClient, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", src, err)
	}

	return nil
}

type StartObserveParams struct {
	BotID int64
	Body *dto.StartObserveRequest
}

func (cl *Client) StartObserve(ctx context.Context, params *StartObserveParams) error {
	const src = "Client.StartObserve"

	url := fmt.Sprintf("%s/api/v1/observe/%d", cl.addr, params.BotID)

	_, err := request[dto.StartObserveRequest, any](ctx, cl.httpClient, http.MethodPost, url, params.Body)
	if err != nil {
		return fmt.Errorf("%s: %w", src, err)
	}

	return nil
}

type StopObserveParams struct {
	BotID int64
}

func (cl *Client) StopObserve(ctx context.Context, params *StopObserveParams) error {
	const src = "Client.StopObserve"

	url := fmt.Sprintf("%s/api/v1/observe/%d/stop", cl.addr, params.BotID)

	_, err := request[any, any](ctx, cl.httpClient, http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("%s: %w", src, err)
	}

	return nil
}
