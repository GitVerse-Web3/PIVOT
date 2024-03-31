package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/fasthttp/websocket"
	relayer "github.com/haisenCx/NOSTR-AGI-GIT"

	"github.com/fiatjaf/eventstore"
	"github.com/fiatjaf/eventstore/postgresql"

	//"github.com/fiatjaf/relayer/v2"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
)

/**
接收来自Contract的agigit pull(modelHash)请求。
对Contract发出takeOrder(x=prime1*prime2)请求。
接收来自User的handshake请求。
向User发送fileTransfer(e=signaturePreservingEncryption(model,prime1+prime2))请求。
在Contract验证通过后，向Contract发送completeOrder(prime1,prime2)请求。
*/

type Relay struct {
	PostgresDatabase string `envconfig:"POSTGRESQL_DATABASE"`

	storage *postgresql.PostgresBackend

	userPaid map[string]bool

	userPrimes map[string][2]*big.Int

	connections map[string]*relayer.WebSocket
}

const numFiles = 4 // 传输文件分块个数

func (r *Relay) Init() error {
	err := envconfig.Process("", r)
	if err != nil {
		return fmt.Errorf("couldn't process envconfig: %w", err)
	}

	// 初始化 userPrimes
	r.userPrimes = make(map[string][2]*big.Int)
	// 初始化 userPaid
	r.userPaid = make(map[string]bool)

	r.connections = make(map[string]*relayer.WebSocket)
	// every hour, delete all very old events
	go func() {
		db := r.Storage(context.TODO()).(*postgresql.PostgresBackend)

		for {
			time.Sleep(60 * time.Minute)
			db.DB.Exec(`DELETE FROM event WHERE created_at < $1`, time.Now().AddDate(0, -3, 0).Unix()) // 3 months
		}
	}()

	return nil
}

/*
- 事件格式
{
	"id": "4376c65d2f232afbe9b882a35baa4f6fe8667c4e684749af565f981833ed6a65",
    "pubkey": "6e468422dfb74a5738702a8823b9b28168abab8655faacb6853cd0ee15deee93",
    "created_at": 1673347337,
    "kind": 1111,
    "tags": [
    ["e", "3da979448d9ba263864c4d6f14984c423a3838364ec255f03c7904b1ae77f206"],
    ["p", "bf2376e17ba4ec269d10fcc996a4746b451152be9031fa48e74553dde5526bce"]
    ],
    "content": "Walled gardens became prisons, and nostr is the first step towards tearing down the prison walls.",
    "sig": "908a15e46fb4d8675bab026fc230a0e3542bfade63da02d542fb78b2a8513fcd0092619a2c8c1221e581946e0191f2af505dfdf8657a414dbca329186f009262"
    }
*/
//根据变量判断是否用户已经支付足够的费用
func (r *Relay) hasUserPaid(pubkey string) bool {
	paid, ok := r.userPaid[pubkey]
	return ok && paid
}

// 监听来自Contract的agigit pull(modelHash)请求
func (r *Relay) AcceptReq(ctx context.Context, id string, filters nostr.Filters, authedPubkey string) bool {
	// 实现对Contract发出agigit pull(modelHash)请求的逻辑
	// 如果用户已经支付足够的费用
	if r.hasUserPaid(authedPubkey) {
		return false
	}
	// 如果用户没有支付足够的费用
	return true
}

// 中继调用合约 takeOrder(x=prime1*prime2)函数
func (r *Relay) TakeOrder(prime1, prime2 big.Int) {
	// 实现对Contract发出takeOrder(x=prime1*prime2)请求的逻辑
	x := new(big.Int).Mul(&prime1, &prime2)
	fmt.Println("Sending takeOrder request to Contract with x: ", x)
}

// 中继监听到takeOrder()[$DOWN]事件
func (r *Relay) listenTakeOrderEvent() {
	// 实现监听到takeOrder()[$DOWN]事件的逻辑
	fmt.Println("Listening to takeOrder()[$DOWN] event")
}

// 中继调用合约completeOrder(prime1,prime2)
func (r *Relay) CompleteOrder(prime1, prime2 big.Int) {
	// 实现在Contract验证通过后，向Contract发送completeOrder(prime1,prime2)请求的逻辑
	fmt.Println("Sending completeOrder request to Contract with prime1 and prime2: ", prime1, prime2)
}

// 中继监听到completeOrder事件
func (r *Relay) listenCompleteOrderEvent() {
	// 实现监听到completeOrder事件的逻辑
	fmt.Println("Listening to completeOrder event")
}

func (r *Relay) FileTransfer(modelPath string, prime1 big.Int, prime2 big.Int, pubKey string) {
	fmt.Println("Sending fileTransfer request to User")
	// 实现向User发送fileTransfer(e=signaturePreservingEncryption(model,prime1+prime2))请求的逻辑
	r.signaturePreservingEncryption(modelPath, *new(big.Int).Add(&prime1, &prime2), pubKey)
}

