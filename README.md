# cartelogic

[![BuildStatus](https://travis-ci.org/scottmotte/cartelogic.png?branch=master)](https://travis-ci.org/scottmotte/cartelogic)

Logic for saving carte data to the redis database.

This library is part of the larger [Carte ecosystem](https://github.com/scottmotte/carte).

## Usage

```go
package main

import (
  "fmt"
  "github.com/scottmotte/cartelogic"
)

func main() {
  cartelogic.Setup("redis://127.0.0.1:6379")

  account := map[string]interface{}{"email": "email@myapp.com"}
  result, logic_error := cartelogic.AccountsCreate(account)
  if logic_error != nil {
    fmt.Println(logic_error)
  }
  fmt.Println(result)
}
```

### Setup

Connects to Redis.

```go
cartelogic.Setup("redis://127.0.0.1.6379")
```

### AccountsCreate

```go
account := map[string]interface{}{"email": "email@myapp.com"}
result, logic_error := cartelogic.AccountsCreate(account)
```

### CardsCreate

```go
card := map[string]interface{}{"front": "<img src='https://some-url.com/some-image.png'>", "back": "John Doe", "api_key": "your_api_key_you_got_when_registering_an_account"}
result, logic_error := cartelogic.CardsCreate(card)
```

## Installation

```
go get github.com/scottmotte/cartelogic
```

## Running Tests

```
go test -v
```

## Database Schema Details (using Redis)

Cartelogic uses a purposely simple database schema - as simple as possible. If you know a simpler approach, even better, please let me know or share as a pull request. 

Cartelogic uses Redis because of its light footprint, ephemeral nature, and lack of migrations.

+ /accounts - collection of keys with all the account ids in there. SADD
+ /accounts/:email - hash with all the data in there. HSET or HMSET
+ /accounts/:api_key - acts as a pointer. key/value pair where the value points to /accounts/:email
+ /accounts/:email/cards - collection of keys with all the cards' ids in there. SADD
+ /accounts/:email/cards/timestamp HSET or HMSET

