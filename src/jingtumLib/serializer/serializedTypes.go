/**
 *
 * 各数据类型的序列化实现
 *
 * @FileName: serializedTypes.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-11 10:44:32
 * @UpdateTime: 2018-07-11 10:44:54
 * Copyright@2018 版权所有
 */

package serializer

import (
	"bytes"
	"container/list"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"jingtumLib/constant"
	jtUtils "jingtumLib/utils"
)

var (
	CURRENCY_NAME_LEN  int = 3
	CURRENCY_NAME_LEN2 int = 6
	typeBoundary       int = 0xff
	typeEnd            int = 0x00
	typeAccount        int = 0x01
	typeCurrency       int = 0x10
	typeIssuer         int = 0x20
)

type ISerializedType interface {
	//Serialize( val interface{}, noMarker bool)
	Serialize(so *Serializer, val interface{}, noMarker bool)
	Parse(so *Serializer) interface{}
}

type PathComputed struct {
	Currency string `json:"currency"`
	Issuer   string `json:"issuer"`
	Value    string `json:"value"`
	Account  string `json:"account"`
	Type     int    `json:"type"`
	TypeHex  string `json:"type_hex"`
}

type PathData struct {
	PathsComputed [][]PathComputed
	Choice        interface{}
}

type SerializedInt16 struct {
}

type SerializedInt32 struct {
}

type SerializedInt64 struct {
}

type SerializedMemo struct {
}

type SerializedArg struct {
}

type SerializedHash128 struct {
}

type SerializedHash256 struct {
}

type SerializedAmount struct {
}

type SerializedCurrency struct {
}

type SerializedObject struct {
}

type SerializedArray struct {
}

type SerializedHash160 struct {
}

type SerializedPathSet struct {
}

type SerializedVector256 struct {
}

type SerializedVariableLength struct {
}

type SerializedAccount struct {
}

func (serInt8 SerializedInt8) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if vuit8, ok := val.(uint8); ok {
		so.Append(jtUtils.GetBytes(vuit8))
	} else if vint, ok := val.(int); ok {
		if vint >= math.MaxUint8 || vint < 0 {
			so.err = errors.New(fmt.Sprintf("Value out of bounds %d.", vint))
			return
		}

		so.Append(jtUtils.GetBytes(val))
	} else {
		so.err = errors.New(fmt.Sprintf("Serialize int8 type error %T, %v.", val, val))
	}
}

func (serInt8 SerializedInt8) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serInt16 SerializedInt16) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if vuint16, ok := val.(uint16); ok {
		so.Append(jtUtils.GetBytes(vuint16))
	} else if vint, ok := val.(int); ok {
		if vint >= math.MaxUint16 || vint < 0 {
			so.err = errors.New(fmt.Sprintf("Value out of bounds %d.", vint))
			return
		}

		so.Append(jtUtils.GetBytes(val))
	} else {
		so.err = errors.New(fmt.Sprintf("Serialize int16 type error %T, %v.", val, val))
	}
}

func (serInt16 SerializedInt16) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serInt32 SerializedInt32) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if vuint32, ok := val.(uint32); ok {
		so.Append(jtUtils.GetBytes(vuint32))
	} else if vint, ok := val.(int); ok {
		if vint >= math.MaxUint32 || vint < 0 {
			so.err = errors.New(fmt.Sprintf("Value out of bounds %d.", vint))
			return
		}
		so.Append(jtUtils.GetBytes(val))
	} else {
		so.err = errors.New(fmt.Sprintf("Serialize int16 type error %T, %v.", val, val))
	}
}

func (serInt32 SerializedInt32) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serInt64 SerializedInt64) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if number, ok := val.(uint64); ok {
		so.Append(jtUtils.GetBytes(number))
		return
	}

	if str, ok := val.(string); ok {
		if !jtUtils.IsHexString(str) {
			so.err = errors.New(fmt.Sprintf("Invalid hex string %s", str))
			return
		}

		if len(str) > 16 {
			panic(fmt.Sprintf("Int64 is too large %v", str))
		}

		b := bytes.NewBufferString("")

		for b.Len() < 16-len(str) {
			b.WriteString("0")
		}

		b.WriteString(str)

		SerializeHex(so, b.String(), true)
		return
	}

	so.err = errors.New(fmt.Sprintf("Invalid type for Int64 %T, %v", val, val))
}

