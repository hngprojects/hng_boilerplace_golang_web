package product

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/hngprojects/hng_boilerplate_golang_web/external/request"
	"github.com/hngprojects/hng_boilerplate_golang_web/internal/models"
	"github.com/hngprojects/hng_boilerplate_golang_web/pkg/repository/storage"
	"github.com/hngprojects/hng_boilerplate_golang_web/services/product"
	"github.com/hngprojects/hng_boilerplate_golang_web/utility"
)

type Controller struct {
	Db        *storage.Database
	Validator *validator.Validate
	Logger    *utility.Logger
	ExtReq    request.ExternalRequest
}

func (base *Controller) CreateProduct(c *gin.Context) {

	var (
		req = models.CreateProductRequestModel{}
	)

	err := c.ShouldBind(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", "Failed to parse request body", err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	err = base.Validator.Struct(&req)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusUnprocessableEntity, "error", "Validation failed", utility.ValidationResponse(err, base.Validator), nil)
		c.JSON(http.StatusUnprocessableEntity, rd)
		return
	}

	respData, code, err := product.CreateProduct(req, base.Db.Postgresql, c)
	if err != nil {
		rd := utility.BuildErrorResponse(http.StatusBadRequest, "error", err.Error(), err, nil)
		c.JSON(http.StatusBadRequest, rd)
		return
	}

	base.Logger.Info("Product created successfully")
	rd := utility.BuildSuccessResponse(http.StatusCreated, "Product created successfully", respData)

	c.JSON(code, rd)
}

func (base *Controller) GetProduct(c *gin.Context) {
	productId := c.Param("product_id")
	respData, code, err := product.GetProduct(productId, base.Db.Postgresql)
	if err != nil {
		resp := gin.H{"error": "Product not found"}
		if code == http.StatusNotFound {
			resp = gin.H{"error": "Invalid product ID"}
		}

		rd := utility.BuildErrorResponse(code, "error", err.Error(), resp, nil)
		c.JSON(code, rd)
		return
	}

	base.Logger.Info("Product found successfully")
	rd := utility.BuildSuccessResponse(http.StatusOK, "Product found successfully", respData)

	c.JSON(code, rd)
}
