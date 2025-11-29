package routes

import (
	"durich-be/internal/controllers"
	"durich-be/internal/domain"
	"durich-be/pkg/http/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterMasterData(router *gin.RouterGroup, ctl controllers.MasterDataController) {
	companyGroup := router.Group("/companies")
	companyGroup.Use(middlewares.TokenAuthMiddleware())
	{
		companyGroup.POST("", middlewares.RoleHandler(domain.RoleAdmin), ctl.CreateCompany)
		companyGroup.GET("", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetCompanies)
		companyGroup.GET("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetCompanyByID)
		companyGroup.PUT("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.UpdateCompany)
		companyGroup.DELETE("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.DeleteCompany)
	}

	estateGroup := router.Group("/estates")
	estateGroup.Use(middlewares.TokenAuthMiddleware())
	{
		estateGroup.POST("", middlewares.RoleHandler(domain.RoleAdmin), ctl.CreateEstate)
		estateGroup.GET("", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetEstates)
		estateGroup.GET("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetEstateByID)
		estateGroup.PUT("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.UpdateEstate)
		estateGroup.DELETE("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.DeleteEstate)
	}

	divisiGroup := router.Group("/divisi")
	divisiGroup.Use(middlewares.TokenAuthMiddleware())
	{
		divisiGroup.POST("", middlewares.RoleHandler(domain.RoleAdmin), ctl.CreateDivisi)
		divisiGroup.GET("", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetDivisiList)
		divisiGroup.GET("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetDivisiByID)
		divisiGroup.PUT("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.UpdateDivisi)
		divisiGroup.DELETE("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.DeleteDivisi)
	}

	blokGroup := router.Group("/bloks")
	blokGroup.Use(middlewares.TokenAuthMiddleware())
	{
		blokGroup.POST("", middlewares.RoleHandler(domain.RoleAdmin), ctl.CreateBlok)
		blokGroup.GET("", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetBloks)
		blokGroup.GET("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetBlokByID)
		blokGroup.PUT("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.UpdateBlok)
		blokGroup.DELETE("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.DeleteBlok)
	}

	jenisDurianGroup := router.Group("/jenis-durian")
	jenisDurianGroup.Use(middlewares.TokenAuthMiddleware())
	{
		jenisDurianGroup.POST("", middlewares.RoleHandler(domain.RoleAdmin), ctl.CreateJenisDurian)
		jenisDurianGroup.GET("", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetJenisDurianList)
		jenisDurianGroup.GET("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetJenisDurianByID)
		jenisDurianGroup.PUT("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.UpdateJenisDurian)
		jenisDurianGroup.DELETE("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.DeleteJenisDurian)
	}

	pohonGroup := router.Group("/pohon")
	pohonGroup.Use(middlewares.TokenAuthMiddleware())
	{
		pohonGroup.POST("", middlewares.RoleHandler(domain.RoleAdmin), ctl.CreatePohon)
		pohonGroup.GET("", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetPohonList)
		pohonGroup.GET("/:id", middlewares.RoleHandler(domain.RoleAdmin, domain.RoleWarehouse), ctl.GetPohonByID)
		pohonGroup.PUT("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.UpdatePohon)
		pohonGroup.DELETE("/:id", middlewares.RoleHandler(domain.RoleAdmin), ctl.DeletePohon)
	}
}
