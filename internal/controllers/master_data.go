package controllers

import (
	"net/http"

	"durich-be/internal/dto/requests"
	"durich-be/internal/services"
	"durich-be/pkg/errors"
	"durich-be/pkg/http/response"
	"durich-be/pkg/utils"

	"github.com/gin-gonic/gin"
)

type MasterDataController struct {
	service services.MasterDataService
}

func NewMasterDataController(service services.MasterDataService) MasterDataController {
	return MasterDataController{service: service}
}

func (c *MasterDataController) CreateCompany(ctx *gin.Context) {
	var req requests.CompanyCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.CreateCompany(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusCreated, "Company created successfully", result)
}

func (c *MasterDataController) GetCompanies(ctx *gin.Context) {
	result, err := c.service.GetCompanies(ctx.Request.Context())
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Companies retrieved successfully", result)
}

func (c *MasterDataController) GetCompanyByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetCompanyByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Company retrieved successfully", result)
}

func (c *MasterDataController) UpdateCompany(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.CompanyUpdateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.UpdateCompany(ctx.Request.Context(), id, req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Company updated successfully", result)
}

func (c *MasterDataController) DeleteCompany(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteCompany(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Company deleted successfully", nil)
}

func (c *MasterDataController) CreateEstate(ctx *gin.Context) {
	var req requests.EstateCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.CreateEstate(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusCreated, "Estate created successfully", result)
}

func (c *MasterDataController) GetEstates(ctx *gin.Context) {
	companyID := ctx.Query("company_id")
	result, err := c.service.GetEstates(ctx.Request.Context(), companyID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Estates retrieved successfully", result)
}

func (c *MasterDataController) GetEstateByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetEstateByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Estate retrieved successfully", result)
}

func (c *MasterDataController) UpdateEstate(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.EstateUpdateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.UpdateEstate(ctx.Request.Context(), id, req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Estate updated successfully", result)
}

func (c *MasterDataController) DeleteEstate(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteEstate(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Estate deleted successfully", nil)
}

func (c *MasterDataController) CreateDivisi(ctx *gin.Context) {
	var req requests.DivisiCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.CreateDivisi(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusCreated, "Divisi created successfully", result)
}

func (c *MasterDataController) GetDivisiList(ctx *gin.Context) {
	estateID := ctx.Query("estate_id")
	result, err := c.service.GetDivisiList(ctx.Request.Context(), estateID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Divisi list retrieved successfully", result)
}

func (c *MasterDataController) GetDivisiByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetDivisiByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Divisi retrieved successfully", result)
}

func (c *MasterDataController) UpdateDivisi(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.DivisiUpdateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.UpdateDivisi(ctx.Request.Context(), id, req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Divisi updated successfully", result)
}

func (c *MasterDataController) DeleteDivisi(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteDivisi(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Divisi deleted successfully", nil)
}

func (c *MasterDataController) CreateBlok(ctx *gin.Context) {
	var req requests.BlokCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.CreateBlok(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusCreated, "Blok created successfully", result)
}

func (c *MasterDataController) GetBloks(ctx *gin.Context) {
	divisiID := ctx.Query("divisi_id")
	result, err := c.service.GetBloks(ctx.Request.Context(), divisiID)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Bloks retrieved successfully", result)
}

func (c *MasterDataController) GetBlokByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetBlokByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Blok retrieved successfully", result)
}

func (c *MasterDataController) UpdateBlok(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.BlokUpdateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.UpdateBlok(ctx.Request.Context(), id, req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Blok updated successfully", result)
}

func (c *MasterDataController) DeleteBlok(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteBlok(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Blok deleted successfully", nil)
}

func (c *MasterDataController) CreateJenisDurian(ctx *gin.Context) {
	var req requests.JenisDurianCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.CreateJenisDurian(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusCreated, "Jenis Durian created successfully", result)
}

func (c *MasterDataController) GetJenisDurianList(ctx *gin.Context) {
	result, err := c.service.GetJenisDurianList(ctx.Request.Context())
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Jenis Durian list retrieved successfully", result)
}

func (c *MasterDataController) GetJenisDurianByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetJenisDurianByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Jenis Durian retrieved successfully", result)
}

func (c *MasterDataController) UpdateJenisDurian(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.JenisDurianUpdateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.UpdateJenisDurian(ctx.Request.Context(), id, req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Jenis Durian updated successfully", result)
}

func (c *MasterDataController) DeleteJenisDurian(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeleteJenisDurian(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Jenis Durian deleted successfully", nil)
}

func (c *MasterDataController) CreatePohon(ctx *gin.Context) {
	var req requests.PohonCreateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.CreatePohon(ctx.Request.Context(), req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusCreated, "Pohon created successfully", result)
}

func (c *MasterDataController) GetPohonList(ctx *gin.Context) {
	result, err := c.service.GetPohonList(ctx.Request.Context())
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Pohon list retrieved successfully", result)
}

func (c *MasterDataController) GetPohonByID(ctx *gin.Context) {
	id := ctx.Param("id")
	result, err := c.service.GetPohonByID(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Pohon retrieved successfully", result)
}

func (c *MasterDataController) UpdatePohon(ctx *gin.Context) {
	id := ctx.Param("id")
	var req requests.PohonUpdateRequest
	if err := utils.BindData(ctx, &req); err != nil {
		response.SendError(ctx, errors.ValidationErrorToAppError(err))
		return
	}
	result, err := c.service.UpdatePohon(ctx.Request.Context(), id, req)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Pohon updated successfully", result)
}

func (c *MasterDataController) DeletePohon(ctx *gin.Context) {
	id := ctx.Param("id")
	err := c.service.DeletePohon(ctx.Request.Context(), id)
	if err != nil {
		response.SendError(ctx, err)
		return
	}
	response.SendSuccess(ctx, http.StatusOK, "Pohon deleted successfully", nil)
}
