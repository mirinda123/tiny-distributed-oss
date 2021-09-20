package objects

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func Handler(w http.ResponseWriter, r *http.Request){
	m := r.Method
	if m == http.MethodPut{

		put(w, r)
		return
	}

	if m == http.MethodGet{
		get(w, r)
		return
	}

	//请求不是put也不是get，就返回错误代码405
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func get(w http.ResponseWriter, r *http.Request) {

	file, err := os.Open(os.Getenv("STORAGE_ROOT") + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	if err != nil{
		log.Println(err)
		log.Println("?")
		w.WriteHeader(http.StatusNotFound)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}(file)

	//使用Copy来将file写入w
	//file是File结构体，同时实现了Writer接口和Reader接口
	_, err = io.Copy(w, file)
	if err != nil {
		log.Println(err)
		return
	}
}

func put(w http.ResponseWriter, r *http.Request) {

	//EscapedPath获取经过转义后的路径，为/objects/<object_name>
	//在本地文件系统的根存储目录的objects子目录下创建文件file
	file, err := os.Create(os.Getenv("STORAGE_ROOT") + "\\objects\\" + strings.Split(r.URL.EscapedPath(), "/")[2])
	log.Println((os.Getenv("STORAGE_ROOT") + "\\objects\\"))
	if err != nil{
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)
	//写入文件
	//r.Body是http的正文内容
	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Println(err)
		return
	}
}

