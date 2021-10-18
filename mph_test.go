package mph

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
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

func testMPH(t *testing.T, keys []string) {
	tab := New(keys)
	for i, k := range keys {
		if got := tab.Query(k); got != int32(i) {
			t.Errorf("Lookup(%q)=%v, want %v", k, got, i)
		}
	}
}

func TestMPH(t *testing.T) {
	keys := loadKeys(t)
	testMPH(t, keys)
}

func TestMPHRandomSubsets(t *testing.T) {
	keys := loadKeys(t)

	iterations := 100 * len(keys)

	for i := 0; i < iterations; i++ {
		perm := rand.Perm(rand.Intn(len(keys)))
		subkeys := make([]string, len(perm))
		for i, v := range perm {
			subkeys[i] = keys[v]
		}

		t.Run(fmt.Sprintf("%d-%d\n", i, len(subkeys)), func(t *testing.T) { testMPH(t, subkeys) })
	}
}

var sink int32

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
	m := make(map[string]int32, len(keys))
	for i, k := range keys {
		m[k] = int32(i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, k := range keys {
			sink += m[k]
		}
	}
}
