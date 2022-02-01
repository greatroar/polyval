package polyval

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"golang.org/x/exp/rand"
)

func unhex(s string) []byte {
	p, err := hex.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return p
}

func elem(s string) fieldElement {
	var z fieldElement
	z.setBytes(unhex(s))
	return z
}

// TestCtmulCommutative tests that ctmul is commutative,
// a required property for multiplication.
func TestCtmulCommutative(t *testing.T) {
	seed := uint64(time.Now().UnixNano())
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 1e6; i++ {
		x, y := rng.Uint64(), rng.Uint64()
		xy1, xy0 := ctmul(x, y)
		yx1, yx0 := ctmul(y, x)
		if xy1 != yx1 || xy0 != yx0 {
			t.Fatalf("%#0.16x*%#0.16x: (%#0.16x, %#0.16x) != (%#0.16x, %#0.16x)",
				x, y, xy1, xy0, yx1, yx0)
		}
	}
}

// TestPolyvalRFCVectors tests polyval using test vectors from
// RFC 8452.
func TestPolyvalRFCVectors(t *testing.T) {
	for i, tc := range []struct {
		H []byte
		X [][]byte
		r []byte
	}{
		// POLYVAL(H, X_1)
		{
			H: unhex("25629347589242761d31f826ba4b757b"),
			X: [][]byte{
				unhex("4f4f95668c83dfb6401762bb2d01a262"),
			},
			r: unhex("cedac64537ff50989c16011551086d77"),
		},
		// POLYVAL(H, X_1, X_2)
		{
			H: unhex("25629347589242761d31f826ba4b757b"),
			X: [][]byte{
				unhex("4f4f95668c83dfb6401762bb2d01a262"),
				unhex("d1a24ddd2721d006bbe45f20d3c9f362"),
			},
			r: unhex("f7a3b47b846119fae5b7866cf5e5b77e"),
		},
	} {
		g, _ := New(tc.H) // generic
		p, _ := New(tc.H) // specialized
		for _, x := range tc.X {
			p.Update(x)
			polymulGeneric(&g.y, &g.h, x)
		}
		want := tc.r
		got := p.Sum(nil)
		if !bytes.Equal(got, want) {
			t.Fatalf("#%d: expected %x, got %x", i, want, got)
		}
		got = g.Sum(nil)
		if !bytes.Equal(got, want) {
			t.Fatalf("#%d: expected %x, got %x", i, want, got)
		}
	}
}

// TestPolyvalVectors tests polyval using the Google-provided
// test vectors.
//
// See https://github.com/google/hctr2/blob/2a80dc7f742127b1f68f02b310975ac7928ae25e/test_vectors/ours/Polyval/Polyval.json
func TestPolyvalVectors(t *testing.T) {
	type vector struct {
		Cipher struct {
			Cipher      string `json:"cipher"`
			BlockCipher struct {
				Cipher  string `json:"cipher"`
				Lengths struct {
					Block int `json:"block"`
					Key   int `json:"key"`
					Nonce int `json:"nonce"`
				} `json:"lengths"`
			} `json:"block_cipher"`
		} `json:"cipher"`
		Description string `json:"description"`
		Input       struct {
			Key     string `json:"key_hex"`
			Tweak   string `json:"tweak_hex"`
			Message string `json:"message_hex"`
			Nonce   string `json:"nonce_hex"`
		} `json:"input"`
		Plaintext  string `json:"plaintext_hex"`
		Ciphertext string `json:"ciphertext_hex"`
		Hash       string `json:"hash_hex"`
	}

	var vecs []vector
	buf, err := os.ReadFile(filepath.Join("testdata", "polyval.json"))
	if err != nil {
		t.Fatal(err)
	}
	err = json.Unmarshal(buf, &vecs)
	if err != nil {
		t.Fatal(err)
	}
	for i, v := range vecs {
		key := unhex(v.Input.Key)
		g, _ := New(key) // generic
		p, _ := New(key) // specialized
		msg := unhex(v.Input.Message)
		for len(msg) > 0 {
			p.Update(msg[0:16])
			polymulGeneric(&g.y, &g.h, msg[0:16])
			msg = msg[16:]
		}
		want := unhex(v.Hash)
		got := p.Sum(nil)
		if !bytes.Equal(want, got) {
			t.Fatalf("#%d: (%s): expected %x, got %x",
				i, v.Description, want, got)
		}
		got = g.Sum(nil)
		if !bytes.Equal(got, want) {
			t.Fatalf("#%d: (%s): expected %x, got %x",
				i, v.Description, want, got)
		}
	}
}

