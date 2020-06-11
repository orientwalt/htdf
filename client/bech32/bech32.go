package bech32

import (
	"fmt"
	sdk "github.com/orientwalt/htdf/types"
	"github.com/spf13/cobra"
)

func Bech32Commands() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bech32",
		Short: "convert accountPrefixAddress(bech32) , validatorPrefixAddress(bech32) and hexAddress",
		Long:  `convert accountPrefixAddress(bech32) , validatorPrefixAddress(bech32) and hexAddress`,
	}
	cmd.AddCommand(
		cmdBech2Hex(),
		cmdHex2Bech(),
		cmdBech2Val(),
		cmdVal2Bech(),
	)

	return cmd
}

func cmdBech2Hex() *cobra.Command {
	return &cobra.Command{
		Use:   "b2h [accountPrefixAddress(bech32)]",
		Short: "convert accountPrefixAddress(bech32) to hex-20",
		Long:  "hscli bech32 b2h htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bech32Addr := args[0]
			hexAddr, err := sdk.AccAddressFromBech32(bech32Addr)
			if err != nil {
				fmt.Printf("AccAddressFromBech32 error|err=%s\n", err)
				return err
			}

			fmt.Printf("accountPrefixAddress=%s|hexAddress=%x\n", bech32Addr, hexAddr)
			return nil
		},
	}
}

func cmdHex2Bech() *cobra.Command {
	return &cobra.Command{
		Use:   "h2b [Hex-20 address]",
		Short: "convert hex-20 to accountPrefixAddress(bech32)",
		Long:  "hscli bech32 h2b 85CED8DDF399D75C9E381E01F3BDDCEFB9132FE9",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hexAddr := args[0]
			bech32Addr, err := sdk.AccAddressFromHex(hexAddr)
			if err != nil {
				fmt.Printf("AccAddressFromHex error|err=%s\n", err)
				return err
			}

			fmt.Printf("hexAddr=%s|accountPrefixAddress=%s\n", hexAddr, bech32Addr.String())
			return nil
		},
	}
}

func cmdBech2Val() *cobra.Command {
	return &cobra.Command{
		Use:   "b2v [accountPrefixAddress(bech32)]",
		Short: "convert accountPrefixAddress(bech32) to validatorPrefixAddress(bech32)",
		Long:  "hscli bech32 b2v htdf1sh8d3h0nn8t4e83crcql80wua7u3xtlfj5dej3",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bech32Addr := args[0]
			hexAddr, err := sdk.AccAddressFromBech32(bech32Addr)
			if err != nil {
				fmt.Printf("AccAddressFromBech32 error|err=%s\n", err)
				return err
			}

			valAddr := sdk.ValAddress(hexAddr)
			fmt.Printf("accountPrefixAddress=%s|validatorPrefixAddress=%s\n", bech32Addr, valAddr.String())
			return nil
		},
	}
}

func cmdVal2Bech() *cobra.Command {
	return &cobra.Command{
		Use:   "v2b [validatorPrefixAddress(bech32)]",
		Short: "convert validatorPrefixAddress(bech32) to accountPrefixAddress(bech32)",
		Long:  "hscli bech32 v2b htdfvaloper12347g0nk9vpae7886xp0sxjdxya27lq4720u04",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			bech32Addr := args[0]
			hexAddr, err := sdk.ValAddressFromBech32(bech32Addr)
			if err != nil {
				fmt.Printf("ValAddressFromBech32 error|err=%s\n", err)
				return err
			}

			accAddr := sdk.AccAddress(hexAddr)
			fmt.Printf("validatorPrefixAddress=%s|AccountPrefixAddress=%s\n", bech32Addr, accAddr.String())
			return nil
		},
	}
}
