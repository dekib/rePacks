package pack

import (
	"github.com/go-playground/validator/v10"
)

// Request DTOs
type (
	CalculatePacksRequest struct {
		ItemsNo int `form:"items_no" validate:"required,min=1"`
	}

	UpdatePackSizesRequest struct {
		Sizes []int `json:"sizes" validate:"required,min=1,dive,gt=0"`
	}
)

// Response DTOs
type (
	PackItem struct {
		Pack     int `json:"pack"`
		Quantity int `json:"quantity"`
	}

	CalculatePacksResponse struct {
		Packs []PackItem `json:"packs"`
	}

	UpdatePackSizesResponse struct {
		Message   string `json:"message"`
		PackSizes []int  `json:"pack_sizes"`
	}
)

var validate = validator.New()

func (r *CalculatePacksRequest) Validate() error {
	return validate.Struct(r)
}

func (r *UpdatePackSizesRequest) Validate() error {
	return validate.Struct(r)
}
