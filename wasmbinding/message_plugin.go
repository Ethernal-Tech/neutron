package wasmbinding

import (
	"encoding/json"

	wasmkeeper "github.com/CosmWasm/wasmd/x/wasm/keeper"
	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/neutron-org/neutron/wasmbinding/bindings"
	icqkeeper "github.com/neutron-org/neutron/x/interchainqueries/keeper"
	icqtypes "github.com/neutron-org/neutron/x/interchainqueries/types"
	ictxkeeper "github.com/neutron-org/neutron/x/interchaintxs/keeper"
	ictxtypes "github.com/neutron-org/neutron/x/interchaintxs/types"
)

func CustomMessageDecorator(ictx *ictxkeeper.Keeper, icq *icqkeeper.Keeper) func(messenger wasmkeeper.Messenger) wasmkeeper.Messenger {
	return func(old wasmkeeper.Messenger) wasmkeeper.Messenger {
		return &CustomMessenger{
			Keeper:        *ictx,
			Wrapped:       old,
			Ictxmsgserver: ictxkeeper.NewMsgServerImpl(*ictx),
			Icqmsgserver:  icqkeeper.NewMsgServerImpl(*icq),
		}
	}
}

type CustomMessenger struct {
	Keeper        ictxkeeper.Keeper
	Wrapped       wasmkeeper.Messenger
	Ictxmsgserver ictxtypes.MsgServer
	Icqmsgserver  icqtypes.MsgServer
}

var _ wasmkeeper.Messenger = (*CustomMessenger)(nil)

func (m *CustomMessenger) DispatchMsg(ctx sdk.Context, contractAddr sdk.AccAddress, contractIBCPortID string, msg wasmvmtypes.CosmosMsg) ([]sdk.Event, [][]byte, error) {
	if msg.Custom != nil {
		var contractMsg bindings.NeutronMsg
		if err := json.Unmarshal(msg.Custom, &contractMsg); err != nil {
			ctx.Logger().Debug("json.Unmarshal: failed to decode incoming custom cosmos message",
				"from_address", contractAddr.String(),
				"message", string(msg.Custom),
				"error", err,
			)
			return nil, nil, sdkerrors.Wrap(err, "failed to decode incoming custom cosmos message")
		}

		if contractMsg.SubmitTx != nil {
			return m.submitTx(ctx, contractAddr, contractMsg.SubmitTx)
		}
		if contractMsg.RegisterInterchainAccount != nil {
			return m.registerInterchainAccount(ctx, contractAddr, contractMsg.RegisterInterchainAccount)
		}
		if contractMsg.RegisterInterchainQuery != nil {
			return m.registerInterchainQuery(ctx, contractAddr, contractMsg.RegisterInterchainQuery)
		}
	}

	return m.Wrapped.DispatchMsg(ctx, contractAddr, contractIBCPortID, msg)
}

func (m *CustomMessenger) submitTx(ctx sdk.Context, contractAddr sdk.AccAddress, submitTx *bindings.SubmitTx) ([]sdk.Event, [][]byte, error) {
	response, err := m.PerformSubmitTx(ctx, contractAddr, submitTx)
	if err != nil {
		ctx.Logger().Debug("PerformSubmitTx: failed to submit interchain transaction",
			"from_address", contractAddr.String(),
			"connection_id", submitTx.ConnectionId,
			"interchain_account_id", submitTx.InterchainAccountId,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to submit interchain transaction")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Debug("json.Marshal: failed to marshal submitTx response to JSON",
			"from_address", contractAddr.String(),
			"connection_id", submitTx.ConnectionId,
			"interchain_account_id", submitTx.InterchainAccountId,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("interchain transaction submitted",
		"from_address", contractAddr.String(),
		"connection_id", submitTx.ConnectionId,
		"interchain_account_id", submitTx.InterchainAccountId,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) PerformSubmitTx(ctx sdk.Context, contractAddr sdk.AccAddress, submitTx *bindings.SubmitTx) (*bindings.SubmitTxResponse, error) {
	tx := ictxtypes.MsgSubmitTx{
		FromAddress:         contractAddr.String(),
		ConnectionId:        submitTx.ConnectionId,
		Memo:                submitTx.Memo,
		InterchainAccountId: submitTx.InterchainAccountId,
	}
	for _, msg := range submitTx.Msgs {
		tx.Msgs = append(tx.Msgs, &types.Any{
			TypeUrl: msg.TypeURL,
			Value:   msg.Value,
		})
	}
	if err := tx.UnpackInterfaces(m.Keeper.Codec); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to unpack interfaces to send interchain transaction")
	}

	if err := tx.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming SubmitTx message")
	}

	response, err := m.Ictxmsgserver.SubmitTx(sdk.WrapSDKContext(ctx), &tx)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to submit interchain transaction")
	}

	return (*bindings.SubmitTxResponse)(response), nil
}

