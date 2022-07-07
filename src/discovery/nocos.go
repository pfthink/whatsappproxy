package discovery

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"strconv"
	"whatsappproxy/config"
	"whatsappproxy/utils"
)

func InitNacos() {
	sc := []constant.ServerConfig{
		{
			Scheme: "https",
			IpAddr: "",
			Port:   443,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId:         config.NamespaceId,
		TimeoutMs:           5000,
		NotLoadCacheAtStart: true,
		LogDir:              "/var/data/logs",
		CacheDir:            "/var/data/cache",
		LogRollingConfig:    &constant.ClientLogRollingConfig{MaxSize: 10},
		LogLevel:            "info",
	}

	client, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &cc,
			ServerConfigs: sc,
		},
	)

	if err != nil {
		panic(err)
	}

	//Register with default cluster and group
	//ClusterName=DEFAULT,GroupName=DEFAULT_GROUP
	netIp, _ := utils.GetLocalIP()
	ip := netIp.String()
	port, err := strconv.ParseUint(config.AppPort, 10, 64)
	client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        port,
		ServiceName: config.ServiceName,
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    map[string]string{},
	})

}
