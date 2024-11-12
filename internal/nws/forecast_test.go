package nws_test

import (
	"context"
	"errors"
	"testing"

	"github.com/avrekhm/forecast/internal/nws"
	"github.com/google/go-cmp/cmp"
)

type mockNWSer struct {
	gridpointsResponse *nws.GridpointsResponse
	gridPointsError    error
	pointsResponse     *nws.PointsResponse
	pointsError        error
}

func (m mockNWSer) GetGridpointsResponse(ctx context.Context, pointsResponse *nws.PointsResponse) (*nws.GridpointsResponse, error) {
	return m.gridpointsResponse, m.gridPointsError
}

func (m mockNWSer) GetPointsResponse(ctx context.Context, lat, long float64) (*nws.PointsResponse, error) {
	return m.pointsResponse, m.pointsError
}

func TestGetForecast(t *testing.T) {
	ctx := context.Background()

	var outOfMeatballsError = errors.New("out of meatballs")
	var noPointError = errors.New("there is no point really")

	testCases := []struct {
		desc             string
		client           nws.NWSer
		expectedForecast *nws.Forecast
		expectedError    error
	}{
		{
			desc: "correct forecast",
			client: mockNWSer{
				gridpointsResponse: &nws.GridpointsResponse{
					Properties: nws.Properties{
						Periods: []nws.Period{
							{
								ShortForecast: "Cloudy with a chance of meatballs",
								Temperature:   82,
							},
						},
					},
				},
				pointsResponse: &nws.PointsResponse{
					Properties: nws.PointsProperties{
						Forecast: "http://notempty.com",
					},
				},
			},
			expectedForecast: &nws.Forecast{
				ShortForecast: "Cloudy with a chance of meatballs",
				Description:   "hot",
			},
		},
		{
			desc: "gridpoints error expected",
			client: mockNWSer{
				gridpointsResponse: nil,
				gridPointsError:    outOfMeatballsError,
				pointsResponse:     nil,
				pointsError:        nil,
			},
			expectedForecast: nil,
			expectedError:    outOfMeatballsError,
		},
		{
			desc: "points error expected",
			client: mockNWSer{
				gridpointsResponse: nil,
				gridPointsError:    nil,
				pointsResponse:     nil,
				pointsError:        noPointError,
			},
			expectedForecast: nil,
			expectedError:    noPointError,
		},
	}

	for _, test := range testCases {
		t.Run(test.desc, func(t *testing.T) {
			forecast, err := nws.GetForecast(ctx, test.client, 1.0, 1.0)

			if err != nil && !errors.Is(err, test.expectedError) {
				t.Errorf("unexpected error %v, expected %v", err, test.expectedError)
			}

			if !cmp.Equal(test.expectedForecast, forecast) {
				t.Errorf("unexpected values returned: %v instead of expected %v", forecast, test.expectedForecast)
			}
		})
	}
}
