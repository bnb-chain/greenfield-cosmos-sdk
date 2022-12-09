package tx

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"

	"cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	ethermint "github.com/evmos/ethermint/types"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/std"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types "github.com/cosmos/cosmos-sdk/types/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/migrations/legacytx"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
)

var msgCodec codec.ProtoCodecMarshaler

var domain = apitypes.TypedDataDomain{
	Name:              "Inscription Tx",
	Version:           "1.0.0",
	VerifyingContract: "inscription",
	Salt:              "0",
}

func init() {
	registry := codectypes.NewInterfaceRegistry()
	std.RegisterInterfaces(registry)
	msgCodec = codec.NewProtoCodec(registry)
}

// signModeEip712Handler defines the SIGN_MODE_DIRECT SignModeHandler
type signModeEip712Handler struct{}

var _ signing.SignModeHandler = signModeEip712Handler{}

// DefaultMode implements SignModeHandler.DefaultMode
func (signModeEip712Handler) DefaultMode() signingtypes.SignMode {
	return signingtypes.SignMode_SIGN_MODE_EIP_712
}

// Modes implements SignModeHandler.Modes
func (signModeEip712Handler) Modes() []signingtypes.SignMode {
	return []signingtypes.SignMode{signingtypes.SignMode_SIGN_MODE_EIP_712}
}

// GetSignBytes implements SignModeHandler.GetSignBytes
func (signModeEip712Handler) GetSignBytes(mode signingtypes.SignMode, signerData signing.SignerData, tx sdk.Tx) ([]byte, error) {
	if mode != signingtypes.SignMode_SIGN_MODE_EIP_712 {
		return nil, fmt.Errorf("expected %s, got %s", signingtypes.SignMode_SIGN_MODE_EIP_712, mode)
	}

	protoTx, ok := tx.(*wrapper)
	if !ok {
		return nil, fmt.Errorf("can only handle a protobuf Tx, got %T", tx)
	}

	cdc := protoTx.cdc
	if protoTx.cdc == nil {
		cdc = msgCodec
	}

	typedChainID, err := ethermint.ParseChainID(signerData.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chainID: %s", signerData.ChainID)
	}

	msgTypes, signDoc, err := GetMsgTypes(cdc, signerData, tx, typedChainID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get msg types")
	}

	typedData, err := WrapTxToTypedData(cdc, typedChainID.Uint64(), protoTx.GetMsgs()[0], signDoc, msgTypes)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to pack tx data in EIP712 object")
	}

	sigHash, err := ComputeTypedDataHash(typedData)
	if err != nil {
		return nil, err
	}

	return sigHash, nil
}

func GetMsgTypes(cdc codectypes.AnyUnpacker, signerData signing.SignerData, tx sdk.Tx, typedChainID *big.Int) (apitypes.Types, *types.SignDocEip712, error) {
	protoTx, ok := tx.(*wrapper)
	if !ok {
		return nil, nil, fmt.Errorf("can only handle a protobuf Tx, got %T", tx)
	}

	signDoc := &types.SignDocEip712{
		AccountNumber: signerData.AccountNumber,
		Sequence:      signerData.Sequence,
		ChainId:       typedChainID.Uint64(),
		TimeoutHeight: protoTx.GetTimeoutHeight(),
		Fee: types.Fee{
			Amount:   protoTx.GetFee(),
			GasLimit: protoTx.GetGas(),
			Payer:    protoTx.FeePayer().String(),
			Granter:  protoTx.FeeGranter().String(),
		},
		Memo: protoTx.GetMemo(),
		Tip:  protoTx.GetTip(),
	}

	msgTypes, err := extractMsgTypes(cdc, protoTx.GetMsgs()[0])
	if err != nil {
		return nil, nil, err
	}

	if signDoc.Tip != nil {
		// patching msgTypes to include Tip
		msgTypes["Tx"] = []apitypes.Type{
			{Name: "account_number", Type: "uint256"},
			{Name: "chain_id", Type: "uint256"},
			{Name: "fee", Type: "Fee"},
			{Name: "memo", Type: "string"},
			{Name: "msg", Type: "Msg"},
			{Name: "sequence", Type: "uint256"},
			{Name: "timeout_height", Type: "uint256"},
			{Name: "tip", Type: "Tip"},
		}
		msgTypes["Tip"] = []apitypes.Type{
			{Name: "amount", Type: "Coin[]"},
			{Name: "tipper", Type: "string"},
		}
	}

	return msgTypes, signDoc, nil
}

