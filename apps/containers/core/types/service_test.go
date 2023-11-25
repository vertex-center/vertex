package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

type ServiceTestSuite struct {
	suite.Suite
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) TestServiceUpgrade() {
	s := ServiceV1{
		ServiceVersioning: ServiceVersioning{
			Version: 1,
		},
		Env: []ServiceEnv{
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
		Methods: ServiceMethods{
			Docker: &ServiceMethodDocker{
				Ports: &map[string]string{
					"22": "22",
					"80": "80",
				},
			},
		},
	}

	bytes, err := yaml.Marshal(&s)
	suite.Require().NoError(err)

	var service Service
	err = yaml.Unmarshal(bytes, &service)
	suite.Require().NoError(err)
	suite.Equal(2, int(service.Version))
	suite.Equal(&map[string]string{
		"22": "PORT_22",
		"80": "PORT_80",
	}, service.Methods.Docker.Ports)
	suite.Equal([]URL{
		{
			Port: "PORT_22",
		},
		{
			Port: "PORT_80",
		},
	}, service.URLs)
}
