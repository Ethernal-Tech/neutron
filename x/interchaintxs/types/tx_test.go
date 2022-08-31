package types_test

import (
	"testing"

	cosmosTypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdktypes "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/neutron-org/neutron/x/interchaintxs/types"
)

const TestAddress = "cosmos10h9stc5v6ntgeygf5xf945njqq5h32r53uquvw"

func TestMsgRegisterInterchainAccountValidate(t *testing.T) {
	tests := []struct {
		name        string
		malleate    func() sdktypes.Msg
		expectedErr error
	}{
		{
			"valid",
			func() sdktypes.Msg {
				return &types.MsgRegisterInterchainAccount{
					FromAddress:         TestAddress,
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
				}
			},
			nil,
		},
		{
			"empty fromAddress",
			func() sdktypes.Msg {
				return &types.MsgRegisterInterchainAccount{
					FromAddress:         "",
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"invalid fromAddress",
			func() sdktypes.Msg {
				return &types.MsgRegisterInterchainAccount{
					FromAddress:         "invalid address",
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"empty connection id",
			func() sdktypes.Msg {
				return &types.MsgRegisterInterchainAccount{
					FromAddress:         TestAddress,
					ConnectionId:        "",
					InterchainAccountId: "1",
				}
			},
			types.ErrEmptyConnectionID,
		},
		{
			"empty interchain account",
			func() sdktypes.Msg {
				return &types.MsgRegisterInterchainAccount{
					FromAddress:         TestAddress,
					ConnectionId:        "connection-id",
					InterchainAccountId: "",
				}
			},
			types.ErrEmptyInterchainAccountID,
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()

		if tt.expectedErr != nil {
			require.ErrorIs(t, msg.ValidateBasic(), tt.expectedErr)
		} else {
			require.NoError(t, msg.ValidateBasic())
		}
	}
}

func TestMsgRegisterInterchainAccountGetSigners(t *testing.T) {
	tests := []struct {
		name     string
		malleate func() sdktypes.Msg
	}{
		{
			"valid_signer",
			func() sdktypes.Msg {
				return &types.MsgRegisterInterchainAccount{
					FromAddress:         TestAddress,
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
				}
			},
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()
		addr, _ := sdktypes.AccAddressFromBech32(TestAddress)
		require.Equal(t, msg.GetSigners(), []sdktypes.AccAddress{addr})
	}
}

func TestMsgSubmitTXValidate(t *testing.T) {
	tests := []struct {
		name        string
		malleate    func() sdktypes.Msg
		expectedErr error
	}{
		{
			"valid",
			func() sdktypes.Msg {
				return &types.MsgSubmitTx{
					FromAddress:         TestAddress,
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
					Msgs: []*cosmosTypes.Any{{
						TypeUrl: "msg",
						Value:   []byte{100}, // just check that values are not nil
					}},
				}
			},
			nil,
		},
		{
			"empty connection id",
			func() sdktypes.Msg {
				return &types.MsgSubmitTx{
					FromAddress:         TestAddress,
					ConnectionId:        "",
					InterchainAccountId: "1",
					Msgs: []*cosmosTypes.Any{{
						TypeUrl: "msg",
						Value:   []byte{100}, // just check that values are not nil
					}},
				}
			},
			types.ErrEmptyConnectionID,
		},
		{
			"empty interchain account id",
			func() sdktypes.Msg {
				return &types.MsgSubmitTx{
					FromAddress:         TestAddress,
					ConnectionId:        "connection-id",
					InterchainAccountId: "",
					Msgs: []*cosmosTypes.Any{{
						TypeUrl: "msg",
						Value:   []byte{100}, // just check that values are not nil
					}},
				}
			},
			types.ErrEmptyInterchainAccountID,
		},
		{
			"no messages",
			func() sdktypes.Msg {
				return &types.MsgSubmitTx{
					FromAddress:         TestAddress,
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
					Msgs:                nil,
				}
			},
			types.ErrNoMessages,
		},
		{
			"empty FromAddress",
			func() sdktypes.Msg {
				return &types.MsgSubmitTx{
					FromAddress:         "",
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
					Msgs: []*cosmosTypes.Any{{
						TypeUrl: "msg",
						Value:   []byte{100}, // just check that values are not nil
					}},
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
		{
			"invalid FromAddress",
			func() sdktypes.Msg {
				return &types.MsgSubmitTx{
					FromAddress:         "invalid_address",
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
					Msgs: []*cosmosTypes.Any{{
						TypeUrl: "msg",
						Value:   []byte{100}, // just check that values are not nil
					}},
				}
			},
			sdkerrors.ErrInvalidAddress,
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()

		if tt.expectedErr != nil {
			require.ErrorIs(t, msg.ValidateBasic(), tt.expectedErr)
		} else {
			require.NoError(t, msg.ValidateBasic())
		}
	}
}

func TestMsgSubmitTXGetSigners(t *testing.T) {
	tests := []struct {
		name     string
		malleate func() sdktypes.Msg
	}{
		{
			"valid_signer",
			func() sdktypes.Msg {
				return &types.MsgSubmitTx{
					FromAddress:         TestAddress,
					ConnectionId:        "connection-id",
					InterchainAccountId: "1",
					Msgs: []*cosmosTypes.Any{{
						TypeUrl: "msg",
						Value:   []byte{100}, // just check that values are not nil
					}},
				}
			},
		},
	}

	for _, tt := range tests {
		msg := tt.malleate()
		addr, _ := sdktypes.AccAddressFromBech32(TestAddress)
		require.Equal(t, msg.GetSigners(), []sdktypes.AccAddress{addr})
	}
}