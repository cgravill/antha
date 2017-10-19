package auto

import (
	"context"

	"github.com/antha-lang/antha/ast"
	driver "github.com/antha-lang/antha/driver/antha_driver_v1"
	runner "github.com/antha-lang/antha/driver/antha_runner_v1"
	lhclient "github.com/antha-lang/antha/driver/lh"
	"github.com/antha-lang/antha/driver/pb/lh"
	"github.com/antha-lang/antha/target/handler"
	"github.com/antha-lang/antha/target/human"
	"github.com/antha-lang/antha/target/mixer"
	"github.com/antha-lang/antha/target/shakerincubator"
	"google.golang.org/grpc"
)

// Common state for tryers
type tryer struct {
	Auto      *Auto
	MaybeArgs []interface{}
	HumanOpt  human.Opt
}

// XXX ADD STUFF HERE
func (a *tryer) Driver(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {
	c := driver.NewDriverClient(conn)
	reply, err := c.DriverType(ctx, &driver.TypeRequest{})
	if err != nil {
		return err
	}

	switch reply.Type {

	case "antha.runner.v1.Runner":
		r := runner.NewRunnerClient(conn)
		reply, err := r.SupportedRunTypes(ctx, &runner.SupportedRunTypesRequest{})
		if err != nil {
			return err
		}
		for _, typ := range reply.Types {
			a.Auto.runners[typ] = append(a.Auto.runners[typ], r)
		}

	case "antha.shakerincubator.v1.ShakerIncubator":
		s := &shakerincubator.ShakerIncubator{}
		a.HumanOpt.CanIncubate = false
		a.Auto.handler[s] = conn
		return a.Auto.Target.AddDevice(s)

	default:
		h := handler.New(
			[]ast.NameValue{
				ast.NameValue{
					Name:  "antha.driver.v1.TypeReply.type",
					Value: reply.Type,
				},
			},
		)
		a.HumanOpt.CanHandle = false
		a.Auto.handler[h] = conn
		return a.Auto.Target.AddDevice(h)
	}

	return nil
}

func (a *tryer) Mixer(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {
	c := lh.NewExtendedLiquidhandlingDriverClient(conn)

	var candidates []interface{}
	candidates = append(candidates, arg)
	candidates = append(candidates, a.MaybeArgs...)

	a.HumanOpt.CanMix = false
	d, err := mixer.New(getMixerOpt(candidates), &lhclient.Driver{C: c})
	if err != nil {
		return err
	}
	return a.Auto.Target.AddDevice(d)
}

func getMixerOpt(maybeArgs []interface{}) (ret mixer.Opt) {
	for _, v := range maybeArgs {
		if o, ok := v.(mixer.Opt); ok {
			return o
		}
	}
	return
}

func (a *tryer) Try(ctx context.Context, conn *grpc.ClientConn, arg interface{}) error {
	var tries []func(context.Context, *grpc.ClientConn, interface{}) error
	tries = append(tries, a.Driver, a.Mixer)

	for _, t := range tries {
		if err := t(ctx, conn, arg); err == nil {
			return nil
		}
	}

	return errNoMatch
}
