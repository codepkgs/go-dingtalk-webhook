package dingtalk

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
	"time"
)

const (
	RequestTimeout = 5000 //请求超时时间
)

// 发送HTTP请求
func (c *Client) do(bytes []byte) ([]byte, error) {
	r := resty.New().
		SetTimeout(time.Duration(RequestTimeout)*time.Millisecond).
		R().
		SetHeader("Content-Type", "application/json")

	r.SetBody(bytes)

	if c.Secret != "" {
		sign := c.sign()
		r.SetQueryParams(sign)
	}

	resp, err := r.Post(c.WebhookAddress)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", resp.StatusCode(), resp.Body())
	}

	return resp.Body(), nil
}
