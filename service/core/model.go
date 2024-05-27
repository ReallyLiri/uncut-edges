package core

type Canvas struct {
	Images []struct {
		Resource struct {
			ID string `json:"@id"`
		} `json:"resource"`
	} `json:"images"`
}

type Sequence struct {
	Canvases []Canvas `json:"canvases"`
}

type Manifest struct {
	Sequences []Sequence `json:"sequences"`
}

type Property struct {
	Key   string
	Value string
}

type Header struct {
	CatalogID   string
	Title       string
	Properties  []Property
	Links       []string
	ManifestURL string
}
