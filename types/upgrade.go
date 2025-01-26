package types

import "fmt"

type Upgrade struct {
	Version            string
	ClusterName        string
	Iteration          int
	RolloutRestartDate string
}

func (u Upgrade) ToConfigMapData() map[string]string {
	return map[string]string{
		"version":              u.Version,
		"cluster_name":         u.ClusterName,
		"iteration":            fmt.Sprint(u.Iteration),
		"rollout_restart_date": u.RolloutRestartDate,
	}
}
