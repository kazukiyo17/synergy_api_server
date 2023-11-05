package rmq

import (
	//"github.com/go-redis/redis/v8"
	"github.com/gomodule/redigo/redis"
	"github.com/kazukiyo17/fake_buddha_server/common/conf"
)

type StreamMQ struct {
	// Redis客户端
	client *redis.Pool
	// 最大消息数量，如果大于这个数量，旧消息会被删除，0表示不管
	maxLen int64
	// Approx是配合MaxLen使用的，表示几乎精确的删除消息，也就是不完全精确，由于stream内部是流，所以设置此参数xadd会更加高效
	approx bool
}

func NewStreamMQ(maxLen int, approx bool) *StreamMQ {
	return &StreamMQ{
		client: conf.C.RedisMQConn,
		maxLen: int64(maxLen),
		approx: approx,
	}
}

func (q *StreamMQ) SendMsg(msg *Msg) error {
	conn := q.client.Get()
	defer conn.Close()
	return conn.Send("XADD", msg.Topic, "MAXLEN", q.maxLen, "*", "body", msg.Body)
	//return conf.C.RedisMQConn.XAdd(ctx, &redis.XAddArgs{
	//	Stream: msg.Topic,
	//	MaxLen: q.maxLen,
	//	Approx: q.approx,
	//	ID:     "*",
	//	Values: []interface{}{"body", msg.Body},
	//}).Err()
}
