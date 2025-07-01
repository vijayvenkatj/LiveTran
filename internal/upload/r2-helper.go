package upload

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/fsnotify/fsnotify"
)

type Uploader interface {
	Upload(ctx context.Context, bucket, key string, data []byte) error
	UploadStream(ctx context.Context, bucket, key string, reader io.Reader, contentType string) error
}

type CloudflareUploader struct {
	client *s3.Client
}

func CreateCloudFlareUploader(ctx context.Context, accessKeyId string, accessKeySecret string, accountId string) (*CloudflareUploader, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
	})

	return &CloudflareUploader{client: client}, nil
}

func (uploader *CloudflareUploader) UploadStream(ctx context.Context, bucket, key string, reader io.Reader, contentType string) error {
	_, err := uploader.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		Body:        reader,
		ContentType: aws.String(contentType),
	})

	// TODO: DELETE AFTER Uploading TS and m3u8
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

func (uploader *CloudflareUploader) WatchAndUpload(ctx context.Context, outputDir string, bucket string) {
	var wg sync.WaitGroup

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
		case <-ctx.Done():
			wg.Wait()
			return

		case event, ok := <-watcher.Events:
			if !ok || event.Name == "" {
				continue
			}

			key := filepath.Base(event.Name)
			path := event.Name

			// Handle TS files (on create only)
			if strings.HasSuffix(event.Name, ".ts") && event.Op&fsnotify.Create == fsnotify.Create {
				if seenTS[event.Name] {
					continue
				}
				seenTS[event.Name] = true

				wg.Add(1)
				go func(path, key string) {
					defer wg.Done()

					time.Sleep(200 * time.Millisecond)
					file, err := os.Open(path)
					if err != nil {
						fmt.Println("Issue uploading:", path)
						return
					}
					defer file.Close()

					contentType := detectContentType(key)

					for i := 0; i < 3; i++ {
						err = uploader.UploadStream(ctx, bucket, key, file, contentType)
						if err == nil {
							fmt.Println("Upload Successful",key)
							break
						}
						time.Sleep(time.Second * time.Duration(i+1)) // exponential backoff
					}
				}(path, key)
			}

			// Handle M3U8 files (on create or update)
			if strings.HasSuffix(event.Name, ".m3u8") &&
				(event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write) {

				wg.Add(1)
				go func(path, key string) {
					defer wg.Done()

					time.Sleep(200 * time.Millisecond)
					file, err := os.Open(path)
					if err != nil {
						fmt.Println("Issue uploading:", path)
						return
					}
					defer file.Close()

					contentType := detectContentType(key)

					for i := 0; i < 3; i++ {
						err = uploader.UploadStream(ctx, bucket, key, file, contentType)
						if err == nil {
							fmt.Println("Upload Successful",key)
							break
						}
						time.Sleep(time.Second * time.Duration(i+1))
					}
				}(path, key)
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("ERROR:", err)
		}
	}
}
