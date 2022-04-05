package influxdb

func CheckServers(influx InfluxSource) (slowServers []string, err error) {
	serversRespTime := `SELECT mean("rtime") FROM "day"."haproxy"
						WHERE "rtime" > 5000
						AND time >= now() - 15m
						AND ("proxy" = 'apache' OR "proxy" = 'varnish')
						GROUP BY "proxy", "host", "provider" fill(null)`
	results, err := RunQuery(influx, serversRespTime, "telegraf", "s")
	if err != nil {
		return
	}
	for _, results := range results {
		for _, server := range results.Series {
			slowServers = append(slowServers, server.Tags["host"])
			log.Warnw("High response time", server.Tags["host"])
		}
	}
	return
}
