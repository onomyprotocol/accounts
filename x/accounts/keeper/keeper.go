package keeper

import (
	"bytes"
	"fmt"
	"time"

	"github.com/onomyprotocol/accounts/x/accounts/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type Keeper struct {
	storeKey sdk.StoreKey
	cdc      *codec.LegacyAmino
	router   sdk.Router
}

type capabilityAuthorize struct {
	Capability types.Capability

	Expiration time.Time
}

func NewKeeper(storeKey sdk.StoreKey, cdc *codec.LegacyAmino, router sdk.Router) Keeper {
	return Keeper{storeKey, cdc, router}
}

func ActorCapabilityKey(operator sdk.AccAddress, authorizer sdk.AccAddress, msg sdk.Msg) []byte {
	return actorCapabilityKey(operator, authorizer, msg.Route(), msg.Type())
}

func actorCapabilityKey(operator sdk.AccAddress, authorizer sdk.AccAddress, route, typ string) []byte {
	return []byte(fmt.Sprintf("c/%x/%x/%s/%s", operator, authorizer, route, typ))
}

func FeeAllowanceKey(operator sdk.AccAddress, authorizer sdk.AccAddress) []byte {
	return []byte(fmt.Sprintf("f/%x/%x", operator, authorizer))
}

// operator sdk.AccAddress, authorizer sdk.AccAddress, msgType sdk.Msg
func (k Keeper) getCapabilityAuthorize(ctx sdk.Context, actor []byte) (grant capabilityAuthorize, found bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(actor)
	if bz == nil {
		return grant, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &grant)
	return grant, true
}

func (k Keeper) Authorize(ctx sdk.Context, operator sdk.AccAddress, authorizer sdk.AccAddress, capability types.Capability, expiration time.Time) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(capabilityAuthorize{capability, expiration})
	actor := ActorCapabilityKey(operator, authorizer, capability.MsgType())
	store.Set(actor, bz)
}

func (k Keeper) update(ctx sdk.Context, operator sdk.AccAddress, authorizer sdk.AccAddress, updated types.Capability) {
	grant, found := k.getCapabilityAuthorize(ctx, ActorCapabilityKey(operator, authorizer, updated.MsgType()))
	if !found {
		return
	}
	grant.Capability = updated
}

func (k Keeper) Revoke(ctx sdk.Context, operator sdk.AccAddress, authorizer sdk.AccAddress, msgType sdk.Msg) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(ActorCapabilityKey(operator, authorizer, msgType))
}

func (k Keeper) GetCapability(ctx sdk.Context, operator sdk.AccAddress, authorizer sdk.AccAddress, msgType sdk.Msg) types.Capability {
	grant, found := k.getCapabilityAuthorize(ctx, ActorCapabilityKey(operator, authorizer, msgType))
	if !found {
		return nil
	}
	if !grant.Expiration.IsZero() && grant.Expiration.Before(ctx.BlockHeader().Time) {
		k.Revoke(ctx, operator, authorizer, msgType)
		return nil
	}
	return grant.Capability
}

func (k Keeper) DispatchActions(ctx sdk.Context, sender sdk.AccAddress, msgs []sdk.Msg) (*sdk.Result, error) {
	var res sdk.Result
	for _, msg := range msgs {
		signers := msg.GetSigners()
		if len(signers) != 1 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "can only dispatch a authorized msg with 1 signer")
		}
		actor := signers[0]
		if !bytes.Equal(actor, sender) {
			capability := k.GetCapability(ctx, sender, actor, msg)
			if capability == nil {
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
			}
			allow, updated, del := capability.Accept(msg, ctx.BlockHeader())
			if !allow {
				return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "unauthorized")
			}
			if del {
				k.Revoke(ctx, sender, actor, msg)
			} else if updated != nil {
				k.update(ctx, sender, actor, updated)
			}
		}
		res, err := k.router.Route(ctx, msg.Route())(ctx, msg)
		if err != nil {
			return res, err
		}
	}
	return &res, nil
}

func (k Keeper) AuthorizeFeeAllowance(ctx sdk.Context, operator sdk.AccAddress, authorizer sdk.AccAddress, allowance types.FeeAllowance) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryBare(allowance)
	store.Set(FeeAllowanceKey(operator, authorizer), bz)
}

func (k Keeper) RevokeFeeAllowance(ctx sdk.Context, operator sdk.AccAddress, authorizer sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(FeeAllowanceKey(operator, authorizer))
}

type FeeAllowanceAuthorize struct {
	Allowance  types.FeeAllowance `json:"allowance"`
	Operator   sdk.AccAddress     `json:"operator"`
	Authorizer sdk.AccAddress     `json:"authorizer"`
}

func (k Keeper) GetFeeAllowances(ctx sdk.Context, operator sdk.AccAddress) []FeeAllowanceAuthorize {
	prefix := fmt.Sprintf("g/%x/", operator)
	prefixBytes := []byte(prefix)
	store := ctx.KVStore(k.storeKey)
	var grants []FeeAllowanceAuthorize
	iter := sdk.KVStorePrefixIterator(store, prefixBytes)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		authorizer, _ := sdk.AccAddressFromHex(string(iter.Key()[len(prefix):]))
		bz := iter.Value()
		var allowance types.FeeAllowance
		k.cdc.MustUnmarshalBinaryBare(bz, &allowance)
		grants = append(grants, FeeAllowanceAuthorize{
			Allowance:  allowance,
			Operator:   operator,
			Authorizer: authorizer,
		})
	}
	return grants
}

func (k Keeper) AllowAuthorizedFees(ctx sdk.Context, operator sdk.AccAddress, authorizer sdk.AccAddress, fee sdk.Coins) bool {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(FeeAllowanceKey(operator, authorizer))
	if len(bz) == 0 {
		return false
	}
	var allowance types.FeeAllowance
	k.cdc.MustUnmarshalBinaryBare(bz, &allowance)
	if allowance == nil {
		return false
	}
	allow, updated, delete := allowance.Accept(fee, ctx.BlockHeader())
	if allow == false {
		return false
	}
	if delete {
		k.RevokeFeeAllowance(ctx, operator, authorizer)
	} else if updated != nil {
		k.AuthorizeFeeAllowance(ctx, operator, authorizer, updated)
	}
	return true
}

func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
