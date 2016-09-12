package custom

import (
	"github.com/andrepinto/goway/proxy"
	"gopkg.in/olivere/elastic.v2"
	"time"
)

type ElasticLog struct {
	uri    string
	client *elastic.Client
	index 	string
	table 	string
}

func NewElasticLog(uri string, index string, table string) *ElasticLog{
	return &ElasticLog{
		uri: uri,
		index: index,
		table: table,
	}
}

func (es *ElasticLog) Start() error {
	client, err := elastic.NewClient(
		elastic.SetURL(es.uri),
		elastic.SetSniff(false),
		elastic.SetHealthcheckInterval(10*time.Second),
		elastic.SetMaxRetries(5))

	if err != nil {
		// Handle error
		return err
	}


	es.client = client
	return nil
}

func(es *ElasticLog) Log(record *proxy.LogRecord){
	es.Insert(record)
}


func (es *ElasticLog) Insert(model interface{}) error {
	_, err := es.client.Index().
		Index(es.index).
		Type(es.table).
		BodyJson(model).
		Do()

	if err != nil {
		return err
	}

	// Flush data
	_, err = es.client.Flush().Index(es.index).Do()

	if err != nil {
		return err
	}

	return nil
}