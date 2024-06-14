package languages

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsSupported(t *testing.T) {
	tests := []struct {
		name      string
		id        int
		supported bool
	}{
		{
			name:      "Go",
			id:        132,
			supported: true,
		},
		{
			name:      "Python",
			id:        303,
			supported: true,
		},
		// Esoteric languages are not meant to be code reviewed, this
		// test should be stable across user requests for more languages.
		{
			name:      "Befunge",
			id:        30,
			supported: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, tc.supported, IsSupported(tc.id))
		})
	}
}
