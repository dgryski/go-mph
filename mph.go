package mph

import (
	"sort"

	"github.com/dgryski/go-metro"
)

type Table struct {
	values []int
	seeds  []int64
}

type entry struct {
	key  string
	idx  int
	hash uint64
}

func New(keys []string) *Table {
	size := uint64(nextPower2(len(keys)))

	h := make([][]entry, size)
	for idx, k := range keys {
		hash := metro.Hash64Str(k, 0)
		i := hash % size
		// idx+1 so we can identify empty entries in the table with 0
		h[i] = append(h[i], entry{k, idx + 1, hash})
	}

	sort.Slice(h, func(i, j int) bool { return len(h[i]) > len(h[j]) })

	values := make([]int, size)
	seeds := make([]int64, size)

	var hidx int
	for hidx = 0; hidx < len(h) && len(h[hidx]) > 1; hidx++ {
		subkeys := h[hidx]

		var seed uint64
		entries := make(map[uint64]int)

	newseed:
		for {
			seed++
			for _, k := range subkeys {
				i := xorshiftMult64(k.hash+seed) % size
				if entries[i] == 0 && values[i] == 0 {
					// looks free, claim it
					entries[i] = k.idx
					continue
				}

				// found a collision, reset and try a new seed
				for k := range entries {
					delete(entries, k)
				}
				continue newseed
			}

			// made it through; everything got placed
			break
		}

		// mark subkey spaces as claimed
		for k, v := range entries {
			values[k] = v
		}

		// and assign this seed value for every subkey
		for _, k := range subkeys {
			i := k.hash % size
			seeds[i] = int64(seed)
		}
	}

	// find the unassigned entries in the table
	var free []int
	for i := range values {
		if values[i] == 0 {
			free = append(free, i)
		} else {
			// decrement idx as this is now the final value for the table
			values[i]--
		}
	}

	for len(h[hidx]) > 0 {
		k := h[hidx][0]
		i := k.hash % size
		hidx++

		// take a free slot
		dst := free[0]
		free = free[1:]

		// claim it; -1 because of the +1 at the start
		values[dst] = k.idx - 1

		// store offset in seed as a negative; -1 so even slot 0 is negative
		seeds[i] = -int64(dst + 1)
	}

	return &Table{
		values: values,
		seeds:  seeds,
	}
}

func (t *Table) Query(k string) int {
	size := uint64(len(t.values))
	hash := metro.Hash64Str(k, 0)
	i := hash & (size - 1)
	seed := t.seeds[i]
	if seed < 0 {
		return t.values[-seed-1]
	}

	i = xorshiftMult64(uint64(seed)+hash) & (size - 1)
	return t.values[i]
}

func xorshiftMult64(x uint64) uint64 {
	x ^= x >> 12 // a
	x ^= x << 25 // b
	x ^= x >> 27 // c
	return x * 2685821657736338717
}

func nextPower2(n int) int {
	i := 1
	for i < n {
		i *= 2
	}
	return i
}
