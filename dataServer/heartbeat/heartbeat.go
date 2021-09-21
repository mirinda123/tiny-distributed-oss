package heartbeat

import (
	"github.com/mirinda123/tiny-distributed-oss/lib/rabbitmq"
	"os"
	"time"
)
func StartHeartbeat(){
	//新建一个queue
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	for {
		//在无限循环中通过publish方法向apiServers发送本节点的监听地址
		q.Publish("apiServers",os.Getenv("LISTEN_ADDRESS"))
		time.Sleep(5 * time.Second)
	}
}