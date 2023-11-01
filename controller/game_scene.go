package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kazukiyo17/fake_buddha_server/common"
	"github.com/spf13/cast"
)

type SceneInfoOut struct {
	SceneID int64  `json:"scene_id"` // 剧本ID
	Url     string `json:"url"`      // 存储地址
	UserID  int64  `json:"user_id"`  // 用户ID
}

func SceneInfo(c *gin.Context) {
	// 参数校验
	sceneId := cast.ToInt64(c.Query("activity_id"))
	if sceneId == 0 {
		common.SendParamError(c, "activityId")
	}

	// 查询剧本信息
	activityInfo, err := activity.NewActivityMgr().GetActivityInfoById(c, activityId)

}
