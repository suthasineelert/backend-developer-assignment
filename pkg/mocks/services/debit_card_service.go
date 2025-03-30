// Code generated by mockery v2.53.3. DO NOT EDIT.

package mocks

import (
	models "backend-developer-assignment/app/models"

	mock "github.com/stretchr/testify/mock"
)

// DebitCardService is an autogenerated mock type for the DebitCardService type
type DebitCardService struct {
	mock.Mock
}

// CreateCardWithDetails provides a mock function with given fields: cardWithDetails
func (_m *DebitCardService) CreateCardWithDetails(cardWithDetails *models.DebitCardWithDetails) error {
	ret := _m.Called(cardWithDetails)

	if len(ret) == 0 {
		panic("no return value specified for CreateCardWithDetails")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.DebitCardWithDetails) error); ok {
		r0 = rf(cardWithDetails)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// DeleteCard provides a mock function with given fields: cardID
func (_m *DebitCardService) DeleteCard(cardID string) error {
	ret := _m.Called(cardID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteCard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(cardID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetCardByID provides a mock function with given fields: cardID
func (_m *DebitCardService) GetCardByID(cardID string) (*models.DebitCard, error) {
	ret := _m.Called(cardID)

	if len(ret) == 0 {
		panic("no return value specified for GetCardByID")
	}

	var r0 *models.DebitCard
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.DebitCard, error)); ok {
		return rf(cardID)
	}
	if rf, ok := ret.Get(0).(func(string) *models.DebitCard); ok {
		r0 = rf(cardID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DebitCard)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(cardID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCardWithDetailByID provides a mock function with given fields: cardID
func (_m *DebitCardService) GetCardWithDetailByID(cardID string) (*models.DebitCardWithDetails, error) {
	ret := _m.Called(cardID)

	if len(ret) == 0 {
		panic("no return value specified for GetCardWithDetailByID")
	}

	var r0 *models.DebitCardWithDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (*models.DebitCardWithDetails, error)); ok {
		return rf(cardID)
	}
	if rf, ok := ret.Get(0).(func(string) *models.DebitCardWithDetails); ok {
		r0 = rf(cardID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*models.DebitCardWithDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(cardID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCardWithDetailByUserID provides a mock function with given fields: userID
func (_m *DebitCardService) GetCardWithDetailByUserID(userID string) ([]*models.DebitCardWithDetails, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetCardWithDetailByUserID")
	}

	var r0 []*models.DebitCardWithDetails
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*models.DebitCardWithDetails, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) []*models.DebitCardWithDetails); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.DebitCardWithDetails)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetCardsByUserID provides a mock function with given fields: userID
func (_m *DebitCardService) GetCardsByUserID(userID string) ([]*models.DebitCard, error) {
	ret := _m.Called(userID)

	if len(ret) == 0 {
		panic("no return value specified for GetCardsByUserID")
	}

	var r0 []*models.DebitCard
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]*models.DebitCard, error)); ok {
		return rf(userID)
	}
	if rf, ok := ret.Get(0).(func(string) []*models.DebitCard); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*models.DebitCard)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateCard provides a mock function with given fields: card, name, color, borderColor
func (_m *DebitCardService) UpdateCard(card *models.DebitCard, name string, color string, borderColor string) error {
	ret := _m.Called(card, name, color, borderColor)

	if len(ret) == 0 {
		panic("no return value specified for UpdateCard")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*models.DebitCard, string, string, string) error); ok {
		r0 = rf(card, name, color, borderColor)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewDebitCardService creates a new instance of DebitCardService. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewDebitCardService(t interface {
	mock.TestingT
	Cleanup(func())
}) *DebitCardService {
	mock := &DebitCardService{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
