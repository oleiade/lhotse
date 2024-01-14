package main

import (
	"errors"
	"fmt"
	"math/rand"
	"strings"

	"github.com/oleiade/gomme"
)

// ByteUnit represents a byte unit.
type ByteUnit int64

const (
	// Byte represents a byte.
	Byte ByteUnit = 1

	// Kilobyte represents a kilobyte.
	Kilobyte = 1024 * Byte

	// Megabyte represents a megabyte.
	Megabyte = 1024 * Kilobyte

	// Gigabyte represents a gigabyte.
	Gigabyte = 1024 * Megabyte
)

// ParseByteUnit parses a byte unit string and returns a ByteUnit.
func ParseByteUnit(u string) ByteUnit {
	switch strings.ToLower(u) {
	case "b":
		return Byte
	case "kb":
		return Kilobyte
	case "mb":
		return Megabyte
	case "gb":
		return Gigabyte
	default:
		return Byte
	}
}

// Size represents a size.
type Size struct {
	LowerBound ByteUnit
	UpperBound ByteUnit
}

// ParseSize parses a string representation of a size and returns a Size object.
//
// The size string should be in the format of a number followed by a unit (e.g., "10kb").
//
// The supported units are "b" for bytes, "kb" for kilobytes, "mb" for megabytes, and "gb" for gigabytes.
//
// The function returns an error if the size string is empty or if it cannot be parsed.
func ParseSize(size string) (sizeObj Size, err error) {
	if size == "" {
		return sizeObj, errors.New("size cannot be empty")
	}

	sizeExpr := gomme.Map(
		gomme.Pair(
			gomme.Int64[string](),
			gomme.Alternative(
				gomme.Assign(Byte, gomme.Token[string]("b")),
				gomme.Assign(Kilobyte, gomme.Token[string]("kb")),
				gomme.Assign(Megabyte, gomme.Token[string]("mb")),
				gomme.Assign(Gigabyte, gomme.Token[string]("gb")),
			),
		),
		func(p gomme.PairContainer[int64, ByteUnit]) (int64, error) {
			return p.Left * int64(p.Right), nil
		},
	)

	parser := gomme.SeparatedPair(
		sizeExpr,
		gomme.Optional(gomme.Token[string]("-")),
		gomme.Optional(sizeExpr),
	)

	bounds := parser(size)
	if bounds.Err != nil {
		return sizeObj, bounds.Err
	}

	if bounds.Output.Left != 0 {
		sizeObj.LowerBound = ByteUnit(bounds.Output.Left)
	}

	if bounds.Output.Right != 0 {
		sizeObj.UpperBound = ByteUnit(bounds.Output.Right)
	}

	return sizeObj, nil
}

// ErrNegativeBound is returned when a bound is negative.
var ErrNegativeBound = errors.New("bounds cannot be negative")

// ErrUpperBoundGreaterThanLowerBound is returned when the upper bound is greater than the lower bound.
var ErrUpperBoundGreaterThanLowerBound = errors.New("upper bound cannot be greater than lower bound")

// Validate checks if the Size struct satisfies the defined constraints.
// It returns an error if any of the constraints are violated.
func (s Size) Validate() error {
	if s.LowerBound < 0 {
		return fmt.Errorf("lower bound is negative: %w", ErrNegativeBound)
	}

	if s.UpperBound < 0 {
		return fmt.Errorf("upper bound is negative: %w", ErrNegativeBound)
	}

	if s.UpperBound > 0 && s.LowerBound > s.UpperBound {
		return ErrUpperBoundGreaterThanLowerBound
	}

	return nil
}

// String returns the string representation of the size.
func (s Size) String() string {
	if !s.HasBounds() {
		return fmt.Sprintf("%d", s.LowerBound)
	}

	return fmt.Sprintf("%d-%d", s.LowerBound, s.UpperBound)
}

// HasBounds returns true if the size has upper and lower bounds.
func (s Size) HasBounds() bool {
	// Neither bound is set
	if s.LowerBound == 0 && s.UpperBound == 0 {
		return false
	}

	if s.LowerBound != 0 && s.UpperBound != 0 {
		return true
	}

	return false
}

// Payload returns a byte slice containing a randomly generated payload.
//
// The size of the payload is determined by the LowerBound field of the Size struct.
//
// The payload is generated using the alphabet stored in the letterRunes array.
//
// If an upper bound is specified in the Size struct, the payload will be extended
// to the upper bound with additional random bytes.
//
//nolint:gosec
func (s Size) Payload() []byte {
	// We store the alphabet used to generate the payload
	// statically on the Stack.
	letterRunes := [...]byte{
		'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h',
		'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
		'q', 'r', 's', 't', 'u', 'v', 'w', 'x',
		'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F',
		'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N',
		'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V',
		'W', 'X', 'Y', 'Z',
	}

	// The byte slice to return will be at least
	// the size of the lower bound.
	bytes := make([]byte, s.LowerBound)

	// Fill the bytes slice with random letters
	for i := range bytes {
		bytes[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	// If we have an upper bound, we extend the bytes
	// slice to the upper bound and fill it with random
	// bytes on the way.
	if s.HasBounds() {
		difference := s.UpperBound - s.LowerBound
		addedSize := rand.Intn(int(difference))
		for i := 0; i < addedSize; i++ {
			//nolint:makezero
			bytes = append(bytes, letterRunes[rand.Intn(len(letterRunes))])
		}
	}

	return bytes
}
