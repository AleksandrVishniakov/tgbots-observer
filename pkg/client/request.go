package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/AleksandrVishniakov/tgbots-util/ctxutil"
	"github.com/AleksandrVishniakov/tgbots-util/http/e"
)

func request[T any, E any](
	ctx context.Context,
	client *http.Client,
	method string,
	url string,
	data *T,
) (*E, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("marshal body: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")
	addRequestIDFromContext(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("make request: %w", err)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		if len(respBody) > 0 {
			var httpError e.HTTPError
			err = json.Unmarshal(respBody, &httpError)
			if err != nil {
				return nil, fmt.Errorf("unmarshal response: %w", err)
			}

			return nil, &httpError
		}

		return nil, fmt.Errorf("api error: status %d", resp.StatusCode)
	}

	if len(respBody) > 0 {
		var res E
		err = json.Unmarshal(respBody, &res)
		if err != nil {
			return nil, fmt.Errorf("unmarshal response: %w", err)
		}

		return &res, nil
	}

	return nil, nil
}

func addRequestIDFromContext(ctx context.Context, r *http.Request) {
	if rid, ok := ctx.Value(ctxutil.ContextKey_RequestID).(string); ok && len(rid)> 0 {
		r.Header.Add("X-Request-Id", rid)
	}
}