// ComputeTypedDataHash computes keccak hash of typed data for signing.
func ComputeTypedDataHash(typedData apitypes.TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		err = errors.Wrap(err, "failed to pack and hash typedData EIP712Domain")
		return nil, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		err = errors.Wrap(err, "failed to pack and hash typedData primary type")
		return nil, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	return crypto.Keccak256(rawData), nil
}

func WrapTxToTypedData(
	cdc codec.Codec,
	chainID uint64,
	msg sdk.Msg,
	signDoc *types.SignDocEip712,
	msgTypes apitypes.Types,
) (apitypes.TypedData, error) {
	var txData map[string]interface{}
	bz, err := cdc.MarshalJSON(signDoc)
	if err != nil {
		return apitypes.TypedData{}, errors.Wrap(sdkerrors.ErrJSONMarshal, "failed to JSON marshal data")
	}
	if err := json.Unmarshal(bz, &txData); err != nil {
		return apitypes.TypedData{}, errors.Wrap(sdkerrors.ErrJSONUnmarshal, "failed to JSON unmarshal data")
	}

	// encode msg
	type msgWrapper struct {
		Msg json.RawMessage `json:"msg"`
	}
	msgValue := make(map[string]interface{})
	legacyMsg := msg.(legacytx.LegacyMsg)
	bz, _ = json.Marshal(msgWrapper{Msg: legacyMsg.GetSignBytes()})
	if err := json.Unmarshal(bz, &msgValue); err != nil {
		panic(err)
	}
	txData["msg"] = msgValue["msg"]

	if txData["tip"] == nil {
		delete(txData, "tip")
	}

	domain.ChainId = math.NewHexOrDecimal256(int64(chainID))

	typedData := apitypes.TypedData{
		Types:       msgTypes,
		PrimaryType: "Tx",
		Domain:      domain,
		Message:     txData,
	}

	return typedData, nil
}

func extractMsgTypes(cdc codectypes.AnyUnpacker, msg sdk.Msg) (apitypes.Types, error) {
	rootTypes := apitypes.Types{
		"EIP712Domain": {
			{
				Name: "name",
				Type: "string",
			},
			{
				Name: "version",
				Type: "string",
			},
			{
				Name: "chainId",
				Type: "uint256",
			},
			{
				Name: "verifyingContract",
				Type: "string",
			},
			{
				Name: "salt",
				Type: "string",
			},
		},
		"Tx": {
			{Name: "account_number", Type: "uint256"},
			{Name: "chain_id", Type: "uint256"},
			{Name: "fee", Type: "Fee"},
			{Name: "memo", Type: "string"},
			{Name: "msg", Type: "Msg"},
			{Name: "sequence", Type: "uint256"},
			{Name: "timeout_height", Type: "uint256"},
		},
		"Fee": {
			{Name: "amount", Type: "Coin[]"},
			{Name: "gas_limit", Type: "uint256"},
			{Name: "payer", Type: "string"},
			{Name: "granter", Type: "string"},
		},
		"Coin": {
			{Name: "denom", Type: "string"},
			{Name: "amount", Type: "uint256"},
		},
		"Msg": {
			{Name: "type", Type: "string"},
			{Name: "value", Type: "MsgValue"},
		},
		"MsgValue": {},
	}

	if err := walkFields(cdc, rootTypes, msg); err != nil {
		return nil, err
	}

	return rootTypes, nil
}

const typeDefPrefix = "_"

