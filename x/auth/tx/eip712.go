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
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types "github.com/cosmos/cosmos-sdk/types/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

var MsgCodec codec.ProtoCodecMarshaler

var domain = apitypes.TypedDataDomain{
	Name:              "Inscription Tx",
	Version:           "1.0.0",
	VerifyingContract: "inscription",
	Salt:              "0",
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

	var msgCodec codec.Codec
	if MsgCodec == nil {
		if protoTx.cdc == nil {
			return nil, fmt.Errorf("no proto codec marshaler")
		}
		msgCodec = protoTx.cdc
	} else {
		msgCodec = MsgCodec
	}

	typedChainID, err := ethermint.ParseChainID(signerData.ChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to parse chainID: %s", signerData.ChainID)
	}

	msgTypes, signDoc, err := GetMsgTypes(msgCodec, signerData, tx, typedChainID)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get msg types")
	}

	typedData, err := WrapTxToTypedData(msgCodec, typedChainID.Uint64(), signDoc, msgTypes)
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

	msgAny, _ := codectypes.NewAnyWithValue(protoTx.GetMsgs()[0])
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
		Msg:  msgAny,
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
	signDoc *types.SignDocEip712,
	msgTypes apitypes.Types,
) (apitypes.TypedData, error) {
	var txData map[string]interface{}
	bz, err := cdc.MarshalJSON(signDoc)
	if err != nil {
		return apitypes.TypedData{}, errors.Wrap(err, "failed to JSON marshal data")
	}
	if err := json.Unmarshal(bz, &txData); err != nil {
		return apitypes.TypedData{}, errors.Wrap(err, "failed to JSON unmarshal data")
	}

	if txData["tip"] == nil {
		delete(txData, "tip")
	}

	// filling nil value
	cleanTypesAndMsgValue(msgTypes, "Msg", txData["msg"].(map[string]interface{}))

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
		},
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

	return traverseFields(cdc, typeMap, typeDefPrefix, typeDefPrefix, t, v, false)
}

type cosmosAnyWrapper struct {
	Type           string      `json:"type"`
	CosmosAnyValue interface{} `json:"cosmos_any_value"`
}

