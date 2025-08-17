package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/romapopov1212/robokp-pdf-service/internal/pdfgen"
	"github.com/romapopov1212/robokp-pdf-service/internal/service"
	"go.uber.org/zap"
)

type Controller struct {
	pdfGenService *pdfgen.Page
	pdfService    *service.PdfService
	router        *gin.Engine
	logger        *zap.Logger
}

func RegisterRoutes(pdfService *service.PdfService, router *gin.Engine, logger *zap.Logger, pdfGenService *pdfgen.Page) Controller {
	cntrl := Controller{
		pdfService:    pdfService,
		router:        router,
		logger:        logger,
		pdfGenService: pdfGenService,
	}
	
	cntrl.router.POST("api/v1/pdf", cntrl.SavePdf)
	cntrl.router.POST("api/v1/pdfGen", cntrl.GeneratePdf)
	
	return cntrl
}
