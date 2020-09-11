package utils

import (
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

// GenUUIDv4 生成一个 UUIDv4
func GenUUIDv4() string {
	u2, err := uuid.NewV4()
	if err != nil {
		log.Errorf("failed to generate UUID: %v", err.Error())
	}
	return u2.String()
}

// GenUUIDv5 生成一个 UUIDv5
func GenUUIDv5(name string) string {
	u2, err := uuid.NewV4()
	if err != nil {
		log.Errorf("failed to generate UUID: %v", err.Error())
	}
	return uuid.NewV5(u2, "test").String()
}

// ParseUUID 解析 UUID
func ParseUUID(s string) uuid.UUID {
	u3, err := uuid.FromString(s)
	if err != nil {
		log.Errorf("failed to parse UUID %q: %v", s, err)
	}
	return u3
}
