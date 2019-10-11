package gateway

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type WindspeedGatewayTestSuite struct {
	suite.Suite
	gateway *WindspeedGateway
}

func TestWindspeedGatewayTestSuite(t *testing.T) {
	suite.Run(t, new(WindspeedGatewayTestSuite))
}

func (suite *WindspeedGatewayTestSuite) SetupTest() {
	suite.gateway = NewWindspeedGateway()
}

func (suite *WindspeedGatewayTestSuite) TestGatewayShouldReturnTemperatureForAGivenDate() {
	// Given

	// When

	// Then
}

func (suite *WindspeedGatewayTestSuite) TestGatewayShouldReturnErrorWhenStatusCodeIsNotOK() {
	// Given

	// When

	// Then
}
