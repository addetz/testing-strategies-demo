package data

type Talk struct {
	Title    string   `json:"title"`
	Speakers []string `json:"speakers"`
	Date     string   `json:"date"`
	Time     string   `json:"time"`
	EventID  string   `json:"event_id"`
}

type Event struct {
	ID        string `json:"ID"`
	Name      string `json:"name"`
	DateStart string `json:"date_start"`
	DateEnd   string `json:"date_end"`
	Location  string `json:"location"`
	Talks     []Talk `json:"-"`
}

type Events struct {
	Events []Event `json:"events"`
}

type Talks struct {
	Talks []Talk `json:"talks"`
}
