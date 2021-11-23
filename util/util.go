package util

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"time"
	"unicode/utf8"

	"../cnutil"
	"../rpc"
)

var Diff1 = StringToBig("0xFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF")

func StringToBig(h string) *big.Int {
	n := new(big.Int)
	n.SetString(h, 0)
	return n
}

func MakeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func GetTargetHex(diff int64) string {
	padded := make([]byte, 32)

	diffBuff := new(big.Int).Div(Diff1, big.NewInt(diff)).Bytes()
	copy(padded[32-len(diffBuff):], diffBuff)
	buff := padded[0:4]
	targetHex := hex.EncodeToString(ReverseBytes(buff))
	return targetHex
}

func GetHashDifficulty(hashBytes []byte) (*big.Int, bool) {
	diff := new(big.Int)
	diff.SetBytes(ReverseBytes(hashBytes))

	// Check for broken result, empty string or zero hex value
	if diff.Cmp(new(big.Int)) == 0 {
		return nil, false
	}
	return diff.Div(Diff1, diff), true
}

func ValidateAddress_Keva(r *rpc.RPCClient, addr string, checkIsMine bool) bool {
	rpcResp, err := r.ValidateAddress(addr)
	if err != nil {
		return false
	}
	var reply *rpc.ValidateAddressReply
	if rpcResp.Result != nil {
		err = json.Unmarshal(*rpcResp.Result, &reply)
		if err != nil {
			return false
		}
		if checkIsMine {
			return reply.IsMine
		}
		return reply.IsValid
	}
	return false
}

func ValidateAddress(addy string, poolAddy string) bool {
	if len(addy) != len(poolAddy) {
		return false
	}
	prefix, _ := utf8.DecodeRuneInString(addy)
	poolPrefix, _ := utf8.DecodeRuneInString(poolAddy)
	if prefix != poolPrefix {
		return false
	}
	return cnutil.ValidateAddress(addy)
}

func ReverseBytes(src []byte) []byte {
	dst := make([]byte, len(src))
	for i := len(src); i > 0; i-- {
		dst[len(src)-i] = src[i-1]
	}
	return dst
}
