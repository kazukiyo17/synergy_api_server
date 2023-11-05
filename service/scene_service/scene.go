package scene_service

import (
	"encoding/json"
	"github.com/kazukiyo17/fake_buddha_server/model"
	"github.com/kazukiyo17/fake_buddha_server/utils/redis"
	"github.com/spf13/cast"
)

type Scene struct {
	ID            int64  `json:"id"`
	Url           string `json:"url"`
	CreatorUserId int64  `json:"creator_user_id"`
}

func (s *Scene) GetSceneInfoById(sceneId int64) (*model.Scene, error) {
	var cacheScene *model.Scene
	sceneIdStr := cast.ToString(sceneId)
	if redis.Exists(sceneIdStr) {
		data, err := redis.Get(sceneIdStr)
		if err != nil {
			//logging.Info(err)
		} else {
			err = json.Unmarshal(data, &cacheScene)
			if err != nil {
				return nil, err
			}
			return cacheScene, nil
		}
	}
	// 缓存中没有,从数据库中获取
	scene, err := model.GetSceneById(sceneId)
	if err != nil {
		return nil, err
	}
	// 存入缓存
	data, _ := json.Marshal(scene)
	err = redis.Set(sceneIdStr, data, 3600)
	if err != nil {
		return nil, err
	}
	return scene, nil
}

// CheckSubScene 检查子场景是否都存在
func (s *Scene) CheckSubScene(sceneId int64) (bool, error) {
	scene, err := s.GetSceneInfoById(sceneId)
	if err != nil {
		return false, err
	}
	if scene == nil {
		return false, nil
	}
	return true, nil
}
