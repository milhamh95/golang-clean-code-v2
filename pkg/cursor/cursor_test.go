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

	cursor, err := cursor.Encode(argCursor)
	require.NoError(t, err)
	require.Equal(t, cursor, wantCursor)
}

func TestDecodeCursor(t *testing.T) {
	argCursor := "eyJpdGVtX2N1cnNvciI6Ik5BPT0iLCJsYXN0X3Bvc2l0aW9uIjoxMH0="

	wantCursor := cursor.Cursor{
		LastPosition: 10,
		ItemCursor:   "NA==",
	}

	cursor, err := cursor.Decode(argCursor)
	require.NoError(t, err)
	require.Equal(t, cursor, wantCursor)
}

func TestEncodeBase64(t *testing.T) {
	want := "MHVqc3N3VGhJR1RVWW0ySzhGak9PZlh0WTFL"

	get := cursor.EncodeBase64("0ujsswThIGTUYm2K8FjOOfXtY1K")
	require.Equal(t, want, get)
}

func TestDecodeBase64(t *testing.T) {
	want := "0ujsswThIGTUYm2K8FjOOfXtY1K"

	get, err := cursor.DecodeBase64("MHVqc3N3VGhJR1RVWW0ySzhGak9PZlh0WTFL")
	require.Equal(t, want, get)
	require.NoError(t, err)
}
