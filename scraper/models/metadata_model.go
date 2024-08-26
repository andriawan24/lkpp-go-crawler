package models

type Metadata struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	NPWP          string `json:"npwp"`
	Address       string `json:"address"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Status        string `json:"status"`
	StartDate     string `json:"start_date"`
	EndDate       string `json:"end_date"`
	PublishedDate string `json:"published_date"`
	Verdict       string `json:"verdict"`
	Number        string `json:"number"`
	Rule          string `json:"rule"`
	Description   string `json:"description"`
}
