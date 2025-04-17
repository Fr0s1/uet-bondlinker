package controller

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"socialnet/util"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"socialnet/config"
)

type FileController struct {
	s3Client *s3.S3
	cfg      *config.Config
}

// NewFileController creates a new file controller
func NewFileController(cfg *config.Config) *FileController {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AWS.Region),
	})

	if err != nil {
		panic(fmt.Sprintf("Failed to create AWS session: %v", err))
	}

	s3Client := s3.New(sess)

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
	validExts := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !validExts[ext] {
		util.RespondWithError(c, http.StatusBadRequest, "Invalid file type (only images allowed)")
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
	_, err = fc.s3Client.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(fc.cfg.AWS.Bucket),
		Key:         aws.String(fileKey),
		Body:        bytes.NewReader(fileContent),
		ContentType: aws.String(http.DetectContentType(fileContent)),
		ACL:         aws.String("public-read"),
	})

	if err != nil {
		util.RespondWithError(c, http.StatusInternalServerError, "Failed to upload file")
		return
	}

	// Generate file URL
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		fc.cfg.AWS.Bucket, fc.cfg.AWS.Region, fileKey)

	util.RespondWithSuccess(c, http.StatusOK, "File uploaded successfully", gin.H{
		"url":      fileURL,
		"filename": fileName,
	})
}
