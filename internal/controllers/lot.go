package controllers

import (
	"durich-be/internal/dto/requests"
	"durich-be/internal/services"
	"durich-be/pkg/authentication"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type LotController struct {
	lotService services.LotService
}

func NewLotController(lotService services.LotService) *LotController {
	return &LotController{
		lotService: lotService,
	}
}

func (c *LotController) Create(ctx *gin.Context) {
	var req requests.LotCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	result, err := c.lotService.Create(ctx.Request.Context(), req, locationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Lot berhasil dibuat",
		"data":    result,
	})
}

func (c *LotController) GetList(ctx *gin.Context) {
	status := ctx.Query("status")
	jenisDurianID := ctx.Query("jenis_durian_id")
	kondisi := ctx.Query("kondisi")

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	// DEBUG LOG
	fmt.Printf("[DEBUG] LotController.GetList - UserID: %s, LocationID: '%s'\n", userAuth.UserID, locationID)

	result, err := c.lotService.GetList(ctx.Request.Context(), status, jenisDurianID, kondisi, locationID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

func (c *LotController) GetDetail(ctx *gin.Context) {
	id := ctx.Param("id")

	result, err := c.lotService.GetDetail(ctx.Request.Context(), id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"status":  "error",
			"message": "Lot tidak ditemukan",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   result,
	})
}

func (c *LotController) AddItems(ctx *gin.Context) {
	id := ctx.Param("id")

	var req requests.LotAddItemsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	result, err := c.lotService.AddItems(ctx.Request.Context(), id, req, locationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Buah berhasil ditambahkan ke Lot",
		"data":    result,
	})
}

func (c *LotController) RemoveItem(ctx *gin.Context) {
	id := ctx.Param("id")

	var req requests.LotRemoveItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	err := c.lotService.RemoveItem(ctx.Request.Context(), id, req, locationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Buah berhasil dikeluarkan dari Lot",
	})
}

func (c *LotController) Finalize(ctx *gin.Context) {
	id := ctx.Param("id")

	// No body required for Finalize
	req := requests.LotFinalizeRequest{}

	userAuth := ctx.MustGet(authentication.Token).(requests.UserAuth)
	locationID := userAuth.LocationID

	result, err := c.lotService.Finalize(ctx.Request.Context(), id, req, locationID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Lot berhasil difinalisasi dan stok aktif",
		"data":    result,
	})
}
