package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"strconv"
)

func ByteToInt(value []byte) (int, error) {
	return strconv.Atoi(string(value))
}

// func float32ToString(value float32) string {
// 	return strconv.FormatFloat(float64(value), 'g', -1, 64)
// }

func float64ToString(value float64) string {
	return strconv.FormatFloat(value, 'g', -1, 64)
}

func UintToString(value uint) string {
	return fmt.Sprintf("%d", db.Count)
}

func IntToString(value int) string {
	return strconv.Itoa(value)
}

func StringToFloat64(value string) float64 {
	f, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0.0
	}
	return f
}

func Round(f float64) float64 {
	return (math.Round(f*100) / 100)
}

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}

func ByteArrayToInt(data []byte) int {
	return int(binary.LittleEndian.Uint64(data))
}
