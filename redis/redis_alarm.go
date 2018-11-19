package redis

import (
	"errors"
	"github.com/go-redis/redis"
	"strconv"
	"time"
	"tratao-push/model"
	"tratao-push/util"
)

func AddAlarm(alarm model.Alarm) (int64, error) {
	incrId := GetId() // auto generator id
	if incrId < 0 {
		return -1, errors.New("generator id error")
	}

	// int64 to string
	id := strconv.FormatInt(incrId, 10)
	alarm.Id = incrId
	account := alarm.Account

	// 为新增的添加设备信息
	if list, err := ListAlarm(account); err == nil {
		if list != nil && len(list) > 0 {
			alarm.DeviceId = list[0].DeviceId
			alarm.DeviceOS = list[0].DeviceOS
			alarm.DeviceCountry = list[0].DeviceCountry
			alarm.DeviceLang = list[0].DeviceLang
			alarm.AppKey = list[0].AppKey
		}
	}
	m := util.Struct2Map(alarm)

	client, cerr := GetClient()
	if cerr != nil {
		return -1, cerr
	}

	// 使用事务类型pipeline
	pipe := client.TxPipeline()
	defer pipe.Close()

	// 存储数据
	pipe.HMSet("alarm:"+id, m)

	// 存储account对应的主键Id
	pipe.SAdd("alarm:"+alarm.Account+":id", id)

	// 存储所有的主键Id
	pipe.SAdd("alarm:ids", id)

	if _, err := pipe.Exec(); err != nil {
		return -1, err
	}
	return incrId, nil
}

func UpdateAlarm(alarm model.Alarm) error {
	id := strconv.FormatInt(alarm.Id, 10)
	client, cerr := GetClient()
	if cerr != nil {
		return cerr
	}

	a, err := GetAlarm(id)
	if err != nil {
		return errors.New("UpdateAlarm error, id: " + id)
	}

	a.BaseCur = alarm.BaseCur
	a.TargetCur = alarm.TargetCur
	a.LbPrice = alarm.LbPrice
	a.UbPrice = alarm.UbPrice
	a.Enabled = alarm.Enabled
	a.UpdateTime = alarm.UpdateTime
	m := util.Struct2Map(a)

	// 更新数据
	if err := client.HMSet("alarm:"+id, m).Err(); err != nil {
		return err
	}

	return nil
}

func UpdateDevice(alarm model.Alarm) error {
	account := alarm.Account

	client, cerr := GetClient()
	if cerr != nil {
		return cerr
	}

	list, err := ListAlarm(account)
	if err != nil {
		return errors.New("ListAlarm error, account: " + account)
	}

	// 使用事务类型pipeline
	pipe := client.TxPipeline()
	defer pipe.Close()

	for _, a := range list {
		a.DeviceId = alarm.DeviceId
		a.DeviceOS = alarm.DeviceOS
		a.DeviceCountry = alarm.DeviceCountry
		a.DeviceLang = alarm.DeviceLang
		a.AppKey = alarm.AppKey
		a.UpdateTime = alarm.UpdateTime
		m := util.Struct2Map(a)

		id := strconv.FormatInt(a.Id, 10)
		pipe.HMSet("alarm:"+id, m)
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}

func ListAlarm(account string) ([]model.Alarm, error) {
	client, cerr := GetClient()
	if cerr != nil {
		return nil, cerr
	}

	// 根据acount查询所有key
	ids, err := client.SMembers("alarm:" + account + ":id").Result()
	if err != nil {
		return nil, errors.New("find account ids error")
	}

	pipe := client.Pipeline()
	defer pipe.Close()

	alarms, err := GetAlarmPipeline(ids, pipe)
	return alarms, nil
}

// 根据ID查询数据
func GetAlarm(id string) (model.Alarm, error) {
	client, cerr := GetClient()
	if cerr != nil {
		return model.Alarm{}, cerr
	}
	return GetAlarmCli(id, client)
}

// 根据ID, Client查询数据
func GetAlarmCli(id string, client *redis.Client) (model.Alarm, error) {
	result, err := client.HMGet("alarm:"+id, model.GetFields()...).Result()
	if err == redis.Nil {
		return model.Alarm{}, errors.New("alarm is not exist.")
	} else if err != nil {
		return model.Alarm{}, errors.New("find alarm error")
	}

	alarm := model.ResultToAlarm(result)
	return alarm, nil
}

// 使用Pipeline批量查询数据
func GetAlarmPipeline(ids []string, pipe redis.Pipeliner) ([]model.Alarm, error) {
	alarms := make([]model.Alarm, 0, len(ids))
	for _, id := range ids {
		pipe.HMGet("alarm:"+id, model.GetFields()...)
	}

	data, err := pipe.Exec()
	if err != nil {
		return nil, err
	}

	for _, d := range data {
		if s, ok := d.(*redis.SliceCmd); ok {
			if r, err := s.Result(); err == nil {
				alarms = append(alarms, model.ResultToAlarm(r))
			}
		}
	}
	return alarms, nil
}

func DelAlarm(ids []string) error {
	client, cerr := GetClient()
	if cerr != nil {
		return cerr
	}

	// 使用事务类型pipeline
	pipe := client.TxPipeline()
	defer pipe.Close()

	for _, id := range ids {

		// 根据ID查询Account
		result, err := client.HMGet("alarm:"+id, "Account").Result()
		if err == nil && err != redis.Nil {

			// 删除account中的主键Id
			if acc, ok := result[0].(string); ok {
				pipe.SRem("alarm:"+acc+":id", id)
			}
		}

		// 删除alarm数据
		pipe.Del("alarm:" + id)

		// 删除主键Id
		pipe.SRem("alarm:ids", id)
	}

	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}

func UpdateAlarmTime(pipe redis.Pipeliner, alarms []model.Alarm) error {
	t := time.Now().Unix()
	for _, alarm := range alarms {
		id := strconv.FormatInt(alarm.Id, 10)
		pipe.HMSet("alarm:"+id, map[string]interface{}{"Ltt": t})
	}
	if _, err := pipe.Exec(); err != nil {
		return err
	}
	return nil
}
