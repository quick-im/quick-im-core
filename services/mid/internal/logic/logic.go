package logic

import (
	"bytes"
	"hash/fnv"
	"sync/atomic"
	"time"
)

const (
	MAX_MSG_SEQ              = 0xFFF
	LOW_22_BITS              = 0x3FFFFF
	LOW_16_BITS              = 0xFFFF
	CONVERSITION_TYPE        = 0xF
	LOW_5                    = 0x1F
	HIGH_5            uint64 = 0xF800000000000000
	HIGH_1                   = 0x8000000000000000
)

var (
	lastTimestamp int64
	spinIDCounter uint64 = 1
	charArray            = [33]byte{
		'2', '3', '4', '5', '6', '7', '8', '9',
		'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H',
		'J', 'K', 'L', 'M', 'N', 'P', 'Q', 'R',
		'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'-',
	}
	// spinIDMutex sync.Mutex
	// // 创建一个 FNV 哈希对象
	// hash = fnv.New32()
	// buf  = make([]byte, 19)
	// bufs = bytes.NewBuffer(buf)
)

func GenerateRongCloudMessageID(conversationType uint64, conversationID string) string {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond) // 生成时间戳
	highBits := uint64(timestamp)
	highBits <<= 12
	highBits |= getNextSpinID(timestamp)
	highBits <<= 4
	conversationType &= CONVERSITION_TYPE
	highBits |= conversationType
	// 将整数转换为字节数组并计算哈希值
	// 创建一个 FNV 哈希对象
	hash := fnv.New32()
	hash.Write([]byte(conversationID))
	conversationHash := uint64(hash.Sum32())
	hash.Reset()
	conversationHash &= LOW_22_BITS
	highBits <<= 6
	highBits |= conversationHash >> 16
	lowBits := (conversationHash & LOW_16_BITS) << (64 - 16)
	messageID := base32Encode(highBits, lowBits)
	return messageID
}
func getNextSpinID(timestamp int64) uint64 {
	// spinIDMutex.Lock()
	// defer spinIDMutex.Unlock()
	var currentSpinID uint64                          // 原子操作读取自旋ID计数器
	if timestamp > atomic.LoadInt64(&lastTimestamp) { // 检查时间戳是否增加
		atomic.StoreInt64(&lastTimestamp, timestamp) // 更新最新时间戳
		atomic.StoreUint64(&spinIDCounter, 1)        // 重置自旋ID计数器
		currentSpinID = atomic.LoadUint64(&spinIDCounter)
	} else {
		currentSpinID = atomic.AddUint64(&spinIDCounter, 1) // 自旋ID递增
		if currentSpinID >= (1 << 12) {                     // 自旋ID超过位数限制时重置
			atomic.StoreUint64(&spinIDCounter, 1)
			currentSpinID = atomic.LoadUint64(&spinIDCounter)
		}
	}
	return currentSpinID
}

func base32Encode(highBits, lowBits uint64) string {
	// var sb strings.Builder
	buf := make([]byte, 19)
	bufs := bytes.NewBuffer(buf)
	bufs.Reset()
	var index uint64
	for i := 0; i < 16; i++ {
		// if i > 12 {
		// 	// 在低16位拿出三个字符
		// 	index = (lowBits & (HIGH_5 >> ((i - 13) * 5))) >> (64 - (i-13+1)*5)
		// 	goto processed
		// }
		// if i == 12 {
		// 	// 到最后4bit（第12个字符）时从低16位拿出最高位拼接进来
		// 	index = (highBits<<1 | (lowBits & HIGH_1)) & LOW_5
		// 	lowBits <<= 1
		// 	goto processed
		// }
		// // 高64位拿出11个字符
		// index = (highBits & (HIGH_5 >> (i * 5))) >> (64 - (i+1)*5)
		// goto processed
		switch {
		case i > 12:
			// 在低16位拿出三个字符
			index = (lowBits & (HIGH_5 >> ((i - 13) * 5))) >> (64 - (i-13+1)*5)
		case i == 12:
			// 到最后4bit（第12个字符）时从低16位拿出最高位拼接进来
			index = (highBits<<1 | (lowBits & HIGH_1)) & LOW_5
			lowBits <<= 1
		default:
			// 高64位拿出11个字符
			index = (highBits & (HIGH_5 >> (i * 5))) >> (64 - (i+1)*5)
		}
		// processed:
		bufs.WriteByte(charArray[index])
		// _ = sb.WriteByte(charArray[index])
		if bufs.Len()%5 == 0 {
			bufs.WriteByte(charArray[32])
			// sb.WriteByte(charArray[32])
		}
	}
	return bufs.String()
}
