package service

import (
	"context"
	"fmt"
	"github.com/romapopov1212/robokp-pdf-service/internal/dto"
	"github.com/romapopov1212/robokp-pdf-service/internal/pdfgen"
	"github.com/romapopov1212/robokp-pdf-service/internal/repository"
	"go.uber.org/zap"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type PdfService struct {
	pdfGen   *pdfgen.Page
	pdfRepo  *repository.PdfRepository
	logger   *zap.Logger
	s3Client *s3.Client
}

func NewPdfService(pdfRepo *repository.PdfRepository, logger *zap.Logger, s3Client *s3.Client, pdfGen *pdfgen.Page) *PdfService {
	return &PdfService{
		pdfRepo:  pdfRepo,
		logger:   logger,
		s3Client: s3Client,
		pdfGen:   pdfGen,
	}
}

func (s *PdfService) SavePdf(
	ctx context.Context, request dto.SavePdfRequest) error {
	
	err := s.pdfRepo.Save(ctx,
		request.UserId,
		request.CartId,
		request.PublicationId,
		request.Logo,
		request.ExecutorParameters,
		request.PresentationParameters,
		request.StyleTemplate,
		request.Count)
	if err != nil {
		s.logger.Error("ошибка при сохранении пдф", zap.Error(err))
		return fmt.Errorf("ошибка при сохрании: %w", err)
	}
	return nil
}
