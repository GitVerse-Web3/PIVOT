package types

type PullParams struct {
	ModelHash       int   `json:"model_hash"`
	RelayPubkeyList []int `json:"relay_pubkey_list"`
	Down            int   `json:"down"`
	Cost            int   `json:"cost"`
	Tip             int   `json:"tip"`
}

type PullResult struct {
	OrderID string `json:"order_id"`
}

type PullEvent struct {
	Address         string `json:"address"`
	OrderID         string `json:"order_id"`
	ModelHash       int    `json:"model_hash"`
	RelayPubkeyList []int  `json:"relay_pubkey_list"`
	Down            int    `json:"down"`
	Cost            int    `json:"cost"`
	Tip             int    `json:"tip"`
}

type RaisePaymentParams struct {
	OrderID string `json:"order_id"`
	Down    int    `json:"down"`
	Cost    int    `json:"cost"`
	Tip     int    `json:"tip"`
}

type RaisePaymentEvent struct {
	Address         string `json:"address"`
	OrderID         string `json:"order_id"`
	ModelHash       int    `json:"model_hash"`
	RelayPubkeyList []int  `json:"relay_pubkey_list"`
	Down            int    `json:"down"`
	Cost            int    `json:"cost"`
	Tip             int    `json:"tip"`
}

type TakeOrderParams struct {
	OrderID      string `json:"order_id"`
	OrderAddress string `json:"order_address"`
	RelayPubkey  int    `json:"relay_pubkey"`
	X            int    `json:"x"`
}

type TakeOrderEvent struct {
	RelayAddress string `json:"relay_address"`
	OrderID      string `json:"order_id"`
	OrderAddress string `json:"order_address"`
}

type GetOrderStatusParams struct {
	OrderID string `json:"order_id"`
}

type GetOrderStatusResult struct {
	Status int `json:"status"`
}

type ConfirmReceivedParams struct {
	OrderID string `json:"order_id"`
}

type CompleteOrderParams struct {
	OrderID      string `json:"order_id"`
	OrderAddress string `json:"order_address"`
	Prime1       int    `json:"prime1"`
	Prime2       int    `json:"prime2"`
}

type CompleteOrderEvent struct {
	OrderAddress string `json:"order_address"`
	OrderID      string `json:"order_id"`
	Prime1       int    `json:"prime1"`
	Prime2       int    `json:"prime2"`
}
