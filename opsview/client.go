package opsview

import (
	"log"

	"github.com/skirvin/opsview-go/rest"
)

var functionalityAvailability = map[string]int{
	"RELOAD": 1,
	"CONFIG_HOST": 2,
}

type Client struct {
	rest.Client
	HostGroups []HostGroup
	Info
	ReloadStatus
}

func NewClient(username string, password string, baseUri string, sslVerify bool, verbose bool) (*Client) {
	var (
		hostGroups  []HostGroup
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
			Verbose:		verbose,
		},
		hostGroups,
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

	hostGroups, err = client.GetHostGroups()
	if err != nil {
		log.Panic(err)
	}
	client.HostGroups = hostGroups

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