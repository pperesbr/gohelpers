package payload

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
)

type FieldType string

const (
	TypeString     FieldType = "string"
	TypeInt        FieldType = "int"
	TypeFloat      FieldType = "float"
	TypeBool       FieldType = "bool"
	TypeUUID       FieldType = "uuid"
	TypeStatic     FieldType = "static"
	TypeObjectList FieldType = "objectList"
	TypeFixedChar  FieldType = "fixedChar"
)

type FieldDef struct {
	Name       string
	Type       FieldType
	FakeQuery  string
	StaticVal  interface{}
	Range      []int
	FloatRange []float64
	Length     int
	Children   []FieldDef
}

func Generate(fields []FieldDef) ([]byte, error) {
	data := buildMap(fields)
	return json.MarshalIndent(data, "", "  ")
}

func GenerateReader(fields []FieldDef) (io.Reader, error) {
	data, err := Generate(fields)
	if err != nil {
		return nil, err
	}
	return bytes.NewBuffer(data), nil
}

func buildMap(fields []FieldDef) map[string]interface{} {
	output := make(map[string]interface{})

	for _, f := range fields {

		if f.StaticVal != nil && f.Type == "" {
			f.Type = TypeStatic
		}

		switch f.Type {

		case TypeUUID:
			output[f.Name] = uuid.NewString()

		case TypeFloat:
			min, max := 0.0, 100.0
			if len(f.FloatRange) == 2 {
				min, max = f.FloatRange[0], f.FloatRange[1]
			}
			output[f.Name] = gofakeit.Float64Range(min, max)

		case TypeFixedChar:
			if f.Length > 0 {
				output[f.Name] = gofakeit.LetterN(uint(f.Length))
			} else {
				output[f.Name] = gofakeit.LetterN(3)
			}

		case TypeString:
			output[f.Name] = gofakeit.Generate(f.FakeQuery)

		case TypeInt:
			min, max := 0, 100
			if len(f.Range) == 2 {
				min, max = f.Range[0], f.Range[1]
			}
			output[f.Name] = gofakeit.Number(min, max)

		case TypeBool:
			output[f.Name] = gofakeit.Bool()

		case TypeStatic:
			output[f.Name] = f.StaticVal

		case TypeObjectList:
			min, max := 1, 3
			if len(f.Range) == 2 {
				min, max = f.Range[0], f.Range[1]
			}
			count := gofakeit.Number(min, max)
			list := make([]map[string]interface{}, count)
			for i := 0; i < count; i++ {
				list[i] = buildMap(f.Children)
			}
			output[f.Name] = list
		default:
			if f.FakeQuery != "" {
				output[f.Name] = gofakeit.Generate(f.FakeQuery)
			}
		}
	}
	return output
}
