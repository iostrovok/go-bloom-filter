package bloomfilter

import (
// "testing"
)

var bloomBanch *BloomFilter

// func init() {
// 	var err error
// 	bloomBanch, err = New(10000*10000*1000, 0.001)
// 	if err != nil {
// 		panic(err)
// 	}
// }

// func BenchmarkAdd(b *testing.B) {
// 	for _, s := range testData {
// 		bloomBanch.Add([]byte(s))
// 	}
// }

// func BenchmarkCheck(b *testing.B) {
// 	for _, s := range testData {
// 		bloomBanch.Check([]byte(s))
// 		if !bloomBanch.Check([]byte(s)) {
// 			panic("wring banch test!")
// 		}
// 	}
// }
