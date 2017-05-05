/*
go install
create an archive file suffixed with .a, and that’s stored in
the pkg directory.

If you want to export variables, constants, and functions to be used with other
programs, the name of the identifier must start with an uppercase letter.
*/

package elasticsearch

import (
	"encoding/json"
  "errors"
  "gopkg.in/olivere/elastic.v3"
	j "github.com/ricardolonga/jsongo"
	"log"
	"os"
  "reflect"
)

var client *elastic.Client

func init() {

  os.Setenv("ELASTICSEARCH_HOSTS", "http://172.17.0.4:9200")
  os.Setenv("ELASTICSEARCH_INDEX", "mx")
  os.Setenv("ELASTICSEARCH_TYPE", "postal_code")

	var err error
  // Create a client
	client, err = elastic.NewClient(elastic.SetURL(os.Getenv("ELASTICSEARCH_HOSTS")))
	if err != nil {
		panic(err)
	}

}

func SearchTerm(term string) (string, error) {

  if len(term) == 0 {
    return "", errors.New("No string supplied")
  }

  //create json object with the string term
  //searchJson has jsongo.O type
  searchJson := j.Object().
  		Put("size", 10).
  		Put("query", j.Object().
  			Put("match", j.Object().
  				Put("_all", j.Object().
  					Put("query", term).
  					Put("operator", "and")))).
  		Put("sort", j.Array().
  			Put(j.Object().
  				Put("colonia", j.Object().
  					Put("order", "asc").
  					Put("mode", "avg"))))
  log.Println(searchJson.Indent())
  log.Println(reflect.TypeOf(searchJson))

  // Search with a term source
	searchResult, err := client.Search().
		Index(os.Getenv("ELASTICSEARCH_INDEX")).
		Type(os.Getenv("ELASTICSEARCH_TYPE")).
		Source(searchJson).
		Do()
	if err != nil {
		panic(err)
	}

	var documents []Document
	for _, hit := range searchResult.Hits.Hits {
		var d Document
    //parses *hit.Source into the instance of the Document struct
		err := json.Unmarshal(*hit.Source, &d)
		if err != nil {
			log.Fatal(err)
		}
    //Puts d into a map for later access
		documents = append(documents, d)
	}

  //Convert documents data to json
	rawJsonDocuments, err := json.Marshal(documents)
	if err != nil {
		log.Fatal(err)
	}

  //return rawJsonDocuments in json format
	return string(rawJsonDocuments), nil

}

//Document type with annotates struct fields for JSON encoding and decoding
type Document struct {
	Ciudad     string `json:"ciudad"`
	Colonia    string `json:"colonia"`
	Cp         string `json:"cp"`
	Delegacion string `json:"delegacion"`
	Location   struct {
		Lat string `json:"lat"`
		Lon string `json:"lon"`
	} `json:"location"`
}
