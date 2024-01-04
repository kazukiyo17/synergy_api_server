package scene_service

import (
	"github.com/kazukiyo17/synergy_api_server/model/scene"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"github.com/kazukiyo17/synergy_api_server/redis_mq"
	"log"
	"strings"
)

const (
	SCENE_EXPIRE_TIME = 1
)

// hasGenerated 是否已经生成
func hasGenerated(sceneId string) (bool, error) {
	// 1 查询redis
	rKey := "cos:" + sceneId
	if redis.Exists(rKey) {
		return true, nil
	}
	// 2. 查redis_mq
	if redis_mq.Check(sceneId) {
		return true, nil
	}
	// 3. 查询db
	cosUrl, err := scene.GetCosUrlBySceneId(sceneId)
	if err != nil {
		return false, err
	}
	if cosUrl == "" {
		return false, nil
	}
	// 2. 存入redis
	err = redis.Set(rKey, cosUrl, SCENE_EXPIRE_TIME)
	if err != nil {
		log.Printf("redis set error: %v", err)
	}
	return true, err
}

// produceChildScene 加入队列
func produceChildScene(sceneId, childSceneId, username string) {
	redis_mq.Produce(childSceneId, username)
}

// getChildSceneIds 获取子场景Id
func getChildSceneIds(sceneId string) ([]string, error) {
	childSceneIds := make([]string, 0)
	// 从redis中获取子场景
	rKey := "child:" + sceneId
	if redis.Exists(rKey) {
		childSceneIdsStr, err := redis.Get(rKey)
		if err == nil {
			childSceneIds = strings.Split(childSceneIdsStr, ",")
		}
	} else {
		sceneIds, err := scene.GetSceneIdByParentSceneId(sceneId)
		if err != nil {
			return make([]string, 0), err
		}
		childSceneIds = append(sceneIds, childSceneIds...)
		// 保存到redis
		childSceneIdsStr := strings.Join(childSceneIds, ",")
		err = redis.Set(rKey, childSceneIdsStr, SCENE_EXPIRE_TIME)
	}
	return childSceneIds, nil
}

func getSceneCreator(sceneId string) (string, error) {
	rKey := "creator:" + sceneId
	if redis.Exists(rKey) {
		return redis.Get(rKey)
	} else {
		creator, err := scene.GetCreatorBySceneId(sceneId)
		if err != nil {
			return "", err
		}
		err = redis.Set(rKey, creator, SCENE_EXPIRE_TIME)
		return creator, err
	}
}

func CheckSceneCreator(sceneId, username string) bool {
	creator, err := getSceneCreator(sceneId)
	if err != nil || creator != username {
		return false
	}
	return true
}

// Check 生成孙子剧本
func Check(sceneId, username string) error {
	// 1. 获取子场景
	childSceneIds, err := getChildSceneIds(sceneId)
	if err != nil {
		return err
	}
	for _, childSceneId := range childSceneIds {
		grandChildSceneIds, err := getChildSceneIds(childSceneId)
		if err != nil {
			continue
		}
		for _, grandChildSceneId := range grandChildSceneIds {
			generated, err := hasGenerated(grandChildSceneId)
			if err != nil {
				log.Fatalln(err)
			}
			if generated {
				continue
			}
			produceChildScene(childSceneId, grandChildSceneId, username)
		}
	}
	return nil
}
