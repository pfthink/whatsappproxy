package config

import (
	"github.com/nacos-group/nacos-sdk-go/common/logger"
	"github.com/pfthink/agollo"
	"github.com/pfthink/agollo/env/config"
)

func InitApollo() agollo.Client {
	// 或者忽略错误处理直接 a.Start()
	c := &config.AppConfig{
		AppID:          "",
		Cluster:        "",
		IP:             "",
		NamespaceName:  "",
		IsBackupConfig: true,
		Secret:         "",
		AuthToken:      "",
	}

	agollo.SetLogger(logger.GetLogger())

	client, err := agollo.StartWithConfig(func() (*config.AppConfig, error) {
		return c, nil
	})
	if err != nil {
		logger.Errorf("init apollo error :%s ", err)
	}
	//utils.ApolloClient = client
	logger.Infof("初始化Apollo配置成功")
	return client
}
