package clientpayment

import (
	"bytes"
	"encoding/json"
	"errors"
	response "eshop/internal/request"
	reqpay "eshop/internal/request/pay"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

const clientType = "application/json"

type Client struct {
	uri        string
	hTTPClient *http.Client
}

func NewClient(uri string) *Client {
	return &Client{uri: uri, hTTPClient: &http.Client{
		Timeout: time.Second * 5,
	}}
}

func (c *Client) Pay() (*response.Common, error) {

	pay := reqpay.Pay{Sum: 999}
	payData, err := json.MarshalIndent(pay, "	", "")
	if err != nil {
		slog.Error("pay data marshal", "err", err)
	}

	resp, err := c.hTTPClient.Post(c.uri+"/pay", clientType, bytes.NewReader(payData))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Server response code status: %s", resp.Status))
	}

	common, err := response.FromBody(resp.Body)
	if err != nil {
		return nil, err
	}
	return common, nil

}
