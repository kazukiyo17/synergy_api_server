package scene

import (
	"encoding/json"
	e "github.com/kazukiyo17/synergy_api_server/common/errcode"
	model "github.com/kazukiyo17/synergy_api_server/model/scene"
	"github.com/kazukiyo17/synergy_api_server/redis"
	"github.com/kazukiyo17/synergy_api_server/redis_mq"
	"github.com/kazukiyo17/synergy_api_server/setting"
	"github.com/kazukiyo17/synergy_api_server/utils/flake"
	"github.com/spf13/cast"
	"log"
	"strconv"
)

type Scene struct {
	Url      string `json:"url"`
	SceneId  string `json:"sceneId"`
	Username string `json:"username"`
	//IsInit   int    `json:"is_init"`
}

// hasGenerated 是否已经生成
//func hasGenerated(sceneId string) (bool, error) {
//	// 1 查询redis
//	rKey := "cos:" + sceneId
//	if redis.Exists(rKey) {
//		return true, nil
//	}
//	// 2. 查redis_mq
//	if redis_mq.Check(sceneId) {
//		return true, nil
//	}
//	// 3. 查询db
//	cosUrl, err := model.GetCosUrlBySceneId(sceneId)
//	if err != nil {
//		return false, err
//	}
//	if cosUrl == "" {
//		return false, nil
//	}
//	// 2. 存入redis
//	redis.Set(rKey, cosUrl, setting.ServerSetting.SceneExpire)
//	return true, err
//}

//// produceChildScene 加入队列
//func produceChildScene(childSceneId, username string) {
//	redis_mq.Produce(childSceneId, username)
//}

// getChildSceneIds 获取子场景Id
//func getChildSceneIds(sceneId string) ([]string, error) {
//	childSceneIds := make([]string, 0)
//	// 从redis中获取子场景
//	rKey := "child:" + sceneId
//	if redis.Exists(rKey) {
//		childSceneIdsStr := redis.Get(rKey)
//		childSceneIds = strings.Split(childSceneIdsStr, ",")
//	} else {
//		sceneIds, err := model.GetSceneIdByParentSceneId(sceneId)
//		if err != nil {
//			return make([]string, 0), err
//		}
//		childSceneIds = append(sceneIds, childSceneIds...)
//		// 保存到redis
//		childSceneIdsStr := strings.Join(childSceneIds, ",")
//		redis.Set(rKey, childSceneIdsStr, setting.ServerSetting.SceneExpire)
//	}
//	return childSceneIds, nil
//}

func getChildScenes(sceneId string) ([]*model.Scene, error) {
	childScenes := make([]*model.Scene, 0)
	// 从redis中获取子场景
	rKey := "childs:" + sceneId
	if redis.Exists(rKey) {
		jsonStr := redis.Get(rKey)
		err := json.Unmarshal([]byte(jsonStr), &childScenes)
		if err == nil {
			return childScenes, nil
		}
	}
	childScenes, err := model.GetSceneByParentSceneId(sceneId)
	if err != nil {
		return make([]*model.Scene, 0), err
	}
	jsonStr, err := json.Marshal(childScenes)
	if err != nil {
		redis.Set(rKey, string(jsonStr), setting.ServerSetting.SceneExpire)
	}
	return childScenes, nil
}

//func getSceneCreator(sceneId string) (string, error) {
//	sceneInfo, err := scene.GetSceneBySceneId(sceneId)
//	if err != nil {
//		return "", err
//	}
//	return sceneInfo.Creator, nil
//}

func GetSceneInfo(sceneId, username string) (*Scene, error) {
	var s = &Scene{}
	rKey := "scene:" + sceneId
	if redis.Exists(rKey) {
		sceneInfo := redis.Get(rKey)
		err := json.Unmarshal([]byte(sceneInfo), &s)
		if err == nil {
			return s, err
		}
	}
	sceneInfo, err := model.GetSceneByCreatorAndSceneId(username, sceneId)
	if err != nil {
		return s, err
	}
	s.SceneId = strconv.FormatInt(sceneInfo.SceneId, 10)
	s.Username = sceneInfo.Creator
	s.Url = sceneInfo.COSUrl
	//s.IsInit = sceneInfo.IsInit
	// 转成json
	sceneJson, err := json.Marshal(s)
	if err == nil {
		redis.Set(rKey, string(sceneJson), setting.ServerSetting.SceneExpire)
	}
	return s, nil
}

