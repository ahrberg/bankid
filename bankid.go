package bankid

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	libraryVersion = "5.1"
	defaultBaseURL = "https://appapi2.bankid.com/rp/v5.1"
	userAgent      = "gobankid/" + libraryVersion
	mediaType      = "application/json"
)

type Client struct {
	BaseUrl    string
	UserAgent  string
	HTTPClient *http.Client
}

// NewCertClient returns a BankID API client with given certificates
func NewCertClient(caCertFilename string, clientCertFilename string, clientCertKeyFilename string) (*Client, error) {
	t, err := NewTransport(caCertFilename, clientCertFilename, clientCertKeyFilename)

	if err != nil {
		return nil, err
	}

	c := &http.Client{
		Timeout:   time.Minute,
		Transport: t,
	}

	return NewClient(c)
}

// NewClient returns a BankID API client using the given http.Client.
//
// Use this function if you want to specify your own client, else use NewCertClient.
func NewClient(httpClient *http.Client) (*Client, error) {

	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		BaseUrl:    defaultBaseURL,
		UserAgent:  userAgent,
		HTTPClient: httpClient,
	}, nil
}

type ErrorResponse struct {
	// HTTP response that caused this error
	Response  *http.Response
	ErrorCode string `json:"errorCode"`
	Details   string `json:"details"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %v %v",
		r.Response.Request.Method, r.Response.Request.URL, r.Response.StatusCode, r.Details, r.ErrorCode)
}

func (c *Client) sendReq(ctx context.Context, req *http.Request, v interface{}) error {
	req = req.WithContext(ctx)

	req.Header.Set("Content-Type", mediaType)
	req.Header.Set("Accept", mediaType)
	req.Header.Set("User-Agent", c.UserAgent)

	res, err := c.HTTPClient.Do(req)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = checkResponse(res)

	if err != nil {
		return err
	}

	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}

	return nil
}

func checkResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errRes := &ErrorResponse{Response: r}

	if err := json.NewDecoder(r.Body).Decode(&errRes); err != nil {
		errRes.Details = "could not parse json from response"
	}

	return errRes
}

// NewTransport creates a http.Transport custom CA and mTLS certificate. Only the specified CA will be used.
func NewTransport(caCertFilename string, clientCertFilename string, clientCertKeyFilename string) (*http.Transport, error) {
	rootCAs := x509.NewCertPool()

	caCert, err := ioutil.ReadFile(caCertFilename)
	if err != nil {
		return nil, err
	}

	rootCAs.AppendCertsFromPEM(caCert)

	clientCert, err := tls.LoadX509KeyPair(clientCertFilename, clientCertKeyFilename)
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		RootCAs:      rootCAs,
		Certificates: []tls.Certificate{clientCert},
	}

	tr := &http.Transport{TLSClientConfig: config}

	return tr, nil
}
