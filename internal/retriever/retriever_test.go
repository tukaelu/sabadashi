package retriever

import (
	"testing"

	"github.com/mackerelio/mackerel-client-go"
	"github.com/stretchr/testify/assert"
)

func TestRetrieve(t *testing.T) {
	fn := func() ([]mackerel.MetricValue, error) {
		var metrics []mackerel.MetricValue
		metrics = append(metrics, mackerel.MetricValue{Time: 1672498800, Value: 1})
		metrics = append(metrics, mackerel.MetricValue{Time: 1672498860, Value: 2})
		metrics = append(metrics, mackerel.MetricValue{Time: 1672498920, Value: 3})
		return metrics, nil
	}
	ret, _ := Retrieve(nil, "foo.bar", fn)

	assert.Equal(t, "foo.bar", ret[0].Name)
	assert.Equal(t, int64(1672498800), ret[0].Time)
	assert.Equal(t, 1, ret[0].Value)
	assert.Equal(t, "foo.bar", ret[1].Name)
	assert.Equal(t, int64(1672498860), ret[1].Time)
	assert.Equal(t, 2, ret[1].Value)
	assert.Equal(t, "foo.bar", ret[2].Name)
	assert.Equal(t, int64(1672498920), ret[2].Time)
	assert.Equal(t, 3, ret[2].Value)
}
