package ojisan

import (
	"github.com/greymd/ojichat/generator"
)

type Ojisan struct {
	Config *generator.Config
}

func New(name string, emojiNum int, level int) Ojisan {
	config := generator.Config{
		TargetName:        name,
		EmojiNum:          emojiNum,
		PunctiuationLevel: level,
	}

	return Ojisan{
		Config: &config,
	}
}

func (ojisan *Ojisan) Generate() (string, error) {
	return generator.Start(*ojisan.Config)
}
