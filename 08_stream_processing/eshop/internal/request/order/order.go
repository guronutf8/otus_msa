package order

type RequestOrder struct {
	User  string `json:"user"`
	Items []Item `json:"items"`
}

type Item struct {
	Title string `json:"title"`
	Count int32  `json:"count"`
}

type Response struct {
	Status  string   `json:"status"`
	Message string   `json:"message"`
	Log     []string `json:"log"`
}
