package gpacker

import (
	"os"
	"testing"
)

const outfile = "0.gpack"

func TestWriteGPackage(t *testing.T) {
	var err error
	var pg *GPackage

	pg = MakeGPackage()

	if pg == nil {
		t.Error("PG is NIL!")
	}

	pg.AppendNewEntry("T", TText, []byte("Test Text"), nil)

	if pg.GetAllEntries() == nil || len(pg.GetAllEntries()) < 1 {
		t.Fail()
	}

	err = pg.WriteToFile(outfile)
	if err != nil {
		t.Error(err)
	}
}

func TestReadGPackage(t *testing.T) {
	var err error
	var pg *GPackage

	pg = MakeGPackage()

	err = pg.ReadFromFile(outfile)
	if err != nil {
		t.Error(err)
	}

	if pg.GetAllEntries() == nil || len(pg.GetAllEntries()) < 1 {
		t.Fail()
	}

	err = os.Remove(outfile)
	if err != nil {
		t.Error(err)
	}
}
