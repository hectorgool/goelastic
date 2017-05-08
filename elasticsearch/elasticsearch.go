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
	"fmt"
)

var client *elastic.Client

func init() {

	var err error
  // Create a client
	client, err = elastic.NewClient(elastic.SetURL(os.Getenv("ELASTICSEARCH_ENTRYPOINT")))
	if err != nil {
		panic(err)
	}

}

func termToJson(term string) (j.O, error) {

	if len(term) == 0 {
    return nil, errors.New("No string supplied")
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

	return searchJson, nil

}

func SearchTerm(term string) (string, error) {

  if len(term) == 0 {
    return "", errors.New("No string supplied")
  }

	//Convert string to json query for elasticsearch
	searchJson, err := termToJson(term)
	if err != nil {
		panic(err)
	}

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
	jsonDocuments, err := json.Marshal(documents)
	if err != nil {
		log.Fatal(err)
	}

  //return jsonDocuments in json format
	return string(jsonDocuments), nil

}

func CreateDocument(id string, document Document) (string, error) {

	// Add a document to the index
	var err error
	_, err = client.Index().
	    Index(os.Getenv("ELASTICSEARCH_INDEX")).
	    Type(os.Getenv("ELASTICSEARCH_TYPE")).
	    Id(id).
	    BodyJson(document).
	    Refresh(true).
	    Do()
	if err != nil {
	    // Handle error
	    panic(err)
	}
	msg := fmt.Sprintf("The document: %s, has been save", id)
	return msg, nil

}

func UpdateDocument(id string, document Document) (string, error) {

	// Update a tweet by the update API of Elasticsearch.
	// We just increment the number of retweets.
	update, err := client.Update().
		Index(os.Getenv("ELASTICSEARCH_INDEX")).
		Type(os.Getenv("ELASTICSEARCH_TYPE")).
		Id(id).
		//Script("ctx._source.retweets += num").
	  //ScriptParams(map[string]interface{}{"num": 1}).
	  //Upsert(map[string]interface{}{"retweets": 0}).
	  Do()
	if err != nil {
	    // Handle error
	    panic(err)
	}
	msg := fmt.Sprintf("New version of tweet %q is now %d", update.Id, update.Version)
	return msg, nil

}

func DeleteDocument(id string) (string, error) {

	// Delete tweet with specified ID
	res, err := client.Delete().
	    Index("twitter").
	    Type("tweet").
	    Id("1").
	    Do()
	if err != nil {
	    // Handle error
	    panic(err)
	}
	//if res.Found {
	return fmt.Sprintf("Document: %s, deleted from from index\n",res.Found ), nil
	//}

}

func Ping() (string, error) {

	info, code, err := client.Ping(os.Getenv("ELASTICSEARCH_ENTRYPOINT")).Do()
	if err != nil {
	    // Handle error
	    panic(err)
	}
	msg := fmt.Sprintf("Elasticsearch returned with code %d and version %s", code, info.Version.Number)
	return msg, nil

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
