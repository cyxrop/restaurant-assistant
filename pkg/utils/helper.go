package utils

import (
	uuid "github.com/satori/go.uuid"

	"restaurant-assistant/pkg/entity"
)

func IsEmptyUUID(uuid uuid.UUID) bool {
	return uuid.String() == "00000000-0000-0000-0000-000000000000"
}

func FormatResponse(r entity.Response) interface{} {
	return r.Format()
}
