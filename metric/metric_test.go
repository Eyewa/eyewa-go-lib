package metric

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestCounter_Add(t *testing.T) {
	option := Option{
		CollectPeriod: 10 * time.Millisecond,
	}

	err := Launch(option)
	assert.Nil(t, err)

	meter := NewMeter("test.meter", nil)
	counter := meter.NewCounter("test.counter")
	counter.Add(1, attribute.Any("test", "testValue"))

	req := httptest.NewRequest(http.MethodGet, "localhost:2222", nil)

	client := http.Client{}
	resp, err := client.Do(req)
	assert.Nil(t, err)

	body, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	fmt.Println(string(body))
}
