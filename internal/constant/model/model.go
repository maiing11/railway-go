package model

type WebResponse[T any] struct {
	Data   T             `json:"data"`
	Paging *PageMetaData `json:"paging,omitempty"`
	Errors string        `json:"errors,omitempty"`
}

type PageMetaData struct {
	Page      int   `json:"page"`
	Size      int   `json:"size"`
	TotalItem int64 `json:"total_item"`
	TotalPage int64 `json:"total_page"`
}

type PageResponse[T any] struct {
	Data         []T          `json:"data,omitempty"`
	PageMetaData PageMetaData `json:"paging,omitempty"`
}

func BuildSuccessResponse[T any](data T, paging *PageMetaData) WebResponse[T] {
	return WebResponse[T]{
		Data:   data,
		Paging: paging,
		Errors: "",
	}
}

func BuildErrorResponse(message string) WebResponse[any] {
	return WebResponse[any]{
		Data:   nil,
		Paging: nil,
		Errors: message,
	}
}
