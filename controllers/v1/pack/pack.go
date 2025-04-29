package pack

import (
	"encoding/json"
	"github.com/dekib/rePacks/internal/errors"
	"github.com/dekib/rePacks/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

type PackController struct {
	packService *services.PackService
}

func NewPackController(packService *services.PackService) *PackController {
	return &PackController{
		packService: packService,
	}
}

func (c *PackController) GetPacks(ctx *gin.Context) {
	var req CalculatePacksRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		var errorResp *errors.ErrorResponse

		switch err.(type) {
		case validator.ValidationErrors:
			errorResp = errors.NewBadRequestError(
				errors.ErrInvalidPackSizes,
				err,
				errors.V{"min_size": 1},
			)
		default:
			errorResp = errors.NewBadRequestError(
				errors.ErrMalformedRequest,
				err,
				nil,
			)
		}

		errorResp.Abort(ctx)
		return
	}

	if err := req.Validate(); err != nil {
		errors.NewBadRequestError(
			errors.ErrValidation,
			err,
			nil,
		).Abort(ctx)
		return
	}

	packs := c.packService.CalculatePacks(req.ItemsNo)
	response := CalculatePacksResponse{
		Packs: make([]PackItem, 0, len(packs)),
	}

	for _, p := range packs {
		response.Packs = append(response.Packs, PackItem{
			Pack:     p.Pack,
			Quantity: p.Quantity,
		})
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *PackController) UpdatePackSizes(ctx *gin.Context) {
	var req UpdatePackSizesRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		var errorResp *errors.ErrorResponse

		switch err.(type) {
		case validator.ValidationErrors:
			errorResp = errors.NewBadRequestError(
				errors.ErrInvalidPackSizes,
				err,
				errors.V{"min_size": 1},
			)
		case *json.UnmarshalTypeError, *json.SyntaxError:
			errorResp = errors.NewBadRequestError(
				errors.ErrMalformedRequest,
				err,
				nil,
			)
		default:
			errorResp = errors.NewBadRequestError(
				errors.ErrBadRequest,
				err,
				nil,
			)
		}

		errorResp.Abort(ctx)
		return
	}

	if err := req.Validate(); err != nil {
		errors.NewBadRequestError(
			errors.ErrValidation,
			err,
			nil,
		).Abort(ctx)
		return
	}

	c.packService.SetPackSizes(req.Sizes)

	ctx.JSON(http.StatusOK, UpdatePackSizesResponse{
		Message:   "Pack sizes updated successfully",
		PackSizes: c.packService.GetPackSizes(),
	})
}
