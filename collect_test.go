package bankid

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	completed = `{
		"orderRef":"131daac9-16c6-4618-beb0-365768f37288", 
		"status":"complete",
		"completionData":{
		"user":{
			"personalNumber":"190000000000",
			"name":"Karl Karlsson",
			"givenName":"Karl",
			"surname":"Karlsson"
		},
		"device":{
			"ipAddress":"192.168.0.1"
		},
		"cert":{
			"notBefore":"1502983274000",
			"notAfter":"1563549674000"
		},
		"signature":"", 
		"ocspResponse":""
		} 
	}`
)

func TestCollectComplete(t *testing.T) {
	srv := serverMock()
	defer srv.Close()

	c, _ := NewClient(nil)

	c.BaseUrl = srv.URL

	res, _ := c.Collect(context.TODO(), &CollectRequest{
		OrderRef: "1234",
	})

	expected := &CollectResponse{}
	json.Unmarshal([]byte(completed), &expected)

	assert.Equal(t, res, expected)
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/collect", collectComplatedMock)

	srv := httptest.NewServer(handler)

	return srv
}

func collectComplatedMock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte(completed))
}
