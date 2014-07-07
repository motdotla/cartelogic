package cartelogic_test

import (
	"../cartelogic"
	"github.com/stvp/tempredis"
	"log"
	"testing"
)

const (
	EMAIL     = "account0@mailinator.com"
	API_KEY   = "custom_api_key"
	FRONT     = "<img src='http://example.com/some-image.png'>"
	BACK      = "John Doe"
	REDIS_URL = "redis://127.0.0.1:11001"
)

func tempredisConfig() tempredis.Config {
	config := tempredis.Config{
		"port":      "11001",
		"databases": "1",
	}
	return config
}

func TestAccountsCreate(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		account := map[string]interface{}{"email": EMAIL}

		cartelogic.Setup(REDIS_URL)
		result, logic_error := cartelogic.AccountsCreate(account)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}
		if result["email"] != EMAIL {
			t.Errorf("Incorrect email " + result["email"].(string))
		}

		if result["api_key"] == nil {
			t.Errorf("api_key is nil and should not be.")
		}
	})
}

func TestAccountsCreateCustomApiKey(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		account := map[string]interface{}{"email": EMAIL, "api_key": API_KEY}

		cartelogic.Setup(REDIS_URL)
		result, logic_error := cartelogic.AccountsCreate(account)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}
		if result["api_key"] != API_KEY {
			t.Errorf("api_key did not equal " + API_KEY)
		}
		exists, _ := cartelogic.KeyExists("accounts/" + API_KEY)
		if exists != true {
			t.Errorf("pointer key using api_key did not exist")
		}
	})
}

func TestAccountsCreateCustomBlankApiKey(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		account := map[string]interface{}{"email": EMAIL, "api_key": ""}

		cartelogic.Setup(REDIS_URL)
		result, logic_error := cartelogic.AccountsCreate(account)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}

		if result["api_key"] == nil || result["api_key"].(string) == "" {
			t.Errorf("It should generate an api_key if blank.")
		}
	})
}

func TestAccountsCreateNilEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		account := map[string]interface{}{}

		cartelogic.Setup(REDIS_URL)
		_, logic_error := cartelogic.AccountsCreate(account)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestAccountsCreateBlankEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		account := map[string]interface{}{"email": ""}

		cartelogic.Setup(REDIS_URL)
		_, logic_error := cartelogic.AccountsCreate(account)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestCardsCreate(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}
		setupAccount(t)

		card := map[string]interface{}{"front": FRONT, "back": BACK, "api_key": API_KEY}
		result, logic_error := cartelogic.CardsCreate(card)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}
		if result["id"] == nil {
			t.Errorf("Error", result)
		}
	})
}

func setupAccount(t *testing.T) {
	account := map[string]interface{}{"email": EMAIL, "api_key": API_KEY}

	cartelogic.Setup(REDIS_URL)
	_, logic_error := cartelogic.AccountsCreate(account)
	if logic_error != nil {
		t.Errorf("Error", logic_error)
	}
}
