package model

type Scene struct {
	ID            int64  `json:"id"`
	Url           string `json:"url"`
	CreatorUserId int64  `json:"creator_user_id"`
	Prompt        string `json:"prompt"`
}

func GetSceneById(sceneId int64) (*Scene, error) {
	var scene Scene
	err := db.Select("id, url, creator_user_id, prompt").Where("id = ?", sceneId).First(&scene).Error
	if err != nil {
		return nil, err
	}
	return &scene, nil
}
