package tx

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"strings"
	"time"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/signer/core/apitypes"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	"github.com/cosmos/cosmos-sdk/crypto/keys/eth/ethsecp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	types "github.com/cosmos/cosmos-sdk/types/tx"
	signingtypes "github.com/cosmos/cosmos-sdk/types/tx/signing"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

var (
	domain = &apitypes.TypedDataDomain{
		Name:              "Greenfield Tx",
		Version:           "1.0.0",
		VerifyingContract: "greenfield",
		Salt:              "0",
	}

	gnfdVerifyingContract = "0x71e835aff094655dEF897fbc85534186DbeaB75d" // keccak256("greenfield")[12:]
)

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
func (h signModeEip712Handler) GetSignBytes(mode signingtypes.SignMode, signerData signing.SignerData, tx sdk.Tx) ([]byte, error) {
	return getSignBytes(mode, signerData, tx, false)
}

// GetSignBytesRuntime implements SignModeHandler.GetSignBytesRuntime
func (h signModeEip712Handler) GetSignBytesRuntime(ctx sdk.Context, mode signingtypes.SignMode, signerData signing.SignerData, tx sdk.Tx) ([]byte, error) {
	if !ctx.IsUpgraded(upgradetypes.Altai) {
		return h.GetSignBytes(mode, signerData, tx)
	}
	return getSignBytes(mode, signerData, tx, true)
}

func getSignBytes(mode signingtypes.SignMode, signerData signing.SignerData, tx sdk.Tx, isAltai bool) ([]byte, error) {
	if mode != signingtypes.SignMode_SIGN_MODE_EIP_712 {
		return nil, fmt.Errorf("expected %s, got %s", signingtypes.SignMode_SIGN_MODE_EIP_712, mode)
	}

	chainID, err := sdk.ParseChainID(signerData.ChainID)
	if err != nil {
		return nil, err
	}

	msgTypes, signDoc, err := GetMsgTypes(signerData, tx, chainID)
	if err != nil {
		return nil, err
	}

	typedDataDomain := *domain
	typedDataDomain.ChainId = math.NewHexOrDecimal256(chainID.Int64())
	if isAltai {
		typedDataDomain.VerifyingContract = gnfdVerifyingContract
	}

	typedData, err := WrapTxToTypedData(signDoc, msgTypes, typedDataDomain)
	if err != nil {
		return nil, err
	}

	return ComputeTypedDataHash(typedData)
}

func GetMsgTypes(signerData signing.SignerData, tx sdk.Tx, typedChainID *big.Int) (apitypes.Types, *types.SignDocEip712, error) {
	protoTx, ok := tx.(*wrapper)
	if !ok {
		return nil, nil, fmt.Errorf("can only handle a protobuf Tx, got %T", tx)
	}

	// construct the signDoc
	msgAnys := make([]*codectypes.Any, 0, len(protoTx.GetMsgs()))
	for _, msg := range protoTx.GetMsgs() {
		msgAny, _ := codectypes.NewAnyWithValue(msg)
		msgAnys = append(msgAnys, msgAny)
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
		Msgs: msgAnys,
	}

	// extract the msg types
	msgTypes := apitypes.Types{
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
			{Name: "sequence", Type: "uint256"},
			{Name: "timeout_height", Type: "uint256"},
		},
		"Fee": {
			{Name: "amount", Type: "Coin[]"},
			{Name: "gas_limit", Type: "uint256"},
			{Name: "granter", Type: "string"},
			{Name: "payer", Type: "string"},
		},
		"Coin": {
			{Name: "amount", Type: "uint256"},
			{Name: "denom", Type: "string"},
		},
	}
	for i, msg := range protoTx.GetMsgs() {
		tmpMsgTypes, err := extractMsgTypes(msg, i+1)
		if err != nil {
			return nil, nil, err
		}

		msgTypes["Tx"] = append(msgTypes["Tx"], apitypes.Type{
			Name: fmt.Sprintf("msg%d", i+1),
			Type: fmt.Sprintf("Msg%d", i+1),
		})

		for key, field := range tmpMsgTypes {
			msgTypes[key] = field
		}
	}

	// patch the msg types to include `Tip` if it's not empty
	if signDoc.Tip != nil {
		msgTypes["Tx"] = append(msgTypes["Tx"], apitypes.Type{Name: "tip", Type: "Tip"})
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
		err = errorsmod.Wrap(err, "failed to pack and hash typedData EIP712Domain")
		return nil, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		err = errorsmod.Wrap(err, "failed to pack and hash typedData primary type")
		return nil, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	return crypto.Keccak256(rawData), nil
}

