package responses

type UserResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaginatedResponse struct {
	Total int         `json:"total"`
	Limit int         `json:"limit"`
	Skip  int         `json:"skip"`
	Data  interface{} `json:"data"`
}
