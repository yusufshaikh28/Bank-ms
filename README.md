

# GoBank — ATM Banking System

A full-stack ATM banking web application built in Go with MongoDB and a clean frontend interface.

## Features

- **Account Creation** — create accounts with account number and PIN
- **Secure Login** — PIN-based authentication
- **Balance Enquiry** — check account balance
- **Deposit & Withdrawal** — with insufficient funds validation
- **MongoDB Backend** — persistent storage via MongoDB

## Tech Stack

- **Backend:** Go (Golang), net/http
- **Database:** MongoDB
- **Frontend:** HTML, CSS, JavaScript

## Running Locally

```bash
go run bank.go
```

Then open `localhost:8080`

## API Endpoints

- `POST /create` — create a new account
- `POST /login` — authenticate with account number and PIN
- `GET /balance` — get account balance
- `POST /deposit` — deposit money
- `POST /withdraw` — withdraw money

---

