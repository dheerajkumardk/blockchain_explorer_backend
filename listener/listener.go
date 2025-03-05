package listener

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dheerajkumardk/blockchain_explorer_backend/database"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

func SubscribeBlocks() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env")
	}

	client, err := ethclient.Dial(os.Getenv("ETH_WSS_URL"))
	if err != nil {
		log.Fatalf("Error connecting websocket: %v", err)
	}

	// create channel to receive latest block headers
	headers := make(chan *types.Header)

	sub, err := client.SubscribeNewHead(context.Background(), headers)
	if err != nil {
		log.Fatalf("Error subscribing %v", err)
	}

	for {
		select {
		case err := <-sub.Err():
			log.Fatal(err)
		case header := <-headers:

			block, err := client.BlockByHash(context.Background(), header.Hash())
			if err != nil {
				log.Fatal(err)
			}

			// create block to insert into db
			newBlock := database.Block{
				BlockNumber: block.Number().Uint64(),
				Timestamp: header.Time,
				FeeRecipient: header.Coinbase.String(),
				BlockReward: "",
				TotalDifficulty: header.Difficulty.Uint64(),
				Size: block.Size(),
				GasUsed: header.GasUsed,
				GasLimit: header.GasLimit,
				BaseFeePerGas: header.BaseFee.String(),
				BurntFees: "",
				Hash: block.Hash().String(),
				ParentHash: header.ParentHash.Hex(),
				StateRoot: header.Root.Hex(),
				WithdrawalsRoot: header.WithdrawalsHash.String(),
				Nonce: header.Nonce.Uint64(),
			}
			// insert block into db
			db := database.BlockDB
			db.Create(&newBlock)

			// console log the data
			fmt.Println("block Number: ", block.Number().Int64())
			fmt.Println("BlockHash: ", block.Hash().String())
			fmt.Println("Timestamp: ", header.Time)
			fmt.Println("Proposed On: ")
			fmt.Println("Transactions: ", block.Transactions())
			fmt.Println("Withdrawals: ", block.Withdrawals())
			fmt.Println("Fee Recipient", header.Coinbase)
			fmt.Println("Block Reward")
			fmt.Println("Total Difficulty", header.Difficulty)
			fmt.Println("Size", block.Size())
			fmt.Println("GasUsed", header.GasUsed)
			fmt.Println("Gas Limit", header.GasLimit)
			fmt.Println("BaseFeePerGas", header.BaseFee)
			fmt.Println("BurntFees", uint64(header.BaseFee.Int64()) + header.GasUsed)
			fmt.Println("Extra Data", header.Extra)
			fmt.Println("Hash", block.Hash().String())
			fmt.Println("Parent Hash", header.ParentHash)
			fmt.Println("StateRoot", header.Root)
			fmt.Println("WithdrawalsRoot", header.WithdrawalsHash)
			fmt.Println("Nonce", header.Nonce.Uint64())


			// transactions
			for index, tx := range block.Transactions() {
				// create txn obj to insert into db
				newTxn := database.Transaction{
					TxHash: tx.Hash().String(),
					BlockNumber: block.Number().Uint64(),
					From: "",
					To: tx.To().String(),
					Value: tx.Value().String(),
					TxnFees: "",
					Timestamp: uint64(tx.Time().Unix()),

				}
				// insert txn into db
				db.Create(&newTxn)

				fmt.Println("txn: ", index)
				fmt.Println("hash: ", tx.Hash().String())
				fmt.Println("time: ", tx.Time())
				fmt.Println("value ", tx.Value().String())
				fmt.Println("gas ", tx.Gas())
				fmt.Println("gas price ", tx.GasPrice().Uint64())
				fmt.Println("nonce ", tx.Nonce())
				fmt.Println("Data ", tx.Data())
				fmt.Println("to", tx.To())

				fmt.Println()
			}

			// withdrawals
			for _, withdrawl := range block.Withdrawals() {
				// create withdraw obj to insert into db
				newWithdrawl := database.Withdrawal{
					Index: withdrawl.Index,
					BlockNumber: block.Number().Uint64(),
					ValidatorIndex: withdrawl.Validator,
					Recipient: withdrawl.Address.String(),
					Amount: withdrawl.Amount,
				}

				// insert withdrawal into db
				db.Create(&newWithdrawl)
			}

		}
	}
}
