package salt

import (
	"context"
	"encoding/json"
	"fmt"
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
		SortList      []string `json:"sortlist"`
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
	CPUArch  string   `json:"cpuarch"`
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
	LSBDistribID       string              `json:"lsb_distrib_id"`
	LSBDistribCodename string              `json:"lsb_distrib_codename"`
	OS                 string              `json:"os"`
	OSFullName         string              `json:"osfullname"`
	OSRelease          string              `json:"osrelease"`
	OSReleaseInfo      []int               `json:"osrelease_info"`
	OSCodename         string              `json:"oscodename"`
	OSMajorRelease     int                 `json:"osmajorrelease"`
	OSFinger           string              `json:"osfinger"`
	OSFamily           string              `json:"os_family"`
	OSArch             string              `json:"osarch"`
	MemTotal           int                 `json:"mem_total"`
	SwapTotal          int                 `json:"swap_total"`
	BiosVersion        string              `json:"biosversion"`
	BiosReleaseDate    string              `json:"biosreleasedate"`
	ProductName        string              `json:"productname"`
	Manufacturer       string              `json:"manufacturer"`
	UUID               string              `json:"uuid"`
	SerialNumber       string              `json:"serialnumber"`
	Virtual            string              `json:"virtual"`
	PS                 string              `json:"ps"`
	Path               string              `json:"path"`
	SystemPath         []string            `json:"systempath"`
	PythonExecutable   string              `json:"pythonexecutable"`
	PythonPath         []string            `json:"pythonpath"`
	PythonVersion      []interface{}       `json:"pythonversion"`
	SaltPath           string              `json:"saltpath"`
	SaltVersion        string              `json:"saltversion"`
	SaltVersionInfo    []int               `json:"saltversioninfo"`
	ZMQVersion         string              `json:"zmqversion"`
	Disks              []string            `json:"disks"`
	SSDs               []string            `json:"ssds"`
	Shell              string              `json:"shell"`
	Lvm                map[string][]string `json:"lvm"`
	MDAdm              []string            `json:"mdadm"`
	Username           string              `json:"username"`
	GroupName          string              `json:"groupname"`
	Pid                int64               `json:"pid"`
	Gid                int                 `json:"gid"`
	Uid                int                 `json:"uid"`
	ZFSSupport         bool                `json:"zfs_support"`
	ZFSFeatureFlags    bool                `json:"zfs_feature_flags"`
	Transactional      bool                `json:"transactional,omitempty"`
	Roles              []string            `json:"roles,omitempty"`
}

type minionResponse struct {
	Return []map[string]Minion `json:"return"`
}

func (c *Client) ListMinions(ctx context.Context) ([]Minion, error) {
	data, err := c.get(ctx, "minions")
	if err != nil {
		return nil, err
	}

	var resp minionResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	minions := make([]Minion, 0)
	for _, v := range resp.Return[0] {
		minions = append(minions, v)
	}

	return minions, nil
}

func (c *Client) GetMinion(ctx context.Context, mid string) (*Minion, error) {
	data, err := c.get(ctx, fmt.Sprintf("minions/%s", mid))
	if err != nil {
		return nil, err
	}

	var resp minionResponse
	err = json.Unmarshal(data, &resp)
	if err != nil {
		return nil, err
	}

	minionTotal := len(resp.Return[0])
	if minionTotal == 0 {
		return nil, fmt.Errorf("minion %s not found", mid)
	} else if minionTotal > 1 {
		return nil, fmt.Errorf("expected one return but received %d", len(resp.Return))
	}

	minion := resp.Return[0][mid]
	return &minion, nil

}
