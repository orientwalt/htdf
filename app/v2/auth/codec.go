package auth

import (
	"github.com/orientwalt/htdf/codec"
	xauth "github.com/orientwalt/htdf/x/auth"
)


var msgCdc = codec.New()

func init() {
	xauth.RegisterCodec(msgCdc)
	codec.RegisterCrypto(msgCdc)
}
