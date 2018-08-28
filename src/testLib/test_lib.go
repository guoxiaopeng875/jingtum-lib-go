/***  初始化
 *** testLib.go
 *** 主要用于用于测试jingtumLib的各个实例
 *** author:              1416205324@qq.com
 *** last_modified_time:  2018-5-25 13:13:23
 */

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	jingtum "jingtumlib"
	"jingtumlib/constant"
)

func main() {
	err := jingtum.Init()
	if err != nil {
		fmt.Println("Init jingtum-lib error,errno", err)
		os.Exit(0)
	}

	/*
	*钱包类测试
	 */
	secret := "snsYqv2FsYLuibE9TGHdG5x5V5Qcn"
	//私钥合法性测试
	isOk := jingtum.IsValidSecret(secret)

	if !isOk {
		fmt.Printf("\nFailure IsValidSecret(%s) is false\n", secret)
	}

	fmt.Printf("\nSuccess IsValidSecret(%v) is true\n\n", secret)

	//根据私钥创建测试
	wt, err := jingtum.FromSecret(secret)

	if err != nil {
		fmt.Printf("Failure FromSecret : %s, err %v\n", secret, err)
	}

	fmt.Printf("Success FromSecret(%s). PublicKey : %s. Wallet address : %s\n\n", wt.GetSecret(), wt.GetPublicKey(), wt.GetAddress())

	//钱包地址合法性验证

	isOk = jingtum.IsValidAddress(wt.GetAddress())

	if !isOk {
		fmt.Printf("Failure IsValidAddress(%s) is false\n", wt.GetAddress())
	}

	fmt.Printf("Success IsValidAddress(%s) is true\n\n", wt.GetAddress())

	//生成新钱包
	newWallet, err := jingtum.Generate()
	isOk = jingtum.IsValidSecret(newWallet.GetSecret())
	if !isOk {
		fmt.Printf("New secret IsValidSecret(%s) is false\n", newWallet.GetSecret())
	}

	isOk = jingtum.IsValidAddress(newWallet.GetAddress())
	if !isOk {
		fmt.Printf("New address IsValidAddress(%s) is false\n", newWallet.GetAddress())
	}

	fmt.Printf("Success new secret (%s). address (%s)\n\n", newWallet.GetSecret(), newWallet.GetAddress())

	secret = "ssc5eiFivvU2otV6bSYmJeZrAsQK3"
	//根据私钥创建测试
	wt, err = jingtum.FromSecret(secret)

	if err != nil {
		fmt.Printf("Failure FromSecret : %s, err %v\n", secret, err)
	}

	fmt.Printf("Success FromSecret(%s). PublicKey : %s. Wallet address : %s\n\n", wt.GetSecret(), wt.GetPublicKey(), wt.GetAddress())

	//123.57.219.57:5020
	//139.129.194.175:5020   合约环境
	remote, err := jingtum.NewRemote("ws://123.57.219.57:5020", true)
	if err != nil {
		fmt.Printf("New remote fail : %s\n", err)
		return
	}

	cerr := remote.Connect(func(err error, result interface{}) {
		if err != nil {
			return
		}

		fmt.Printf("%s", result)
	})

	if cerr != nil {
		fmt.Printf("Connect service fail : %v\n", err)
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	//请求账号信息
	options := make(map[string]interface{})
	options["account"] = "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"
	req, err := remote.RequestAccountInfo(options)

	if err != nil {
		fmt.Printf("RequestAccountInfo fail : %v", err)
		return
	}

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			fmt.Printf("Requst account info : %v\n", err)
			wg.Done()
			return
		}

		fmt.Printf("Requst submit result : %v", result)
		wg.Done()
	})

	//支付请求
	var v struct {
		account string
		secret  string
	}
	// v.account = "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"
	// v.secret = "saNUs41BdTWSwBRqSTbkNdjnAVR8h"
	// to := "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk" //"jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"
	v.account = "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"
	v.secret = "saNUs41BdTWSwBRqSTbkNdjnAVR8h"
	to := "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"
	amount := jingtum.Amount{}
	amount.Currency = "CCT"
	amount.Value = "5"
	amount.Issuer = "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"
	tx, err := remote.BuildPaymentTx(v.account, to, amount)
	if err != nil {
		fmt.Printf("Build paymanet tx fail : %s\n", err)
		return
	}
	wg.Add(1)
	tx.SetSecret(v.secret)
	tx.AddMemo("支付5SWT")
	tx.Submit(func(err error, result interface{}) {
		if err != nil {
			fmt.Printf("Payment fail : %v\n", err)
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)

		fmt.Printf("Payment result : %s\n", jsonByte)
		wg.Done()
	})

	//请求服务信息
	req, err = remote.RequestServerInfo()
	if err != nil {
		fmt.Printf("Fail request server info %s", err.Error())
	}

	req.Submit(func(err error, result interface{}) {
		if err != nil {
			fmt.Printf("Fail request server info %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		fmt.Printf("Success request server info %s", jsonByte)
		wg.Done()
	})

	//请求市场挂单
	options = make(map[string]interface{}) //{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk"}
	gets := jingtum.Amount{}
	gets.Currency = "SWT"
	pays := jingtum.Amount{}
	pays.Currency = "CNY"
	pays.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
	options["gets"] = gets
	options["pays"] = pays
	req, err = remote.RequestOrderBook(options)
	if err != nil {
		fmt.Printf("Fail request order book %s", err.Error())
		return
	}

	if err != nil {
		fmt.Printf("Fail request order book %s", err.Error())
		return
	}

	wg.Add(1)
	req.Submit(func(err error, result interface{}) {
		if err != nil {
			fmt.Printf("Fail request order book %s", err.Error())
			wg.Done()
			return
		}

		jsonByte, _ := json.Marshal(result)
		fmt.Printf("Success request order book %s", jsonByte)
		wg.Done()
	})

	wg.Add(1)
	//账本监听
	remote.On(constant.EventLedgerClosed, func(data interface{}) {
		jsonBytes, _ := json.Marshal(data)
		fmt.Printf("Success listener ledger closed : %s", string(jsonBytes))
		wg.Done()
	})

	wg.Add(1)
	//部署合约
	options = map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "amount": float64(100), "payload": fmt.Sprintf("%X", "result={}; function Init(t) result=scGetAccountBalance(t) return result end; function foo(t) result=scGetAccountBalance(t) return result end"), "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
	tx, err = remote.DeployContractTx(options)
	if err != nil {
		fmt.Printf("Fail request deploy contract %s", err.Error())
		wg.Done()
	}
	tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
	tx.Submit(func(err error, data interface{}) {
		if err != nil {
			fmt.Printf("Fail request deploy contract %s", err.Error())
		} else {
			jsonBytes, _ := json.Marshal(data)
			fmt.Printf("Success deploy contract : %s", string(jsonBytes))
		}
		wg.Done()
	})

	//执行合约
	options = map[string]interface{}{"account": "jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h", "destination": "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72", "foo": "foo", "params": []string{"jHJJXehDxPg8HLYytVuMVvG3Z5RfhtCz7h"}}
	tx, err = remote.CallContractTx(options)
	if err != nil {
		fmt.Printf("Fail request call contract Tx %s", err.Error())
		wg.Done()
	}
	wg.Add(1)
	tx.SetSecret("saNUs41BdTWSwBRqSTbkNdjnAVR8h")
	tx.Submit(func(err error, data interface{}) {
		if err != nil {
			fmt.Printf("Fail request call contract Tx %s", err.Error())
		} else {
			jsonBytes, _ := json.Marshal(data)
			fmt.Printf("Success call contract Tx : %s", string(jsonBytes))
		}
		wg.Done()
	})

	//请求账号信息
	options = map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "trust", "quality_out": 100, "quality_in": 10}
	//options := map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "authorize", "target": "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"}
	limit := jingtum.Amount{}
	limit.Currency = "CCA"
	limit.Value = "0.0001"
	limit.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
	options["limit"] = limit
	treq, err := remote.BuildRelationTx(options)
	if err != nil {
		fmt.Printf("BuildRelationTx fail : %s", err.Error())
		return
	}

	// wg = sync.WaitGroup{}
	wg.Add(1)
	treq.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
	treq.Submit(func(err error, result interface{}) {
		if err != nil {
			fmt.Printf("Build Relation Tx : %s", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		fmt.Printf("Success Build Relation Tx result : %s", jsonBytes)
		wg.Done()
	})

	//置账号属性
	options = map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "delegate", "delegate_key": "jGXjV57AKG7dpEv8T6x5H6nmPvNK5tZj72"}
	limit = jingtum.Amount{}
	limit.Currency = "SWT"
	limit.Value = "100.0001"
	limit.Issuer = "jBciDE8Q3uJjf111VeiUNM775AMKHEbBLS"
	options["limit"] = limit
	tran, err := remote.BuildAccountSetTx(options)
	if err != nil {
		fmt.Printf("Build AccountSet Tx fail : %s", err.Error())
		return
	}
	// wg = sync.WaitGroup{}
	wg.Add(1)
	tran.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
	tran.Submit(func(err error, result interface{}) {
		if err != nil {
			fmt.Printf("Build AccountSet Tx : %s", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		fmt.Printf("Success Build AccountSet Tx result : %s", jsonBytes)
		wg.Done()
	})

	// 挂单
	options = map[string]interface{}{"account": "j3N35VHut94dD1Y9H1KoWmGZE2kNNRFcVk", "type": "property", "set_flag": "asfRequireDest", "clear": "asfDisableMaster"}
	gets = jingtum.Amount{}
	gets.Currency = "SWT"
	pays = jingtum.Amount{}
	pays.Currency = "CNY"
	options["gets"] = gets
	options["pays"] = pays
	reqt, err := remote.BuildAccountSetTx(options)
	if err != nil {
		fmt.Printf("BuildOfferCreateTx fail : %s", err.Error())
		return
	}
	// wg = sync.WaitGroup{}
	wg.Add(1)
	reqt.SetSecret("ss2QPCgioAmWoFSub4xdScnSBY7zq")
	reqt.Submit(func(err error, result interface{}) {
		if err != nil {
			fmt.Printf("Build Offer Create Tx : %s", err.Error())
			wg.Done()
			return
		}
		jsonBytes, _ := json.Marshal(result)
		fmt.Printf("Success Build Offer Create Tx result : %s", jsonBytes)
		wg.Done()
	})
	wg.Wait()

	defer jingtum.Exits()
}
