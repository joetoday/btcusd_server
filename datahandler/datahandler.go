package datahandler

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type Rate struct {
	Price           string `json:"price"`
	TransactionTime string `json:"transactiontime"`
}

func StoreRate(price string, fetchdate string) {

	db, err := sql.Open("sqlite3", "./BtcUsd.db")
	if err != nil {
		message := fmt.Sprintf("Connection Failed StoreRate() - %v", err.Error())
		fmt.Println(message)
	}

	// Prepare statement for inserting data
	stmtIns, err := db.Prepare("INSERT INTO Rates (price, fetchtime) VALUES ( ?, ? )") // ? = placeholder
	if err != nil {
		message := fmt.Sprintf("Insert Failed StoreRate() - %v", err.Error())
		fmt.Println(message)
	}
	res, err := stmtIns.Exec(price, fetchdate)
	// defer stmtIns.Close()
	if err != nil {
		message := fmt.Sprintf("Statement.Exec Failed StoreRate() - %v", err.Error())
		fmt.Println(message)
	}

	rowCnt, err := res.RowsAffected()
	if err != nil {
		message := fmt.Sprintf("rowCnt StoreRate() - %v", err.Error())
		fmt.Println(message)
	}

	if rowCnt > 0 {
		fmt.Printf("BTCUSD rate saved successful!!!\n")
	}

}

func FetchLast(limit string) ([]Rate, error) {

	db, err := sql.Open("sqlite3", "./BtcUsd.db")
	if err != nil {
		message := fmt.Sprintf("Connection Failed FetchLast() - %v", err.Error())
		fmt.Println(message)
	}

	var price, fetchtime string
	var prices []Rate

	// query := "SELECT price, fetchdate from Rate ORDER BY _id DESC LIMIT" + string(limit)
	rows, err := db.Query("SELECT price, fetchtime from Rates ORDER BY _id DESC LIMIT ?", limit)
	if err != nil {
		return prices, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&price, &fetchtime)
		if err != nil {
			return prices, err
		}
		prices = append(prices, Rate{Price: price, TransactionTime: fetchtime})
	}
	err = rows.Err()
	if err != nil {
		return prices, err
	}

	return prices, nil

}

func FetchTime(timestamp string) ([]byte, error) {

	db, err := sql.Open("sqlite3", "./BtcUsd.db")
	if err != nil {
		message := fmt.Sprintf("Connection Failed FetchTime() - %v", err.Error())
		fmt.Println(message)
	}

	var price string

	// query := "SELECT price, fetchdate from Rate ORDER BY _id DESC LIMIT" + string(limit)
	err = db.QueryRow("SELECT price from Rates WHERE fetchtime = ?", timestamp).Scan(&price)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(Rate{Price: price, TransactionTime: timestamp})

	if err != nil {
		return nil, err
	}

	// fmt.Println(string(data))
	return data, nil

}

func FetchAll() ([]Rate, error) {

	db, err := sql.Open("sqlite3", "./BtcUsd.db")
	if err != nil {
		message := fmt.Sprintf("Connection Failed FetchAll() - %v", err.Error())
		fmt.Println(message)
	}

	var price, fetchtime string
	var prices []Rate

	// query := "SELECT price, fetchdate from Rate ORDER BY _id DESC LIMIT" + string(limit)
	rows, err := db.Query("SELECT price, fetchtime from Rates")
	if err != nil {
		return prices, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&price, &fetchtime)
		if err != nil {
			return prices, err
		}
		prices = append(prices, Rate{Price: price, TransactionTime: fetchtime})
	}
	err = rows.Err()
	if err != nil {
		return prices, err
	}

	return prices, nil

}

func FetchRange(timestamp1, timestamp2 string) ([]Rate, error) {

	db, err := sql.Open("sqlite3", "./BtcUsd.db")
	if err != nil {
		message := fmt.Sprintf("Connection Failed FetchRange() - %v", err.Error())
		fmt.Println(message)
	}

	var price, timestamp string
	var prices []Rate

	// query := "SELECT price, fetchdate from Rate ORDER BY _id DESC LIMIT" + string(limit)
	rows, err := db.Query("SELECT price, fetchtime FROM Rates WHERE  fetchtime BETWEEN ? AND ?", timestamp1, timestamp2)
	if err != nil {
		return prices, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&price, &timestamp)
		if err != nil {
			return prices, err
		}
		prices = append(prices, Rate{Price: price, TransactionTime: timestamp})
	}
	err = rows.Err()
	if err != nil {
		return prices, err
	}

	return prices, nil

}

func FetchAverage(timestamp1, timestamp2 string) (string, error) {

	db, err := sql.Open("sqlite3", "./BtcUsd.db")
	if err != nil {
		message := fmt.Sprintf("Connection Failed FetchAverage() - %v", err.Error())
		fmt.Println(message)
	}

	var price string

	// query := "SELECT price, fetchdate from Rate ORDER BY _id DESC LIMIT" + string(limit)
	rows, err := db.Query("SELECT price FROM Rates WHERE  fetchtime BETWEEN ? AND ?", timestamp1, timestamp2)
	if err != nil {
		return "", err
	}
	defer rows.Close()

	// var total, count float64
	total, count := 0.0, 0.0
	for rows.Next() {
		count += 1.0
		err := rows.Scan(&price)
		if err != nil {
			return "", err
		}
		p, _ := strconv.ParseFloat(price, 32)
		total += p
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%.2f", total/count), nil

}

func FetchRecent(timestamp1, timestamp2 string) ([]byte, error) {

	db, err := sql.Open("sqlite3", "./BtcUsd.db")
	if err != nil {
		message := fmt.Sprintf("Connection Failed FetchRecent() - %v", err.Error())
		fmt.Println(message)
	}

	var price, fetchtime string

	// query := "SELECT price, fetchdate from Rate ORDER BY _id DESC LIMIT" + string(limit)
	err = db.QueryRow("SELECT price, fetchtime FROM Rates WHERE  fetchtime BETWEEN ? AND ? ORDER BY _id DESC LIMIT 1", timestamp1, timestamp2).Scan(&price, &fetchtime)
	if err != nil {
		return nil, err
	}

	data, err := json.Marshal(Rate{Price: price, TransactionTime: fetchtime})

	if err != nil {
		return nil, err
	}

	fmt.Println(string(data))
	return data, nil
}
