package subcommand

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/stretchr/testify/assert"
)

func TestListExternalMonitorMetricNames(t *testing.T) {
	sv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resJson := `
		{
			"monitors": [
				{
					"id": "1XXXXXXXXXX",
					"type": "host"
				},
				{
					"id": "2XXXXXXXXXX",
					"type": "external",
					"service": "test",
					"url": "https://example.com"
				},
				{
					"id": "3XXXXXXXXXX",
					"type": "external",
					"service": "test",
					"url": "https://example.com"
				}
			]
		}
		`
		_, _ = w.Write([]byte(resJson))
	}))
	defer sv.Close()

	client, _ := mackerel.NewClientWithOptions("dummy", sv.URL, false)

	c := &serviceCommand{
		baseCommand: baseCommand{
			client: client,
		},
		name: "test",
	}

	metricNames, _ := c.listExternalMonitorMetricNames()

	assert.Equal(t, 2, len(metricNames))
	assert.Equal(t, "__externalhttp.responsetime.2XXXXXXXXXX", metricNames[0])
	assert.Equal(t, "__externalhttp.responsetime.3XXXXXXXXXX", metricNames[1])
}
