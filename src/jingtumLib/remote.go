package jingtumLib

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"jingtumLib/constant"
	jtLRU "jingtumLib/lruCache"
	"jingtumLib/utils"
)

var (
	//MaxReciveLen 接收最长报文
	MaxReciveLen = 4096000
)

//Remote 是跟井通底层交互最主要的类，它可以组装交易发送到底层、订阅事件及从底层拉取数据。
type Remote struct {
	requests  map[uint64]*ReqCtx
	LocalSign bool
	Paths     *jtLRU.LRU
	server    *Server
}

//ReqCtx 请求包装类
type ReqCtx struct {
	command  string
	data     map[string]interface{}
	callback func(err error, data interface{})
	cid      uint64
	filter   Filter
}

//Remoter 提供以下方法：
type Remoter interface {
	//Connect 连接
	Connect(callback func(err error, result interface{})) error

	//GetNowTime 获取当前时间
	GetNowTime() string

	//断开连接
	Disconnect()

	//RequestServerInfo 请求底层服务器信息
	RequestServerInfo() (*Request, error)

	//RequestLedgerClosed 获取最新账本信息
	RequestLedgerClosed() (*Request, error)

	//获取某一账本具体信息
	RequestLedger(options map[string]interface{}) (*Request, error)

	//RequestTx 询某一交易具体信息
	RequestTx(hash string) (*Request, error)

	//请求账号信息
	RequestAccountInfo(options map[string]interface{}) (*Request, error)

	//RequestAccountTums 得账号可接收和发送的货币
	RequestAccountTums(options map[string]interface{}) (*Request, error)

	//RequestAccountRelations 得账号关系
	RequestAccountRelations(options map[string]interface{}) (*Request, error)

	//RequestAccountOffers 获得账号挂单
	RequestAccountOffers(options map[string]interface{}) (*Request, error)

	//RequestAccountTx 获得账号交易列表
	RequestAccountTx(options map[string]interface{}) (*Request, error)

	//RequestOrderBook 获得市场挂单列表
	RequestOrderBook(options map[string]interface{}) (*Request, error)

	//BuildPaymentTx 创建支付对象
	BuildPaymentTx(account string, to string, amount constant.Amount) (*Transaction, error)
}

//NewRemote 创建Remote，url 为空是从配置文件获取server 地址
func NewRemote(url string, localSign bool) (*Remote, error) {
	remote := new(Remote)

	if url == "" {
		url = JTConfig.Read("Service", "Host")

		if url == "" {
			Error("Config Service:Host is null.")
			return remote, errors.New("Config|service:Host setting error")
		}

		port := JTConfig.Read("Service", "Port")

		if port == "" {
			Error("Config Service:Port is null.")
			return remote, errors.New("Config|service:Port setting error")
		}

		url += ":" + port
	}

	remote.requests = make(map[uint64]*ReqCtx)
	lru, err := jtLRU.NewLRU(100, time.Duration(5)*time.Minute, nil)
	if err != nil {
		return remote, err
	}
	remote.Paths = lru
	remote.LocalSign = localSign
	server, err := NewServer(remote, url)
	if err != nil {
		return remote, err
	}

	remote.server = server

	return remote, nil
}

//Connect 连接函数
func (remote *Remote) Connect(callback func(err error, result interface{})) error {
	if remote.server == nil {
		callback(constant.ERR_SERVER_NOT_READY, nil)
		return constant.ERR_SERVER_NOT_READY
	}

	return remote.server.connect(callback)
}

//GetNowTime 获取当前时间。格式(2006-01-02 15:04:05)
func (remote *Remote) GetNowTime() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05")
}

//Disconnect 关闭连接
func (remote *Remote) Disconnect() {
	if remote.server != nil && remote.server.Disconnect() {
		//清除请求缓存
		for id := range remote.requests {
			delete(remote.requests, id)
		}
	}
}

//RequestServerInfo 请求底层服务器信息
func (remote *Remote) RequestServerInfo() (*Request, error) {
	req := NewRequest(remote, constant.CommandServerInfo, func(data interface{}) interface{} {
		info := data.(map[string]interface{})["info"].(map[string]interface{})
		retData := map[string]interface{}{"version": "skywelld-" + info["build_version"].(string), "peers": info["peers"], "state": info["server_state"], "public_key": info["pubkey_node"], "complete_ledgers": info["complete_ledgers"], "ledger": info["validated_ledger"].(map[string]interface{})["hash"]}

		return retData
	})

	return req, nil
}

