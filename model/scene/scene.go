package scene

import (
	"github.com/kazukiyo17/synergy_api_server/model"
	"strconv"
)

type Scene struct {
	model.Model
	SceneId       int64  `json:"scene_id" gorm:"type:bigint;index"`
	ChooseContent string `json:"choose_content" gorm:"type:varchar(255)"`
	Creator       string `json:"creator" gorm:"type:varchar(255)"`
	ParentSceneId int64  `json:"parent_scene_id" gorm:"type:bigint;index"`
	COSUrl        string `json:"cos_url" gorm:"type:varchar(255)"`
	ShortDesc     string `json:"desc" gorm:"type:varchar(600)"`
	IsInit        int    `json:"is_init" gorm:"type:int(11)"`
}

func GetCosUrlBySceneId(sceneId string) (cosUrl string, err error) {
	scene := &Scene{}
	err = model.DB.Model(&Scene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	cosUrl = scene.COSUrl
	return cosUrl, nil
}

func GetSceneByParentSceneId(parentSceneId string) ([]*Scene, error) {
	scenes := make([]*Scene, 0)
	err := model.DB.Model(&Scene{}).Where("parent_scene_id = ?", parentSceneId).Find(&scenes).Error
	if err != nil {
		return nil, err
	}
	return scenes, nil
}

func GetCreatorBySceneId(sceneId string) (creator string, err error) {
	scene := &Scene{}
	err = model.DB.Model(&Scene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	creator = scene.Creator
	return creator, nil
}

func GetSceneBySceneId(sceneId string) (scene *Scene, err error) {
	scene = &Scene{}
	err = model.DB.Model(&Scene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return nil, err
	}
	return scene, nil
}

func GetSceneByCreatorAndSceneId(creator, sceneId string) (scene *Scene, err error) {
	scene = &Scene{}
	err = model.DB.Model(&Scene{}).Where("creator = ? AND scene_id = ?", creator, sceneId).First(&scene).Error
	if err != nil {
		return nil, err
	}
	return scene, nil
}

func CopyStartScene(username string, sceneId int64) (url string, err error) {
	scene, err := GetSceneBySceneId("491213694122852611")
	if err != nil {
		return "", err
	}
	scene.Creator = username
	scene.SceneId = sceneId
	err = model.DB.Model(&Scene{}).Create(&scene).Error
	return scene.COSUrl, err
}

func GetStartScene(username string) (string, error) {
	scene := &Scene{}
	// parent scene id = 0
	err := model.DB.Model(&Scene{}).Where("parent_scene_id = ? AND creator = ?", 0, username).First(&scene).Error
	return scene.COSUrl, err
}

func GetCosUrlByCreatorAndChooseContent(creator, chooseContent string) (string, string, error) {
	scene := &Scene{}
	err := model.DB.Model(&Scene{}).Where("creator = ? AND choose_content = ? and is_init = 1", creator, chooseContent).First(&scene).Error
	return strconv.FormatInt(scene.SceneId, 10), scene.COSUrl, err
}

func SaveUngeneratedScene(sceneId, parentSceneId int64, choose, username string, isInit int) (error, *Scene) {
	scene := &Scene{
		SceneId:       sceneId,
		ChooseContent: choose,
		Creator:       username,
		ParentSceneId: parentSceneId,
		IsInit: isInit,
	}
	err := model.DB.Model(&Scene{}).Create(&scene).Error
	return err, scene
}

func GetSceneBySceneIdAndCreator(sceneId, creator string) (scene *Scene, err error) {
	scene = &Scene{}
	err = model.DB.Model(&Scene{}).Where("scene_id = ? AND creator = ?", sceneId, creator).First(&scene).Error
	if err != nil {
		return nil, err
	}
	return scene, nil
}

func GetSceneByCreatorAndInitId(username string, initId int) (scene *Scene, err error) {
	scene = &Scene{}
	err = model.DB.Model(&Scene{}).Where("creator = ? AND is_init = ?", username, initId).First(&scene).Error
	if err != nil {
		return nil, err
	}
	return scene, nil
}