func walkFields(cdc codectypes.AnyUnpacker, typeMap apitypes.Types, in interface{}) (err error) {
	defer doRecover(&err)

	t := reflect.TypeOf(in)
	v := reflect.ValueOf(in)

	for {
		if t.Kind() == reflect.Ptr ||
			t.Kind() == reflect.Interface {
			t = t.Elem()
			v = v.Elem()

			continue
		}

		break
	}

	return traverseFields(cdc, typeMap, typeDefPrefix, t, v)
}

type cosmosAnyWrapper struct {
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func traverseFields(
	cdc codectypes.AnyUnpacker,
	typeMap apitypes.Types,
	prefix string,
	t reflect.Type,
	v reflect.Value,
) error {
	n := t.NumField()

	if prefix == typeDefPrefix {
		if len(typeMap["MsgValue"]) == n {
			return nil
		}
	} else {
		typeDef := sanitizeTypedef(prefix)
		if len(typeMap[typeDef]) == n {
			return nil
		}
	}

	for i := 0; i < n; i++ {
		var field reflect.Value
		if v.IsValid() {
			field = v.Field(i)
		}

		fieldType := t.Field(i).Type
		fieldName := jsonNameFromTag(t.Field(i).Tag)

		if fieldType == cosmosAnyType {
			any, ok := field.Interface().(*codectypes.Any)
			if !ok {
				return errors.Wrapf(sdkerrors.ErrPackAny, "%T", field.Interface())
			}

			anyWrapper := &cosmosAnyWrapper{
				Type: any.TypeUrl,
			}

			if err := cdc.UnpackAny(any, &anyWrapper.Value); err != nil {
				return errors.Wrap(err, "failed to unpack Any in msg struct")
			}

			fieldType = reflect.TypeOf(anyWrapper)
			field = reflect.ValueOf(anyWrapper)

			// then continue as normal
		}

		// If it's a nil pointer, do not include in types
		if fieldType.Kind() == reflect.Ptr && field.IsNil() {
			continue
		}
		if fieldType.Kind() == reflect.String && field.String() == "" {
			continue
		}

		for {
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()

				if field.IsValid() {
					field = field.Elem()
				}

				continue
			}

			if fieldType.Kind() == reflect.Interface {
				fieldType = reflect.TypeOf(field.Interface())
				continue
			}

			if field.Kind() == reflect.Ptr {
				field = field.Elem()
				continue
			}

			break
		}

		var isCollection bool
		if fieldType.Kind() == reflect.Array || fieldType.Kind() == reflect.Slice {
			if field.Len() == 0 {
				// skip empty collections from type mapping
				continue
			}

			fieldType = fieldType.Elem()
			field = field.Index(0)
			isCollection = true
		}

		for {
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()

				if field.IsValid() {
					field = field.Elem()
				}

				continue
			}

			if fieldType.Kind() == reflect.Interface {
				fieldType = reflect.TypeOf(field.Interface())
				continue
			}

			if field.Kind() == reflect.Ptr {
				field = field.Elem()
				continue
			}

			break
		}

		fieldPrefix := fmt.Sprintf("%s.%s", prefix, fieldName)

		ethTyp := typToEth(fieldType)
		if len(ethTyp) > 0 {
			// Support array of uint64
			if isCollection && fieldType.Kind() != reflect.Slice && fieldType.Kind() != reflect.Array {
				ethTyp += "[]"
			}

			if prefix == typeDefPrefix {
				typeMap["MsgValue"] = append(typeMap["MsgValue"], apitypes.Type{
					Name: fieldName,
					Type: ethTyp,
				})
			} else {
				typeDef := sanitizeTypedef(prefix)
				typeMap[typeDef] = append(typeMap[typeDef], apitypes.Type{
					Name: fieldName,
					Type: ethTyp,
				})
			}

			continue
		}

		if fieldType.Kind() == reflect.Struct {
			var fieldTypedef string

			if isCollection {
				fieldTypedef = sanitizeTypedef(fieldPrefix) + "[]"
			} else {
				fieldTypedef = sanitizeTypedef(fieldPrefix)
			}

			if prefix == typeDefPrefix {
				typeMap["MsgValue"] = append(typeMap["MsgValue"], apitypes.Type{
					Name: fieldName,
					Type: fieldTypedef,
				})
			} else {
				typeDef := sanitizeTypedef(prefix)
				typeMap[typeDef] = append(typeMap[typeDef], apitypes.Type{
					Name: fieldName,
					Type: fieldTypedef,
				})
			}

			if err := traverseFields(cdc, typeMap, fieldPrefix, fieldType, field); err != nil {
				return err
			}

			continue
		}
	}

	return nil
}

