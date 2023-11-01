package scene

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
			json.Unmarshal(data, &cacheScene)
			return cacheScene, nil
		}
	}
	// 缓存中没有
	scene, err := model.GetSceneById(sceneId)
	if err != nil {
		return nil, err
	}
	// 存入缓存
	data, _ := json.Marshal(scene)
	redis.Set(sceneIdStr, data)
	return scene, nil
}
