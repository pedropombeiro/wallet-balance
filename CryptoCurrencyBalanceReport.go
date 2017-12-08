package main

import (
	"github.com/PombeirP/wallet-balance/fetchers"
)

// cryptoCurrencyTickerSymbol represents the ticker symbol for a crypto-currency
type cryptoCurrencyTickerSymbol string

const (
	btc  cryptoCurrencyTickerSymbol = "BTC"
	eth  cryptoCurrencyTickerSymbol = "ETH"
	ltc  cryptoCurrencyTickerSymbol = "LTC"
	dash cryptoCurrencyTickerSymbol = "DASH"
	uno  cryptoCurrencyTickerSymbol = "UNO"
	bcc  cryptoCurrencyTickerSymbol = "BCC"
)

// CryptoCurrencyBalanceReport provides functionality to check for the aggregate balance of crypto-currency addresses
type CryptoCurrencyBalanceReport struct {
	Symbol          cryptoCurrencyTickerSymbol
	UsdExchangeRate float64
	Balance         float64
	Error           error
}

// NewCryptoCurrencyBalanceReport creates a crypto-currency balance report instance for the given crypto-currency
func NewCryptoCurrencyBalanceReport(symbol cryptoCurrencyTickerSymbol, balance, exchangeRate float64, err error) *CryptoCurrencyBalanceReport {
	return &CryptoCurrencyBalanceReport{Symbol: symbol, Balance: balance, UsdExchangeRate: exchangeRate, Error: err}
}

// FetchInfoForCryptoCurrency retrieves the exchange rate and the aggregate balances for the provided addresses
func FetchInfoForCryptoCurrency(config *cryptoBalanceCheckerConfig, infoFetcher fetchers.CryptoCurrencyInfoFetcher, done chan<- *CryptoCurrencyBalanceReport) {
	balancesFetched := make(chan bool)
	exchangeRateFetched := make(chan bool)

	var report *CryptoCurrencyBalanceReport
	var balance, usdExchangeRate float64
	var err1, err2 error

	go infoFetcher.FetchBalance(config.Addresses, config.APIKey, &balance, &err1, balancesFetched)
	go infoFetcher.FetchExchangeRate(config.APIKey, "usd", &usdExchangeRate, &err2, exchangeRateFetched)

	<-exchangeRateFetched
	<-balancesFetched

	err := err1
	if err2 != nil && err1 == nil {
		err = err2
	}
	report = NewCryptoCurrencyBalanceReport(config.Symbol, balance, usdExchangeRate, err)

	done <- report
}