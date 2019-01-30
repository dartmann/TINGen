package main

import (
	"testing"
)

// Tests #calcCheckDigitTIN of TINGen. The example TINs and their check digits are taken from official documents.
// See: https://ec.europa.eu/taxation_customs/tin/specs/FS-TIN%20Algorithms-Public.docx
// See: https://download.elster.de/download/schnittstellen/Pruefung_der_Steuer_und_Steueridentifikatsnummer.pdf
func Test_calcCheckDigitTIN(t *testing.T) {
	// TIN1: 0247629135, Check digit 1: 8
	// TIN2: 2695437182, Check digit 2: 7
	// TIN3: 8609574271, Check digit 3: 9
	// TIN4: 6592997048, Check digit 4: 9
	tINsWithoutCheckDigits := [][]int{{0, 2, 4, 7, 6, 2, 9, 1, 3, 5}, {2, 6, 9, 5, 4, 3, 7, 1, 8, 2}, {8, 6, 0, 9, 5, 7, 4, 2, 7, 1}, {6, 5, 9, 2, 9, 9, 7, 0, 4, 8}}
	checkDigitsOfTINs := []int{8, 7, 9, 9}
	for i, val := range tINsWithoutCheckDigits {
		shouldBeCheckDigit := checkDigitsOfTINs[i]
		calculatedCheckDigit := calcCheckDigitTIN(val)
		if calculatedCheckDigit != shouldBeCheckDigit {
			t.Errorf("Wrong check digit: %d for given TIN: %d. Should be: %d", calculatedCheckDigit, val, shouldBeCheckDigit)
		}
	}
}
