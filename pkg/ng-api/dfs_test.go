package ngapi

import (
	"testing"
)

func TestGetDemandFlexibilityServiceRequirements(t *testing.T) {
	result := GetDemandFlexibilityServiceRequirements()
	if result == nil {
		t.Fail()
	}
}
