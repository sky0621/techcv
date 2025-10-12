package uuidv7

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

// NewString generates a UUID version 7 as a canonical string.
func NewString() (string, error) {
	var uuid [16]byte

	ts := uint64(time.Now().UTC().UnixMilli())
	uuid[0] = byte(ts >> 40)
	uuid[1] = byte(ts >> 32)
	uuid[2] = byte(ts >> 24)
	uuid[3] = byte(ts >> 16)
	uuid[4] = byte(ts >> 8)
	uuid[5] = byte(ts)

	if _, err := rand.Read(uuid[6:]); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	uuid[6] = (uuid[6] & 0x0f) | 0x70
	uuid[8] = (uuid[8] & 0x3f) | 0x80

	buf := make([]byte, 32)
	hex.Encode(buf, uuid[:])

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		buf[0:8],
		buf[8:12],
		buf[12:16],
		buf[16:20],
		buf[20:32],
	), nil
}
