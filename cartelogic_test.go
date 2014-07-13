package cartelogic_test

import (
	"../cartelogic"
	"github.com/joho/godotenv"
	"github.com/orchestrate-io/gorc"
	"log"
	"os"
	"testing"
)

const (
	EMAIL   = "account0@mailinator.com"
	API_KEY = "custom_api_key"
	FRONT   = "<img src='http://example.com/some-image.png'>"
	BACK    = "John Doe"
)

func TestAccountsCreateCustomApiKey(t *testing.T) {
	setup(t)
	tearDown(t)

	account := map[string]interface{}{"email": EMAIL, "api_key": API_KEY}

	cartelogic.Setup(os.Getenv("ORCHESTRATE_API_KEY"))
	result, logic_error := cartelogic.AccountsCreate(account)
	if logic_error != nil {
		t.Errorf("Error", logic_error)
	}

	if result["api_key"] == nil || result["api_key"].(string) == "" {
		t.Errorf("It should generate an api_key if blank.")
	}
}

func TestAccountsCreateNilEmail(t *testing.T) {
	setup(t)

	account := map[string]interface{}{}

	cartelogic.Setup(os.Getenv("ORCHESTRATE_API_KEY"))
	_, logic_error := cartelogic.AccountsCreate(account)
	if logic_error.Code != "required" {
		t.Errorf("Error", "Logic error should have been 'required'")
	}
}

func TestAccountsCreateBlankEmail(t *testing.T) {
	setup(t)

	account := map[string]interface{}{"email": ""}

	cartelogic.Setup(os.Getenv("ORCHESTRATE_API_KEY"))
	_, logic_error := cartelogic.AccountsCreate(account)
	if logic_error.Code != "required" {
		t.Errorf("Error", "logic_error should have been 'required'")
	}
}

func TestCardsCreate(t *testing.T) {
	setup(t)
	setupAccount(t)

	card := map[string]interface{}{"front": FRONT, "back": BACK, "api_key": API_KEY}
	result, logic_error := cartelogic.CardsCreate(card)
	if logic_error != nil {
		t.Errorf("Error", logic_error)
	}
	if result["account_id"] == nil {
		t.Errorf("Error", result)
	}
}

func setup(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func tearDown(t *testing.T) {
	orchestrate_api_key := os.Getenv("ORCHESTRATE_API_KEY")
	o := gorc.NewClient(orchestrate_api_key)
	o.DeleteCollection("cards")
	o.DeleteCollection("accounts")
}

func setupAccount(t *testing.T) {
	account := map[string]interface{}{"email": EMAIL, "api_key": API_KEY}

	cartelogic.Setup(os.Getenv("ORCHESTRATE_API_KEY"))
	cartelogic.AccountsCreate(account)
}
