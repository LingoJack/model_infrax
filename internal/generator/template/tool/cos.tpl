package tool

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// https://cloud.tencent.com/document/product/436/35059

// CosClient COS 客户端，复用 HTTP client 和配置
type CosClient struct {
	client *cos.Client
	cosUrl string
	ak     string
	sk     string
}

var (
	cosClientInstance *CosClient
	cosClientOnce     sync.Once
)

// NewCosClient 创建 COS 客户端实例
func NewCosClient(cosUrl, ak, sk string) (*CosClient, error) {
	u, err := url.Parse(cosUrl)
	if err != nil {
		return nil, fmt.Errorf("[NewCosClient] parse cos url error: %v", err)
	}

	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
		Timeout: 30 * time.Second, // 设置超时时间
	})

	return &CosClient{
		client: client,
		cosUrl: cosUrl,
		ak:     ak,
		sk:     sk,
	}, nil
}

// GetCosClient 获取全局单例 COS 客户端（懒加载）
func GetCosClient(cosUrl, ak, sk string) (*CosClient, error) {
	var err error
	cosClientOnce.Do(func() {
		cosClientInstance, err = NewCosClient(cosUrl, ak, sk)
	})
	return cosClientInstance, err
}

// UploadPresignedUrl 获取上传预签名 URL
func (c *CosClient) UploadPresignedUrl(ctx context.Context, cosKey string, expiration time.Duration) (string, error) {
	presignedURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodPut, cosKey, c.ak, c.sk, expiration, nil)
	if err != nil {
		return "", fmt.Errorf("[UploadPresignedUrl] get presigned url error: %v", err)
	}
	return presignedURL.String(), nil
}

// DownloadPresignedUrl 获取下载预签名 URL
func (c *CosClient) DownloadPresignedUrl(ctx context.Context, cosKey string, expiration time.Duration) (string, error) {
	presignedURL, err := c.client.Object.GetPresignedURL(ctx, http.MethodGet, cosKey, c.ak, c.sk, expiration, nil)
	if err != nil {
		return "", fmt.Errorf("[DownloadPresignedUrl] get presigned url error: %v", err)
	}
	return presignedURL.String(), nil
}

// Upload 上传文件
func (c *CosClient) Upload(ctx context.Context, cosKey string, reader io.Reader) error {
	_, err := c.client.Object.Put(ctx, cosKey, reader, nil)
	if err != nil {
		return fmt.Errorf("[Upload] upload reader error: %v", err)
	}
	return nil
}

// Get 下载文件
func (c *CosClient) Get(ctx context.Context, cosKey string) ([]byte, error) {
	rsp, err := c.client.Object.Get(ctx, cosKey, nil)
	if err != nil {
		return nil, fmt.Errorf("[Get] get file error: %v", err)
	}
	defer rsp.Body.Close()

	byts, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, fmt.Errorf("[Get] read file error: %v", err)
	}
	return byts, nil
}