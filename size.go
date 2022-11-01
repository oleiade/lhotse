package main

import (
	"fmt"
	"math/rand"

	"github.com/oleiade/gomme"
)

// FIXME: could be uint to ensure type safety
type ByteUnit int64

const (
	Byte     ByteUnit = 1
	Kilobyte          = 1024 * Byte
	Megabyte          = 1024 * Kilobyte
	Gigabyte          = 1024 * Megabyte
)

// ParseByteUnit parses a byte unit string and returns a ByteUnit.
func ParseByteUnit(u string) ByteUnit {
	switch u {
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

type Size struct {
	LowerBound ByteUnit
	UpperBound ByteUnit
}

// FIXME: bounds could be uint to ensure type safety
func ParseSize(size string) (s Size, err error) {
	alt := gomme.Map(
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
		alt,
		gomme.Optional(gomme.Token[string]("-")),
		gomme.Optional(alt),
	)

	bounds := parser(size)
	if bounds.Err != nil {
		return s, bounds.Err
	}

	if bounds.Output.Left != 0 {
		s.LowerBound = ByteUnit(bounds.Output.Left)
	}

	if bounds.Output.Right != 0 {
		s.UpperBound = ByteUnit(bounds.Output.Right)
	}

	return s, nil
}

var ErrNegativeBound = fmt.Errorf("bounds cannot be negative")
var ErrUpperBoundGreaterThanLowerBound = fmt.Errorf("upper bound cannot be greater than lower bound")

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

// TODO: we should support setting the kind of content we want: ascii, utf8, binary, etc...
func (s Size) Payload() []byte {
	// We store the alphabet used to generate the payload
	// stactically on the Stack.
	var letterRunes = [...]byte{
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
			bytes = append(bytes, letterRunes[rand.Intn(len(letterRunes))])
		}
	}

	return bytes
}
