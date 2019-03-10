package feature

import (
	"github.com/gidor/tigre2shp/config"
)

func Load() []FeatureDescription {
	var out = make([]FeatureDescription, 0)
	cfg, ok := config.Get()
	features, ok := config.GetArray(cfg, "Features")
	if !ok {
		return out
	}
	for _, value := range features {
		var item FeatureDescription
		item.Code, _ = config.GetInt(value, "FCode")
		item.Name, _ = config.GetString(value, "Nome")

		geom, _ := config.GetString(value, "Geometry")
		switch geom {
		case "Line":
			item.GeometryType = Line
		case "Point":
			item.GeometryType = Point
		case "Polygon":
			item.GeometryType = Polygon
		default:
			item.GeometryType = Point
		}
		// item.Attribute = make(AttrbuteDescription)
		item.Attribute = make([]FieldDescription, 0)
		atts, _ := config.GetArray(value, "Attributes")
		for _, att := range atts {
			format, _ := config.GetString(att, "Format")
			at_len, _ := config.GetInt(att, "Length")

			iat := new(FieldDescription)
			iat.FName, _ = config.GetString(att, "Name")
			iat.Domain, _ = config.GetInt(att, "Domain")
			iat.Defalt, _ = config.GetString(att, "Default")
			at_null, _ := config.GetInt(att, "NotNull")
			iat.Nullable = (at_null == 0)
			switch format {
			case "INTERO":
				iat.Type = Numeric
				iat.Len = 9
				iat.Decimals = 0
			case "TESTO":
				iat.Type = Character
				iat.Len = at_len
				iat.Decimals = 0
			case "DECIMALE":
				iat.Type = Numeric
				iat.Len = 17
				iat.Decimals = 8
			case "DATA":
				iat.Type = Date
				iat.Len = 8
				iat.Decimals = 0
			default:
				iat.Type = Character
				iat.Len = 20
				iat.Decimals = 0
			}
			item.Attribute = append(item.Attribute, *iat)
		}
		out = append(out, item)
	}
	return out

}
