package types

import "github.com/cosmos/cosmos-sdk/types/kv"

const (
	// ModuleName is the module name
	ModuleName = "gashub"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey defines the module's message routing key
	RouterKey = ModuleName
)

var (
	ParamsKey = []byte{0x00} // key for x/gashub module params

	MsgGasParamsPrefix = []byte{0x01} // key for msg gas params
)

func GetMsgTypeUrl(key []byte) string {
	// key is in the format:
	// 0x02<urlLen (1 Byte)><url_Bytes>

	// Remove prefix and address length.
	kv.AssertKeyAtLeastLength(key, 3)
	url := key[2:]
	kv.AssertKeyLength(url, int(key[1]))

	return string(url)
}

func GetMsgGasParamsKey(msgTypeUrl string) []byte {
	return append(MsgGasParamsPrefix, LengthPrefix([]byte(msgTypeUrl))...)
}

func LengthPrefix(bz []byte) []byte {
	bzLen := len(bz)
	if bzLen == 0 {
		return bz
	}

	return append([]byte{byte(bzLen)}, bz...)
}
