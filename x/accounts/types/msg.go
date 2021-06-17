package types

import (
	"encoding/json"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgExecAuthorizedAction struct {
	Signer sdk.AccAddress `json:"signer"`
	Msgs   []sdk.Msg      `json:"msg"`
}

func (msg MsgExecAuthorizedAction) Route() string {
	return "delegation"
}

func (msg MsgExecAuthorizedAction) Type() string {
	return "exec_authorized"
}

func (msg MsgExecAuthorizedAction) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgExecAuthorizedAction) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgExecAuthorizedAction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Signer}
}

type MsgAuthorize struct {
	Authorizer sdk.AccAddress `json:"authorizer"`
	Operator   sdk.AccAddress `json:"operator"`
	Capability Capability     `json:"capability"`
	Expiration time.Time      `json:"expiration"`
}

func NewMsgAuthorize(authorizer sdk.AccAddress, operator sdk.AccAddress, capability Capability, expiration time.Time) MsgAuthorize {
	return MsgAuthorize{Authorizer: authorizer, Operator: operator, Capability: capability, Expiration: expiration}
}

func (msg MsgAuthorize) Route() string {
	return "delegation"
}

func (msg MsgAuthorize) Type() string {
	return "authorize"
}

func (msg MsgAuthorize) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgAuthorize) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgAuthorize) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authorizer}
}

type MsgRevoke struct {
	Authorizer sdk.AccAddress `json:"authorizer"`
	Operator   sdk.AccAddress `json:"operator"`
	MsgType    sdk.Msg        `json:"msg_type"`
}

func (msg MsgRevoke) Route() string {
	return "delegation"
}

func (msg MsgRevoke) Type() string {
	return "revoke"
}

func (msg MsgRevoke) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgRevoke) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgRevoke) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authorizer}
}

type MsgAuthorizeFeeAllowance struct {
	Authorizer sdk.AccAddress `json:"authorizer"`
	Operator   sdk.AccAddress `json:"operator"`
	Allowance  FeeAllowance   `json:"allowance"`
}

func NewMsgAuthorizeFeeAllowance(authorizer sdk.AccAddress, operator sdk.AccAddress, allowance FeeAllowance) MsgAuthorizeFeeAllowance {
	return MsgAuthorizeFeeAllowance{Authorizer: authorizer, Operator: operator, Allowance: allowance}
}

func (msg MsgAuthorizeFeeAllowance) Route() string {
	return "delegation"
}

func (msg MsgAuthorizeFeeAllowance) Type() string {
	return "authorize-fee-allowance"
}

func (msg MsgAuthorizeFeeAllowance) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgAuthorizeFeeAllowance) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgAuthorizeFeeAllowance) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authorizer}
}

type MsgRevokeFeeAllowance struct {
	Authorizer sdk.AccAddress `json:"authorizer"`
	Operator   sdk.AccAddress `json:"operator"`
}

func (msg MsgRevokeFeeAllowance) Route() string {
	return "delegation"
}

func (msg MsgRevokeFeeAllowance) Type() string {
	return "revoke-fee-allowance"
}

func (msg MsgRevokeFeeAllowance) ValidateBasic() sdk.Error {
	return nil
}

func (msg MsgRevokeFeeAllowance) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return sdk.MustSortJSON(b)
}

func (msg MsgRevokeFeeAllowance) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Authorizer}
}
