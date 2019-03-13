package conv

import (
	"path/filepath"

	"github.com/gidor/tigre2shp/feature"
	"github.com/jonas-p/go-shp"
)

type shpFile struct {
	descriptor feature.FeatureDescription
	filename   string
	handler    *shp.Writer
}

type DataSet struct {
	Path  string
	Files []shpFile
}

// NewDataset
func NewDataset(path string) *DataSet {
	out := new(DataSet)
	out.Path = path
	out.Files = make([]shpFile, 0)
	return out
}

// SetFeatures
func (d *DataSet) SetFeatures(features []feature.FeatureDescription) {
	for _, descr := range features {
		d.addShp(descr)
	}
}

func (d *DataSet) addShp(desc feature.FeatureDescription) {
	var (
		gt      shp.ShapeType
		shpItem shpFile
	)

	switch desc.GeometryType {
	case feature.Point:
		gt = shp.POINT
	case feature.Polygon:
		gt = shp.POLYGON
	case feature.Line:
		gt = shp.POLYLINE
	default:
		gt = shp.POINT
	}
	shpItem.descriptor = desc
	shpItem.filename = filepath.Join(d.Path, desc.Name+".shp")
	shpItem.handler, _ = shp.Create(shpItem.filename, gt)
	// fields to write
	fields := make([]shp.Field, 0)
	for _, att := range desc.Attribute {
		switch att.Type {
		case feature.Character:
			fields = append(fields, shp.StringField(att.FName, uint8(att.Len)))
		case feature.Numeric:
			if att.Decimals == 0 {
				fields = append(fields, shp.NumberField(att.FName, uint8(att.Len)))
			} else {
				fields = append(fields, shp.FloatField(att.FName, uint8(att.Len), uint8(att.Decimals)))
			}
		case feature.Date:
			fields = append(fields, shp.DateField(att.FName))
		case feature.Logical:
			fields = append(fields, shp.NumberField(att.FName, 1))
		default:
			fields = append(fields, shp.StringField(att.FName, 20))
		}
	}
	shpItem.handler.SetFields(fields)
	d.Files = append(d.Files, shpItem)
}

func (d *DataSet) GetShp(id uint32) (*shpFile, bool) {
	i := int(id)
	for _, item := range d.Files {
		if item.descriptor.Code == i {
			return &item, true
		}
	}
	return nil, false
}

func (d *DataSet) closeShp() {
	for _, item := range d.Files {
		item.handler.Close()
	}
	d.Files = d.Files[0:0]
}
