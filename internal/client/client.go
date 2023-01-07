package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/SemenRyzhkov/go-market-app/internal/entity"
)

const (
	getAccrualPath = "api/orders"
)

type (
	Client struct {
		host       string
		httpClient *http.Client
	}
)

func NewClient(host string, timeout time.Duration) *Client {
	client := &http.Client{
		Timeout: timeout,
	}
	return &Client{
		host:       host,
		httpClient: client,
	}
}

func (c *Client) do(method, endpoint string, number int) (*http.Response, error) {
	baseURL := fmt.Sprintf("%s/%s/%d", c.host, endpoint, number)

	req, err := http.NewRequest(method, baseURL, nil)
	if err != nil {
		return nil, err
	}
	return c.httpClient.Do(req)
}

func (c *Client) GetOrderResponse(number int) (entity.OrderResponse, error) {
	res, err := c.do(http.MethodGet, getAccrualPath, number)
	if err != nil {
		return entity.OrderResponse{}, err
	}

	defer res.Body.Close()
	var orderResponse entity.OrderResponse
	err = json.NewDecoder(res.Body).Decode(&orderResponse)
	if err != nil {
		return entity.OrderResponse{}, err
	}
	return orderResponse, nil
}
