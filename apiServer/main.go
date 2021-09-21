package apiServer

import (
	"github.com/mirinda123/tiny-distributed-oss/apiServer/heartbeat"
	"github.com/mirinda123/tiny-distributed-oss/apiServer/objects"
	"github.com/mirinda123/tiny-distributed-oss/apiServer/locate"
	"net/http"
)

func main(){
	//注意接口服务的object,heartbeat,locate这三个包和数据服务的三个包是不一样的，具体的见书
	go heartbeat.ListenHeartbeat()
	http.HandleFunc("/objects/", objects.Handler)
	http.HandleFunc("/loacate/", locate.Handler)
}