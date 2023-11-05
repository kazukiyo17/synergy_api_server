package scene

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/fake_buddha_server/common"
	"github.com/kazukiyo17/fake_buddha_server/service/scene_service"
	"github.com/spf13/cast"
	"net/http"
)

type SceneInfoOut struct {
	SceneID int64  `json:"scene_id"` // 剧本ID
	Url     string `json:"url"`      // 存储地址
	UserID  int64  `json:"user_id"`  // 用户ID
}

func SceneInfo(c *gin.Context) {
	// 参数校验
	var sceneId int64 = cast.ToInt64(c.Query("activity_id"))
	if sceneId == 0 {
		common.SendParamError(c, "activityId")
	}

	// 查询信息
	sceneService := scene_service.Scene{
		ID: sceneId,
	}
	scene, err := sceneService.GetSceneInfoById(sceneId)
	if err != nil {
		common.SendParamError(c, err.Error())
		return
	}

	// 返回信息
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "success",
		"data": SceneInfoOut{
			SceneID: scene.ID,
			Url:     scene.Url,
			UserID:  scene.CreatorUserId,
		},
	})

}
