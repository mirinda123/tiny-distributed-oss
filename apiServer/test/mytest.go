package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type PutStream struct{
	writer *io.PipeWriter
	c chan error
}
func main(){
	stream := NewPutStream()
	data := url.Values{}
	data.Set("name", "foo")
	data.Set("surname", "bar")

	//write, err := stream.writer.Write([]byte(data.Encode()))
	//if err != nil {
	//	fmt.Println(err)
	//}
	//fmt.Println(write)
	time.Sleep(10 * time.Second)
}


func NewPutStream() *PutStream{
	//reader和write是关联起来的，从writer中写入的东西，可以从reader中读出来
	reader, writer := io.Pipe()
	c := make(chan error)

	go func() {
		buffer := make([]byte, 100)
		read, err2 := reader.Read(buffer)
		if err2 != nil{
			fmt.Println(err2)
		}
		fmt.Println(read)
		fmt.Println(1)
		request, _ := http.NewRequest(http.MethodPost, "https://www.oschina.net/", reader)

		client := http.Client{}
		r, err := client.Do(request)

		if err == nil && r.StatusCode != http.StatusOK {
			err = fmt.Errorf("dataServer return http code %d", r.StatusCode)
		}

		fmt.Println(r)
		res,_ :=http.Get("https://www.baidu.com/")
		fmt.Println(res)
	}()


	return &PutStream{writer, c}
}