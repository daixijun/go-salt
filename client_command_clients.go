package salt

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

type commandRequest struct {
	Client     CommandClient          `json:"client"`
	Target     interface{}            `json:"tgt,omitempty"` // string or slices
	Function   string                 `json:"fun"`
	Arguments  []string               `json:"arg,omitempty"`   // []string or [][]string
	Match      string                 `json:"match,omitempty"` // Wheel Client
	KwArg      map[string]interface{} `json:"kwarg,omitempty"`
	TargetType TargetType             `json:"tgt_type,omitempty"`
	Timeout    int                    `json:"timeout,omitempty"`
	FullReturn bool                   `json:"full_return,omitempty"`
}

type RunOption func(*commandRequest)

func WithTargetType(t TargetType) RunOption {
	return func(r *commandRequest) {
		r.TargetType = t
	}
}
