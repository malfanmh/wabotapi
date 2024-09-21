package model

import (
	"testing"
)

func TestFormatExpiryDate(t *testing.T) {
	type args struct {
		inputTime string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"valid",
			args{inputTime: "2024-09-19 14:23:32"},
			"19 September 2024, pukul 14:23 WIB",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatExpiryDate(tt.args.inputTime); got != tt.want {
				t.Errorf("FormatExpiryDate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatRP(t *testing.T) {
	type args struct {
		amount float64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"1000",
			args{amount: 1000},
			"Rp1.000,00",
		},
		{
			"1000,10",
			args{amount: 1000.10},
			"Rp1.000,10",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FormatRP(tt.args.amount); got != tt.want {
				t.Errorf("FormatRP() = %v, want %v", got, tt.want)
			}
		})
	}
}
