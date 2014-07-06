package cartelogic_test

import (
	"../cartelogic"
	"github.com/stvp/tempredis"
	"log"
	"testing"
)

const (
	NAME      = "People Deck"
	EMAIL     = "deck0@mailinator.com"
	API_KEY   = "custom_api_key"
	REDIS_URL = "redis://127.0.0.1:11001"
)

func tempredisConfig() tempredis.Config {
	config := tempredis.Config{
		"port":      "11001",
		"databases": "1",
	}
	return config
}

func TestDecksCreate(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		deck := map[string]interface{}{"email": EMAIL, "name": NAME}

		cartelogic.Setup(REDIS_URL)
		result, logic_error := cartelogic.DecksCreate(deck)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}
		if result["email"] != EMAIL {
			t.Errorf("Incorrect email " + result["email"].(string))
		}

		if result["name"] != NAME {
			t.Errorf("Incorrect name " + result["name"].(string))
		}
		if result["api_key"] == nil {
			t.Errorf("api_key is nil and should not be.")
		}
		if result["id"] == nil {
			t.Errorf("id is nil and should not be.")
		}
	})
}

func TestDecksCreateCustomApiKey(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		deck := map[string]interface{}{"email": EMAIL, "name": NAME, "api_key": API_KEY}

		cartelogic.Setup(REDIS_URL)
		result, logic_error := cartelogic.DecksCreate(deck)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}

		if result["api_key"] != API_KEY {
			t.Errorf("api_key did not equal " + API_KEY)
		}
	})
}

func TestDecksCreateCustomBlankApiKey(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		deck := map[string]interface{}{"email": EMAIL, "name": NAME, "api_key": ""}

		cartelogic.Setup(REDIS_URL)
		result, logic_error := cartelogic.DecksCreate(deck)
		if logic_error != nil {
			t.Errorf("Error", logic_error)
		}

		if result["api_key"] == nil || result["api_key"].(string) == "" {
			t.Errorf("It should generate an api_key if blank.")
		}
	})
}

func TestDecksCreateBlankName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		deck := map[string]interface{}{"email": EMAIL, "name": ""}

		cartelogic.Setup(REDIS_URL)
		_, logic_error := cartelogic.DecksCreate(deck)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestDecksCreateNilName(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		deck := map[string]interface{}{"email": EMAIL}

		cartelogic.Setup(REDIS_URL)
		_, logic_error := cartelogic.DecksCreate(deck)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestDecksCreateBlankEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		deck := map[string]interface{}{"email": "", "name": NAME}

		cartelogic.Setup(REDIS_URL)
		_, logic_error := cartelogic.DecksCreate(deck)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}

func TestDecksCreateNilEmail(t *testing.T) {
	tempredis.Temp(tempredisConfig(), func(err error) {
		if err != nil {
			log.Println(err)
		}

		deck := map[string]interface{}{"name": NAME}

		cartelogic.Setup(REDIS_URL)
		_, logic_error := cartelogic.DecksCreate(deck)
		if logic_error.Code != "required" {
			t.Errorf("Error", err)
		}
	})
}
