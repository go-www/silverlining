package percent

import (
	"testing"
)

func TestDecodeHex(t *testing.T) {
	tests := []struct {
		name string
		d    byte
		want byte
	}{
		{
			name: "0",
			d:    '0',
			want: 0x0,
		},
		{
			name: "1",
			d:    '1',
			want: 0x1,
		},
		{
			name: "F",
			d:    'F',
			want: 0xF,
		},
		{
			name: "a",
			d:    'a',
			want: 0xa,
		},
		{
			name: "f",
			d:    'f',
			want: 0xf,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeHexOne(tt.d); got != tt.want {
				t.Errorf("DecodeHex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecodeHexTwo(t *testing.T) {
	tests := []struct {
		name string
		a, b byte
		want byte
	}{
		{
			name: "00",
			a:    '0',
			b:    '0',
			want: 0x00,
		},

		{
			name: "01",
			a:    '0',
			b:    '1',
			want: 0x01,
		},

		{
			name: "0F",
			a:    '0',
			b:    'F',
			want: 0x0F,
		},

		{
			name: "0a",
			a:    '0',
			b:    'a',
			want: 0x0a,
		},

		{
			name: "FF",
			a:    'F',
			b:    'F',
			want: 0xFF,
		},

		{
			name: "aF",
			a:    'a',
			b:    'F',
			want: 0xAF,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DecodeHexTwo(tt.a, tt.b); got != tt.want {
				t.Errorf("DecodeHexTwo() = %v, want %v", got, tt.want)
			}
		})
	}
}
