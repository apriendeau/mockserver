package mockserver_test

import (
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/apriendeau/mockserver"
	"github.com/stretchr/testify/assert"
)

type testItem struct {
	Item  string
	Count int
	Thing string
}

type unsupportedItem struct {
	Item    string
	Channel chan bool
}

func parseResp(body io.Reader) (string, error) {
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func TestNewMockServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	server := mockserver.Simple(200, "application/json", `{"hello":"world"}`)
	defer server.Close()
	resp, err := http.Get(server.URL + "/testing")
	assert.NoError(err)
	assert.Equal(200, resp.StatusCode)
	b, err := parseResp(resp.Body)
	assert.NoError(err)
	assert.Equal("{\"hello\":\"world\"}", b)
}

func TestNewMockJsonServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	o := testItem{
		Item:  "testing",
		Count: 1234,
	}
	server, err := mockserver.JSON(201, o)
	assert.NoError(err)
	defer server.Close()
	resp, err := http.Get(server.URL + "/testing")
	assert.NoError(err)
	assert.Equal(201, resp.StatusCode)
	b, err := parseResp(resp.Body)
	assert.NoError(err)
	samp, err := json.Marshal(o)
	assert.NoError(err)
	assert.Equal(string(samp), b)
}

func TestJSONMockFail(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	bad := unsupportedItem{
		Item: "hello",
	}
	_, err := mockserver.JSON(201, bad)
	assert.Error(err)
}

func TestNewMockXMLServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	o := testItem{
		Item:  "testing",
		Count: 1234,
	}
	server, err := mockserver.XML(201, o)
	assert.NoError(err)
	defer server.Close()
	resp, err := http.Get(server.URL + "/testing")
	assert.NoError(err)
	assert.Equal(201, resp.StatusCode)
	b, err := parseResp(resp.Body)
	assert.NoError(err)
	samp, err := xml.Marshal(o)
	assert.NoError(err)
	assert.Equal(string(samp), b)
}

func TestXMLMockFail(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	bad := unsupportedItem{
		Item: "hello",
	}
	_, err := mockserver.XML(201, bad)
	assert.Error(err)
}
