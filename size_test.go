package main

import (
	"reflect"
	"testing"
)

func TestParseSize(t *testing.T) {
	t.Parallel()

	type args struct {
		size string
	}
	tests := []struct {
		name    string
		args    args
		wantS   Size
		wantErr bool
	}{
		{
			name: "parsing the 1b expression to 1 byte should succeed",
			args: args{size: "1b"},
			wantS: Size{
				LowerBound: 1 * Byte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 1kb expression to 1 kilobyte should succeed",
			args: args{size: "1kb"},
			wantS: Size{
				LowerBound: 1 * Kilobyte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 1mb expression to 1 megabyte should succeed",
			args: args{size: "1mb"},
			wantS: Size{
				LowerBound: 1 * Megabyte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 1gb expression to 1 gigabyte should succeed",
			args: args{size: "1gb"},
			wantS: Size{
				LowerBound: 1 * Gigabyte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 123b expression to 123 bytes should succeed",
			args: args{size: "123b"},
			wantS: Size{
				LowerBound: 123 * Byte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 123kb expression to 123 kilobytes should succeed",
			args: args{size: "123kb"},
			wantS: Size{
				LowerBound: 123 * Kilobyte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 123mb expression to 123 megabytes should succeed",
			args: args{size: "123mb"},
			wantS: Size{
				LowerBound: 123 * Megabyte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 123gb expression to 123 gigabytes should succeed",
			args: args{size: "123gb"},
			wantS: Size{
				LowerBound: 123 * Gigabyte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 1b-123b expression to a size of 1 to 10 bytes should succeed",
			args: args{size: "1b-10b"},
			wantS: Size{
				LowerBound: 1 * Byte,
				UpperBound: 10 * Byte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 1kb-123kb expression to a size of 1 to 123 kilobytes should succeed",
			args: args{size: "1kb-123kb"},
			wantS: Size{
				LowerBound: 1 * Kilobyte,
				UpperBound: 123 * Kilobyte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 1mb-123mb expression to a size of 1 to 123 megabytes should succeed",
			args: args{size: "1mb-123mb"},
			wantS: Size{
				LowerBound: 1 * Megabyte,
				UpperBound: 123 * Megabyte,
			},
			wantErr: false,
		},
		{
			name: "parsing the 1gb expression to a size of 1 to 123 gigabyte should succeed",
			args: args{size: "1gb-123gb"},
			wantS: Size{
				LowerBound: 1 * Gigabyte,
				UpperBound: 123 * Gigabyte,
			},
			wantErr: false,
		},
		{
			name: "parsing a size that spans across ranges should succeed",
			args: args{size: "1b-1gb"},
			wantS: Size{
				LowerBound: 1 * Byte,
				UpperBound: 1 * Gigabyte,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotS, err := ParseSize(tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseSize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("ParseSize() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}

func TestSize_Validate(t *testing.T) {
	t.Parallel()

	type fields struct {
		LowerBound ByteUnit
		UpperBound ByteUnit
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "validating a size with a lower bound of 1 byte and an upper bound of 0 byte should succeed",
			fields: fields{
				LowerBound: 1 * Byte,
				UpperBound: 0 * Byte,
			},
			wantErr: false,
		},
		{
			name: "validating a size with a lower bound of 1 byte and an upper bound of 1 byte should succeed",
			fields: fields{
				LowerBound: 1 * Byte,
				UpperBound: 1 * Byte,
			},
			wantErr: false,
		},
		{
			name: "validating a size with a lower bound of 1 byte and an upper bound of 2 bytes should succeed",
			fields: fields{
				LowerBound: 1 * Byte,
				UpperBound: 2 * Byte,
			},
			wantErr: false,
		},
		{
			name: "validating a size with a lower bound of 1 byte and an upper bound of 1 kilobyte should succeed",
			fields: fields{
				LowerBound: 1 * Byte,
				UpperBound: 1 * Kilobyte,
			},
			wantErr: false,
		},
		{
			name: "validating a size with a lower bound with a bigger unit but greater absolution value than upper bound should succeed",
			fields: fields{
				LowerBound: 1 * Megabyte,
				UpperBound: 100000 * Kilobyte,
			},
			wantErr: false,
		},
		{
			name: "validating a size with a lower bound < 0 should fail",
			fields: fields{
				LowerBound: -1 * Byte,
				UpperBound: 0 * Byte,
			},
			wantErr: true,
		},
		{
			name: "validating a size with a lower bound > upper bound should fail",
			fields: fields{
				LowerBound: 2 * Byte,
				UpperBound: 1 * Byte,
			},
			wantErr: true,
		},
		{
			name: "validating a size with an upper bound < 0 should fail",
			fields: fields{
				LowerBound: 0 * Byte,
				UpperBound: -1 * Byte,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := Size{
				LowerBound: tt.fields.LowerBound,
				UpperBound: tt.fields.UpperBound,
			}
			if err := s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Size.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSize_Payload(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		size Size
	}{
		{
			name: "payload with lower bound of 1 byte",
			size: Size{
				LowerBound: 1 * Byte,
				UpperBound: 0 * Byte,
			},
		},
		{
			name: "payload with lower bound of 1 kilobyte",
			size: Size{
				LowerBound: 1 * Kilobyte,
				UpperBound: 0 * Byte,
			},
		},
		{
			name: "payload with lower bound of 1 megabyte",
			size: Size{
				LowerBound: 1 * Megabyte,
				UpperBound: 0 * Byte,
			},
		},
		{
			name: "payload with lower bound of 1 gigabyte",
			size: Size{
				LowerBound: 1 * Gigabyte,
				UpperBound: 0 * Byte,
			},
		},
		{
			name: "payload with lower bound of 123 bytes",
			size: Size{
				LowerBound: 123 * Byte,
				UpperBound: 0 * Byte,
			},
		},
		{
			name: "payload with lower bound of 123 kilobytes",
			size: Size{
				LowerBound: 123 * Kilobyte,
				UpperBound: 0 * Byte,
			},
		},
		{
			name: "payload with lower bound of 123 megabytes",
			size: Size{
				LowerBound: 123 * Megabyte,
				UpperBound: 0 * Byte,
			},
		},
		{
			name: "payload with lower bound of 123 gigabytes",
			size: Size{
				LowerBound: 123 * Gigabyte,
				UpperBound: 0 * Byte,
			},
		},
		{
			name: "payload with lower bound of 1 byte and upper bound of 10 bytes",
			size: Size{
				LowerBound: 1 * Byte,
				UpperBound: 10 * Byte,
			},
		},
		{
			name: "payload with lower bound of 1 kilobyte and upper bound of 123 kilobytes",
			size: Size{
				LowerBound: 1 * Kilobyte,
				UpperBound: 123 * Kilobyte,
			},
		},
		{
			name: "payload with lower bound of 1 megabyte and upper bound of 123 megabytes",
			size: Size{
				LowerBound: 1 * Megabyte,
				UpperBound: 123 * Megabyte,
			},
		},
		{
			name: "payload with lower bound of 1 gigabyte and upper bound of 123 gigabytes",
			size: Size{
				LowerBound: 1 * Gigabyte,
				UpperBound: 123 * Gigabyte,
			},
		},
		{
			name: "payload with lower bound of 1 byte and upper bound of 1 gigabyte",
			size: Size{
				LowerBound: 1 * Byte,
				UpperBound: 1 * Gigabyte,
			},
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			payload := tt.size.Payload()
			payloadSize := len(payload)

			if payloadSize < int(tt.size.LowerBound) || (tt.size.HasBounds() && payloadSize > int(tt.size.UpperBound)) {
				t.Errorf("Payload size mismatch for %s: got %d bytes, expected between %d and %d bytes", tt.name, payloadSize, tt.size.LowerBound, tt.size.UpperBound)
			}
		})
	}
}

func TestParseByteUnit(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		u    string
		want ByteUnit
	}{
		{
			name: "parsing 'b' should return Byte",
			u:    "b",
			want: Byte,
		},
		{
			name: "parsing 'kb' should return Kilobyte",
			u:    "kb",
			want: Kilobyte,
		},
		{
			name: "parsing 'mb' should return Megabyte",
			u:    "mb",
			want: Megabyte,
		},
		{
			name: "parsing 'gb' should return Gigabyte",
			u:    "gb",
			want: Gigabyte,
		},
		{
			name: "parsing unknown unit should return Byte",
			u:    "unknown",
			want: Byte,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := ParseByteUnit(tt.u)
			if got != tt.want {
				t.Errorf("ParseByteUnit() = %v, want %v", got, tt.want)
			}
		})
	}
}
