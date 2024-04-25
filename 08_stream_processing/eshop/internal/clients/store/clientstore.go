package store

import (
	"bytes"
	"encoding/json"
	"errors"
	dbOrder "eshop/internal/db/order"
	response "eshop/internal/request"
	reqStore "eshop/internal/request/store"
	"fmt"
	"io"
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

func (c *Client) Reserve(items []dbOrder.Item) (bool, string, error) {
	itemsRequest := []reqStore.Item{}
	for _, item := range items {
		itemsRequest = append(itemsRequest, reqStore.Item{
			Title: item.Title,
			Count: item.Count,
		})
	}
	jsonData, err := json.MarshalIndent(reqStore.Reserve{Items: itemsRequest}, "	", "")
	if err != nil {
		return false, "", err
	}

	resp, err := c.hTTPClient.Post(c.uri+"/reserve", clientType, bytes.NewReader(jsonData))
	if err != nil {
		return false, "", err
	}

	if resp.StatusCode != http.StatusOK {
		return false, resp.Status, nil
	}

	common := &response.Common{}
	all, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, "", err
	}

	if err := json.Unmarshal(all, &common); err != nil {
		return false, "", err
	}

	return common.Status, common.Message, nil
}

func (c *Client) Rollback(items []dbOrder.Item) error {
	itemsRequest := []reqStore.Item{}
	for _, item := range items {
		itemsRequest = append(itemsRequest, reqStore.Item{
			Title: item.Title,
			Count: item.Count,
		})
	}
	jsonData, err := json.MarshalIndent(reqStore.Reserve{Items: itemsRequest}, "	", "")
	if err != nil {
		return err
	}

	resp, err := c.hTTPClient.Post(c.uri+"/rollback", clientType, bytes.NewReader(jsonData))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("response status: %s", resp.Status))
	}
	return nil

}
