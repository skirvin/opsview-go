package opsview

import (
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/skirvin/opsview-go/rest"
)

type ReloadStatus struct {
	ServerStatus       string    `json:"server_status,omitempty"`
	ConfigurationStaus string    `json:"configuration_status,omitempty"`
	AverageDuration    string    `json:"average_duration,omitempty"`
	LastUpdated        string    `json:"lastupdated,omitempty"`
	AuditLogEntries    string    `json:"auditlog_entries,omitempty"`
	Messages           []Message `json:"messages,omitempty"`
}

type Message struct {
	Detail           string `json:"detail,omitempty"`
	MonitoringServer string `json:"monitoringserver,omitempty"`
	Severity         string `json:"severity,omitempty"`
}

func (c *Client) GetReloadStatusDetails() (ReloadStatus, error) {
	var (
		uri          = "/rest/reload"
		reloadStatus ReloadStatus
	)

	available, err := c.FunctionalityAvailable("RELOAD")
	if err != nil {
		return reloadStatus, err
	}

	if available {
		if(c.Verbose) {
			log.Printf("Getting OpsView Reload Status from %s\n", c.BaseURI)
		}

		data, err := c.RestAPICall(rest.GET, uri, nil)
		if err != nil {
			return reloadStatus, err
		}

		if(c.Verbose) {
			log.Printf("Reload State: %s", data)
		}

		if err := json.Unmarshal([]byte(data), &reloadStatus); err != nil {
			return reloadStatus, err
		}
	}

	return reloadStatus, nil
}

func (c *Client) DoReload() (ReloadStatus, error) {
	var (
		uri          = "/rest/reload"
		qp           = make(map[string]interface{})
		reloadStatus ReloadStatus
	)

	available, err := c.FunctionalityAvailable("RELOAD")
	if err != nil {
		return reloadStatus, err
	}

	if available {
		log.Printf("Initiating ASync OpsView Reload on %s\n", c.BaseURI)
		if err != nil {
			return reloadStatus, err
		}

		qp["asynchronous"] = "1"
		c.SetQueryString(qp)

		var data []byte
		data, err = c.RestAPICall(rest.POST, uri, nil)

		if err != nil {
			return reloadStatus, err
		}

		if err := json.Unmarshal([]byte(data), &reloadStatus); err != nil {
			return reloadStatus, err
		}

		c.WaitForReload()
		if err != nil {
			return reloadStatus, err
		}
	}

	return reloadStatus, nil
}

func (c *Client) IsReloading() (bool, error) {
	c.GetReloadStatusDetails()
	var status, err = strconv.Atoi(c.ReloadStatus.ServerStatus)

	if err != nil {
		return true, err
	}

	switch status {
	case 1, 2, 3:
		return true, nil
	}

	return false, nil
}

func (c *Client) WaitForReload() error {
	isReloading, err := c.IsReloading()
	if err != nil {
		return err
	}

	if isReloading {
		time.Sleep(time.Second * 30)
		c.WaitForReload()
	}
	return nil
}
