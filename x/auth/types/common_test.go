package types_test

import (
	"github.com/orientwalt/htdf/simapp"
)

var (
	app                   = simapp.Setup(false)
	appCodec, legacyAmino = simapp.MakeCodecs()
)
