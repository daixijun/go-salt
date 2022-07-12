package salt

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
	IPCidr     TargetType = "ipcidr"
)
