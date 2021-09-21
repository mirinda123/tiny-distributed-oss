package objectstream

import (
	"fmt"
	"io"
	"net/http"
)

type PutStream struct{
	writer *io.PipeWriter
	c chan error
}

func NewPutStream(server, object string) *PutStream{
	//reader和write是关联起来的，从writer中写入的东西，可以从reader中读出来
	reader, writer := io.Pipe()
	c := make(chan error)

	go func() {
		request, _ := http.NewRequest("PUT", "http://"+server+"/objects/"+object, reader)
		client := http.Client{}

		//执行下一句话的时候，会阻塞，因为要等待writer传数据到reader中，再从reader中读取，通过http发送出去
		r, err := client.Do(request)
		if err == nil && r.StatusCode != http.StatusOK {
			err = fmt.Errorf("dataServer return http code %d", r.StatusCode)
		}
		c <- err
	}()


	return &PutStream{writer, c}
}

func (w *PutStream) Write(p []byte) (n int, err error) {
	return w.writer.Write(p)
}

func (w *PutStream) Close() error {
	w.writer.Close()
	return <-w.c
}