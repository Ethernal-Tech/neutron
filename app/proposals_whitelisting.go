package app

import (
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	"github.com/cosmos/cosmos-sdk/x/params/types/proposal"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
	icahosttypes "github.com/cosmos/ibc-go/v3/modules/apps/27-interchain-accounts/host/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v3/modules/apps/transfer/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
)

func IsConsumerProposalWhitelisted(content govtypes.Content) bool {

	switch c := content.(type) {
	case *proposal.ParameterChangeProposal:
		return isConsumerParamChangeWhitelisted(c.Changes)
	case *govtypes.TextProposal,
		*upgradetypes.SoftwareUpgradeProposal,
		*upgradetypes.CancelSoftwareUpgradeProposal,
		*ibcclienttypes.ClientUpdateProposal,
		*ibcclienttypes.UpgradeProposal:
		return true

	default:
		return false
	}

}

func isConsumerParamChangeWhitelisted(paramChanges []proposal.ParamChange) bool {
	for _, paramChange := range paramChanges {
		_, found := WhitelistedParams[paramChangeKey{Subspace: paramChange.Subspace, Key: paramChange.Key}]
		if !found {
			return false
		}
	}
	return true
}

type paramChangeKey struct {
	Subspace, Key string
}

var WhitelistedParams = map[paramChangeKey]struct{}{
	//bank
	{Subspace: banktypes.ModuleName, Key: "SendEnabled"}: {},
	//staking
	{Subspace: stakingtypes.ModuleName, Key: "UnbondingTime"}:     {},
	{Subspace: stakingtypes.ModuleName, Key: "MaxValidators"}:     {},
	{Subspace: stakingtypes.ModuleName, Key: "MaxEntries"}:        {},
	{Subspace: stakingtypes.ModuleName, Key: "HistoricalEntries"}: {},
	{Subspace: stakingtypes.ModuleName, Key: "BondDenom"}:         {},
	//distribution
	{Subspace: distrtypes.ModuleName, Key: "communitytax"}:        {},
	{Subspace: distrtypes.ModuleName, Key: "baseproposerreward"}:  {},
	{Subspace: distrtypes.ModuleName, Key: "bonusproposerreward"}: {},
	{Subspace: distrtypes.ModuleName, Key: "withdrawaddrenabled"}: {},
	//mint
	{Subspace: minttypes.ModuleName, Key: "MintDenom"}:           {},
	{Subspace: minttypes.ModuleName, Key: "InflationRateChange"}: {},
	{Subspace: minttypes.ModuleName, Key: "InflationMax"}:        {},
	{Subspace: minttypes.ModuleName, Key: "InflationMin"}:        {},
	{Subspace: minttypes.ModuleName, Key: "GoalBonded"}:          {},
	{Subspace: minttypes.ModuleName, Key: "BlocksPerYear"}:       {},
	//ibc transfer
	{Subspace: ibctransfertypes.ModuleName, Key: "SendEnabled"}:    {},
	{Subspace: ibctransfertypes.ModuleName, Key: "ReceiveEnabled"}: {},
	//ica
	{Subspace: icahosttypes.SubModuleName, Key: "HostEnabled"}:   {},
	{Subspace: icahosttypes.SubModuleName, Key: "AllowMessages"}: {},
}
