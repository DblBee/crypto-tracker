package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var newUser = User{Name: "Bryan"}

var bitcoin_asset = Asset{
	Symbol: "BTC",
	Name:   "bitcoin",
}

var sol_asset = Asset{
	Symbol: "SOL",
	Name:   "solana",
}

var eth_asset = Asset{
	Symbol: "ETH",
	Name:   "ethereum",
}

type CoinGeckoResponse struct {
	Bitcoin struct {
		Usd int `json:"usd"`
	} `json:"bitcoin"`
	Solana struct {
		Usd float64 `json:"usd"`
	} `json:"solana"`
	Ethereum struct {
		Usd float64 `json:"usd"`
	} `json:"ethereum"`
}

type User struct {
	ID   int `gorm:"primary_key`
	Name string
}

type Asset struct {
	ID     int    `gorm:"primary_key`
	Symbol string `gorm:"unique;not null"`
	Name   string
}

type Transaction struct {
	ID       int       `gorm:"primaryKey;"`
	TS       time.Time `gorm:"primaryKey;"`
	UserId   int
	User     User `gorm:"foreignKey:UserId"`
	AssetId  int
	Asset    Asset `gorm:"foreignKey:AssetId"`
	Amount   uint
	PriceUsd float32
}

func seed_reference_data(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		db_err := tx.Save(&newUser).Error
		if db_err != nil {
			fmt.Print(db_err.Error())
		}

		db_err = tx.Save(&newUser).Error
		if db_err != nil {
			fmt.Print(db_err.Error())
		}

		db_err = tx.Save(&bitcoin_asset).Error
		if db_err != nil {
			fmt.Print(db_err.Error())
		}

		db_err = tx.Save(&sol_asset).Error
		if db_err != nil {
			fmt.Print(db_err.Error())
		}

		db_err = tx.Create(&eth_asset).Error
		if db_err != nil {
			fmt.Print(db_err.Error())
		}

		return nil
	})

	if err != nil {
		fmt.Print(err.Error())
	}
}

func get_asset_id_map(db *gorm.DB) (*[3]Asset, error) {
	assets_from_tiger := [3]Asset{bitcoin_asset, sol_asset, eth_asset}
	err := db.Find(&assets_from_tiger).Error
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	return &assets_from_tiger, nil
}

func get_user(db *gorm.DB) (*User, error) {
	var user = User{}
	err := db.First(&user).Error

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	fmt.Printf("Found User %v\n", user.ID)
	return &user, nil
}

func fecth_prices() (*CoinGeckoResponse, error) {
	fmt.Println("Calling Coin Gecko API...")
	// COIN_GECKO_API_KEY
	coin_gecko_api_key := os.Getenv("COIN_GECKO_API_KEY")

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://api.coingecko.com/api/v3/simple/price?vs_currencies=usd&ids=bitcoin,solana,ethereum", nil)

	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", coin_gecko_api_key)
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
		return nil, err
	}

	var responseObject CoinGeckoResponse
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response as struct %+v\n", responseObject)

	return &responseObject, nil
}

func ingest_once(db *gorm.DB) {
	var data, err = fecth_prices()
	if err != nil {
		panic(err.Error())
	}
	var user, err1 = get_user(db)
	if err1 != nil {
		panic(err1.Error())
	}

	ts := time.Now().UTC()

	var asset_id_by_symbol, err2 = get_asset_id_map(db)
	if err2 != nil {
		panic(err2.Error())
	}

	fmt.Println(data.Bitcoin.Usd)
	fmt.Println(data.Solana.Usd)
	fmt.Println(data.Ethereum.Usd)
	fmt.Println(user.Name)
	fmt.Println(user.ID)
	fmt.Println(asset_id_by_symbol[0])
	fmt.Println(asset_id_by_symbol[1])
	fmt.Println(asset_id_by_symbol[2])

	err3 := db.Transaction(func(tx *gorm.DB) error {
		btc_transaction := Transaction{
			TS:       ts,
			User:     *user,
			Asset:    asset_id_by_symbol[0],
			UserId:   user.ID,
			AssetId:  asset_id_by_symbol[0].ID,
			PriceUsd: float32(data.Bitcoin.Usd),
			Amount:   1,
			ID:       1,
		}

		tx.Create(&btc_transaction)

		sol_transaction := Transaction{
			TS:       ts,
			User:     *user,
			Asset:    asset_id_by_symbol[1],
			UserId:   user.ID,
			AssetId:  asset_id_by_symbol[1].ID,
			PriceUsd: float32(data.Solana.Usd),
			Amount:   1,
		}

		tx.Create(&sol_transaction)

		eth_transaction := Transaction{
			TS:       ts,
			User:     *user,
			Asset:    asset_id_by_symbol[2],
			UserId:   user.ID,
			AssetId:  asset_id_by_symbol[2].ID,
			PriceUsd: float32(data.Ethereum.Usd),
			Amount:   1,
		}
		tx.Create(&eth_transaction)

		return nil
	})

	if err3 != nil {
		panic(err3.Error())
	}

}

func main() {
	err := godotenv.Load()

	dsn := os.Getenv("TIMESCALE_SERVICE_URL")

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	err = db.AutoMigrate(&User{})
	err = db.AutoMigrate(&Asset{})
	err = db.AutoMigrate(&Transaction{})

	if err != nil {
		panic("failed to migrate database")
	}

	// create the hypertable
	tx := db.Exec("SELECT create_hypertable('transactions', 'ts')")

	if tx.Error != nil {
		fmt.Fprintf(os.Stderr, "Hypertable exists %v\n", tx.Error.Error())
	}

	seed_reference_data(db)

	for {
		ingest_once(db)
		time.Sleep(30 * time.Second)
	}
}
