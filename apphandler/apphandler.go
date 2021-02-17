package apphandler

import (
	"btcusd_server/datahandler"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
)

func LatestPrice() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := "1"
		bRate, err := datahandler.FetchLast(limit)
		if err != nil {
			// log.Fatal(err)
			// fmt.Println(err.Error())
			message := fmt.Sprintf("Request Failed - %v", err.Error())
			w.Write([]byte(message))
		} else {
			if len(bRate) == 0 {
				w.Write([]byte("Price record unavailable"))
			} else {
				json.NewEncoder(w).Encode(bRate)
			}
		}

	}
}

func PriceByTime() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		thedatetime := chi.URLParam(r, "thedatetime")
		message := ""
		thetime, err := time.Parse("2006-01-02 15:04:05", thedatetime)
		if err != nil {
			message = fmt.Sprintf("Request Failed - %v", "Price records not available")
		} else {
			// message := fmt.Sprintf("Time is %v", thetime.Format("2006-01-02 15:04:05"))
			// w.Write([]byte(message))
			// w.Write([]byte(tenminutesago))
			thetimestr := thetime.Format("2006-01-02 15:04:05")
			bRate, err := datahandler.FetchTime(thetimestr)
			if err != nil {
				tenminutesago := thetime.Add(-10 * time.Minute).Format("2006-01-02 15:04:05")
				bRate, err = datahandler.FetchRecent(tenminutesago, thetimestr)
				if err != nil {
					message = fmt.Sprintf("Request Failed - %v", "Price records not available")
				} else {
					rate := datahandler.Rate{}
					json.Unmarshal(bRate, &rate)
					message = fmt.Sprintf("The Recent price at %v is %v USD", rate.TransactionTime, rate.Price)
				}
			} else {
				rate := datahandler.Rate{}
				json.Unmarshal(bRate, &rate)

				message = fmt.Sprintf("The price of BTC at %s is %s USD", thetime, rate.Price)
				// json.NewEncoder(w).Encode(rate)
				// w.Write([]byte(message))
			}

		}

		w.Write([]byte(message))

	}
}

func PriceByRange() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		thedatetime := chi.URLParam(r, "timestamp1")
		thedatetime2 := chi.URLParam(r, "timestamp2")
		message := ""
		thetime1, err := time.Parse("2006-01-02 15:04:05", thedatetime)
		if err != nil {
			message = fmt.Sprintf("First time parameter invalid - %v", err.Error())
		} else {

			thetime2, err := time.Parse("2006-01-02 15:04:05", thedatetime2)
			if err != nil {
				message = fmt.Sprintf("Second time parameter invalid - %v", err.Error())
			} else {
				thetime1str := thetime1.Format("2006-01-02 15:04:05")
				thetime2str := thetime2.Format("2006-01-02 15:04:05")

				price, err := datahandler.FetchAverage(thetime1str, thetime2str)
				if err != nil {
					message = fmt.Sprintf("Request Failed - %v", err.Error())
				} else {
					message = fmt.Sprintf("The average price between %v and %v is %v USD", thetime1str, thetime2str, price)
				}
			}

		}

		w.Write([]byte(message))

	}
}

func Recents() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		limit := chi.URLParam(r, "limit")

		_, err := strconv.Atoi(limit)
		if err != nil {
			message := fmt.Sprintf("Request Failed - %v", err.Error())
			w.Write([]byte(message))
		}
		bRate, err := datahandler.FetchLast(limit)
		if err != nil {
			// log.Fatal(err)
			// fmt.Println(err.Error())
			message := fmt.Sprintf("Request Failed - %v", err.Error())
			w.Write([]byte(message))
		} else {
			// rate := datahandler.Rate{}
			// json.Unmarshal(bRate, &rate)
			if len(bRate) == 0 {
				w.Write([]byte("Price records unavailable"))
			} else {
				json.NewEncoder(w).Encode(bRate)
			}
		}

	}
}

func History() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		bRate, err := datahandler.FetchAll()
		if err != nil {
			message := fmt.Sprintf("Request Failed - %v", err.Error())
			w.Write([]byte(message))
		} else {
			if len(bRate) == 0 {
				w.Write([]byte("Price records not available"))
			} else {
				json.NewEncoder(w).Encode(bRate)
			}
		}

	}
}

