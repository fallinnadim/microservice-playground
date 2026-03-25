package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/fallinnadim/order-service/internal/adapter/inbound/http/response"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
)

type PaymentRequest struct {
	OrderID string
	UserID  string
	Amount  int
}

type PaymentResponse struct {
	Status  string
	Message string
}

type PaymentAdapter struct {
	baseURL string
	client  *http.Client
}

func NewPaymentAdapter(baseURL string, client *http.Client) *PaymentAdapter {
	return &PaymentAdapter{
		baseURL, client,
	}
}

func (p *PaymentAdapter) Pay(ctx context.Context, req PaymentRequest) (*PaymentResponse, error) {
	url := p.baseURL + "/api/v1/payment"

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if deadline, ok := ctx.Deadline(); ok {
		remaining := time.Until(deadline)

		httpReq.Header.Set("X-Timeout-Remaining",
			strconv.Itoa(int(remaining.Milliseconds())),
		)
	}
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(httpReq.Header))
	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("payment service error: %s", string(respBody))
	}

	var respz response.SuccessResponse[PaymentResponse]
	if err := json.Unmarshal(respBody, &respz); err != nil {
		return nil, err
	}
	return &respz.Data, nil
}
