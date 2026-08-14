package xpress

import (
	"encoding/binary"
	"testing"
)

// FuzzXpress runs the fuzzed input through both decompressors and asserts
// that neither panics and that neither ever produces more than
// MAX_DECOMPRESSED_FILE bytes of output. The seed corpus is the compressed
// side of the MS-XCA test vectors.
func FuzzXpress(f *testing.F) {
	for _, v := range xpressVectors {
		f.Add(v.data)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		out, err := XpressDecompress(data)
		if err == nil && len(out) > MAX_DECOMPRESSED_FILE {
			t.Fatalf("plain LZ77: %d bytes of output exceeds MAX_DECOMPRESSED_FILE", len(out))
		}

		// The Huffman variant requires the uncompressed size; take it from
		// the start of the input and clamp it to the output limit.
		size := 0
		if len(data) >= 4 {
			s := binary.LittleEndian.Uint32(data)
			if uint64(s) > MAX_DECOMPRESSED_FILE {
				s = MAX_DECOMPRESSED_FILE
			}
			size = int(s)
		}
		out, err = XpressHuffmanDecompress(data, size)
		if err == nil && len(out) > MAX_DECOMPRESSED_FILE {
			t.Fatalf("LZ77+Huffman: %d bytes of output exceeds MAX_DECOMPRESSED_FILE", len(out))
		}
	})
}