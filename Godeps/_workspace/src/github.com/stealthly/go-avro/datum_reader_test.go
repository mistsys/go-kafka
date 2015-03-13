package avro

import (
	"fmt"
	"testing"
)

//primitives
type Primitive struct {
	BooleanField bool
	IntField     int32
	LongField    int64
	FloatField   float32
	DoubleField  float64
	BytesField   []byte
	StringField  string
	NullField    interface{}
}

//TODO replace with encoder <-> decoder tests when decoder is available
//primitive values predefined test data
var (
	primitive_bool   bool        = true
	primitive_int    int32       = 7498
	primitive_long   int64       = 7921326876135578931
	primitive_float  float32     = 87612736.5124367
	primitive_double float64     = 98671578.12563891
	primitive_bytes  []byte      = []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09}
	primitive_string string      = "A very long and cute string here!"
	primitive_null   interface{} = nil
)

func TestPrimitiveBinding(t *testing.T) {
	datumReader := NewSpecificDatumReader()
	reader, err := NewDataFileReader("test/primitives.avro", datumReader)
	if err != nil {
		t.Fatal(err)
	}
	for {
		p := &Primitive{}
		ok, err := reader.Next(p)
		if !ok {
			if err != nil {
				t.Fatal(err)
			}
			break
		} else {
			assert(t, p.BooleanField, primitive_bool)
			assert(t, p.IntField, primitive_int)
			assert(t, p.LongField, primitive_long)
			assert(t, p.FloatField, primitive_float)
			assert(t, p.DoubleField, primitive_double)
			assert(t, p.BytesField, primitive_bytes)
			assert(t, p.StringField, primitive_string)
			assert(t, p.NullField, primitive_null)
		}
	}
}

//complex
type complex struct {
	StringArray []string
	LongArray   []int64
	EnumField   GenericEnum
	MapOfInts   map[string]int32
	UnionField  string
	FixedField  []byte
	RecordField *testRecord
}

type testRecord struct {
	LongRecordField   int64
	StringRecordField string
	IntRecordField    int32
	FloatRecordField  float32
}

//TODO replace with encoder <-> decoder tests when decoder is available
//predefined test data for complex types
var (
	complex_union         string  = "union value"
	complex_fixed         []byte  = []byte{0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04, 0x01, 0x02, 0x03, 0x04}
	complex_record_long   int64   = 1925639126735
	complex_record_string string  = "I am a test record"
	complex_record_int    int32   = 666
	complex_record_float  float32 = 7171.17
)

func TestComplexBinding(t *testing.T) {
	datumReader := NewSpecificDatumReader()
	reader, err := NewDataFileReader("test/complex.avro", datumReader)
	if err != nil {
		t.Fatal(err)
	}
	for {
		c := &complex{}
		ok, err := reader.Next(c)
		if !ok {
			if err != nil {
				t.Fatal(err)
			}
			break
		} else {
			arrayLength := 5
			if len(c.StringArray) != arrayLength {
				t.Errorf("Expected string array length %d, actual %d", arrayLength, len(c.StringArray))
			}
			for i := 0; i < arrayLength; i++ {
				if c.StringArray[i] != fmt.Sprintf("string%d", i+1) {
					t.Errorf("Invalid string: expected %v, actual %v", fmt.Sprintf("string%d", i+1), c.StringArray[i])
				}
			}

			if len(c.LongArray) != arrayLength {
				t.Errorf("Expected long array length %d, actual %d", arrayLength, len(c.LongArray))
			}
			for i := 0; i < arrayLength; i++ {
				if c.LongArray[i] != int64(i+1) {
					t.Errorf("Invalid long: expected %v, actual %v", i+1, c.LongArray[i])
				}
			}

			enumValues := []string{"A", "B", "C", "D"}
			for i := 0; i < len(enumValues); i++ {
				if enumValues[i] != c.EnumField.Symbols[i] {
					t.Errorf("Invalid enum value in sequence: expected %v, actual %v", enumValues[i], c.EnumField.Symbols[i])
				}
			}

			if c.EnumField.Get() != enumValues[2] {
				t.Errorf("Invalid enum value: expected %v, actual %v", enumValues[2], c.EnumField.Get())
			}

			if len(c.MapOfInts) != arrayLength {
				t.Errorf("Invalid map length: expected %d, actual %d", arrayLength, len(c.MapOfInts))
			}

			for k, v := range c.MapOfInts {
				if k != fmt.Sprintf("key%d", v) {
					t.Errorf("Invalid key for a map value: expected %v, actual %v", fmt.Sprintf("key%d", v), k)
				}
			}

			if c.UnionField != complex_union {
				t.Errorf("Invalid union value: expected %v, actual %v", complex_union, c.UnionField)
			}

			assert(t, c.FixedField, complex_fixed)
			assert(t, c.RecordField.LongRecordField, complex_record_long)
			assert(t, c.RecordField.StringRecordField, complex_record_string)
			assert(t, c.RecordField.IntRecordField, complex_record_int)
			assert(t, c.RecordField.FloatRecordField, complex_record_float)
		}
	}
}

