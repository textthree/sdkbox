package alioss

// android：https://help.aliyun.com/zh/oss/developer-reference/authorize-access?spm=a2c4g.11186623.0.0.34485e0fRLKejO
// sts：https://help.aliyun.com/zh/oss/developer-reference/use-temporary-access-credentials-provided-by-sts-to-access-oss?spm=a2c4g.11186623.0.0.6c783f8cb2tU2S#section-7tz-fgu-oji

import (
	"bytes"
	"errors"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/textthree/provider"
	"github.com/textthree/provider/config"
	"github.com/textthree/provider/core"
	"sort"
	"strings"
	"sync"
	"time"
)

type Service interface {
	init()
	UploadFromByteArrayToOss(string, []byte) error
	UploadFromLocalFile(objectkey, localFilePath string) (err error)
	ListDir(path string, maxRows int) []OssDirList
	ListObjects(path string, maxRows int) []string
}

type AliossService struct {
	Service
	c      core.Container
	lock   sync.Mutex
	cfgSvc config.Service
}

// 返回 bool 代表初始化是否成功
func (self *AliossService) initClient() (*oss.Bucket, bool) {
	cfg := self.cfgSvc.GetAliOss()
	// oss.Timeout(30, 120) 表示设置HTTP连接超时时间为 30 秒，HTTP读写超时时间为 120 秒。0 表示永不超时（不推荐使用）
	client, err := oss.New(cfg.Endpoint, cfg.Ak, cfg.Sk, oss.Timeout(30, 120))
	if err != nil {
		provider.Clog().Error("初始化 OSS 失败:", err)
		return nil, false
	}
	bucketName := cfg.Bucket
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		provider.Clog().Error("Bucket 不正确:", err)
		return nil, false
	}
	return bucket, true
}

// 通过字节数组上传，在代码中根据表单名称获取字节数组，示例：
/*
	file, handle, err := c.GetFile(name)
	defer file.Close()
	fileByteArray, err = ioutil.ReadAll(file)
*/
func (self *AliossService) UploadFromByteArrayToOss(objectName string, fileConent []byte) (err error) {
	if bucket, success := self.initClient(); success {
		err = bucket.PutObject(objectName, bytes.NewReader(fileConent))
		if err != nil {
			err = errors.New("从字节数组文件上传出错:" + err.Error())
		}
	}
	return
}

func (self *AliossService) UploadFromLocalFile(objectkey, localFilePath string) (err error) {
	if bucket, success := self.initClient(); success {
		err = bucket.PutObjectFromFile(objectkey, localFilePath)
		if err != nil {
			err = errors.New("本地文件上传出错:" + err.Error())
		}
	}
	return
}

// 列举指定路径下的目录，不递归，按修改时间倒序
func (self *AliossService) ListDir(path string, maxRows int) (list []OssDirList) {
	var bucket *oss.Bucket
	var initSuccess bool
	if bucket, initSuccess = self.initClient(); initSuccess {
		return
	}
	result, err := bucket.ListObjects(
		oss.Prefix(path),
		oss.Delimiter("/"),   // 结尾带 / 的代表是目录
		oss.MaxKeys(maxRows), // 最多列举 maxRows 个
	)
	if err != nil {
		provider.Clog().Error(err.Error())
		return
	}
	for _, commonPrefix := range result.CommonPrefixes {
		item := OssDirList{
			Path:    commonPrefix,
			DirName: strings.Split(commonPrefix, "/")[1],
		}
		list = append(list, item)
	}
	// 按时间排序
	for index, item := range list {
		result, err = bucket.ListObjects(
			oss.Prefix(item.Path),
			oss.MaxKeys(100), // 最多列举 100 个
		)
		var lastModifiyTIme time.Time
		if len(result.Objects) > 0 {
			prev := result.Objects[0].LastModified
			for _, v := range result.Objects {
				if v.LastModified.After(prev) {
					prev = v.LastModified
				}
			}
			lastModifiyTIme = prev
		}
		list[index].LastModifyTime = lastModifiyTIme
	}
	sort.Slice(list, func(i, j int) bool {
		return list[i].LastModifyTime.After(list[j].LastModifyTime)
	})
	return
}

// 列举指定路径下的所有文件的路径，字典序
func (self *AliossService) ListObjects(path string, maxRows int) (list []string) {
	var bucket *oss.Bucket
	var initSuccess bool
	if bucket, initSuccess = self.initClient(); initSuccess {
		return
	}
	result, err := bucket.ListObjects(
		oss.Prefix(path),
		oss.MaxKeys(maxRows), // 最多列举 maxRows 个
	)
	if err != nil {
		provider.Clog().Error(err.Error())
		return
	}
	for _, object := range result.Objects {
		if object.Size > 0 {
			url := object.Key
			list = append(list, url)
		}
	}
	return
}
