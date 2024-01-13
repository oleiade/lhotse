package main

import (
	"reflect"
	"testing"
	"time"
)

func TestParseLatency(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		duration      string
		wantLatency   Latency
		wantErr       bool
		expectedError string
	}{
		{
			name:     "parsing a valid duration without upper bound should succeed",
			duration: "1s",
			wantLatency: Latency{
				LowerBound: 1 * time.Second,
				UpperBound: 0,
			},
			wantErr: false,
		},
		{
			name:     "parsing a valid duration with upper bound should succeed",
			duration: "1s-10s",
			wantLatency: Latency{
				LowerBound: 1 * time.Second,
				UpperBound: 10 * time.Second,
			},
			wantErr: false,
		},
		{
			name:          "parsing an invalid duration should return an error",
			duration:      "invalid",
			wantLatency:   Latency{},
			wantErr:       true,
			expectedError: "failed parsing duration's lower bound",
		},
		{
			name:          "parsing a duration with upper bound less than lower bound should return an error",
			duration:      "10s-1s",
			wantLatency:   Latency{},
			wantErr:       true,
			expectedError: "upper bound is less than lower bound",
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotLatency, err := ParseLatency(tt.duration)

			if (err != nil) != tt.wantErr {
				t.Errorf("ParseLatency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil && err.Error() != tt.expectedError {
				t.Errorf("ParseLatency() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(gotLatency, tt.wantLatency) {
				t.Errorf("ParseLatency() = %v, want %v", gotLatency, tt.wantLatency)
			}
		})
	}
}