//RequestLedgerClosed 获取最新账本信息
func (remote *Remote) RequestLedgerClosed() (*Request, error) {
	req := NewRequest(remote, constant.CommandLedgerClosed, func(data interface{}) interface{} {
		retData := map[string]interface{}{"ledger_hash": data.(map[string]interface{})["ledger_hash"], "ledger_index": data.(map[string]interface{})["ledger_index"]}
		return retData
	})
	return req, nil
}

//RequestLedger 获取某一账本具体信息.
func (remote *Remote) RequestLedger(options map[string]interface{}) (*Request, error) {
	isFilter := true
	req := NewRequest(remote, constant.CommandLedger, func(data interface{}) interface{} {
		ledger, ok := data.(map[string]interface{})["ledger"]
		if !ok {
			if closed, ok := data.(map[string]interface{})["closed"]; ok {
				ledger, ok = closed.(map[string]interface{})["ledger"]
			}
		}
		if !isFilter {
			return ledger
		}

		if ledger == nil {
			return nil
		}

		retData := map[string]interface{}{"accepted": ledger.(map[string]interface{})["accepted"], "ledger_hash": ledger.(map[string]interface{})["hash"], "ledger_index": ledger.(map[string]interface{})["ledger_index"], "parent_hash": ledger.(map[string]interface{})["parent_hash"], "close_time": ledger.(map[string]interface{})["close_time_human"], "total_coins": ledger.(map[string]interface{})["total_coins"]}
		return retData
	})

	if ledgerIndex, ok := options["ledger_index"].(string); ok && utils.MatchString("^\\d+$", ledgerIndex) {
		ledgerIndexNum, err := strconv.Atoi(ledgerIndex)
		if err != nil {
			return nil, err
		}
		req.message["ledger_index"] = ledgerIndexNum
	}

	if ledgerHash, ok := options["ledger_hash"].(string); ok && utils.MatchString("^[A-F0-9]{64}$", ledgerHash) {
		req.message["ledger_hash"] = ledgerHash
	}

	if transactions, ok := options["transactions"].(bool); ok {
		req.message["transactions"] = transactions
		isFilter = false
	}

	return req, nil
}

//RequestTx 查询某一交易具体信息
func (remote *Remote) RequestTx(hash string) (*Request, error) {
	if hash == "" || !utils.MatchString("^[A-F0-9]{64}$", hash) {
		return nil, fmt.Errorf("Invalid tx hash")
	}

	req := NewRequest(remote, constant.CommandTX, nil)
	req.message["transaction"] = hash
	return req, nil
}

func getRelationType(relationType string) *constant.Integer {
	switch relationType {
	case "trustline":
		return constant.NewInteger(0)
	case "authorize":
		return constant.NewInteger(1)
	case "freeze":
		return constant.NewInteger(3)

	}
	return nil
}

func requestAccount(req *Request, options map[string]interface{}) {
	if retype, ok := options["type"].(string); ok {
		relationType := getRelationType(retype)
		if relationType != nil {
			req.message["relation_type"] = relationType.IntValue()
		}
	}

	if account, ok := options["account"].(string); ok {
		req.message["account"] = account
	}

	ledger, _ := options["ledger"]
	req.SelectLedger(ledger)

	if peer, ok := options["peer"].(string); ok {
		if utils.IsValidAddress(peer) {
			req.message["peer"] = peer
		}
	}

	if limit, ok := options["limit"].(int); ok {
		if limit < 0 {
			limit = 0
		}

		if limit > 1000000000 {
			limit = 1000000000
		}

		req.message["limit"] = limit

	}

	if marker, ok := options["marker"]; ok {
		req.message["marker"] = marker
	}
}

//RequestAccountInfo 请求账号信息
func (remote *Remote) RequestAccountInfo(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, "", nil)
	req.command = constant.CommandAccountInfo
	requestAccount(req, options)
	return req, nil
}

