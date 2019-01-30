package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"fyne.io/fyne"

	"fyne.io/fyne/app"
	"fyne.io/fyne/widget"
)

// Program to generate a TIN (in German: Steueridentifikationsnummer)
// - It is a 11 digit number and the 11th digit is a check digit
// - No leading 0 allowed, except it is a test TIN number
// - The leading 10 digits must contain exactly one number twice and the rest individual
//		E.g.: 0247629135 (Test TIN), checkdigit is 8, number 2 is existing twice and the rest is individual
// - Check digit is calculated from the leading 10 digits via #calcCheckDigitTIN
func main() {
	createGUI()
}

// Creates the GUI.
func createGUI() {
	isTest := false

	app := app.New()
	app.SetIcon(nil)
	w := app.NewWindow("TINGen")
	w.CenterOnScreen()
	w.Resize(fyne.Size.Add(fyne.NewSize(320, 80), fyne.NewSize(0, 0)))

	l := widget.NewLabelWithStyle("", fyne.TextAlignCenter, fyne.TextStyle{Monospace: true})
	l.Resize(fyne.Size.Add(fyne.NewSize(300, 50), fyne.NewSize(0, 0)))
	l.Alignment = fyne.TextAlignCenter

	c := widget.NewCheck("Test TIN", func(checked bool) {
		isTest = !isTest
	})

	b := widget.NewButton("Create TIN", func() {
		tinIntSlice := generateTIN(isTest)
		tinStrSlice := make([]string, 11)
		for i, val := range tinIntSlice {
			tinStrSlice[i] = strconv.Itoa(val)
		}
		l.SetText(strings.Join(tinStrSlice, ""))
	})
	w.SetContent(widget.NewVBox(l, c, b))
	w.ShowAndRun()
}

// Generates a German TIN (in German: Steueridentifikationsnummer).
// See: https://de.wikipedia.org/wiki/Steuerliche_Identifikationsnummer
// See: https://download.elster.de/download/schnittstellen/Pruefung_der_Steuer_und_Steueridentifikatsnummer.pdf
func generateTIN(testTIN bool) []int {
	// Possible set of digits
	possibleDigits := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	// Resultset for the created TIN
	tin := make([]int, len(possibleDigits)+1)

	// Decide wether one digit should exist twice or thrice in the TIN
	isTwice := isTwice()
	// Get the indices for the lucky digit
	luckyIndices := getLuckyIndices(isTwice)
	fmt.Println(luckyIndices)
	// Determine digit which should exist twice or thrice
	luckyDigit := determineLuckyDigit(luckyIndices, testTIN)
	// Remove the lucky digit from the set of possible digits (value set to -1)
	removeDigitFromPossibleDigits(luckyDigit, possibleDigits)

	for i := 0; i < len(possibleDigits); i++ {
		// If we hit a lucky index we set the lucky digit
		// Because #determineLuckyDigit does not allow the value 0 for a non-testing TIN we do not need to check for 0 here
		if (isTwice && (i == luckyIndices[0] || i == luckyIndices[1])) || (!isTwice && (i == luckyIndices[0] || i == luckyIndices[1] || i == luckyIndices[2])) {
			tin[i] = luckyDigit
		} else {
			// Else we pseudorandomly take one remaining possible digit and remove it afterwards
			// from the set of possible digits (value set to -1)
			for {
				possibleDigit, _ := rand.Int(rand.Reader, big.NewInt(int64(len(possibleDigits))))
				selectedDigit := possibleDigits[int(possibleDigit.Uint64())]
				if selectedDigit != -1 {
					isDigitSave := false
					// If we do not have a TIN for testing purposes we need to assure that the first index does not get the value 0.
					if !testTIN {
						if selectedDigit == 0 && i != 0 || selectedDigit != 0 && i == 0 || selectedDigit != 0 && i != 0 {
							isDigitSave = true
						}
					} else {
						isDigitSave = true
					}
					if isDigitSave {
						tin[i] = selectedDigit
						removeDigitFromPossibleDigits(selectedDigit, possibleDigits)
						break
					}
				}

				//if selectedDigit != -1 /*&& (!testTIN && selectedDigit != 0 && i != 0)*/ {
				//	if (!testTIN && selectedDigit != 0 && i == 0) || (!testTIN && selectedDigit == 0 && i != 0) {

				//	}
				//	tin[i] = selectedDigit
				//	removeDigitFromPossibleDigits(selectedDigit, possibleDigits)
				//	break
				//}
			}
		}
	}
	tin[len(tin)-1] = calcCheckDigitTIN(tin)
	return tin
}

