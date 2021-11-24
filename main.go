package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/juanefec/geth/erc20"
)

const (
	USDT = "USDT"
	UBI  = "UBI"
	ETH  = "ETH"
)

var contracts = map[string]string{
	USDT: "0xdac17f958d2ee523a2206206994597c13d831ec7",
	UBI:  "0xdd1ad9a21ce722c151a836373babe42c868ce9a4",
}

var infuraProject string
var account string

func main() {
	flag.StringVar(&infuraProject, "p", "", "usage: -p b98dac6bacd66b879..")
	flag.StringVar(&account, "a", "0x5871db3E25C5BCaD9EAa08cD1d022b4F150AEE34", "usage: -p 0xb98dac6bacd66b879..")
	flag.Parse()

	if infuraProject == "" {
		log.Fatal("no infura project id")
	}

	eth, err := getBalance(account, "")
	if err != nil {
		log.Fatal(err)
	}

	usdt, err := getBalance(account, contracts[USDT])
	if err != nil {
		log.Fatal(err)
	}

	ubi, err := getBalance(account, contracts[UBI])
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%v\n%v\n%v\n", eth, usdt, ubi)

}

func getBalance(addr, token string) (string, error) {
	if addr == "" {
		return "", errors.New("address can't be empty")
	}
	account := common.HexToAddress(addr)
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/" + infuraProject)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	if token == "" {
		balance, err := client.BalanceAt(context.Background(), account, nil)
		if err != nil {
			log.Fatal(err)
		}
		fbalance := new(big.Float)
		fbalance.SetString(balance.String())
		ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
		return fmt.Sprintf("%v %v", ETH, ethValue), nil
	}

	contract := common.HexToAddress(token)
	instance, err := erc20.NewToken(contract, client)
	if err != nil {
		log.Fatal(err)
	}

	bal, err := instance.BalanceOf(&bind.CallOpts{}, account)
	if err != nil {
		log.Fatal(err)
	}

	symbol, err := instance.Symbol(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	decimals, err := instance.Decimals(&bind.CallOpts{})
	if err != nil {
		log.Fatal(err)
	}

	fbal := new(big.Float)
	fbal.SetString(bal.String())
	value := new(big.Float).Quo(fbal, big.NewFloat(math.Pow10(int(decimals))))

	return fmt.Sprintf("%v %v", symbol, value), nil
}
