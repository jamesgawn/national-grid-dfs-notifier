package natgridapi

import (
	"testing"
)

func TestGetDemandFlexibilityServiceRequirements(t *testing.T) {
	result := GetDemandFlexibilityServiceRequirements()
	if result == nil {
		t.Fail()
	}
}

func TestGetDFSRequirementsForSupplier(t *testing.T) {
	result := GetDFSRequirementsForSupplier("OCTOPUS ENERGY LIMITED")
	if result == nil {
		t.Fail()
	}
}
