package snow

import (
	"github.com/bwmarrin/snowflake"
	"github.com/golang-module/carbon/v2"
)

type Flake struct {
	Start string
	Size  int64
	Node  *snowflake.Node
}

type FlakeID struct {
	ID     int64
	Str    string
	Base64 string
}

func NewFlake(start string, size int64) (*Flake, error) {
	if start == "" {
		now := carbon.Now()
		start = now.Format("Y-m-d")
	}
	snowflake.Epoch = carbon.Parse(start).TimestampMilli()
	node, err := snowflake.NewNode(size)
	if err != nil {
		return nil, err
	}
	return &Flake{
		Start: start,
		Size:  size,
		Node:  node,
	}, nil
}

func (f *Flake) Make() *FlakeID {
	return &FlakeID{
		ID:     f.Node.Generate().Int64(),
		Str:    f.Node.Generate().String(),
		Base64: f.Node.Generate().Base64(),
	}
}
