package utils

import (
	"math/rand"
	"time"
)

var ExcelChar = []string{"", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

func ConvertNumToCol(num int) (string, error) {
	if num < 27 {
		return ExcelChar[num], nil
	}
	k := num % 26
	if k == 0 {
		k = 26
	}
	v := (num - k) / 26
	col, err := ConvertNumToCol(v)
	if err != nil {
		return "", err
	}
	cols := col + ExcelChar[k]
	return cols, nil
}

var letterBytes = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandStringBytes 生成一段随机字符
func RandStringBytes(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
