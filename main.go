package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func main() {
	// Load environment variables
	bucketName := os.Getenv("S3_BUCKET_NAME")
	if bucketName == "" {
		log.Fatalf("S3_BUCKET_NAME environment variable is not set")
	}

	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1" // Default region
	}

	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithHTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}))

	// Create S3 client
	s3Client := s3.NewFromConfig(cfg)

	// Create uploader
	uploader := manager.NewUploader(s3Client)

	// File content
	content := "Hello, World!"
	fileKey := "hello-world.txt"

	// Upload file
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
		Body:   strings.NewReader(content),
	})
	if err != nil {
		log.Fatalf("failed to upload file: %v", err)
	}

	fmt.Printf("File '%s' successfully uploaded to bucket '%s'\n", fileKey, bucketName)
}
