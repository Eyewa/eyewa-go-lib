package utils

import "time"

// ConvertIntToPointer convert an int to *int
func ConvertIntToPointer(i int) *int {
	return (*int)(&i)
}

// ConvertStringToPointer convert an string to *string
func ConvertStringToPointer(s string) *string {
	return (*string)(&s)
}

// ConvertFloat64ToPointer convert an float64 to *float64
func ConvertFloat64ToPointer(f float64) *float64 {
	return (*float64)(&f)
}

// NowRFC3339 return time in Dubai/Asia timezone in time.RFC3339 format
func NowRFC3339() string {
	return time.Now().UTC().Add(4 * time.Hour).Format(time.RFC3339)
}

// Now return time in Dubai/Asia timezone.
func Now() time.Time {
	return time.Now().UTC().Add(4 * time.Hour)
}
