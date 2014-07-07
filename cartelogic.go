package cartelogic

import (
	"errors"
	"github.com/dchest/uniuri"
	"github.com/garyburd/redigo/redis"
	"github.com/handshakejs/handshakejserrors"
	"github.com/scottmotte/redisurlparser"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	BASE_10 = 10
)

var (
	redisurl redisurlparser.RedisURL
	pool     *redis.Pool
)

func Setup(redis_url_string string) {
	redisurl, err := redisurlparser.Parse(redis_url_string)
	if err != nil {
		log.Fatal(err)
	}

	pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisurl.Host+":"+redisurl.Port)
			if err != nil {
				log.Fatal(err)
				return nil, err
			}

			if redisurl.Password != "" {
				if _, err := c.Do("AUTH", redisurl.Password); err != nil {
					c.Close()
					log.Fatal(err)
					return nil, err
				}
			}
			return c, err
		},
	}
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
		return account, logic_error
	}
	account["email"] = email

	generated_api_key := uniuri.NewLen(20)
	if account["api_key"] == nil {
		account["api_key"] = generated_api_key
	}
	if account["api_key"].(string) == "" {
		account["api_key"] = generated_api_key
	}

	key := "accounts/" + account["email"].(string)

	err := validateAccountDoesNotExist(key)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"not_unique", "email", "email must be unique"}
		return account, logic_error
	}
	err = addAccountToAccounts(account["email"].(string))
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}
	err = saveAccount(key, account)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}

	return account, nil
}

func CardsCreate(card map[string]interface{}) (map[string]interface{}, *handshakejserrors.LogicError) {
	front, logic_error := checkFrontPresent(card)
	if logic_error != nil {
		return card, logic_error
	}
	card["front"] = front

	back, logic_error := checkBackPresent(card)
	if logic_error != nil {
		return card, logic_error
	}
	card["back"] = back

	api_key, logic_error := checkApiKeyPresent(card)
	if logic_error != nil {
		return card, logic_error
	}
	card["api_key"] = api_key

	email, err := getEmailAssociatedWithApiKey(card["api_key"].(string))
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"incorrect", "api_key", "the api_key is incorrect"}
		return card, logic_error
	}
	card["email"] = email

	// set card id
	current_ms_epoch_time_as_int64 := (time.Now().Unix() * 1000)
	current_ms_epoch_time := strconv.FormatInt(current_ms_epoch_time_as_int64, 10)
	card["id"] = current_ms_epoch_time

	account_key := "accounts/" + card["email"].(string)
	key := account_key + "/cards/" + card["id"].(string)

	err = validateAccountExists(account_key)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"not_found", "account", "account could not be found"}
		return card, logic_error
	}
	err = addCardToCards(account_key, card["email"].(string))
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return card, logic_error
	}
	err = saveCard(key, card)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}

	return card, nil
}

func KeyExists(key string) (bool, error) {
	conn := Conn()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return false, err
	}
	return exists, nil
}

func validateAccountExists(key string) error {
	conn := Conn()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}
	if !exists {
		err = errors.New("That account does not exist.")
		return err
	}

	return nil
}

func validateAccountDoesNotExist(key string) error {
	exists, err := KeyExists(key)
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}
	if exists == true {
		err = errors.New("That email already exists.")
		return err
	}

	return nil
}

func addAccountToAccounts(email string) error {
	conn := Conn()
	defer conn.Close()
	_, err := conn.Do("SADD", "accounts", email)
	if err != nil {
		return err
	}

	return nil
}

func saveAccount(key string, account map[string]interface{}) error {
	args := []interface{}{key}
	for k, v := range account {
		args = append(args, k, v)
	}

	conn := Conn()
	defer conn.Close()
	pointer_key := "accounts/" + account["api_key"].(string)
	_, err := conn.Do("SET", pointer_key, key)
	if err != nil {
		return err
	}

	_, err = conn.Do("HMSET", args...)
	if err != nil {
		return err
	}

	return nil
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

func addCardToCards(account_key string, email string) error {
	conn := Conn()
	defer conn.Close()
	_, err := conn.Do("SADD", account_key+"/cards", email)
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}

	return nil
}

func saveCard(key string, card map[string]interface{}) error {
	unixtime := (time.Now().Unix() * 1000)
	card["id"] = strconv.FormatInt(unixtime, BASE_10)

	args := []interface{}{key}
	args = append(args, "front", card["front"].(string))
	args = append(args, "back", card["back"].(string))
	args = append(args, "id", card["id"].(string))

	conn := Conn()
	defer conn.Close()
	_, err := conn.Do("HMSET", args...)
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}

	return nil
}

func getEmailAssociatedWithApiKey(api_key string) (string, error) {
	pointer_key := "accounts/" + api_key

	conn := Conn()
	defer conn.Close()
	account_key, err := redis.String(conn.Do("GET", pointer_key))
	log.Printf(account_key)
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return "", err
	}

	split_result := strings.Split(account_key, "/")
	result := split_result[len(split_result)-1]

	return result, nil
}

func Conn() redis.Conn {
	return pool.Get()
}
