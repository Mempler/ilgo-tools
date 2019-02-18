package gpacker

import (
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"

	"github.com/bnch/uleb128"
)

type Entry struct {
	EntryName      string
	EntryType      entrytype
	EntryData      []byte
	AdditionalData []byte
}

func (e *Entry) WriteToStream(w io.Writer) (err error) {
	_, err = w.Write(uleb128.Marshal(len(e.EntryName)))
	if err != nil {
		return
	}

	_, err = w.Write([]byte(e.EntryName))
	if err != nil {
		return
	}

	_, err = w.Write([]byte{byte(e.EntryType)})
	if err != nil {
		return
	}

	buff := bytes.NewBuffer(nil)
	gzw, err := gzip.NewWriterLevel(buff, 9)
	if err != nil {
		return
	}

	_, err = gzw.Write(e.EntryData)
	if err != nil {
		return
	}

	gzw.Flush()
	gzw.Close()

	_, err = w.Write(uleb128.Marshal(buff.Len()))
	if err != nil {
		return
	}

	_, err = w.Write(buff.Bytes())
	if err != nil {
		return
	}

	buff.Reset()

	gzw, err = gzip.NewWriterLevel(buff, 9)
	if err != nil {
		return
	}

	_, err = gzw.Write(e.AdditionalData)
	if err != nil {
		return
	}

	gzw.Flush()
	gzw.Close()

	_, err = w.Write(uleb128.Marshal(buff.Len()))
	if err != nil {
		return
	}

	_, err = w.Write(buff.Bytes())
	if err != nil {
		return
	}

	return
}

func (e *Entry) ReadFromStream(r io.Reader, version uint64) (err error) {
	length := uleb128.UnmarshalReader(r)

	b := make([]byte, length)
	_, err = r.Read(b)
	if err != nil {
		return
	}

	e.EntryName = string(b)

	b = make([]byte, 1)
	_, err = r.Read(b)
	if err != nil {
		return
	}

	e.EntryType = entrytype(b[0])

	length = uleb128.UnmarshalReader(r)

	b = make([]byte, length)
	_, err = r.Read(b)
	if err != nil {
		return
	}
	buff := bytes.NewReader(b)
	gzr, err := gzip.NewReader(buff)
	if err != nil {
		return
	}

	b, err = ioutil.ReadAll(gzr)
	if err != nil {
		return
	}

	gzr.Close()

	e.EntryData = b

	length = uleb128.UnmarshalReader(r)

	b = make([]byte, length)
	_, err = r.Read(b)
	if err != nil {
		return
	}
	buff = bytes.NewReader(b)
	gzr, err = gzip.NewReader(buff)
	if err != nil {
		return
	}

	b, err = ioutil.ReadAll(gzr)
	if err != nil {
		return
	}

	gzr.Close()

	e.AdditionalData = b

	return
}
