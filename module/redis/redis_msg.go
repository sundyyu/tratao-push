package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"xcurrency-push/model"
	"xcurrency-push/util"
)

func AddPushMsg(msg model.PushMsg) (int64, error) {
	incrId := GetId() // auto generator id
	if incrId < 0 {
		return -1, errors.New("generator id error")
	}

	msg.Id = incrId
	m := util.Struct2Map(msg)

	client, cerr := GetClient()
	if cerr != nil {
		return -1, cerr
	}

	// 使用事务类型pipeline
	pipe := client.TxPipeline()
	defer pipe.Close()

	// int64 to string
	id := strconv.FormatInt(incrId, 10)

	// 存储数据
	pipe.HMSet("alarm:pushmsg:"+id, m)

	// 存储所有的主键Id
	pipe.SAdd("alarm:pushmsg:ids", id)

	if _, err := pipe.Exec(); err != nil {
		return -1, err
	}
	return incrId, nil
}

// 使用Pipeline批量查询数据
func GetPushMsgPipeline(ids []string, pipe redis.Pipeliner) ([]model.PushMsg, error) {
	msgs := make([]model.PushMsg, 0, len(ids))
	for _, id := range ids {
		pipe.HMGet("alarm:pushmsg:"+id, model.GetPushMsgFields()...)
	}

	data, err := pipe.Exec()
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		if s, ok := d.(*redis.SliceCmd); ok {
			if r, err := s.Result(); err == nil {
				msgs = append(msgs, model.ResultToPushMsg(r))
			}
		}
	}
	return msgs, nil
}

// 使用Pipeline批量删除数据
func DelPushMsgArr(ids []string, pipe redis.Pipeliner) error {
	for _, id := range ids {
		pipe.Del("alarm:pushmsg:" + id)
		pipe.SRem("alarm:pushmsg:ids", id)
	}

	_, err := pipe.Exec()
	if err != nil {
		return err
	}
	return nil
}
