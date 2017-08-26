package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/pkg/errors"
	"golang.org/x/net/html/charset"
	"log"
	"net/http"
	"strconv"
	"strings"
	
	"github.com/Chiliec/golang-bootcamp-perm/models"
)

const (
	defaultCurrency string  = "HKD"
	defaultValue    float64 = 1000
)

var (
	inputtedCurrency string
	inputtedValue    float64
)

func init() {
	flag.StringVar(&inputtedCurrency, "currency", defaultCurrency, "3-letter currency symbol")
	flag.Float64Var(&inputtedValue, "value", defaultValue, "amount of money at that currency")
	flag.Parse()
}

func main() {
	finalString, err := getFinalString()
	if err != nil {
		finalString = "К сожалению, произошла ошибка: " + err.Error()
	}
	log.Println(finalString)
}

func getFinalString() (string, error) {
	exchangeRates, err := getExchangeRates()
	if err != nil || exchangeRates == nil {
		return "Не получили курсы валют", err
	}
	for _, currency := range exchangeRates.Currencies {
		if currency.CharCode == inputtedCurrency {
			valueStringWithoutComma := strings.Replace(currency.Value, ",", ".", -1)
			valueInFloat, err := strconv.ParseFloat(valueStringWithoutComma, 64)
			if err != nil {
				return "Не смогли распарсить курс", err
			}
			valueInRubles := valueInFloat * inputtedValue
			return fmt.Sprintf("За %f %s сегодня дают %.2f рублей", inputtedValue, currency.Name, valueInRubles), nil
		}
	}
	return "", errors.New("Нет такой валюты или что-то пошло не так")
}

func getExchangeRates() (*models.ExchangeRates, error) {
	resp, err := http.Get("http://www.cbr.ru/scripts/XML_daily.asp")
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, errors.New("Пустое тело ответа!")
	}
	defer resp.Body.Close()
	var exchangeRates models.ExchangeRates
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&exchangeRates)
	if err != nil {
		return nil, err
	}
	return &exchangeRates, nil
}
