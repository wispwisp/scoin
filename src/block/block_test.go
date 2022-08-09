package block

import "testing"

// TestBlockHash calls block.Hash, checking for a valid result.
func TestBlockHash(t *testing.T) {
	b := Block{
		Index:    0,
		PrevHash: "none",
		Nonce:    42,
	}
	expected := "2ac2055ff07497cd024e140c209258bb258914c182b4baec47f1f0fa6d821503"
	if got := b.Hash(); got != expected {
		t.Errorf("Hash() = %q, want %q", got, expected)
	}
}
