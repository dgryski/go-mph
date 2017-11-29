package mph

import (
	"bufio"
	"flag"
	"os"
	"testing"
)

var keysFile = flag.String("keys", "", "load keys datafile")

func loadKeys(tb testing.TB) []string {

	if *keysFile != "" {
		return loadBigKeys(tb, *keysFile)
	}

	return []string{
		"foo",
		"bar",
		"baz",
		"qux",
		"zot",
		"frob",
		"zork",
		"zeek",
	}
}

func loadBigKeys(tb testing.TB, filename string) []string {
	f, err := os.Open(filename)
	if err != nil {
		tb.Fatalf("unable to keys file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var k []string
	for scanner.Scan() {
		k = append(k, scanner.Text())
	}

	return k
}

func TestMPH(t *testing.T) {

	keys := loadKeys(t)

	tab := New(keys)

	for i, k := range keys {
		got := tab.Query(k)

		if got != i {
			t.Errorf("Lookup(%q)=%v, want %v", k, got, i)
		}
	}
}

var sink int

func BenchmarkMPH(b *testing.B) {
	keys := loadKeys(b)
	tab := New(keys)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		for _, k := range keys {
			sink += tab.Query(k)
		}
	}
}

func BenchmarkMap(b *testing.B) {
	keys := loadKeys(b)

	m := make(map[string]int, len(keys))

	for i, k := range keys {
		m[k] = i
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {

		for _, k := range keys {
			sink += m[k]
		}
	}
}
