package utils

import "github.com/google/uuid"

func StringToUUID(stringUUID string) (uuid.UUID, error) {
	return uuid.Parse(stringUUID)
}

func UUIDToString(uuid uuid.UUID) string {
	return uuid.String()
}
