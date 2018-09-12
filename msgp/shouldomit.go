package msgp

import "time"

// Implement this interface to enable struct type to be omitted by `omitempty`.
type OmitEmptyAware interface {
	// This method should return true if the object should not be serialized with `omitempty` marker.
	ShouldOmitAsEmpty() bool
}

func ShouldOmitAsEmptyBytes(v []byte) bool {
	return len(v) == 0
}
func ShouldOmitAsEmptyString(v string) bool {
	return v == ""
}
func ShouldOmitAsEmptyFloat32(v float32) bool {
	return v == 0
}
func ShouldOmitAsEmptyFloat64(v float64) bool {
	return v == 0
}
func ShouldOmitAsEmptyComplex64(v complex64) bool {
	return v == 0
}
func ShouldOmitAsEmptyComplex128(v complex128) bool {
	return v == 0
}
func ShouldOmitAsEmptyIntf(v interface{}) bool {
	return v == nil
}
func ShouldOmitAsEmptyUint(v uint) bool {
	return v == 0
}
func ShouldOmitAsEmptyUint8(v uint8) bool {
	return v == 0
}
func ShouldOmitAsEmptyUint16(v uint16) bool {
	return v == 0
}
func ShouldOmitAsEmptyUint32(v uint32) bool {
	return v == 0
}
func ShouldOmitAsEmptyUint64(v uint64) bool {
	return v == 0
}
func ShouldOmitAsEmptyByte(v byte) bool {
	return v == 0
}
func ShouldOmitAsEmptyInt(v int) bool {
	return v == 0
}
func ShouldOmitAsEmptyInt8(v int8) bool {
	return v == 0
}
func ShouldOmitAsEmptyInt16(v int16) bool {
	return v == 0
}
func ShouldOmitAsEmptyInt32(v int32) bool {
	return v == 0
}
func ShouldOmitAsEmptyInt64(v int64) bool {
	return v == 0
}
func ShouldOmitAsEmptyBool(v bool) bool {
	return v == false
}
func ShouldOmitAsEmptyTime(v time.Time) bool {
	return v.IsZero()
}
func ShouldOmitAsEmptyExtension(v *RawExtension) bool {
	return v == nil
}
