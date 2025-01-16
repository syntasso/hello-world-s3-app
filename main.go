package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

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

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithHTTPClient(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}))

	// Load AWS configuration
	maxRetries := 120
	for i := 0; i < maxRetries; i++ {
		time.Sleep(time.Second)
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
		if err == nil {
			fmt.Printf("File uploaded successfully: %s\n", fileKey)
			time.Sleep(10000 * time.Second)
		}

		fmt.Printf("Failed to upload file: %v\n", err)
		if i == maxRetries-1 {
			log.Fatalf("Failed to upload file after %d retries", maxRetries)
		}
	}
}
