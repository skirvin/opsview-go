package opsview

import (
	"log"

	"git.monitoring.bskyb.com/monitoring/opsview-go/rest"
)

var functionalityAvailability = map[string]int{
	"RELOAD": 391,
}

type Client struct {
	rest.Client
	Info
	ReloadStatus
}

func NewClient(username string, password string, baseUri string, sslVerify bool) (*Client) {
	var (
		info		Info
		reloadInfo	ReloadStatus
		err			error
	)

	client := &Client{
		rest.Client{
			Username:		username,
			Password:   	password,
			BaseURI:		baseUri,
			SSLVerify:		sslVerify,
		},
		info,
		reloadInfo,
	}

	client.RefreshLogin()
	client.SetHeaders(client.GetAuthHeaderMap())

	info, err = client.GetInfo()
	if err != nil {
		log.Panic(err)
	}

	client.Info = info

	reloadInfo, err = client.GetReloadStatusDetails()
	if err != nil {
		log.Panic(err)
	}

	client.ReloadStatus = reloadInfo

	return client
}

func (c *Client) FunctionalityAvailable(functionalityKey string) (bool, error) {
	version, err := c.GetVersion()

	if err != nil {
		return false, err
	}

	if version > functionalityAvailability[functionalityKey] {
		return true, nil
	}

	return false, nil
}