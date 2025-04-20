package controller

import (
	"bytes"
	"context"
	"fmt"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"socialnet/config"
	"socialnet/util"
	"strings"
	"time"
)

var validExts = map[string]struct{}{".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {}, ".webp": {}}

type FileController struct {
	s3Client *s3.Client
	cfg      *config.Config
}

// NewFileController creates a new file controller
func NewFileController(cfg *config.Config) *FileController {
	c, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKeyID, cfg.AWS.SecretAccessKey, "")),
		awsconfig.WithRegion(cfg.AWS.Region),
		awsconfig.WithBaseEndpoint(cfg.AWS.Endpoint),
	)

	if err != nil {
		panic(fmt.Sprintf("Failed to create AWS session: %v", err))
	}

	s3Client := s3.NewFromConfig(c)

	return &FileController{
		s3Client: s3Client,
		cfg:      cfg,
	}
}

// UploadFile handles file upload and returns the file URL
func (fc *FileController) UploadFile(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		util.RespondWithError(c, http.StatusBadRequest, "No file provided")
		return
	}
	defer file.Close()

	// Check file size (limit to 5MB)
	if header.Size > 5*1024*1024 {
		util.RespondWithError(c, http.StatusBadRequest, "File too large (max 5MB)")
		return
	}

	// Validate file type
	ext := strings.ToLower(filepath.Ext(header.Filename))
	_, found := validExts[ext]
	if !found {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid file type (only .jpg, .png, .gif and .webp)")
		return
	}

	// Read file content
	fileContent, err := io.ReadAll(file)
	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to read file")
		return
	}

	// Generate unique filename
	uniqueID := uuid.New().String()
	fileName := fmt.Sprintf("%s-%s%s", time.Now().Format("20060102"), uniqueID, ext)
	fileKey := fmt.Sprintf("uploads/%s", fileName)

	// Upload to S3
	contentType := http.DetectContentType(fileContent)
	_, err = fc.s3Client.PutObject(c, &s3.PutObjectInput{
		Bucket:      &fc.cfg.AWS.Bucket,
		Key:         &fileKey,
		Body:        bytes.NewReader(fileContent),
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})

	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to upload file")
		return
	}

	cdnURL := fc.cfg.AWS.CdnURL
	if cdnURL == "" {
		cdnURL = fmt.Sprintf("https://%s.s3.%s.amazonaws.com", fc.cfg.AWS.Bucket, fc.cfg.AWS.Region)
	}
	// Generate file URL
	fileURL := fmt.Sprintf("%s/%s", cdnURL, fileKey)

	util.RespondWithSuccess(c, http.StatusOK, "File uploaded successfully", gin.H{
		"url":      fileURL,
		"filename": fileName,
	})
}
