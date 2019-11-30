package rsv

import (
	"testing"

	"github.com/mchaynes/rsv/internal"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalRow_Embedded(t *testing.T) {
	type C struct {
		D string `idx:"0"`
		E *int   `idx:"1"`
	}
	type B struct {
		C C
	}
	type A struct {
		B B
	}
	expected := A{
		B: B{
			C: C{
				D: "D",
				E: internal.IntPtr(10),
			},
		},
	}
	row := []string{"D", "10"}
	actual := A{}
	err := UnmarshalRow(row, &actual)
	assert.NoError(t, err, "embedded struct shouldn't have returned an error")
	assert.Equal(t, expected, actual)
}

func TestUnmarshalRow_Numbers(t *testing.T) {
	type A struct {
		// Mark F64 & F32 to prove that ptrs to floats works
		F64 *float64 `idx:"0"`
		F32 *float32 `idx:"1"`

		I   int   `idx:"2"`
		I64 int64 `idx:"3"`
		I32 int32 `idx:"4"`

		UI   uint   `idx:"5"`
		UI8  uint8  `idx:"6"`
		UI16 uint16 `idx:"7"`
		UI32 uint32 `idx:"8"`
		UI64 uint64 `idx:"9"`
	}
	expected := A{
		F64:  internal.Float64Ptr(64),
		F32:  internal.Float32Ptr(32),
		I:    1,
		I64:  64,
		I32:  32,
		UI:   1,
		UI8:  8,
		UI16: 16,
		UI32: 32,
		UI64: 64,
	}
	row := []string{"64", "32", "1", "64", "32", "1", "8", "16", "32", "64"}
	actual := A{}
	err := UnmarshalRow(row, &actual)
	assert.NoError(t, err, "should have no errors")
	assert.Equal(t, expected, actual)
}

func TestUnmarshalRow_Omitempty(t *testing.T) {
	type A struct {
		F64 *float64 `idx:"0,omitempty"`
		S   *string  `idx:"1,omitempty"`
	}
	expected := A{
		F64: internal.Float64Ptr(64),
		S:   internal.StringPtr("s"),
	}
	row := []string{"64", "s"}
	actual := A{}
	err := UnmarshalRow(row, &actual)
	assert.NoError(t, err, "should have no errors")
	assert.Equal(t, expected, actual)

	// Make all fields nil
	expected = A{}
	row = []string{"", ""}
	actual = A{}
	err = UnmarshalRow(row, &actual)
	assert.NoError(t, err, "should have no errors")
	assert.Equal(t, expected, actual)
}

func TestUnmarshalRow_Panics(t *testing.T) {
	defer func() {
		if i := recover(); i == nil {
			t.Error("test should have panicked")
		}
	}()
	_ = UnmarshalRow([]string{}, struct{}{})
}
