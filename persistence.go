package similar

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	log "github.com/sirupsen/logrus"
)

type CSVPersistence struct {
	headers       []string
	persistWriter *csv.Writer
	persistFile   *os.File
	filepath      string
	logger        *log.Entry
}

func NewCSVPersistence(path string, headers ...string) (p *CSVPersistence, err error) {

	_, err = os.Stat(path)

	p = &CSVPersistence{headers: headers, filepath: path, logger: log.WithFields(log.Fields{
		"module": "persistence",
		"file":   path,
	})}

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {

			p.logger.Debug("persistence file", path, "is not existed, create it")

			p.persistFile, err = os.OpenFile(
				path,
				os.O_CREATE|os.O_WRONLY|os.O_TRUNC,
				0600,
			)

			if err != nil {
				return nil, err
			}

			p.persistWriter = csv.NewWriter(p.persistFile)

			if err = p.persistWriter.Write(p.headers); err != nil {
				return nil, err
			}

		} else {
			return nil, err
		}

	} else {

		p.logger.Debug("persistence file", path, "is existed")

		p.persistFile, err = os.OpenFile(
			path,
			os.O_CREATE|os.O_WRONLY|os.O_APPEND,
			0600,
		)

		if err != nil {
			return nil, err
		}

		p.persistWriter = csv.NewWriter(p.persistFile)

	}
	return
}

func (p *CSVPersistence) Restore(onRestore func(cells ...string)) error {

	readOnlyFile, err := os.OpenFile(
		p.filepath,
		os.O_RDONLY, // READ ONLY
		0600,
	)

	if err != nil {
		return err
	}

	defer readOnlyFile.Close()

	reader := csv.NewReader(readOnlyFile)
	headers, err := reader.Read()

	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	if !reflect.DeepEqual(headers, p.headers) {
		err = fmt.Errorf("the header of target file is not correct: %v, required: %v", headers, p.headers)
		return err
	}

	for {

		record, err := reader.Read()

		if err != nil {
			if err == io.EOF {
				// file end
				break
			} else {
				return err
			}
		}

		onRestore(record...)

	}
	return nil
}

func (p *CSVPersistence) Append(cells ...string) error {
	if len(cells) != len(p.headers) {
		return fmt.Errorf("total count of cells is not match the headers")
	}
	return p.persistWriter.Write(cells)
}

func (p *CSVPersistence) Flush() {
	p.persistWriter.Flush()
}

func (p *CSVPersistence) Close() error {
	p.persistWriter.Flush()
	return p.persistFile.Close()
}
