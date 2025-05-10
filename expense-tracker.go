package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
	"io/ioutil"
)

type Expense struct {
	Amount   float64   `json:"amount"`
	Category string    `json:"category"`
	Date     time.Time `json:"date"`
}

const dataFile = "expenses.json"

func loadExpenses() ([]Expense, error) {
	data, err := ioutil.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			return []Expense{}, nil
		}
		return nil, err
	}
	var expenses []Expense
	err = json.Unmarshal(data, &expenses)
	return expenses, err
}

func saveExpenses(expenses []Expense) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dataFile, data, 0644)
}

func addExpense(args []string) {
	if len(args) < 2 {
		fmt.Println("Usage: add <amount> <category>")
		return
	}
	amount, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Println("Invalid amount.")
		return
	}
	category := args[1]
	expenses, err := loadExpenses()
	if err != nil {
		fmt.Println("Error loading expenses:", err)
		return
	}
	expense := Expense{
		Amount:   amount,
		Category: category,
		Date:     time.Now(),
	}
	expenses = append(expenses, expense)
	if err := saveExpenses(expenses); err != nil {
		fmt.Println("Error saving expense:", err)
	} else {
		fmt.Println("Expense added.")
	}
}

func listExpenses() {
	expenses, err := loadExpenses()
	if err != nil {
		fmt.Println("Error loading expenses:", err)
		return
	}
	if len(expenses) == 0 {
		fmt.Println("No expenses recorded.")
		return
	}
	var total float64
	for _, e := range expenses {
		fmt.Printf("- $%.2f [%s] (%s)\n", e.Amount, e.Category, e.Date.Format("2006-01-02"))
		total += e.Amount
	}
	fmt.Printf("Total: $%.2f\n", total)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Commands: add <amount> <category> | list")
		return
	}
	cmd := os.Args[1]
	switch cmd {
	case "add":
		addExpense(os.Args[2:])
	case "list":
		listExpenses()
	default:
		fmt.Println("Unknown command.")
	}
}
