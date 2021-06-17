package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	abci "github.com/tendermint/tendermint/abci/types"
)

type SendCapability struct {
	// SpendLimit specifies the maximum amount of tokens that can be spent
	// by this capability and will be updated as tokens are spent. If it is
	// empty, there is no spend limit and any amount of coins can be spent.
	SpendLimit sdk.Coins
}

var _ Capability = SendCapability{}

func (cap SendCapability) MsgType() sdk.Msg {
	return bank.MsgSend{}
}

func (cap SendCapability) Accept(msg sdk.Msg, block abci.Header) (allow bool, updated Capability, delete bool) {
	switch msg := msg.(type) {
	case bank.MsgSend:
		left, invalid := cap.SpendLimit.SafeSub(msg.Amount)
		if invalid {
			return false, nil, false
		}
		if left.IsZero() {
			return true, nil, true
		}
		return true, SendCapability{SpendLimit: left}, false
	}
	return false, nil, false
}
