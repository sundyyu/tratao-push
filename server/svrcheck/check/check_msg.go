package check

import (
	"github.com/json-iterator/go"
	"xcurrency-push/config"
	"xcurrency-push/rabbitmq"
	"xcurrency-push/redis"
	"xcurrency-push/util"
)

type CheckMessage struct {
}

/**
 * @des 推送信息检查
 * @author 于朝鹏
 * @date 2018年11月13日 16:10
 */
func (checkMessage *CheckMessage) Update() {
	defer func() {
		if err := recover(); err != nil {
			util.LogErrorM(err, "Recover [check_msg.go CheckMessage] error.")
		}
	}()

	cfg := config.GetConfig()
	queue := cfg.GetString("check.pushMsgQueue")

	conn := rabbitmq.GetConn()
	ch := rabbitmq.GetChannelQueue(conn, queue)
	defer conn.Close()
	defer ch.Close()

	util.LogInfo("===check push message start. ")
	if client, err := redis.GetClient(); err == nil {
		pipe := client.Pipeline()
		defer redis.ClosePipeline(pipe)

		if ids, err := client.SMembers("alarm:pushmsg:ids").Result(); err == nil {

			// 批量查询
			if msgs, err := redis.GetPushMsgPipeline(ids, pipe); err == nil {
				for _, msg := range msgs {
					var json = jsoniter.ConfigCompatibleWithStandardLibrary
					if byt, err := json.Marshal(&msg); err == nil {

						// 推送数据到RabbitMQ队列
						rabbitmq.DoSendQueue(ch, byt, queue)
					}

					// 从Redis删除数据
					redis.DelPushMsgArr(ids, pipe)
				}
			}
		}
	}
}

/**
 * @desc 定期器循环运行预警信息检查
 * @author 于朝鹏
 * @date 2018年10月8日 15:20
 */
func (check *CheckMessage) Loop() {
	util.LogInfo("check pushmsg loop start.")

	cfg := config.GetConfig()
	t := cfg.GetInt("check.pushMsgTickTime") // 定时执行时间（秒）
	TickLoop(check, t)
}
