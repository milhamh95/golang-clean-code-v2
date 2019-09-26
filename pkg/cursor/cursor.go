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

// Encode is a function to transform cursor object to string
func Encode(cursor Cursor) (encodedCursor string, err error) {
	cursorByte, err := json.Marshal(cursor)
	if err != nil {
		return
	}
	encodedCursor = base64.StdEncoding.EncodeToString(cursorByte)
	return
}

// Decode is a function to transform string to cursor object
func Decode(encodedCursor string) (decodedcursor Cursor, err error) {
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

// EncodeBase64 encode string to base64 form
func EncodeBase64(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}

// DecodeBase64 decode base64 value to string
func DecodeBase64(value string) (res string, err error) {
	cursorByte, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return
	}

	res = string(cursorByte)
	return
}