// Check 生成孙子剧本
func Check(sceneId, username string) (int, *Scene) {
	// sceneId 是否为数字
	_, err := cast.ToInt64E(sceneId)
	if err != nil {
		return e.INVALID_PARAMS, nil
	}
	// 用户是否有权限
	sceneInfo, err := GetSceneInfo(sceneId, username)
	if err != nil {
		return e.ERROR, nil
	}
	if sceneInfo.Username != username {
		return e.AUTH_CHECK_ERROR, nil
	}
	// 获取子场景,
	childScenes, err := getChildScenes(sceneId)
	if err != nil {
		return e.ERROR,	nil
	}
	// 生成子场景
	for _, childScene := range childScenes {
		if childScene.COSUrl != "" {
			continue
		}
		jsonStr, err := json.Marshal(childScene)
		if err != nil {
			log.Println(err)
			continue
		}
		redis_mq.Produce(strconv.FormatInt(childScene.SceneId, 10), string(jsonStr))
		//produceChildScene(strconv.FormatInt(childScene.SceneId, 10), string(jsonStr))
	}
	return e.SUCCESS, sceneInfo
}

func GenerateInitScene(username string) {
	chooses := make([]string, 0)
	chooses = append(chooses, "我决定独自前往，无论前方有多少困难")
	chooses = append(chooses, "尝试探探周远山的口风")
	for index, choose := range chooses {
		sceneId, err := flake.Generate()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		//sceneId := int64(index + 1)
		err, scene := model.SaveUngeneratedScene(sceneId, int64(0), choose, username, index+1)
		if err != nil {
			log.Fatalln(err)
			continue
		}
		// scene 转 json
		sceneJson, err := json.Marshal(scene)
		if err != nil {
			log.Fatalln(err)
			continue
		}
		redis_mq.Produce(strconv.FormatInt(sceneId, 10), string(sceneJson))
	}
}

//func GetInitChooseScene(username, choose string) (int, *StartScene) {
//	rKey := "init:" + username + ":" + choose
//	var initScene = StartScene{}
//	if redis.Exists(rKey) {
//		jsonStr := redis.Get(rKey)
//		err := json.Unmarshal([]byte(jsonStr), &initScene)
//		if err != nil {
//			return e.ERROR, &initScene
//		}
//		return e.SUCCESS, &initScene
//	}
//	sceneId, cosUrl, err := model.GetCosUrlByCreatorAndChooseContent(username, choose)
//	if err != nil {
//		return e.ERROR, &initScene
//	}
//	initScene.SceneId = sceneId
//	initScene.Url = cosUrl
//	initScene.Username = username
//	jsonStr, err := json.Marshal(initScene)
//	if err != nil {
//		return e.ERROR, &initScene
//	}
//	redis.Set(rKey, string(jsonStr), setting.ServerSetting.SceneExpire)
//	return e.SUCCESS, &initScene
//
//}


func GetInitScene(username, initId string) (int, *Scene){
	rKey := "init:" + username + initId
	var initScene = Scene{}
	if redis.Exists(rKey) {
		jsonStr := redis.Get(rKey)
		err := json.Unmarshal([]byte(jsonStr), &initScene)
		if err == nil {
			return e.ERROR, &initScene
		}
	}
	scene , err := model.GetSceneByCreatorAndInitId(username, cast.ToInt(initId))
	if err != nil {
		return e.ERROR, &initScene
	}
	sceneId := strconv.FormatInt(scene.SceneId, 10)
	//initScene.SceneId = strconv.FormatInt(scene.SceneId, 10)
	//initScene.Url = scene.COSUrl
	//initScene.Username = username
	jsonStr, err := json.Marshal(initScene)
	if err == nil {
		redis.Set(rKey, string(jsonStr), setting.ServerSetting.SceneExpire)
		redis.Set("scene:" + sceneId, string(jsonStr), setting.ServerSetting.SceneExpire)
	}
	// 生成
	code, sceneInfo := Check(sceneId, username)
	return code, sceneInfo
}