package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/ibc-go/v4/modules/apps/transfer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (msg *MsgTransfer) ValidateBasic() error {
	if msg.PayerFee == nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInsufficientFee, "fee can't be nil")
	}

	sdkMsg := types.NewMsgTransfer(msg.SourcePort, msg.SourceChannel, msg.Token, msg.Sender, msg.Receiver, msg.TimeoutHeight, msg.TimeoutTimestamp)
	return sdkMsg.ValidateBasic()
}

func (msg *MsgTransfer) GetSigners() []sdk.AccAddress {
	fromAddress, _ := sdk.AccAddressFromBech32(msg.Sender)
	return []sdk.AccAddress{fromAddress}
}

func (msg *MsgTransfer) Route() string {
	return types.RouterKey
}

func (msg *MsgTransfer) Type() string {
	return types.TypeMsgTransfer
}

func (msg MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(&msg))
}