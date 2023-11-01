package conf

import (
	"context"
	redis "github.com/go-redis/redis/v8"
	"github.com/jinzhu/gorm"
	"sync"
)

var (
	once sync.Once
	C    = new(Config) // C 全局配置
)

type Config struct {
	//Server       *serverConfig                    // 服务配置
	MysqlConnMap map[string]*gorm.DB              // 数据库实例
	RedisConnMap map[string]redis.UniversalClient // redis实例
	//ESConnMap    map[string]*elasticv6.Client     // es实例
	//consumerAPI  api.ConsumerAPI                  // 北极星客户端，只在conf中使用，不可外部使用
	//secretProxy  ywSecret.SecretClientProxy       // 密钥客户端
}

// Load 加载服务配置,加载失败直接退出
func Load(ctx context.Context) {
	once.Do(func() {
		cfg := new(Config)
		//cfg.secretProxy = ywSecret.NewSecretClientProxy()
		//logger := pcgLog.Start(ctx)

		//srvCfg := new(serverConfig)
		//err := readConfig(ctx, rainbowServiceGroupKey, serverConfigKey, srvCfg)
		//if err != nil {
		//	logger.Errorf(-1, "failed to load %s, err: %v", serverConfigKey, err)
		//	panic("failed to read serverConfig")
		//}
		//cfg.Server = srvCfg
		//logger.Info("read serverConfig success")

		// 作为主调端使用，直接创建ConsumerAPI
		//consumer, err := api.NewConsumerAPI()
		//if err != nil {
		//	logger.Errorf(errcode.L5FailError, "NewConsumerAPI err: %v", err)
		//	panic("failed to NewConsumerAPI")
		//}
		//cfg.consumerAPI = consumer
		//defer cfg.consumerAPI.Destroy()

		err = cfg.loadMysql(ctx)
		if err != nil {
			logger.Errorf(errcode.DBFailError, "loadMysql err: %v", err)
			panic("failed to loadMysql")
		}

		err = cfg.loadRedis(ctx)
		if err != nil {
			logger.Errorf(errcode.RedisFailError, "loadRedis err: %v", err)
			panic("failed to loadRedis")
		}

		err = cfg.loadES(ctx)
		if err != nil {
			logger.Errorf(errcode.ESFailError, "loadES err: %v", err)
			panic("failed to loadES")
		}

		C = cfg
	})
}

func (c *Config) loadMysql(ctx context.Context) error {
	mysqlConnMap := make(map[string]*gorm.DB, len(dbConfKeys))
	for _, rainbowKey := range dbConfKeys {
		options, err := c.getMysqlOption(ctx, rainbowKey)
		if err != nil {
			//pcgLog.Errorf(ctx, errcode.DBFailError, "[NewMysqlClientProxy]rainbowKey:%v, err: %v",
			//	rainbowKey, err)
			return err
		}

		//trpcName := "trpc.ywos.mysql." + rainbowKey
		// db *gorm.DB
		dbClient, err := gorm.Open("mysql", options...)
		if err != nil {
			//pcgLog.Errorf(ctx, errcode.DBFailError, "[NewMysqlClientProxy]rainbowKey:%v, err: %v", rainbowKey, err)
			// retry
			dbClient, err = pcgGorm.NewClientProxy(trpcName, options...)
			if err != nil {
				pcgLog.Fatalf(ctx, "[NewMysqlClientProxy]rainbowKey:%s, err: %v", rainbowKey, err)
				return err
			}
		}
		mysqlConnMap[rainbowKey] = dbClient
	}
	c.MysqlConnMap = mysqlConnMap
	pcgLog.Info(ctx, "mysql config load finish")

	return nil
}
