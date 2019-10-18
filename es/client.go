package es

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic"
	"libs/utils"
	"time"
)

func NewElasticClient(addr []string) *ElasticClient {
	eclient := &ElasticClient{}
	eclient.Init(addr)
	return eclient
}

type ElasticClient struct {
	client *elastic.Client
}

//初始化
func (ec *ElasticClient) Init(addr []string) {
	var err error
	if utils.IsMac() {
		//ec.client, err = elastic.NewClient(elastic.SetURL(addr...),
		//	elastic.SetBasicAuth("elastic", "changeme"), elastic.SetSniff(false))
		ec.client, err = elastic.NewClient(elastic.SetURL(addr...), elastic.SetSniff(false))
	} else {
		ec.client, err = elastic.NewClient(elastic.SetURL(addr...))
	}
	if err != nil {
		panic("create elastic client failed, err:" + err.Error())
	}
}

//ping连接测试
func (ec *ElasticClient) PingNode(addr string) (time.Duration, *elastic.PingResult, error) {
	start := time.Now()
	info, _, err := ec.client.Ping(addr).Do(context.Background())
	if err != nil {
		return 0, nil, err
	}
	end := time.Since(start)
	return end, info, nil
}

//校验index是否存在
func (ec *ElasticClient) IndexExists(index ...string) (bool, error) {
	exists, err := ec.client.IndexExists(index...).Do(context.Background())
	if err != nil {
		return false, err
	}
	return exists, nil
}

//创建index
func (ec *ElasticClient) CreateIndex(index, mapping string) (bool, error) {
	service := ec.client.CreateIndex(index)
	if mapping != "" {
		service.BodyString(mapping)
	}
	result, err := service.Do(context.Background())
	if err != nil {
		return false, err
	}
	return result.Acknowledged, nil
}

//删除index
func (ec *ElasticClient) DelIndex(index ...string) (bool, error) {
	response, err := ec.client.DeleteIndex(index...).Do(context.Background())
	if err != nil {
		return false, err
	}
	return response.Acknowledged, nil
}

//批量插入
func (ec *ElasticClient) BatchWithID(index string, data map[interface{}]interface{}) (*elastic.BulkResponse, error) {
	bulkRequest := ec.client.Bulk()
	for id, v := range data {
		doc := elastic.NewBulkIndexRequest().Index(index).Id(ec.getString(id)).Doc(v)
		bulkRequest = bulkRequest.Add(doc)
	}
	//插入elastic search
	return bulkRequest.Do(context.TODO())
}

//批量插入
func (ec *ElasticClient) Batch(index string, data []interface{}) (*elastic.BulkResponse, error) {
	bulkRequest := ec.client.Bulk()
	for _, v := range data {
		doc := elastic.NewBulkIndexRequest().Index(index).Doc(v)
		bulkRequest = bulkRequest.Add(doc)
	}
	//插入elastic search
	return bulkRequest.Do(context.TODO())
}

//获取指定Id 的文档
func (ec *ElasticClient) GetDoc(index, id string) ([]byte, error) {
	result, err := ec.client.Get().Index(index).Id(id).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if !result.Found {
		return nil, errors.New("not find the document")
	}
	source, err := result.Source.MarshalJSON()
	if err != nil {
		return nil, err
	}
	return source, nil
}

//term查询
func (ec *ElasticClient) TermQuery(index, fieldName, fieldValue string, size, page int) *elastic.SearchResult {
	query := elastic.NewTermQuery(fieldName, fieldValue)
	service := ec.client.Search().Index(index).Query(query)
	if size > 0 {
		service.Size(size)
	}
	if page > 0 {
		service.From((page - 1) * size)
	}
	searchResult, err := service.Pretty(true).Do(context.Background())

	if err != nil {
		panic(err)
	}
	fmt.Printf("query cost %d millisecond.\n", searchResult.TookInMillis)

	return searchResult
}

//query搜索
func (ec *ElasticClient) Search(index string, query elastic.Query, size, page int) (*elastic.SearchResult, error) {
	service := ec.client.Search(index).Query(query).Pretty(true)
	if size > 0 {
		service.Size(size)
	}
	if page > 0 {
		service.From((page - 1) * size)
	}
	result, err := service.Do(context.Background())
	if err != nil {
		return result, err
	}

	return result, nil
}

//aggregation搜索
func (ec *ElasticClient) AggsSearch(index, aggName string, agg elastic.Aggregation, size, page int) (*elastic.SearchResult, error) {
	service := ec.client.Search(index).Pretty(true)
	if size > 0 {
		service.Size(size)
	}
	if page > 0 {
		service.From((page - 1) * size)
	}
	return service.Aggregation(aggName, agg).Do(context.Background())
}

func (ec *ElasticClient) QueryAggregationSearch(index string, query elastic.Query, aggName string, agg elastic.Aggregation) (*elastic.SearchResult, error) {
	return ec.client.Search(index).Query(query).Aggregation(aggName, agg).Pretty(true).Do(context.Background())
}

func (ec *ElasticClient) QueryAggregationsSearch(index string, query elastic.Query, aggMap map[string]elastic.Aggregation) (*elastic.SearchResult, error) {
	searchService := ec.client.Search(index).Query(query)
	for aggName, agg := range aggMap {
		searchService.Aggregation(aggName, agg)
	}
	return searchService.Pretty(true).Do(context.Background())
}

func (ec *ElasticClient) getString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprintf("%v", result)
		}
	}
	return ""
}

func (ec *ElasticClient) QueryCount(index string, query elastic.Query) (count int64, err error) {
	count, err = ec.client.Count(index).Query(query).Do(context.Background())
	return
}
