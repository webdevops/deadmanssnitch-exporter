package dmsclient

import (
	"encoding/json"
	"strings"
	"time"
)

type (
	ResponseSnitchList []ResponseSnitch
	ResponseSnitch     struct {
		Token       string        `json:"token"`
		Href        string        `json:"href"`
		Name        string        `json:"name"`
		Tags        []string      `json:"tags"`
		Notes       string        `json:"notes,omitempty"`
		Status      string        `json:"status"`
		CheckedInAt *time.Time    `json:"checked_in_at"`
		CreatedAt   time.Time     `json:"created_at"`
		Interval    string        `json:"interval"`
		AlertType   string        `json:"alert_type"`
		AlertEmail  []interface{} `json:"alert_email"`
	}
)

func (c *Client) ListSnitches() (list ResponseSnitchList, error error) {
	response, err := c.rest().R().Get("/snitches")
	if err := c.checkResponse(response, err); err != nil {
		error = err
		return
	}

	err = json.Unmarshal(response.Body(), &list)
	if err != nil {
		error = err
		return
	}

	return
}

func (s *ResponseSnitch) IsHealthy() bool {
	switch strings.ToLower(s.Status) {
	case "healthy":
		return true
	}

	return false
}
