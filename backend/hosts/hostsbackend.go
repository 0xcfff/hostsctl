package hosts

import (
	"fmt"
	"runtime"

	"github.com/0xcfff/dnssync/model"
	"github.com/spf13/afero"
)

const (
	PathHostsFile = "/etc/hosts"
)

type hostsBackend struct {
	etcHostsPath string
	fs           afero.Fs
}

func (backend *hostsBackend) ReadState() (*model.BackendState, error) {

	file, err := backend.fs.Open(backend.etcHostsPath)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %w", err)
	}

	defer file.Close()

	// reader := utils.NewHashCalcReader(file, nil)

	stats, err := file.Stat()

	if err != nil {
		return nil, fmt.Errorf("Error reading stats %w", err)
	}

	fmt.Println("Size: ", stats.Size())

	return &model.BackendState{}, nil
}

func (backend *hostsBackend) UpdateState(changeSet model.BackendStateChangeSet) (*model.BackendState, error) {
	return &model.BackendState{}, nil
}

func DefaultBackend() model.Backend {
	return NewBackend(nil, nil)
}

func NewBackend(hostsFilePath *string, fs afero.Fs) model.Backend {

	backend := hostsBackend{}

	osName := runtime.GOOS

	if osName == "windows" {
		panic("not implemented")

	} else {
		if hostsFilePath == nil || *hostsFilePath == "" {
			backend.etcHostsPath = PathHostsFile
		} else {
			backend.etcHostsPath = *hostsFilePath
		}

		backend.fs = fs
		if backend.fs == nil {
			backend.fs = afero.NewOsFs()
		}
	}

	return &backend
}
