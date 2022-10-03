package hosts

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/spf13/afero"
)

type HostsFileSource struct {
	etcHostsPath string
	fs           afero.Fs
}

func NewHostsFileSource(hostsFilePath string, fs afero.Fs) *HostsFileSource {
	fileSource := HostsFileSource{}

	if hostsFilePath == "" {
		fileSource.etcHostsPath = EtcHostsPath()
	} else {
		fileSource.etcHostsPath = hostsFilePath
	}

	fileSource.fs = fs
	if fileSource.fs == nil {
		fileSource.fs = afero.NewOsFs()
	}

	return &fileSource
}

func (s *HostsFileSource) openRead() (afero.File, error) {
	return s.fs.Open(s.etcHostsPath)
}

func (s *HostsFileSource) openWrite() (afero.File, error) {
	return s.fs.OpenFile(s.etcHostsPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0o644)
}

func (s *HostsFileSource) LoadFile() (*HostsFile, error) {
	f, err := s.fs.Open(s.etcHostsPath)
	if err != nil {
		return nil, fmt.Errorf("Can't open hosts file %s, %w", s.etcHostsPath, err)
	}

	defer f.Close()

	buff, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("Can't read hosts file %s, %w", s.etcHostsPath, err)
	}

	hostsFile := NewHostsFile(buff)
	return hostsFile, nil
}

func (f *HostsFileSource) LoadAndParse(mode ParseMode, flags ParseFlags) (*HostsFileContent, error) {
	r, err := f.openRead()
	if err != nil {
		return nil, err
	}

	defer r.Close()

	return ParseHostsFile(r, mode, flags)
}

func EtcHostsPath() string {
	result := "/etc/hosts"
	if runtime.GOOS == "windows" {
		result = fmt.Sprintf("%s\\Drivers\\etc\\hosts", os.Getenv("SYSTEM32"))
	}
	return result
}
