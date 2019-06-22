package common

import "strings"

type Target struct {
	Who     string `json:"who"`
	Targets string `json:"targets"`
}

func (t *Target) Sanitize() error {
	t.Who = strings.ToLower(t.Who)
	t.Targets = strings.ToLower(t.Targets)
	return nil
}

const (
	TargetWho_AllAlly    = "all ally"
	TargetWho_TargetAlly = "target ally"
	TargetWho_OtherAlly  = "other allys"
	TargetWho_Target     = "target"
)
