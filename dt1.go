// Copyright 2019-2024 go-sccp authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package sccp

import (
	"fmt"
	"io"
)

// DT1 represents a SCCP Message Data form 1 (DT1).
type DT1 struct {
	Type                   MsgType
	DestinationLocalRef    []byte
	SegmentingReassembling byte
	Ptr1                   uint8
	DataLength             uint8
	Data                   []byte
}

// NewDT1 creates a new DT1.
func NewDT1(destinationLocalRef []byte, segmentingReassembling byte, data []byte) *DT1 {
	u := &DT1{
		Type:                   MsgTypeDT1,
		DestinationLocalRef:    destinationLocalRef,
		SegmentingReassembling: segmentingReassembling,
		Ptr1:                   1,
		Data:                   data,
	}
	u.SetLength()

	return u
}

// MarshalBinary returns the byte sequence generated from a DT1 instance.
func (d *DT1) MarshalBinary() ([]byte, error) {
	b := make([]byte, d.MarshalLen())
	if err := d.MarshalTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// MarshalTo puts the byte sequence in the byte array given as b.
// SCCP is dependent on the Pointers when serializing, which means that it might fail when invalid Pointers are set.
func (d *DT1) MarshalTo(b []byte) error {
	l := len(b)
	if l < 7 {
		return io.ErrUnexpectedEOF
	}

	b[0] = uint8(d.Type)
	copy(b[1:4], d.DestinationLocalRef)
	b[4] = d.SegmentingReassembling
	b[5] = d.Ptr1
	if l < int(d.Ptr1) {
		return io.ErrUnexpectedEOF
	}

	b[d.Ptr1+5] = d.DataLength

	// Succeed if the rest of buffer is longer than u.DataLength.
	if offset := int(d.Ptr1 + 6); len(b[offset:]) >= int(d.DataLength) {
		copy(b[offset:], d.Data)
		return nil
	}

	return io.ErrUnexpectedEOF
}

// ParseDT1 decodes given byte sequence as a SCCP DT1.
func ParseDT1(b []byte) (*DT1, error) {
	u := &DT1{}
	if err := u.UnmarshalBinary(b); err != nil {
		return nil, err
	}

	return u, nil
}

// UnmarshalBinary sets the values retrieved from byte sequence in a SCCP DT1.
func (d *DT1) UnmarshalBinary(b []byte) error {
	l := len(b)
	if l <= 7 {
		return io.ErrUnexpectedEOF
	}

	d.Type = MsgType(b[0])
	d.DestinationLocalRef = b[1:4]
	d.SegmentingReassembling = b[4]
	d.Ptr1 = b[5]
	if l < int(d.Ptr1) {
		return io.ErrUnexpectedEOF
	}

	// Succeed if the rest of buffer is longer than u.DataLength.
	d.DataLength = b[int(d.Ptr1+5)]
	if offset, dataLen := int(d.Ptr1+6), int(d.DataLength); l >= offset+dataLen {
		d.Data = b[offset : offset+dataLen]
		return nil
	}

	return io.ErrUnexpectedEOF
}

// MarshalLen returns the serial length.
func (d *DT1) MarshalLen() int {
	l := 7
	l += len(d.Data)

	return l
}

// SetLength sets the length in Length field.
func (d *DT1) SetLength() {
	d.DataLength = uint8(len(d.Data))
}

// String returns the DT1 values in human readable format.
func (d *DT1) String() string {
	return fmt.Sprintf("{Type: %d, DestinationLocalRef: %v, SegmentingReassembling: %v, DataLength: %d, Data: %x}",
		d.Type,
		d.DestinationLocalRef,
		d.SegmentingReassembling,
		d.DataLength,
		d.Data,
	)
}

// MessageType returns the Message Type in int.
func (d *DT1) MessageType() MsgType {
	return MsgTypeDT1
}

// MessageTypeName returns the Message Type in string.
func (d *DT1) MessageTypeName() string {
	return "DT1"
}
