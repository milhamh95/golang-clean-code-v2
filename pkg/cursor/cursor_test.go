package cursor_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/milhamhidayat/golang-clean-code-v2/pkg/cursor"
)

func TestEncodeCursor(t *testing.T) {
	argCursor := cursor.Cursor{
		LastPosition: 10,
		ItemCursor:   "NA==",
	}

	wantCursor := "eyJpdGVtX2N1cnNvciI6Ik5BPT0iLCJsYXN0X3Bvc2l0aW9uIjoxMH0="

	cursor, err := cursor.EncodeCursor(argCursor)
	require.NoError(t, err)
	require.Equal(t, cursor, wantCursor)
}

func TestDecodeCursor(t *testing.T) {
	argCursor := "eyJpdGVtX2N1cnNvciI6Ik5BPT0iLCJsYXN0X3Bvc2l0aW9uIjoxMH0="

	wantCursor := cursor.Cursor{
		LastPosition: 10,
		ItemCursor:   "NA==",
	}

	cursor, err := cursor.DecodeCursor(argCursor)
	require.NoError(t, err)
	require.Equal(t, cursor, wantCursor)
}
