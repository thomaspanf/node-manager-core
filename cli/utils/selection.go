package utils

import (
	"fmt"
	"strconv"
	"strings"
)

// An option that can be selected from a list of choices in the CLI
type SelectionOption[DataType any] struct {
	// The underlying element this option represents
	Element *DataType

	// The human-readable ID of the option, used for non-interactive selection
	ID string

	// The text to display for this option when listing the available options
	Display string
}

// Parse a comma-separated list of indices to select in a multi-index operation
func ParseIndexSelection[DataType any](selectionString string, options []SelectionOption[DataType]) ([]*DataType, error) {
	// Select all
	if selectionString == "" {
		selectedElements := make([]*DataType, len(options))
		for i, option := range options {
			selectedElements[i] = option.Element
		}
		return selectedElements, nil
	}

	// Trim spaces
	elements := strings.Split(selectionString, ",")
	trimmedElements := make([]string, len(elements))
	for i, element := range elements {
		trimmedElements[i] = strings.TrimSpace(element)
	}

	// Process elements
	optionLength := uint64(len(options))
	seenIndices := map[uint64]bool{}
	selectedElements := []*DataType{}
	for _, element := range trimmedElements {
		before, after, found := strings.Cut(element, "-")
		if !found {
			// Handle non-ranges
			index, err := strconv.ParseUint(element, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing index '%s': %w", element, err)
			}
			index--

			// Make sure it's in the list of options
			if index >= optionLength {
				return nil, fmt.Errorf("selection '%s' is too large", element)
			}

			// Add it if it's new
			_, exists := seenIndices[index]
			if !exists {
				seenIndices[index] = true
				selectedElements = append(selectedElements, options[index].Element)
			}
		} else {
			// Handle ranges
			start, err := strconv.ParseUint(before, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing range start in '%s': %w", element, err)
			}
			end, err := strconv.ParseUint(after, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing range end in '%s': %w", element, err)
			}
			start--
			end--

			// Make sure the start and end are in the list of options
			if end <= start {
				return nil, fmt.Errorf("range end for '%s' is not greater than the start", element)
			}
			if start >= optionLength {
				return nil, fmt.Errorf("range start for '%s' is too large", element)
			}
			if end >= optionLength {
				return nil, fmt.Errorf("range end for '%s' is too large", element)
			}

			// Add each index if it's new
			for index := start; index <= end; index++ {
				_, exists := seenIndices[index]
				if !exists {
					seenIndices[index] = true
					selectedElements = append(selectedElements, options[index].Element)
				}
			}
		}
	}
	return selectedElements, nil
}

// Parse a comma-separated list of option IDs to select in a multi-index operation
func ParseOptionIDs[DataType any](selectionString string, options []SelectionOption[DataType]) ([]*DataType, error) {
	elements := strings.Split(selectionString, ",")

	// Trim spaces
	trimmedElements := make([]string, len(elements))
	for i, element := range elements {
		trimmedElements[i] = strings.TrimSpace(element)
	}

	// Remove duplicates
	uniqueElements := make([]string, 0, len(elements))
	seenIndices := map[string]bool{}
	for _, element := range trimmedElements {
		_, exists := seenIndices[element]
		if !exists {
			uniqueElements = append(uniqueElements, element)
			seenIndices[element] = true
		}
	}

	// Validate
	selectedElements := make([]*DataType, len(uniqueElements))
	for i, element := range uniqueElements {
		// Make sure it's in the list of options
		found := false
		for _, option := range options {
			if option.ID == element {
				found = true
				selectedElements[i] = option.Element
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("element '%s' is not a valid option", element)
		}
	}
	return selectedElements, nil
}