func WrapTxToTypedData(
	signDoc *types.SignDocEip712,
	msgTypes apitypes.Types,
	typedDataDomain apitypes.TypedDataDomain,
) (apitypes.TypedData, error) {
	msgCodec := jsonpb.Marshaler{
		EmitDefaults: true,
		OrigName:     true,
	}
	bz, err := msgCodec.MarshalToString(signDoc)
	if err != nil {
		return apitypes.TypedData{}, errorsmod.Wrap(err, "failed to JSON marshal data")
	}

	var txData map[string]interface{}
	if err := json.Unmarshal([]byte(bz), &txData); err != nil {
		return apitypes.TypedData{}, errorsmod.Wrap(err, "failed to JSON unmarshal data")
	}

	if txData["tip"] == nil {
		delete(txData, "tip")
	}

	// filling nil value and do the other clean up
	msgData := txData["msgs"].([]interface{})
	for i := range signDoc.GetMsgs() {
		txData[fmt.Sprintf("msg%d", i+1)] = msgData[i]
		cleanTypesAndMsgValue(msgTypes, fmt.Sprintf("Msg%d", i+1), msgData[i].(map[string]interface{}))
	}
	delete(txData, "msgs")

	// sort the msg types
	for _, val := range msgTypes {
		sort.Slice(val, func(i, j int) bool {
			return val[i].Name < val[j].Name
		})
	}

	typedData := apitypes.TypedData{
		Types:       msgTypes,
		PrimaryType: "Tx",
		Domain:      typedDataDomain,
		Message:     txData,
	}

	return typedData, nil
}

func extractMsgTypes(msg sdk.Msg, index int) (apitypes.Types, error) {
	rootTypes := apitypes.Types{
		fmt.Sprintf("Msg%d", index): {
			{Name: "type", Type: "string"},
		},
	}

	if err := walkFields(rootTypes, msg, index); err != nil {
		return nil, err
	}

	return rootTypes, nil
}

const typeDefPrefix = "_"

func walkFields(typeMap apitypes.Types, in interface{}, index int) (err error) {
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

	return traverseFields(typeMap, typeDefPrefix, index, t, v)
}

type anyWrapper struct {
	Type  string `json:"type"`
	Value []byte `json:"value"`
}

