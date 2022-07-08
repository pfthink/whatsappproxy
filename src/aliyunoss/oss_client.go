package aliyunoss

import (
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"sync"
	"time"
	"whatsappproxy/config"
	"whatsappproxy/utils"
)

var (
	syOnce sync.Once
)

func InitOssClient() (Bucket *oss.Bucket) {
	if Bucket == nil {
		syOnce.Do(func() {
			client, err := oss.New(config.OssEndpoint, config.AccessKeyId, config.AccessKeySecret)
			if err != nil {
				logger.Errorf("Init aliyunOss client fail,error:%s", err)
			}
			Bucket, err = client.Bucket(config.BucketName)
			if err != nil {
				logger.Errorf("Get aliyunOss bucket fail,error:%s", err)
			}
		})
		logger.Info("阿里云OSS初始化成功")
	}
	return Bucket
}

func UploadByUrl(imageUrl string) (cdn, bucketPath, fileKey string, err error) {
	res, err := http.Get(imageUrl)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	bucketUrl := fmt.Sprintf("https://%s.%s/", config.BucketName, config.OssEndpoint)
	l, err := url.Parse(imageUrl)
	fileSuffix := path.Ext(l.Path)
	path := config.OssPath + generateFileName() + fileSuffix
	err = utils.Bucket.PutObject(path, res.Body)
	if err != nil {
		logger.Errorf("UploadByUrl Error:%s", err)
	}
	err = utils.Bucket.SetObjectACL(path, oss.ACLPublicRead)
	if err != nil {
		logger.Errorf("Set fileAcl Error:%s", err)
	}
	return config.OssCDN, bucketUrl, path, err
}

func generateFileName() string {
	t := time.Now()
	return fmt.Sprintf("%s%s%s%s%s%s%s-%s", strconv.Itoa(t.Year()),
		strconv.Itoa(int(t.Month())), strconv.Itoa(t.Day()),
		strconv.Itoa(t.Hour()), strconv.Itoa(t.Minute()),
		strconv.Itoa(t.Second()), strconv.Itoa(int(t.Unix())),
		utils.RandomStr(16))
}
