package store

type Reserve struct {
	Items []Item `json:"items"`
}

type Item struct {
	Title string `bson:"title" json:"title"`
	Count int32  `bson:"count" json:"count"`
}
