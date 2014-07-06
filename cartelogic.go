package cartelogic

import (
	//"bytes"
	//"crypto/sha1"
	//"encoding/hex"
	"errors"
	"github.com/dchest/uniuri"
	"github.com/garyburd/redigo/redis"
	"github.com/handshakejs/handshakejserrors"
	"github.com/scottmotte/redisurlparser"
	"log"
	//"math/rand"
	"strconv"
	"time"
)

var (
	DB_ENCRYPTION_SALT        string
	AUTHCODE_LIFE_IN_MS       int64
	AUTHCODE_LENGTH           int
	KEY_EXPIRATION_IN_SECONDS int
	PBKDF2_HASH_ITERATIONS    int
	PBKDF2_HASH_BITES         int
	redisurl                  redisurlparser.RedisURL
	pool                      *redis.Pool
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

func main() {
}

func DecksCreate(deck map[string]interface{}) (map[string]interface{}, *handshakejserrors.LogicError) {
	var name string
	if str, ok := deck["name"].(string); ok {
		name = str
	} else {
		name = ""
	}
	if name == "" {
		logic_error := &handshakejserrors.LogicError{"required", "name", "name cannot be blank"}
		return deck, logic_error
	}
	deck["name"] = name

	var email string
	if str, ok := deck["email"].(string); ok {
		email = str
	} else {
		email = ""
	}
	if email == "" {
		logic_error := &handshakejserrors.LogicError{"required", "email", "email cannot be blank"}
		return deck, logic_error
	}
	deck["email"] = email

	generated_api_key := uniuri.NewLen(20)
	if deck["api_key"] == nil {
		deck["api_key"] = generated_api_key
	}
	if deck["api_key"].(string) == "" {
		deck["api_key"] = generated_api_key
	}

	current_ms_epoch_time := (time.Now().Unix() * 1000)
	deck["id"] = strconv.FormatInt(current_ms_epoch_time, 10)
	key := "decks/" + deck["id"].(string)

	err := validateDeckDoesNotExist(key)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"not_unique", "name", "name must be unique"}
		return deck, logic_error
	}
	err = addDeckToDecks(deck["name"].(string))
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}
	err = saveDeck(key, deck)
	if err != nil {
		logic_error := &handshakejserrors.LogicError{"unknown", "", err.Error()}
		return nil, logic_error
	}

	return deck, nil
}

func validateDeckDoesNotExist(key string) error {
	conn := Conn()
	defer conn.Close()
	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		log.Printf("ERROR " + err.Error())
		return err
	}
	if exists == true {
		err = errors.New("That id already exists.")
		return err
	}

	return nil
}

func addDeckToDecks(name string) error {
	conn := Conn()
	defer conn.Close()
	_, err := conn.Do("SADD", "decks", name)
	if err != nil {
		return err
	}

	return nil
}

func saveDeck(key string, deck map[string]interface{}) error {
	args := []interface{}{key}
	for k, v := range deck {
		args = append(args, k, v)
	}

	conn := Conn()
	defer conn.Close()
	_, err := conn.Do("HMSET", args...)
	if err != nil {
		return err
	}

	return nil
}

func Conn() redis.Conn {
	return pool.Get()
}
