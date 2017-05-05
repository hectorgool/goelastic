/*
go build
go run -v main.go
*/
package main

import (
  "github.com/hectorgool/gomicrosearch2/elasticsearch"
  "fmt"
)

func main() {

  if result, err := elasticsearch.SearchTerm("cortes de villa"); err != nil {
    fmt.Printf("Error: %s\n", err)
  } else {
    fmt.Printf("ElasticSearch result: '%s'\n", result)
  }

}
