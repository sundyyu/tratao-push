package zkclient

import (
	"github.com/samuel/go-zookeeper/zk"
	"sort"
	"strconv"
	"strings"
	"time"
	"xcurrency-push/config"
	"xcurrency-push/util"
)

// var Path = ""
// var Node = ""
// var NodeSuffix = ""

var Flags int32 = 0
var Data = []byte("trataodata")
var Acls = zk.WorldACL(zk.PermAll)

func GetConn() *zk.Conn {
	cfg := config.GetConfig()
	host := cfg.GetString("zookeeper.addrs")
	hosts := strings.Split(host, ",")
	timeout := cfg.GetInt("zookeeper.timeout")

	option := zk.WithEventCallback(nil)
	conn, _, err := zk.Connect(hosts, time.Second*time.Duration(timeout), option)
	if err != nil {
		util.LogError(err)
		return nil
	}
	return conn
}

func CreateNode(conn *zk.Conn, path string, data []byte, flags int32) string {
	exist, _, err := conn.Exists(path)
	if err != nil || exist {
		util.LogError(err)
		return ""
	}

	nodeId, err := conn.Create(path, data, flags, Acls)
	if err != nil {
		util.LogError(err)
		return ""
	}
	return nodeId
}

func CreateSeqNode(conn *zk.Conn, path string, data []byte) string {
	nodeId, err := conn.CreateProtectedEphemeralSequential(path, data, Acls)
	if err != nil {
		util.LogError(err)
		return ""
	}
	return nodeId
}

func WatchChildren(zkconn *zk.Conn, path string, ch chan []string) {
	go func() {
		cfg := config.GetConfig()
		t := cfg.GetInt("zookeeper.watchtime")

		for {
			c, _, _, err := zkconn.ChildrenW(path)
			if err != nil {
				util.LogError(err)
			}
			ch <- c
			time.Sleep(time.Second * time.Duration(t))
		}
	}()
}

// Node数组转为Int数组
func NodeArr2IntArr(strs []string, node string) []int {
	if strs == nil {
		return []int{}
	}
	r := make([]int, 0, len(strs))
	for _, str := range strs {
		r = append(r, NodeIdToInt(str, node))
	}
	sort.Ints(r)
	return r
}

// NodeId转为Int数字
func NodeIdToInt(str string, node string) int {
	sarr := strings.Split(str, node)
	if len(sarr) > 1 {
		s := sarr[len(sarr)-1]
		if i, err := strconv.Atoi(s); err == nil {
			return i
		}
	}
	return -1
}

// 获得互斥锁，否则一直被阻塞
func AcquireMetux(conn *zk.Conn) {
	cfg := config.GetConfig()
	path := cfg.GetString("node.root")
	node := cfg.GetString("node.nodeMutex")
	Acquire(conn, path, node)
}

// 获得分片锁，否则一直被阻塞
func AcquirePart(conn *zk.Conn) {
	cfg := config.GetConfig()
	path := cfg.GetString("node.root")
	node := cfg.GetString("node.nodePart")
	Acquire(conn, path, node)
}

func Acquire(conn *zk.Conn, path string, node string) {
	rootPath := "/" + path + "/" + node
	nodeName := node + "-"
	CreateNode(conn, rootPath, Data, Flags)
	nid := CreateSeqNode(conn, rootPath+"/"+nodeName, Data)
	ch := make(chan []string, 1)
	WatchChildren(conn, rootPath, ch)
	acch := make(chan int, 1)
	for {
		child := <-ch
		c := NodeArr2IntArr(child, nodeName)
		n := NodeIdToInt(nid, nodeName)
		util.LogInfo(n, c, child)
		if n == c[0] {
			acch <- 1
			util.LogInfo(node + " get lock and acquire success.")
			break
		}
	}
	<-acch
}
