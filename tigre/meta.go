package tigre

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type ScannEOF error

/*
	Test
	solo test
*/
func Test() {
	fmt.Println("fatto")
	fmt.Println("fatto")
	fmt.Println("fatto")
	fmt.Println("fatto")
}

//  MetaFile  descriptio of a single metafile
type MetaFile struct {
	PathOggetti string //  path to Object file
	PathPunti   string //  path to Point file
	PathDati    string //  path to Data file
}

type loadingMeta struct {
	meta  MetaFile
	dati  metaDati
	punti metaPunti
}

func (m *MetaFile) Load(base string) (oggetti []Oggetto) {
	work := new(loadingMeta)
	work.meta = *m
	work.readDati(base)
	work.readPunti(base)
	oggetti = work.readOggetti(base)
	return
}

//  MetaFile  descriptio of a single metafile
type Dataset struct {
	Path      string     //  path to directory or zip containig etafiles
	MetaFiles []MetaFile //  path to Point file
	data      string     //  path to Data file
	index     int
}

func (m *Dataset) scandir() {
	pattern := path.Join(m.Path, "o*.RTE")
	// files, err := glob(dirMeta)
	files, err := filepath.Glob(pattern)
	if err != nil {
	}
	for i := 0; i < len(files); i++ {
		_, file := filepath.Split(files[i])
		file = file[1:]
		meta := MetaFile{PathOggetti: "o" + file, PathPunti: "p" + file, PathDati: "d" + file}
		m.MetaFiles = append(m.MetaFiles, meta)
	}
	return
}
func (m *Dataset) Get() []Oggetto {
	var result = make([]Oggetto, 0, 50)
	for i := 0; i < len(m.MetaFiles); i++ {
		w := loadingMeta{meta: m.MetaFiles[i]}
		w.readDati(m.Path)
		w.readPunti(m.Path)
		ogg := w.readOggetti(m.Path)
		result = append(result, ogg...)
	}
	return result
}

func NewDataset(base string) (ds *Dataset) {
	ds = new(Dataset)
	ds.Path = base
	ds.MetaFiles = make([]MetaFile, 0, 30)
	ds.scandir()
	return
}

// Oggetto dati relativi ad una feature
type Oggetto struct {
	Mappa       uint32
	IdOggetto   uint32
	TipoOggetto uint32
	Punti       []Punto
	Attributi   Dati
}

/*
Oggetti
oggetto      | tipo | posizione | lunghezza |valori |
mappa        |intero|   1       |4          |0001 9999
idOggetto    |intero|   5       |4          |0001 9999
tipoOggetto  |intero|   9       |4          |0001 9999
seqPt        |intero|   13      |3          |001 999
idPunto      |intero|   16      |5          |00001 90999
tipoPt       |intero|   21      |1          | G W Y
codSecondary |intero|   22      |1          | 0
rettiline    |intero|   23      |1          |
*/

