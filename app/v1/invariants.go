package v1

import (
	"fmt"
	"time"

	sdk "github.com/orientwalt/htdf/types"
)

func (p *ProtocolV1) assertRuntimeInvariants(ctx sdk.Context) {
	start := time.Now()
	invarRoutes := p.crisisKeeper.Routes()
	for _, ir := range invarRoutes {
		if err := ir.Invar(ctx); err != nil {
			panic(fmt.Errorf("invariant broken: %s\n"+
				"\tCRITICAL please submit the following transaction:\n"+
				"\t\t hscli tx crisis invariant-broken %v %v", err, ir.ModuleName, ir.Route))
		}
	}
	end := time.Now()
	diff := end.Sub(start)
	// ctx.WithLogger(ctx.Logger().With("module", "invariants").Info(
	// 	"Asserted all invariants", "duration", diff, "height", ctx.BlockHeight()))
	ctx.Logger().With("module", "invariants").Info(
		"Asserted all invariants", "duration", diff, "height", ctx.BlockHeight())
}
