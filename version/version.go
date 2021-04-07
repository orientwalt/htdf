package version

import (
	"fmt"
	"strconv"

	"github.com/orientwalt/htdf/codec"
	"github.com/orientwalt/htdf/params"
	"github.com/spf13/cobra"
)

// DO NOT EDIT THIS AppVersion
const AppVersion = 0

//-------------------------------------------
// ProtocolVersion - protocol version of (software)upgrade
const ProtocolVersion = 1 // start from version 1 by yqq 2021-04-07
var Version = params.Version

// GitCommit set by build flags
var GitCommit = ""

// return version of CLI/node and commit hash
func GetVersion() string {
	v := Version
	if GitCommit != "" {
		v = v + "-" + GitCommit + "-" + strconv.Itoa(ProtocolVersion)
	}
	return v
}

// ServeVersionCommand
func ServeVersionCommand(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Show executable binary version",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(GetVersion())
			return nil
		},
	}
	return cmd
}
