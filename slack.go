package slackWebhook

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/parnurzeal/gorequest"
)

type StatusColor string

const (
	ErrorColor StatusColor = "#a40000"
	OKColor    StatusColor = "#2aba88"
)

type Field struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type Action struct {
	Type  string `json:"type"`
	Text  string `json:"text"`
	URL   string `json:"url"`
	Style string `json:"style"`
}

type Attachment struct {
	Fallback     *string     `json:"fallback"`
	Color        StatusColor `json:"color"`
	PreText      *string     `json:"pretext"`
	AuthorName   *string     `json:"author_name"`
	AuthorLink   *string     `json:"author_link"`
	AuthorIcon   *string     `json:"author_icon"`
	Title        *string     `json:"title"`
	TitleLink    *string     `json:"title_link"`
	Text         *string     `json:"text"`
	ImageURL     *string     `json:"image_url"`
	Fields       []*Field    `json:"fields"`
	Footer       *string     `json:"footer"`
	FooterIcon   *string     `json:"footer_icon"`
	Timestamp    *int64      `json:"ts"`
	MarkdownIn   *[]string   `json:"mrkdwn_in"`
	Actions      []*Action   `json:"actions"`
	CallbackID   *string     `json:"callback_id"`
	ThumbnailURL *string     `json:"thumb_url"`
}

type Payload struct {
	Parse       string       `json:"parse,omitempty"`
	Username    string       `json:"username,omitempty"`
	IconURL     string       `json:"icon_url,omitempty"`
	IconEmoji   string       `json:"icon_emoji,omitempty"`
	Channel     string       `json:"channel,omitempty"`
	Text        string       `json:"text,omitempty"`
	LinkNames   string       `json:"link_names,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
	UnfurlLinks bool         `json:"unfurl_links,omitempty"`
	UnfurlMedia bool         `json:"unfurl_media,omitempty"`
	Markdown    bool         `json:"mrkdwn,omitempty"`
}

func (attachment *Attachment) AddField(field Field) *Attachment {
	attachment.Fields = append(attachment.Fields, &field)
	return attachment
}

func (attachment *Attachment) AddAction(action Action) *Attachment {
	attachment.Actions = append(attachment.Actions, &action)
	return attachment
}

func redirectPolicyFunc(req gorequest.Request, via []gorequest.Request) error {
	return fmt.Errorf("Incorrect token (redirection)")
}

func Send(webhookUrl string, payload Payload) error {

	client := &http.Client{}
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(payload)

	req, err := http.NewRequest("POST", webhookUrl, buf)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(body))

		return errors.New("400 not found")
	}

	return nil
}
