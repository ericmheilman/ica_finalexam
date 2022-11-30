package cli

import (
	"github.com/b9lab/toll-road/x/tollroad/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cast"
	"github.com/spf13/cobra"
)

func CmdCreateUserVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-user-vault [index] [owner] [road-operator-index] [token] [balance]",
		Short: "Create a new userVault",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes
			// Get value arguments
			argOwner := args[1]
			argRoadOperatorIndex := args[2]
			argBalance, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateUserVault(
				clientCtx.GetFromAddress().String(),
				argOwner,
				argRoadOperatorIndex,
				argBalance,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateUserVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-user-vault [index] [owner] [road-operator-index] [token] [balance]",
		Short: "Update a userVault",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			// Get indexes

			// Get value arguments
			argOwner := args[1]
			argRoadOperatorIndex := args[2]
			argBalance, err := cast.ToUint64E(args[4])
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateUserVault(
				clientCtx.GetFromAddress().String(),
				argOwner,
				argRoadOperatorIndex,
				argBalance,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteUserVault() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-user-vault [index]",
		Short: "Delete a userVault",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			indexIndex := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteUserVault(
				clientCtx.GetFromAddress().String(),
				indexIndex,
			)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
