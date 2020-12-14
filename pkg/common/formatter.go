package common

import "restaurant-assistant/pkg/entity"

func FormatResponse(r entity.Response) interface{} {
	return r.Format()
}
