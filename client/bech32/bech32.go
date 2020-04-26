package bech32

import (
	"fmt"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/spf13/cobra"
)

func Bech32Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bech32",
		Short: "convert address between bech32 and hex",
		Long:  `convert address between bech32 and hex`,
	}
	cmd.AddCommand(
		cmdBech2Hex(),
		cmdHex2Bech(),
	)

	return cmd
}

func cmdBech2Hex() *cobra.Command {
	return &cobra.Command{
		Use: "b2h [bech32addr]",
		//Aliases: []string{"b2h"},
		Short: "convert bech32 to hex-20",
		Long:  "hscli bech32 b2h htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bech32Addr := args[0]
			hexAddr, err := sdk.AccAddressFromBech32(bech32Addr)
			if err != nil {
				fmt.Printf("AccAddressFromBech32 error|err=%s\n", err)
				return err
			}

			fmt.Printf("bechAddress=%s|hexAddress=%x\n", bech32Addr, hexAddr)
			return nil
		},
	}
}

func cmdHex2Bech() *cobra.Command {
	return &cobra.Command{
		Use: "h2b [Hex-20 address]",
		//Aliases: []string{"b2h"},
		Short: "convert hex-20 to bech32",
		Long:  "hscli bech32 h2b 85CED8DDF399D75C9E381E01F3BDDCEFB9132FE9",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hexAddr := args[0]
			bech32Addr, err := sdk.AccAddressFromHex(hexAddr)
			if err != nil {
				fmt.Printf("AccAddressFromHex error|err=%s\n", err)
				return err
			}

			fmt.Printf("hexAddr=%s|bech32Addr=%s\n", hexAddr, bech32Addr.String())
			return nil
		},
	}
}
