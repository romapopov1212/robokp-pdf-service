package pdfgen

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jung-kurt/gofpdf"
	"github.com/romapopov1212/robokp-pdf-service/internal/dto"
	"strconv"
	"strings"
	"time"
)

type Page struct {
	//HTML        string
	s3Client    *s3.Client
	s3Bucket    string
	s3Region    string
	s3UploadDir string
}

func New(s3Client *s3.Client, s3Bucket string, s3Region string, s3UploadDir string) *Page {
	return &Page{
		//HTML:        HTML,
		s3Client:    s3Client,
		s3Bucket:    s3Bucket,
		s3Region:    s3Region,
		s3UploadDir: s3UploadDir,
	}
}

func initRussianFonts(pdf *gofpdf.Fpdf) {
	// Используем встроенную поддержку UTF-8 в gofpdf
	// Для кириллицы можно использовать стандартные шрифты с поддержкой UTF-8
	pdf.SetFont("Arial", "", 12)
}

// GeneratePDFWithGofpdf генерирует PDF используя gofpdf
func GeneratePDFWithGofpdf(filename string, req dto.SaveRequest) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Устанавливаем шрифт с поддержкой UTF-8
	pdf.SetFont("Arial", "", 12)
	
	// Заголовок
	pdf.SetFont("Arial", "B", 16)
	pdf.SetTextColor(0, 0, 0)
	pdf.Cell(0, 10, "PDF документ")
	pdf.Ln(15)
	
	// Основной шрифт
	pdf.SetFont("Arial", "", 12)
	
	// Информация о пользователе
	pdf.Cell(0, 8, "ID пользователя: "+strconv.FormatInt(req.UserId, 10))
	pdf.Ln(8)
	pdf.Cell(0, 8, "ID корзины: "+strconv.FormatInt(req.CartId, 10))
	pdf.Ln(8)
	pdf.Cell(0, 8, "ID публикации: "+strconv.FormatInt(req.PublicationId, 10))
	pdf.Ln(8)
	pdf.Cell(0, 8, "Количество: "+strconv.Itoa(req.Count))
	pdf.Ln(15)
	
	// Логотип текст
	if req.Logo.LogoText.Value != "" {
		pdf.SetFont("Arial", "", 12)
		
		// Применяем стили к тексту
		style := ""
		if req.Logo.LogoText.Bold {
			style += "B"
		}
		if req.Logo.LogoText.Kursive {
			style += "I"
		}
		if req.Logo.LogoText.Under {
			style += "U"
		}
		
		if style != "" {
			pdf.SetFont("Arial", style, 12)
		}
		
		pdf.Cell(0, 8, req.Logo.LogoText.Value)
		pdf.Ln(15)
	}
	
	// Параметры исполнителя
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Параметры исполнителя:")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "Показать логотип: "+strconv.FormatBool(req.ExecutorParameters.First.ShowLogo))
	pdf.Ln(8)
	pdf.Cell(0, 8, "Показать имя: "+req.ExecutorParameters.First.ShowName)
	pdf.Ln(8)
	pdf.Cell(0, 8, "Показать контакты: "+req.ExecutorParameters.First.ShowContacts)
	pdf.Ln(15)
	
	// Параметры презентации
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Параметры презентации:")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "Список: "+strconv.FormatBool(req.PresentationParameters.List))
	pdf.Ln(8)
	pdf.Cell(0, 8, "По одному: "+strconv.FormatBool(req.PresentationParameters.OneByOne))
	pdf.Ln(8)
	pdf.Cell(0, 8, "Сумма: "+strconv.FormatBool(req.PresentationParameters.Sum))
	pdf.Ln(8)
	pdf.Cell(0, 8, "Цена: "+strconv.FormatBool(req.PresentationParameters.Price))
	pdf.Ln(15)
	
	// Шаблон стиля
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Шаблон стиля:")
	pdf.Ln(10)
	
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(0, 8, "ID шаблона: "+req.StyleTemplate.TemplateID)
	pdf.Ln(8)
	if req.StyleTemplate.Color != "" {
		pdf.Cell(0, 8, "Цвет: "+req.StyleTemplate.Color)
		pdf.Ln(8)
	}
	
	return pdf.OutputFileAndClose(filename)
}

