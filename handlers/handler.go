package handlers

import (
	"log"

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
			"error": "Error retrieving blocks",
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
			"error": "Error retrieving block",
			"details": err.Error(),
		})
	}
	if block.Hash == "" {
		c.Status(500).Send([]byte("No block found with given BlockHash"))
	}
	return c.JSON(block)
}

func GetBlockTransactions(c *fiber.Ctx) error {
	blockNumber := c.Params("blockNumber")
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

