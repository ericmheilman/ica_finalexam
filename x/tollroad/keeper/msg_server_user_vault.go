package keeper

import (
	"context"
	"strconv"

	"github.com/b9lab/toll-road/x/tollroad/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func (k msgServer) CreateUserVault(goCtx context.Context, msg *types.MsgCreateUserVault) (*types.MsgCreateUserVaultResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	// Check if the value already exists
	_, isFound := k.GetUserVault(
		ctx,
		msg.RoadOperatorIndex,
	)
	if isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "index already set")
	}

	var userVault = types.UserVault{
		Creator:           msg.Creator,
		Index:             msg.RoadOperatorIndex,
		RoadOperatorIndex: msg.RoadOperatorIndex,
		Token:             msg.Token,
		Balance:           msg.Balance,
	}

	//balanceStr := strconv.Itoa(int(msg.Balance))
	//balanceInt, _ := strconv.Atoi(balanceStr)

	k.SetUserVault(
		ctx,
		userVault,
	)
	creatorAddress, _ := sdk.AccAddressFromBech32(msg.Creator)

	k.Keeper.UserVault(goCtx, nil)

	err := k.bank.SendCoinsFromAccountToModule(ctx, creatorAddress, types.ModuleName, sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(msg.Balance))))

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
		msg.RoadOperatorIndex,
	)

	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	//k.GetRoadOperator()
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	systemInfo, found := k.Keeper.GetSystemInfo(ctx)
	if !found {
		panic("SystemInfo not found")
	}

	newIndex := strconv.FormatUint(systemInfo.NextOperatorId, 10)

	var userVault = types.UserVault{
		Creator:           msg.Creator,
		Index:             msg.RoadOperatorIndex,
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
		msg.RoadOperatorIndex,
	)
	if !isFound {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, "index not set")
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.Creator != valFound.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemoveUserVault(
		ctx,
		msg.RoadOperatorIndex,
	)

	err := k.bank.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.AccAddress(msg.Creator), sdk.NewCoins(sdk.NewCoin(msg.Token, sdk.NewIntFromUint64(valFound.Balance))))

	if err != nil {
		return nil, err
	}

	/*
		if err != nil {
			panic(err)
		}

	*/
	return &types.MsgDeleteUserVaultResponse{}, nil
}
