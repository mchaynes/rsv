package rsv

import "fmt"

type ErrFailedToParse struct {
	Value string
	Err   error
}

func (e ErrFailedToParse) Error() string {
	return fmt.Sprintf("%q cannot be parsed: %v", e.Value, e.Err)
}

type errIdxTagNotSet struct{}

func (e errIdxTagNotSet) Error() string {
	return "idx tag not set"
}
