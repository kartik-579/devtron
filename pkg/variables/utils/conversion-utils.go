package utils

import (
	mapset "github.com/deckarep/golang-set"
)

// ToInterfaceArray converts an array of string to an array of interface{}
func ToInterfaceArrayAny[T any](arr []T) []interface{} {
	interfaceArr := make([]interface{}, len(arr))
	for i, v := range arr {
		interfaceArr[i] = v
	}
	return interfaceArr
}

// ToInterfaceArray converts an array of string to an array of interface{}
func ToTypedArrayAny[T any](arr []interface{}) []T {
	typedArr := make([]T, len(arr))
	for i, v := range arr {
		typedArr[i] = v.(T)
	}
	return typedArr
}

// ToInterfaceArray converts an array of string to an array of interface{}
func ToInterfaceArray(arr []string) []interface{} {
	interfaceArr := make([]interface{}, len(arr))
	for i, v := range arr {
		interfaceArr[i] = v
	}
	return interfaceArr
}

// ToStringArray converts an array of interface{} back to an array of string
func ToStringArray(interfaceArr []interface{}) []string {
	stringArr := make([]string, len(interfaceArr))
	for i, v := range interfaceArr {
		stringArr[i] = v.(string)
	}
	return stringArr
}

// ToIntArray converts an array of interface{} back to an array of int
func ToIntArray(interfaceArr []interface{}) []int {
	intArr := make([]int, len(interfaceArr))
	for i, v := range interfaceArr {
		intArr[i] = v.(int)
	}
	return intArr
}

// ToIntArray converts an array of interface{} back to an array of int32
func ToInt32Array(interfaceArr []interface{}) []int32 {
	intArr := make([]int32, len(interfaceArr))
	for i, v := range interfaceArr {
		intArr[i] = v.(int32)
	}
	return intArr
}

func FilterDuplicatesInStringArray(items []string) []string {
	itemsSet := mapset.NewSetFromSlice(ToInterfaceArray(items))
	uniqueItems := ToStringArray(itemsSet.ToSlice())
	return uniqueItems
}

func FilterDuplicates[T any](items []T) []T {
	set := mapset.NewSetFromSlice(ToInterfaceArrayAny(items))
	return ToTypedArrayAny[T](set.ToSlice())
}
