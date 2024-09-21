package model

import (
	"fmt"
	"strings"
	"time"
)

func FormatExpiryDate(inputTime string) string {
	// Format the time in Indonesian custom format
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, inputTime)
	if err != nil {
		fmt.Println("Error parsing time:", err)
		t = time.Now().Add(24 * time.Hour)
	}

	// Format the time in Indonesian custom format
	formattedTime := t.Format("02 January 2006, pukul 15:04 WIB")

	// Replace English month names with Indonesian names
	indonesianMonths := map[string]string{
		"January": "Januari", "February": "Februari", "March": "Maret",
		"April": "April", "May": "Mei", "June": "Juni",
		"July": "Juli", "August": "Agustus", "September": "September",
		"October": "Oktober", "November": "November", "December": "Desember",
	}

	// Replace English month names with Indonesian ones
	for en, id := range indonesianMonths {
		formattedTime = strings.ReplaceAll(formattedTime, en, id)
	}

	return formattedTime
}

func replace(str, old, new string) string {
	return strings.ReplaceAll(str, old, new)
}

func FormatRP(amount float64) string {
	formatted := fmt.Sprintf("%.2f", amount)

	parts := strings.Split(formatted, ".")
	integerPart := parts[0]
	fractionalPart := parts[1]

	var result []string
	for i, r := range reverse(integerPart) {
		if i > 0 && i%3 == 0 {
			result = append(result, ".")
		}
		result = append(result, string(r))
	}

	final := "Rp" + reverse(strings.Join(result, "")) + "," + fractionalPart
	return final
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