// TestMarshal tests Polyval's MarshalBinary and UnmarshalBinary
// methods.
func TestMarshal(t *testing.T) {
	h, _ := New(make([]byte, 16))
	block := make([]byte, 16)
	seed := uint64(time.Now().UnixNano())
	rng := rand.New(rand.NewSource(seed))
	for i := 0; i < 500; i++ {
		rng.Read(block)

		// Save the current digest and state.
		prevSum := h.Sum(nil)
		prev, _ := h.MarshalBinary()

		// Update the state and save the digest.
		h.Update(block)
		curSum := h.Sum(nil)

		// Read back the first state and check that we get the
		// same results.
		var h2 Polyval
		h2.UnmarshalBinary(prev)
		if got := h2.Sum(nil); !bytes.Equal(got, prevSum) {
			t.Fatalf("#%d: exepected %x, got %d", i, prevSum, got)
		}
		h2.Update(block)
		if got := h2.Sum(nil); !bytes.Equal(got, curSum) {
			t.Fatalf("#%d: exepected %x, got %d", i, curSum, got)
		}
	}
}

// TestDoubleRFCVectors tests double over the set of vectors from
// RFC 8452.
//
// See https://datatracker.ietf.org/doc/html/rfc8452#appendix-A
func TestDoubleRFCVectors(t *testing.T) {
	for i, tc := range []struct {
		input  fieldElement
		output fieldElement
	}{
		{
			input:  elem("01000000000000000000000000000000"),
			output: elem("02000000000000000000000000000000"),
		},
		{
			input:  elem("9c98c04df9387ded828175a92ba652d8"),
			output: elem("3931819bf271fada0503eb52574ca572"),
		},
	} {
		want := tc.output
		if got := tc.input.double(); got != want {
			t.Fatalf("#%d: expected %#x, got %#x", i, want, got)
		}
	}
}

