package geoborder

import (
	"testing"

	"fmt"

	"github.com/golang/geo/s2"
)

func TestSearch(t *testing.T) {
	index := NewIndex(13 /* ~ 1km */)

	fmt.Println("first add")

	index.AddPolygon(1, []s2.LatLng{
		s2.LatLng{Lat: 55.854, Lng: 37.419},
		s2.LatLng{Lat: 55.877, Lng: 37.809},
		s2.LatLng{Lat: 55.637, Lng: 37.854},
		s2.LatLng{Lat: 55.610, Lng: 37.37},
	}) // Moscow Center

	fmt.Println("first add finished")

	fmt.Println("second add")

	index.AddPolygon(2, []s2.LatLng{
		s2.LatLng{37.1846008, 55.6806816},
		s2.LatLng{37.2518921, 55.6078327},
		s2.LatLng{37.3988342, 55.6628694},
		s2.LatLng{37.3329163, 55.7271101},
		s2.LatLng{37.1846008, 55.6806816},
	}) // Odintsovo

	fmt.Println("second add finished")

	// Points
	// 55.736, 37.63435 Moscow Center
	// 55.78892, 37.78198 Izmaylovo

	fmt.Println("first search")

	found, _ := index.Search(37.63435, 55.736)
	if len(found) != 1 && found[0] != 1 {
		t.Fatal("error while searching at moscow center")
	}

	fmt.Println("first search finished")

	fmt.Println("second test")

	found, _ = index.Search(37.78198, 55.78892)
	if len(found) != 0 {
		t.Fatal("error while searching at izmaylovo")
	}

	fmt.Println("second test finished")

}

func TestSearchNearest(t *testing.T) {
	index := NewIndex(13 /* ~ 1km */)

	index.AddPolygon(1, []s2.LatLng{
		s2.LatLng{37.5595093, 55.7649858},
		s2.LatLng{37.5787354, 55.6849398},
		s2.LatLng{37.7565765, 55.7290434},
		s2.LatLng{37.6714325, 55.8020528},
		s2.LatLng{37.5595093, 55.7649858},
	}) // Moscow Center

	index.AddPolygon(2, []s2.LatLng{
		s2.LatLng{37.1846008, 55.6806816},
		s2.LatLng{37.2518921, 55.6078327},
		s2.LatLng{37.3988342, 55.6628694},
		s2.LatLng{37.3329163, 55.7271101},
		s2.LatLng{37.1846008, 55.6806816},
	}) // Odintsovo

	// Points
	// 55.736, 37.63435 Moscow Center
	// 55.78892, 37.78198 Izmaylovo

	found, _ := index.SearchNearest(37.78198, 55.78892)
	if len(found) != 1 && found[0] != 1 {
		t.Fatal("nearest from izmaylovo is not moscow city center")
	}
}
