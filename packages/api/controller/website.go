package controller

import (
	"github.com/dupmanio/dupman/packages/api/model"
	"github.com/dupmanio/dupman/packages/api/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WebsiteController struct {
	websiteSvc *service.WebsiteService
}

func NewWebsiteController(websiteSvc *service.WebsiteService) (*WebsiteController, error) {
	return &WebsiteController{websiteSvc: websiteSvc}, nil
}

func (ctrl *WebsiteController) Create(ctx *gin.Context) {
	// @todo: needs implementation.
	ctrl.websiteSvc.Create(&model.Website{
		URL:    "https://example.com",
		UserID: uuid.New(),
	})
}
