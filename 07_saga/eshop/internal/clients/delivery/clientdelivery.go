package clientdelivery

import (
	"bytes"
	"encoding/json"
	"errors"
	reqDelivery "eshop/internal/request/delivery"
	"fmt"
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

func (c *Client) Reserve() (*reqDelivery.Response, error) {

	resp, err := c.hTTPClient.Post(c.uri+"/reserve_slot", clientType, bytes.NewReader([]byte{}))
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("Server response code status: %s", resp.Status))
	}

	respData, err := reqDelivery.FromBody(resp.Body)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func (c *Client) Rollback(slot int32) error {
	jsonData, err := json.MarshalIndent(reqDelivery.Rollback{Slot: slot}, "	", "")
	if err != nil {
		return err
	}

	resp, err := c.hTTPClient.Post(c.uri+"/rollback", clientType, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("rollback service deliver response: %s", resp.Status))
	}
	return nil
}
