package main

import "testing"

func TestDisplayVersion(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			name: "release",
			in:   "v0.1.0",
			want: "v0.1.0",
		},
		{
			name: "homebrew head",
			in:   "HEAD-ad1b265",
			want: "dev (ad1b265)",
		},
		{
			name: "empty",
			in:   " ",
			want: "dev",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := displayVersion(tt.in); got != tt.want {
				t.Fatalf("displayVersion(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}
