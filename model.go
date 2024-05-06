package main

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

type JSONData struct {
	Sequences []Sequence `json:"sequences"`
}
