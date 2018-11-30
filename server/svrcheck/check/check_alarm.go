package check

import (
	gredis "github.com/go-redis/redis"
	"github.com/json-iterator/go"
	"github.com/streadway/amqp"
	"sync"
	"time"
	"xcurrency-push/config"
	"xcurrency-push/model"
	"xcurrency-push/module/rabbitmq"
	"xcurrency-push/module/redis"
	"xcurrency-push/util"
)

type CheckAlarm struct {
	UExrate Exrate
}

/**
 * @des 预警信息检查，采用Redis的pipeline进行高效数据读取，并用多线程进行数据处理
 * 另外，通过chan控制最大并发数，通过sync.WaitGroup保证所有线程执行完毕
 * @author 于朝鹏
 * @date 2018年11月13日 16:10
 */
func (checkAlarm *CheckAlarm) Update() {
	defer func() {
		if err := recover(); err != nil {
			util.LogErrorM(err, "Recover [check_alarm.go CheckAlarm] error.")
		}
	}()

	cfg := config.GetConfig()

	max := cfg.GetInt("check.maxSend") // 最大发送线程数
	p := cfg.GetInt("check.pipeExecCount")

	send := make(chan int, max) // 并发线程数量控制
	wg := new(sync.WaitGroup)   // 并发线程执行控制

	conn := rabbitmq.GetConn()
	chanArr := rabbitmq.GetChannelArr(conn, max)
	defer rabbitmq.CloseChannelArr(chanArr)

	for i := 0; i < max; i++ {
		send <- i
	}

	util.LogInfo("===check alarm start. max rabbit channel:", max)
	if client, err := redis.GetClient(); err == nil {
		pipe := client.Pipeline()
		defer redis.ClosePipeline(pipe)

		// t1 := time.Now().UnixNano()
		if ids, err := client.SMembers("alarm:ids").Result(); err == nil {

			// 以p为基数，把总数量分成n组，每组通过Pipeline批量查询
			arr := GroupAlarm(ids, p)
			for _, id := range arr {

				// 批量查询
				if alarms, err := redis.GetAlarmPipeline(id, pipe); err == nil {
					// 过滤符合条件的数据
					fm := filterAlarm(alarms, checkAlarm)
					// 更新Alarm的最后预警时间
					redis.UpdateAlarmTime(pipe, fm)

					x := <-send
					wg.Add(1)

					// 多线程推送数据到rabbitmq队列
					go DoSend(fm, chanArr[x], send, x, wg, nil)
				}
			}

			// 等待所有线程执行完毕
			wg.Wait()

			// t2 := time.Now().UnixNano()
			// util.LogInfo("===check alarm finish. time:", (t2-t1)/1e6, " millisecond")
			// util.LogInfo("***********************************************************")
		}
	}
}

/**
 * @desc 以p为基数，把总数量分成n组
 * @author 于朝鹏
 * @date 2018年11月14日 15:52
 */
func GroupAlarm(ids []string, p int) [][]string {
	if len(ids) == 0 {
		return [][]string{}
	}

	s := len(ids) // 总数
	c := p        // 基数
	if c > s {
		c = s
	}
	n := s / c // 组数
	r := s % c // 余数

	arr := make([][]string, 0, n+1)

	for i := 0; i < n; i++ {
		var id []string
		if i > 0 && r > 0 && i == n-1 {
			id = ids[i*c : ((i+1)*c + r)]
		} else {
			id = ids[i*c : ((i + 1) * c)]
		}
		arr = append(arr, id)
	}
	return arr
}

/**
 * @desc 过滤符合条件的数据
 * @author 于朝鹏
 * @date 2018年11月14日 15:13
 */
func filterAlarm(alarms []model.Alarm, checkAlarm *CheckAlarm) []model.Alarm {
	a := make([]model.Alarm, 0, 500)
	for _, alarm := range alarms {
		lprice := alarm.LbPrice
		uprice := alarm.UbPrice
		baseCur := alarm.BaseCur
		targetCur := alarm.TargetCur
		enable := alarm.Enabled
		devId := alarm.DeviceId

		// 数据验证
		if enable && len(devId) >= 0 && len(baseCur) > 0 && len(targetCur) > 0 {

			// 时间间隔判断
			ltt := alarm.Ltt
			ct := time.Now().Unix()
			// 60秒1次
			if ltt > 0 && (ct-ltt) < 1*60 {
				continue
			}

			// 价格判断
			price, err := checkAlarm.UExrate.GetPrice(baseCur, targetCur)
			if err != nil {
				util.LogError(err)
			}
			util.LogInfo("GetPrice: ["+baseCur+"/"+targetCur+"]", lprice, uprice, price)

			if checkPrice(lprice, uprice, price) {
				alarm.Price = price
				a = append(a, alarm)
			}
		}
	}
	return a
}

/**
 * @desc 发送数据到RabbitMQ队列
 * @author 于朝鹏
 * @date 2018年11月2日 11:40
 */
func DoSend(alarms []model.Alarm, ch *amqp.Channel, send chan int, x int, wg *sync.WaitGroup, pipe gredis.Pipeliner) {
	defer func() {
		if err := recover(); err != nil {
			util.LogErrorM(err, "Recover [check_alarm.go CheckAlarm DoSend] error.")
		}
	}()
	defer wg.Done()
	defer util.WriteChanNum(send, x)

	// 更新Alarm的最后预警时间
	for _, alarm := range alarms {
		var json = jsoniter.ConfigCompatibleWithStandardLibrary
		if byt, err := json.Marshal(&alarm); err == nil {

			// 推送数据到RabbitMQ队列
			rabbitmq.DoSend(ch, byt)
		}
	}
}

/**
 * @desc 价格范围检查，判断用户是否在预警价格区间
 * @author 于朝鹏
 * @date 2018年11月14日 15:01
 */
func checkPrice(lprice float64, uprice float64, price float64) bool {
	if lprice > uprice {
		util.LogError("The lowest price cannot be greater than the highest price")
		return false
	}
	if lprice == 0 && uprice == 0 {
		return false
	}
	if lprice == 0 {
		return price < uprice
	}
	if uprice == 0 {
		return price > lprice
	}
	if lprice == uprice {
		return price == lprice
	}
	return (price > lprice && price < uprice)
}

/**
 * @desc 定期器循环运行预警信息检查
 * @author 于朝鹏
 * @date 2018年10月8日 15:20
 */
func (alarm *CheckAlarm) Loop() {
	util.LogInfo("check alarm loop start.")

	cfg := config.GetConfig()
	t := cfg.GetInt("check.alarmTickTime") // 定时执行时间（秒）
	TickLoop(alarm, t)
}
