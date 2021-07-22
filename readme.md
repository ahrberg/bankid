# Swedish BankID API client for Go

The client is tested but not battle tested. For detailed API documentation go to [https://www.bankid.com](https://www.bankid.com).

## Install

No versions supported at the moment only master targeting BankID API version 5.1.

Install

```sh
go get github.com/ahrberg/bankid
```

## Usage

```go
import "github.com/ahrberg/bankid"

// Client configured with certificates
// use bankid.NewClient to specify a custom client
client := bankid.NewCertClient("./ca_cert.pem", "./client_cert.pem", "./client_key.pem")

// Make an auth request
ctx = context.Background()

params := AuthRequest{
		EndUserIp: "192.168.0.1",
    }

res, err := client.Auth(ctx, &params)
```

## Errors

Response errors are returned using the following `struct`. For error details see BankID [error documentation](https://www.bankid.com/utvecklare/guider/teknisk-integrationsguide/graenssnittsbeskrivning/felfall).

```go
type ErrorResponse struct {
	// HTTP response that caused this error
	Response  *http.Response
	ErrorCode string `json:"errorCode"`
	Details   string `json:"details"`
}
```
