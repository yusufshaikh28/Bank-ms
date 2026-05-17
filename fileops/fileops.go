package fileops

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

// CreateAccount inserts a new account
func CreateAccount(ctx context.Context, accountNumber int, pin string, balance float64) error {
	if AccountCollection == nil {
		return errors.New("MongoDB not connected")
	}

	count, _ := AccountCollection.CountDocuments(ctx, bson.M{"accountNumber": accountNumber})
	if count > 0 {
		return fmt.Errorf("account already exists")
	}

	_, err := AccountCollection.InsertOne(ctx, bson.M{
		"accountNumber": accountNumber,
		"pin":           pin,
		"balance":       balance,
	})
	return err
}

// ValidateLogin checks if account number + PIN is correct
func ValidateLogin(ctx context.Context, accountNumber int, pin string) (bool, error) {
	if AccountCollection == nil {
		return false, errors.New("MongoDB not connected")
	}

	var result struct {
		Pin string `bson:"pin"`
	}

	err := AccountCollection.FindOne(ctx, bson.M{"accountNumber": accountNumber}).Decode(&result)
	if err != nil {
		return false, err
	}

	return result.Pin == pin, nil
}

// GetBalance retrieves account balance
func GetBalance(ctx context.Context, accountNumber int, pin string) (float64, error) {
	valid, err := ValidateLogin(ctx, accountNumber, pin)
	if err != nil || !valid {
		return 0, fmt.Errorf("invalid credentials")
	}

	var result struct {
		Balance float64 `bson:"balance"`
	}
	err = AccountCollection.FindOne(ctx, bson.M{"accountNumber": accountNumber}).Decode(&result)
	if err != nil {
		return 0, err
	}

	return result.Balance, nil
}

// UpdateBalance adds/subtracts money from account
func UpdateBalance(ctx context.Context, accountNumber int, pin string, amount float64) error {
	valid, err := ValidateLogin(ctx, accountNumber, pin)
	if err != nil || !valid {
		return fmt.Errorf("invalid credentials")
	}

	_, err = AccountCollection.UpdateOne(ctx,
		bson.M{"accountNumber": accountNumber},
		bson.M{"$inc": bson.M{"balance": amount}},
	)
	return err
}
