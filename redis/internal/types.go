package internal

import "fmt"

type Kind int

const (
	NullKind Kind = iota
	SimpleStringKind
	BulkStringKind
	IntKind
	ArrayKind
	MapKind
	SimpleErrorKind
)

func (k Kind) String() string {
	switch k {
	case NullKind:
		return "Null"
	case SimpleStringKind:
		return "SimpleString"
	case BulkStringKind:
		return "BulkString"
	case IntKind:
		return "Int"
	case ArrayKind:
		return "Array"
	case SimpleErrorKind:
		return "SimpleError"
	case MapKind:
		return "Map"
	default:
		return "Unknown"
	}
}

type Data struct {
	kind  Kind
	value interface{}
}

func (d Data) String() string {
	switch d.kind {
	case NullKind:
		return "Data{Null}"
	case SimpleStringKind:
		return fmt.Sprintf("Data{%q}", d.value)
	case BulkStringKind:
		return fmt.Sprintf("Data{%q}", d.value)
	case IntKind:
		return fmt.Sprintf("Data{%d}", d.value)
	case ArrayKind:
		return fmt.Sprintf("Data{%v}", d.value)
	case MapKind:
		return fmt.Sprintf("Data(%v)", d.value)
	case SimpleErrorKind:
		return fmt.Sprintf("Data{Error: %q}", d.value)
	default:
		return "Unknown"
	}
}

func (d Data) GetKind() Kind {
	return d.kind
}

func (d Data) GetString() (string, error) {
	if !(d.kind == SimpleStringKind || d.kind == BulkStringKind || d.kind == SimpleErrorKind) {
		return "", fmt.Errorf("cannot GetString of kind: %s", d.kind)
	}
	s, ok := d.value.(string)
	if !ok {
		return "", fmt.Errorf("error value is not a string: %v", d.value)
	}
	return s, nil
}

func (d Data) GetInt() (int, error) {
	if d.kind != IntKind {
		return 0, fmt.Errorf("cannot GetInt of kind: %s", d.kind)
	}
	i, ok := d.value.(int)
	if !ok {
		return 0, fmt.Errorf("error value is not an int: %v", d.value)
	}
	return i, nil
}

// TODO: Return pointer to array of pointers
func (d Data) GetArray() ([]Data, error) {
	if d.kind != ArrayKind {
		return nil, fmt.Errorf("cannot GetArray of kind: %s", d.kind)
	}
	a, ok := d.value.([]Data)
	if !ok {
		return nil, fmt.Errorf("error value is not an array: %v", d.value)
	}
	return a, nil
}

func (d Data) GetMap() (map[Data]Data, error) {
	if d.kind != MapKind {
		return nil, fmt.Errorf("cannot GetMap of kind: %s", d.kind)
	}
	m, ok := d.value.(map[Data]Data)
	if !ok {
		return nil, fmt.Errorf("error value is not a map: %v", d.value)
	}
	return m, nil
}

func NewSimpleStringData(s string) *Data {
	return &Data{kind: SimpleStringKind, value: s}
}

func NewBulkStringData(s string) *Data {
	return &Data{kind: BulkStringKind, value: s}
}

func NewIntData(i int) *Data {
	return &Data{kind: IntKind, value: i}
}

func NewArrayData(a []Data) *Data {
	return &Data{kind: ArrayKind, value: a}
}

func NewMapData(m map[Data]Data) *Data {
	return &Data{kind: MapKind, value: m}
}

func NewSimpleError(e string) *Data {
	return &Data{kind: SimpleErrorKind, value: e}
}

func NewNullData() *Data {
	return &Data{kind: NullKind}
}
