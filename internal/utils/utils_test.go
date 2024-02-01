package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetMD5Hash(t *testing.T) {
	tests := []struct {
		name string
		text string
		want string
	}{
		{
			name: "empty string",
			text: "",
			want: "d41d8cd98f00b204e9800998ecf8427e",
		},
		{
			name: "alphanumeric string",
			text: "abc123",
			want: "e99a18c428cb38d5f260853678922e03",
		},
		{
			name: "special characters",
			text: "#$!",
			want: "712415625c07d00e6e93aa884801b84e",
		},
		{
			name: "long string",
			text: "ThisIsALongStringWithSomeDifferentCharacters!!",
			want: "dc503be185eb8856b99d4c9bae6e91bc",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, GetMD5Hash(tt.text))
		})
	}
}
