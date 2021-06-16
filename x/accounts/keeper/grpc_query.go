package keeper

import (
	"github.com/user/accounts/x/accounts/types"
)

var _ types.QueryServer = Keeper{}
