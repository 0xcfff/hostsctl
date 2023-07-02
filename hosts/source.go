package hosts

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/0xcfff/hostsctl/hosts/dom"
	"github.com/spf13/afero"
)

// Holds information about hosts mapping config file location
type Source struct {
	etcHostsPath string
	fs           afero.Fs
}

var (
	// Default OS-level hosts mapping configuration file
	EtcHosts *Source = NewSource("", nil)
)

func NewSource(hostsFilePath string, fs afero.Fs) *Source {
	fileSource := Source{}

	if hostsFilePath == "" {
		fileSource.etcHostsPath = defaultEtcHostsPath()
	} else {
		fileSource.etcHostsPath = hostsFilePath
	}

	fileSource.fs = fs
	if fileSource.fs == nil {
		fileSource.fs = afero.NewOsFs()
	}

	return &fileSource
}

func (src *Source) Path() string {
	return src.etcHostsPath
}

func (src *Source) openRead() (afero.File, error) {
	return src.fs.Open(src.etcHostsPath)
}

func (src *Source) openWrite() (afero.File, error) {
	// return src.fs.OpenFile(src.etcHostsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
	return src.fs.OpenFile(src.etcHostsPath, os.O_CREATE|os.O_RDWR, 0o644)
}

func (src *Source) Load() (*dom.Document, error) {
	f, err := src.openRead()
	if err != nil {
		return nil, fmt.Errorf("can't open hosts file %s, %w", src.Path(), err)
	}

	defer f.Close()

	doc, err := dom.Read(f)
	if err != nil {
		return nil, fmt.Errorf("can't parse hosts file %s, %w", src.Path(), err)
	}

	return doc, nil
}

func (src *Source) Save(doc *dom.Document, fm dom.FmtMode) error {
	f, err := src.openWrite()
	if err != nil {
		return fmt.Errorf("can't open hosts file %s, %w", src.Path(), err)
	}

	defer f.Close()

	err = dom.Write(f, doc, fm)
	if err != nil {
		return fmt.Errorf("can't parse hosts file %s, %w", src.Path(), err)
	}

	position, err := f.Seek(0, io.SeekCurrent)
	if err != nil {
		return fmt.Errorf("can't read position %s, %w", src.Path(), err)
	}

	err = f.Truncate(position)
	if err != nil {
		return fmt.Errorf("can't truncate %s, %w", src.Path(), err)
	}

	return nil
}

func defaultEtcHostsPath() string {
	result := "/etc/hosts"
	if runtime.GOOS == "windows" {
		result = fmt.Sprintf("%s\\Drivers\\etc\\hosts", os.Getenv("SYSTEM32"))
	}
	return result
}
