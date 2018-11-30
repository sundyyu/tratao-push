package main

import (
	"fmt"
	"xcurrency-push/config"
	"xcurrency-push/part"
	"xcurrency-push/zookeeper/zkclient"
)

func main2() {
	config.LoadConfig("../config/cfg.yaml")

	forever := make(chan int, 1)
	// alarm check

	conn := zkclient.GetConn()
	if conn == nil {
		return
	}
	defer conn.Close()
	zkclient.AcquireMetux(conn)

	<-forever
}

func main() {
	config.LoadConfig("../config/cfg.yaml")

	servers := []string{"server0", "server1", "server2"}
	s := part.ServerIndex(servers)
	part.AddCount(s)
	fmt.Println(s)

	part.ServerPart()
}
