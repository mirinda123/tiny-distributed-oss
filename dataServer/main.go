package main

import (
	"github.com/mirinda123/tiny-distributed-oss/dataServer/heartbeat"
	"github.com/mirinda123/tiny-distributed-oss/dataServer/locate"
	"github.com/mirinda123/tiny-distributed-oss/objects"
	"log"
	"net/http"
	"os"
)


func main() {


	//注意接口服务的object,heartbeat,locate这三个包和数据服务的三个包是不一样的，具体的见书

	go heartbeat.StartHeartbeat()
	go locate.StartLocate()


	//lujing := os.Getenv("STORAGE_ROOT") + "\\objects\\" + "123.txt"
	//os.Create(lujing)

	//访问objects的时候，交给Handler处理
	http.HandleFunc("/objects/",objects.Handler)

	//当非服务器非正常退出，log.Fatal会打印错误程序并退出
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"),nil))
}

