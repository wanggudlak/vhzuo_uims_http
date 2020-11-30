package conf

import "uims/pkg/env"

const (
	UIMS_NODE_ID = 1
	CASS_NODE_ID = 2
)

type NodeConf struct {
	Epoch int64
}

type NodeID int

type NodeMap map[NodeID]*NodeConf

func NewNodeMap() *NodeMap {
	return &NodeMap{
		UIMS_NODE_ID: {
			Epoch: int64(env.DefaultGetInt("UIMS_UUID_START_TIME", 1288834974657)),
		},
		CASS_NODE_ID: {
			Epoch: int64(env.DefaultGetInt("CASS_UUID_START_TIME", 1480166465631)),
		},
	}
}
