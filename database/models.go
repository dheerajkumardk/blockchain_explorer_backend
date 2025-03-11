package database

import "gorm.io/gorm"

// Struct to store block data in db
type Block struct {
	gorm.Model
	BlockNumber     uint64 `json:"blockNumber"`
	Timestamp       uint64 `json:"timestamp"`
	FeeRecipient    string `json:"feeRecipient"`
	BlockReward     string `json:"blockReward"`
	TotalDifficulty uint64 `json:"totalDifficulty"`
	Size            uint64 `json:"size"`
	GasUsed         uint64 `json:"gasUsed"`
	GasLimit        uint64 `json:"gasLimit"`
	BaseFeePerGas   string `json:"baseFee"`
	BurntFees       string `json:"burntFee"`
	Hash            string `json:"hash"`
	ParentHash      string `json:"parentHash"`
	StateRoot       string `json:"stateRoot"`
	WithdrawalsRoot string `json:"withdrawalsRoot"`
	Nonce           uint64 `json:"nonce"`
}

// Struct to store txn data in db
type Transaction struct {
	gorm.Model
	TxHash      string `json:"txHash"`
	BlockNumber uint64 `json:"blockNumber"`
	From        string `json:"from"`
	To          string `json:"to"`
	Value       string `json:"value"`
	TxnFees     string `json:"txnFees"`
	Timestamp   uint64 `json:"timestamp"`
}

// Struct to store withdrawals in db
type Withdrawal struct {
	gorm.Model
	Index          uint64 `json:"index"`
	BlockNumber    uint64 `json:"blockNumber"`
	ValidatorIndex uint64 `json:"validatorIndex"`
	Recipient      string `json:"recipient"`
	Amount         uint64 `json:"amount"`
}

// Struct to store account data in db
type Account struct {
	gorm.Model
	Address     string `json:"address"`
	AddressType string `json:"addressType"`
	ETHBalance  string `json:"ethBalance"`
	Nonce       uint64 `json:"nonce"`
}

// Struct to store account address & transaction hash
type AccountTransaction struct {
	gorm.Model
	Address string `json:"address"`
	TxHash  string `json:"string"`
	Role string `json:"role"`
}
