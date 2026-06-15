package api

import (
	"license/internal/gitlab/types"
	"license/internal/gitlab/service"

	"github.com/gin-gonic/gin"
)

// Controller defines the controller structure
type Controller struct {
}

// NewController creates a new controller instance
func NewController() *Controller {
	return &Controller{}
}

// Generate handles license generation
func (controller *Controller) Generate(ctx *gin.Context) {

	Name := ctx.PostForm("Name")
	Email := ctx.PostForm("Email")
	Company := ctx.PostForm("Company")
	var license = types.LicenseInfo{
		Name:    Name,
		Email:   Email,
		Company: Company,
	}
	// Expiration time
	expireTime := ctx.PostForm("ExpireTime")
	// Generate license
	service.Generate(ctx, license, expireTime)
}
