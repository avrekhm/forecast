package nws

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"
)

type Client struct {
	HTTPClient http.Client
	URL        string
}

var httpClient = http.Client{
	Timeout: 5 * time.Second,
}

type NWSer interface {
	GetGridpointsResponse(ctx context.Context, pointsResponse *PointsResponse) (*GridpointsResponse, error)
	GetPointsResponse(ctx context.Context, lat, long float64) (*PointsResponse, error)
}

type Forecast struct {
	ShortForecast string `json:"short_forecast"`
	Description   string `json:"description"`
}

func NewClient(ctx context.Context) *Client {
	return &Client{
		HTTPClient: httpClient,
		URL:        "https://api.weather.gov",
	}
}

func GetForecast(ctx context.Context, s NWSer, lat, long float64) (*Forecast, error) {
	var forecast *Forecast

	pointsResponse, err := s.GetPointsResponse(ctx, lat, long)
	if err != nil {
		return nil, err
	}

	gridPointsResponse, err := s.GetGridpointsResponse(ctx, pointsResponse)
	if err != nil {
		return nil, err
	}

	forecast, err = processGridpointsResponse(gridPointsResponse)
	if err != nil {
		return nil, err
	}

	return forecast, nil
}

// GetGridpointsResponse gets a response from the NWS /gridpoints endpoint, which contains the
// weather forecast for a specific grid point under a specific weather forecast office code
func (c *Client) GetGridpointsResponse(ctx context.Context, pointsResponse *PointsResponse) (*GridpointsResponse, error) {
	var gridPointsResponse GridpointsResponse

	if pointsResponse == nil {
		return nil, errors.New("missing points response")
	}

	forecastURL := pointsResponse.Properties.Forecast

	resp, err := c.getHTTPResponse(ctx, forecastURL)
	if err != nil {
		return nil, fmt.Errorf("error %v for URL %s", err, forecastURL)
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&gridPointsResponse); err != nil {
		return nil, err
	}

	return &gridPointsResponse, nil
}

// GetPointsResponse gets a response from the NWS /points endpoint, which contains the URL with the forecast
// for the NWS grid square that contains the provided coordinates
func (c *Client) GetPointsResponse(ctx context.Context, lat, long float64) (*PointsResponse, error) {
	var pointsResponse PointsResponse

	// The NWS /points endpoint returns a redirect with an error message if more than 4 decimal digits are provided.
	// Instead of dealing with the redirect, here we truncate and round to 4 decimal digits.
	lat = truncateToFourDecimals(lat)
	long = truncateToFourDecimals(long)

	url := c.URL + fmt.Sprintf("/points/%f,%f",
		lat,
		long,
	)

	resp, err := c.getHTTPResponse(ctx, url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&pointsResponse); err != nil {
		return nil, err
	}

	return &pointsResponse, nil
}

func (c *Client) getHTTPResponse(_ context.Context, url string) (*http.Response, error) {
	// not using context here but it could potentially be used for things like auth info and observability tags

	if url == "" {
		return nil, errors.New("empty URL")
	}

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(http.StatusText(resp.StatusCode))
	}

	return resp, nil
}

func processGridpointsResponse(response *GridpointsResponse) (*Forecast, error) {
	if response == nil {
		return nil, errors.New("nil gridpoints response")
	}

	// assume that the first of the "periods" under properties is the current forecast
	currentForecast := response.Properties.Periods[0]

	if currentForecast.ShortForecast == "" {
		return nil, errors.New("missing short forecast")
	}

	// only dealing with Fahrenheit here for now
	temperature := Fahrenheit(currentForecast.Temperature)

	return &Forecast{
		ShortForecast: currentForecast.ShortForecast,
		Description:   temperature.Description(),
	}, nil
}

// The NWS API always returns temperature as int.
type Fahrenheit int

func (t Fahrenheit) Description() string {
	var label string

	if t >= 80 {
		label = "hot"
	} else if t > 40 && t < 80 {
		label = "moderate"
	} else {
		label = "cold"
	}

	return label
}

func truncateToFourDecimals(num float64) float64 {
	return float64(math.Round(num*10000)) / 10000
}
