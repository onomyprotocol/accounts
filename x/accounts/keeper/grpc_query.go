package keeper

import (
	"github.com/onomyprotocol/accounts/x/accounts/types"
)

var _ types.QueryServer = Keeper{}
