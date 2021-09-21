package locate

import (
	"encoding/json"
	"github.com/mirinda123/tiny-distributed-oss/lib/rabbitmq"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)
//locate成功，则返回拥有该对象的一个数据服务节点的地址
func Handler(w http.ResponseWriter, r *http.Request){
	m := r.Method
	if m != http.MethodGet{
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	info := Locate(strings.Split(r.URL.EscapedPath(),"/")[2])
	if len(info) == 0{
		w.WriteHeader(http.StatusNotFound)
		return
	}

	b, err := json.Marshal(info)
	if err != nil{
		panic(err)
	}

	w.Write(b)

}
//name是需要定位的对象的名字
func Locate(name string) string{
	//创建一个新的临时消息队列
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))

	//向dataServers exchange群发这个name
	q.Publish("dataServers", name)
	c := q.Consume()
	//临时的消息队列只有1秒
	go func(){
		time.Sleep(time.Second)
		q.Close()
	}()
	msg := <-c
	s, _ := strconv.Unquote(string(msg.Body))
	return s
}

//如果1s内没有任何反馈，则消息队列关闭后，应该返回一个空字符串
func Exist(name string) bool {
	return Locate(name) != ""
}