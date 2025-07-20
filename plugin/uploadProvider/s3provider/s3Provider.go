package s3provider

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/coderconquerer/social-todo/common"
	"github.com/coderconquerer/social-todo/configs"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

type S3Provider struct {
	s3prefix   string
	BucketName string `json:"bucket_name"`
	Region     string `json:"region"`
	AccessKey  string `json:"access_key"`
	SecretKey  string `json:"secret_key"`
	Domain     string `json:"domain"`
	Session    *session.Session
}

func NewS3ProviderWithConfig(config *configs.AWSConfig) *S3Provider {
	provider := &S3Provider{
		BucketName: config.Bucket,
		Region:     config.Region,
		AccessKey:  config.PublicKey,
		SecretKey:  config.SecretKey,
		Domain:     config.Domain,
		s3prefix:   config.S3Prefix,
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

func (s *S3Provider) GetPrefix() string {
	return s.s3prefix
}

func (s *S3Provider) Name() string {
	return "s3-upload-provider"
}

func (s *S3Provider) Get() interface{} {
	return s
}

func (s *S3Provider) InitFlags() {
	// Bind config values from CLI (can be overwritten via flag.Parse())
	flag.StringVar(&s.s3prefix, "aws-s3-prefix", "aws-s3-s3prefix", "S3 s3prefix")
	flag.StringVar(&s.BucketName, "aws-s3-bucket", "", "S3 Bucket Name")
	flag.StringVar(&s.Region, "aws-s3-region", "ap-southeast-1", "AWS Region")
	flag.StringVar(&s.AccessKey, "aws-s3-public-key", "", "AWS Access Key")
	flag.StringVar(&s.SecretKey, "aws-s3-secret-key", "", "AWS Secret Key")
	flag.StringVar(&s.Domain, "aws-s3-domain", "", "S3 File Access Domain")
}

func (s *S3Provider) Configure() error {
	if s.s3prefix == "" {
		s.s3prefix = "aws-s3-prefix"
	}
	if s.BucketName == "" || s.Region == "" || s.AccessKey == "" || s.SecretKey == "" {
		return fmt.Errorf("missing required S3 configuration")
	}
	return nil
}

func (s *S3Provider) Run() error {
	return nil
}

func (s *S3Provider) Stop() <-chan bool {
	// No background processes, return closed channel
	ch := make(chan bool)
	go func() {
		ch <- true
	}()
	return ch
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
