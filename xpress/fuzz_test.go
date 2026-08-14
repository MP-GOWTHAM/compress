package xpress

import (
	"encoding/binary"
	"testing"
)

// FuzzXpress runs the fuzzed input through both decompressors and asserts
// that neither panics and that neither ever produces more than MaxSize
// bytes of output. The seed corpus is the compressed side of the MS-XCA
// test vectors.
func FuzzXpress(f *testing.F) {
	for _, v := range xpressVectors {
		f.Add(v.data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		out, err := AppendDecompressed(nil, data)
		if err == nil && len(out) > MaxSize {
			t.Fatalf("plain LZ77: %d bytes of output exceeds MaxSize", len(out))
		}
	})
}

// FuzzXpressHuffman is like FuzzXpress but for the Huffman variant, which
// requires the uncompressed size as a separate argument. The size is taken
// from the start of the input and clamped to MaxSize.
func FuzzXpressHuffman(f *testing.F) {
	for _, v := range xpressVectors {
		f.Add(v.data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		size := 0
		if len(data) >= 4 {
			s := binary.LittleEndian.Uint32(data)
			if uint64(s) > MaxSize {
				s = MaxSize
			}
			size = int(s)
		}
		out, err := AppendHDecompressed(nil, data, size)
		if err == nil && len(out) > MaxSize {
			t.Fatalf("LZ77+Huffman: %d bytes of output exceeds MaxSize", len(out))
		}
	})
}
