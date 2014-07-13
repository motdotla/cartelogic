package cartelogic

import (
	"code.google.com/p/go-uuid/uuid"
	"github.com/dchest/uniuri"
	"github.com/handshakejs/handshakejserrors"
	"github.com/orchestrate-io/gorc"
)

const (
	BASE_10  = 10
	ACCOUNTS = "accounts"
)

var (
	o *gorc.Client
)

func Setup(orchestrate_api_key string) {
	o = gorc.NewClient(orchestrate_api_key)
}

func AccountsCreate(account map[string]interface{}) (map[string]interface{}, *handshakejserrors.LogicError) {
	var email string
	if str, ok := account["email"].(string); ok {
		email = str
	} else {
		email = ""
	}
	if email == "" {
		logic_error := &handshakejserrors.LogicError{"required", "email", "email cannot be blank"}
		return nil, logic_error
	}
	account["email"] = email

	generated_api_key := uniuri.NewLen(20)
	if account["api_key"] == nil {
		account["api_key"] = generated_api_key
	}
	if account["api_key"].(string) == "" {
		account["api_key"] = generated_api_key
	}

	results, err := o.Search(ACCOUNTS, "email:"+account["email"].(string), 10, 0)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}
	if results.TotalCount > 0 {
		logic_error := &handshakejserrors.LogicError{"not_unique", "email", "email must be unique"}
		return nil, logic_error
	}

	key := uuid.New()
	_, errr := o.Put(ACCOUNTS, key, account)
	if errr != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", errr.Error()}
		return nil, logic_error
	}

	return account, nil
}

func CardsCreate(card map[string]interface{}) (map[string]interface{}, *handshakejserrors.LogicError) {
	front, logic_error := checkFrontPresent(card)
	if logic_error != nil {
		return nil, logic_error
	}
	card["front"] = front

	back, logic_error := checkBackPresent(card)
	if logic_error != nil {
		return nil, logic_error
	}
	card["back"] = back

	api_key, logic_error := checkApiKeyPresent(card)
	if logic_error != nil {
		return nil, logic_error
	}
	card["api_key"] = api_key

	results, err := o.Search(ACCOUNTS, "api_key:"+card["api_key"].(string), 10, 0)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}
	if results.TotalCount <= 0 {
		logic_error := &handshakejserrors.LogicError{"incorrect", "api_key", "the api_key is incorrect"}
		return nil, logic_error
	}
	account_id := results.Results[0].Path.Key
	card["account_id"] = account_id

	key := uuid.New()
	_, errr := o.Put("cards", key, card)
	if errr != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", errr.Error()}
		return nil, logic_error
	}

	return card, nil
}

func checkFrontPresent(card map[string]interface{}) (string, *handshakejserrors.LogicError) {
	var front string
	if str, ok := card["front"].(string); ok {
		front = str
	} else {
		front = ""
	}
	if front == "" {
		logic_error := &handshakejserrors.LogicError{"required", "front", "front cannot be blank"}
		return front, logic_error
	}

	return front, nil
}

func checkBackPresent(card map[string]interface{}) (string, *handshakejserrors.LogicError) {
	var back string
	if str, ok := card["back"].(string); ok {
		back = str
	} else {
		back = ""
	}
	if back == "" {
		logic_error := &handshakejserrors.LogicError{"required", "back", "back cannot be blank"}
		return back, logic_error
	}

	return back, nil
}

func checkApiKeyPresent(card map[string]interface{}) (string, *handshakejserrors.LogicError) {
	var api_key string
	if str, ok := card["api_key"].(string); ok {
		api_key = str
	} else {
		api_key = ""
	}
	if api_key == "" {
		logic_error := &handshakejserrors.LogicError{"required", "api_key", "api_key cannot be blank"}
		return api_key, logic_error
	}

	return api_key, nil
}
