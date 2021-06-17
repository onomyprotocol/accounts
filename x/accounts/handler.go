package accounts

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/onomyprotocol/accounts/x/accounts/keeper"
	"github.com/onomyprotocol/accounts/x/accounts/types"
)

// NewHandler ...
func NewHandler(k keeper.Keeper) sdk.Handler {
	// this line is used by starport scaffolding # handler/msgServer

	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgAuthorize:
			k.Authorize(ctx, msg.Operator, msg.Authorizer, msg.Capability, msg.Expiration)
			return sdk.Result{}
		case types.MsgExecAuthorizedAction:
			return k.DispatchActions(ctx, msg.Signer, msg.Msgs)
		case types.MsgRevoke:
			k.Revoke(ctx, msg.Operator, msg.Authorizer, msg.MsgType)
			return sdk.Result{}
		case types.MsgAuthorizeFeeAllowance:
			k.AuthorizeFeeAllowance(ctx, msg.Operator, msg.Authorizer, msg.Allowance)
			return sdk.Result{}
		case types.MsgRevokeFeeAllowance:
			k.RevokeFeeAllowance(ctx, msg.Operator, msg.Authorizer)
			return sdk.Result{}
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}
