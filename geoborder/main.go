package geoborder

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/google/btree"
)

type IndexItem struct {
	cellID         s2.CellID
	polygonsInCell []uint32
}

func (ii IndexItem) Less(than btree.Item) bool {
	return uint64(ii.cellID) < uint64(than.(IndexItem).cellID)
}

type Index struct {
	storageLevel int
	bt           *btree.BTree
	polygons     map[uint32]*s2.Polygon
}

func NewIndex(storageLevel int) *Index {
	return &Index{
		storageLevel: storageLevel,
		bt:           btree.New(35),
		polygons:     make(map[uint32]*s2.Polygon),
	}
}

func (i *Index) AddPolygon(polygonID uint32, vertices []s2.LatLng) error {

	points := func() (res []s2.Point) {
		for _, vertex := range vertices {
			point := s2.PointFromLatLng(vertex)
			res = append(res, point)
		}
		return
	}()

	loop := s2.LoopFromPoints(points)
	loop.Normalize()
	polygon := s2.PolygonFromLoops([]*s2.Loop{loop})

	coverer := s2.RegionCoverer{MinLevel: i.storageLevel, MaxLevel: i.storageLevel}
	cells := coverer.Covering(loop)

	for _, cell := range cells {
		ii := IndexItem{cellID: cell}
		item := i.bt.Get(ii)
		if item != nil {
			ii = item.(IndexItem)
		}
		ii.polygonsInCell = append(ii.polygonsInCell, polygonID)
		i.bt.ReplaceOrInsert(ii)
	}

	i.polygons[polygonID] = polygon

	return nil
}

func (i *Index) Search(lon, lat float64) ([]uint32, error) {
	latlng := s2.LatLngFromDegrees(lat, lon)
	cellID := s2.CellIDFromLatLng(latlng)
	cellIDOnStorageLevel := cellID.Parent(i.storageLevel)

	item := i.bt.Get(IndexItem{cellID: cellIDOnStorageLevel})
	if item != nil {
		return item.(IndexItem).polygonsInCell, nil
	}

	return []uint32{}, nil
}

func in(haystack []s2.CellID, needle s2.CellID) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}
	return false
}

func (i *Index) searchNextLevel(radiusEdges []s2.CellID, alreadyVisited *[]s2.CellID) (found []IndexItem, searched []s2.CellID) {
	for _, radiusEdge := range radiusEdges {
		neighbors := radiusEdge.AllNeighbors(i.storageLevel)
		for _, neighbor := range neighbors {

			if in(*alreadyVisited, neighbor) {
				continue
			}
			*alreadyVisited = append(*alreadyVisited, neighbor)
			searched = append(searched, neighbor)

			item := i.bt.Get(IndexItem{cellID: neighbor})
			if item != nil {
				found = append(found, item.(IndexItem))
			}
		}
	}
	return
}

func (i *Index) filter(lon, lat float64, found []IndexItem) []uint32 {

	var minDistance s1.ChordAngle
	var minPolygon uint32

	cell := s2.CellFromLatLng(s2.LatLngFromDegrees(lat, lon))

	for _, f := range found {
		for _, polygonID := range f.polygonsInCell {
			polygon := i.polygons[polygonID]
			for i := 0; i < polygon.NumEdges(); i++ {
				edge := polygon.Edge(i)
				distance := cell.DistanceToEdge(edge.V0, edge.V1)
				if distance == 0 {
					minDistance = distance
					minPolygon = polygonID
				} else if distance < minDistance {
					minDistance = distance
					minPolygon = polygonID
				}
			}
		}
	}

	return []uint32{minPolygon}
}

func (i *Index) SearchNearest(lon, lat float64) ([]uint32, error) {

	latlng := s2.LatLngFromDegrees(lat, lon)
	cellID := s2.CellIDFromLatLng(latlng)
	cellIDOnStorageLevel := cellID.Parent(i.storageLevel)

	item := i.bt.Get(IndexItem{cellID: cellIDOnStorageLevel})
	if item != nil {
		return item.(IndexItem).polygonsInCell, nil
	}

	alreadyVisited := []s2.CellID{cellID}
	var found []IndexItem
	searched := []s2.CellID{cellID}

	for {
		found, searched = i.searchNextLevel(searched, &alreadyVisited)
		if len(searched) == 0 {
			break
		}
		if len(found) > 0 {
			return i.filter(lon, lat, found), nil
		}
	}

	return []uint32{}, nil
}
