package influxdb

import (
	"net/url"

	client "github.com/influxdata/influxdb/client/v2"
)

type InfluxSource interface {
	Query(client.Query) (*client.Response, error)
}

// Influx holds InfluxDB client
type InfluxClient struct {
	client.Client
}

// NewInfluxDSN returns InfluxDB client ready for use
func NewInfluxDSN(dsn string) InfluxClient {
	u, err := url.Parse(dsn)
	if err != nil {
		log.Fatal(err)
	}
	pass := ""
	if p, ok := u.User.Password(); ok {
		pass = p
	}
	user := u.User.Username()
	u.User = nil

	influxC, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:               u.String(),
		Username:           user,
		Password:           pass,
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	return InfluxClient{influxC}
}

func RunQuery(influx InfluxSource, command, database, precision string) (res []client.Result, err error) {
	q := client.NewQuery(command, database, precision)
	r, err := influx.Query(q)
	if err != nil {
		return
	}
	if r.Error() != nil {
		err = r.Error()
		return
	}
	return r.Results, nil
}
