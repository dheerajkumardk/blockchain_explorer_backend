package routes

import (
	"github.com/dheerajkumardk/blockchain_explorer_backend/handlers"
	"github.com/gofiber/fiber/v2"
)

// lists all routes provided by the explorer
func SetupRoutes(app *fiber.App) {
	// Get all blocks - limit, paginated
	app.Get("/blocks", handlers.GetAllBlocks)

	// Retrieve info abt a block
	app.Get("/blocks/:blockNumber", handlers.GetBlockByBlockNumber)

	// Get all txnx by blockNumber
	app.Get("blocks/:block/transactions", handlers.GetBlockTransactions)

	// Get all txns -> paginated list [page, limit, sort by timestamp]
	app.Get("/transactions", handlers.GetAllTransactions)

	// Retrieve txn info by hash
	app.Get("/transactions/:txHash", handlers.GetTransactionInfo)

	// Get all txns for an address
	app.Get("/accounts/:address/transactions", handlers.AccountTransactions)

	// Get all withdrawals
	app.Get("/withdrawals", handlers.GetAllWithdrawals)

	// Retrieve withdrawal by index
	app.Get("/withdrawals/:index", handlers.GetWithdrawalInfo)

	// Get info abt an account
	app.Get("/accounts/:address", handlers.GetAccountInfo)

	// Get balance of an account
	app.Get("/accounts/:address/balance", handlers.GetAccountBalance)

	// Search End-points
	// Search for blocks (blockNumber), transactions(txnHash), or accounts(Address)
	// app.Get("/search")

	// Post a new block
	// app.Post("/blocks")

	// Post a new transaction
	// app.Post("/transactions")

	// Post a new withdrawl
	// app.Post("/withdrawals")
}
