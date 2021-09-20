package main

import (
	"github.com/mirinda123/tiny-distributed-oss/objects"
	"log"
	"net/http"
	"os"
)

func main() {
	os.Setenv("STORAGE_ROOT", "D:")
	os.Setenv("LISTEN_ADDRESS",":12345")
	//访问objects的时候，交给Handler处理
	http.HandleFunc("/objects",objects.Handler)

	//当非服务器非正常退出，log.Fatal会打印错误程序并退出
	log.Fatal(http.ListenAndServe(os.Getenv("LISTEN_ADDRESS"),nil))
}

