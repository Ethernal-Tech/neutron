package wasmbinding

import (
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	icatypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/types"
)

// GetInterchainQueryResult is a function, not method, so the message_plugin can use it
func (qp *QueryPlugin) GetInterchainQueryResult(ctx sdk.Context, queryID uint64) (string, error) {
	// TODO
	return "", nil
}

func (qp *QueryPlugin) GetInterchainAccountAddress(ctx sdk.Context, ownerAddress, connectionId string) (string, error) {
	if ownerAddress == "" || connectionId == "" {
		return "", wasmvmtypes.InvalidRequest{Err: "invalid request params for interchain account address"}
	}

	portID, err := icatypes.NewControllerPortID(ownerAddress)
	if err != nil {
		return "", sdkerrors.Wrapf(err, "could not find account for ownerAddress=%s: %v", ownerAddress, err)
	}

	addr, found := qp.ICAControllerKeeper.GetInterchainAccountAddress(ctx, connectionId, portID)
	if !found {
		return "", sdkerrors.Wrapf(err, "no account found for portID=%s: %v", portID, err)
	}

	return addr, nil
}