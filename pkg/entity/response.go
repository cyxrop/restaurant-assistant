package entity

type Response interface {
	Format() interface{}
}

type BaseResponse struct {
	Message string `json:"message"`
}

func NewBaseResponse(message string) *BaseResponse {
	return &BaseResponse{
		Message: message,
	}
}

func (r *BaseResponse) Format() interface{} {
	return &BaseResponse{
		Message: r.Message,
	}
}
