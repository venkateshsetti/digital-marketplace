package coingecko

import (
	"digital-marketplace/config"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Coingecko struct {
	config *config.AppConfig
}

type CoinsList []struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func CoingeckoService(config *config.AppConfig) *Coingecko {
	return &Coingecko{
		config: config,
	}
}

func (c *Coingecko) GetCoinList() (*CoinsList, error) {
	var data CoinsList
	coinlistUrl := c.config.Coingecko.Url + "/coins/list"
	response, err := http.Get(coinlistUrl)
	if err != nil {
		log.Println("Error:", err)
		return nil, err
	}
	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return nil, err
	}
	log.Println("Successfully fetched the tokens list")
	return &data, nil
}

func (c *Coingecko) GetTokenPrice(tokenID string) (float64, error) {
	var data map[string]map[string]float64
	tokenPriceUrl := c.config.Coingecko.Url + "/simple/price?ids=" + tokenID + "&vs_currencies=inr"
	response, err := http.Get(tokenPriceUrl)
	if err != nil {
		log.Println("Error:", err)
		return 0.0, err
	}
	defer response.Body.Close()
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Println("Error reading response body:", err)
		return 0.0, err

	}
	log.Println(string(bodyBytes), tokenID)
	err = json.Unmarshal(bodyBytes, &data)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		return 0.0, err

	}
	return data[tokenID]["inr"], nil
}
