package main

import (
	// "flag"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	// "os"
	"time"
	// "xcurrency-push/util"
)

var hosts = []string{"192.168.1.104:2181"}

var path = "/trataolock"
var node = "/trataonode-"

var flags int32 = 0
var data = []byte("hello,this is a zk go test demo!!!")
var acls = zk.WorldACL(zk.PermAll)

var zkconn *zk.Conn
var nodeId string

func main() {
	// c := flag.String("config", "", "配置文件参数")
	// flag.Parse()
	//
	// p := *c
	// _, err := os.Stat(p)
	// if err == nil {
	// 	util.LogInfoF("config file %s exists", p)
	// } else if os.IsNotExist(err) {
	// 	util.LogInfoF("config file %s not exists", p)
	// } else {
	// 	util.LogInfoF("config file %s stat error: %v", p, err)
	// }
	// node += p

	forever := make(chan int, 1)

	option := zk.WithEventCallback(callback)

	conn, _, err := zk.Connect(hosts, time.Second*5, option)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	zkconn = conn

	exist, _, _, err := conn.ExistsW(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	// go watchCreataNode(event)

	if !exist {
		create(conn, path, data, flags)
		fmt.Println("created path.")
	}

	time.Sleep(time.Second * 2)

	exist, _, _, err = conn.ExistsW(path + node)
	if err != nil {
		fmt.Println(err)
		return
	}
	// go watchCreataNode(event)

	if !exist {
		createNode(conn, path+node, data)
		fmt.Println("created node: " + path + node)
	}

	go watchCreataNode(path)

	<-forever
}

func callback(event zk.Event) {
	// fmt.Println("path:", event.Path)
	// fmt.Println("type:", event.Type.String())
	// fmt.Println("state:", event.State.String())
	// fmt.Println("-------------------")

	// if event.Type == zk.EventNodeCreated || event.Type == zk.EventNodeDeleted {
	// 	child, _, err := zkconn.Children(path)
	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// 	fmt.Println("ChildrenW:", child)
	// }

}

func create(conn *zk.Conn, path string, data []byte, f int32) {
	_, err := conn.Create(path, data, f, acls)
	if err != nil {
		fmt.Println(err)
		// return
	}
}

func createNode(conn *zk.Conn, path string, data []byte) {
	// conn.Create(path, data, zk.FlagEphemeral, acls)
	str, err := conn.CreateProtectedEphemeralSequential(path, data, acls)
	if err != nil {
		fmt.Println(err)
		return
	}
	nodeId = str

	// _, err := conn.CreateProtectedEphemeralSequential(path, data, acls)
	// if err != nil {
	// 	fmt.Println(err)
	// 	// return
	// }
}

func watchCreataNode(path string) {

	for {
		ch, _, _, _ := zkconn.ChildrenW(path)
		// event := <-ech
		// if event.Type == zk.EventNodeCreated || event.Type == zk.EventNodeChildrenChanged {
		// 	fmt.Println("---------WatchNode----------")
		// 	fmt.Println("path:", event.Path)
		// 	fmt.Println("type:", event.Type.String())
		// 	fmt.Println("state:", event.State.String())
		fmt.Println("children:", ch)
		// 	fmt.Println("curNodeId:", nodeId)
		// 	fmt.Println("---------WatchNode----------")
		// }
		time.Sleep(time.Second * 5)
	}
}
