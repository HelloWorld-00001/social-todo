package uploadProvider

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/coderconquerer/go-login-app/internal/common"
	"github.com/coderconquerer/go-login-app/pkg/config"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Provider struct {
	BucketName string `json:"bucket_name"`
	Region     string `json:"region"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Domain     string `json:"domain"`
	Session    *session.Session
}

func NewS3ProviderWithConfig(config *config.AWSConfig) *S3Provider {
	provider := &S3Provider{
		BucketName: config.Bucket,
		Region:     config.Region,
		AccessKey:  config.PublicKey,
		SecretKey:  config.SecretKey,
		Domain:     config.Domain,
	}

	s3Session, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region),
		Credentials: credentials.NewStaticCredentials(
			config.PublicKey, config.SecretKey, ""),
	})

	if err != nil {
		log.Fatalln(fmt.Errorf("error creating S3 session: %v", err))
	}
	provider.Session = s3Session
	return provider
}

func NewS3Provider(bucketName, region, accessKey, secretKey, domain string) *S3Provider {
	provider := &S3Provider{
		BucketName: bucketName,
		Region:     region,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		Domain:     domain,
	}

	s3Session, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
		Credentials: credentials.NewStaticCredentials(
			accessKey, secretKey, ""),
	})

	if err != nil {
		log.Fatalln(fmt.Errorf("error creating S3 session: %v", err))
	}
	provider.Session = s3Session
	return provider
}

func (s *S3Provider) SaveFileUpload(ctx context.Context, data []byte, destination string) (*common.Image, error) {
	fileBytes := bytes.NewReader(data)
	fileType := http.DetectContentType(data)

	input := &s3.PutObjectInput{
		Bucket:      aws.String(s.BucketName),
		Key:         aws.String(destination),
		Body:        fileBytes,
		ACL:         aws.String("private"),
		ContentType: aws.String(fileType),
	}

	_, err := s3.New(s.Session).PutObject(input)

	if err != nil {
		return nil, err
	}

	// Construct the full URL to access the file
	url := fmt.Sprintf("%s/%s", s.Domain, destination)

	return &common.Image{
		Url:       url,
		CloudName: "S3",
	}, nil
}
