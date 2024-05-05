package input

import "fmt"

// ==============
// === Errors ===
// ==============

type InvalidArgCountError struct {
	ExpectedCount int
	ActualCount   int
}

func (e *InvalidArgCountError) Error() string {
	return fmt.Sprintf("Incorrect argument count - expected %d but have %d", e.ExpectedCount, e.ActualCount)
}

func NewInvalidArgCountError(expectedCount int, actualCount int) *InvalidArgCountError {
	return &InvalidArgCountError{
		ExpectedCount: expectedCount,
		ActualCount:   actualCount,
	}
}

// ==================
// === Validation ===
// ==================

// Validate command argument count
func ValidateArgCount(argCount int, expectedCount int) *InvalidArgCountError {
	if argCount != expectedCount {
		return NewInvalidArgCountError(expectedCount, argCount)
	}
	return nil
}
