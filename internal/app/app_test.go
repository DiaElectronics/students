package app

import (
	"testing"

	gomock "github.com/golang/mock/gomock"
)

type currencyCodeTest struct {
	code   string
	result error
}

func TestCurrencyValidationCodes(t *testing.T) {
	a, err := NewApplication(nil)
	if err != nil {
		t.Error(err)
		return
	}
	testData := []currencyCodeTest{
		{"ROMA", ErrCodeTooLong},
		{"USD", nil},
		{"EU", ErrCodeTooShort},
		{"", ErrCodeIsEmpty},
		{"543", ErrCodeFormat},
		{"R54", nil},
	}
	for _, val := range testData {
		err := a.ValidateCurrencyCode([]byte(val.code))
		if err != val.result {
			t.Errorf("errors do not match. expected %#v, got %#v", val.result, err)
			return
		}
	}
}

func TestAppWithDal(t *testing.T) {
	ctrl := gomock.NewController(t)
	finishFunc := ctrl.Finish
	defer finishFunc()
	myDal := NewMockDal(ctrl)
	val := float64(1.0) / float64(80.0)

	rateToSave := Rate{CodeFrom: "USD", CodeTo: "RUR", Value: val}
	normalizedRate := Rate{CodeFrom: "RUR", CodeTo: "USD", Value: float64(1.0) / val}
	myDal.EXPECT().SaveRate(normalizedRate).Return(nil).Times(1)
	app, err := NewApplication(myDal)
	if err != nil {
		t.Error(err)
		return
	}
	err = app.SetRate(rateToSave)
	if err != nil {
		t.Error(err)
	}
}
