package publisher

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Utsav-pixel/go-sensor-engine/internal/sim"
)

type HTTPPublisher struct {
	endpoint string
	client   *http.Client
}

func NewHTTPPublisher(endpoint string) *HTTPPublisher {
	return &HTTPPublisher{
		endpoint: endpoint,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (h *HTTPPublisher) Publish(ctx context.Context, telemetry sim.Telemetry) error {
	payload, err := json.Marshal(telemetry)

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP request failed with status: %d", resp.StatusCode)
	}

	return nil
}
