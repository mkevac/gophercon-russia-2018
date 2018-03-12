package geoborder

import (
	"testing"

	"fmt"

	"github.com/golang/geo/s2"
)

func TestSearch(t *testing.T) {
	index := NewIndex(13 /* ~ 1km */)

	index.AddPolygon(1, []s2.LatLng{
		s2.LatLngFromDegrees(55.77116, 37.65289),
		s2.LatLngFromDegrees(55.7729, 37.588),
		s2.LatLngFromDegrees(55.73406, 37.58422),
		s2.LatLngFromDegrees(55.73522, 37.65666),
	}) // Moscow Center

	index.AddPolygon(2, []s2.LatLng{
		s2.LatLngFromDegrees(55.69113, 37.32192),
		s2.LatLngFromDegrees(55.69345, 37.24708),
		s2.LatLngFromDegrees(55.65628, 37.24159),
		s2.LatLngFromDegrees(55.65899, 37.31849),
	}) // Odintsovo

	// Points
	// 55.75648, 37.62199 Moscow Center
	// 55.79047, 37.78816 Izmaylovo

	found, _ := index.Search(37.62199, 55.75648)
	if len(found) != 1 || found[0] != 1 {
		t.Fatal("error while searching at moscow center")
	}

	found, _ = index.Search(37.78816, 55.79047)
	if len(found) != 0 {
		t.Fatal("error while searching at izmaylovo")
	}
}

func TestSearchNearest(t *testing.T) {
	index := NewIndex(13 /* ~ 1km */)

	index.AddPolygon(1, []s2.LatLng{
		s2.LatLngFromDegrees(55.77116, 37.65289),
		s2.LatLngFromDegrees(55.7729, 37.588),
		s2.LatLngFromDegrees(55.73406, 37.58422),
		s2.LatLngFromDegrees(55.73522, 37.65666),
	}) // Moscow Center

	index.AddPolygon(2, []s2.LatLng{
		s2.LatLngFromDegrees(55.69113, 37.32192),
		s2.LatLngFromDegrees(55.69345, 37.24708),
		s2.LatLngFromDegrees(55.65628, 37.24159),
		s2.LatLngFromDegrees(55.65899, 37.31849),
	}) // Odintsovo

	// Points
	// 55.75648, 37.62199 Moscow Center
	// 55.79047, 37.78816 Izmaylovo

	found, _ := index.SearchNearest(37.78816, 55.79047)
	fmt.Println("found", found)
	if len(found) != 1 || found[0] != 1 {
		t.Fatal("nearest from izmaylovo is not moscow city center")
	}
}
