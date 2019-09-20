package opsview

import (
	"encoding/json"
	"fmt"
	"github.com/skirvin/opsview-go/rest"
	"log"
)

type HostAttribute struct {
	Arg1	string `json:"arg1,omitempty"`
	Arg2	string `json:"arg2,omitempty"`
	Arg3	string `json:"arg3,omitempty"`
	Arg4	string `json:"arg4,omitempty"`
	Name	string `json:"name,omitempty"`
	Value	string `json:"value,omitempty"`
}

type ServiceCheck struct {
	Reference		string `json:"ref,omitempty"`
	Name			string `json:"name,omitempty"`
	Exception		string `json:"exception,omitempty"`
	TimedException	TimedException `json:"timed_exception,omitempty"`
	EventHandler	string `json:"event_handler,omitempty"`
	Remove			string `json:"remove_servicecheck,omitempty"`
}

type TimedException struct {
	TimePeriod		Object `json:"timeperiod,omitempty"`
	Args			string `json:"args,omitempty"`
}

type Host struct {
	ID								string `json:"id,omitempty"`
	Name							string `json:"name,omitempty"`
	Alias							string `json:"alias,omitempty"`
	IP								string `json:"ip,omitempty"`
	OtherAddresses					string `json:"other_addresses,omitempty"`
	Uncommited						string `json:"uncommitted,omitempty"`
	SNMPEnabled						string `json:"enable_snmp,omitempty"`
	SNMPVersion						string `json:"snmp_version,omitempty"`
	SNMPv3Enabled					string `json:"snmpv3_privprotocol,omitempty"`
	SNMPv3AuthenticationProtocol	string `json:"snmpv3_authprotocol,omitempty"`
	SNMPv3Username					string `json:"snmpv3_privpassword,omitempty"`
	SNMPv3Password					string `json:"snmpv3_privpassword,omitempty"`
	FlapDetectionEnabled			string `json:"flap_detection_enabled,omitempty"`
	HostTemplates					[]Reference `json:"hosttemplates,omitempty"`
	Keywords						[]Reference `json:"keywords,omitempty"`
	CheckPeriod						Object `json:"check_period,omitempty"`
	HostAttributes					[]HostAttribute `json:"hostattributes,omitempty"`
	NotificationPeriod				Object `json:"notification_period,omitempty"`
	NotificationOptions				string `json:"notification_options,omitempty"`
	RancidVendor					string `json:"rancid_vendor,omitempty"`
	HostGroup						Object `json:"hostgroup,omitempty"`
	EventHandler					string `json:"event_handler,omitempty"`
	MonitoredBy						Object `json:"monitored_by,omitempty"`
	Parents							[]Reference `json:"parents,omitempty"`
	RetryCheckInterval				string `json:"retry_check_interval,omitempty"`
	Icon							Object `json:"icon,omitempty"`
	UseMRTG							[]Reference `json:"use_mrtg,omitempty"`
	ServiceChecks					[]ServiceCheck `json:"servicechecks,omitempty"`
	UseRancid						string `json:"use_rancid,omitempty"`
	NMISNodeType					string `json:"nmis_node_type,omitempty"`
	UseNMIS							string `json:"use_nmis,omitempty"`
	RancidConnectionType			string `json:"rancid_connection_type,omitempty"`
	RancidUsername					string `json:"rancid_username,omitempty"`
	RancidPassword					string `json:"rancid_password,omitempty"`
	CheckCommand					Object `json:"check_command,omitempty"`
	CheckAttempts					string `json:"check_attempts,omitempty"`
	CheckInterval					string `json:"check_interval,omitempty"`
	NotificationInterval			string `json:"notification_interval,omitempty"`
	BusinessComponents				[]Reference `json:"business_components,omitempty"`
}

type HostWrapper struct {
	Summary	Summary `json:"summary,omitempty"`
	List	[]Host `json:"list,omitempty"`
}

func (c *Client) HostExists(name string) (bool, error) {
	var (
		uri  = fmt.Sprintf("/rest/config/host/exists?name=%s", name)
	)

	log.Printf("Checking Host exists on %s\n", c.BaseURI)
	data, err := c.RestAPICall(rest.GET, uri, nil)
	if err != nil {
		return false, err
	}

	if(c.Verbose) {
		fmt.Printf("Host status: %+v", data)
	}

	return true, err
}

func (c *Client) GetHostsByHostTemplates(names []string) ([]Host, error) {

	var (
		uri          = "/rest/config/host"
		hostWrapper HostWrapper
		hosts		[]Host
		queryParams = make(map[string]interface{})
	)

	queryParams["rows"] = "all"
	queryParams["group_by"] = "host"

	for _, name := range names {
		queryParams["s.hosttemplates.name"] = name
	}

	available, err := c.FunctionalityAvailable("CONFIG_HOST")
	if err != nil {
		return hosts, err
	}

	if available {
		if(c.Verbose) {
			log.Printf("Getting Hosts by Host Template(s): %s\n", names)
		}

		c.SetQueryString(queryParams)
		data, err := c.RestAPICall(rest.GET, uri, nil)
		if err != nil {
			return hosts, err
		}

		if(c.Verbose) {
			log.Printf("Host(s) Config: %s", data)
		}
		if err := json.Unmarshal([]byte(data), &hostWrapper); err != nil {
			return hosts, err
		}
		hosts = hostWrapper.List
	}

	return hosts, nil
}
