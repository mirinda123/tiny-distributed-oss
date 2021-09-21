package objects

import (
	"fmt"
	"github.com/mirinda123/tiny-distributed-oss/apiServer/heartbeat"
	"github.com/mirinda123/tiny-distributed-oss/apiServer/locate"
	"github.com/mirinda123/tiny-distributed-oss/lib/objectstream"
	"io"
	"log"
	"net/http"
	"strings"
)

//接口服务的put和get负责将http请求转发给数据服务层
func Handler(w http.ResponseWriter, r *http.Request) {
	m := r.Method
	if m == http.MethodPut {
		put(w, r)
		return
	}
	if m == http.MethodGet {
		get(w, r)
		return
	}
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func put(w http.ResponseWriter, r *http.Request){
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	c, err := storeObject(r.Body, object)
	if err != nil{
		log.Println(err)
	}
	w.WriteHeader(c)
}

func storeObject(r io.Reader, object string) (int, error) {
	//以object为参数生成stream
	//stream是一个指向objectstream.PutStream结构体的指针
	//该结构体实现了Write方法，所以是io.Write接口
	//objectstream包是对http包的封装，来把http函数的调用转换成读写流的形式
	stream, err := putStream(object)
	if err != nil{
		return http.StatusServiceUnavailable, err
	}

	io.Copy(stream, r)

	//Close方法返回的error，用来通知调用者在数据传输的时候发生的错误
	err = stream.Close()
	if err != nil{
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil

}

//生成一个流
func putStream(object string) (*objectstream.PutStream, error){
	server := heartbeat.ChooseRandomDataServer()
	if server == ""{
		return nil, fmt.Errorf("cannot find any data server")
	}

	return objectstream.NewPutStream(server, object), nil
}


func get(w http.ResponseWriter, r *http.Request) {
	object := strings.Split(r.URL.EscapedPath(), "/")[2]
	//返回一个Getstream结构体
	stream, e := getStream(object)
	if e != nil {
		log.Println(e)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w, stream)
}

func getStream(object string) (io.Reader, error) {
	server := locate.Locate(object)
	if server == "" {
		return nil, fmt.Errorf("object %s locate fail", object)
	}
	return objectstream.NewGetStream(server, object)
}