func traverseFields(
	typeMap apitypes.Types,
	prefix string,
	index int,
	t reflect.Type,
	v reflect.Value,
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

		fieldType, field, fieldName = unwrapField(fieldType, field, fieldName)

		var isCollection bool
		if fieldType.Kind() == reflect.Array || fieldType.Kind() == reflect.Slice {
			isCollection = true
			if field.Len() == 0 {
				if !isOmitEmpty {
					fieldType = reflect.TypeOf("str")
				} else {
					// skip empty collections from type mapping if is omitEmpty
					continue
				}
			} else {
				fieldType = fieldType.Elem()
				field = field.Index(0)

				fieldType, field, fieldName = unwrapField(fieldType, field, fieldName)
			}
		}

		fieldPrefix := fmt.Sprintf("%s.%s", prefix, fieldName)

		ethTyp := typToEth(fieldType)
		if len(ethTyp) > 0 {
			if isCollection {
				if fieldType.Kind() != reflect.Slice && fieldType.Kind() != reflect.Array {
					ethTyp += "[]"
					if ethTyp == "uint8[]" {
						ethTyp = "bytes"
					}
				} else if ethTyp == "uint8[]" {
					ethTyp = "bytes[]"
				}
			}

			if prefix == typeDefPrefix {
				tag := fmt.Sprintf("Msg%d", index)
				typeMap[tag] = append(typeMap[tag], apitypes.Type{
					Name: fieldName,
					Type: ethTyp,
				})
			} else {
				typeDef := sanitizeTypedef(prefix, index)
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
				fieldTypedef = sanitizeTypedef(fieldPrefix, index) + "[]"
			} else {
				fieldTypedef = sanitizeTypedef(fieldPrefix, index)
			}

			if prefix == typeDefPrefix {
				tag := fmt.Sprintf("Msg%d", index)
				typeMap[tag] = append(typeMap[tag], apitypes.Type{
					Name: fieldName,
					Type: fieldTypedef,
				})
			} else {
				typeDef := sanitizeTypedef(prefix, index)
				typeMap[typeDef] = append(typeMap[typeDef], apitypes.Type{
					Name: fieldName,
					Type: fieldTypedef,
				})
			}

			if err := traverseFields(typeMap, fieldPrefix, index, fieldType, field); err != nil {
				return err
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
		encName := field.Name
		encType := field.Type
		if strings.HasSuffix(encName, "Any") {
			if encType[len(encType)-1:] == "]" {
				anySet := msgValue[encName[:len(encName)-3]].([]interface{})
				newAnySet := make([]interface{}, len(anySet))
				for i, item := range anySet {
					anyValue := item.(map[string]interface{})
					newValue := make(map[string]interface{})
					bz, _ := json.Marshal(anyValue)
					newValue["type"] = anyValue["@type"]
					newValue["value"] = bz
					newAnySet[i] = newValue
				}
				msgValue[encName[:len(encName)-3]] = newAnySet
				typedData[primaryType][i].Name = encName[:len(encName)-3]
				typedData[primaryType][i].Type = "TypeAny[]"
				delete(typedData, encType[:len(encType)-2])
			} else {
				anyValue := msgValue[encName[:len(encName)-3]].(map[string]interface{})
				newValue := make(map[string]interface{})
				bz, _ := json.Marshal(anyValue)
				newValue["type"] = anyValue["@type"]
				base64Str := base64.StdEncoding.EncodeToString(bz) // base64 encode to keep consistency with js-sdk
				newValue["value"] = []byte(base64Str)
				msgValue[encName[:len(encName)-3]] = newValue
				typedData[primaryType][i].Name = encName[:len(encName)-3]
				typedData[primaryType][i].Type = "TypeAny"
				delete(typedData, encType)
			}
			typedData["TypeAny"] = anyApiTypes
			continue
		}
		encValue := msgValue[encName]
		switch {
		case encType[len(encType)-1:] == "]":
			if typedData[encType[:len(encType)-2]] != nil {
				for i := 0; i < len(msgValue[encName].([]interface{})); i++ {
					cleanTypesAndMsgValue(typedData, encType[:len(encType)-2], msgValue[encName].([]interface{})[i].(map[string]interface{}))
				}
			} else if encType == "bytes[]" {
				// convert string to type
				byteList, ok := encValue.([]interface{})
				if !ok {
					continue
				}
				newBytesList := make([]interface{}, len(byteList))
				for j, item := range byteList {
					newBytesList[j] = []byte(item.(string))
				}
				msgValue[encName] = newBytesList
			}
		case typedData[encType] != nil:
			subType, ok := msgValue[encName].(map[string]interface{})
			if !ok {
				// Delete nil struct
				typedData[primaryType] = append(typedData[primaryType][:i], typedData[primaryType][i+1:]...)
				delete(typedData, encType)
				delete(msgValue, encName)
			} else {
				cleanTypesAndMsgValue(typedData, encType, subType)
			}
		case encValue == nil:
			// For nil primitive value, fill in default value
			switch encType {
			case "bool":
				msgValue[encName] = false
			case "string":
				msgValue[encName] = ""
			default:
				if strings.HasPrefix(encType, "uint") || strings.HasPrefix(encType, "int") {
					msgValue[encName] = "0"
				}
			}
		case encType == "bytes":
			if reflect.TypeOf(encValue).Kind() == reflect.String {
				msgValue[encName] = []byte(encValue.(string))
			}
		}
	}

	// Delete nil struct
	for key := range msgValue {
		var isExist bool
		for _, field := range typedData[primaryType] {
			if field.Name == key {
				isExist = true
			}
		}
		if !isExist {
			delete(msgValue, key)
		}
	}
}

func unwrapField(fieldType reflect.Type, field reflect.Value, fieldName string) (reflect.Type, reflect.Value, string) {
	// Unpack any type into the wrapper to keep align between the structs' json tag and field name
	var anyValue *anyWrapper
	if fieldType == cosmosAnyType {
		typeAny, ok := field.Interface().(*codectypes.Any)
		if !ok {
			panic(sdkerrors.ErrPackAny)
		}

		anyValue = &anyWrapper{
			Type:  typeAny.TypeUrl,
			Value: typeAny.Value,
		}

		fieldType = reflect.TypeOf(anyValue)
		field = reflect.ValueOf(anyValue)
		fieldName += "Any"
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

		break
	}

	return fieldType, field, fieldName
}

// _.foo_bar.baz -> TypeFooBarBaz
//
// this is needed for Geth's own signing code which doesn't
// tolerate complex type names
func sanitizeTypedef(str string, index int) string {
	buf := new(bytes.Buffer)
	parts := strings.Split(str, ".")
	caser := cases.Title(language.English, cases.NoLower)

	for _, part := range parts {
		if part == "_" {
			buf.WriteString(fmt.Sprintf("TypeMsg%d", index))
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
	cosmosIntType    = reflect.TypeOf(sdkmath.Int{})
	cosmosDecType    = reflect.TypeOf(sdk.Dec{})
	cosmosAnyType    = reflect.TypeOf(&codectypes.Any{})
	timeType         = reflect.TypeOf(time.Time{})
	timePtrType      = reflect.TypeOf(&time.Time{})
	timeDurationType = reflect.TypeOf(time.Duration(1))
	enumType         = reflect.TypeOf(stakingtypes.AuthorizationType(0))
	edType           = reflect.TypeOf(ed25519.PubKey{})
	secpType         = reflect.TypeOf(ethsecp256k1.PubKey{})

	anyApiTypes = []apitypes.Type{
		{
			Name: "type",
			Type: "string",
		},
		{
			Name: "value",
			Type: "bytes",
		},
	}
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
	case reflect.Slice, reflect.Array:
		ethName := typToEth(typ.Elem())
		if len(ethName) > 0 {
			return ethName + "[]"
		}
	case reflect.Ptr:
		if typ.Elem().ConvertibleTo(bigIntType) ||
			typ.Elem().ConvertibleTo(timePtrType) ||
			typ.Elem().ConvertibleTo(edType) ||
			typ.Elem().ConvertibleTo(secpType) ||
			typ.Elem().ConvertibleTo(cosmosDecType) ||
			typ.Elem().ConvertibleTo(cosmosIntType) {
			return str
		}
	case reflect.Struct:
		if typ.ConvertibleTo(hashType) ||
			typ.ConvertibleTo(addressType) ||
			typ.ConvertibleTo(bigIntType) ||
			typ.ConvertibleTo(edType) ||
			typ.ConvertibleTo(secpType) ||
			typ.ConvertibleTo(timeType) ||
			typ.ConvertibleTo(cosmosDecType) ||
			typ.ConvertibleTo(cosmosIntType) {
			return str
		}
	}

	return ""
}

func doRecover(err *error) {
	if r := recover(); r != nil {
		if e, ok := r.(error); ok {
			e = errorsmod.Wrap(e, "panicked with error")
			*err = e
			return
		}

		*err = fmt.Errorf("%v", r)
	}
}
