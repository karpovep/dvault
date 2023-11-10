package main

import (
	"encoding/base64"
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"log"
	"strings"

	"github.com/miguelmota/go-ethereum-hdwallet"
)

func main() {
	mnemonic := "tag volcano eight thank tide danger coast health above argue embrace heavy"
	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatal(err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, true)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(account.Address.Hex()) // 0xC49926C4124cEe1cbA0Ea94Ea31a6c12318df947

	///////////////
	data := []byte(strings.ReplaceAll("{\"userId\":\"BGAFyGpnGPZiIXE6dwc8QSkcw6u_zQOqSVXpsrUNv3-bZnLa0NRq3mHjgveYiKc-p4mdlBm-zx1snsIIfBGI-hg\",\"timestamp\":100500}", "\\", ""))
	hash := crypto.Keccak256Hash(data)
	signature, err := wallet.SignHash(account, hash.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	publicKeyBytes, err := wallet.PublicKeyBytes(account)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("publicKeyBytes", publicKeyBytes)
	fmt.Println("string(publicKeyBytes)", string(publicKeyBytes))
	fmt.Println("base64 publicKeyBytes", base64.RawURLEncoding.EncodeToString(publicKeyBytes))

	//fmt.Println(hash.Hex())
	//sigPublicKey, err := crypto.Ecrecover(hash.Bytes(), signature)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//fmt.Println(sigPublicKey)
	//matches := bytes.Equal(sigPublicKey, publicKeyBytes)
	//fmt.Println(matches) // true

	//sigPublicKeyECDSA, err := crypto.SigToPub(hash.Bytes(), signature)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//sigPublicKeyBytes := crypto.FromECDSAPub(sigPublicKeyECDSA)
	//fmt.Println(sigPublicKeyBytes)
	//matches = bytes.Equal(sigPublicKeyBytes, publicKeyBytes)
	//fmt.Println(matches) // true
	//
	signatureNoRecoverID := signature[:len(signature)-1] // remove recovery ID
	verified := crypto.VerifySignature(publicKeyBytes, hash.Bytes(), signatureNoRecoverID)
	fmt.Println(verified) // true

	//GENERATE JWT
	base64EncodedSignature := base64.RawURLEncoding.EncodeToString(signature)
	fmt.Println("signature: ", base64EncodedSignature)
	base64EncodedHeader := base64.RawURLEncoding.EncodeToString([]byte("{\"alg\":\"HS256\"}"))
	base64EncodedPayoad := base64.RawURLEncoding.EncodeToString(data)
	fmt.Println("'----------------- JWT -----------------------'")
	fmt.Println(base64EncodedHeader + "." + base64EncodedPayoad + "." + base64EncodedSignature)
}
