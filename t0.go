package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	// для авторизации/регистрации
	"github.com/miguelmota/go-ethereum-hdwallet"
)

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

// Авторизация по Seed-фразе
func AuthMnemonic(seedPhr string) (string, string, error) {
	wallet, err := hdwallet.NewFromMnemonic(seedPhr)
	if err != nil {
		//panic(err)
		return "", "", err
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return "", "", err
	}

	//M+`в нижнем регистре(без видущего нуля)`
	strAdrs := account.Address.String()                                    // 0x512B699Ab21542B8875609593e845818f301903B
	addrss := fmt.Sprintf("M%s", strings.ToLower(strAdrs[1:len(strAdrs)])) // Mx512b699ab21542b8875609593e845818f301903b
	privKeyStr, err := wallet.PrivateKeyHex(account)
	if err != nil {
		return "", "", err
	}
	return addrss, privKeyStr, nil
}

func main() {
	cfg := loadStartJSON()

	if cfg.SeedPhr == "" {
		fmt.Printf("{\"mxaddress\":\"%s\",\"privkey\":\"%s\",\"error\":\"%s\"}\n", "", "", "no seed")
		return
	}

	mxAddrs, PrivKeyWallet, err := AuthMnemonic(cfg.SeedPhr)
	if err != nil {
		fmt.Printf("{\"mxaddress\":\"%s\",\"privkey\":\"%s\",\"error\":\"%s\"}\n", "", "", err.Error())
		return
	}
	fmt.Printf("{\"mxaddress\":\"%s\",\"privkey\":\"%s\",\"error\":\"%s\"}\n", mxAddrs, PrivKeyWallet, "")
}
