package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TemplateTestSuite struct {
	suite.Suite
}

func TestTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(TemplateTestSuite))
}

func (suite *TemplateTestSuite) TestTemplateUpgradeV2() {
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

	template := t.Upgrade()

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

func (suite *TemplateTestSuite) TestTemplateUpgradeV3() {
	t := TemplateV2{
		TemplateVersioning: TemplateVersioning{
			Version: 2,
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
				Port: "PORT_22",
			},
			{
				Port: "PORT_80",
			},
		},
		Methods: TemplateMethods{
			Docker: &TemplateMethodDocker{
				Ports: &map[string]string{
					"22": "PORT_22",
					"80": "PORT_80",
				},
			},
		},
	}

	template := t.Upgrade()

	suite.Equal(3, int(template.Version))
	suite.Equal([]TemplatePort{
		{
			Name: "PORT_22",
			Port: "22",
		},
		{
			Name: "PORT_80",
			Port: "80",
		},
	}, template.Ports)
	suite.Empty(template.Env)
	suite.Empty(template.URLs)
}
