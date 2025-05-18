package models

type SteamspyResponse struct {
	AppID        int    `json:"appid"`
	Name         string `json:"name"`
	Price        string `json:"price"`
	InitialPrice string `json:"initialprice"`
	Discount     string `json:"discount"`
	Positive     int    `json:"positive"`
	Negative     int    `json:"negative"`
}

type SteamspyDetails struct {
	Genres string         `json:"genre"`
	Tags   map[string]int `json:"tags"`
}
