package listener

import (
	"context"
	"fmt"
	_ "fmt"
	"log"
	"os"

	"github.com/dheerajkumardk/blockchain_explorer_backend/database"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
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
				BlockNumber:     block.Number().Uint64(),
				Timestamp:       header.Time,
				FeeRecipient:    header.Coinbase.String(),
				BlockReward:     "",
				TotalDifficulty: header.Difficulty.Uint64(),
				Size:            block.Size(),
				GasUsed:         header.GasUsed,
				GasLimit:        header.GasLimit,
				BaseFeePerGas:   header.BaseFee.String(),
				BurntFees:       "",
				Hash:            block.Hash().String(),
				ParentHash:      header.ParentHash.Hex(),
				StateRoot:       header.Root.Hex(),
				WithdrawalsRoot: header.WithdrawalsHash.String(),
				Nonce:           header.Nonce.Uint64(),
			}

			db := database.BlockDB
			if db == nil {
				log.Print("DB is nil")
			} else {
				log.Print("All OK!")
			}
			// insert block into db
			result := db.Create(&newBlock)
			if result.Error != nil {
				log.Printf("Failed to insert block %d %v", block.Number(), result.Error)
			} else {
				log.Printf("Stored block %d\n", block.Number())
			}

			// console log the data
			// fmt.Println("block Number: ", block.Number().Int64())
			// fmt.Println("BlockHash: ", block.Hash().String())
			// fmt.Println("Timestamp: ", header.Time)
			// fmt.Println("Proposed On: ")
			// fmt.Println("Transactions: ", block.Transactions())
			// fmt.Println("Withdrawals: ", block.Withdrawals())
			// fmt.Println("Fee Recipient", header.Coinbase)
			// fmt.Println("Block Reward")
			// fmt.Println("Total Difficulty", header.Difficulty)
			// fmt.Println("Size", block.Size())
			// fmt.Println("GasUsed", header.GasUsed)
			// fmt.Println("Gas Limit", header.GasLimit)
			// fmt.Println("BaseFeePerGas", header.BaseFee)
			// fmt.Println("BurntFees", uint64(header.BaseFee.Int64())+header.GasUsed)
			// fmt.Println("Extra Data", header.Extra)
			// fmt.Println("Hash", block.Hash().String())
			// fmt.Println("Parent Hash", header.ParentHash)
			// fmt.Println("StateRoot", header.Root)
			// fmt.Println("WithdrawalsRoot", header.WithdrawalsHash)
			// fmt.Println("Nonce", header.Nonce.Uint64())

			// transactions
			for _, tx := range block.Transactions() {
				// get from address
				from, err := types.Sender(types.LatestSignerForChainID(tx.ChainId()), tx)
				if err != nil {
					fmt.Println("Failed to get `from` address", err)
				}

				// to handle nil on contract creation txn
				var toAddress common.Address
				if tx.To() == nil {
					toAddress = crypto.CreateAddress(from, tx.Nonce())
					fmt.Println("Contract creation txn: ", tx.Hash().String(), " :: toAddress -> ", toAddress.String())
				} else {
					toAddress = *tx.To()
				}
				// create txn obj to insert into db
				newTxn := database.Transaction{
					TxHash:      tx.Hash().String(),
					BlockNumber: block.Number().Uint64(),
					From:        from.String(),
					To:          toAddress.String(),
					Value:       tx.Value().String(),
					TxnFees:     "",
					Timestamp:   uint64(tx.Time().Unix()),
				}
				// insert txn into db
				result = db.Create(&newTxn)
				if result.Error != nil {
					log.Printf("Failed to store txn %s", tx.Hash().String())
				}

				// fmt.Println("txn: ", index)
				// fmt.Println("hash: ", tx.Hash().String())
				// fmt.Println("time: ", tx.Time())
				// fmt.Println("value ", tx.Value().String())
				// fmt.Println("gas ", tx.Gas())
				// fmt.Println("gas price ", tx.GasPrice().Uint64())

				// Account
				// Update sender account, always update nonce for EOA
				if err := updateOrCreateAccount(client, db, from, true); err != nil {
					log.Printf("Failed to update sender account %s: %v\n", from.String(), err)
					continue
				}

				// Update receiver account, if not a contract
				if tx.To() != nil {
					if err := updateOrCreateAccount(client, db, *tx.To(), false); err != nil {
						log.Printf("Failed to update receiver account %s: %v", tx.To().String(), err)
						continue
					}
				}
			}

			// withdrawals
			for _, withdrawl := range block.Withdrawals() {
				// create withdraw obj to insert into db
				newWithdrawl := database.Withdrawal{
					Index:          withdrawl.Index,
					BlockNumber:    block.Number().Uint64(),
					ValidatorIndex: withdrawl.Validator,
					Recipient:      withdrawl.Address.String(),
					Amount:         withdrawl.Amount,
				}

				// insert withdrawal into db
				result = db.Create(&newWithdrawl)
				if result.Error != nil {
					log.Printf("Failed to store withdrawal %d", withdrawl.Index)
				}
			}

		}
	}
}

func isContractAddress(client *ethclient.Client, address common.Address) (bool, error) {
	bytecode, err := client.CodeAt(context.Background(), address, nil)
	if err != nil {
		return false, err
	}
	return len(bytecode) > 0, nil
}

func updateOrCreateAccount(client *ethclient.Client, db *gorm.DB, address common.Address, updateNonce bool) error {
	// fetch latest eth balance
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		return err
	}

	// fetch nonce, only for EOA
	var nonce uint64
	if updateNonce {
		nonce, err = client.NonceAt(context.Background(), address, nil)
		if err != nil {
			return err
		}
	}

	// check if an EOA or a contract
	isContract, err := isContractAddress(client, address)
	if err != nil {
		return err
	}
	// determine address type
	addressType := "EOA"
	if isContract {
		addressType = "Contract"
	}

	account := database.Account{
		Address:     address.String(),
		AddressType: addressType,
		ETHBalance:  balance.String(),
		Nonce:       nonce,
	}
	// Create/Update the account in the database
	if err := db.Where(database.Account{Address: address.String()}).Assign(account).FirstOrCreate(&account).Error; err != nil {
		return err
	}
	return nil
}
