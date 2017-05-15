package slack

import (
	"encoding/json"
	"fmt"
	"github.com/satori/go.uuid"
	"markusreschke.name/selfhostedsimplepolling/poll"
	"strconv"
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
	Ts             float64           `json:"ts,omitempty,string"`
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

type Team struct {
	ID     string `json:"id,omitempty"`
	Domain string `json:"domain,omitempty"`
}

type Channel struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type User struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type ActionResponse struct {
	Actions      []Action `json:"actions,omitempty"`
	CallbackID   string   `json:"callback_id,omitempty"`
	Team         Team     `json:"team,omitempty"`
	Channel      Channel  `json:"channel,omitempty"`
	User         User     `json:"user,omitempty,string"`
	ActionTS     float64  `json:"action_ts,omitempty,string"`
	MessageTS    float64  `json:"message_ts,omitempty,string"`
	AttachmentID int      `json:"attachment_id,omitempty,string"`
	Token        string   `json:"token,omitempty"`
	AppUnfurl    bool     `json:"is_app_unfurl,omitempty"`
	ResponseURL  string   `json:"response_url,omitempty"`
}

func NewActionResponseFromPayload(jsonPaylaod string) (ActionResponse, error) {
	var resp ActionResponse
	err := json.Unmarshal([]byte(jsonPaylaod), &resp)
	return resp, err
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

func NewPollMessage(callbackID uuid.UUID, question string, options ...string) SlackMessage {
	var msg SlackMessage
	msg.Text = question
	var buttonAttachment Attachment
	buttonAttachment.Fallback = "Poll not available"
	buttonAttachment.CallbackID = callbackID.String()
	for index, option := range options {
		var button Action
		button.Name = option + "_button"
		button.Text = option
		button.Type = "button"
		button.Value = strconv.Itoa(index)
		buttonAttachment.AddAction(button)
	}
	msg.AddAttachment(buttonAttachment)
	return msg
}

func UpdatePollMessage(poll poll.Poll, callback ActionResponse, results map[string]uint64) SlackMessage {
	var msg SlackMessage
	msg.Text = poll.Question()
	var buttonAttachment Attachment
	buttonAttachment.Fallback = "Poll not available"
	buttonAttachment.CallbackID = poll.ID()
	for index, option := range poll.Options() {
		var button Action
		button.Name = option + "_button"
		button.Text = option + " " + fmt.Sprintf("%d", results[strconv.Itoa(index)])
		button.Type = "button"
		button.Value = strconv.Itoa(index)
		buttonAttachment.AddAction(button)
	}
	msg.AddAttachment(buttonAttachment)
	return msg
}
