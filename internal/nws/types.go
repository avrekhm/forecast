package nws

// /points response types
type PointsResponse struct {
	ID         string           `json:"id"`
	Properties PointsProperties `json:"properties"`
}

type PointsProperties struct {
	ID       string `json:"@id"`
	GridID   string `json:"gridId"`
	GridX    int    `json:"gridX"`
	GridY    int    `json:"gridY"`
	Forecast string `json:"forecast"`
}

// /gridpoints/{wfo} response types

type GridpointsResponse struct {
	Properties Properties `json:"properties,omitempty"`
}

type Properties struct {
	Units   string   `json:"units"`
	Periods []Period `json:"periods"`
}

type Period struct {
	Number           int    `json:"number"`
	Name             string `json:"name"`
	StartTime        string `json:"startTime"`
	EndTime          string `json:"endTime"`
	Temperature      int    `json:"temperature"`
	TemperatureUnit  string `json:"temperatureUnit"`
	ShortForecast    string `json:"shortForecast"`
	DetailedForecast string `json:"detailedForecast"`
}
