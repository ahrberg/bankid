package bankid

import (
	"context"
	"encoding/base64"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestBaseUrl   = "https://appapi2.test.bankid.com/rp/v5.1"
	CaCert        = "./test_cert/test_ca_cert.pem"
	ClientCert    = "./test_cert/client_cert.pem"
	ClientCertKey = "./test_cert/client_key.pem"
)

var (
	c   *Client
	ctx = context.Background()
)

func TestError(t *testing.T) {

	setup()

	params := AuthRequest{
		// Not setting EndUserIp to get error
	}

	_, err := c.Auth(ctx, &params)

	aErr, _ := err.(*ErrorResponse)

	assert.Equal(t, aErr.ErrorCode, "invalidParameters")
	assert.Equal(t, aErr.Details, "Invalid endUserIp")
}

func TestAuth(t *testing.T) {
	setup()

	params := AuthRequest{
		EndUserIp: "192.168.0.1",
	}

	res, err := c.Auth(ctx, &params)

	assert.Nil(t, err, "Expect no error")
	assert.NotEmpty(t, res.OrderRef, "Expect OrderRef")
	assert.NotEmpty(t, res.AutoStartToken, "Expect AutoStartToken")
	assert.NotEmpty(t, res.QrStartSecret, "Expect QrStartSecret")
	assert.NotEmpty(t, res.QrStartToken, "Expect QrStartToken")
}

func TestSign(t *testing.T) {
	setup()

	params := SignRequest{
		EndUserIp:       "192.168.0.1",
		UserVisibleData: base64.StdEncoding.EncodeToString([]byte("Sign this please")),
	}

	res, err := c.Sign(ctx, &params)

	if err != nil {
		log.Fatal(err)
	}

	assert.Nil(t, err, "Expect no error")
	assert.NotEmpty(t, res.OrderRef, "Expect OrderRef")
	assert.NotEmpty(t, res.AutoStartToken, "Expect AutoStartToken")
	assert.NotEmpty(t, res.QrStartSecret, "Expect QrStartSecret")
	assert.NotEmpty(t, res.QrStartToken, "Expect QrStartToken")
}

func TestCollect(t *testing.T) {
	setup()

	params := SignRequest{
		EndUserIp:       "192.168.0.1",
		UserVisibleData: base64.StdEncoding.EncodeToString([]byte("Sign this please")),
	}

	signRes, err := c.Sign(ctx, &params)

	if err != nil {
		log.Fatal(err)
	}

	collectParams := CollectRequest{
		OrderRef: signRes.OrderRef,
	}

	collectRes, err := c.Collect(ctx, &collectParams)

	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, collectRes.Status, "pending")
	assert.Equal(t, collectRes.HintCode, "outstandingTransaction")
	assert.Equal(t, collectRes.CompletionData, "")
	assert.Equal(t, collectRes.OrderRef, signRes.OrderRef)
}

func TestCancel(t *testing.T) {
	setup()

	params := SignRequest{
		EndUserIp:       "192.168.0.1",
		UserVisibleData: base64.StdEncoding.EncodeToString([]byte("Sign this please")),
	}

	signRes, err := c.Sign(ctx, &params)

	if err != nil {
		log.Fatal(err)
	}

	cancelParams := CancelRequest{
		OrderRef: signRes.OrderRef,
	}

	cancelRes, err := c.Cancel(ctx, &cancelParams)

	assert.Nil(t, err, "Expect no error")
	assert.Equal(t, cancelRes, &CancelResponse{})
}

func setup() {

	var err error

	c, err = NewCertClient(CaCert, ClientCert, ClientCertKey)

	if err != nil {
		log.Fatal(err)
	}

	c.BaseUrl = TestBaseUrl
}
