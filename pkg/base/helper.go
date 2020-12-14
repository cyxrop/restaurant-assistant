package base

import uuid "github.com/satori/go.uuid"

func IsEmptyUUID(uuid uuid.UUID) bool {
	return uuid.String() == "00000000-0000-0000-0000-000000000000"
}
