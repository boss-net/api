package factory

import (
	"fmt"

	"github.com/boss-net/api/cloud/oss"
	"github.com/boss-net/api/cloud/oss/aliyun"
	"github.com/boss-net/api/cloud/oss/azureblob"
	"github.com/boss-net/api/cloud/oss/gcsblob"
	"github.com/boss-net/api/cloud/oss/huanweiobs"
	"github.com/boss-net/api/cloud/oss/local"
	"github.com/boss-net/api/cloud/oss/s3"
	"github.com/boss-net/api/cloud/oss/tencentcos"
	"github.com/boss-net/api/cloud/oss/volcenginetos"
)

var OSSFactory = map[string]func(oss.OSSArgs) (oss.OSS, error){
	"local":      local.NewLocalStorage,
	"local_file": local.NewLocalStorage,

	"s3":     s3.NewS3Storage,
	"aws_s3": s3.NewS3Storage,
	"aws-s3": s3.NewS3Storage,

	"azure":      azureblob.NewAzureBlobStorage,
	"azure_blob": azureblob.NewAzureBlobStorage,
	"azure-blob": azureblob.NewAzureBlobStorage,

	"aliyun":     aliyun.NewAliyunOSSStorage,
	"aliyun-oss": aliyun.NewAliyunOSSStorage,
	"aliyun_oss": aliyun.NewAliyunOSSStorage,

	"tencent":     tencentcos.NewTencentCOSStorage,
	"tencent_cos": tencentcos.NewTencentCOSStorage,
	"tencent-cos": tencentcos.NewTencentCOSStorage,

	"gcs":            gcsblob.NewGoogleCloudStorage,
	"google-storage": gcsblob.NewGoogleCloudStorage,
	"google_storage": gcsblob.NewGoogleCloudStorage,

	"huawei":     huanweiobs.NewHuaweiOBSStorage,
	"huawei-obs": huanweiobs.NewHuaweiOBSStorage,
	"huawei_obs": huanweiobs.NewHuaweiOBSStorage,

	"volcengine":     volcenginetos.NewVolcengineTOSStorage,
	"volcengine_tos": volcenginetos.NewVolcengineTOSStorage,
	"volcengine-tos": volcenginetos.NewVolcengineTOSStorage,
}

func Load(name string, args oss.OSSArgs) (oss.OSS, error) {
	f, ok := OSSFactory[name]
	if !ok {
		msg := fmt.Sprintf("[ %s ] is not in the provider list", name)
		return nil, oss.ErrProviderNotFound.WithDetail(msg)
	}
	return f(args)
}
