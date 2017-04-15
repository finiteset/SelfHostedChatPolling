package slack

import (
	"encoding/json"
)

type SlackMessage struct {
	Text string `json:"text"`
}

func (m *SlackMessage) ToJSON() ([]byte, error) {
	return json.Marshal(m)
}