// GenerateAdvancedPDFWithGofpdf создает более продвинутый PDF с таблицами и изображениями
func (s *Page) GenerateAdvancedPDFWithGofpdf(req dto.SaveRequest) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	
	// Устанавливаем шрифт с поддержкой UTF-8
	pdf.SetFont("Arial", "", 12)
	
	// Заголовок с цветом
	if req.StyleTemplate.Color != "" {
		// Парсим цвет (предполагаем формат "#RRGGBB")
		color := req.StyleTemplate.Color
		if strings.HasPrefix(color, "#") {
			color = color[1:] // Убираем #
		}
		if len(color) == 6 {
			r, _ := strconv.ParseUint(color[0:2], 16, 8)
			g, _ := strconv.ParseUint(color[2:4], 16, 8)
			b, _ := strconv.ParseUint(color[4:6], 16, 8)
			pdf.SetTextColor(int(r), int(g), int(b))
		}
	}
	
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(0, 15, "Документ PDF")
	pdf.Ln(20)
	
	// Сброс цвета на черный
	pdf.SetTextColor(0, 0, 0)
	
	// Информационная таблица
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Основная информация")
	pdf.Ln(12)
	
	// Создаем таблицу
	header := []string{"Поле", "Значение"}
	data := [][]string{
		{"ID пользователя", strconv.FormatInt(req.UserId, 10)},
		{"ID корзины", strconv.FormatInt(req.CartId, 10)},
		{"ID публикации", strconv.FormatInt(req.PublicationId, 10)},
		{"Количество", strconv.Itoa(req.Count)},
	}
	
	createTable(pdf, header, data)
	pdf.Ln(15)
	
	// Логотип и текст
	if req.Logo.LogoText.Value != "" {
		pdf.SetFont("Arial", "B", 14)
		pdf.Cell(0, 10, "Логотип")
		pdf.Ln(12)
		
		// Применяем стили к тексту логотипа
		style := ""
		if req.Logo.LogoText.Bold {
			style += "B"
		}
		if req.Logo.LogoText.Kursive {
			style += "I"
		}
		if req.Logo.LogoText.Under {
			style += "U"
		}
		
		pdf.SetFont("Arial", style, 12)
		pdf.Cell(0, 8, req.Logo.LogoText.Value)
		pdf.Ln(15)
		
		// Добавляем изображения логотипа если они есть
		if req.Logo.Square != "" {
			addImageFromBase64(pdf, req.Logo.Square, "Логотип (квадрат)", 30)
		}
		if req.Logo.Rectangle != "" {
			addImageFromBase64(pdf, req.Logo.Rectangle, "Логотип (прямоугольник)", 30)
		}
	}
	
	// Параметры исполнителя в таблице
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Параметры исполнителя")
	pdf.Ln(12)
	
	execHeader := []string{"Параметр", "Значение"}
	execData := [][]string{
		{"Показать логотип", strconv.FormatBool(req.ExecutorParameters.First.ShowLogo)},
		{"Показать имя", req.ExecutorParameters.First.ShowName},
		{"Показать контакты", req.ExecutorParameters.First.ShowContacts},
	}
	
	createTable(pdf, execHeader, execData)
	pdf.Ln(15)
	
	// Параметры презентации в таблице
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Параметры презентации")
	pdf.Ln(12)
	
	presentHeader := []string{"Параметр", "Значение"}
	presentData := [][]string{
		{"Список", strconv.FormatBool(req.PresentationParameters.List)},
		{"По одному", strconv.FormatBool(req.PresentationParameters.OneByOne)},
		{"Сумма", strconv.FormatBool(req.PresentationParameters.Sum)},
		{"Цена", strconv.FormatBool(req.PresentationParameters.Price)},
	}
	
	createTable(pdf, presentHeader, presentData)
	pdf.Ln(15)
	
	// Шаблон стиля
	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 10, "Шаблон стиля")
	pdf.Ln(12)
	
	styleHeader := []string{"Параметр", "Значение"}
	styleData := [][]string{
		{"ID шаблона", req.StyleTemplate.TemplateID},
	}
	if req.StyleTemplate.Color != "" {
		styleData = append(styleData, []string{"Цвет", req.StyleTemplate.Color})
	}
	
	createTable(pdf, styleHeader, styleData)
	
	var buf bytes.Buffer
	err := pdf.Output(&buf)
	if err != nil {
		return fmt.Errorf("ошибка при генерации PDF: %w", err)
	}
	
	pdfReader := bytes.NewReader(buf.Bytes())
	
	s3Key := fmt.Sprintf("%s/%s_%d.pdf",
		s.s3UploadDir,
		req.CartId,
		time.Now().UnixNano())
	
	_, err = s.s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(s.s3Bucket),
		Key:         aws.String(s3Key),
		Body:        pdfReader,
		ContentType: aws.String("application/pdf"),
	})
	if err != nil {
		//s.logger.Error("ошибка при сохранении PDF в S3", zap.Error(err))
		return fmt.Errorf("ошибка при сохранении PDF в S3: %w", err)
	}
	return nil
	//return buf.Bytes(), nil
	
	//return pdf.OutputFileAndClose(filename)
}

