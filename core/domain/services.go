package domain

import (
	"digital-marketplace/config"
	"digital-marketplace/core"
	"digital-marketplace/core/infrastructure/mongo"
	"digital-marketplace/core/infrastructure/redis"
	"digital-marketplace/core/utils/coingecko"
	"log"
	)

type Domain struct {
	config    *config.AppConfig
	mongo     *mongo.MongoDB
	coingecko *coingecko.Coingecko
	redis     *redis.RedisRepository
}

func NewDomain(config *config.AppConfig, mongo *mongo.MongoDB, coingecko *coingecko.Coingecko, redis *redis.RedisRepository) *Domain {
	return &Domain{
		config:    config,
		mongo:     mongo,
		coingecko: coingecko,
		redis:     redis,
	}
}

func (d *Domain) GetInventory(tokenId string) ([]core.InventoryItem, error) {
	result, err := d.mongo.GetInventory()
	if err != nil {
		log.Println("Error: ", err)
	}
	tokenPrice, err := d.coingecko.GetTokenPrice(tokenId)
	if err != nil {
		log.Println("Error: ", err)
		return []core.InventoryItem{},err
	} 
	for i, data := range result {
		inTokenPrice := data.Price / tokenPrice
		tokenSymbol, _ := d.redis.GetTokenData(tokenId)
		log.Printf("Item Price is %f %s", inTokenPrice, tokenSymbol)
		result[i].Price = inTokenPrice
	}
	return result, nil

}

func (d *Domain) ExecuteOrder(user_id string, item_id string, quantity int32) (*core.ExecuteOrderResponse, error) {
	res, err := d.mongo.PurchaseItem(user_id, item_id, quantity)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return res, nil
}
