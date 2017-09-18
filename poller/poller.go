package poller

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/marclop/elasticsearch-cli/client"
	"github.com/marclop/elasticsearch-cli/elasticsearch"
)

const defaultPollingEndpoint = "/_cat/indices"

// IndexPoller polls the ElasticSearch API to discover which indices exist
type IndexPoller struct {
	client         client.Client
	endpoint       string
	channel        chan []string
	pollRate       time.Duration
	controlChannel chan bool
}

// NewIndexPoller is the factory to create a new IndexPoller
func NewIndexPoller(client client.Client, c chan []string, poll int) *IndexPoller {
	return &IndexPoller{
		channel:        c,
		client:         client,
		endpoint:       defaultPollingEndpoint,
		pollRate:       time.Duration(poll) * time.Second,
		controlChannel: make(chan bool, 1),
	}
}

// Run the IndexPoller indefinitely, which will get the cluster indexList
// And will send the results back to the channel
func (w *IndexPoller) Run() {
	w.channel <- w.run()
	ticker := time.NewTicker(w.pollRate)
	for {
		select {
		case <-w.controlChannel:
			defer close(w.channel)
			ticker.Stop()
			return
		case <-ticker.C:
			w.channel <- w.run()
		}
	}
}

// Stop closes the channels
func (w *IndexPoller) Stop() {
	defer close(w.controlChannel)
	w.controlChannel <- true
}

func (w *IndexPoller) run() []string {
	res, err := w.client.HandleCall("GET", w.endpoint, "")
	if err != nil {
		log.Print("[ERROR]: ", err)
		return nil
	}
	defer res.Body.Close()

	return w.parseIndices(res)
}

func (w *IndexPoller) parseIndices(res *http.Response) []string {
	var indexList []string

	if strings.Contains(res.Header.Get("Content-Type"), "application/json") {
		var indices elasticsearch.Indices

		err := json.NewDecoder(res.Body).Decode(&indices)
		if err != nil {
			log.Print("[ERROR]: ", err)
		}

		for _, index := range indices {
			indexList = append(indexList, index.Index)
		}

		return indexList
	}

	bytesBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Print("[ERROR]: ", err)
	}
	indicesRaw := strings.TrimSpace(string(bytesBody))
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
