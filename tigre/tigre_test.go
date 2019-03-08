package tigre_test

import (
	"fmt"
	"testing"
)

func Test_a_test(t *testing.T) {
	var lines [15]string
	lines[0] = "20036000700100002                                                  "
	lines[1] = "20036000700200125                                                  "
	lines[2] = "2003600070030003.50                                                "
	lines[3] = "20036000800100002                                                  "
	lines[4] = "20036000800200125                                                  "
	lines[5] = "2003600080030005.97                                                "
	lines[6] = "20036000900100002                                                  "
	lines[7] = "20036000900200125                                                  "
	lines[8] = "2003600090030002.63                                                "
	lines[9] = "20036001000100002                                                  "
	lines[10] = "2003600100032009                                                   "
	lines[11] = "20036001000408                                                     "
	lines[12] = "20036001000530                                                     "
	lines[13] = "200360010006Z                                                      "
	lines[14] = "200360011004A                                                      "
	var mappa, idOggetto, seqAttributo uint32
	var dato string

	for i, s := range lines {
		fmt.Println(i, s)
		t.Log(i, s)
		n, _ := fmt.Sscanf(s, "%4d%4d%3d%s", &mappa, &idOggetto, &seqAttributo, &dato)
		fmt.Println(i, n, mappa, idOggetto, seqAttributo, dato)
	}
}
