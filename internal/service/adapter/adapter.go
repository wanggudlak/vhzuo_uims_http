package adapter

import (
	"bytes"
	"github.com/json-iterator/go"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
)

const BufferSize = 4096

type Adapter struct {
	pool sync.Pool
}

func New() *Adapter {
	return &Adapter{
		pool: sync.Pool{
			New: func() interface{} {
				return bytes.NewBuffer(make([]byte, BufferSize))
			},
		},
	}
}

//func (adp *Adapter) ReadRequestBody(r *http.Request, req interface{}) error {
//	method := r.Method
//	if "GET" == method {
//		return nil
//	}
//	buffer := adp.pool.Get().(*bytes.Buffer)
//	buffer.Reset()
//	defer func() {
//		if buffer != nil {
//			adp.pool.Put(buffer)
//			buffer = nil
//		}
//	}()
//
//	_, err := io.Copy(buffer, r.Body)
//	if err != nil {
//		return err
//	}
//	// 为了让body被读取一次后不要关闭，以供gin框架的绑定是再次读取并关闭
//	r.Body = ioutil.NopCloser(bytes.NewBuffer(buffer.Bytes()))
//
//	if err = jsoniter.Unmarshal(buffer.Bytes(), req); err != nil {
//		return err
//	}
//	adp.pool.Put(buffer)
//	buffer = nil
//	return nil
//}

func (adp *Adapter) ReadRequestBody(r *http.Request, req interface{}) error {
	method := r.Method
	if "GET" == method {
		return nil
	}
	buffer := adp.pool.Get().(*bytes.Buffer)
	buffer.Reset()
	defer func() {
		if buffer != nil {
			adp.pool.Put(buffer)
			buffer = nil
		}
	}()

	_, err := io.Copy(buffer, r.Body)
	if err != nil {
		return err
	}

	// 为了让body被读取一次后不要关闭，以供gin框架的绑定是再次读取然后由http.go server.go中的finishRequest中关闭请求
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buffer.Bytes()))

	if err = jsoniter.Unmarshal(buffer.Bytes(), req); err != nil {
		return err
	}
	adp.pool.Put(buffer)
	buffer = nil
	return nil
}
