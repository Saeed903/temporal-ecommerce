package app

import (
	"time"

	"github.com/stretchr/testify/suite"

	"go.temporal.io/sdk/testsuite"
)
type UnitTestSuite struct {
	suite.Suite
	testsuite.WorkflowTestSuite

	env *testsuite.TestWorkflowEnvironment 
}

func (s *UnitTestSuite) SetupTest() {
	s.env = s.NewTestWorkflowEnvironment()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.env.AssertExpectations(s.T())
}

func (s *UnitTestSuite) Test_AddToCart() {
	cart := CartState{Items: make([]CartItem, 0)}

	s.env.RegisterDelayedCallback(func() {
		res, err := s.env.QueryWorkflow("getCart")
		s.NoError(err)
		err = res.Get(&cart)
		s.NoError(err)
		s.Equal(len(cart.Items), 0)

		update := AddToCartSignal{
			Route: RouteTypes.ADD_TO_CART,
			Item: CartItem{ProductId: 1, Quantity: 1},
		} 
		s.env.SignalWorkflow(SignalChannels.ADD_TO_CART_CHANNEL, update)
	}, time.Millisecond * 1)

	s.env.RegisterDelayedCallback(func ()  {
		res, err := s.env.QueryWorkflow("getCart")
		s.NoError(err)
		err = res.Get(&cart)
		s.NoError(err)
		s.Equal(1, len(cart.Items))
	}, time.Microsecond * 2)

	s.env.ExecuteWorkflow(CartWorkflow, cart)

	s.True(s.env.IsWorkflowCompleted())
}