func (r *Relay) GetConnection(pubKey string) (*relayer.WebSocket, bool) {
	conn, ok := r.connections[pubKey]
	return conn, ok
}

// 定义消息结构
type Message struct {
	Index   int     `json:"index"`
	Sum     string  `json:"sum"`
	Percent float64 `json:"percent"`
	Data    []byte  `json:"data"`
}

func (r *Relay) signaturePreservingEncryption(modelPath string, sum big.Int, pubKey string) {
	// 读取文件内容
	data, err := os.ReadFile(modelPath)
	if err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	// 计算每个文件应该有多少字节
	fileSize := len(data) / numFiles
	if len(data)%numFiles != 0 {
		fileSize++
	}

	// 将文件的内容分割成多个数据块
	conn, ok := r.GetConnection(pubKey)
	if !ok {
		log.Fatalf("no connection found for pubkey: %s", pubKey)
	}

	for i := 0; i < len(data); i += fileSize {
		end := i + fileSize
		if end > len(data) {
			end = len(data)
		}
		chunk := data[i:end]

		// 计算当前的传输百分比
		percent := float64(end) / float64(len(data))

		// 创建一个消息
		msg := Message{
			Index:   i / fileSize,
			Sum:     sum.String(),
			Percent: percent,
			Data:    chunk,
		}

		// 将消息序列化为 JSON
		msgBytes, err := json.Marshal(msg)
		if err != nil {
			log.Fatalf("failed to marshal message: %v", err)
		}

		// 发送消息
		err = conn.WriteMessage(websocket.TextMessage, msgBytes)
		if err != nil {
			log.Fatalf("failed to send message: %v", err)
		}
	}
	// 如果正常发送完毕
	fmt.Println("File transfer complete")

}

func (r *Relay) AcceptEvent(ctx context.Context, evt *nostr.Event, conn *relayer.WebSocket) bool {
	fmt.Print("AcceptEvent: ")
	// 只有当客户端是来自Contract的takeOrder()[$FEE]后
	if r.userPaid[evt.PubKey] {
		return false
	}

	//防止event太大
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 100000 {
		return false
	}

	//假设kind1111为agigitPull
	fmt.Println("evt.Kind: ", evt.Kind)
	if evt.Kind == 1111 {
		prime1, err := generatePrime(128)
		if err != nil {
			fmt.Errorf("couldn't generate prime1: %w", err)
		}

		prime2, err := generatePrime(128)
		if err != nil {
			fmt.Errorf("couldn't generate prime2: %w", err)
		}
		defer func() {
			r.CompleteOrder(*prime1, *prime2)
		}()
		fmt.Println("prime1: ，prime2:，pubkey: ", prime1, prime2, evt.PubKey)
		r.userPrimes[evt.PubKey] = [2]*big.Int{prime1, prime2}
		r.connections[evt.PubKey] = conn
		modelPath := "testfile.txt"
		r.FileTransfer(modelPath, *prime1, *prime2, evt.PubKey)

	}

	return true
}
func generatePrime(bits int) (*big.Int, error) {
	p, err := rand.Prime(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (r *Relay) Name() string {
	return "BasicRelay"
}

func (r *Relay) Storage(ctx context.Context) eventstore.Store {
	return r.storage
}

// type IpAddress string
// type PublicKey string
// type CommitMessage string
// type ModelHash string

// // Commit
// type Commit struct {
// 	ModelHash     string
// 	CommitMessage string
// 	Author        PublicKey
// 	Signature     string
// }

// func (c Commit) GetFullModel() []byte {

// 	return []byte{}
// }

// // Event
// type Event struct {
// 	Model  ModelHash
// 	Relays []Relay
// }

// func (e Event) getCommit() Commit {

// 	return Commit{}
// }

// // DeltaModel
// type DeltaModel struct {
// 	ParentModel Event
// 	IsEncrypted bool
// 	ModelDelta  []byte
// 	Commit
// }

// // MergeModel
// type MergeModel struct {
// 	ParentModels []Event
// 	Commit
// }

func main() {
	r := Relay{}
	if err := envconfig.Process("", &r); err != nil {
		log.Fatalf("failed to read from env: %v", err)
		return
	}
	r.storage = &postgresql.PostgresBackend{DatabaseURL: r.PostgresDatabase}
	server, err := relayer.NewServer(&r)
	if err != nil {
		log.Fatalf("failed to create server: %v", err)
	}
	if err := server.Start("0.0.0.0", 7447); err != nil {
		log.Fatalf("server terminated: %v", err)
	}
}
