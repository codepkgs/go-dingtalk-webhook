package dingtalk

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	WebhookAddress string
	Secret         string
}

var ErrWebhookAddress = fmt.Errorf(`dingtalk webhook address must begin with "http://" or "https://"`)
var ErrWebhookSecret = fmt.Errorf(`dingtalk webhook secret must begin with "SEC"`)

func NewClient(webhookAddress, secret string) (*Client, error) {
	if !strings.HasPrefix(webhookAddress, "http://") && !strings.HasPrefix(webhookAddress, "https://") {
		return nil, ErrWebhookAddress
	}

	if secret != "" {
		if !strings.HasPrefix(secret, "SEC") {
			return nil, ErrWebhookSecret
		}
	}

	return &Client{
		WebhookAddress: webhookAddress,
		Secret:         secret,
	}, nil
}

// sign 计算签名 -- 在Secret不为空的情况下，需要计算签名
func (c *Client) sign() map[string]string {
	timestamp := time.Now().UnixMilli()
	contents := []byte(fmt.Sprintf("%s%s%s", strconv.FormatInt(timestamp, 10), "\n", c.Secret))
	h := hmac.New(sha256.New, []byte(c.Secret))
	h.Write(contents)
	hs := h.Sum(nil)
	r := make([]byte, base64.StdEncoding.EncodedLen(len(hs)))
	base64.StdEncoding.Encode(r, hs)
	sign := url.QueryEscape(string(r))
	return map[string]string{
		"timestamp": strconv.FormatInt(timestamp, 10),
		"sign":      sign,
	}
}
