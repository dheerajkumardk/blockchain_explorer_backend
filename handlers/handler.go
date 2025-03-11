package handlers

import (
	"fmt"
	"strconv"

	"github.com/dheerajkumardk/blockchain_explorer_backend/database"
	"github.com/gofiber/fiber/v2"
)

func GetAllBlocks(c *fiber.Ctx) error {
	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database is not initialized",
		})
	}
	var blocks []database.Block
	if err := db.Find(&blocks).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Error retrieving blocks",
			"details": err.Error(),
		})
	}
	return c.JSON(blocks)
}

func GetBlockByBlockNumber(c *fiber.Ctx) error {
	db := database.BlockDB
	blockNumber := c.Params("blockNumber")
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database is not initialized",
		})
	}
	var block database.Block
	if err := db.Find(&block, "block_number = ?", blockNumber).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error":   "Error retrieving block",
			"details": err.Error(),
		})
	}
	if block.Hash == "" {
		c.Status(500).Send([]byte("No block found with given BlockHash"))
	}
	return c.JSON(block)
}

func GetBlockTransactions(c *fiber.Ctx) error {
	blockNumber, _ := strconv.ParseUint(c.Params("blockNumber"), 10, 64)

	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database is not initialized",
		})
	}
	var transactions []database.Transaction
	if err := db.Find(&transactions, "block_number = ?", blockNumber).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error retrieving transactions",
		})
	}
	if len(transactions) == 0 {
		c.Status(500).Send([]byte("No transactions found for given BlockNumber"))
	}
	return c.JSON(transactions)
}

func GetAllTransactions(c *fiber.Ctx) error {
	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database is not initialized",
		})
	}
	var transactions []database.Transaction
	if err := db.Find(&transactions).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error retrieving transactions",
		})
	}
	if len(transactions) == 0 {
		c.Status(500).Send([]byte("No transactions found"))
	}
	return c.JSON(transactions)
}

func GetTransactionInfo(c *fiber.Ctx) error {
	txHash := c.Params("txHash")
	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Database is not initialised",
		})
	}
	var transaction database.Transaction
	if err := db.Find(&transaction, "tx_hash = ?", txHash).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error retrieving transaction Info",
		})
	}
	if transaction.BlockNumber == 0 {
		c.Status(500).Send([]byte("No transaction found for given hash"))
	}
	return c.JSON(transaction)
}

func AccountTransactions(c *fiber.Ctx) error {
	address := c.Params("address")
	var transactions []database.Transaction

	// find the account by address
	var account database.Account
	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Datbase is not initialised",
		})
	}
	if err := db.Where("address = ?", address).First(&account).Error; err != nil {
		return err
	}

	// Find all txns linked to the account
	if err := db.Joins("JOIN account_transactions ON account_transactions.tx_Hash = transactions.tx_Hash").Where("account_transactions.address = ?", account.Address).Find(&transactions).Error; err != nil {
		return err
	}

	return c.JSON(transactions)
}

func GetAllWithdrawals(c *fiber.Ctx) error {
	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Datbase is not initialised",
		})
	}
	var withdrawals []database.Withdrawal
	if err := db.Find(&withdrawals).Error; err != nil {
		c.Status(500).JSON(fiber.Map{
			"error": "Error retreiving withdrawals",
		})
	}
	if len(withdrawals) == 0 {
		c.Status(500).Send([]byte("No withdrawals found"))
	}
	return c.JSON(withdrawals)
}

func GetWithdrawalInfo(c *fiber.Ctx) error {
	index, _ := strconv.ParseUint(c.Params("index"), 10, 64)
	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Datbase is not initialised",
		})
	}
	var withdrawal database.Withdrawal
	if err := db.Find(&withdrawal, "`index` = ?", index).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error retreiving withdrawal info",
		})
	}
	fmt.Println("Index -> ", index)
	fmt.Println("Withdrawal -> ", withdrawal)
	if withdrawal.Index == 0 {
		return c.Status(500).Send([]byte("Withdrawal for given index not found"))
	}
	return c.JSON(withdrawal)
}

func GetAccountInfo(c *fiber.Ctx) error {
	address := c.Params("address")
	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Datbase is not initialised",
		})
	}
	var account database.Account
	if err := db.Find(&account, "address = ?", address).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error retreiving txns for given address",
		})
	}
	if account.Address == "" {
		return c.Status(500).Send([]byte("No info found for given address"))
	}
	return c.JSON(account)
}

func GetAccountBalance(c *fiber.Ctx) error {
	address := c.Params("address")
	db := database.BlockDB
	if db == nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Datbase is not initialised",
		})
	}
	var account database.Account
	if err := db.Find(&account, "address = ?", address).Error; err != nil {
		return c.Status(500).JSON(fiber.Map{
			"error": "Error retreiving txns for given address",
		})
	}
	if account.Address == "" {
		c.Status(500).Send([]byte("No info found for given address"))
	}
	return c.JSON(account.ETHBalance)
}


// Search endpoints

// Not exposing post End-points
