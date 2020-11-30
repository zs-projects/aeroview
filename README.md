# Aeroview

## Use case 1: Load huge datasets 

1. Load to `map[string][]byte`.
2. Load to `Memtable` ( Read from disk )
3. Load to `Memtable` using `mmap`
4. Load only header from `Memtable`
5. Load only header from `Memtable` + `mmap` on the values and `on-page-miss` callback
6. Same as 3, 4 ,5 for `perfect hashing` -> Examples with `Recsplit` and `CHD`

Measurements ( for medium ): 

- Peak Memory
- CPU time 
- Lookup time

### Timeline : 

1. Implement experiment 1 et 2 et 3 ( With complete documentation )
2. Write medium article
3. Implement 4 5 -> Write medium article
4. Article to make `Memtable` more user friendly.
5. Article for perfect hashing
