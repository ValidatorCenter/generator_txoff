package main

import (
	"encoding/json"
	"fmt"
	"os"

	m "github.com/ValidatorCenter/minter-go-sdk"
)

type ConfigAuth struct {
	MxAddress string `json:"mxaddress"`
	PrivKey   string `json:"privkey"`
	TxtError  string `json:"error"`
}

// Структура для файла "start.json"
type ConfigStart struct {
	NodeAPI    string `json:"url_api"`
	SeedPhr    string `json:"seed"`
	PubKeyNode string `json:"validator"`
	IsMainnet  bool   `json:"mainnet"`
}

// загрузка файла start.json
func loadStartJSON() *ConfigStart {
	file, _ := os.Open("start.json")
	decoder := json.NewDecoder(file)
	cfgStr := new(ConfigStart)
	err := decoder.Decode(&cfgStr)
	if err != nil {
		fmt.Println("ERR:", err)
		return cfgStr
	}
	return cfgStr
}

// загрузка файла gen_txoff.json
func loadAuthJSON() *ConfigAuth {
	file, _ := os.Open("gen_txoff.json")
	decoder := json.NewDecoder(file)
	cfgAuth := new(ConfigAuth)
	err := decoder.Decode(&cfgAuth)
	if err != nil {
		fmt.Println("ERR:", err)
		return cfgAuth
	}
	return cfgAuth
}

func main() {
	var err error

	cfg := loadStartJSON()
	cgn := loadAuthJSON()

	if cfg.PubKeyNode == "" {
		fmt.Println("ERROR: start.json - not pubkey validator")
		return
	}

	if cgn.TxtError != "" {
		fmt.Println("ERROR:", cgn.TxtError)
		return
	}

	if cgn.PrivKey == "" {
		fmt.Println("ERROR: gen_txoff.json - not privkey wallet address")
		return
	}

	sdk := m.SDK{
		MnAddress:     cfg.NodeAPI, // Есть необходимость получения Nonce в TxSetCandidateRLP()
		AccPrivateKey: cgn.PrivKey,
		ChainMainnet:  cfg.IsMainnet,
		AccAddress:    cgn.MxAddress,
	}

	//sdk.AccAddress нужен для получения Nonce
	if sdk.AccAddress == "" {
		sdk.AccAddress, err = m.GetAddressPrivateKey(sdk.AccPrivateKey)
		if err != nil {
			fmt.Println("ERROR: ", err)
			return
		}
	}

	// Separate Minter wallet (masternode owner) w/o output transactions is a must!

	CoinMinter := ""

	if sdk.ChainMainnet {
		CoinMinter = "BIP"
	} else {
		CoinMinter = "MNT"
	}

	Gas := int64(10) // Комиссия: 0.1bip, но мы готовы заплатить в 10 раз больше -> 1bip, но что бы транза прошла!

	sndDt := m.TxSetCandidateData{
		PubKey:   cfg.PubKeyNode,
		Activate: false, //true-"on", false-"off"
		GasCoin:  CoinMinter,
		GasPrice: Gas,
	}

	strRLP, err := sdk.TxSetCandidateRLP(&sndDt)
	if err != nil {
		panic(err)
	}

	fmt.Println("CoinMinter:", CoinMinter)

	fmt.Println("TX RLP:", strRLP)
}
