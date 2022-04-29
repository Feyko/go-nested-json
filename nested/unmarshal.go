package nested

func Unmarshal(data []byte, v any) error {
	unmarshaler := NewUnmarshaler(v)
	return unmarshaler.UnmarshalJSON(data)
}

func NewUnmarshaler(v any) Unmarshaler {
	return Unmarshaler{out: v}
}

type Unmarshaler struct {
	data []byte
	out  any
}

func (u Unmarshaler) UnmarshalJSON(data []byte) error {
	u.data = data
	return nil
}
