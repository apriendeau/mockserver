package mockserver_test

import (
	"net/http"
	"testing"

	"github.com/apriendeau/mockserver"
	"github.com/stretchr/testify/assert"
)

func TestTemplateServer(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	o := testItem{
		Item:  "boat",
		Thing: "wood",
		Count: 2,
	}
	tmpl, err := mockserver.NewTemplate("base", "{{.Count}} {{.Item}}s are made of {{.Thing}}", o)
	assert.NoError(err)
	server := tmpl.Server(200, "plain/text")
	defer server.Close()

	resp, err := http.Get(server.URL + "/testing")
	assert.NoError(err)
	assert.Equal(200, resp.StatusCode)

	b, err := parseResp(resp.Body)
	assert.NoError(err)
	assert.Equal("2 boats are made of wood", b)
}
