/***  初始化
 *** testLib.go
 *** 主要用于用于测试jingtumLib的各个实例
 *** author:              1416205324@qq.com
 *** last_modified_time:  2018-5-25 13:13:23
 */

package main

import (
	"fmt"
	jingtum "jingtumLib"
	"os"
)

func main() {
	err := jingtum.Init()
	if err != nil {
		fmt.Println("Init jingtum-lib error,errno", err)
		os.Exit(0)
	}

	/*
		isNumber := jingtum.Number("5445")
		fmt.Println(isNumber)

		jingtum.Generate()
	*/
	_, remote := jingtum.NewRemote()
	err = remote.Connect()
	if err != nil {
		fmt.Println("Connect service", remote.Wsconn.Host, remote.Wsconn.Port, "fail.", err)
		return
	}
	fmt.Println("Connect service", remote.Wsconn.Host, remote.Wsconn.Port, "succ.")

	//请求底层服务器信息
	err, response := remote.RequestServerInfo()
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Server Info succ.len=", len(response), "data=", response)

	//获取最新账本信息
	err, response = remote.RequestLedgerClosed()
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Ledger Closed succ.len=", len(response), "data=", response)

	//获取某一账本具体信息
	var ledger_index string = "8488670"
	var ledger_hash string = ""
	var transactions bool = false
	err, response = remote.RequestLedger(ledger_index, ledger_hash, transactions)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Ledger succ.len=", len(response), "data=", response)

	//获取某一账本具体信息
	var hash string = "084C7823C318B8921A362E39C67A6FB15ADA5BCCD0C7E9A3B13485B1EF2A4313"
	err, response = remote.RequestTx(hash)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Tx succ.len=", len(response), "data=", response)

	//请求账号信息
	account := "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	err, response = remote.RequestAccountInfo(account)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Account Info succ.len=", len(response), "data=", response)

	//获得账号可接收和发送的货币
	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	err, response = remote.RequestAccountTums(account)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Account Tums succ.len=", len(response), "data=", response)

	//获得账号交易列表
	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	var limit int = 100
	err, response = remote.RequestAccountTx(account, limit)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Account Tx succ.len=", len(response), "data=", response)

	//获得账号交易列表
	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	atype := "trust"
	err, response = remote.RequestAccountRelations(account, atype)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Account Relations succ.len=", len(response), "data=", response)

	atype = "authorize"
	err, response = remote.RequestAccountRelations(account, atype)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Account Relations succ.len=", len(response), "data=", response)

	//获得账号挂单
	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	err, response = remote.RequestAccountOffers(account)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Account Offers succ.len=", len(response), "data=", response)

	//获得账号挂单
	account = "jD86doF9mBbAfTgK62L6mpqg4YJ1Yhm5wq"
	gets := "SWT"
	pays := "CNY"
	err, response = remote.RequestOrderBook(account, gets, pays)
	if err != nil {
		fmt.Println("Get data:", response)
		return
	}
	fmt.Println("Get Response Account Order Book succ.len=", len(response), "data=", response)

	/*
		//获得账号挂单
		err, response = remote.BuildPaymentTx()
		if err != nil {
			fmt.Println("Get data:", response)
			return
		}
		fmt.Println("Get Response Build Payment Tx succ.len=", len(response), "data=", response)
	*/
	defer jingtum.Exits()
}