// TestDoubleRustVectors tests double over the set of vectors
// from RustCrypto.
//
// See https://github.com/RustCrypto/universal-hashes/blob/5361f44a1162bd0d84e6560b6e30c7cb445e683f/polyval/src/double.rs#L58
func TestDoubleRustVectors(t *testing.T) {
	r := elem("01000000000000000000000000000000")

	for i, v := range []fieldElement{
		elem("02000000000000000000000000000000"),
		elem("04000000000000000000000000000000"),
		elem("08000000000000000000000000000000"),
		elem("10000000000000000000000000000000"),
		elem("20000000000000000000000000000000"),
		elem("40000000000000000000000000000000"),
		elem("80000000000000000000000000000000"),
		elem("00010000000000000000000000000000"),
		elem("00020000000000000000000000000000"),
		elem("00040000000000000000000000000000"),
		elem("00080000000000000000000000000000"),
		elem("00100000000000000000000000000000"),
		elem("00200000000000000000000000000000"),
		elem("00400000000000000000000000000000"),
		elem("00800000000000000000000000000000"),
		elem("00000100000000000000000000000000"),
		elem("00000200000000000000000000000000"),
		elem("00000400000000000000000000000000"),
		elem("00000800000000000000000000000000"),
		elem("00001000000000000000000000000000"),
		elem("00002000000000000000000000000000"),
		elem("00004000000000000000000000000000"),
		elem("00008000000000000000000000000000"),
		elem("00000001000000000000000000000000"),
		elem("00000002000000000000000000000000"),
		elem("00000004000000000000000000000000"),
		elem("00000008000000000000000000000000"),
		elem("00000010000000000000000000000000"),
		elem("00000020000000000000000000000000"),
		elem("00000040000000000000000000000000"),
		elem("00000080000000000000000000000000"),
		elem("00000000010000000000000000000000"),
		elem("00000000020000000000000000000000"),
		elem("00000000040000000000000000000000"),
		elem("00000000080000000000000000000000"),
		elem("00000000100000000000000000000000"),
		elem("00000000200000000000000000000000"),
		elem("00000000400000000000000000000000"),
		elem("00000000800000000000000000000000"),
		elem("00000000000100000000000000000000"),
		elem("00000000000200000000000000000000"),
		elem("00000000000400000000000000000000"),
		elem("00000000000800000000000000000000"),
		elem("00000000001000000000000000000000"),
		elem("00000000002000000000000000000000"),
		elem("00000000004000000000000000000000"),
		elem("00000000008000000000000000000000"),
		elem("00000000000001000000000000000000"),
		elem("00000000000002000000000000000000"),
		elem("00000000000004000000000000000000"),
		elem("00000000000008000000000000000000"),
		elem("00000000000010000000000000000000"),
		elem("00000000000020000000000000000000"),
		elem("00000000000040000000000000000000"),
		elem("00000000000080000000000000000000"),
		elem("00000000000000010000000000000000"),
		elem("00000000000000020000000000000000"),
		elem("00000000000000040000000000000000"),
		elem("00000000000000080000000000000000"),
		elem("00000000000000100000000000000000"),
		elem("00000000000000200000000000000000"),
		elem("00000000000000400000000000000000"),
		elem("00000000000000800000000000000000"),
		elem("00000000000000000100000000000000"),
		elem("00000000000000000200000000000000"),
		elem("00000000000000000400000000000000"),
		elem("00000000000000000800000000000000"),
		elem("00000000000000001000000000000000"),
		elem("00000000000000002000000000000000"),
		elem("00000000000000004000000000000000"),
		elem("00000000000000008000000000000000"),
		elem("00000000000000000001000000000000"),
		elem("00000000000000000002000000000000"),
		elem("00000000000000000004000000000000"),
		elem("00000000000000000008000000000000"),
		elem("00000000000000000010000000000000"),
		elem("00000000000000000020000000000000"),
		elem("00000000000000000040000000000000"),
		elem("00000000000000000080000000000000"),
		elem("00000000000000000000010000000000"),
		elem("00000000000000000000020000000000"),
		elem("00000000000000000000040000000000"),
		elem("00000000000000000000080000000000"),
		elem("00000000000000000000100000000000"),
		elem("00000000000000000000200000000000"),
		elem("00000000000000000000400000000000"),
		elem("00000000000000000000800000000000"),
		elem("00000000000000000000000100000000"),
		elem("00000000000000000000000200000000"),
		elem("00000000000000000000000400000000"),
		elem("00000000000000000000000800000000"),
		elem("00000000000000000000001000000000"),
		elem("00000000000000000000002000000000"),
		elem("00000000000000000000004000000000"),
		elem("00000000000000000000008000000000"),
		elem("00000000000000000000000001000000"),
		elem("00000000000000000000000002000000"),
		elem("00000000000000000000000004000000"),
		elem("00000000000000000000000008000000"),
		elem("00000000000000000000000010000000"),
		elem("00000000000000000000000020000000"),
		elem("00000000000000000000000040000000"),
		elem("00000000000000000000000080000000"),
		elem("00000000000000000000000000010000"),
		elem("00000000000000000000000000020000"),
		elem("00000000000000000000000000040000"),
		elem("00000000000000000000000000080000"),
		elem("00000000000000000000000000100000"),
		elem("00000000000000000000000000200000"),
		elem("00000000000000000000000000400000"),
		elem("00000000000000000000000000800000"),
		elem("00000000000000000000000000000100"),
		elem("00000000000000000000000000000200"),
		elem("00000000000000000000000000000400"),
		elem("00000000000000000000000000000800"),
		elem("00000000000000000000000000001000"),
		elem("00000000000000000000000000002000"),
		elem("00000000000000000000000000004000"),
		elem("00000000000000000000000000008000"),
		elem("00000000000000000000000000000001"),
		elem("00000000000000000000000000000002"),
		elem("00000000000000000000000000000004"),
		elem("00000000000000000000000000000008"),
		elem("00000000000000000000000000000010"),
		elem("00000000000000000000000000000020"),
		elem("00000000000000000000000000000040"),
		elem("00000000000000000000000000000080"),
		elem("010000000000000000000000000000c2"),
	} {
		want := v
		got := r.double()
		if got != want {
			t.Fatalf("#%d: expected %#x, got %#x", i, want, got)
		}
		r = got
	}
}

var (
	byteSink []byte
	elemSink fieldElement
)

func BenchmarkDouble(b *testing.B) {
	x := fieldElement{
		hi: rand.Uint64(),
		lo: rand.Uint64(),
	}
	for i := 0; i < b.N; i++ {
		x = x.double()
	}
	elemSink = x
}

func BenchmarkPolyval(b *testing.B) {
	b.SetBytes(16)
	p, _ := New(unhex("01000000000000000000000000000000"))
	x := make([]byte, p.BlockSize())
	for i := 0; i < b.N; i++ {
		p.Update(x)
	}
	byteSink = p.Sum(nil)
}