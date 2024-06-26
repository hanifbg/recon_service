package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

type Transaction struct {
	TrxID              string    `csv:"trxID"`
	Amount             float64   `csv:"amount"`
	Type               string    `csv:"type"` // DEBIT or CREDIT
	TransactionTimeStr string    `csv:"transactionTime"`
	TransactionTime    time.Time `csv:"-"`
}

type BankStatement struct {
	UniqueIdentifier string    `csv:"unique_identifier"`
	Amount           float64   `csv:"amount"`
	DateStr          string    `csv:"date"`
	Date             time.Time `csv:"-"`
}

type ReconciliationSummary struct {
	TotalTransactionsProcessed  int
	TotalMatchedTransactions    int
	TotalUnmatchedTransactions  int
	UnmatchedSystemTransactions []Transaction
	UnmatchedBankStatements     []BankStatement
	TotalDiscrepancies          float64
}

func main() {
	filePath := "system_transactions.csv"
	transactionsFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer transactionsFile.Close()

	var transactions []Transaction
	if err = gocsv.UnmarshalFile(transactionsFile, &transactions); err != nil {
		fmt.Println(err)
	}

	transactions[0].TransactionTime, err = time.Parse("2006-01-02 15:04:05", transactions[0].TransactionTimeStr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(transactions)

	filePath = "bank_statements.csv"
	statementsFile, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer statementsFile.Close()

	var statements []BankStatement
	if err = gocsv.UnmarshalFile(statementsFile, &statements); err != nil {
		fmt.Println(err)
	}

	statements[0].Date, err = time.Parse("2006-01-02 15:04:05", statements[0].DateStr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(statements)

	//recon
	startDate := "2023-01-01"
	endDate := "2023-12-31"
	summary := ReconciliationSummary{}

	start, _ := time.Parse("2006-01-02", startDate)
	end, _ := time.Parse("2006-01-02", endDate)
	// Filter transactions and statements by date range
	var filteredTransactions []Transaction
	var filteredBankStatements []BankStatement
	for _, t := range transactions {
		if t.TransactionTime.After(start) && t.TransactionTime.Before(end) {
			filteredTransactions = append(filteredTransactions, t)
		}
	}
	for _, b := range statements {
		if b.Date.After(start) && b.Date.Before(end) {
			filteredBankStatements = append(filteredBankStatements, b)
		}
	}

	summary.TotalTransactionsProcessed = len(filteredTransactions)

	// Matching logic
	matchedTransactions := make(map[string]bool)
	for _, t := range filteredTransactions {
		found := false
		for _, b := range filteredBankStatements {
			if t.Amount == b.Amount && t.TransactionTime.Equal(b.Date) {
				matchedTransactions[t.TrxID] = true
				found = true
				break
			}
		}
		if !found {
			summary.UnmatchedSystemTransactions = append(summary.UnmatchedSystemTransactions, t)
		} else {
			summary.TotalMatchedTransactions++
		}
	}

	// Identify unmatched bank statements
	for _, b := range filteredBankStatements {
		matched := false
		for _, t := range filteredTransactions {
			if t.Amount == b.Amount && t.TransactionTime.Equal(b.Date) {
				matched = true
				break
			}
		}
		if !matched {
			summary.UnmatchedBankStatements = append(summary.UnmatchedBankStatements, b)
		}
	}

	// Calculate total discrepancies
	for _, t := range filteredTransactions {
		if !matchedTransactions[t.TrxID] {
			for _, b := range filteredBankStatements {
				if t.Amount != b.Amount && t.TransactionTime.Equal(b.Date) {
					summary.TotalDiscrepancies += abs(t.Amount - b.Amount)
				}
			}
		}
	}

	summary.TotalUnmatchedTransactions = len(summary.UnmatchedSystemTransactions) + len(summary.UnmatchedBankStatements)

	fmt.Printf("Total Transactions Processed: %d\n", summary.TotalTransactionsProcessed)
	fmt.Printf("Total Matched Transactions: %d\n", summary.TotalMatchedTransactions)
	fmt.Printf("Total Unmatched Transactions: %d\n", summary.TotalUnmatchedTransactions)
	fmt.Printf("Unmatched System Transactions: %v\n", summary.UnmatchedSystemTransactions)
	fmt.Printf("Unmatched Bank Statements: %v\n", summary.UnmatchedBankStatements)
	fmt.Printf("Total Discrepancies: %f\n", summary.TotalDiscrepancies)
}

func abs(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}
