/*
go fmt .
go build
go run -v main.go
./goelastic
*/
package main

import (
	"fmt"
	"github.com/hectorgool/goelastic/elasticsearch"
)

func main() {

	if result, err := elasticsearch.SearchTerm("cortes de villa"); err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("ElasticSearch result: '%s'\n", result)
	}

}
