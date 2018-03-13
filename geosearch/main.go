package geosearch

import (
	"github.com/golang/geo/s1"
	"github.com/golang/geo/s2"
	"github.com/google/btree"
)

const EarthRadiusM = 6371010.0

type userList struct {
	cellID s2.CellID
	list   []uint32
}

func (ul userList) Less(than btree.Item) bool {
	return uint64(ul.cellID) < uint64(than.(userList).cellID)
}

type Index struct {
	storageLevel int
	bt           *btree.BTree
}

func NewIndex(storageLevel int) *Index {
	return &Index{
		storageLevel: storageLevel,
		bt:           btree.New(35),
	}
}

func (i *Index) AddUser(userID uint32, lon, lat float64) error {
	latlng := s2.LatLngFromDegrees(lat, lon)
	cellID := s2.CellIDFromLatLng(latlng)
	cellIDOnStorageLevel := cellID.Parent(i.storageLevel)

	ul := userList{cellID: cellIDOnStorageLevel, list: nil}

	item := i.bt.Get(ul)
	if item != nil {
		ul = item.(userList)
	}
	ul.list = append(ul.list, userID)

	i.bt.ReplaceOrInsert(ul)

	return nil
}

func (i *Index) Search(lon, lat float64, radius uint32) ([]uint32, error) {

	latlng := s2.LatLngFromDegrees(lat, lon)
	centerPoint := s2.PointFromLatLng(latlng)

	centerAngle := float64(radius) / EarthRadiusM
	cap := s2.CapFromCenterAngle(centerPoint, s1.Angle(centerAngle))

	rc := s2.RegionCoverer{MaxLevel: i.storageLevel, MinLevel: i.storageLevel}
	cu := rc.Covering(cap)

	var result []uint32

	for _, cellID := range cu {
		item := i.bt.Get(userList{cellID: cellID})
		if item != nil {
			result = append(result, item.(userList).list...)
		}
	}

	return result, nil
}

func (i *Index) SearchFaster(lon, lat float64, radius uint32) ([]uint32, error) {

	latlng := s2.LatLngFromDegrees(lat, lon)
	centerPoint := s2.PointFromLatLng(latlng)

	centerAngle := float64(radius) / EarthRadiusM
	cap := s2.CapFromCenterAngle(centerPoint, s1.Angle(centerAngle))

	rc := s2.RegionCoverer{MaxLevel: i.storageLevel}
	cu := rc.Covering(cap)

	var result []uint32

	for _, cellID := range cu {
		if cellID.Level() < i.storageLevel {
			begin := cellID.ChildBeginAtLevel(i.storageLevel)
			end := cellID.ChildEndAtLevel(i.storageLevel)
			i.bt.AscendRange(userList{cellID: begin}, userList{cellID: end.Next()}, func(item btree.Item) bool {
				result = append(result, item.(userList).list...)
				return true
			})
		} else {
			item := i.bt.Get(userList{cellID: cellID})
			if item != nil {
				result = append(result, item.(userList).list...)
			}
		}
	}

	return result, nil
}
