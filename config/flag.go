package config

import (
	"flag"
)

var (
	ListenAddress              string
	BankHolidayJsonEndpointURL string
)

func InitFlags() {
	flag.StringVar(&ListenAddress, "listenAddress", ":8080",
		"The listen address of REST HTTP server ")
	flag.StringVar(&BankHolidayJsonEndpointURL, "BankHolidayJsonEndpointURL", "https://www.gov.uk/bank-holidays.json",
		"The URL of the endpoint that hosts the bank holiday json data")
}
