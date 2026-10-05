// Task 4: Bank account
// 1. Create BankAccount with an exported Owner and an unexported balance .
// 2. Add Deposit(amount float64) , Withdraw(amount float64) error (return an
// error if there's not enough money) and Balance() float64 .
// 3. Choose the receiver type for each method and justify your choice on
// the defence. Should Balance() have a value or a pointer receiver?
// 4. Put three accounts in a slice and add 5% interest to each one in a loop.
// Make sure the balances really change.
package main

import (
	"errors"
	"fmt"
)

type BankAccount struct {
	Owner   string  // exported
	balance float64 // unexported: can only be changed through methods (inside this package)
}

// Deposit: pointer receiver, it modifies the account.
func (a *BankAccount) Deposit(amount float64) {
	if amount <= 0 {
		return
	}
	a.balance += amount
}

// Withdraw: pointer receiver, it modifies the account.
func (a *BankAccount) Withdraw(amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	if amount > a.balance {
		return fmt.Errorf("insufficient funds: have %.2f, need %.2f", a.balance, amount)
	}
	a.balance -= amount
	return nil
}

// Balance: pointer receiver too. It only reads, so a value receiver would work and would be cheap
// (the struct is a string + float64). But Deposit and Withdraw already need pointers, and the Go FAQ
// consistency rule says: if some methods need a pointer receiver, make all of them pointer receivers.
// Mixing them is confusing and changes the method sets (only *BankAccount would have everything).
// Also a bank account is a thing with identity/state, not a small value like Point, and later it
// could get a sync.Mutex, which must never be copied.
func (a *BankAccount) Balance() float64 { return a.balance }

func main() {
	accounts := []BankAccount{
		{Owner: "Ali"},
		{Owner: "Bota"},
		{Owner: "Dana"},
	}
	accounts[0].Deposit(100)
	accounts[1].Deposit(200)
	accounts[2].Deposit(300)

	if err := accounts[0].Withdraw(500); err != nil {
		fmt.Println("error:", err) // error: insufficient funds: have 100.00, need 500.00
	}
	_ = accounts[1].Withdraw(50) // Bota: 150

	// 5% interest. Index loop on purpose: `for _, acc := range accounts` would give COPIES
	// and the balances in the slice would stay the same.
	for i := range accounts {
		accounts[i].Deposit(accounts[i].Balance() * 0.05)
	}

	for i := range accounts {
		fmt.Printf("%s: %.2f\n", accounts[i].Owner, accounts[i].Balance())
	}
	// Ali: 105.00
	// Bota: 157.50
	// Dana: 315.00
}
