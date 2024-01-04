package scene

import (
	"github.com/kazukiyo17/synergy_api_server/model"
	"strconv"
)

type ModelScene struct {
	model.Model
	SceneId       int64  `json:"tag_id" gorm:"type:bigint;index"`
	ChooseContent string `json:"choose_content" gorm:"type:varchar(255)"`
	//CreatorId     int64  `json:"creator_id" gorm:"type:bigint;index"`
	Creator       string `json:"creator" gorm:"type:varchar(255)"`
	ParentSceneId int64  `json:"parent_scene_id" gorm:"type:bigint;index"`
	COSUrl        string `json:"cos_url" gorm:"type:varchar(255)"`
	ShortDesc     string `json:"desc" gorm:"type:varchar(600)"`
}

func GetCosUrlBySceneId(sceneId string) (cosUrl string, err error) {
	scene := &ModelScene{}
	err = model.DB.Model(&ModelScene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	cosUrl = scene.COSUrl
	return cosUrl, nil
}

func GetSceneIdByParentSceneId(parentSceneId string) (sceneIds []string, err error) {
	scenes := make([]*ModelScene, 0)
	err = model.DB.Model(&ModelScene{}).Where("parent_scene_id = ?", parentSceneId).Find(&scenes).Error
	if err != nil {
		return nil, err
	}
	for _, scene := range scenes {
		sceneIds = append(sceneIds, strconv.FormatInt(scene.SceneId, 10))
	}
	return sceneIds, nil
}

func GetCreatorBySceneId(sceneId string) (creator string, err error) {
	scene := &ModelScene{}
	err = model.DB.Model(&ModelScene{}).Where("scene_id = ?", sceneId).First(&scene).Error
	if err != nil {
		return "", err
	}
	creator = scene.Creator
	return creator, nil
}
