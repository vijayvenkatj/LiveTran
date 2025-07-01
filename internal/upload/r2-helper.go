package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fsnotify/fsnotify"
)

type Uploader interface {
	Upload(ctx context.Context, bucket, key string, data []byte) error
}

type CloudflareUploader struct {
	client *s3.Client
}



func CreateCloudFlareUploader(accessKeyId string, accessKeySecret string, accountId string) (*CloudflareUploader,error) {

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




func (uploader *CloudflareUploader) UploadStream(ctx context.Context, bucket, key string, reader io.Reader, contentType string) error {
	_, err := uploader.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})
	return err
}

func detectContentType(key string) string {
	switch {
	case strings.HasSuffix(key, ".m3u8"):
		return "application/vnd.apple.mpegurl"
	case strings.HasSuffix(key, ".ts"):
		return "video/MP2T"
	default:
		return "application/octet-stream"
	}
}



func (uploader *CloudflareUploader) WatchAndUpload(ctx context.Context,outputDir string) {
	
	// Create a FSNOTIFY Watcher

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return 
	}
	defer watcher.Close()

	err = watcher.Add(outputDir)
	if err != nil {
		return
	}

	seenTS := make(map[string]bool)

	for {
		select {
		case <- ctx.Done():
			return

		case event, ok := <- watcher.Events:

			if !ok {
				return 
			}

			if event.Name == "" {
				continue
			}

			// Handle .ts files: only on CREATE
			if strings.HasSuffix(event.Name, ".ts") && event.Op&fsnotify.Create == fsnotify.Create {
				if seenTS[event.Name] {
					continue
				}
				seenTS[event.Name] = true

				file, err := os.Open(event.Name)
				if err != nil {
					fmt.Println("Issue uploading: ", event.Name)
					continue
				}

				go uploader.UploadStream(ctx, "testing", event.Name, file, detectContentType(event.Name))

			// Handle .m3u8 files: on CREATE or WRITE
			} else if strings.HasSuffix(event.Name, ".m3u8") && (event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write) {
				
				file, err := os.Open(event.Name)
				if err != nil {
					fmt.Println("Issue uploading: ", event.Name)
					continue
				}

				go uploader.UploadStream(ctx, "testing", event.Name, file, detectContentType(event.Name))
			}

		case err, ok := <-watcher.Errors:
			
			if !ok {
				return 
			}

			fmt.Println("ERROR:",err)
		}
	}
}