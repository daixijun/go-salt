package salt

import (
	"net"
	"strings"
)

type CommandClient string

const (
	LocalClient       CommandClient = "local"
	LocalAsyncClient  CommandClient = "local_async"
	RunnerClient      CommandClient = "runner"
	RunnerAsyncClient CommandClient = "runner_async"
	WheelClient       CommandClient = "wheel"
	WheelAsyncClient  CommandClient = "wheel_async"
	SSHClient         CommandClient = "ssh"
)

type TargetType string

const (
	Glob       TargetType = "glob"
	Pcre       TargetType = "pcre"
	List       TargetType = "list"
	Grain      TargetType = "grain"
	GrainPcre  TargetType = "grain_pcre"
	Pillar     TargetType = "pillar"
	PillarPcre TargetType = "pillar_pcre"
	NodeGroup  TargetType = "nodegroup"
	Range      TargetType = "range"
	Compound   TargetType = "compound"
	IPCIDR     TargetType = "ipcidr"
)

type commandRequest struct {
	Client     CommandClient          `json:"client"`
	Target     string                 `json:"tgt,omitempty"` // string or slices
	Function   string                 `json:"fun"`
	Arguments  []string               `json:"arg,omitempty"`   // []string or [][]string
	Match      string                 `json:"match,omitempty"` // Wheel Client
	KwArg      map[string]interface{} `json:"kwarg,omitempty"`
	TargetType TargetType             `json:"tgt_type,omitempty"`
	Timeout    int                    `json:"timeout,omitempty"`
	FullReturn bool                   `json:"full_return,omitempty"`
}

type RunOption func(*commandRequest)

func WithGlobTarget(t string) RunOption {
	return func(r *commandRequest) {
		r.Target = t
		r.TargetType = Glob
	}
}

func WithPcreTarget(t string) RunOption {
	return func(r *commandRequest) {
		r.Target = t
		r.TargetType = Pcre
	}
}

func WithListTarget(t []string) RunOption {
	return func(r *commandRequest) {
		r.Target = strings.Join(t, ",")
		r.TargetType = List
	}
}

func WithPillarTarget(t string) RunOption {
	return func(r *commandRequest) {
		r.Target = t
		r.TargetType = Pillar
	}
}

func WithPillarPcreTarget(t string) RunOption {
	return func(r *commandRequest) {
		r.Target = t
		r.TargetType = PillarPcre
	}
}

func WithGrainTarget(t string) RunOption {
	return func(r *commandRequest) {
		r.Target = t
		r.TargetType = Grain
	}
}

func WithGrainPcreTarget(t string) RunOption {
	return func(r *commandRequest) {
		r.Target = t
		r.TargetType = GrainPcre
	}
}

func WithNodeGroupTarget(t string) RunOption {
	return func(r *commandRequest) {
		r.Target = t
		r.TargetType = NodeGroup
	}
}

func WithIPCIDRTarget(t net.IPNet) RunOption {
	return func(r *commandRequest) {
		r.Target = t.String()
		r.TargetType = IPCIDR
	}
}

func WithKeywordArguments(kw map[string]interface{}) RunOption {
	return func(r *commandRequest) {
		r.KwArg = kw
	}
}
