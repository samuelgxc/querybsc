package oss

import (
	"fmt"
	"github.com/minio/minio-go"
	"os"
	"testing"
)

func TestMinioOss(t *testing.T) {
	StartMinio()

	filePath := `D:\00Down\cache\IMG_20171030_231638.jpg`
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileInfo, _ := file.Stat()

	filename := "test" + "/" + fileInfo.Name()
	_, err = NewMinioClient().PutObject(GetMinioBucket(), filename, file, fileInfo.Size(), minio.PutObjectOptions{ContentType: GetContentType(filename)})
	if err != nil {
		fmt.Println(err)
	} else {
		url := GetMinioDomain() + "/" + GetMinioBucket() + "/" + filename
		fmt.Println(url)
	}
}
