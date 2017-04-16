package slack

import (
	"encoding/json"
)

type SlackMessage struct {
	Text        string       `json:"text,omitempty"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type Attachment struct {
	Fallback       string            `json:"fallback,omitempty"`
	Color          string            `json:"color,omitempty"`
	Pretext        string            `json:"pretext,omitempty"`
	AuthorName     string            `json:"author_name,omitempty"`
	AuthorLink     string            `json:"author_link,omitempty"`
	AuthorIcon     string            `json:"author_icon,omitempty"`
	Title          string            `json:"title,omitempty"`
	Text           string            `json:"text,omitempty"`
	ImageURL       string            `json:"image_url,omitempty"`
	ThumbURL       string            `json:"thumb_url,omitempty"`
	Footer         string            `json:"footer,omitempty"`
	FooterIcon     string            `json:"footer_icon,omitempty"`
	Ts             int64             `json:"ts,omitempty"`
	Fields         []AttachmentField `json:"fields,omitempty"`
	CallbackID     string            `json:"callback_id,omitempty"`
	AttachmentType string            `json:"attachment_type,omitempty"`
	Actions        []Action          `json:"actions,omitempty"`
}

type Action struct {
	Name  string `json:"name,omitempty"`
	Text  string `json:"text,omitempty"`
	Type  string `json:"type,omitempty"`
	Value string `json:"value,omitempty"`
}

type AttachmentField struct {
	Title string `json:"title,omitempty"`
	Value string `json:"value,omitempty"`
	Short string `json:"short,omitempty"`
}

func (m *SlackMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}

func (m *SlackMessage) AddAttachment(a Attachment) {
	m.Attachments = append(m.Attachments, a)
}

func (a *Attachment) AddField(f AttachmentField) {
	a.Fields = append(a.Fields, f)
}

func (a *Attachment) AddAction(action Action) {
	a.Actions = append(a.Actions, action)
}