//complex within complex
type complexOfComplex struct {
	ArrayStringArray  [][]string
	RecordArray       []testRecord
	IntOrStringArray  []interface{}
	RecordMap         map[string]testRecord2
	IntOrStringMap    map[string]interface{}
	NullOrRecordUnion *testRecord3
}

type testRecord2 struct {
	DoubleRecordField float64
	FixedRecordField  []byte
}

type testRecord3 struct {
	StringArray     []string
	EnumRecordField GenericEnum
}

func TestComplexOfComplexBinding(t *testing.T) {
	datumReader := NewSpecificDatumReader()
	reader, err := NewDataFileReader("test/complex_of_complex.avro", datumReader)
	if err != nil {
		t.Fatal(err)
	}
	for {
		c := &complexOfComplex{}
		ok, err := reader.Next(c)
		if !ok {
			if err != nil {
				t.Fatal(err)
			}
			break
		} else {
			arrayLength := 5
			if len(c.ArrayStringArray) != arrayLength {
				t.Errorf("Expected array of arrays length %d, actual %d", arrayLength, len(c.ArrayStringArray))
			}

			for i := 0; i < arrayLength; i++ {
				for j := 0; j < arrayLength; j++ {
					if c.ArrayStringArray[i][j] != fmt.Sprintf("string%d%d", i, j) {
						t.Errorf("Expected array element %d, actual %d", fmt.Sprintf("string%d%d", i, j), c.ArrayStringArray[i][j])
					}
				}
			}

			recordArrayLength := 2
			if len(c.RecordArray) != recordArrayLength {
				t.Errorf("Expected record array length %d, actual %d", recordArrayLength, len(c.RecordArray))
			}

			for i := 0; i < recordArrayLength; i++ {
				rec := c.RecordArray[i]

				assert(t, rec.LongRecordField, int64(i))
				assert(t, rec.StringRecordField, fmt.Sprintf("TestRecord%d", i))
				assert(t, rec.IntRecordField, int32(1000+i))
				assert(t, rec.FloatRecordField, float32(i)+0.05)
			}

			intOrString := []interface{}{int32(32), "not an integer", int32(49)}

			if len(c.IntOrStringArray) != len(intOrString) {
				t.Errorf("Expected union array length %d, actual %d", len(intOrString), len(c.IntOrStringArray))
			}

			for i := 0; i < len(intOrString); i++ {
				assert(t, c.IntOrStringArray[i], intOrString[i])
			}

			recordMapLength := 2
			if len(c.RecordMap) != recordMapLength {
				t.Errorf("Expected map length %d, actual %d", recordMapLength, len(c.RecordMap))
			}

			rec1 := c.RecordMap["a key"]
			assert(t, rec1.DoubleRecordField, float64(32.5))
			assert(t, rec1.FixedRecordField, []byte{0x00, 0x01, 0x02, 0x03})
			rec2 := c.RecordMap["another key"]
			assert(t, rec2.DoubleRecordField, float64(33.5))
			assert(t, rec2.FixedRecordField, []byte{0x01, 0x02, 0x03, 0x04})

			stringMapLength := 3
			if len(c.IntOrStringMap) != stringMapLength {
				t.Errorf("Expected string map length %d, actual %d", stringMapLength, len(c.IntOrStringMap))
			}
			assert(t, c.IntOrStringMap["a key"], "a value")
			assert(t, c.IntOrStringMap["one more key"], int32(123))
			assert(t, c.IntOrStringMap["another key"], "another value")

			if len(c.NullOrRecordUnion.StringArray) != arrayLength {
				t.Errorf("Expected record union string array length %d, actual %d", arrayLength, len(c.NullOrRecordUnion.StringArray))
			}
			for i := 0; i < arrayLength; i++ {
				assert(t, c.NullOrRecordUnion.StringArray[i], fmt.Sprintf("%d", i))
			}

			enumValues := []string{"A", "B", "C", "D"}
			for i := 0; i < len(enumValues); i++ {
				if enumValues[i] != c.NullOrRecordUnion.EnumRecordField.Symbols[i] {
					t.Errorf("Invalid enum value in sequence: expected %v, actual %v", enumValues[i], c.NullOrRecordUnion.EnumRecordField.Symbols[i])
				}
			}

			if c.NullOrRecordUnion.EnumRecordField.Get() != enumValues[3] {
				t.Errorf("Invalid enum value: expected %v, actual %v", enumValues[3], c.NullOrRecordUnion.EnumRecordField.Get())
			}
		}
	}
}
