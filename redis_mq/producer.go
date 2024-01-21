package redis_mq

import (
	"fmt"
	"log"
)

func Produce(chooseId, username string) {
	msg, err := redisMQClient.PutMsg(chooseId, username)
	if err != nil {
		fmt.Println("PutMsg err:", err)
		return
	}
	//log.Printf("-------------------------------------------------")
	log.Printf( "======================put msg: %v", msg)
}