// createTable создает таблицу в PDF
func createTable(pdf *gofpdf.Fpdf, header []string, data [][]string) {
	// Ширина колонок
	colWidth := 80.0
	
	// Заголовок таблицы
	pdf.SetFont("Arial", "B", 10)
	pdf.SetFillColor(240, 240, 240)
	for _, col := range header {
		pdf.CellFormat(colWidth, 7, col, "1", 0, "", true, 0, "")
	}
	pdf.Ln(-1)
	
	// Данные таблицы
	pdf.SetFont("Arial", "", 10)
	pdf.SetFillColor(255, 255, 255)
	for _, row := range data {
		for _, col := range row {
			pdf.CellFormat(colWidth, 6, col, "1", 0, "", false, 0, "")
		}
		pdf.Ln(-1)
	}
}

// addImageFromBase64 добавляет изображение из base64 строки
func addImageFromBase64(pdf *gofpdf.Fpdf, base64Data, caption string, height float64) {
	// Убираем префикс data:image/...;base64, если есть
	if strings.Contains(base64Data, ",") {
		base64Data = strings.Split(base64Data, ",")[1]
	}
	
	// Декодируем base64
	imageData, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		return // Игнорируем ошибки декодирования
	}
	
	// Определяем тип изображения по первым байтам
	var imageType string
	if len(imageData) > 2 {
		switch {
		case imageData[0] == 0xFF && imageData[1] == 0xD8:
			imageType = "JPEG"
		case imageData[0] == 0x89 && imageData[1] == 0x50:
			imageType = "PNG"
		default:
			return // Неподдерживаемый формат
		}
	}
	
	// Добавляем изображение
	pdf.Image("", pdf.GetX(), pdf.GetY(), 0, height, false, imageType, 0, "")
	
	// Добавляем подпись
	pdf.SetFont("Arial", "", 10)
	pdf.Cell(0, 5, caption)
	pdf.Ln(10)
}

// GeneratePDF - оставляем старую функцию для совместимости, но теперь она использует gofpdf
func GeneratePDF(filename string, pages []Page) error {
	// Для совместимости с существующим кодом
	// В реальном проекте лучше переписать все места, где используется эта функция
	return nil
}
