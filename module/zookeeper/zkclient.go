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

func CreateNode(conn *zk.Conn, path string, data []byte, flags int32) (string, error) {
	exist, _, err := conn.Exists(path)
	if err != nil {
		util.LogError(err)
		return "", err
	}
	if exist {
		byt, _, _ := conn.Get(path)
		return string(byt), nil
	}

	nodeId, err := conn.Create(path, data, flags, Acls)
	if err != nil {
		util.LogError(err)
		return "", err
	}
	return nodeId, nil
}

func CreateSeqNode(conn *zk.Conn, path string, data []byte) (string, error) {
	nodeId, err := conn.CreateProtectedEphemeralSequential(path, data, Acls)
	if err != nil {
		util.LogError(err)
		return "", err
	}
	return nodeId, nil
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
	path := cfg.GetString("zkserver.root")
	node := cfg.GetString("zkserver.nodeMutex")
	Acquire(conn, path, node)
}

// 获得分片锁，否则一直被阻塞
func AcquirePart(conn *zk.Conn) {
	cfg := config.GetConfig()
	path := cfg.GetString("zkserver.root")
	node := cfg.GetString("zkserver.nodePart")
	Acquire(conn, path, node)
}

func Acquire(conn *zk.Conn, path string, node string) {
	rootPath := "/" + path + "/" + node
	nodeName := node + "-"
	CreateNode(conn, rootPath, Data, Flags)
	nid, _ := CreateSeqNode(conn, rootPath+"/"+nodeName, Data)
	ch := make(chan []string, 1)
	WatchChildren(conn, rootPath, ch)
	for {
		child := <-ch
		c := NodeArr2IntArr(child, nodeName)
		n := NodeIdToInt(nid, nodeName)
		// util.LogInfo(n, c, child)
		if n == c[0] {
			util.LogInfo(node + " get lock and acquire success.")
			break
		}
	}
}

func WatchNodePart(conn *zk.Conn, ch chan []string) (string, error) {
	cfg := config.GetConfig()
	path := cfg.GetString("zkserver.root")
	node := cfg.GetString("zkserver.nodePart")
	return WatchNode(conn, path, node, ch)
}

func WatchNodePart4Int(conn *zk.Conn, ch chan []int) (int, error) {
	cfg := config.GetConfig()
	path := cfg.GetString("zkserver.root")
	node := cfg.GetString("zkserver.nodePart")
	return WatchNode4Int(conn, path, node, ch)
}

// 监控子节点，并返回节点名称
func WatchNode(conn *zk.Conn, path string, node string, ch chan []string) (string, error) {
	rootPath := "/" + path + "/" + node
	nodeName := node + "-"
	if _, err := CreateNode(conn, rootPath, Data, Flags); err != nil {
		return "", err
	}
	nid, err := CreateSeqNode(conn, rootPath+"/"+nodeName, Data)
	if err != nil {
		return "", err
	}
	WatchChildren(conn, rootPath, ch)
	return nid, nil
}

// 监控子节点，并返回节点对应数字
func WatchNode4Int(conn *zk.Conn, path string, node string, ich chan []int) (int, error) {
	rootPath := "/" + path + "/" + node
	nodeName := node + "-"
	if _, err := CreateNode(conn, rootPath, Data, Flags); err != nil {
		return -1, err
	}

	nid, err := CreateSeqNode(conn, rootPath+"/"+nodeName, Data)
	if err != nil {
		return -1, err
	}
	n := NodeIdToInt(nid, nodeName)
	ch := make(chan []string, 1)
	WatchChildren(conn, rootPath, ch)
	go func() {
		for {
			child := <-ch
			c := NodeArr2IntArr(child, nodeName)
			// util.LogInfo(n, c)
			ich <- c
		}
	}()
	return n, nil
}
