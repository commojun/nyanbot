package masterdata

import (
	"github.com/commojun/nyanbot/masterdata/key_value"
	"github.com/commojun/nyanbot/masterdata/table"
)

func Initialize() error {
	_, err := table.Initialize()
	if err != nil {
		return err
	}

	_, err = key_value.Initialize()
	if err != nil {
		return err
	}

	return nil
}
