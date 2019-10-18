package ossutil

import "github.com/aliyun/aliyun-oss-go-sdk/oss"

type Oss struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	EndPoint        string `mapstructure:"end_point"`
	OuterEndPoint   string `mapstructure:"outer_end_point"`
	BucketName      string `mapstructure:"bucket_name"`
	ImageDomain     string `mapstructure:"image_domain"`
}

func (o *Oss) Upload(srcName, dstName string, isInternal bool) error {
	endPoint := o.EndPoint
	if !isInternal {
		endPoint = o.OuterEndPoint
	}
	client, err := oss.New(endPoint, o.AccessKeyID, o.AccessKeySecret)
	if err != nil {
		return err
	}
	// 获取存储空间。
	bucket, err := client.Bucket(o.BucketName)
	if err != nil {
		return err
	}
	// 上传文件。
	err = bucket.PutObjectFromFile(dstName, srcName)
	if err != nil {
		return err
	}
	return nil
}
