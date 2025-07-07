package common

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type Image struct {
	Url       string `json:"url"`
	Width     uint   `json:"width,omitempty"`
	Height    uint   `json:"height,omitempty"`
	FileSize  string `json:"file_size,omitempty"`
	FileName  string `json:"file_name,omitempty"`
	IsDelete  bool   `json:"is_delete,omitempty"`
	CloudName string `json:"cloud_name,omitempty"`
}

// Value converts Image to JSON for SQL storage
func (img *Image) Value() (driver.Value, error) {
	if img == nil {
		return nil, nil
	}
	return json.Marshal(img)
}

// Scan reads JSON from SQL and fills Image
func (img *Image) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var data []byte

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return NewBadRequestErrorResponse(errors.New("invalid type for Image Scan"), "cannot process invalid image", "")
	}

	return json.Unmarshal(data, img)
}

// Images represents a slice of Image
type Images []Image

// Value converts Images to JSON
func (imgs *Images) Value() (driver.Value, error) {
	if imgs == nil {
		return nil, nil
	}
	return json.Marshal(imgs)
}

// Scan reads JSON array from SQL into Images
func (imgs *Images) Scan(value interface{}) error {
	if value == nil {
		*imgs = nil
		return nil
	}
	var data []byte

	switch v := value.(type) {
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return NewBadRequestErrorResponse(errors.New("invalid type for Image Scan"), "cannot process invalid image", "")
	}

	return json.Unmarshal(data, imgs)
}