// readDati read Punti
func (m *loadingMeta) readOggetti(base string) (oggetti []Oggetto) {
	oggetti = make([]Oggetto, 0, 10)
	f, err := os.Open(path.Join(base, m.meta.PathOggetti))

	if err != nil {
		log.Fatal(err)
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	var (
		currogg *Oggetto
	)
	// Scan for next line.
	for scanner.Scan() {
		var (
			mappa, idOggetto, tipoOggetto, seqPt, idPunto, codSecondary, rettiline uint32
			tipoPt                                                                 string
		)
		buffer := scanner.Text()
		n, err := fmt.Sscanf(buffer, "%4d%4d%4d%3d%5d%1s%1d%1d", &mappa, &idOggetto, &tipoOggetto, &seqPt, &idPunto, &tipoPt, &codSecondary, &rettiline)
		if n < 8 {
			if err != nil {
				log.Println(buffer)
				log.Fatal(err)
			}
		}
		if currogg == nil {
			currogg = new(Oggetto)
			currogg.Mappa = mappa
			currogg.IdOggetto = idOggetto
			currogg.TipoOggetto = tipoOggetto
			currogg.Punti = make([]Punto, 0, 2)
			currogg.Punti = append(currogg.Punti, m.getPunti(idPunto))
			// att, _ := m.getDati(idOggetto)
			currogg.Attributi, _ = m.getDati(idOggetto)
		} else if currogg.IdOggetto == idOggetto {
			if tipoPt == "G" || tipoPt == "W" {
				currogg.Punti = append(currogg.Punti, m.getPunti(idPunto))
			}
		} else {
			oggetti = append(oggetti, *currogg)
			currogg = new(Oggetto)
			currogg.Mappa = mappa
			currogg.IdOggetto = idOggetto
			currogg.TipoOggetto = tipoOggetto
			currogg.Punti = make([]Punto, 0, 2)
			currogg.Punti = append(currogg.Punti, m.getPunti(idPunto))
			currogg.Attributi, _ = m.getDati(idOggetto)
		}
	}
	return
}

type Dati map[uint32]string
type metaDati map[uint32]Dati

/*
Dati
oggetto      | tipo | posizione | lunghezza |valori |
mappa        |intero|   1       |4          |0001 9999
idOggetto    |intero|   5       |4          |0001 9999
seqAttributo |intero|   9       |3          |001 999
dato         |alfanu|  12       |55         |
*/

// readDati read Punti
func (m *loadingMeta) readDati(base string) {
	f, err := os.Open(path.Join(base, m.meta.PathDati))
	// f, err := os.Open(m.meta.PathDati)
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	defer f.Close()
	m.dati = make(metaDati)
	// Scan for next line.
	for scanner.Scan() {
		var (
			mappa, idOggetto, seqAttributo uint32
			dato                           string
		)
		buffer := scanner.Text()
		n, err := fmt.Sscanf(buffer, "%4d%4d%3d%s", &mappa, &idOggetto, &seqAttributo, &dato)
		if n < 4 {
			if err != nil {
				log.Println(buffer)
				log.Fatal(err)
			}
		}
		elem, ok := m.dati[idOggetto]
		if !ok {
			elem = make(Dati)
			m.dati[idOggetto] = elem
		}
		elem[seqAttributo] = strings.Trim(dato, " ")
	}
	return

}

func (m *loadingMeta) getDati(idOggetto uint32) (elem Dati, ok bool) {
	elem, ok = m.dati[idOggetto]
	return
}

type metaPunti map[uint32]Punto

// mappa    uint32  //|intero|   1       |4          |0001 9999
// idPunto  uint32  //|intero|   5       |5          |00001 99999
type Punto struct {
	X float64 //|reale |   10      |17         |001 999
	Y float64 //|reale |   27      |17         |001 999
}

/*
Punti
oggetto      | tipo | posizione | lunghezza |valori |
mappa        |intero|   1       |4          |0001 9999
idPunto      |intero|   5       |5          |00001 99999
ascissa	     |reale	|   10      |17         |001 999
ordinata     |reale	|   27      |17         |001 999
*/

// read Punti
func (m *loadingMeta) getPunti(idPunto uint32) Punto {
	return m.punti[idPunto]
}

func (m *loadingMeta) readPunti(base string) {
	f, err := os.Open(path.Join(base, m.meta.PathPunti))
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)
	defer f.Close()
	m.punti = make(metaPunti)
	// Scan for next line.
	for scanner.Scan() {
		var (
			mappa, idPunto uint32
			x, y           float64
		)
		buffer := scanner.Text()
		n, err := fmt.Sscanf(buffer, "%4d%5d%17f%17f", &mappa, &idPunto, &x, &y)
		if n < 4 {
			if err != nil {
				log.Println(buffer)
				log.Fatal(err)
			}
		}
		m.punti[idPunto] = Punto{X: x, Y: y}
	}
	return
}
