package controllers

import (
	"net/http"
	"strconv"

	"durich-be/internal/dto/requests"
	"durich-be/internal/services"
	"durich-be/pkg/errors"
	"durich-be/pkg/http/response"
	"durich-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type BuahRawController struct {
	service services.BuahRawService
}

func NewBuahRawController(service services.BuahRawService) BuahRawController {
	return BuahRawController{service: service}
}

func (c *BuahRawController) Create(ctx *gin.Context) {
	var req requests.BuahRawCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	id, err := c.service.Create(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusCreated, "Berhasil menyimpan data panen", gin.H{"id": id})
}

func (c *BuahRawController) BulkCreate(ctx *gin.Context) {
	var req requests.BuahRawBulkCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	ids, err := c.service.BulkCreate(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusCreated, "Berhasil menyimpan data panen", gin.H{"inserted_ids": ids})
}

func (c *BuahRawController) GetList(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(ctx.DefaultQuery("limit", "10"))

	filter := map[string]interface{}{}
	if v := ctx.Query("tgl_panen"); v != "" {
		filter["tgl_panen"] = v
	}
	if v := ctx.Query("blok_panen_id"); v != "" {
		filter["blok_panen_id"] = v
	}
	if v := ctx.Query("jenis_durian_id"); v != "" {
		filter["jenis_durian_id"] = v
	}
	if v := ctx.Query("is_sorted"); v != "" {
		b, _ := strconv.ParseBool(v)
		filter["is_sorted"] = b
	}

	res, err := c.service.GetList(ctx.Request.Context(), filter, limit, page)
	if err != nil {
		response.SendError(ctx, err)
		return
	}

	response.SendSuccess(ctx, http.StatusOK, "Success", res)
}

func (c *BuahRawController) GetDetail(ctx *gin.Context) {
	id := ctx.Param("id")
	res, err := c.service.GetDetail(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Success", res)
}

func (c *BuahRawController) Update(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.BuahRawUpdateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}

	err := c.service.Update(ctx.Request.Context(), id, req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Berhasil update data buah", nil)
}

func (c *BuahRawController) Delete(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.Delete(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Data buah berhasil dihapus", nil)
}
