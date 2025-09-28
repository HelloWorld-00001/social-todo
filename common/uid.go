package common

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
)

const secretByte byte = 57

type Uid struct {
	localId  uint32
	objectId int
	sharedId uint32
}

func (u Uid) LocalId() uint32 {
	return u.localId
}

func (u Uid) ObjectId() int {
	return u.objectId
}

func (u Uid) SharedId() uint32 {
	return u.sharedId
}

func NewUid(localId uint32, objectId int, sharedId uint32) *Uid {
	return &Uid{localId, objectId, sharedId}
}

func (u Uid) String() string {
	buf := new(bytes.Buffer)

	// Write fields in binary
	_ = binary.Write(buf, binary.LittleEndian, u.localId)
	_ = binary.Write(buf, binary.LittleEndian, int64(u.objectId)) // 8 bytes for int
	_ = binary.Write(buf, binary.LittleEndian, u.sharedId)

	data := buf.Bytes()

	// todo: design a more reliable algorithm
	for i := range data {
		offset := (secretByte + byte(i%255)) % 255
		data[i] ^= offset
	}

	// Base64 encode
	return base64.URLEncoding.EncodeToString(data)
}

// Optional: Reverse the process
func GetUidFromString(encoded string) (*Uid, error) {
	data, err := base64.URLEncoding.DecodeString(encoded)
	if err != nil {
		return nil, err
	}

	// Reverse XOR
	for i := range data {
		offset := (secretByte + byte(i%255)) % 255
		data[i] ^= offset
	}

	buf := bytes.NewReader(data)

	var localId uint32
	var objectId int64
	var sharedId uint32

	if err := binary.Read(buf, binary.LittleEndian, &localId); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &objectId); err != nil {
		return nil, err
	}
	if err := binary.Read(buf, binary.LittleEndian, &sharedId); err != nil {
		return nil, err
	}

	return &Uid{
		localId:  localId,
		objectId: int(objectId),
		sharedId: sharedId,
	}, nil
}

func (u Uid) MarshalJSON() ([]byte, error) {
	encoded := u.String()
	return []byte(`"` + encoded + `"`), nil
}
