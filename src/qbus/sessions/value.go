package sessions

import "encoding/json"

// Value helps converting json.RawMessage to real data type underly
type Value struct {
	value json.RawMessage
	ok    bool
}

func NewValue(value json.RawMessage) *Value {
	return &Value{value: value}
}

func NewSessionValue(value json.RawMessage, ok bool) *Value {
	return &Value{value: value, ok: ok}
}

func (v *Value) IsExist() bool {
	return v.ok
}

func (v *Value) Value() json.RawMessage {
	return v.value
}

func (v *Value) Unmarshal(vt interface{}) (err error) {
	if !v.ok {
		err = ErrNotFound
		return
	}

	err = json.Unmarshal(v.value, &vt)
	return
}

func (v *Value) String() string {
	b, err := json.Marshal(v.value)
	if err != nil {
		return ""
	}

	return string(b)
}

func (v *Value) Bool() (ok bool, err error) {
	err = v.Unmarshal(&ok)
	if err == nil {
		return
	}

	// fallback to weakness
	switch v.String() {
	case "t", "true", "on", "y", "yes":
		ok = true
		err = nil
		return
	}

	return
}

func (v *Value) Int() (i int, err error) {
	err = v.Unmarshal(&i)
	return
}

func (v *Value) Int8() (i8 int8, err error) {
	err = v.Unmarshal(&i8)
	return
}

func (v *Value) Int16() (i16 int16, err error) {
	err = v.Unmarshal(&i16)
	return
}

func (v *Value) Int32() (i32 int32, err error) {
	err = v.Unmarshal(&i32)
	return
}

func (v *Value) Int64() (i64 int64, err error) {
	err = v.Unmarshal(&i64)
	return
}

func (v *Value) Uint() (u uint, err error) {
	err = v.Unmarshal(&u)
	return
}

func (v *Value) Uint8() (u8 uint8, err error) {
	err = v.Unmarshal(&u8)
	return
}

func (v *Value) Uint16() (u16 uint16, err error) {
	err = v.Unmarshal(&u16)
	return
}

func (v *Value) Uint32() (u32 uint32, err error) {
	err = v.Unmarshal(&u32)
	return
}

func (v *Value) Uint64() (u64 uint64, err error) {
	err = v.Unmarshal(&u64)
	return
}

func (v *Value) Float32() (f32 float32, err error) {
	err = v.Unmarshal(&f32)
	return
}

func (v *Value) Float64() (f64 float64, err error) {
	err = v.Unmarshal(&f64)
	return
}
