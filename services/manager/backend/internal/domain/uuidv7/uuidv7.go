// Package uuidv7 provides utilities for generating UUID version 7 identifiers.
package uuidv7

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"
)

const (
	uuidSize                 = 16
	timestampByte0Shift      = 40
	timestampByte1Shift      = 32
	timestampByte2Shift      = 24
	timestampByte3Shift      = 16
	timestampByte4Shift      = 8
	versionBitsMask          = 0x0f
	versionByteValue         = 0x70
	variantBitsMask          = 0x3f
	variantByteValue         = 0x80
	hexBufferLength          = 32
	uuidSection1End          = 8
	uuidSection2End          = 12
	uuidSection3End          = 16
	uuidSection4End          = 20
	unixMilliNegativeMessage = "negative unix milli timestamp"
)

// NewString generates a UUID version 7 as a canonical string.
func NewString() (string, error) {
	var uuid [uuidSize]byte

	tsMillis := time.Now().UTC().UnixMilli()
	if tsMillis < 0 {
		return "", errors.New(unixMilliNegativeMessage)
	}
	timestamp := uint64(tsMillis)
	uuid[0] = byte(timestamp >> timestampByte0Shift)
	uuid[1] = byte(timestamp >> timestampByte1Shift)
	uuid[2] = byte(timestamp >> timestampByte2Shift)
	uuid[3] = byte(timestamp >> timestampByte3Shift)
	uuid[4] = byte(timestamp >> timestampByte4Shift)
	uuid[5] = byte(timestamp)

	if _, err := rand.Read(uuid[6:]); err != nil {
		return "", fmt.Errorf("generate random bytes: %w", err)
	}

	uuid[6] = (uuid[6] & versionBitsMask) | versionByteValue
	uuid[8] = (uuid[8] & variantBitsMask) | variantByteValue

	buf := make([]byte, hexBufferLength)
	hex.Encode(buf, uuid[:])

	return fmt.Sprintf("%s-%s-%s-%s-%s",
		buf[0:uuidSection1End],
		buf[uuidSection1End:uuidSection2End],
		buf[uuidSection2End:uuidSection3End],
		buf[uuidSection3End:uuidSection4End],
		buf[uuidSection4End:],
	), nil
}