//RequestAccountTums 获得账号可接收和发送的货币
func (remote *Remote) RequestAccountTums(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, "", nil)
	req.command = constant.CommandAccountCurrencies
	requestAccount(req, options)
	return req, nil
}

//RequestAccountRelations 获得账号关系
func (remote *Remote) RequestAccountRelations(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, "", nil)
	rtype, ok := options["type"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid realtion type")
	}

	if _, okType := constant.RelationTypes[rtype]; !okType {
		return nil, fmt.Errorf("invalid realtion type %s", rtype)
	}

	switch rtype {
	case "trust":
		req.command = constant.CommandAccountLines
	case "authorize", "freeze":
		req.command = constant.CommandAccountRelation
	default:
		return nil, fmt.Errorf("relation should not go here %s", rtype)
	}

	requestAccount(req, options)
	return req, nil
}

//RequestAccountOffers 获得账号挂单
func (remote *Remote) RequestAccountOffers(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, "", nil)
	req.command = constant.CommandAccountOffers
	requestAccount(req, options)
	return req, nil
}

//RequestAccountTx 获得账号交易列表
func (remote *Remote) RequestAccountTx(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, constant.CommandAccountTX, func(data interface{}) interface{} {
		//過濾交易列表
		return data
	})

	if _, ok := options["limit"]; !ok {
		options["limit"] = 200
	}

	if account, ok := options["account"].(string); ok {
		if !utils.IsValidAddress(account) {
			return nil, fmt.Errorf("account parameter is invalid %s", account)
		}
		req.message["account"] = account
	}

	if ledgerMin, ok := options["ledger_min"].(int); ok {
		req.message["ledger_index_min"] = ledgerMin
	} else {
		req.message["ledger_index_min"] = 0
	}
	if ledgerMax, ok := options["ledger_max"].(int); ok {
		req.message["ledger_index_max"] = ledgerMax
	} else {
		req.message["ledger_index_max"] = -1
	}

	if limit, ok := options["limit"].(int); ok {
		req.message["limit"] = limit
	}

	if offset, ok := options["offset"].(int); ok {
		req.message["offset"] = offset
	}

	if marker, ok := options["offset"].(map[string]interface{}); ok {
		if _, ok = marker["ledger"].(int); ok {
			if _, ok = marker["seq"].(int); ok {
				req.message["marker"] = marker
			}
		}
	}
	if forward, ok := options["forward"].(bool); ok {
		//true 正向；false反向
		req.message["forward"] = forward
	}
	return req, nil
}

//RequestOrderBook 获得市场挂单列表
func (remote *Remote) RequestOrderBook(options map[string]interface{}) (*Request, error) {
	req := NewRequest(remote, constant.CommandBookOffers, nil)

	if takerGets, ok := options["taker_gets"]; ok {
		getsAmount, ok := takerGets.(constant.Amount)
		if !ok {
			return nil, fmt.Errorf("invalid taker_gets type. See also constant.Amount")
		}
		if !utils.IsValidAmount0(&getsAmount) {
			return nil, fmt.Errorf("invalid taker gets amount")
		}
		req.message["taker_gets"] = getsAmount
	} else if pays, ok := options["pays"]; ok {
		paysAmount, ok := pays.(constant.Amount)
		if !ok {
			return nil, fmt.Errorf("invalid pays type. See also constant.Amount")
		}
		if !utils.IsValidAmount0(&paysAmount) {
			return nil, fmt.Errorf("invalid taker gets amount")
		}
		req.message["taker_gets"] = paysAmount
	}

	if takerPays, ok := options["taker_pays"]; ok {
		paysAmount, ok := takerPays.(constant.Amount)
		if !ok {
			return nil, fmt.Errorf("invalid taker_pays type. See also constant.Amount")
		}
		if !utils.IsValidAmount0(&paysAmount) {
			return nil, fmt.Errorf("invalid taker pays amount")
		}
		req.message["taker_pays"] = paysAmount

	} else if gets, ok := options["gets"]; ok {
		getsAmount, ok := gets.(constant.Amount)
		if !ok {
			return nil, fmt.Errorf("invalid gets type. See also constant.Amount")
		}
		if !utils.IsValidAmount0(&getsAmount) {
			return nil, fmt.Errorf("invalid gets amount")
		}
		req.message["taker_pays"] = getsAmount
	}

	if limit, ok := options["limit"].(int); ok {
		req.message["limit"] = limit
	}

	if taker, ok := options["taker"]; ok {
		req.message["taker"] = taker
	} else {
		req.message["taker"] = constant.AccountOne
	}
	return req, nil
}

