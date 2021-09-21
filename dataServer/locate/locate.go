package locate

import (
	"github.com/mirinda123/tiny-distributed-oss/lib/rabbitmq"
	"os"
	"strconv"
)

func Locate(name string) bool{

	//os.Stat来访问磁盘上对应的文件
	_, err := os.Stat(name)
	return !os.IsNotExist(err)
}

func StartLocate(){
	q := rabbitmq.New(os.Getenv("RABBITMQ_SERVER"))
	defer q.Close()
	q.Bind("dataServers")

	//c是一个channel
	c := q.Consume()

	for msg := range c{
		//消息正文经过JSON编码，对象名字上有一对双引号，使用Unquote来去掉字符串前后的双引号
		//Body应该是json中的Body:xxxx
		object, err := strconv.Unquote(string(msg.Body))
		if err != nil {
			panic(err)
		}
		if Locate(os.Getenv("STORAGE_ROOT") + "\\objects\\" + object){
			q.Send(msg.ReplyTo, os.Getenv("LISTEN_ADDRESS"))
		}
	}
}