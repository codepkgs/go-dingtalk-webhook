package dingtalk

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

func getReturn(bytes []byte) (*SendResult, error) {
	var r SendResult
	err := json.Unmarshal(bytes, &r)
	if err != nil {
		return nil, err
	} else {
		return &r, nil
	}
}

type SendResult struct {
	ErrCode int64  `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

type BtnOrientation int

const (
	Vertical   BtnOrientation = iota // 按钮垂直排列
	Horizontal                       // 按钮水平排列
)

type ActionCardButton struct {
	Title      string     `json:"title"`      // 单条信息文本
	ActionURL  string     `json:"actionURL"`  // 点击单条信息到跳转链接
	ActionType ActionType `json:"actionType"` // 打开方式。在应用内侧边栏打开还是在独立的浏览器中打开
}

type FeedCardLink struct {
	Title      string     `json:"title"`      // 单条信息文本
	MessageURL string     `json:"messageURL"` // 点击单条信息到跳转链接
	PicURL     string     `json:"picURL"`     // 单条信息后面图片的URL
	ActionType ActionType `json:"actionType"` // 打开方式。在应用内侧边栏打开还是在独立的浏览器中打开
}

type ActionType int

const (
	APP ActionType = iota // 在APP侧边栏打开
	WEB                   // 在独立的浏览器中打开
)

// Text 发送普通文本消息
func (c *Client) Text(content string, atMobiles []string, isAtAll bool) (*SendResult, error) {
	t := struct {
		At struct {
			AtMobiles []string `json:"atMobiles,omitempty"`
			AtUserIds []string `json:"atUserIds,omitempty"`
			IsAtAll   bool     `json:"isAtAll,omitempty"`
		} `json:"at"`
		Text struct {
			Content string `json:"content"`
		} `json:"text"`
		Msgtype string `json:"msgtype"`
	}{
		Msgtype: "text",
		Text: struct {
			Content string `json:"content"`
		}{Content: content},
		At: struct {
			AtMobiles []string `json:"atMobiles,omitempty"`
			AtUserIds []string `json:"atUserIds,omitempty"`
			IsAtAll   bool     `json:"isAtAll,omitempty"`
		}{AtMobiles: atMobiles, AtUserIds: nil, IsAtAll: isAtAll},
	}

	body, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(body)
	if err != nil {
		return nil, err
	}

	return getReturn(resp)
}

// Markdown 发送Markdown消息
func (c *Client) Markdown(title, content string, atMobiles []string, isAtAll bool) (*SendResult, error) {
	t := struct {
		Msgtype  string `json:"msgtype"`
		Markdown struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		} `json:"markdown"`
		At struct {
			AtMobiles []string `json:"atMobiles"`
			AtUserIds []string `json:"atUserIds"`
			IsAtAll   bool     `json:"isAtAll"`
		} `json:"at"`
	}{
		Msgtype: "markdown",
		At: struct {
			AtMobiles []string `json:"atMobiles"`
			AtUserIds []string `json:"atUserIds"`
			IsAtAll   bool     `json:"isAtAll"`
		}{AtMobiles: atMobiles, AtUserIds: nil, IsAtAll: isAtAll},
		Markdown: struct {
			Title string `json:"title"`
			Text  string `json:"text"`
		}{Title: title, Text: content},
	}

	body, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(body)
	if err != nil {
		return nil, err
	}

	var r SendResult
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	} else {
		return &r, nil
	}
}

// Link 发送Link类型的消息
func (c *Client) Link(title, content, messageUrl, picUrl string) (*SendResult, error) {
	t := struct {
		Msgtype string `json:"msgtype"`
		Link    struct {
			Text       string `json:"text"`
			Title      string `json:"title"`
			PicUrl     string `json:"picUrl"`
			MessageUrl string `json:"messageUrl"`
		} `json:"link"`
	}{
		Msgtype: "link",
		Link: struct {
			Text       string `json:"text"`
			Title      string `json:"title"`
			PicUrl     string `json:"picUrl"`
			MessageUrl string `json:"messageUrl"`
		}{Text: content, Title: title, PicUrl: picUrl, MessageUrl: messageUrl},
	}

	body, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(body)
	if err != nil {
		return nil, err
	}

	var r SendResult
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return nil, err
	} else {
		return &r, nil
	}
}

// ActionCard 发送独立跳转的ActionCard消息
func (c *Client) ActionCard(title, content string, btnOrientation BtnOrientation, btns []ActionCardButton) (*SendResult, error) {
	webFormat := "dingtalk://dingtalkclient/page/link?url=%s&pc_slide=false"
	var bs []struct {
		Title     string `json:"title"`
		ActionURL string `json:"actionURL"`
	}

	for _, b := range btns {
		i := struct {
			Title     string `json:"title"`
			ActionURL string `json:"actionURL"`
		}{
			Title: b.Title,
		}
		if b.ActionType == 0 {
			i.ActionURL = b.ActionURL
		} else {
			i.ActionURL = fmt.Sprintf(webFormat, url.QueryEscape(b.ActionURL))
		}
		bs = append(bs, i)
	}

	t := struct {
		Msgtype    string `json:"msgtype"`
		ActionCard struct {
			Title          string `json:"title"`
			Text           string `json:"text"`
			BtnOrientation string `json:"btnOrientation"`
			Btns           []struct {
				Title     string `json:"title"`
				ActionURL string `json:"actionURL"`
			} `json:"btns"`
		} `json:"actionCard"`
	}{
		Msgtype: "actionCard",
		ActionCard: struct {
			Title          string `json:"title"`
			Text           string `json:"text"`
			BtnOrientation string `json:"btnOrientation"`
			Btns           []struct {
				Title     string `json:"title"`
				ActionURL string `json:"actionURL"`
			} `json:"btns"`
		}{Title: title, Text: content, BtnOrientation: strconv.Itoa(int(btnOrientation)), Btns: bs},
	}

	body, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(body)
	if err != nil {
		return nil, err
	}

	return getReturn(resp)
}

// FeedCard 发送FeedCard类型的消息
func (c *Client) FeedCard(links []FeedCardLink) (*SendResult, error) {
	webFormat := "dingtalk://dingtalkclient/page/link?url=%s&pc_slide=false"

	var nlinks []struct {
		Title      string `json:"title"`
		MessageURL string `json:"messageURL"`
		PicURL     string `json:"picURL"`
	}

	for _, link := range links {
		i := struct {
			Title      string `json:"title"`
			MessageURL string `json:"messageURL"`
			PicURL     string `json:"picURL"`
		}{
			Title:  link.Title,
			PicURL: link.PicURL,
		}
		if link.ActionType == 0 {
			i.MessageURL = link.MessageURL
		} else {
			i.MessageURL = fmt.Sprintf(webFormat, url.QueryEscape(link.MessageURL))
		}
		nlinks = append(nlinks, i)
	}

	t := struct {
		Msgtype  string `json:"msgtype"`
		FeedCard struct {
			Links []struct {
				Title      string `json:"title"`
				MessageURL string `json:"messageURL"`
				PicURL     string `json:"picURL"`
			} `json:"links"`
		} `json:"feedCard"`
	}{
		Msgtype: "feedCard",
		FeedCard: struct {
			Links []struct {
				Title      string `json:"title"`
				MessageURL string `json:"messageURL"`
				PicURL     string `json:"picURL"`
			} `json:"links"`
		}{Links: nlinks},
	}

	body, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}
	resp, err := c.do(body)
	if err != nil {
		return nil, err
	}

	return getReturn(resp)
}
