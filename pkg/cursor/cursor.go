package cursor

import (
	"encoding/base64"
	"encoding/json"
)

// Cursor represent cursor model
type Cursor struct {
	ItemCursor   string `json:"item_cursor"`
	LastPosition int    `json:"last_position"`
}

// EncodeCursor is a function to transform cursor object to string
func EncodeCursor(cursor Cursor) (encodedCursor string, err error) {
	cursorByte, err := json.Marshal(cursor)
	if err != nil {
		return
	}
	encodedCursor = base64.StdEncoding.EncodeToString(cursorByte)
	return
}

// DecodeCursor is a function to transform string to cursor object
func DecodeCursor(encodedCursor string) (decodedcursor Cursor, err error) {
	cursorByte, err := base64.StdEncoding.DecodeString(encodedCursor)
	if err != nil {
		return
	}

	err = json.Unmarshal(cursorByte, &decodedcursor)
	if err != nil {
		return
	}

	return
}
