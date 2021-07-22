package bankid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type AuthRequest struct {
	// The user IP address as seen by RP. String. IPv4 and IPv6 is allowed.
	// Note the importance of using the correct IP address.
	// It must be the IP address representing the user agent (the end user device)
	// as seen by the RP. If there is a proxy for inbound traffic, special considerations
	// may need to be taken to get the correct address.
	// In some use cases the IP address is not available,
	// for instance for voice based services. In this case,
	// the internal representation of those systems IP address is ok to use.
	EndUserIp string `json:"endUserIp"`
	// The personal number of the user. String. 12 digits.
	// Century must be included. If the personal number is excluded,
	// the client must be started with the autoStartToken returned in the response.
	PersonalNumber string `json:"personalNumber,omitempty"`
	// Requirements on how the auth or sign order must be performed.
	Requirement string `json:"requirement,omitempty"`
}

type AuthResponse struct {
	OrderRef       string `json:"orderRef"`
	AutoStartToken string `json:"autoStartToken"`
	QrStartToken   string `json:"qrStartToken"`
	QrStartSecret  string `json:"qrStartSecret"`
}

func (c *Client) Auth(ctx context.Context, params *AuthRequest) (*AuthResponse, error) {
	url := fmt.Sprintf("%s/auth", c.BaseUrl)

	j, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))

	if err != nil {
		return nil, err
	}

	res := AuthResponse{}

	if err := c.sendReq(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
