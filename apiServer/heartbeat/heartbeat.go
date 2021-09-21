package heartbeat

import (
	"github.com/mirinda123/tiny-distributed-oss/lib/rabbitmq"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

var dataServers = make(map[string]time.Time)
var mutex sync.Mutex

//创建消息队列来绑定apiServers 交换机
func ListenHeartbeat(){
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()

	q.Bind("apiServer")
	c := q.Consume()

	go removeExpiredDataServer()

	//接收心跳信息
	for msg := range c{
		dataServer, err := strconv.Unquote(string(msg.Body))
		if err != nil{
			panic(err)
		}
		mutex.Lock()
		//将消息的正文，即数据服务节点的监听地址，作为map的键
		//收到消息的时间作为值
		dataServers[dataServer] = time.Now()
		mutex.Unlock()
	}
}
//每隔5s扫描一遍dataServers，清除其中超过10s没收到心跳信息的数据服务节点
func removeExpiredDataServer() {
	for{
		time.Sleep(5 * time.Second)
		mutex.Lock()
		for s, t := range dataServers{
			if t.Add(10 * time.Second).Before(time.Now()){
				delete(dataServers, s)
			}
		}
		mutex.Unlock()
	}
}
//遍历dataServers 并返回当前所有的数据服务节点
func GetDataServers() []string{
	mutex.Lock()
	defer mutex.Unlock()
	ds := make([]string, 0)

	for s, _ := range dataServers{
		ds = append(ds, s)
	}
	return ds
}

//在档期啊你的所有数据服务节点中随机选出一个并返回
func ChooseRandomDataServer() string{
	ds := GetDataServers()
	n := len(ds)
	if n == 0{
		return ""
	}
	return ds[rand.Intn(n)]
}