package tool

import (
	"agent_group/logic/constant"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// https://cloud.tencent.com/document/product/436/35059

func CosUploadPresignedUrl(ctx context.Context, cosKey, cosUrl, ak, sk string, expiration time.Duration) (presignedUrl string, err error) {
	u, err := url.Parse(cosUrl)
	if err != nil {
		err = fmt.Errorf("[PresignedUrl]  parse cos url error: %v", err)
		return
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
	})

	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodPut, cosKey, ak, sk, expiration, nil)
	if err != nil {
		err = fmt.Errorf("[PresignedUrl]  get presigned url error: %v", err)
		return
	}
	presignedUrl = presignedURL.String()

	return
}

func CosDownloadPresignedUrl(ctx context.Context, cosKey, cosUrl, ak, sk string, expiration time.Duration) (presignedUrl string, err error) {
	u, err := url.Parse(cosUrl)
	if err != nil {
		err = fmt.Errorf("[PresignedUrl]  parse cos url error: %v", err)
		return
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
	})

	presignedURL, err := client.Object.GetPresignedURL(ctx, http.MethodGet, cosKey, ak, sk, expiration, nil)
	if err != nil {
		err = fmt.Errorf("[PresignedUrl]  get presigned url error: %v", err)
		return
	}
	presignedUrl = presignedURL.String()

	return
}

func CosUpload(ctx context.Context, cosKey, cosUrl, ak, sk string, reader io.Reader) (err error) {
	u, err := url.Parse(cosUrl)
	if err != nil {
		err = fmt.Errorf("[UploadFile]  parse cos url error: %v", err)
		return
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
	})

	_, err = client.Object.Put(ctx, cosKey, reader, nil)
	if err != nil {
		err = fmt.Errorf("[UploadFile]  upload reader error: %v", err)
		return
	}

	return
}

func CosGet(ctx context.Context, cosKey, cosUrl, ak, sk string) (byts [] byte, err error) {
	u, err := url.Parse(cosUrl)
	if err != nil {
		err = fmt.Errorf("[GetFile]  parse cos url error: %v", err)
		return
	}
	b := &cos.BaseURL{BucketURL: u}
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  ak,
			SecretKey: sk,
		},
	})

	rsp, err := client.Object.Get(ctx, cosKey, nil)
	if err != nil {
		err = fmt.Errorf("[GetFile]  get file error: %v", err)
		return
	}
	defer rsp.Body.Close()
	byts, err = io.ReadAll(rsp.Body)
	if err != nil {
		err = fmt.Errorf("[GetFile]  read file error: %v", err)
		return
	}
	return
}
