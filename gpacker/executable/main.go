package main

import (
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"net/http"

	"github.com/mempler/ilgo-tools/gpacker"
)

var compress bool
var outFile string
var outDirectory string

type Lfiles []string

func (i *Lfiles) String() string {
	return ""
}

func (i *Lfiles) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var files Lfiles

func init() {
	flag.BoolVar(&compress, "compress", true, "If true Compress, if false Decompress")
	flag.StringVar(&outFile, "outfile", "out.gpack", "Output file")
	flag.StringVar(&outDirectory, "outfolder", "out", "Output file")
	flag.Var(&files, "files", "Files to Compress / File to Decompress")

	flag.Parse()
}

func main() {
	gp := gpacker.MakeGPackage()

	if !compress && len(files) > 1 {
		panic("Cant Decompress more then 1 file")
	}

	if len(files) < 1 {
		flag.PrintDefaults()
	}

	for _, fn := range files {
		if compress {
			f, err := os.Open(fn)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			// Get the content
			contentType, err := GetFileContentType(f)
			if err != nil {
				panic(err)
			}

			var x = gpacker.TBinary
			if strings.HasPrefix(contentType, "text") {
				x = gpacker.TText
			}

			if strings.HasPrefix(contentType, "image") {
				x = gpacker.TImage
			}

			if strings.HasPrefix(contentType, "font") {
				x = gpacker.TFont
			}

			out, err := ioutil.ReadAll(f)
			if err != nil {
				panic(err)
			}

			gp.AppendNewEntry(filepath.Base(fn), x, out, nil)
		} else {
			err := gp.ReadFromFile(files[0])
			if err != nil {
				panic(err)
			}

			os.MkdirAll(outDirectory, 0777)

			for _, entry := range gp.GetAllEntries() {
				ioutil.WriteFile(outDirectory+"/"+entry.EntryName, entry.EntryData, 0777)
				ioutil.WriteFile(outDirectory+"/"+entry.EntryName+"_addition", entry.AdditionalData, 0777)
			}
		}
	}

	if !compress {
		gp.WriteToFile(outFile)
	}
}

func GetFileContentType(out *os.File) (string, error) {
	buffer := make([]byte, 512)

	_, err := out.Read(buffer)
	if err != nil {
		return "", err
	}

	contentType := http.DetectContentType(buffer)

	out.Seek(0, 0)

	return contentType, nil
}