func jsonNameFromTag(tag reflect.StructTag) string {
	jsonTags := tag.Get("json")
	parts := strings.Split(jsonTags, ",")
	return parts[0]
}

// _.foo_bar.baz -> TypeFooBarBaz
//
// this is needed for Geth's own signing code which doesn't
// tolerate complex type names
func sanitizeTypedef(str string) string {
	buf := new(bytes.Buffer)
	parts := strings.Split(str, ".")
	caser := cases.Title(language.English, cases.NoLower)

	for _, part := range parts {
		if part == "_" {
			buf.WriteString("Type")
			continue
		}

		subparts := strings.Split(part, "_")
		for _, subpart := range subparts {
			buf.WriteString(caser.String(subpart))
		}
	}

	return buf.String()
}

var (
	hashType      = reflect.TypeOf(common.Hash{})
	addressType   = reflect.TypeOf(common.Address{})
	bigIntType    = reflect.TypeOf(big.Int{})
	cosmIntType   = reflect.TypeOf(sdkmath.Int{})
	cosmDecType   = reflect.TypeOf(sdk.Dec{})
	cosmosAnyType = reflect.TypeOf(&codectypes.Any{})
	timeType      = reflect.TypeOf(time.Time{})

	edType = reflect.TypeOf(ed25519.PubKey{})
)

// typToEth supports only basic types and arrays of basic types.
// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-712.md
func typToEth(typ reflect.Type) string {
	const str = "string"

	switch typ.Kind() {
	case reflect.String:
		return str
	case reflect.Bool:
		return "bool"
	case reflect.Int:
		return "int64"
	case reflect.Int8:
		return "int8"
	case reflect.Int16:
		return "int16"
	case reflect.Int32:
		return "int32"
	case reflect.Int64:
		return "int64"
	case reflect.Uint:
		return "uint64"
	case reflect.Uint8:
		return "uint8"
	case reflect.Uint16:
		return "uint16"
	case reflect.Uint32:
		return "uint32"
	case reflect.Uint64:
		return "uint64"
	case reflect.Slice:
		ethName := typToEth(typ.Elem())
		if len(ethName) > 0 {
			return ethName + "[]"
		}
	case reflect.Array:
		ethName := typToEth(typ.Elem())
		if len(ethName) > 0 {
			return ethName + "[]"
		}
	case reflect.Ptr:
		if typ.Elem().ConvertibleTo(bigIntType) ||
			typ.Elem().ConvertibleTo(timeType) ||
			typ.Elem().ConvertibleTo(edType) ||
			typ.Elem().ConvertibleTo(cosmDecType) ||
			typ.Elem().ConvertibleTo(cosmIntType) {
			return str
		}
	case reflect.Struct:
		if typ.ConvertibleTo(hashType) ||
			typ.ConvertibleTo(addressType) ||
			typ.ConvertibleTo(bigIntType) ||
			typ.ConvertibleTo(edType) ||
			typ.ConvertibleTo(timeType) ||
			typ.ConvertibleTo(cosmDecType) ||
			typ.ConvertibleTo(cosmIntType) {
			return str
		}
	}

	return ""
}

func doRecover(err *error) {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			e = errors.Wrap(e, "panicked with error")
			*err = e
			return
		}

		*err = fmt.Errorf("%v", r)
	}
}