func removeDigitFromPossibleDigits(digit int, possibleDigits []int) {
	for i := 0; i < len(possibleDigits); i++ {
		if possibleDigits[i] == digit {
			possibleDigits[i] = -1
		}
	}
}

// Calculates two or three indices for the digit which exists twice or thrice in the TIN.
// If one digit exists twice, the third index gets the value -1.
// If a digit occures thrice, this function assures that the third index is not besides the other indices.
func getLuckyIndices(isTwice bool) []int {
	const len = 10
	// Field for the indices, either two, or three
	indices := make([]int, 3)

	index1, _ := rand.Int(rand.Reader, big.NewInt(len))
	indices[0] = int(index1.Uint64())

	index2, _ := rand.Int(rand.Reader, big.NewInt(len))
	for index1.Cmp(index2) == 0 {
		index2, _ = rand.Int(rand.Reader, big.NewInt(len))
	}
	indices[1] = int(index2.Uint64())

	if !isTwice {
		index3, _ := rand.Int(rand.Reader, big.NewInt(len))

		i1 := index1.Uint64()
		i2 := index2.Uint64()
		i3 := index3.Uint64()

		indicesAreNeighbours := checkIndicesForBeingNeighbours(i1, i2, i3)

		for index1.Cmp(index3) == 0 || index2.Cmp(index3) == 0 || indicesAreNeighbours {
			index3, _ = rand.Int(rand.Reader, big.NewInt(len))
			i3 = index3.Uint64()
			indicesAreNeighbours = checkIndicesForBeingNeighbours(i1, i2, i3)
		}
		indices[2] = int(index3.Uint64())
	} else {
		indices[2] = -1
	}

	return indices
}

// Helper function to check if the given indices are positioned besides each other.
func checkIndicesForBeingNeighbours(i1 uint64, i2 uint64, i3 uint64) bool {
	a := i1 < i2 && i2 < i3 && (i3-i2 == 1 && i2-i1 == 1)
	b := i1 < i3 && i3 < i2 && (i2-i3 == 1 && i3-i1 == 1)
	c := i2 < i1 && i1 < i3 && (i3-i1 == 1 && i1-i2 == 1)
	d := i3 < i1 && i1 < i2 && (i2-i1 == 1 && i1-i3 == 1)
	e := i3 < i2 && i2 < i1 && (i1-i2 == 1 && i2-i3 == 1)
	f := i2 < i3 && i3 < i1 && (i1-i3 == 1 && i3-i2 == 1)
	return a || b || c || d || e || f
}

// Determines a digit between 0 and 9 which will exist twice or thrice in the TIN.
// If the TIN is not for testing purposes and a lucky index is 0 the lucky digit is not 0.
func determineLuckyDigit(luckyIndices []int, testTIN bool) int {
	luckyDigitBig, _ := rand.Int(rand.Reader, big.NewInt(10))
	luckyDigit := int(luckyDigitBig.Uint64())
	if !testTIN && (luckyIndices[0] == 0 || luckyIndices[1] == 0 || luckyIndices[2] == 0) {
		for luckyDigit == 0 {
			luckyDigitBig, _ = rand.Int(rand.Reader, big.NewInt(10))
			luckyDigit = int(luckyDigitBig.Uint64())
		}
	}
	return luckyDigit
}

// Pseudorandomly decides wether one digit in the TIN should exists twice or thrice.
// If #isTwice returns true, one digit should exist twice in the TIN, otherwise this digit
// exists thrice.
func isTwice() bool {
	isTwice, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		panic(err)
	}
	if isTwice.Uint64() == 0 {
		return true
	}
	return false
}

// Calculates the check digit for a given TIN.
func calcCheckDigitTIN(tin []int) int {
	const m = 10
	const n = 11
	sum := 0
	product := m

	for _, val := range tin {
		sum = (val + product) % m
		if sum == 0 {
			sum = m
		}
		product = (2 * sum) % n
	}

	checkDigit := n - product
	if checkDigit == 10 {
		return 0
	}
	return checkDigit
}
