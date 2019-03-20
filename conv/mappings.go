package conv

import (
	"fmt"

	"github.com/gidor/tigre2shp/config"
)

type Mapping struct {
	In   uint32
	Out  uint32
	Atts map[string]string
}

// LoadMapings carica mapping da sqlcfg
func LoadMapings() map[uint32]Mapping {
	out := make(map[uint32]Mapping)
	mapping := config.FMappings()
	for _, m := range mapping {
		mm := new(Mapping)
		mm.In = uint32(m.FcIn)
		mm.Out = uint32(m.FcOut)
		mm.Atts = make(map[string]string)
		// atts := config.AttMappings(mm.In, mm.Out)
		// for _, k := range atts {
		for _, k := range config.AttMappings(mm.In, mm.Out) {
			mm.Atts[k.Attributo] = k.Valore
		}
		out[mm.In] = *mm
	}
	// store default mapping as code 0
	mm := new(Mapping)
	mm.Atts = make(map[string]string)
	defaults := config.DefaultValues()

	for _, k := range defaults {
		mm.Atts[k.Attributo] = k.Valore
	}
	out[0] = *mm

	return out
}

// noLoadMapings deprecata carica da json cfg
func noLoadMapings() map[uint32]Mapping {
	cfg, ok := config.Get()
	if !ok {
		fmt.Errorf("Errore leggendo configuration")
	}
	out := make(map[uint32]Mapping)
	mappings, ok := config.GetArray(cfg, "Mappings")
	if !ok {
		fmt.Errorf("Errore leggendo Mappings")
	}
	for _, m := range mappings {
		mm := new(Mapping)
		i, _ := config.GetInt(m, "In")
		mm.In = uint32(i)
		i, _ = config.GetInt(m, "Out")
		mm.Out = uint32(i)
		mm.Atts = make(map[string]string)
		atts, _ := config.GetMap(m, "Atts")
		for k := range atts {
			v, _ := config.GetString(atts, k)
			mm.Atts[k] = v
		}
		out[mm.In] = *mm
	}
	// store default mapping as code 0
	mm := new(Mapping)
	mm.Atts = make(map[string]string)
	defaults, _ := config.Defaults("defaults")

	atts, _ := config.GetMap(defaults, "Atts")
	for k := range atts {
		v, _ := config.GetString(atts, k)
		mm.Atts[k] = v
	}
	out[0] = *mm

	return out
}
