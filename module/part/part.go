package part

import (
	"math/rand"
	"sort"
	"strconv"
	"sync"
	"time"
	"xcurrency-push/config"
	"xcurrency-push/module/redis"
	"xcurrency-push/module/zookeeper"
	"xcurrency-push/util"
)

const (
	ALARMID_PREFIX = "alarm:ids:"
	ALARM_SERVER   = "alarm:server"
)

var PartArr []int
var NextIndex int
var ServerCount map[string]int64 = map[string]int64{}
var addLock *sync.Mutex = new(sync.Mutex)
var indexLock *sync.Mutex = new(sync.Mutex)

/**
* @desc 服务器分片计算
*       把多个分片的alarm:ids（如alarm:ids:0, alarm:ids:1, alarm:ids:2）分配到正在运行的多个服务器中，使多个服务器平分alarm:ids的索引值
* @param [servers] Redis中存储的启动过的服务器信息，每启动一台服务器，都会记录一条信息（如 CKServer0， CKServer1）
* @param [zkNodeNums] 正在运行的服务器zookeeper节点数值（如 0，1，2）
* @param [PartArr] 分片计算后，当前服务器被分配的alarm:ids的索引值（如0，2）
* @author 于朝鹏
 */
func ServerPart() {
	conn := zkclient.GetConn()
	if conn == nil {
		return
	}
	defer conn.Close()

	ch := make(chan []int, 1)
	if n, err := zkclient.WatchNodePart4Int(conn, ch); err == nil {
		if client, err := redis.GetClient(); err == nil {
			cfg := config.GetConfig()
			serverName := cfg.GetString("check.serverName")
			if isMem, err := client.SIsMember(ALARM_SERVER, serverName).Result(); err == nil {
				if !isMem { // 首次运行添加服务器节点
					client.SAdd(ALARM_SERVER, serverName)
				}
			}

			for {
				zkNodeNums := <-ch
				if servers, err := client.SMembers(ALARM_SERVER).Result(); err == nil {
					sort.Strings(servers)
					PartArr := PartArr[0:0]
					NextIndex = 0

					for i := 0; i < len(servers); i++ {
						if n == NextNodeNum(zkNodeNums) {
							PartArr = append(PartArr, i)
						}
					}
					util.LogInfo("node:", n, zkNodeNums, "server:", serverName, servers, "part:", PartArr)
				}
			}
		}
	}
}

/**
* @desc 按顺序依次返回对应的zookeeper节点数值，并循环进行（如 0，1，2，0，1，2，0，1，2）
* @param [zkNodeNums] 正在运行的服务器zookeeper节点数值（如 0，1，2）
* @author 于朝鹏
 */
func NextNodeNum(zkNodeNums []int) int {
	l := len(zkNodeNums)
	if NextIndex >= l {
		NextIndex = 0
	}
	n := zkNodeNums[NextIndex]
	NextIndex += 1
	return n
}

/**
* @desc 需要加到哪台服务器的索引计算
*       根据Redis存储的服务器数量，计算每个服务器alarm:ids存储的数量，数量越小，则优先被添加数据
* @param [servers] Redis中存储的启动过的服务器信息，每启动一台服务器，都会记录一条信息 如（CKServer0，CKServer1）
* @author 于朝鹏
 */
func ServerIndex(servers []string) int {
	indexLock.Lock()
	defer indexLock.Unlock()

	s := 0
	if client, err := redis.GetClient(); err == nil {
		if r, err := client.SCard(ALARM_SERVER).Result(); err == nil {
			s = int(r)
		}
	}
	if s <= 1 {
		return 0
	}
	for i := (s - 1); i >= 0; i-- {
		if (i - 1) < 0 {
			break
		}
		if GetCount(i) < GetCount(i-1) {
			return i
		}
	}
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(s)
}

func GetCount(i int) int64 {
	var count int64
	if v, ok := ServerCount[ALARMID_PREFIX+strconv.Itoa(i)]; ok {
		count = v
	} else {
		if client, err := redis.GetClient(); err == nil {
			count, _ = client.SCard(ALARMID_PREFIX + strconv.Itoa(i)).Result()
			ServerCount[ALARMID_PREFIX+strconv.Itoa(i)] = count
		}
	}
	return count
}

func AddCount(i int) {
	addLock.Lock()
	defer addLock.Unlock()
	var count int64
	if v, ok := ServerCount[ALARMID_PREFIX+strconv.Itoa(i)]; ok {
		count = v
	} else {
		if client, err := redis.GetClient(); err == nil {
			count, _ = client.SCard(ALARMID_PREFIX + strconv.Itoa(i)).Result()
		}
	}
	ServerCount[ALARMID_PREFIX+strconv.Itoa(i)] = count + 1
}
