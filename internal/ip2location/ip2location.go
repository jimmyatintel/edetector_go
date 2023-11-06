package ip2location

import (
	"edetector_go/config"

	"github.com/ip2location/ip2location-go/v9"
)

func ToCountry(ip string) (string, error) {
	db, err := ip2location.OpenDB(config.Viper.GetString("IP2LOCATION_PATH"))
	if err != nil {
		return "", err
	}
	defer db.Close()
	code, err := db.Get_country_short(ip)
	if err != nil {
		return "", err
	}
	return code.Country_short, nil
}

func ToCountries(ips []string) ([]string, error) {
	db, err := ip2location.OpenDB(config.Viper.GetString("IP2LOCATION_PATH"))
	if err != nil {
		return nil, err
	}
	defer db.Close()
	country_codes := []string{}
	for _, ip := range ips {
		code, err := db.Get_country_short(ip)
		if err != nil {
			return nil, err
		}
		country_codes = append(country_codes, code.Country_short)
	}
	return country_codes, nil
}
