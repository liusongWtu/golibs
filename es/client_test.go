package es

import (
	"encoding/json"
	"fmt"
	"github.com/olivere/elastic"
	"reflect"
	"testing"
	"time"
)

type Tweet struct {
	User     string                `json:"user"`
	Age      int                   `json:"age"`
	Message  string                `json:"message"`
	Retweets int                   `json:"retweets"`
	Image    string                `json:"image,omitempty"`
	Created  time.Time             `json:"created,omitempty"`
	Tags     []string              `json:"tags,omitempty"`
	Location string                `json:"location,omitempty"`
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

var addr = []string{"http://127.0.0.1:9200"}

var mapping = `{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 1
	},
	"mappings":{
		"doc":{
			"properties":{
				"user":{
					"type":"keyword"
				},
				"age":{
					"type": "integer"
				},
				"message":{
					"type":"text",
					"store": true,
					"fielddata": true
				},
				"image":{
					"type":"keyword"
				},
				"created":{
					"type":"date"
				},
				"tags":{
					"type":"keyword"
				},
				"location":{
					"type":"geo_point"
				},
				"suggest_field":{
					"type":"completion"
				}
			}
		}
	}
}`

func TestPingNode(t *testing.T) {
	client := NewElasticClient(addr)
	client.PingNode("http://127.0.0.1:9200")
}

func TestIndexExists(t *testing.T) {
	client := NewElasticClient(addr)
	result,err := client.IndexExists("car_source", "test")
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("all index exists: ", result)
}

func TestDeleteIndex(t *testing.T) {
	client := NewElasticClient(addr)
	result,err := client.DelIndex("twitter")
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("all index deleted: ", result)
}

func TestCreateIndex(t *testing.T) {
	client := NewElasticClient(addr)
	result,err := client.CreateIndex("twitter", mapping)
	if err!=nil {
		fmt.Println(err)
	}
	fmt.Println("mapping created: ", result)
}

func TestBatch(t *testing.T) {
	tweet1 := Tweet{User: "Jame1",Age: 23, Message: "Take One", Retweets: 1, Created: time.Now()}
	tweet2 := Tweet{User: "Jame2",Age: 32, Message: "Take Two", Retweets: 0, Created: time.Now()}
	tweet3 := Tweet{User: "Jame3",Age: 32, Message: "Take Three", Retweets: 0, Created: time.Now()}
	data :=make(map[string]interface{})
	data["1"] = tweet1
	data["2"] = tweet2
	data["3"] = tweet3
	client := NewElasticClient(addr)
	client.Batch("twitter", data)
}

func TestGetDoc(t *testing.T) {
	var tweet Tweet
	client := NewElasticClient(addr)
	data,err := client.GetDoc("twitter", "1")
	if err!=nil {
		fmt.Println(err)
	}
	if err := json.Unmarshal(data, &tweet); err == nil {
		fmt.Printf("data: %v\n", tweet)
	}
}

func TestTermQuery(t *testing.T) {
	var tweet Tweet
	client := NewElasticClient(addr)
	result := client.TermQuery("twitter", "user", "Take Two",0,0)
	//获得数据, 方法一
	for _, item := range result.Each(reflect.TypeOf(tweet)) {
		if t, ok := item.(Tweet); ok {
			fmt.Printf("tweet : %v\n", t)
		}
	}
	//获得数据, 方法二
	fmt.Println("num of raws: ", result.Hits.TotalHits)
	if result.Hits.TotalHits.Value > 0 {
		for _, hit := range result.Hits.Hits {
			err := json.Unmarshal([]byte(fmt.Sprintf("%v", hit.Source)), &tweet)
			if err != nil {
				fmt.Printf("source convert json failed, err: %v\n", err)
			}
			fmt.Printf("data: %v\n", tweet)
		}
	}
}

func TestSearch(t *testing.T) {
	boolQuery := elastic.NewBoolQuery()
	boolQuery.Must(elastic.NewMatchQuery("user", "Jame10"))
	boolQuery.Filter(elastic.NewRangeQuery("age").Gt("30"))

	client := NewElasticClient(addr)
	result,err := client.Search("twitter",boolQuery,0,0)
	if err!=nil {
		fmt.Println(err)
	}
	var tweet Tweet
	for _, item := range result.Each(reflect.TypeOf(tweet)) {
		if t, ok := item.(Tweet); ok {
			fmt.Printf("tweet : %v\n", t)
		}
	}
}

func TestAggsSearch(t *testing.T) {
	client := NewElasticClient(addr)

	minAgg := elastic.NewMinAggregation().Field("age")
	minResult, err :=client.AggsSearch("twitter", "minAgg",minAgg,0,0)
	if err!=nil {
		fmt.Println(err)
	}
	minAggRes, _ := minResult.Aggregations.Min("minAgg")
	fmt.Printf("min: %v\n", *minAggRes.Value)


	rangeAgg := elastic.NewRangeAggregation().Field("age").AddRange(0,30).AddRange(30,60).Gt(60)
	rangeResult, err := client.AggsSearch("twitter", "rangeAgg",rangeAgg,0,0)
	if err!=nil {
		fmt.Println(err)
	}
	rangeAggRes, _ := rangeResult.Aggregations.Range("rangeAgg")
	for _, item := range rangeAggRes.Buckets {
		fmt.Printf("key: %s, value: %v\n", item.Key, item.DocCount)
	}
}

