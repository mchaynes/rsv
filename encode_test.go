package rsv

import (
	"testing"

	"github.com/mchaynes/rsv/internal"

	"github.com/stretchr/testify/assert"
)

func TestMarshalRow(t *testing.T) {
	type A struct {
		Zero string  `idx:"0"`
		One  int     `idx:"1"`
		Two  *string `idx:"2"`
	}
	input := A{
		Zero: "0",
		One:  1,
		Two:  internal.StringPtr("2"),
	}
	expected := []interface{}{"0", int64(1), "2"}
	actual, err := MarshalRow(&input)
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)

}
