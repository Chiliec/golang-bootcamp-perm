package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/Chiliec/golang-bootcamp-perm/models"
	"golang.org/x/net/html/charset"
	"net/http"
	"strconv"
	"strings"
)

const (
	defaultCurrency string  = "HKD"
	defaultValue    float64 = 1000
)

var (
	currency string
	value    float64
)

func init() {
	flag.StringVar(&currency, "currency", defaultCurrency, "3-letter currency symbol")
	flag.Float64Var(&value, "value", defaultValue, "amount of money at that currency")
	flag.Parse()
}

func main() {
	valCurs := getExchangeRates()
	for _, valute := range valCurs.Valutes {
		if valute.CharCode == currency {
			valueStringWithoutComma := strings.Replace(valute.Value, ",", ".", -1)
			valueInFloat, err := strconv.ParseFloat(valueStringWithoutComma, 64)
			checkErr(err)
			valueInRubles := valueInFloat * value
			fmt.Printf("За %f %s сегодня дают %.2f рублей", value, valute.Name, valueInRubles)
			return
		}
	}
	fmt.Println("Нет такой валюты или что-то пошло не так")
}

func getExchangeRates() models.ValCurs {
	resp, err := http.Get("http://www.cbr.ru/scripts/XML_daily.asp")
	checkErr(err)
	defer resp.Body.Close()
	var valCurs models.ValCurs
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&valCurs)
	checkErr(err)
	return valCurs
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
