package keeper

import (
	"context"

	"github.com/b9lab/toll-road/x/tollroad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateUserVault(goCtx context.Context, msg *types.MsgCreateUserVault) (*types.MsgCreateUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value already exists
	_, isFound := k.GetUserVault(
		ctx,
		msg.Creator,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if isFound {
		return nil, types.ErrIndexSet
	}

	var userVault = types.UserVault{

		//	Creator:           msg.Creator,
		Owner: msg.Creator,
		// 		Owner:             msg.Owner,

		RoadOperatorIndex: msg.RoadOperatorIndex,
		Token:             msg.Token,

		Balance: msg.Balance,
	}

	k.SetUserVault(
		ctx,
		userVault,
	)

	if msg.Balance == 0 {
		return nil, types.ErrZeroTokens
	}

	creatorAddress, _ := sdk.AccAddressFromBech32(msg.Creator)
	err := k.bank.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(creatorAddress), types.ModuleName, sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(msg.Balance))))

	if err != nil {
		return nil, err
	}

	return &types.MsgCreateUserVaultResponse{}, nil
}

func (k msgServer) UpdateUserVault(goCtx context.Context, msg *types.MsgUpdateUserVault) (*types.MsgUpdateUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetUserVault(
		ctx,
		msg.Creator,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// if balance = 0
	if msg.Balance == 0 {
		return nil, types.ErrZeroTokens
	}

	//creatorAddress, _ := sdk.AccAddressFromBech32(msg.Creator)

	// if balance is lower
	if msg.Balance < valFound.Balance {
		difference := valFound.Balance - msg.Balance
		creatorAddress, _ := sdk.AccAddressFromBech32(msg.Creator)
		err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(creatorAddress), sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(difference))))
		if err != nil {
			panic("bank error")
		}
	}

	// if balance is higher
	if msg.Balance > valFound.Balance {
		difference := msg.Balance - valFound.Balance
		creatorAddress, _ := sdk.AccAddressFromBech32(msg.Creator)
		err := k.bank.SendCoinsFromAccountToModule(ctx, sdk.AccAddress(creatorAddress), types.ModuleName, sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(difference))))
		if err != nil {
			return nil, err
		}
	}

	var userVault = types.UserVault{
		Owner:             msg.Creator,
		RoadOperatorIndex: msg.RoadOperatorIndex,
		Token:             msg.Token,
		Balance:           msg.Balance,
	}

	k.SetUserVault(ctx, userVault)

	return &types.MsgUpdateUserVaultResponse{}, nil
}

func (k msgServer) DeleteUserVault(goCtx context.Context, msg *types.MsgDeleteUserVault) (*types.MsgDeleteUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Check if the value exists
	valFound, isFound := k.GetUserVault(
		ctx,
		msg.Creator,
		msg.RoadOperatorIndex,
		msg.Token,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Owner {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveUserVault(
		ctx,
		msg.Creator,
		msg.RoadOperatorIndex,
		msg.Token,
	)

	creatorAddress, _ := sdk.AccAddressFromBech32(msg.Creator)
	err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(creatorAddress), sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(valFound.Balance))))

	if err != nil {
		panic("bank error")
	}

	return &types.MsgDeleteUserVaultResponse{}, nil
}
