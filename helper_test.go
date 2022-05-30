package main

import (
	"fmt"
	"math"
	"strconv"
	"testing"
)

///////////////////
// Float64ToString
///////////////////

func BenchmarkFloat64ToString_fmtSprintf(b *testing.B) {
	var i float64
	for i = 0; i < float64(b.N); i++ {
		val := fmt.Sprintf("%f", i)
		_ = val
	}
}

func BenchmarkFloat64ToString_strconvFormatFloat(b *testing.B) {
	var i float64
	for i = 0; i < float64(b.N); i++ {
		val := strconv.FormatFloat(i, 'g', -1, 64)
		_ = val
	}
}

///////////////////
// Float32ToString
///////////////////

func BenchmarkFloat32ToString_fmtSprintf(b *testing.B) {
	var i float32
	for i = 0; i < float32(b.N); i++ {
		val := fmt.Sprintf("%f", i)
		_ = val
	}
}

func BenchmarkFloat32ToString_strconvFormatFloat(b *testing.B) {
	var i float64
	for i = 0; i < float64(b.N); i++ {
		val := strconv.FormatFloat(float64(i), 'g', -1, 32)
		_ = val
	}
}

///////////////////
// IntToString
///////////////////

func BenchmarkIntToString_strconvItoa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		val := strconv.Itoa(i)
		_ = val
	}
}

func BenchmarkIntToString_FormatInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		val := strconv.FormatInt(int64(i), 10)
		_ = val
	}
}

///////////////////
// Float64Round
///////////////////

func BenchmarkFloat64Round_mathRound(b *testing.B) {
	var i float64
	for i = 0; i < float64(b.N); i++ {
		val := math.Round(i*100) / 100
		_ = val
	}
}

func BenchmarkFloat64Round_fmtSprintf(b *testing.B) {
	var i float64
	for i = 0; i < float64(b.N); i++ {
		val := fmt.Sprintf("%.2f", i)
		_ = val
	}
}

///////////////////
// Float32Round
///////////////////

// func BenchmarkFloat32Round_mathRound(b *testing.B) {
// 	var i float32
// 	for i = 0; i < float32(b.N); i++ {
// 		val := math.Round(float64(i*100)) / 100
// 		_ = val
// 	}
// }

func BenchmarkFloat32Round_fmtSprintf(b *testing.B) {
	var i float32
	for i = 0; i < float32(b.N); i++ {
		val := fmt.Sprintf("%.2f", i)
		_ = val
	}
}

// func BenchmarkFloat32Round_intFloat(b *testing.B) {
// 	var i float32
// 	for i = 0; i < float32(b.N); i++ {
// 		val := float32(int(i*100)) / 100
// 		_ = val
// 	}
// }

///////////////////
// StringsConcat
///////////////////

// fmt.Sprintf("<P=%d,%d><C=FFFFFF><B=%d,%d><C><P=0,0>", c.OSDX, c.OSDY, c.OSDWidth, c.OSDHeight)

func BenchmarkStringsConcat_fmtSprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		val := fmt.Sprintf("<P=%d,%d><C=FFFFFF><B=%d,%d><C><P=0,0>", i, i*2, i*3, i*4)
		_ = val
	}
}

func BenchmarkStringsConcat_strconvItoa(b *testing.B) {
	for i := 0; i < b.N; i++ {
		val := "<P=" + strconv.Itoa(i) + "," + strconv.Itoa(i*2) + "><C=FFFFFF><B=" + strconv.Itoa(i*3) + "," + strconv.Itoa(i*4) + "><C><P=0,0>"
		_ = val
	}
}
