package opsview

import (
	"log"
	"time"
	"strings"
	"encoding/json"

	"github.com/skirvin/opsview-go/rest"
)

type XOpsviewToken struct {
	Token	string	`json:"token,omitempty"`
}

type TimeOut struct {
	Duration time.Duration
}

type AuthBody struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type AuthHeader struct {
	XOpsviewUsername	string	`json:"X-Opsview-Username,omitempty"`
	XOpsviewToken		string	`json:"X-Opsview-Token,omitempty"`
}

func (c *Client) IsSessionValid() (bool) {
	var (
		timeout = TimeOut{Duration: 3600}
	)

	return time.Since(c.LastAuthenticated).Minutes() > timeout.Duration.Minutes()
}

func (c *Client) GetAuthHeaderMap() map[string]string {
	return map[string]string{
		"X-OPSVIEW-USERNAME":	c.Username,
		"X-OPSVIEW-TOKEN":		c.Options.Headers["X-OPSVIEW-TOKEN"],
	}
}

func (c *Client) RefreshLogin() error {
	var authToken = c.Options.Headers["X-OPSVIEW-TOKEN"]

	if authToken == "" || len(strings.TrimSpace(authToken)) == 0 || authToken == "none" {
		if(c.Verbose) {
			log.Print("Getting new X-OPSVIEW-TOKEN")
		}
		body, err := c.Login()
		if err != nil {
			return err
		}
		c.Options.Headers["X-OPSVIEW-TOKEN"] = body.Token
	}

	if !c.IsSessionValid() {
		body, err := c.Login()
		if err != nil {
			return err
		}
		c.Options.Headers["X-OPSVIEW-TOKEN"] = body.Token
	}
	return nil
}

func (c *Client) Login() (XOpsviewToken, error) {
	var (
		uri     = "/rest/login"
		body    = AuthBody{Username: c.Username, Password: c.Password}
		token	XOpsviewToken
	)

	c.SetHeaders(
		map[string]string{
			"Accept": "application/json",
			"Content-Type": "application/json",
		},
	)

	data, err := c.RestAPICall(rest.POST, uri, body)
	if err != nil {
		return token, err
	}

	if err := json.Unmarshal([]byte(data), &token); err != nil {
		return token, err
	}

	c.SetHeaders(c.GetAuthHeaderMap())
	c.LastAuthenticated = time.Now()
	return token, err
}