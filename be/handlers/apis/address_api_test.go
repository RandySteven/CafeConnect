package apis

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type AddressApiTestSuite struct {
	suite *suite.Suite
}

func (a *AddressApiTestSuite) T() *testing.T {
	return nil
}

func (a *AddressApiTestSuite) SetT(t *testing.T) {
}

func (a *AddressApiTestSuite) SetS(suite suite.TestingSuite) {
}

func (a *AddressApiTestSuite) SetupTest() {

}

func TestAddressApiTestSuite(t *testing.T) {
	suite.Run(t, new(AddressApiTestSuite))
}
