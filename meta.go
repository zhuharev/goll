package goll

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type Meta struct {
	Total   int
	TotalUp int
}

func newMeta(bts []byte) (m Meta, e error) {
	if bts == nil {
		return
	}
	if len(bts) != 8 {
		e = fmt.Errorf("len not 8, %d", len(bts))
		return
	}

	r := bytes.NewReader(bts)
	var tmpUint32 uint32
	e = binary.Read(r, binary.BigEndian, &tmpUint32)
	if e != nil {
		return
	}
	m.Total = int(tmpUint32)

	e = binary.Read(r, binary.BigEndian, &tmpUint32)
	if e != nil {
		return
	}
	m.TotalUp = int(tmpUint32)

	return
}

func (m Meta) encode() ([]byte, error) {
	w := bytes.NewBuffer(nil)
	e := binary.Write(w, binary.BigEndian, int32(m.Total))
	if e != nil {
		return nil, e
	}
	e = binary.Write(w, binary.BigEndian, int32(m.TotalUp))
	if e != nil {
		return nil, e
	}
	return w.Bytes(), nil
}
