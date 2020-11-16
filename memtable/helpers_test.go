package memtable

import "fmt"

func getSampleMap(size int) map[string][]byte {
	mp := make(map[string][]byte, size)
	for i := 0; i < size; i++ {
		key := fmt.Sprintf("key-%v", i)
		value := []byte(fmt.Sprintf("value-%v", i))
		mp[key] = value

	}
	return mp
}
