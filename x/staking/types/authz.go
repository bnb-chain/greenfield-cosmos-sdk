package types

import (
	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

const maxValidatorsNumber = 41

var _ authz.Authorization = &StakeAuthorization{}

// NewStakeAuthorization creates a new StakeAuthorization object.
func NewStakeAuthorization(allowed []sdk.ValAddress, denied []sdk.ValAddress, authzType AuthorizationType, amount *sdk.Coin) (*StakeAuthorization, error) {
	allowedValidators, deniedValidators, err := validateAllowAndDenyValidators(allowed, denied)
	if err != nil {
		return nil, err
	}

	a := StakeAuthorization{}
	if allowedValidators != nil {
		a.Validators = &StakeAuthorization_AllowList{AllowList: &StakeAuthorization_Validators{Address: allowedValidators}}
	} else {
		a.Validators = &StakeAuthorization_DenyList{DenyList: &StakeAuthorization_Validators{Address: deniedValidators}}
	}

	if amount != nil {
		a.MaxTokens = amount
	}
	a.AuthorizationType = authzType

	return &a, nil
}

// MsgTypeURL implements Authorization.MsgTypeURL.
func (a StakeAuthorization) MsgTypeURL() string {
	authzType, err := normalizeAuthzType(a.AuthorizationType)
	if err != nil {
		panic(err)
	}
	return authzType
}

func (a StakeAuthorization) ValidateBasic() error {
	if a.MaxTokens != nil && a.MaxTokens.IsNegative() {
		return errors.Wrapf(authz.ErrNegativeMaxTokens, "negative coin amount: %v", a.MaxTokens)
	}
	if a.AuthorizationType == AuthorizationType_AUTHORIZATION_TYPE_UNSPECIFIED {
		return authz.ErrUnknownAuthorizationType
	}
	if allowList := a.GetAllowList().GetAddress(); allowList != nil {
		if len(allowList) > maxValidatorsNumber {
			return errors.Wrapf(authz.ErrTooManyValidators, "allow list number: %d, limit: %d", len(allowList), maxValidatorsNumber)
		}
	}
	if denyList := a.GetDenyList().GetAddress(); denyList != nil {
		if len(denyList) > maxValidatorsNumber {
			return errors.Wrapf(authz.ErrTooManyValidators, "deny list number: %d, limit: %d", len(denyList), maxValidatorsNumber)
		}
	}

	return nil
}

// Accept implements Authorization.Accept.
func (a StakeAuthorization) Accept(ctx sdk.Context, msg sdk.Msg) (authz.AcceptResponse, error) {
	var validatorAddress string
	var amount sdk.Coin

	switch msg := msg.(type) {
	case *MsgDelegate:
		validatorAddress = msg.ValidatorAddress
		amount = msg.Amount
	case *MsgUndelegate:
		validatorAddress = msg.ValidatorAddress
		amount = msg.Amount
	case *MsgBeginRedelegate:
		validatorAddress = msg.ValidatorDstAddress
		amount = msg.Amount
	default:
		return authz.AcceptResponse{}, sdkerrors.ErrInvalidRequest.Wrap("unknown msg type")
	}

	isValidatorExists := false
	allowedList := a.GetAllowList().GetAddress()
	for _, validator := range allowedList {
		if validator == validatorAddress {
			isValidatorExists = true
			break
		}
	}

	denyList := a.GetDenyList().GetAddress()
	for _, validator := range denyList {
		if validator == validatorAddress {
			return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrapf("cannot delegate/undelegate to %s validator", validator)
		}
	}

	if len(allowedList) > 0 && !isValidatorExists {
		return authz.AcceptResponse{}, sdkerrors.ErrUnauthorized.Wrapf("cannot delegate/undelegate to %s validator", validatorAddress)
	}

	if a.MaxTokens == nil {
		return authz.AcceptResponse{
			Accept: true, Delete: false,
			Updated: &StakeAuthorization{Validators: a.GetValidators(), AuthorizationType: a.GetAuthorizationType()},
		}, nil
	}

	limitLeft, err := a.MaxTokens.SafeSub(amount)
	if err != nil {
		return authz.AcceptResponse{}, err
	}
	if limitLeft.IsZero() {
		return authz.AcceptResponse{Accept: true, Delete: true}, nil
	}
	return authz.AcceptResponse{
		Accept: true, Delete: false,
		Updated: &StakeAuthorization{Validators: a.GetValidators(), AuthorizationType: a.GetAuthorizationType(), MaxTokens: &limitLeft},
	}, nil
}

func validateAllowAndDenyValidators(allowed []sdk.ValAddress, denied []sdk.ValAddress) ([]string, []string, error) {
	if len(allowed) == 0 && len(denied) == 0 {
		return nil, nil, sdkerrors.ErrInvalidRequest.Wrap("both allowed & deny list cannot be empty")
	}

	if len(allowed) > 0 && len(denied) > 0 {
		return nil, nil, sdkerrors.ErrInvalidRequest.Wrap("cannot set both allowed & deny list")
	}

	allowedValidators := make([]string, len(allowed))
	if len(allowed) > 0 {
		for i, validator := range allowed {
			allowedValidators[i] = validator.String()
		}
		return allowedValidators, nil, nil
	}

	deniedValidators := make([]string, len(denied))
	for i, validator := range denied {
		deniedValidators[i] = validator.String()
	}

	return nil, deniedValidators, nil
}

// Normalized Msg type URLs
func normalizeAuthzType(authzType AuthorizationType) (string, error) {
	switch authzType {
	case AuthorizationType_AUTHORIZATION_TYPE_DELEGATE:
		return sdk.MsgTypeURL(&MsgDelegate{}), nil
	case AuthorizationType_AUTHORIZATION_TYPE_UNDELEGATE:
		return sdk.MsgTypeURL(&MsgUndelegate{}), nil
	case AuthorizationType_AUTHORIZATION_TYPE_REDELEGATE:
		return sdk.MsgTypeURL(&MsgBeginRedelegate{}), nil
	default:
		return "", sdkerrors.Wrapf(authz.ErrUnknownAuthorizationType, "cannot normalize authz type with %T", authzType)
	}
}
