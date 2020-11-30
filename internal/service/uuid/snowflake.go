package uuid

import (
	"github.com/bwmarrin/snowflake"
	"strconv"
	"sync"
	"uims/conf"
)

type ID int64

var nodePool = sync.Pool{
	New: func() interface{} { return newNodes() },
}

type Nodes map[conf.NodeID]*snowflake.Node

func newNodes() *Nodes {
	nodes := Nodes{}
	for nodeID, confV := range *conf.UUIDConf {
		if confV.Epoch > 0 {
			snowflake.Epoch = confV.Epoch
		}
		pNode, err := snowflake.NewNode(int64(nodeID))
		if err != nil {
			nodes[nodeID] = pNode
		}
	}
	return &nodes
}

// GenerateForUIMS 为UIMS系统用雪花算法生成唯一id
func GenerateForUIMS() (id ID) {
	var err error
	currNodeID := conf.NodeID(conf.UIMS_NODE_ID)
	pNodes := nodePool.Get().(*Nodes)
	if pND, ok := (*pNodes)[currNodeID]; ok && pND != nil {
		id = ID(pND.Generate())
	} else {
		pND, err = snowflake.NewNode(int64(currNodeID))
		if pND != nil {
			(*pNodes)[currNodeID] = pND
			id = ID(pND.Generate())
		}
	}
	nodePool.Put(pNodes)

	if err != nil {
		panic(err)
	}

	return
}

// GenerateForCASS 为结算系统用雪花算法生成唯一id
func GenerateForCASS() (id ID) {
	var err error
	currNodeID := conf.NodeID(conf.CASS_NODE_ID)
	pNodes := nodePool.Get().(*Nodes)
	if pND, ok := (*pNodes)[currNodeID]; ok && pND != nil {
		id = ID(pND.Generate())
	} else {
		pND, err = snowflake.NewNode(int64(currNodeID))
		if pND != nil {
			(*pNodes)[currNodeID] = pND
			id = ID(pND.Generate())
		}
	}
	nodePool.Put(pNodes)

	if err != nil {
		panic(err)
	}

	return
}

// String returns a string of the snowflake ID
func (id ID) String() string {
	return strconv.FormatInt(int64(id), 10)
}
