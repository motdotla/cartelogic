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
  cartelogic.Setup("your-orchestrate-api-key")

  account := map[string]interface{}{"email": "email@myapp.com"}
  result, logic_error := cartelogic.AccountsCreate(account)
  if logic_error != nil {
    fmt.Println(logic_error)
  }
  fmt.Println(result)
}
```

### Setup

Connect to [Orchestrate.io](http://orchestrate.io/).

```go
cartelogic.Setup("your-orchestrate-api-key")
```

### AccountsCreate

```go
account := map[string]interface{}{"email": "email@myapp.com"}
result, logic_error := cartelogic.AccountsCreate(account)
```

### CardsCreate

```go
card := map[string]interface{}{"front": "<img src='https://some-url.com/some-image.png'>", "back": "John Doe", "api_key": "your_api_key_you_got_when_creating_an_account"}
result, logic_error := cartelogic.CardsCreate(card)
```

## Installation

```
go get github.com/scottmotte/cartelogic
```

## Running Tests

```
cp .env.example .env
```

Edit the contents of `.env.`

```
go test -v
```
