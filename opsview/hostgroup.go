package opsview

import (
	"encoding/json"
	"fmt"
	"github.com/skirvin/opsview-go/rest"
	"log"
)

type Parent struct {
	Name		string	`json:"name,omitempty"`
	MatPath		string	`json:"matpath,omitempty"`
	Reference	string	`json:"ref,omitempty"`
}

type HostGroup struct {
	ID			string	`json:"id,omitempty"`
	Name		string	`json:"name,omitempty"`
	Hosts		[]Reference	`json:"hosts,omitempty"`
	Parent		Parent	`json:"parent,omitempty"`
	Children	[]Reference	`json:"children,omitempty"`
	Uncommited	string	`json:"uncommitted,omitempty"`
	Reference	string	`json:"ref,omitempty"`
	MatPath		string	`json:"matpath,omitempty"`
	IsLeaf		string	`json:"is_leaf,omitempty"`
}

type HostGroupWrapper struct {
	Summary	Summary `json:"summary,omitempty"`
	List	[]HostGroup `json:"list,omitempty"`
}


func (c *Client) GetHostGroups() ([]HostGroup, error) {
	var (
		hostgroupWrapper HostGroupWrapper
		opsviewHostGroups []HostGroup
		uri = "/rest/config/hostgroup"
	)

	data, err := c.RestAPICall(rest.GET, uri, nil)
	if err != nil {

		return opsviewHostGroups, err
	}

	if(c.Verbose) {
		log.Printf("HostGroups: %s", data)
	}

	if err := json.Unmarshal([]byte(data), &hostgroupWrapper); err != nil {
		return opsviewHostGroups, err
	}
	opsviewHostGroups = hostgroupWrapper.List

	return opsviewHostGroups, nil
}

func (c *Client) GetHostGroupByName(name string) (HostGroup, error) {
	for _, hostGroup := range c.HostGroups {
		if(hostGroup.Name == name) {
			return hostGroup, nil
		}
	}

	return HostGroup{}, fmt.Errorf("%s", "not found")
}

func (c *Client) GetHostsStatusByHostGroup(name string) ([]Host, error) {

	var (
		uri          = "/rest/status/hostgroup"
		hostDetails []Host
		queryParams = make(map[string]interface{})
	)

	queryParams["fromhostgroupid"] = name

	available, err := c.FunctionalityAvailable("CONFIG_HOST")
	if err != nil {
		return hostDetails, err
	}

	if available {
		log.Printf("Getting OpsView Reload Status from %s\n", c.BaseURI)

		c.SetQueryString(queryParams)
		data, err := c.RestAPICall(rest.GET, uri, nil)
		if err != nil {
			return hostDetails, err
		}

		log.Printf("Reload State: %s", data)
		if err := json.Unmarshal([]byte(data), &hostDetails); err != nil {
			return hostDetails, err
		}
	}

	return hostDetails, nil
}