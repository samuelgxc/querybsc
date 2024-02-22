package oss

import (
	"github.com/minio/minio-go"
	"shangchenggo/common/config"
	"shangchenggo/tools"
	"mime/multipart"
)

var (
	localOssClient *minio.Client //局域网oss
	localOssDomain string        //局域网oss domain
	localOssBucket string        //局域网oss 存储桶
)

//连接局域网oss
func StartMinio() error {
	localOss := config.Val.Local_oss()
	//localOss := `{"scheme":"http","endpoint":"192.168.0.170:9000","accesskey":"AKIAIOSFODNN7EXAMPLE","secretkey":"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY","bucket":"mcc1200"}`
	ossConfig := tools.JsonDecode(localOss)
	endpoint := tools.ToString(ossConfig["endpoint"])
	accesskey := tools.ToString(ossConfig["accesskey"])
	secretkey := tools.ToString(ossConfig["secretkey"])
	bucket := tools.ToString(ossConfig["bucket"])
	scheme := tools.ToString(ossConfig["scheme"])

	secure := false
	domian := "http://" + endpoint
	if scheme == "https" {
		secure = true
		domian = "https://" + endpoint
	}
	minioClient, err := minio.New(endpoint, accesskey, secretkey, secure)
	if err != nil {
		return err
	}
	ok, err := minioClient.BucketExists(bucket)
	if err != nil {
		return err
	} else if !ok {
		err = minioClient.MakeBucket(bucket, "")
		if err != nil {
			return err
		}
	}

	localOssClient = minioClient
	localOssDomain = domian
	localOssBucket = bucket
	return nil
}

func NewMinioClient() *minio.Client {
	return localOssClient
}

func GetMinioDomain() string {
	return localOssDomain
}

func GetMinioBucket() string {
	return localOssBucket
}

//保存到minio oss
func SaveMinioObject(header *multipart.FileHeader, filepath string) (url string, err error) {
	url = GetMinioDomain() + "/" + GetMinioBucket() + "/" + filepath
	file, err := header.Open()
	if err != nil {
		return url, err
	}
	defer file.Close()
	_, err = NewMinioClient().PutObject(GetMinioBucket(), filepath, file, header.Size, minio.PutObjectOptions{ContentType: GetContentType(header.Filename)})
	if err != nil {
		return url, err
	}
	return url, nil
}