//Submit 提交请求
func (remote *Remote) Submit(command string, data map[string]interface{}, filter Filter, callback func(err error, data interface{})) {
	rc := new(ReqCtx)
	rc.command = command
	rc.data = data
	rc.callback = callback
	rc.filter = filter
	rc.cid = remote.server.GetCid()
	remote.requests[rc.cid] = rc

	remote.server.sendMessage(rc)
}

func (remote *Remote) handleResponse(data *constant.ResponseData) {
	request, ok := remote.requests[data.ID]

	if !ok {
		Errorf("Request id error %d", data.ID)
		return
	}

	delete(remote.requests, data.ID)

	if data.Status == "success" {
		result := request.filter(data.Result)
		request.callback(nil, result)
	} else if data.Status == "error" {
		errMsg := data.ErrorMessage
		if errMsg != "" {
			request.callback(errors.New(errMsg), nil)
		}
	}
}

func (remote *Remote) handlePathFind(data *constant.ResponseData) {
	//    this.emit('path_find', data);
}

func (remote *Remote) handleTransaction(data *constant.ResponseData) {
	//    var self = this;
	//    var tx = data.transaction.hash;
	//    if (self._cache.get(tx)) return;
	//    self._cache.set(tx, 1);
	//    this.emit('transactions', data);
}

func (remote *Remote) handleServerStatus(data *constant.ResponseData) {
	// TODO check data format
	//    this._updateServerStatus(data);
	//    this.emit('server_status', data);
}

func (remote *Remote) handleLedgerClosed(data *constant.ResponseData) {
	//    var self = this;
	//    if (data.ledger_index > self._status.ledger_index) {
	//        self._status.ledger_index = data.ledger_index;
	//        self._status.ledger_time = data.ledger_time;
	//        self._status.reserve_base = data.reserve_base;
	//        self._status.reserve_inc = data.reserve_inc;
	//        self._status.fee_base = data.fee_base;
	//        self._status.fee_ref = data.fee_ref;
	//        self.emit('ledger_closed', data);
	//    }
}

//消息处理方法
func (remote *Remote) handleMessage(msg []byte) {
	var data constant.ResponseData
	err := json.Unmarshal(msg, &data)
	if err != nil {
		Errorf("Received msg json Unmarshal error : %v", err)
		return
	}

	if data.Error != "" {
		// cid, err := strconv.ParseUint(data["id"].(string), 10, 64)
		// if err != nil {
		// 	Errorf("Received msg parse id error : %v", err)
		// 	return
		// }
		remote.requests[data.ID].callback(errors.New(data.ErrorMessage), nil)
	} else {
		resType := data.Type
		switch resType {
		case "ledgerClosed":
			remote.handleLedgerClosed(&data)
		case "serverStatus":
			remote.handleServerStatus(&data)
		case "response":
			remote.handleResponse(&data)
		case "transaction":
			remote.handleTransaction(&data)
		case "path_find":
			remote.handlePathFind(&data)
		}
	}
}

//BuildPaymentTx 创建支付对象
func (remote *Remote) BuildPaymentTx(account string, to string, amount constant.Amount) (*Transaction, error) {
	tx, err := NewTransaction(remote, nil)
	if err != nil {
		return nil, err
	}

	if !utils.IsValidAddress(account) {
		return nil, constant.ERR_PAYMENT_INVALID_SRC_ADDR
	}

	if !utils.IsValidAddress(to) {
		return nil, constant.ERR_PAYMENT_INVALID_DST_ADDR
	}

	if !utils.IsValidAmount(&amount) {
		return nil, constant.ERR_PAYMENT_INVALID_AMOUNT
	}

	tx.AddTxJSON("TransactionType", "Payment")
	tx.AddTxJSON("Account", account)

	toamount, err := utils.ToAmount(amount)

	if err != nil {
		return nil, err
	}

	tx.AddTxJSON("Amount", toamount)
	tx.AddTxJSON("Destination", to)

	return tx, nil
}
