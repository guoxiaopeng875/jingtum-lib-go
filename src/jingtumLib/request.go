/**
 * Request类主管GET请求，包括获得服务器、账号、挂单、路径等信息。请求时不需要提供密
 * 钥，且对所有用户公开。所有的请求是异步的，会提供一个回调函数。每个回调函数有两个参
 * 数，一个是错误，另一个是结果。
 *
 * @FileName: request.go
 * @Auther : 13851485286
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-05-28 10:44:32
 * @UpdateTime: 2018-05-28 10:44:54
 */

package jingtumLib

import (
	_ "errors"

	"jingtumLib/constant"
	"jingtumLib/utils"

	_ "common/github.com/blog4go"
)

type Filter func(interface{}) interface{}

type Request struct {
	remote  *Remote
	message map[string]interface{}
	command string
}

func NewRequest(remote *Remote) *Request {
	request := new(Request)
	request.remote = remote
	request.message = make(map[string]interface{})
	return request
}

//提交请求
func (req *Request) Submit(callback func(err error, data interface{})) {
	if err, ok := req.message[constant.TXJSON_ERROR_KEY].(error); ok {
		callback(err, nil)
		return
	}

	req.remote.Submit(req.command, req.message, nil, callback)
}

//SelectLedger 选择账本
func (req *Request) SelectLedger(ledger interface{}) {

	if ledger == nil {
		req.message["ledger_index"] = "validated"
		return
	}
	switch ledger.(type) {
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		req.message["ledger_index"] = ledger
	case string:
		_, ok := constant.LedgerStates[ledger.(string)]

		if ok {
			req.message["ledger_index"] = ledger
		} else if utils.MatchString("^[A-F0-9]+$", ledger.(string)) {
			req.message["ledger_hash"] = ledger.(string)
		} else {
			req.message["ledger_index"] = "validated"
		}
	default:
		req.message["ledger_index"] = "validated"
	}
}
