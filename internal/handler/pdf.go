package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/romapopov1212/robokp-pdf-service/internal/dto"
	"go.uber.org/zap"
	"net/http"
)

func (h *Controller) GeneratePdf(c *gin.Context) {
	var req dto.SaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "невалидный запрос"})
		return
	}
	
	//filename := "generated_" + strconv.FormatInt(req.UserId, 10) + ".pdf"
	
	err := h.pdfGenService.GenerateAdvancedPDFWithGofpdf(req)
	if err != nil {
		h.logger.Error("ошибка генерации PDF", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка генерации PDF"})
		return
	}
	
	//c.FileAttachment(filename, "document.pdf") // если раскоментить будет ошибка 404 потому что он попыта
}

func (h *Controller) SavePdf(c *gin.Context) {
	var req dto.SaveRequest
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	
	logoJson, err := json.Marshal(req.Logo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации logo"})
		return
	}
	
	execParamsJson, err := json.Marshal(req.ExecutorParameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации executor_parameters"})
		return
	}
	
	presentationJson, err := json.Marshal(req.PresentationParameters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации presentation_parameters"})
		return
	}
	
	styleTemplateJson, err := json.Marshal(req.StyleTemplate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ошибка сериализации style_template"})
		return
	}
	
	err = h.pdfService.SavePdf(c.Request.Context(), dto.SavePdfRequest{
		UserId:                 req.UserId,
		CartId:                 req.CartId,
		PublicationId:          req.PublicationId,
		Logo:                   logoJson,
		ExecutorParameters:     execParamsJson,
		PresentationParameters: presentationJson,
		StyleTemplate:          styleTemplateJson,
		Count:                  req.Count,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "не удалось сохранить данные pdf"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"status": "успешно сохранено"})
}
