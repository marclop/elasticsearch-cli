package poller

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/elastic/elasticsearch-cli/client"
	"github.com/elastic/elasticsearch-cli/elasticsearch"
	"github.com/elastic/elasticsearch-cli/utils"
)

// IndexPoller polls the ElasticSearch API to discover which indices exist
type IndexPoller struct {
	client   client.ClientInterface
	endpoint string
	channel  chan []string
	pollRate time.Duration
}

type PollerInterface interface {
	Run()
}

// NewIndexPoller is the factory to create a new IndexPoller
func NewIndexPoller(client client.ClientInterface, c chan []string, poll int) PollerInterface {
	return &IndexPoller{
		channel:  c,
		client:   client,
		endpoint: "/_cat/indices",
		pollRate: time.Duration(poll) * time.Second,
	}
}

// Run the IndexPoller indefinetly, which will get the cluster indexList
// And will send the results back to the channel
func (w *IndexPoller) Run() {
	for {
		w.run()
		time.Sleep(w.pollRate)
	}
}

func (w *IndexPoller) run() {
	res, err := w.client.HandleCall("GET", w.endpoint, "")
	if err != nil {
		log.Print("[ERROR]: ", err)
		time.Sleep(w.pollRate / 2)
		return
	}
	defer res.Body.Close()

	w.channel <- w.parseIndices(res)
}

func (w *IndexPoller) parseIndices(res *http.Response) []string {
	var indexList []string

	if strings.TrimSpace(strings.Split(res.Header.Get("Content-Type"), ";")[0]) == "application/json" {
		var indices elasticsearch.Indices
		err := json.NewDecoder(res.Body).Decode(&indices)
		if err != nil {
			panic(err)
		}
		for _, index := range indices {
			indexList = append(indexList, index.Index)
		}

		return indexList
	}

	indicesRaw := strings.TrimSpace(utils.ReaderToString(res.Body))
	indexLines := strings.Split(indicesRaw, "\n")

	for _, indexLine := range indexLines {
		if len(strings.Fields(indexLine)) == 0 {
			continue
		}
		indexName := strings.Fields(indexLine)[2]
		indexList = append(indexList, indexName)
	}

	return indexList
}