func (m *CustomMessenger) registerInterchainAccount(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainAccount) ([]sdk.Event, [][]byte, error) {
	response, err := m.PerformRegisterInterchainAccount(ctx, contractAddr, reg)
	if err != nil {
		ctx.Logger().Debug("PerformRegisterInterchainAccount: failed to register interchain account",
			"from_address", contractAddr.String(),
			"connection_id", reg.ConnectionId,
			"interchain_account_id", reg.InterchainAccountId,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to register interchain account")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Debug("json.Marshal: failed to marshal register interchain account response to JSON",
			"from_address", contractAddr.String(),
			"connection_id", reg.ConnectionId,
			"interchain_account_id", reg.InterchainAccountId,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("registered interchain account",
		"from_address", contractAddr.String(),
		"connection_id", reg.ConnectionId,
		"interchain_account_id", reg.InterchainAccountId,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) PerformRegisterInterchainAccount(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainAccount) (*bindings.RegisterInterchainAccountResponse, error) {
	msg := ictxtypes.MsgRegisterInterchainAccount{
		FromAddress:         contractAddr.String(),
		ConnectionId:        reg.ConnectionId,
		InterchainAccountId: reg.InterchainAccountId,
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming RegisterInterchainAccount message")
	}

	response, err := m.Ictxmsgserver.RegisterInterchainAccount(sdk.WrapSDKContext(ctx), &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to register interchain account")
	}

	return (*bindings.RegisterInterchainAccountResponse)(response), nil
}

func (m *CustomMessenger) registerInterchainQuery(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainQuery) ([]sdk.Event, [][]byte, error) {
	response, err := m.PerformRegisterInterchainQuery(ctx, contractAddr, reg)
	if err != nil {
		ctx.Logger().Debug("PerformRegisterInterchainQuery: failed to register interchain query",
			"from_address", contractAddr.String(),
			"query_type", reg.QueryType,
			"kv_keys", icqtypes.KVKeys(reg.Keys).String(),
			"transactions_filter", reg.TransactionsFilter,
			"zone_id", reg.ZoneId,
			"connection_id", reg.ConnectionId,
			"update_period", reg.UpdatePeriod,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "failed to register interchain query")
	}

	data, err := json.Marshal(response)
	if err != nil {
		ctx.Logger().Debug("json.Marshal: failed to marshal register interchain query response to JSON",
			"from_address", contractAddr.String(),
			"kv_keys", icqtypes.KVKeys(reg.Keys).String(),
			"transactions_filter", reg.TransactionsFilter,
			"zone_id", reg.ZoneId,
			"connection_id", reg.ConnectionId,
			"update_period", reg.UpdatePeriod,
			"error", err,
		)
		return nil, nil, sdkerrors.Wrap(err, "marshal json failed")
	}

	ctx.Logger().Debug("registered interchain query",
		"from_address", contractAddr.String(),
		"query_type", reg.QueryType,
		"kv_keys", icqtypes.KVKeys(reg.Keys).String(),
		"transactions_filter", reg.TransactionsFilter,
		"zone_id", reg.ZoneId,
		"connection_id", reg.ConnectionId,
		"update_period", reg.UpdatePeriod,
		"query_id", response.Id,
	)
	return nil, [][]byte{data}, nil
}

func (m *CustomMessenger) PerformRegisterInterchainQuery(ctx sdk.Context, contractAddr sdk.AccAddress, reg *bindings.RegisterInterchainQuery) (*bindings.RegisterInterchainQueryResponse, error) {
	msg := icqtypes.MsgRegisterInterchainKVQuery{
		Keys:         reg.Keys,
		ZoneId:       reg.ZoneId,
		ConnectionId: reg.ConnectionId,
		UpdatePeriod: reg.UpdatePeriod,
		Sender:       contractAddr.String(),
	}
	if err := msg.ValidateBasic(); err != nil {
		return nil, sdkerrors.Wrap(err, "failed to validate incoming RegisterInterchainKVQuery message")
	}

	response, err := m.Icqmsgserver.RegisterInterchainKVQuery(sdk.WrapSDKContext(ctx), &msg)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to register interchain query")
	}

	return (*bindings.RegisterInterchainQueryResponse)(response), nil
}
