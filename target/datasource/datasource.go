package datasource

import (
	"context"
	"fmt"

	"github.com/antha-lang/antha/ast"
	"github.com/antha-lang/antha/target"
	"github.com/antha-lang/antha/target/human"
)

var (
	_ target.Device = DataSource{}
)

type DataSource struct {
}

func (ds DataSource) CanCompile(req ast.Request) bool {
	can := ast.Request{
		Selector: []ast.NameValue{
			target.DriverSelectorV1DataSource,
		},
	}
	return can.Contains(req)
}

func (ds DataSource) MoveCost(from target.Device) int {
	return human.HumanByXCost + 1 // same as a mixer
}

func (ds DataSource) Compile(ctx context.Context, cmds []ast.Node) (insts []target.Inst, err error) {
	ret := make([]target.Inst, 0, 1)

	// only support one command... for now?

	if len(cmds) != 1 {
		return ret, fmt.Errorf("Multiple GetData commands not supported (%d received)", len(cmds))
	}

	// XXX: rest of parameters
	node := cmds[0]

	c, ok := node.(*ast.Command)
	if !ok {
		return nil, fmt.Errorf("cannot compile %T", node)
	}

	m, ok := c.Inst.(*ast.AwaitInst)
	if !ok {
		return nil, fmt.Errorf("cannot compile %T", c.Inst)
	}

	ret = append(ret, &target.AwaitData{
		Dev:     ds,
		Tags:    m.Tags,
		AwaitID: m.AwaitID,
	})

	return ret, nil
}
