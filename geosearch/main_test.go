package geosearch

import (
	"testing"
)

func prepare() *Index {
	i := NewIndex(13 /* ~ 1km */)
	i.AddUser(1, 14.1313, 14.1313)
	i.AddUser(2, 14.1314, 14.1314)
	i.AddUser(3, 14.1311, 14.1311)
	i.AddUser(10, 14.2313, 14.2313)
	i.AddUser(11, 14.0313, 14.0313)
	return i
}

func TestSearch(t *testing.T) {
	indx := prepare()

	found, _ := indx.Search(14.1313, 14.1313, 2000)
	if len(found) != 3 {
		t.Fatal("error while searching with radius 2000")
	}

	found, _ = indx.Search(14.1313, 14.1313, 20000)
	if len(found) != 5 {
		t.Fatal("error while searching with radius 20000")
	}
}

func TestSearchFaster(t *testing.T) {
	indx := prepare()

	found, _ := indx.SearchFaster(14.1313, 14.1313, 2000)
	if len(found) != 3 {
		t.Fatal("error while searching with radius 2000")
	}

	found, _ = indx.SearchFaster(14.1313, 14.1313, 20000)
	if len(found) != 5 {
		t.Fatal("error while searching with radius 20000")
	}
}

var res []uint32

func BenchmarkSearch(b *testing.B) {
	indx := prepare()

	for i := 0; i < b.N; i++ {
		res, _ = indx.Search(14.1313, 14.1313, 50000)
	}
}

func BenchmarkSearchFaster(b *testing.B) {
	indx := prepare()

	for i := 0; i < b.N; i++ {
		res, _ = indx.SearchFaster(14.1313, 14.1313, 50000)
	}
}
