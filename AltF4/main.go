package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)
// Структура для хранения символа каждого стакана
type Coin struct {
	Symbol string `json:"symbol"`
	Price string `json:"price"`
}

// Структура для хранения данные БИД/АСК по стаканам
type Order struct {
	Symbol string `json:"symbol"`
	Id int `json:"lastUpdateId"`
	Total float64	`json:"total"`
	Bids [][]string `json:"bids"`
	Asks [][]string `json:"asks"`
}

// Структура для вывода
type Res struct {
	Orders []*Order `json:"orders"`
}

func main() {
	response, err := http.Get("https://api.binance.com/api/v3/ticker/price")
	if err != nil {
		log.Fatal(err)
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var coins []*Coin
	var ord *Order
	var orders []*Order
	var symbols[] string //хранение список символов чтобы итерировать по циклу
	json.Unmarshal(responseData, &coins)

	for i := 0; i < len(coins); i++{
		symbols = append(symbols, coins[i].Symbol)
	}
	// тут можно вывести данные по всем стаканам просто поменяв 5 на длину массива symbols
	for i := 0; i < 5; i++{
		s := "https://api.binance.com/api/v3/depth?symbol=" + symbols[i]
		response2, err := http.Get(s)
		if err != nil {
			log.Fatal(err)
		}
		responseData2, err := ioutil.ReadAll(response2.Body)
		if err != nil {
			log.Fatal(err)
		}
		json.Unmarshal(responseData2, &ord)

		//копирование начальных 15 ордеров
		bids := make([][]string, 15)
		copy(bids, ord.Bids)
		asks := make([][]string, 15)
		copy(asks, ord.Bids)

		total := 0.0
		for j := 0; j < 15; j++{
			//Вычисление суммарного объема
			bprice, err := strconv.ParseFloat(bids[j][0], 8)
			bquantity, err := strconv.ParseFloat(bids[j][1], 8)
			total += bquantity * bprice

			aprice, err := strconv.ParseFloat(asks[j][0], 8)
			aquantity, err := strconv.ParseFloat(asks[j][1], 8)
			total += aquantity * aprice

			if err != nil {
				log.Fatal(err)
			}
		}

		stakan := Order{symbols[i], ord.Id, total, bids, asks}
		orders = append(orders, &stakan)
	}
	res := Res{orders}
	answer, _ := json.MarshalIndent(res, "", "  ")
	fmt.Println(string(answer))
}