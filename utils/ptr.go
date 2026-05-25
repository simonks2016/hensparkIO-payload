package utils

func Ptr[T comparable](data T) *T {
	return &data
}

func StrPtr(data string) *string {
	return Ptr[string](data)
}

func IntPtr(data int) *int {
	return Ptr[int](data)
}
func Int8Ptr(data int8) *int8 {
	return Ptr[int8](data)
}
func Int16Ptr(data int16) *int16 {
	return Ptr[int16](data)
}
func Int32Ptr(data int32) *int32 {
	return Ptr[int32](data)
}
func Int64Ptr(data int64) *int64 {
	return Ptr[int64](data)
}
func BoolPtr(data bool) *bool {
	return Ptr[bool](data)
}
func Float32Ptr(data float32) *float32 {
	return Ptr[float32](data)
}
func Float64Ptr(data float64) *float64 {
	return Ptr[float64](data)
}
