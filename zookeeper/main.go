package main

import (
	"xcurrency-push/config"
	"xcurrency-push/util"
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
	zkclient.CreateNode(conn, zkclient.Path, zkclient.Data, zkclient.Flags)
	nid := zkclient.CreateSeqNode(conn, zkclient.Node, zkclient.Data)
	ch := make(chan []string, 1)
	zkclient.WatchChildren(conn, zkclient.Path, ch)
	for {
		child := <-ch
		c := zkclient.NodeArr2IntArr(child)
		n := zkclient.NodeIdToInt(nid)
		util.LogInfo(n, c)
		if n == c[0] {
			util.LogInfo("start run check.")
			break
		}
	}

	<-forever
}
