package main

import (
	"reflect"
	"testing"
)

func TestParseSize(t *testing.T) {
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
			name: "parsing a size that spans accross ranges should suceed",
			args: args{size: "1b-1gb"},
			wantS: Size{
				LowerBound: 1 * Byte,
				UpperBound: 1 * Gigabyte,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
		t.Run(tt.name, func(t *testing.T) {
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