func traverseFields(
	cdc codectypes.AnyUnpacker,
	typeMap apitypes.Types,
	prefix string,
	prePrefix string,
	t reflect.Type,
	v reflect.Value,
	isCosmosAnyValue bool,
) error {
	n := t.NumField()

	for i := 0; i < n; i++ {
		var field reflect.Value
		if v.IsValid() {
			field = v.Field(i)
		}

		fieldType := t.Field(i).Type
		fieldName := jsonNameFromTag(t.Field(i).Tag)
		isOmitEmpty := isOmitEmpty(t.Field(i).Tag)

		if fieldName == "" {
			// For protobuf one_of interface, there's no json tag.
			// So we need to unwrap it first.
			if isProtobufOneOf(t.Field(i).Tag) {
				fieldType = reflect.TypeOf(field.Interface())
				field = reflect.ValueOf(field.Interface())
				if fieldType.Kind() == reflect.Ptr {
					fieldType = fieldType.Elem()
					if field.IsValid() {
						field = field.Elem()
					}
				}
				field = field.Field(0)
				fieldName = jsonNameFromTag(fieldType.Field(0).Tag)
				fieldType = fieldType.Field(0).Type
			} else {
				panic("empty json tag")
			}
		}

		// Unpack any type into the wrapper to keep align between the structs' json tag and field name
		if fieldType == cosmosAnyType {
			typeAny, ok := field.Interface().(*codectypes.Any)
			if !ok {
				return errors.Wrapf(sdkerrors.ErrPackAny, "%T", field.Interface())
			}

			anyWrapper := &cosmosAnyWrapper{
				Type: typeAny.TypeUrl,
			}

			if err := cdc.UnpackAny(typeAny, &anyWrapper.CosmosAnyValue); err != nil {
				return errors.Wrap(err, "failed to unpack Any in msg struct")
			}

			fieldType = reflect.TypeOf(anyWrapper)
			field = reflect.ValueOf(anyWrapper)
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
				field = reflect.ValueOf(field.Interface())
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
				if !isOmitEmpty {
					fieldType = reflect.TypeOf("str")
					goto FINISH
				}
				// skip empty collections from type mapping if is omitEmpty
				continue
			}

			fieldType = fieldType.Elem()
			field = field.Index(0)

			// Unpack any type into the wrapper to keep align between the structs' json tag and field name
			if fieldType == cosmosAnyType {
				typeAny, ok := field.Interface().(*codectypes.Any)
				if !ok {
					return errors.Wrapf(sdkerrors.ErrPackAny, "%T", field.Interface())
				}

				anyWrapper := &cosmosAnyWrapper{
					Type: typeAny.TypeUrl,
				}

				if err := cdc.UnpackAny(typeAny, &anyWrapper.CosmosAnyValue); err != nil {
					return errors.Wrap(err, "failed to unpack Any in msg struct")
				}

				fieldType = reflect.TypeOf(anyWrapper)
				field = reflect.ValueOf(anyWrapper)
			}
		FINISH:
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
		fieldPrePrefix := fmt.Sprintf("%s.%s", prePrefix, fieldName)

		ethTyp := typToEth(fieldType)
		if len(ethTyp) > 0 {
			// Support array of uint64
			if isCollection && fieldType.Kind() != reflect.Slice && fieldType.Kind() != reflect.Array {
				ethTyp += "[]"
				if ethTyp == "uint8[]" {
					ethTyp = "string"
				}
			}

			switch {
			case prefix == typeDefPrefix:
				typeMap["Msg"] = append(typeMap["Msg"], apitypes.Type{
					Name: fieldName,
					Type: ethTyp,
				})
			case isCosmosAnyValue:
				typeDef := sanitizeTypedef(prePrefix)
				typeMap[typeDef] = append(typeMap[typeDef], apitypes.Type{
					Name: fieldName,
					Type: ethTyp,
				})
			default:
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
			var fieldPreTypedef string

			if isCollection {
				fieldTypedef = sanitizeTypedef(fieldPrefix) + "[]"
				fieldPreTypedef = sanitizeTypedef(fieldPrePrefix) + "[]"
			} else {
				fieldTypedef = sanitizeTypedef(fieldPrefix)
				fieldPreTypedef = sanitizeTypedef(fieldPrePrefix)
			}

			switch {
			case prefix == typeDefPrefix:
				typeMap["Msg"] = append(typeMap["Msg"], apitypes.Type{
					Name: fieldName,
					Type: fieldTypedef,
				})

				if err := traverseFields(cdc, typeMap, fieldPrefix, prefix, fieldType, field, false); err != nil {
					return err
				}
			case fieldName == "cosmos_any_value":
				if err := traverseFields(cdc, typeMap, fieldPrefix, prefix, fieldType, field, true); err != nil {
					return err
				}
			case isCosmosAnyValue:
				typeDef := sanitizeTypedef(prePrefix)
				typeMap[typeDef] = append(typeMap[typeDef], apitypes.Type{
					Name: fieldName,
					Type: fieldPreTypedef,
				})

				if err := traverseFields(cdc, typeMap, fieldPrePrefix, prePrefix, fieldType, field, false); err != nil {
					return err
				}
			default:
				typeDef := sanitizeTypedef(prefix)
				typeMap[typeDef] = append(typeMap[typeDef], apitypes.Type{
					Name: fieldName,
					Type: fieldTypedef,
				})

				if err := traverseFields(cdc, typeMap, fieldPrefix, prefix, fieldType, field, false); err != nil {
					return err
				}
			}

			continue
		}
	}

	return nil
}

// jsonNameFromTag gets json name
func jsonNameFromTag(tag reflect.StructTag) string {
	jsonTags := tag.Get("json")
	parts := strings.Split(jsonTags, ",")
	return parts[0]
}

// isOmitEmpty returns if the struct has emitempty tag
func isOmitEmpty(tag reflect.StructTag) bool {
	jsonTags := tag.Get("json")
	parts := strings.Split(jsonTags, ",")
	for _, tag := range parts {
		if tag == "omitempty" {
			return true
		}
	}
	return false
}

// isProtobufOneOf returns if the struct is protobuf_oneof type
func isProtobufOneOf(tag reflect.StructTag) bool {
	return tag.Get("protobuf_oneof") != ""
}

func cleanTypesAndMsgValue(typedData apitypes.Types, primaryType string, msgValue map[string]interface{}) {
	// 1. the proto codec will set *types.Any's type struct name to be "@type". Need remove prefix "@"
	if msgValue["@type"] != nil {
		msgValue["type"] = msgValue["@type"]
		delete(msgValue, "@type")
	}

	// 2. clean msg value.
	for i, field := range typedData[primaryType] {
		encType := field.Type
		encValue := msgValue[field.Name]
		switch {
		case encType[len(encType)-1:] == "]":
			if typedData[encType[:len(encType)-2]] != nil {
				for i := 0; i < len(msgValue[field.Name].([]interface{})); i++ {
					cleanTypesAndMsgValue(typedData, encType[:len(encType)-2], msgValue[field.Name].([]interface{})[i].(map[string]interface{}))
				}
			}
		case typedData[encType] != nil:
			subType, ok := msgValue[field.Name].(map[string]interface{})
			if !ok {
				// Delete nil struct
				typedData[primaryType] = append(typedData[primaryType][:i], typedData[primaryType][i+1:]...)
				delete(typedData, encType)
				delete(msgValue, field.Name)
			} else {
				cleanTypesAndMsgValue(typedData, encType, subType)
			}
		case encValue == nil:
			// If the field's type is *types.Any and there are only 2 sub-fields, the sub-field's name need to fix
			if field.Name == "cosmos_any_value" {
				if len(msgValue) != 2 {
					panic("unexpected msg value")
				}
				for key := range msgValue {
					if key != "type" {
						typedData[primaryType] = append(typedData[primaryType], apitypes.Type{
							Name: key,
							Type: field.Type,
						})
						typedData[primaryType] = append(typedData[primaryType][:i], typedData[primaryType][i+1:]...)
						break
					}
				}
				break
			}
			// For nil primitive value, fill in default value
			switch encType {
			case "bool":
				msgValue[field.Name] = false
			case "string":
				msgValue[field.Name] = ""
			default:
				if strings.HasPrefix(encType, "uint") || strings.HasPrefix(encType, "int") {
					msgValue[field.Name] = 0
				}
			}
		}
	}

	// Delete nil struct
	for key := range msgValue {
		for _, field := range typedData[primaryType] {
			if field.Name == key {
				goto FINISH
			}
		}
		delete(msgValue, key)
	FINISH:
	}
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
	hashType         = reflect.TypeOf(common.Hash{})
	addressType      = reflect.TypeOf(common.Address{})
	bigIntType       = reflect.TypeOf(big.Int{})
	cosmIntType      = reflect.TypeOf(sdkmath.Int{})
	cosmDecType      = reflect.TypeOf(sdk.Dec{})
	cosmosAnyType    = reflect.TypeOf(&codectypes.Any{})
	timeType1        = reflect.TypeOf(time.Time{})
	timeType2        = reflect.TypeOf(&time.Time{})
	timeDurationType = reflect.TypeOf(time.Duration(1))

	// enum type
	enumType = reflect.TypeOf(stakingtypes.AuthorizationType(0))

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
		if typ.ConvertibleTo(enumType) {
			return str
		}
		return "int32"
	case reflect.Int64:
		if typ.ConvertibleTo(timeDurationType) {
			return str
		}
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
			typ.Elem().ConvertibleTo(timeType2) ||
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
			typ.ConvertibleTo(timeType1) ||
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
