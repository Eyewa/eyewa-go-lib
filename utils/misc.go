package utils

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
