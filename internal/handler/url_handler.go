package handler

import (
	"net/http"

	"github.com/Somaycon/gourl.git/internal/model"
	"github.com/Somaycon/gourl.git/internal/service"
	"github.com/gin-gonic/gin"
)

type UrlHandler struct {
	service *service.UrlService
}

func NewUrlHandler(service *service.UrlService) *UrlHandler {
	return &UrlHandler{
		service: service,
	}
}

func (h *UrlHandler) Create(ctx *gin.Context) {
	var request model.UrlRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL or not allowed",
		})
		return
	}
	response, err := h.service.Shorten(ctx.Request.Context(), request.Url)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *UrlHandler) Redirect(ctx *gin.Context) {
	short := ctx.Param("short")
	originalUrl, err := h.service.GetOrininalUrl(ctx.Request.Context(), short)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "URL not found",
		})
		return
	}
	ctx.Redirect(http.StatusFound, originalUrl)
}
