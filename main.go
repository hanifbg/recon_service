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
}
