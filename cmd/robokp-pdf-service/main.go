package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	conf "github.com/romapopov1212/robokp-pdf-service/internal/config"
	db2 "github.com/romapopov1212/robokp-pdf-service/internal/db"
	"github.com/romapopov1212/robokp-pdf-service/internal/handler"
	"github.com/romapopov1212/robokp-pdf-service/internal/pdfgen"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/romapopov1212/robokp-pdf-service/internal/repository"
	"github.com/romapopov1212/robokp-pdf-service/internal/service"
	"go.uber.org/zap"
	"log"
)

func main() {
	configPath := flag.String("config", "./config", "path to the config file")
	
	flag.Parse()
	
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("error init logger: %v", err)
	}
	
	logger.Info("starting app")
	
	cfg, err := conf.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config")
	}
	
	db, _, err := db2.NewDatabaseConnection(cfg.Database)
	if err != nil {
		log.Fatalf("error init database connection: %v", err)
	}
	
	repo, err := repository.New(db)
	if err != nil {
		log.Fatalf("error create table: %v", err)
	}
	
	router := gin.Default()
	
	awsCfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(cfg.AWS.Region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.AWS.AccessKeyID,
			cfg.AWS.SecretAccessKey,
			"",
		)),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:           cfg.AWS.EndpointUri,
					SigningRegion: cfg.AWS.Region,
				}, nil
			}),
		),
	)
	
	if err != nil {
		log.Fatalf("error loading AWS config: %v", err)
	}
	
	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})
	
	pd := pdfgen.New(s3Client, cfg.AWS.Bucket, cfg.AWS.Region, cfg.AWS.UploadDir)
	
	fmt.Println(cfg.AWS.Bucket, cfg.AWS.Region, cfg.AWS.UploadDir)
	
	srv := service.NewPdfService(repo, logger, s3Client, pd)
	
	handler.RegisterRoutes(srv, router, logger, pd)
	
	servAddr := cfg.Address
	
	logger.Info("stating server", zap.String("address", servAddr))
	if err := router.Run(servAddr); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}
