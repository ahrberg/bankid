package bankid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CancelRequest struct {
	// The orderRef returned from auth or sign.
	OrderRef string `json:"orderRef"`
}

type CancelResponse struct {
}

func (c *Client) Cancel(ctx context.Context, params *CancelRequest) (*CancelResponse, error) {
	url := fmt.Sprintf("%s/cancel", c.BaseUrl)

	j, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))

	if err != nil {
		return nil, err
	}

	res := CancelResponse{}

	if err := c.sendReq(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
