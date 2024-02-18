package utils

import "github.com/google/uuid"

func ParseUUID(stringUUID string) (uuid.UUID, error) {
	return uuid.Parse(stringUUID)
}
