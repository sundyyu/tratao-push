package main

import (
	"xcurrency-push/config"
	"xcurrency-push/zookeeper/zkclient"
)

func main() {
	config.LoadConfig("../config/cfg.yaml")

	forever := make(chan int, 1)
	conn := zkclient.GetConn()
	if conn == nil {
		return
	}
	defer conn.Close()

	zkclient.AcquireMetux(conn)
	zkclient.AcquirePart(conn)

	<-forever
}
