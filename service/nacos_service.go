package service

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/EricJSanchez/gotool/sys"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/spf13/cast"
	"github.com/spf13/viper"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
)

var onceNacos sync.Once

type Nacos struct {
	NamingClient naming_client.INamingClient
	ConfigClient map[string]config_client.IConfigClient
}

// InitClient 初始化
func (n *Nacos) InitClient() error {
	onceNacos.Do(func() {
		clientConfig := *constant.NewClientConfig(
			constant.WithNamespaceId(sys.Cfg("app").GetString("nacos.namespace_id")),
			constant.WithTimeoutMs(uint64(sys.Cfg("app").GetInt("nacos.timeout"))),
			constant.WithNotLoadCacheAtStart(true),
			constant.WithLogDir(sys.Cfg("app").GetString("log_path")+"nacos/cache"),
			constant.WithCacheDir(sys.Cfg("app").GetString("log_path")+"nacos/logs"),
			//constant.WithUsername(sys.Cfg("app").GetString("nacos.username")),
			//constant.WithPassword(sys.Cfg("app").GetString("nacos.password")),
		)
		serverConfigs := []constant.ServerConfig{
			{
				IpAddr:      sys.Cfg("app").GetString("nacos.addr"),
				ContextPath: "/nacos",
				Port:        uint64(sys.Cfg("app").GetInt("nacos.port")),
				Scheme:      sys.Cfg("app").GetString("nacos.scheme"),
			},
		}
		var err error
		if n.NamingClient, err = clients.NewNamingClient(
			vo.NacosClientParam{
				ClientConfig:  &clientConfig,
				ServerConfigs: serverConfigs,
			},
		); err != nil {
			panic("init nacos error:" + cast.ToString(err))
			//sys.Log().WithError(err).Error("init nacos error")
			//return
		}

		// 下面开始初始化配置文件监听，根据data_id和group
		ns := sys.Cfg("app").GetStringSlice("nacos.group_data_ids")
		sys.NacosConfig = make(map[string]string, len(ns))
		n.ConfigClient = make(map[string]config_client.IConfigClient, len(ns))
		for _, nv := range ns {
			tmpConf := strings.Split(nv, ":")
			if len(tmpConf) != 2 {
				fmt.Println("配置有误：" + nv)
				continue
			}
			// 创建动态配置客户端
			if n.ConfigClient[tmpConf[0]], err = clients.NewConfigClient(
				vo.NacosClientParam{
					ClientConfig:  &clientConfig,
					ServerConfigs: serverConfigs,
				},
			); err != nil {
				fmt.Println("init nacos error:" + cast.ToString(err))
				return
			}
			n.NewInitConfig(tmpConf[0], tmpConf[1])
		}
	})
	return nil
}

// NewInitConfig 初始化配置文件监听
func (n *Nacos) NewInitConfig(dataId, group string) {
	var err error
	if sys.NacosConfig[dataId], err = n.ConfigClient[dataId].GetConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	}); err != nil {
		fmt.Println("NewInitConfig error:" + cast.ToString(err))
		panic(err)
	}
	//helper.Pr(dataId, sys.NacosConfig)
	err = n.ConfigClient[dataId].ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			//nacos变更,更新本地
			fmt.Println(dataId+" nacos changed:", data)
			sys.ResetCfgKey(dataId)
			sys.NacosConfig[dataId] = data
		},
	})
	if err != nil {
		fmt.Println(err)
	}
}

func GetLocalIp() (string, error) {
	addRs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addRs {
		// 检查ip地址判断是否回环地址
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), err
			}
		}
	}
	return "", errors.New("no ip found")
}

// Register 服务注册
func (n *Nacos) Register() {
	var (
		ip  string
		err error
	)
	if ip, err = GetLocalIp(); err != nil {
		fmt.Println("get local ip error", err)
		return
	}
	if _, err = n.NamingClient.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          ip,
		Port:        uint64(sys.Cfg("app").GetInt("port")),
		ServiceName: sys.Cfg("app").GetString("service_name"),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		//Metadata:    map[string]string{"idc": "shanghai"},
	}); err != nil {
		return
	}
}

// DeRegister 服务撤销
func (n *Nacos) DeRegister() {
	var (
		ip  string
		err error
	)
	if ip, err = GetLocalIp(); err != nil {
		return
	}
	if _, err = n.NamingClient.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          ip,
		Port:        uint64(sys.Cfg("app").GetInt("port")),
		ServiceName: sys.Cfg("app").GetString("service_name"),
		Ephemeral:   true,
	}); err != nil {
		return
	}
}

// GetUri 根据权重获取某服务的健康节点
func (n Nacos) GetUri(serviceName string, router string) (string, error) {
	if instance, err := n.NamingClient.SelectOneHealthyInstance(
		vo.SelectOneHealthInstanceParam{
			ServiceName: serviceName,
		}); err != nil {
		return "", err
	} else {
		return "http://" + instance.Ip + ":" + strconv.FormatUint(instance.Port, 10) + router, nil
	}
}

func (n Nacos) GetViper(data string, ty string) (vp *viper.Viper, err error) {
	vp = viper.New()
	vp.SetConfigType(ty)
	err = vp.ReadConfig(bytes.NewBuffer([]byte(data)))
	if err != nil {
		log.Println("GetViper err:", err)
		return
	}
	return
}

func (n Nacos) GetTomlViper(data string) (vp *viper.Viper, err error) {
	return n.GetViper(data, "toml")
}
