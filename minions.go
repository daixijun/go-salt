package salt

import (
	"encoding/json"
)

type Minion struct {
	ID        string      `json:"id"`
	CWD       string      `json:"cwd"`
	NodeName  string      `json:"nodename"`
	MachineID string      `json:"machine_id"`
	Master    string      `json:"master"`
	ServerID  int64       `json:"server_id"`
	Localhost string      `json:"localhost"`
	Host      string      `json:"host"`
	Domain    string      `json:"domain"`
	IPGW      interface{} `json:"ip_gw"`
	IP4GW     interface{} `json:"ip4_gw"`
	IP6GW     interface{} `json:"ip6_gw"`
	DNS       struct {
		Nameservers   []string `json:"nameservers"`
		IP4Namespaces []string `json:"ip4_namespaces"`
		IP6Namespaces []string `json:"ip6_namespaces"`
		Sortlist      []string `json:"sortlist"`
		Domain        string   `json:"domain"`
		Search        []string `json:"search"`
		Options       []string `json:"options"`
	} `json:"dns"`
	FQDN            string              `json:"fqdn"`
	FQDNs           []string            `json:"fqdns"`
	HWAddrInterface map[string]string   `json:"hwaddr_interface"`
	IP4Interfaces   map[string][]string `json:"ip4_interfaces"`
	IP6Interfaces   map[string][]string `json:"ip6_interfaces"`
	IPInterfaces    map[string][]string `json:"ip_interfaces"`
	IPv4            []string            `json:"ipv4"`
	IPv6            []string            `json:"ipv6"`
	FQDNIP4         []string            `json:"fqdn_ip4"`
	FQDNIP6         []string            `json:"fqdn_ip6"`
	Kernel          string              `json:"kernel"`
	KernelRelease   string              `json:"kernelrelease"`
	KernelVersion   string              `json:"kernelversion"`
	KernelParams    [][]string          `json:"kernelparams"`
	LocaleInfo      struct {
		DefaultLanguage  string `json:"defaultlanguage"`
		DefaultEncoding  string `json:"defaultencoding"`
		DetectedEncoding string `json:"detectedencoding"`
		TimeZone         string `json:"timezone"`
	} `json:"locale_info"`
	NumGPUS int `json:"num_gpus"`
	GPUS    []struct {
		Vendor string `json:"vendor"`
		Model  string `json:"model"`
	} `json:"gpus"`
	NumCPUS  int      `json:"num_cpus"`
	CPUarch  string   `json:"cpuarch"`
	CPUModel string   `json:"cpu_model"`
	CPUFlags []string `json:"cpu_flags"`
	Selinux  struct {
		Enabled  bool   `json:"enabled"`
		Enforced string `json:"enforced"`
	} `json:"selinux"`
	Systemd struct {
		Version  string `json:"version"`
		Features string `json:"features"`
	} `json:"systemd"`
	Init               string              `json:"init"`
	LsbDistribID       string              `json:"lsb_distrib_id"`
	LsbDistribCodename string              `json:"lsb_distrib_codename"`
	OS                 string              `json:"os"`
	Osfullname         string              `json:"osfullname"`
	Osrelease          string              `json:"osrelease"`
	OsreleaseInfo      []int               `json:"osrelease_info"`
	Oscodename         string              `json:"oscodename"`
	Osmajorrelease     int                 `json:"osmajorrelease"`
	OSFinger           string              `json:"osfinger"`
	OsFamily           string              `json:"os_family"`
	Osarch             string              `json:"osarch"`
	MemTotal           int                 `json:"mem_total"`
	SwapTotal          int                 `json:"swap_total"`
	Biosversion        string              `json:"biosversion"`
	Biosreleasedate    string              `json:"biosreleasedate"`
	Productname        string              `json:"productname"`
	Manufacturer       string              `json:"manufacturer"`
	UUID               string              `json:"uuid"`
	Serialnumber       string              `json:"serialnumber"`
	Virtual            string              `json:"virtual"`
	PS                 string              `json:"ps"`
	Path               string              `json:"path"`
	Systempath         []string            `json:"systempath"`
	Pythonexecutable   string              `json:"pythonexecutable"`
	Pythonpath         []string            `json:"pythonpath"`
	Pythonversion      []interface{}       `json:"pythonversion"`
	Saltpath           string              `json:"saltpath"`
	Saltversion        string              `json:"saltversion"`
	Saltversioninfo    []int               `json:"saltversioninfo"`
	ZMQVersion         string              `json:"zmqversion"`
	Disks              []string            `json:"disks"`
	Ssds               []string            `json:"ssds"`
	Shell              string              `json:"shell"`
	Lvm                map[string][]string `json:"lvm"`
	Mdadm              []string            `json:"mdadm"`
	Username           string              `json:"username"`
	Groupname          string              `json:"groupname"`
	Pid                int64               `json:"pid"`
	Gid                int                 `json:"gid"`
	Uid                int                 `json:"uid"`
	ZFSSupport         bool                `json:"zfs_support"`
	ZFSFeatureFlags    bool                `json:"zfs_feature_flags"`
}

type MinionResponse struct {
	Return []map[string]Minion `json:"return"`
}

func (c *Client) Minions() (*MinionResponse, error) {
	// data, err := c.doRequest(ctx, "GET", "minions", nil)
	// if err != nil {
	// 	return nil, err
	// }
	//
	// var minions MinionResponse
	// if err := json.Unmarshal(data, &minions); err != nil {
	// 	return nil, err
	// }
	// return &minions, nil
	return c.Minion("")
}

func (c *Client) Minion(mid string) (*MinionResponse, error) {
	var uri string
	if mid == "" {
		uri = "minions"
	} else {
		uri = "minions/" + mid
	}
	data, err := c.doRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}

	var minion MinionResponse
	if err := json.Unmarshal(data, &minion); err != nil {
		return nil, err
	}
	return &minion, nil
}
