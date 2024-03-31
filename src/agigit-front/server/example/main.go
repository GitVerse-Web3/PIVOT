package main

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/coming-chat/go-aptos/aptosaccount"
	"github.com/coming-chat/go-aptos/aptosclient"
	"github.com/coming-chat/go-aptos/aptostypes"
)

const rpcUrl = "https://fullnode.devnet.aptoslabs.com"

func main() {
	// ctx := context.Background()
	// client, err := aptosclient.Dial(ctx, rpcUrl)
	// if err != nil {
	// 	printError(err)
	// }
	// ledgerInfo, err := client.LedgerInfo()
	// if err != nil {
	// 	panic(err)
	// }
	// content, err := json.Marshal(ledgerInfo)
	// if err != nil {
	// 	printError(err)
	// }
	// fmt.Println(string(content))

	// transactions, err := client.GetTransactions(1, 10)
	// if err != nil {
	// 	printError(err)
	// }
	// printLine("get tx list")
	// for _, tx := range transactions {
	// 	fmt.Printf("version: %d, type: %s, hash: %s, round: %d, time: %d\n", tx.Version, tx.Type, tx.Hash, tx.Round, tx.Timestamp)
	// }

	// events, err := client.GetEventsByCreationNumber("0x8ab4fe552d76c49cdf7a4707338520ddf5de705044665fbcc4b6ea60a6b026d4", "2", 0, 10)
	// if err != nil {
	// 	printError(err)
	// }
	// for _, e := range events {
	// 	fmt.Println(e.SequenceNumber)
	// }
	// Import account with private key
	account, _ := aptosaccount.NewAccountWithMnemonic("candy maple cake sugar pudding cream honey rich smooth crumble sweet treat")

	//privateKey, _ := hex.DecodeString("0xE5951E406CC872770BA6203FAF04736CEDF1FA3952E723C9513AEC0C5A0A9FB4")
	//account := aptosaccount.NewAccount(privateKey)
	fromAddress := "0x" + hex.EncodeToString(account.AuthKey[:])
	fmt.Printf("from address = %s\n", fromAddress)
	toAddress := "0x2c3c7c96c148777979c1964b555f1c559b8f6efea74b2971257af23d5193ba3a"
	amount := "100"
	fmt.Printf("to address = %s\n", toAddress)

	// Initialize the client
	restUrl := "https://fullnode.devnet.aptoslabs.com"
	client, _ := aptosclient.Dial(context.Background(), restUrl)

	// Get Sender's account data and ledger info
	accountData, _ := client.GetAccount(fromAddress)
	ledgerInfo, _ := client.LedgerInfo()

	// Get gas price
	gasPrice, _ := client.EstimateGasPrice()

	// Build paylod
	payload := &aptostypes.Payload{
		Type:          "entry_function_payload",
		Function:      "0x2c3c7c96c148777979c1964b555f1c559b8f6efea74b2971257af23d5193ba3a::order::pull",
		TypeArguments: []string{"0x1::aptos_coin::AptosCoin"},
		Arguments: []interface{}{
			toAddress, amount,
		},
	}

	// Build transaction
	transaction := &aptostypes.Transaction{
		Sender:                  fromAddress,
		SequenceNumber:          accountData.SequenceNumber,
		MaxGasAmount:            2000,
		GasUnitPrice:            gasPrice,
		Payload:                 payload,
		ExpirationTimestampSecs: ledgerInfo.LedgerTimestamp + 600, // 10 minutes timeout
	}

	// Get signing message from remote server
	// Note: Later we will implement the local use of BCS encoding to create signing messages
	signingMessage, _ := client.CreateTransactionSigningMessage(transaction)

	// Sign message and complete transaction information
	signatureData := account.Sign(signingMessage, "")
	signatureHex := "0x" + hex.EncodeToString(signatureData)
	publicKey := "0x" + hex.EncodeToString(account.PublicKey)
	transaction.Signature = &aptostypes.Signature{
		Type:      "ed25519_signature",
		PublicKey: publicKey,
		Signature: signatureHex,
	}

	// Submit transaction
	newTx, _ := client.SubmitTransaction(transaction)
	fmt.Printf("tx hash = %v\n", newTx.Hash)
}

func printLine(content string) {
	fmt.Printf("================= %s =================\n", content)
}

func printError(err error) {
	var restError *aptostypes.RestError
	if b := errors.As(err, &restError); b {
		fmt.Printf("code: %d, message: %s, aptos_ledger_version: %d\n", restError.Code, restError.Message, restError.AptosLedgerVersion)
	} else {
		fmt.Printf("err: %s\n", err.Error())
	}
}
