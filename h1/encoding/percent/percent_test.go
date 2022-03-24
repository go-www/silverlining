package percent

import (
	"testing"
)

func TestDecode(t *testing.T) {
	tests := []struct {
		name   string
		buffer string
		want   string
	}{
		{
			name:   "empty",
			buffer: "",
			want:   "",
		},
		{
			name:   "no percent",
			buffer: "abc",
			want:   "abc",
		},
		{
			name:   "SELECT%20*%20FROM%20table",
			buffer: "SELECT%20*%20FROM%20table",
			want:   "SELECT * FROM table",
		},
		{
			name:   "!#$%&'()*+,/:;=?@[]",
			buffer: "%21%23%24%25%26%27()*%2B%2C%2F%3A%3B%3D%3F%40%5B%5D",
			want:   "!#$%&'()*+,/:;=?@[]",
		},
		{
			name:   "🎉✨🎇🎁",
			buffer: "%F0%9F%8E%89%E2%9C%A8%F0%9F%8E%87%F0%9F%8E%81",
			want:   "🎉✨🎇🎁",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Decode([]byte(tt.buffer)); string(got) != tt.want {
				t.Errorf("Decode() = %v, want %v", string(got), tt.want)
			}
		})
	}
}
