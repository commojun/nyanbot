package main

import (
	"fmt"

	"github.com/commojun/nyanbot/app/fortune"
)

func main() {
	f := fortune.New()

	fmt.Println(f.DrawByStringSeed("test"))
	fmt.Println(f.DrawByStringSeed("hoge"))
	fmt.Println(f.DrawByStringSeed("fuga"))
	fmt.Println(f.DrawByStringSeed("fuga"))

	fmt.Println(f.Draw(1))
	fmt.Println(f.Draw(2))
	fmt.Println(f.Draw(3))
	fmt.Println(f.Draw(4))
}
