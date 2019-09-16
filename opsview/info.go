package opsview

import (
	"log"
	"errors"
	"strconv"
	"strings"
	"encoding/json"

	"github.com/skirvin/opsview-go/rest"
)

type Object struct {
	Name	string `json:"name,omitempty"`
}

type Reference struct {
	Reference	string `json:"ref,omitempty"`
	Name		string `json:"name,omitempty"`
}

type Summary struct {
	Page 		string	`json:"page,omitempty"`
	TotalRows 	string	`json:"totalrows,omitempty"`
	AllRows 	string	`json:"allrows,omitempty"`
	TotalPages 	string	`json:"totalpages,omitempty"`
	Rows 		string	`json:"rows,omitempty"`
}

type Info struct {
	Version					string	`json:"opsview_version,omitempty"`
	Build					string	`json:"opsview_build,omitempty"`
	Edition					string	`json:"opsview_edition,omitempty"`
	Name					string	`json:"opsview_name,omitempty"`
	HostLimit				string	`json:"hosts_limit,omitempty"`
	ServerTimezone			string	`json:"server_timezone,omitempty"`
	ServerTimeZoneOffset 	string	`json:"server_timezone_offset,omitempty"`
	UUID					string	`json:"uuid,omitempty"`
}

func (c *Client) GetInfo() (Info, error) {
	var (
		opsviewInfo	Info
	)

	data, err := c.RestAPICall(rest.GET, "/rest/info", nil)
	if err != nil {
		return opsviewInfo, err
	}

	if(c.Verbose) {
		log.Printf("Info %s", data)
	}

	if err := json.Unmarshal([]byte(data), &opsviewInfo); err != nil {
		return opsviewInfo, err
	}

	return opsviewInfo, nil
}

func (c *Client) GetVersion() (int, error) {
	version, err := strconv.Atoi(strings.Replace(c.Info.Version, ".", "", -1))
	if err != nil {
		return 0, errors.New("unable to determine OpsView version")
	}

	return version, nil
}