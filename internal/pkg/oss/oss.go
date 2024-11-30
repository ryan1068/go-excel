package oss

import (
	"cst/internal/pkg/config"
	"fmt"
	aliyunOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"path"
	"time"
)

type oss struct {
	cfg *config.Config
}

func New(cfg *config.Config) *oss {
	return &oss{
		cfg: cfg,
	}
}

// client 获取oss客户端
func (o *oss) client() (*aliyunOss.Client, error) {
	client, err := aliyunOss.New(o.cfg.Oss.Endpoint, o.cfg.Oss.AccessKeyId, o.cfg.Oss.AccessKeySecret)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// bucket 获取bucket
func (o *oss) bucket() (*aliyunOss.Bucket, error) {
	client, err := o.client()
	if err != nil {
		return nil, err
	}

	bucket, err := client.Bucket(o.cfg.Oss.BucketName)
	if err != nil {
		return nil, err
	}
	return bucket, nil
}

// UploadFile 上传文件
func (o *oss) UploadFile(filePath string) (string, error) {
	bucket, err := o.bucket()
	if err != nil {
		return "", err
	}

	objectKey := fmt.Sprintf("%v/%v/%v", "mirco", time.Now().Format("20060102"), path.Base(filePath))
	err = bucket.PutObjectFromFile(objectKey, filePath)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%v.%v/%v", o.cfg.Oss.BucketName, o.cfg.Oss.Endpoint, objectKey), nil
}

// DownloadFile 下载文件
func (o *oss) DownloadFile(objectKey, filePath string) (string, error) {
	bucket, err := o.bucket()
	if err != nil {
		return "", err
	}

	err = bucket.GetObjectToFile(objectKey, filePath)
	if err != nil {
		return "", err
	}

	return filePath, nil
}

// DeleteObject 删除文件
func (o *oss) DeleteObject(objectKey string) error {
	bucket, err := o.bucket()
	if err != nil {
		return err
	}

	err = bucket.DeleteObject(objectKey)
	if err != nil {
		return err
	}

	return nil
}
