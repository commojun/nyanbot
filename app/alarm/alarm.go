package alarm

import (
	"io/ioutil"
	"time"

	"github.com/commojun/nyanbot/app/constant"
	"github.com/commojun/nyanbot/app/linebot"
	"github.com/jszwec/csvutil"
)

type Alarm struct {
	Entities *[]Entity
}

func Load() (*Alarm, error) {
	path := constant.AlarmCsvPath
	return LoadWithPath(path)
}

func LoadWithPath(path string) (*Alarm, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return &Alarm{}, err
	}

	var entities []Entity
	err = csvutil.Unmarshal(buf, &entities)
	if err != nil {
		return &Alarm{}, err
	}

	alm := Alarm{
		Entities: &entities,
	}
	return &alm, nil
}

func (alm *Alarm) Send() error {
	bot, err := linebot.New()
	if err != nil {
		return err
	}

	for _, entity := range *alm.Entities {
		chk, err := entity.Check(time.Now())
		if err != nil {
			return err
		}
		if chk == false {
			continue
		}
		bot.TextMessage(entity.Message)
		if err != nil {
			return err
		}
	}

	return nil

}
