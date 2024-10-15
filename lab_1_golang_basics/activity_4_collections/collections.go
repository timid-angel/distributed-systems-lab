package main

import "fmt"

func main() {
	// fixed-size arrays
	arr := [3]int{1, 2, 3}
	fmt.Println("Array:", arr)

	// slices (dynamic-sized lists)
	slice := []int{4, 5, 6}
	slice = append(slice, 7)
	fmt.Println("Slice:", slice)

	// hashmaps
	hashMap := make(map[string]int)
	hashMap["Alice"] = 25
	hashMap["Bob"] = 30
	fmt.Println("Map:", hashMap)
	fmt.Println("Alice's Age:", hashMap["Alices"])

	for i, v := range slice {
		fmt.Printf("Index: %d, Value: %d\n", i, v)
	}

	for key, value := range hashMap {
		fmt.Printf("%s is %d years old\n", key, value)
	}
}
