package gpacker

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"sync"
)

// CurrentVersion is the Current used Version ofc.
const CurrentVersion uint64 = 0x7323

// UniqueDataType is a semi MimeType to check if our .gpack file is actually an GPACKAGE
const UniqueDataType uint64 = 0x4547414B43415047

type GPackage struct {
	_Mut     sync.Mutex
	_Entries []Entry
}

func MakeGPackage() *GPackage {
	return &GPackage{}
}

func (gp *GPackage) WriteToFile(filename string) (err error) {
	buff := bytes.NewBuffer(nil)

	err = binary.Write(buff, binary.LittleEndian, CurrentVersion)
	if err != nil {
		return
	}
	err = binary.Write(buff, binary.LittleEndian, UniqueDataType)
	if err != nil {
		return
	}
	err = binary.Write(buff, binary.LittleEndian, int64(len(gp._Entries)))
	if err != nil {
		return
	}

	for _, entry := range gp._Entries {
		err = entry.WriteToStream(buff)
		if err != nil {
			return
		}
	}

	return ioutil.WriteFile(filename, buff.Bytes(), 0644)
}

func (gp *GPackage) ReadFromFile(filename string) (err error) {

	f, err := os.OpenFile(filename, os.O_RDONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()

	vers := uint64(0)
	binary.Read(f, binary.LittleEndian, &vers)

	unique := uint64(0)
	binary.Read(f, binary.LittleEndian, &unique)
	if unique != UniqueDataType {
		return errors.New("File is not an GPackage")
	}

	entryLength := uint64(0)
	binary.Read(f, binary.LittleEndian, &entryLength)

	for i := uint64(0); i < entryLength; i++ {
		e := Entry{}

		err = e.ReadFromStream(f, vers)
		if err != nil {
			return
		}

		gp.AppendEntry(e)
	}

	return
}

func (gp *GPackage) AppendEntry(entry Entry) {
	gp._Mut.Lock()
	gp._Entries = append(gp._Entries, entry)
	gp._Mut.Unlock()
}

func (gp *GPackage) AppendNewEntry(name string, etype entrytype, data []byte, additionaldata []byte) {
	gp._Mut.Lock()
	gp._Entries = append(gp._Entries, Entry{name, etype, data, additionaldata})
	gp._Mut.Unlock()
}

func (gp *GPackage) GetAllEntries() []Entry {
	return gp._Entries
}
