package bluesky

type BlueskyPost struct {
	Type      string `json:"$type"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
}

type Facet struct {
	Index    FacetIndex     `json:"index"`
	Features []FacetFeature `json:"features"`
}

type FacetIndex struct {
	ByteStart int `json:"byteStart"`
	ByteEnd   int `json:"byteEnd"`
}

type FacetFeature struct {
	Type string `json:"$type"`
	Did  string `json:"did,omitempty"`
	URI  string `json:"uri,omitempty"`
}

type BlueskyAuth struct {
	Identifier string `json:"identifier"`
	Password   string `json:"password"`
}

type BlueskyAuthResponse struct {
	AccessJwt  string `json:"accessJwt"`
	RefreshJwt string `json:"refreshJwt"`
	Handle     string `json:"handle"`
	Did        string `json:"did"`
}
