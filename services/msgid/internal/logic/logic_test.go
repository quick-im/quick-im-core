package logic

import (
	"testing"
	"time"
)

func Benchmark_generateRongCloudMessageID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GenerateRongCloudMessageID(10000000000, "3")
	}
}

func Benchmark_getNextSpinID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = getNextSpinID(time.Now().Unix())
	}
}

func Benchmark_base32Encode(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = base32Encode(7079856935964509836, 125182123728)
	}
}

func TestGenerateRongCloudMessageID(t *testing.T) {
	println(GenerateRongCloudMessageID(1, "tt.args.conversationID"))
}
