package pgclient

import (
	"xcurrency-push/model"
	"xcurrency-push/util"
)

// 插入PushMsg数据
func InsertPushMsg(p model.PushMsg) error {
	sql := "INSERT INTO trataopush (account,title,body,token,os,lang,country,createtime) VALUES($1,$2,$3,$4,$5,$6,$7,$8)"
	err := ExecBySQL(sql, p.Account, p.Title, p.Body, p.DeviceId, p.DeviceOS, p.DeviceLang, p.DeviceCountry, p.CreateTime)
	if err != nil {
		return err
	}
	return nil
}

// 查询PushMsg数据
func QueryPushMsg() ([]model.PushMsg, error) {
	db, err := GetConn()
	if err != nil {
		return nil, err
	}
	defer ReleaseConn(db)
	//查询数据
	rows, err := db.Query("SELECT * FROM trataopush")
	if err != nil {
		util.LogInfo(err)
		return nil, err
	}

	list := make([]model.PushMsg, 0, 10)
	for rows.Next() {
		pushmsg := model.PushMsg{}
		err := rows.Scan(&pushmsg.Id, &pushmsg.Account, &pushmsg.Title, &pushmsg.Body,
			&pushmsg.DeviceId, &pushmsg.DeviceOS, &pushmsg.DeviceLang, &pushmsg.DeviceCountry,
			&pushmsg.CreateTime)
		if err == nil {
			list = append(list, pushmsg)
		}
	}
	return list, nil
}
