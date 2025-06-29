package upload

import (
	"bytes"
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type Uploader interface {
	Upload(ctx context.Context, bucket, key string, data []byte) error
}

type CloudflareUploader struct {
	client *s3.Client
}



func createCloudFlareUploader(accessKeyId string, accessKeySecret string, accountId string) (*CloudflareUploader,error) {

	// TODO: Add timeout logic here
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		fmt.Println(err)
		return nil,err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return &CloudflareUploader{client:client},nil
}




func (uploader *CloudflareUploader) Upload(ctx context.Context, bucket, key string, data []byte) (error) {
	
	_, err := uploader.client.PutObject(ctx,&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key: aws.String(key),
		Body: bytes.NewReader(data),
	})

	return err
}


func (uploader *CloudflareUploader) WatchAndUpload(outputDir string) {
	
}