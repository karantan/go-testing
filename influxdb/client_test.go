package influxdb

import (
	"os"
	"testing"

	client "github.com/influxdata/influxdb/client/v2"
	"github.com/influxdata/influxdb/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockedInfluxSource struct {
	mock.Mock
}

func (m mockedInfluxSource) Query(query client.Query) (r *client.Response, err error) {
	args := m.Called(query)
	return args.Get(0).(*client.Response), args.Error(1)
}

// Functional tests
func TestFunctionalCheckServers(t *testing.T) {
	assert := assert.New(t)
	influx := NewInfluxDSN(os.Getenv("DEV_INFLUX_DSN")) // see Dockerfile for this env
	want := []string{"server2"}
	got, err := CheckServers(influx)
	assert.Equal(want, got)
	assert.NoError(err)
}

// Unit tests
func TestUnitCheckServers(t *testing.T) {
	assert := assert.New(t)
	influx := *new(mockedInfluxSource)

	serversRespTime := `SELECT mean("rtime") FROM "day"."haproxy"
						WHERE "rtime" > 5000
						AND time >= now() - 15m
						AND ("proxy" = 'apache' OR "proxy" = 'varnish')
						GROUP BY "proxy", "host", "provider" fill(null)`
	// setup expectations
	q := client.NewQuery(serversRespTime, "telegraf", "s")
	r := &client.Response{
		Results: []client.Result{
			{
				Series: []models.Row{
					{
						Tags: map[string]string{"host": "server2"},
					},
				},
				Messages: []*client.Message{},
				Err:      "",
			},
		},
		Err: "",
	}
	influx.On("Query", q).Return(r, nil)

	// call the code we are testing
	want := []string{"server2"}
	got, err := CheckServers(influx)

	// assert that the expectations were met
	influx.AssertExpectations(t)
	assert.Equal(want, got)
	assert.NoError(err)

}
