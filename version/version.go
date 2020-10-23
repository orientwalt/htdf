package version

import (
	"strconv"
)

// Version - Htdf Version
const ProtocolVersion = 0
const Version = "0.14.1"

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
// func ServeVersionCommand(cdc *codec.Codec) *cobra.Command {
// 	cmd := &cobra.Command{
// 		Use:   "version",
// 		Short: "Show executable binary version",
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			fmt.Println(GetVersion())
// 			return nil
// 		},
// 	}
// 	return cmd
// }
