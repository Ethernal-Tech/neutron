package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/interchainqueries module sentinel errors
var (
	ErrInvalidQueryID            = sdkerrors.Register(ModuleName, 1100, "invalid query id")
	ErrEmptyResult               = sdkerrors.Register(ModuleName, 1101, "empty result")
	ErrInvalidClientID           = sdkerrors.Register(ModuleName, 1102, "invalid client id")
	ErrInvalidUpdatePeriod       = sdkerrors.Register(ModuleName, 1103, "invalid update period")
	ErrInvalidConnectionID       = sdkerrors.Register(ModuleName, 1104, "invalid connection id")
	ErrInvalidZoneID             = sdkerrors.Register(ModuleName, 1105, "invalid zone id")
	ErrInvalidQueryType          = sdkerrors.Register(ModuleName, 1106, "invalid query type")
	ErrInvalidSubmittedResult    = sdkerrors.Register(ModuleName, 1107, "invalid result")
	ErrProtoMarshal              = sdkerrors.Register(ModuleName, 1108, "failed to marshal protobuf bytes")
	ErrProtoUnmarshal            = sdkerrors.Register(ModuleName, 1109, "failed to unmarshal protobuf bytes")
	ErrInvalidType               = sdkerrors.Register(ModuleName, 1110, "invalid type")
	ErrInternal                  = sdkerrors.Register(ModuleName, 1111, "internal error")
	ErrInvalidProof              = sdkerrors.Register(ModuleName, 1112, "merkle proof is invalid")
	ErrInvalidHeader             = sdkerrors.Register(ModuleName, 1113, "header is invalid")
	ErrInvalidHeight             = sdkerrors.Register(ModuleName, 1114, "height is invalid")
	ErrNoQueryResult             = sdkerrors.Register(ModuleName, 1115, "no query result")
	ErrInvalidTransactionsFilter = sdkerrors.Register(ModuleName, 1116, "invalid transactions filter")
)