func (serInt64 SerializedInt64) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serMemo SerializedMemo) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serMemo SerializedMemo) Serialize(so *Serializer, val interface{}, noMarker bool) {
	memo, ok := val.(*MemoDataInfo)
	if !ok {
		so.err = errors.New(fmt.Sprintf("Serialize Memo type error %T.", val))
		return
	}

	fileds := jtUtils.GetFieldNames(val)
	for i := 0; i < len(fileds); i++ {
		_, ok := constant.INVERSE_FIELDS_MAP[fileds[i]]
		if !ok {
			so.err = errors.New(fmt.Sprintf("JSON contains unknown field : %s.", fileds[i]))
			return
		}
	}

	jtUtils.SortByFieldName(fileds)

	isJson := memo.MemoFormat == "json"

	for _, fn := range fileds {
		value := jtUtils.GetFieldValue(val, fn)
		switch fn {
		case "MemoType", "MemoFormat":
			value = jtUtils.StringToHex(value.(string))
		case "MemoData":
			if _, ok := value.(string); ok {
				value = jtUtils.StringToHex(value.(string))
				break
			}
			if isJson {
				mjson, _ := json.Marshal(value)
				value = jtUtils.StringToHex(string(mjson))
				break
			}
			so.err = errors.New(fmt.Sprintf("MemoData can only be a JSON object with a valid json MemoFormat. %v", value))
			return
		}

		Serialize(so, fn, value)
	}

	if !noMarker {
		STInt8.Serialize(so, 0xe1, false)
	}
}

func (serArg SerializedArg) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serArg SerializedArg) Serialize(so *Serializer, val interface{}, noMarker bool) {
	fileds := jtUtils.GetFieldNames(val)
	for i := 0; i < len(fileds); i++ {
		kvp := constant.INVERSE_FIELDS_MAP[fileds[i]]
		if kvp == nil {
			panic(fmt.Sprintf("JSON contains unknown field: %v", fileds[i]))
		}
	}
	jtUtils.SortByFieldName(fileds)

	for _, fn := range fileds {
		if fn == "Parameter" {
			break
		}
		value := jtUtils.GetFieldValue(val, fn)
		Serialize(so, fn, value)
	}

	if !noMarker {
		STInt8.Serialize(so, 0xe1, false)
	}
}

func (serHash128 SerializedHash128) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serHash128 SerializedHash128) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if v, ok := val.(string); ok && jtUtils.MatchString("^[0-9A-F]{0,32}$", v) && len(v) <= 32 {
		SerializeHex(so, v, true)
		return
	}
}

func (serHash256 SerializedHash256) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serHash256 SerializedHash256) Serialize(so *Serializer, val interface{}, noMarker bool) {
	if v, ok := val.(string); ok && jtUtils.MatchString("^[0-9A-F]{0,32}$", v) && len(v) <= 64 {
		SerializeHex(so, v, true)
		return
	}
}

func (serAmount SerializedAmount) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serAmount SerializedAmount) Serialize(so *Serializer, val interface{}, noMarker bool) {
	amount := val.(TumAmount)
	if !amount.IsValid {
		panic("Not a valid Amount object.")
	}

	if amount.IsNative {
		valueHex := hex.EncodeToString(amount.Value.Bytes())

		if len(valueHex) > 16 {
			panic("Amount value out of bounds.")
		}
		b := bytes.NewBufferString("")
		for b.Len() < 16 {
			b.WriteString("0")
		}
		b.WriteString(valueHex)

		valueBytes, err := jtUtils.HexToBytes(b.String())
		if err != nil {
			panic("Hex to bytes error.")
		}

		valueBytes[0] &= 0x3f

		if amount.IsNegative {
			valueBytes[0] |= 0x40
		}

		so.Append(valueBytes)
	} else {
		//For other non-native currency
		//1. Serialize the currency value with offset
		//Put offset
		var hi, lo int64 = 0, 0
		hi |= 1 << 31
		if !amount.IsZeroM() {
			// Second bit: non-negative?
			if amount.IsNegative {
				hi |= 1 << 30
			}
			// Next eight bits: offset/exponent
			hi |= ((int64(97) + int64(amount.Offset)) & 0xff) << 22
			// Remaining 54 bits: mantissa
			hi |= (amount.Value.Int64() >> 32) & 0x3fffff
			lo = amount.Value.Int64() & 0xffffffff
		}

		// Convert from a bitArray to an array of bytes.
		arr := []int64{hi, lo}
		l := len(arr)
		var bl int64
		if l == 0 {
			bl = 0
		} else {
			x := arr[l-1]
			roundX := x / 0x10000000000
			if roundX == 0 {
				roundX = 32
			}
			bl = int64((l-1)*32) + int64(roundX)
		}

		var tmparray []byte
		var tmp int64 = 0
		for i := 0; int64(i) < int64(bl/8); i++ {
			if (i & 3) == 0 {
				tmp = arr[i/4]
			}

			tmparray = append(tmparray, byte(tmp>>24))
			tmp <<= 8
		}

		if len(tmparray) > 8 {
			panic("Invalid byte array length in AMOUNT value representation")
		}

		so.Append(tmparray)
		tumBytes := amount.TumToBytes()
		so.Append(tumBytes)
		so.Append(jtUtils.DecodeAddress(amount.Issuer))
	}
}

func (serCurrency SerializedCurrency) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serCurrency SerializedCurrency) Serialize(so *Serializer, val interface{}, noMarker bool) {
	currencty := val.(string)
	so.Append(serCurrency.fromJsonToBytes(currencty))
}

