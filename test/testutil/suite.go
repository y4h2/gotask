package testutil

import "github.com/stretchr/testify/suite"

type Suite struct {
	suite.Suite
}

func (s *Suite) Given(str string) {
	s.T().Log("Given " + str)
}

func (s *Suite) When(str string) {
	s.T().Log("When " + str)
}

func (s *Suite) Then(str string) {
	s.T().Log("Then " + str)
}
