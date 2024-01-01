package redis_mq

import "fmt"

func Produce(chooseId string, sceneId string) {
	msg, err := redisMQClient.PutMsg(chooseId, sceneId)
	if err != nil {
		fmt.Println("PutMsg err:", err)
		return
	}
	fmt.Println("PutMsg:", msg)
}

func Check(chooseId string) bool {
	msg, err := redisMQClient.CheckPeddingList(chooseId)
	if err != nil {
		return false
	}
	return msg
}