func (serCurrency SerializedCurrency) fromJsonToBytes(currencty string) []byte {
	var result []byte
	if currencty != "" {
		if jtUtils.IsHexString(currencty) && len(currencty) == 40 {
			var err error
			result, err = jtUtils.HexToBytes(currencty)

			if err != nil {
				panic("Invalid currencty.")
			}

		} else if jtUtils.IsValidCurrency(currencty) {
			if len(currencty) >= CURRENCY_NAME_LEN && len(currencty) <= CURRENCY_NAME_LEN2 {
				var end = 14
				var clen = len(currencty) - 1
				for x := clen; x >= 0; x-- {
					result[end-x] = byte(currencty[clen-x] & 0xff)
				}
			}
		} else {
			panic(fmt.Sprintf("Input tum code invalid %v", currencty))
		}
	} else {
		panic(fmt.Sprintf("Input tum code invalid %v", currencty))
	}

	return result
}

func (serObject SerializedObject) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serObject SerializedObject) Serialize(so *Serializer, val interface{}, noMarker bool) {
	txData, ok := val.(map[string]interface{})
	if !ok {
		so.err = errors.New("Serialive object type must be map[string]interface{}.")
		return
	}

	var fieldNames []string

	for k, _ := range txData {
		_, ok := constant.INVERSE_FIELDS_MAP[k]
		if !ok {
			so.err = errors.New(fmt.Sprintf("Not fund field name %s.", k))
			return
		}

		fieldNames = append(fieldNames, k)
	}

	jtUtils.SortByFieldName(fieldNames)

	for _, field := range fieldNames {
		value := txData[field]
		if value == nil {
			continue
		}

		Serialize(so, field, value)
	}

	if !noMarker {
		STInt8.Serialize(so, 0xe1, false)
	}
}

func (serArray SerializedArray) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serArray SerializedArray) Serialize(so *Serializer, val interface{}, noMarker bool) {
	array, ok := val.(*list.List)
	if !ok {
		so.err = errors.New(fmt.Sprintf("Serialize array type error %T", val))
		return
	}

	for e := array.Front(); e != nil; e = e.Next() {
		fields := jtUtils.GetFieldNames(e.Value)

		if len(fields) != 1 {
			so.err = errors.New("Cannot serialize an array containing non-single-key objects.")
			return
		}

		value := jtUtils.GetFieldValue(e.Value, fields[0])

		Serialize(so, fields[0], value)
	}

	STInt8.Serialize(so, 0xf1, false)
}

func (serHash160 SerializedHash160) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serHash160 SerializedHash160) Serialize(so *Serializer, val interface{}, noMarker bool) {
	valStr := val.(string)
	SerializeHex(so, valStr, true)
}

func (serPathSet SerializedPathSet) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serPathSet SerializedPathSet) Serialize(so *Serializer, val interface{}, noMarker bool) {
	path := val.([][]PathComputed)
	for i := 0; i < len(path); i++ {
		if i > 0 {
			STInt8.Serialize(so, typeBoundary, false)
		}

		pathes := path[i]
		for j, l2 := 0, len(pathes); j < l2; j++ {
			entry := pathes[j]
			typev := 0
			if entry.Account != "" {
				typev |= typeAccount
			}

			if entry.Currency != "" {
				typev |= typeCurrency
			}

			if entry.Issuer != "" {
				typev |= typeIssuer
			}

			STInt8.Serialize(so, typev, false)

			if entry.Account != "" {
				so.Append(jtUtils.DecodeAddress(entry.Account))
			}

			if entry.Currency != "" {
				sc := new(SerializedCurrency)
				so.Append(sc.fromJsonToBytes(entry.Currency))
			}

			if entry.Issuer != "" {
				so.Append(jtUtils.DecodeAddress(entry.Issuer))
			}
		}
	}

	STInt8.Serialize(so, typeEnd, false)
}

func (serVector256 SerializedVector256) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serVector256 SerializedVector256) Serialize(so *Serializer, val interface{}, noMarker bool) {
	array := val.([]string)
	SerializeVarint(so, uint(len(array)*32))

	for _, v := range array {
		STHash256.Serialize(so, v, false)
	}
}

func (serVL SerializedVariableLength) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serVL SerializedVariableLength) Serialize(so *Serializer, val interface{}, noMarker bool) {
	vl, ok := val.(string)
	if !ok {
		so.err = errors.New(fmt.Sprintf("Serialized Variable type error %T.", val))
		return
	}
	SerializeHex(so, vl, false)
}

func (serAccount SerializedAccount) Parse(so *Serializer) interface{} {
	return errors.New("Not implemented error.")
}

func (serAccount SerializedAccount) Serialize(so *Serializer, val interface{}, noMarker bool) {
	bytes := jtUtils.DecodeAddress(val.(string))
	SerializeVarint(so, uint(len(bytes)))
	so.Append(bytes)
}
