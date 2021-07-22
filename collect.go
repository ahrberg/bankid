package bankid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type CollectRequest struct {
	// The orderRef returned from auth or sign.
	OrderRef string `json:"orderRef"`
}

type CollectResponse struct {
	OrderRef       string `json:"orderRef"`
	Status         string `json:"status"`
	HintCode       string `json:"hintCode"`
	CompletionData string `json:"completionData"`
}

func (c *Client) Collect(ctx context.Context, params *CollectRequest) (*CollectResponse, error) {
	url := fmt.Sprintf("%s/collect", c.BaseUrl)

	j, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))

	if err != nil {
		return nil, err
	}

	res := CollectResponse{}

	if err := c.sendReq(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
