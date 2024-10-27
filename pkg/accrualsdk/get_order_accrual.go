package accrualsdk

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *AccrualClient) GetOrderAccrual(ctx context.Context, orderNumber string) (*OrderAccrualInfo, error) {
	path := fmt.Sprintf("orders/%s", orderNumber)
	var queries map[string]string

	content, statusCode, err := c.makeLimitedGetRequest(ctx, path, queries)
	if err != nil {
		return nil, fmt.Errorf("AccrualClient.GetOrderAccrual: %w", err)
	}
	if statusCode == http.StatusNoContent {
		return nil, nil
	}

	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("AccrualClient.GetOrderAccrual unexpected status code %d", statusCode)
	}

	info := &OrderAccrualInfo{}
	if err := json.Unmarshal(content, &info); err != nil {
		return nil, fmt.Errorf("AccrualClient.GetOrderAccrual json.Unmarshal (data): %w", err)
	}

	return info, nil
}
