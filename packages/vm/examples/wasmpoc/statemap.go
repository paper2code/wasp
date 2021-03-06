package wasmpoc

import (
	"github.com/iotaledger/wart/host/interfaces"
	"github.com/iotaledger/wart/host/interfaces/objtype"
	"github.com/iotaledger/wasp/packages/kv"
	"github.com/iotaledger/wasp/packages/util"
)

type StateMap struct {
	vm     *wasmVMPocProcessor
	items  *kv.MustDictionary
	types  map[int32]int32
}

func NewStateMap(h *wasmVMPocProcessor, items  *kv.MustDictionary) *StateMap {
	return &StateMap{vm: h, items: items, types: make(map[int32]int32)}
}

func (m *StateMap) GetInt(keyId int32) int64 {
	if !m.valid(keyId, objtype.OBJTYPE_INT) {
		return 0
	}
	key := []byte(m.vm.GetKey(keyId))
	value, _ := kv.DecodeInt64(m.items.GetAt(key))
	return value
}

func (m *StateMap) GetLength() int32 {
	m.vm.SetError("Invalid length")
	return 0
}

func (m *StateMap) GetObjectId(keyId int32, typeId int32) int32 {
	m.vm.SetError("Invalid access")
	return 0
}

func (m *StateMap) GetString(keyId int32) string {
	if !m.valid(keyId, objtype.OBJTYPE_STRING) {
		return ""
	}
	key := []byte(m.vm.GetKey(keyId))
	return string(m.items.GetAt(key))
}

func (m *StateMap) SetInt(keyId int32, value int64) {
	if keyId == interfaces.KeyLength {
		m.vm.SetError("Invalid clear")
		return
	}
	if !m.valid(keyId, objtype.OBJTYPE_INT) {
		return
	}
	key := []byte(m.vm.GetKey(keyId))
	m.items.SetAt(key, util.Uint64To8Bytes(uint64(value)))
}

func (m *StateMap) SetString(keyId int32, value string) {
	if !m.valid(keyId, objtype.OBJTYPE_STRING) {
		return
	}
	key := []byte(m.vm.GetKey(keyId))
	m.items.SetAt(key, []byte(value))
}

func (m *StateMap) valid(keyId int32, typeId int32) bool {
	fieldType, ok := m.types[keyId]
	if !ok {
		m.types[keyId] = typeId
		return true
	}
	if fieldType != typeId {
		m.vm.SetError("Invalid access")
		return false
	}
	return true
}
