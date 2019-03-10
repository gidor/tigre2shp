package feature

import (
	"github.com/twpayne/go-geom"
)

// Feature attribte data type
const (
	Character = 1 // Character data type
	Numeric   = 2 // Numeric data type
	Date      = 3 // Date data type
	Logical   = 4 // Logical data type
	Memo      = 5 // Memo data type
)

// GeometryTypes
const (
	Point        = 1
	Line         = 2
	Polygon      = 3
	MultiPoint   = 4
	MultiLine    = 5
	MultiPolygon = 6
)

// Dimension
const (
	xy   = 2
	xyz  = 3
	xyzm = 4
)

// FieldDescription Desribe a field
type FieldDescription struct {
	FName    string
	Type     int
	Len      int
	Decimals int
	Domain   int
	Defalt   string
	Nullable bool
}

// Description of a Feature
type Description struct {
	Name string
}

// AttrbuteDescription a map of FieldDescription
// type AttrbuteDescription map[string]FieldDescription
type AttrbuteDescription []FieldDescription

type AttributeMap map[string]interface{}

type FeatureDescription struct {
	Code         int
	Name         string
	GeometryType int
	Attribute    []FieldDescription
}

// FeatureImpl T base implementation tehe data rappraetning the Feature
type FeatureImpl struct {
	geometry geom.T
	atts     AttributeMap
	ftype    string
}

// SetFeatureType setter
func (m *FeatureImpl) SetFeatureType(value string) {
	m.ftype = value
	return
}

// Geometry getter
func (m *FeatureImpl) FeatureType() string {
	return m.ftype
}

// SetGeometry setter
func (m *FeatureImpl) SetGeometry(val geom.T) {
	m.geometry = val
	return
}

// Geometry getter
func (m *FeatureImpl) Geometry() geom.T {
	return m.geometry
}

// Attribute getter
func (m *FeatureImpl) Attribute(key string) interface{} {
	if m.atts == nil {
		return nil
	}
	return m.atts[key]
}

// SetAttribute getter
func (m *FeatureImpl) SetAttribute(key string, val interface{}) {
	if m.atts == nil {
		m.atts = make(AttributeMap)
	}
	m.atts[key] = val
	return
}

// Attributes getter
func (m *FeatureImpl) Attributes() AttributeMap {
	return m.atts
}

// Srid getter
func (m *FeatureImpl) Srid() string {
	return ""
}

type T interface {
	Geometry() geom.T
	Attributes() AttributeMap
	FeatureType() string
	Srid() string
}
