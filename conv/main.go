package conv

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/gen2brain/dlgs"
	"github.com/gidor/tigre2shp/feature"
	"github.com/gidor/tigre2shp/tigre"
	"github.com/jonas-p/go-shp"
)

func init() {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// wr, _ := os.Open(filepath.Join(dir, "tigre2sp.log"))
	wr, _ := os.OpenFile(filepath.Join(dir, "tigre2sp.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 066)
	log.SetOutput(wr)
	log.SetPrefix("tigre2shp")
}

func convertvalue(fd feature.FieldDescription, bs string) (setval interface{}) {
	switch fd.Type {
	case feature.Character:
		setval = bs
	case feature.Numeric:
		if fd.Decimals == 0 {
			val, _ := strconv.ParseInt(bs, 10, 32)
			setval = int(val)
		} else {
			val, _ := strconv.ParseFloat(bs, 64)
			setval = float64(val)
		}
	case feature.Date:
		setval = bs
	}
	return
}

func buildshape(gt int, item tigre.Oggetto) (geom shp.Shape) {
	switch gt {
	case feature.Point:
		g := new(shp.Point)
		g.X = item.Punti[0].X
		g.Y = item.Punti[0].Y
		geom = g
	case feature.Polygon:
	case feature.Line:
		ps := make([]shp.Point, 0)
		for _, p := range item.Punti {
			ps = append(ps, shp.Point{X: p.X, Y: p.Y})
		}
		g := shp.NewPolyLine([][]shp.Point{ps})
		geom = g
	default:
		g := new(shp.Point)
		g.X = item.Punti[0].X
		g.Y = item.Punti[0].Y
		geom = g
	}
	return
}

func renderAttribute(val string, atts map[uint32]string) string {
	r, _ := regexp.Compile("<([0-9]+)>")
	ss := r.FindAllString(val, 22)
	if len(ss) > 0 {
		// se ho trovato match
		for i := 0; i < len(ss); i++ {
			seq, _ := strconv.Atoi(strings.Trim(ss[i], "<>"))
			part, ok := atts[uint32(seq)]
			if !ok {
				part = ""
			}
			val = strings.ReplaceAll(val, ss[i], part)
		}
	}
	return val
}

func process(outdataset *DataSet, oggetti []tigre.Oggetto, mappings map[uint32]Mapping) {
	defaults := mappings[0].Atts
	for _, item := range oggetti {
		mapping, ok := mappings[item.TipoOggetto]
		if !ok {
			fmt.Printf("no mapping for item %d\n", item.TipoOggetto)
			log.Printf("no mapping for item %d\n", item.TipoOggetto)
			continue
		}
		shpfile, ok := outdataset.GetShp(mapping.Out)
		if !ok {
			log.Printf("no shapefile for item %d\n", mapping.Out)
			fmt.Errorf("no shapefile for item %d", mapping.Out)
			continue
		}
		// write geometry
		index := shpfile.handler.Write(buildshape(shpfile.descriptor.GeometryType, item))
		// write data
		for i, f := range shpfile.descriptor.Attribute {
			val, ok := mapping.Atts[f.FName]
			if !ok {
				val, ok = defaults[f.FName]
			}
			if ok {
				// ho mapping per l'attributo
				val = renderAttribute(val, item.Attributi)
			}
			if len(val) > 0 {
				shpfile.handler.WriteAttribute(int(index), i, convertvalue(f, val))
				if f.FName == "DATA_POSA" {
					fmt.Println(val)
				}
			} else {
				// non ho mapping  verifico default
				if len(f.Default) > 0 {
					shpfile.handler.WriteAttribute(int(index), i, convertvalue(f, f.Default))
					if f.FName == "DATA_POSA" {
						fmt.Println(convertvalue(f, f.Default))
					}
				}
			}
		}
	}

}

func selectDir(msg string) (dir string, err error) {
	ok := false
	for i := 0; i < 3; i++ {
		println(ok)
		println(i)
		dir, ok, err = dlgs.File(msg, "", true)
		if err != nil {
			panic(err)
		}
		if ok {
			return
		}
	}
	if !ok {
		println("3 tenttivi")
		os.Exit(1)
	}
	return
}

func cancellamiglob(base string) ([]string, error) {
	var files []string
	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		fmt.Println(path)
		if path == base {
			return nil
		} else if !info.IsDir() {
			files = append(files, path)
		} else {
			return filepath.SkipDir
		}
		return nil
	})
	return files, err
}

func Main() {
	// conf, _ := config.Get()
	// fmt.Println(conf)
	var (
		dirMeta, dirShp string
		err             error
	)
	defs := feature.Load()
	mappingss := LoadMapings()
	// fmt.Println(defs)
	if len(os.Args) == 1 {
		dirMeta, err = selectDir("Seleziona directory Metafile Tigre")
		if err != nil {
			println(err)
			os.Exit(1)
		}
		dirShp, err = selectDir("Seleziona directory Output")
		if err != nil {
			println(err)
			os.Exit(1)
		}
	} else if len(os.Args) == 2 {
		println("due dir ")
		os.Exit(1)
	} else {
		dirMeta = os.Args[1]
		dirShp = os.Args[2]
	}
	log.Printf("===========================================\n")
	log.Printf("Conversione tigre 2 shapefile\n Meta in %s \n shapefile out %s\n", dirMeta, dirShp)
	outdataset := NewDataset(dirShp)
	outdataset.SetFeatures(defs)
	dataset := tigre.NewDataset(dirMeta)
	ogg := dataset.Get()
	process(outdataset, ogg, mappingss)

	// fmt.Println(ogg)
	// fmt.Println(dirShp)
	outdataset.closeShp()
	// tigre.Test()
	// shp.Open(dirShp)
	fmt.Println("Finito")
}
