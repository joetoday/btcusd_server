package main

import (
	"btcusd_server/apphandler"
	"btcusd_server/datahandler"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"

	_ "github.com/mattn/go-sqlite3"
)

var (
	message string
)

type Exchange struct {
	Currency1 string `json:"curr1"`
	Currency2 string `json:"curr2"`
	Price     string `json:"lprice"`
}

func connectAndInitializeDB() {
	db, err := sql.Open("sqlite3", "./BtcUsd.db")
	if err != nil {
		message = fmt.Sprintf("Connection Failed - %v", err.Error())
		fmt.Println(message)
	}
	statement, err := db.Prepare("CREATE TABLE IF NOT EXISTS Rates (_id INTEGER PRIMARY KEY AUTOINCREMENT, price TEXT, fetchtime INT)")
	if err != nil {
		message = fmt.Sprintf("Table creation failed - %v", err)
		fmt.Println(message)
	}
	_, err = statement.Exec()
	if err != nil {
		message := fmt.Sprintf("Statement.Exec Failed connectAndInitializeDB() - %v", err.Error())
		fmt.Println(message)
	}

	fmt.Println("Database connection successful")
}

func GetCurrentRate() (Exchange, error) {
	exch := Exchange{}

	resp, err := http.Get("https://cex.io/api/last_price/BTC/USD")
	if err != nil {
		return exch, err
	}

	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return exch, err
	}

	// Convert the body to string
	sb := []byte(string(body))

	json.Unmarshal(sb, &exch)
	return exch, nil
}

func main() {
	today := time.Now()
	fmt.Printf("\nWelcome to BTCUSD @ %v !\n\n", today)

	connectAndInitializeDB()

	// defer db.Close()

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	r.Use(render.SetContentType(render.ContentTypeJSON))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(([]byte("Welcome to BTCUSD exchange rate")))
	})

	r.Get("/latestprice", apphandler.LatestPrice())
	r.Get("/pricebytime/{thedatetime}", apphandler.PriceByTime())
	r.Get("/pricebyrange/{timestamp1}/{timestamp2}", apphandler.PriceByRange())

	r.Get("/recents/{limit}", apphandler.Recents())
	r.Get("/history", apphandler.History())
	r.Get("/pricetoday/{hour}:{minute}:{second}", apphandler.PriceToday())
	r.Get("/pricerangetoday/{hour}:{minute}:{second}/{hour2}:{minute2}:{second2}", apphandler.PriceRangeToday())

	go func() {
		// run this indefinitely in the background
		for {
			currentTime := time.Now().Truncate(time.Second).Format("2006-01-02 15:04:05.999999999")
			// fmt.Printf("%v :  A minute tick\n", currentTime)

			exch, err := GetCurrentRate()
			if err != nil {
				log.Println(err)
			}
			fmt.Printf("BTCUSD rate @ %v  =>  %v\n", currentTime, exch)

			datahandler.StoreRate(exch.Price, currentTime)

			time.Sleep((60 * time.Second))
		}
	}()

	port := ":8080" // localhost
	// port := ":80" // dockerhost
	fmt.Println("Serving on " + port)
	http.ListenAndServe(port, r)
}
