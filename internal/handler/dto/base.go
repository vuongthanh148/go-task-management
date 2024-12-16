package dto

type BaseResponse struct {
	Message string `json:"message"`
}

type Paginate struct {
	Page  uint32 `json:"page" form:"page"`
	Limit uint32 `json:"limit" form:"limit"`
}
