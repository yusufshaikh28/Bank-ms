package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"example.com/bank/fileops"
)

func main() {
	// Connect MongoDB
	fileops.ConnectMongo()

	// Serve frontend
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// API routes
	http.HandleFunc("/create", createAccount)
	http.HandleFunc("/login", login)
	http.HandleFunc("/balance", getBalance)
	http.HandleFunc("/deposit", depositMoney)
	http.HandleFunc("/withdraw", withdrawMoney)

	fmt.Println("🚀 Server running at http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// --- Handlers ---

func createAccount(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountNumber int     `json:"accountNumber"`
		Pin           string  `json:"pin"`
		Balance       float64 `json:"balance"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := fileops.CreateAccount(ctx, req.AccountNumber, req.Pin, req.Balance)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "account created"})
}

func login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountNumber int    `json:"accountNumber"`
		Pin           string `json:"pin"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	valid, err := fileops.ValidateLogin(ctx, req.AccountNumber, req.Pin)
	if err != nil || !valid {
		http.Error(w, "invalid login", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "login successful"})
}

func getBalance(w http.ResponseWriter, r *http.Request) {
	accStr := r.URL.Query().Get("account")
	pin := r.URL.Query().Get("pin")

	if accStr == "" || pin == "" {
		http.Error(w, "missing account number or pin", http.StatusBadRequest)
		return
	}
	accountNumber, _ := strconv.Atoi(accStr)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	balance, err := fileops.GetBalance(ctx, accountNumber, pin)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	json.NewEncoder(w).Encode(map[string]float64{"balance": balance})
}

func depositMoney(w http.ResponseWriter, r *http.Request) {
	accStr := r.URL.Query().Get("account")
	pin := r.URL.Query().Get("pin")

	if accStr == "" || pin == "" {
		http.Error(w, "missing account number or pin", http.StatusBadRequest)
		return
	}
	accountNumber, _ := strconv.Atoi(accStr)

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "invalid deposit", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := fileops.UpdateBalance(ctx, accountNumber, pin, req.Amount); err != nil {
		http.Error(w, "deposit failed", http.StatusUnauthorized)
		return
	}

	getBalance(w, r)
}

func withdrawMoney(w http.ResponseWriter, r *http.Request) {
	accStr := r.URL.Query().Get("account")
	pin := r.URL.Query().Get("pin")

	if accStr == "" || pin == "" {
		http.Error(w, "missing account number or pin", http.StatusBadRequest)
		return
	}
	accountNumber, _ := strconv.Atoi(accStr)

	var req struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "invalid withdrawal", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// check balance first
	balance, err := fileops.GetBalance(ctx, accountNumber, pin)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if req.Amount > balance {
		http.Error(w, "insufficient funds", http.StatusBadRequest)
		return
	}

	if err := fileops.UpdateBalance(ctx, accountNumber, pin, -req.Amount); err != nil {
		http.Error(w, "withdraw failed", http.StatusUnauthorized)
		return
	}

	getBalance(w, r)
}
