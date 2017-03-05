package elasticsearch

// Indices Repnresents the JSON response for /_cat/indices when Content-Type is set to application/JSON (<5.x)
type Indices []indexCat

type indexCat struct {
	Health    string `json:"health"`
	Status    string `json:"status"`
	Index     string `json:"index"`
	Primaries string `json:"pri"`
	Replicas  string `json:"rep"`
	Docs      indexDocs
	Store     indexStore
}

type indexDocs struct {
	Count   string `json:"count"`
	Deleted string `json:"deleted"`
}

type indexStore struct {
	Primary string `json:"pri.store.size"`
	Size    string `json:"size"`
}
