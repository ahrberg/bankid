package bankid

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type SignRequest struct {
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
	// The text to be displayed and signed. String. The text can be formatted using
	// CR, LF and CRLF for new lines. The text must be encoded as UTF-8 and then
	// base 64 encoded. 1--40 000 characters after base 64 encoding.
	UserVisibleData string `json:"userVisibleData,omitempty"`
	// Data not displayed to the user. String. The value must be base 64-encoded.
	// 1-200 000 characters after base 64-encoding.
	UserNonVisibleData string `json:"userNonVisibleData,omitempty"`
	// If present, and set to “simpleMarkdownV1”, this parameter indicates that
	// userVisibleData holds formatting characters which, if used correctly,
	// will make the text displayed with the user nicer to look at.
	UserVisibleDataFormat string `json:"userVisibleDataFormat,omitempty"`
}

type SignResponse struct {
	OrderRef       string `json:"orderRef"`
	AutoStartToken string `json:"autoStartToken"`
	QrStartToken   string `json:"qrStartToken"`
	QrStartSecret  string `json:"qrStartSecret"`
}

func (c *Client) Sign(ctx context.Context, params *SignRequest) (*SignResponse, error) {
	url := fmt.Sprintf("%s/sign", c.BaseUrl)

	j, err := json.Marshal(params)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))

	if err != nil {
		return nil, err
	}

	res := SignResponse{}

	if err := c.sendReq(ctx, req, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
