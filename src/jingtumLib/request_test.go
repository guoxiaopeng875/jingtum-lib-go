/**
 * 请求测试类
 *
 * @FileName: request_test.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-26 10:44:32
 * @UpdateTime: 2018-07-26 10:44:54
 */
package jingtumLib

import (
	"encoding/json"
	"sync"
	"testing"
)

//Test_RequestLedger 获取某一账本
func Test_RequestLedger(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	options := map[string]interface{}{"transactions": true, "ledger_index": 969054, "ledger_hash": "AEE4B16B543D8C8924F09C1DB822C6419780B86019F5F5FF8DC2938E7E0E89D2"}

	req, err := remote.RequestLedger(options)
	if err != nil {
		t.Fatalf("Fail request ledger closed %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request ledger %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request ledger %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

// Test_RequestLedgerClosed 获取最新账本
func Test_RequestLedgerClosed(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	req, err := remote.RequestLedgerClosed()
	if err != nil {
		t.Fatalf("Fail request ledger closed %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request ledger closed %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request ledger closed %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestServerInfo 获取服务器信息
func Test_RequestServerInfo(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	req, err := remote.RequestServerInfo()
	if err != nil {
		t.Fatalf("Fail request server info %s", err.Error())
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Fail request server info %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		t.Logf("Success request server info %s", jsonByte)
		wg.Done()
	})

	wg.Wait()
}

//Test_RequestAccountInfo 账号信息测试
func Test_RequestAccountInfo(t *testing.T) {
	remote, err := NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		t.Fatalf("New remote fail : %s", err.Error())
		return
	}

	defer remote.Disconnect()

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("New remote fail : %s", err.Error())
			return
		}

		jsonBytes, _ := json.Marshal(result)

		t.Logf("Connect success : %s", jsonBytes)
	})

	if cerr != nil {
		t.Fatalf("Connect service fail : %s", err.Error())
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	//请求账号信息
	options := make(map[string]interface{})
	options["account"] = "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"
	req, err := remote.RequestAccountInfo(options)

	if err != nil {
		t.Fatalf("RequestAccountInfo fail : %s", err.Error())
		wg.Done()
		return
	}

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			t.Fatalf("Requst account info : %s", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		t.Logf("Requst submit result : %s", jsonBytes)
		wg.Done()
	})

	wg.Wait()
}
