package dto

import "time"

type StartObserveRequest struct {
	Token           string        `json:"token"`
	PollingInterval time.Duration `json:"pollingInterval"`
}
