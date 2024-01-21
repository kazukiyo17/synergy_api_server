package redis_mq

import "fmt"

func Produce(chooseId, username string) {
	msg, err := redisMQClient.PutMsg(chooseId, username)
	if err != nil {
		fmt.Println("PutMsg err:", err)
		return
	}
	fmt.Println("PutMsg:", msg)
}
