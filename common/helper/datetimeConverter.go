package helper

import (
	"encoding/base64"
	"time"
)

// Custom layout (template) for formatting/parsing time
const timeLayout = "2006-01-02T15:04:05.999999"      // ISO-like for encoding/decoding
const MySqlTimeLayout = "2006-01-02 15:04:05.999999" // MySQL-compatible

// EncodeTimeToBase64URL converts time to string using layout and encodes it to Base64 URL-safe
func EncodeTimeToBase64URL(t time.Time) string {
	timeStr := t.Format(timeLayout)
	return base64.URLEncoding.EncodeToString([]byte(timeStr))
}

// DecodeBase64URLToTime decodes Base64 URL-safe string back to time.Time using layout
func DecodeBase64URLToTime(encoded string) (time.Time, error) {
	decodedBytes, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return time.Time{}, err
	}
	return time.Parse(timeLayout, string(decodedBytes))
}
