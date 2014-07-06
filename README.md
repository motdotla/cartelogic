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

  deck := map[string]interface{}{"email": "email@myapp.com", "name": "People Deck"}
  result, logic_error := cartelogic.DecksCreate(deck)
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

### DecksCreate

```go
deck := map[string]interface{}{"email": "email@myapp.com", "name": "People Deck"}
result, logic_error := cartelogic.AppsCreate(deck)
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

* decks - collection of keys with all the deck ids in there. SADD
* decks/ID - hash with all the data in there. HSET or HMSET
* decks/ID/cards - collection of keys with all the cards' emails in there. SADD
* decks/ID/cards/emailaddress HSET or HMSET

