package utils_test

import (
	"io"
	"strings"
	"testing"

	"github.com/gsaaraujo/ecommerce-api-scenario-1/internal/utils"
	"github.com/stretchr/testify/suite"
)

type ParseJSONBodySuite struct {
	suite.Suite
}

func (p *ParseJSONBodySuite) Test1() {
	p.Run("when parsing a JSON body, returns body parsed", func() {
		reader := strings.NewReader(`{"name": "John Doe"}`)
		readCloser := io.NopCloser(reader)

		body, err := utils.ParseJSONBody[map[string]any](readCloser)
		p.Require().NoError(err)

		p.Equal("John Doe", body["name"])
	})
}

func TestParseJSONBody(t *testing.T) {
	suite.Run(t, new(ParseJSONBodySuite))
}
