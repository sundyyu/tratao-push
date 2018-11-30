package redis

import (
	"fmt"
	// "github.com/go-redis/redis"
	// "reflect"
	// "encoding/json"
	"strconv"
	"testing"
	"time"
	"xcurrency-push/model"
	"xcurrency-push/util"
)

func TestRedis(t *testing.T) {

	client, err := GetClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	pipe := client.Pipeline()
	defer pipe.Close()

	//#====================================Pipeline返回数据类型测试========================
	// pipe.HMGet("alarm:1", "Account", "BaseCur")
	// pipe.HMGet("alarm:2", "Account", "BaseCur")
	// pipe.HMGet("alarm:3", "Account", "BaseCur")
	//
	// data, err := pipe.Exec()
	//
	// if err != nil {
	// 	fmt.Println(data)
	// 	return
	// }
	//
	// for _, d := range data {
	// 	if s, ok := d.(*redis.SliceCmd); ok {
	// 		fmt.Println(s.Result())
	// 	}
	// }

	//#====================================分组Pipeline测试========================
	// ids := make([]string, 10)
	// for j := 0; j < 105; j++ {
	// 	ids = append(ids, strconv.Itoa(j))
	// }
	//
	// c := 10
	// size := len(ids)
	// if size < c {
	// 	c = size
	// }
	//
	// l := size / c
	// r := size % c
	// for i := 0; i < l; i++ {
	//
	// 	if i > 0 && r > 0 && i == l-1 {
	// 		id := ids[i*c : ((i+1)*c + r)]
	// 		fmt.Println(id)
	// 	} else {
	// 		id := ids[i*c : ((i + 1) * c)]
	// 		fmt.Println(id)
	// 	}
	// }

	//#====================================Pipeline查询测试========================
	// r, err := GetAlarmPipeline([]string{"1", "2"}, pipe)
	// if err == nil {
	// 	fmt.Println(r)
	// }

	// fmt.Println(data)
	// fmt.Println(m)

	// if cerr != nil {
	// 	panic(cerr)
	// }
	//

	//#====================================添加测试数据========================
	fmt.Println("test add alarm start...")
	t1 := time.Now().UnixNano()

	// "id":1,
	// "account":"test9",
	// "basecur":"USD",
	// "targetcur":"CNY",
	// "lbprice":4.0,
	// "ubprice":7.3,
	// "enable":true,
	// "quote":"5",
	// "ptype":"LB"

	n := 500000
	for i := 1; i < n; i++ {
		alarm := model.Alarm{}
		alarm.Account = "test" + strconv.Itoa(i)
		alarm.BaseCur = "USD"
		alarm.TargetCur = "CNY"
		alarm.LbPrice = 4.0
		alarm.Enabled = true
		alarm.Price = 6.9
		alarm.CreateTime = time.Now().Unix()
		alarm.UpdateTime = time.Now().Unix()

		alarm.Id = GetId()

		if i <= 50000 {
			alarm.UbPrice = 8.0
		} else {
			alarm.UbPrice = 5.0
		}

		// AddAlarm(alarm)

		id := strconv.FormatInt(alarm.Id, 10)
		alarmMap := util.Struct2Map(alarm)

		// 存储数据
		pipe.HMSet("alarm:"+id, alarmMap)

		// 存储account对应的主键Id
		pipe.SAdd("alarm:"+alarm.Account+":id", id)

		// 存储所有的主键Id
		pipe.SAdd("alarm:ids", id)

		if i%5000 == 0 || i == (n-1) {
			pipe.Exec()
			fmt.Println("add count :", i)
		}

	}

	t2 := time.Now().UnixNano()
	fmt.Println("test add alarm finish.  time: ", (t2-t1)/1e6)

	// ssh-keygen -t rsa -C "sundy.yu@tratao.com"

}
