package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type TemplateTestSuite struct {
	suite.Suite
}

func TestTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}

func (suite *TemplateTestSuite) TestTemplateUpgrade() {
	t := TemplateV1{
		TemplateVersioning: TemplateVersioning{
			Version: 1,
		},
		Env: []TemplateEnv{
			{
				Type:    "port",
				Name:    "PORT_22",
				Default: "22",
			},
			{
				Type:    "port",
				Name:    "PORT_80",
				Default: "80",
			},
		},
		URLs: []URL{
			{
				Port: "22",
			},
			{
				Port: "80",
			},
		},
		Methods: TemplateMethods{
			Docker: &TemplateMethodDocker{
				Ports: &map[string]string{
					"22": "22",
					"80": "80",
				},
			},
		},
	}

	bytes, err := yaml.Marshal(&t)
	suite.Require().NoError(err)

	var template Template
	err = yaml.Unmarshal(bytes, &template)
	suite.Require().NoError(err)
	suite.Equal(2, int(template.Version))
	suite.Equal(&map[string]string{
		"22": "PORT_22",
		"80": "PORT_80",
	}, template.Methods.Docker.Ports)
	suite.Equal([]URL{
		{
			Port: "PORT_22",
		},
		{
			Port: "PORT_80",
		},
	}, template.URLs)
}