func PriceToday() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		hour, minute, second := chi.URLParam(r, "hour"), chi.URLParam(r, "minute"), chi.URLParam(r, "second")

		valid, err := IsTimeValid(hour, minute, second)
		if err != nil {
			// log.Fatalln(err.Error())
			w.Write([]byte(valid))
		} else {

			if valid != "" {
				w.Write([]byte(valid))
			} else {

				today := time.Now().Truncate(time.Hour).Format("2006-01-02")

				timestamp := fmt.Sprintf("%s %s:%s:%s", today, hour, minute, second)
				fmt.Println("Param - ", timestamp)

				// w.Write([]byte("todays - " + timestamp))

				bRate, err := datahandler.FetchTime(timestamp)
				if err != nil {
					message := fmt.Sprintf("Request Failed - %v", "Price records not available")
					w.Write([]byte(message))
				} else {
					rate := datahandler.Rate{}
					json.Unmarshal(bRate, &rate)

					message := fmt.Sprintf("The price of BTC at %s is %s USD", timestamp, rate.Price)
					// json.NewEncoder(w).Encode(rate)
					w.Write([]byte(message))
				}

			}
		}

	}
}

func PriceRangeToday() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		timestamp1, timestamp2 := "", ""
		hour, minute, second := chi.URLParam(r, "hour"), chi.URLParam(r, "minute"), chi.URLParam(r, "second")
		valid, err := IsTimeValid(hour, minute, second)
		if err != nil {
			message := fmt.Sprintf("Request Failed - %v", err.Error())
			w.Write([]byte(message))
		} else {
			if valid != "" {
				w.Write([]byte(valid))
			} else {
				today := time.Now().Truncate(time.Hour).Format("2006-01-02")
				timestamp1 = fmt.Sprintf("%s %s:%s:%s", today, hour, minute, second)
				fmt.Println("Timestamp1 - ", timestamp1)
			}
		}

		hour, minute, second = chi.URLParam(r, "hour2"), chi.URLParam(r, "minute2"), chi.URLParam(r, "second2")
		valid, err = IsTimeValid(hour, minute, second)
		if err != nil {
			message := fmt.Sprintf("Request Failed - %v", err.Error())
			w.Write([]byte(message))
		} else {
			if valid != "" {
				w.Write([]byte(valid))
			} else {
				today := time.Now().Truncate(time.Hour).Format("2006-01-02")
				timestamp2 = fmt.Sprintf("%s %s:%s:%s", today, hour, minute, second)
				fmt.Println("Timestamp2 - ", timestamp2)
			}
		}

		if timestamp1 == "" || timestamp2 == "" {
			w.Write([]byte("Something is wrong. One of the timestamps may be invalid"))
		} else {
			bRate, err := datahandler.FetchRange(timestamp1, timestamp2)
			if err != nil {
				message := fmt.Sprintf("Request Failed - %v", "Price records not available")
				w.Write([]byte(message))
			} else {
				if len(bRate) == 0 {
					w.Write([]byte("Price records not available"))
				} else {
					json.NewEncoder(w).Encode(bRate)
				}
			}
		}

	}
}

func ReformatParam(iHour, iMin int) string {
	hour, minute := string(iHour), string(iMin)
	if len(hour) < 2 {
		hour = "0" + hour
	}

	if len(minute) < 2 {
		minute = "0" + minute
	}

	return hour + ":" + minute
}

func IsTimeValid(hour, minute, second string) (string, error) {
	if hour == "" || minute == "" || second == "" || len(hour) != 2 || len(minute) != 2 || len(second) != 2 {
		// log.Fatalln("Invalid Request")
		return "Invalid Request- The timestamp paameters is invalid. Please use this time format (hh-mm-ss)", nil
	} else {

		iHour, err := strconv.Atoi(hour)
		if err != nil {
			// log.Fatalln(err.Error())
			return "Request Failed - The hours component of the timestamp parameters is invalid", err
		}
		if iHour > 24 {
			return "Request Failed - The hours component of the timestamp parameters is invalid", nil
		}

		iMin, err := strconv.Atoi(minute)
		if err != nil {
			// log.Fatalln(err.Error())
			return "Request Failed - The minutes component of the timestamp parameters is invalid", err
		}
		if iMin > 60 {
			return "Request Failed - The minutes component of the timestamp parameters is invalid", nil
		}

		iSec, err := strconv.Atoi(second)
		if err != nil {
			// log.Fatalln(err.Error())
			return "Request Failed - The seconds component of the timestamp parameters is invalid", err
		}
		if iSec > 60 {
			return "Request Failed - The seconds component of the timestamp parameters is invalid", nil
		}

		return "", nil
	}
}
