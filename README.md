# btcusd_server
Fetching BTCUSD rates every minutes

# README

Steps to query BTCUSD Rates server

## Installation
Extract the contents of the compresses application file (btcusd_server.zip)


```bash
cd btcusd_server

docker build -t test-server .
docker run --rm -p 80:80 test-server
```

## Example


```bash
curl -sS -X GET -k "http://localhost:80/latestprice"

{"price":"48689.9","transactiontime":"2021-02-17 18:05:07"}
```

```bash
curl -sS -X GET -k "http://localhost:80/pricebytime/2021-02-17 18:46:44"

The price of BTC at 2021-02-17 18:46:44 is 48529.9 USD
```
If the exact does not exist, the price at the closest recent time before the searched time is returned
```bash
curl -sS -X GET -k "http://localhost:80/pricebytime/2021-02-16 22:20:01"

The Recent price at 2021-02-16T22:19:04+08:00 is 49058.5 USD
```
```bash
curl -sS -X GET -k "http://localhost:80/pricebyrange/2021-02-16 21:06:53/2021-02-16 22:19:04"

The average price between 2021-02-16 21:06:53 and 2021-02-16 22:19:04 is 49135.04 USD
```
##
##

## Other Examples

Recent n prices and time
```bash
curl -sS -X GET -k "http://localhost:80/recents/5"

[{"price":"50972","transactiontime":"2021-02-17 22:52:59"}
{"price":"50898.9","transactiontime":"2021-02-17 22:51:59"},
{"price":"50891.7","transactiontime":"2021-02-17 22:50:43"},
{"price":"50852","transactiontime":"2021-02-17 22:44:29"},
{"price":"50875","transactiontime":"2021-02-17 22:42:22"}]
```

Historical price records in the database
```bash
curl -sS -X GET -k "http://localhost:80/history"

[{"price":"50848.6","transactiontime":"2021-02-17 21:41:02"},
{"price":"51159.5","transactiontime":"2021-02-17 22:19:37"},
{"price":"50875","transactiontime":"2021-02-17 22:42:22"},
{"price":"50852","transactiontime":"2021-02-17 22:44:29"},
{"price":"50891.7","transactiontime":"2021-02-17 22:50:43"},
{"price":"50898.9","transactiontime":"2021-02-17 22:51:59"}]
```
Price at a particular time today
```bash
curl -sS -X GET -k "http://localhost:80/pricetoday/22:20:01"

The prices at a time range today 2021-02-16T22:19:04+08:00 is 49058.5 USD
```
```bash
curl -sS -X GET -k "http://localhost:80/pricerangetoday/21:06:53/22:19:04"

[{"price":"50898.9","transactiontime":"2021-02-17 22:51:59"},
{"price":"50972","transactiontime":"2021-02-17 22:52:59"},
{"price":"50945.2","transactiontime":"2021-02-17 22:54:00"},
{"price":"50806.4","transactiontime":"2021-02-17 22:55:01"},
{"price":"51032.8","transactiontime":"2021-02-17 22:56:01"},
{"price":"51032.8","transactiontime":"2021-02-17 22:56:23"},
{"price":"51003.7","transactiontime":"2021-02-17 22:57:39"},
{"price":"51003.7","transactiontime":"2021-02-17 22:58:39"}]
```
