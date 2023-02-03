package main

import (
	"github.com/maldan/go-cdb/cdb_proto/harray"
	"github.com/maldan/go-cmhp/cmhp_print"
)

type Test struct {
	FirstName string `json:"firstName" id:"0"`
	LastName  string `json:"lastName" id:"1" len:"32"`
	Phone     string `json:"phone" id:"2" len:"64"`
	Sex       string `json:"sex" id:"3" len:"64"`
	Rock      string `json:"rock" id:"4" len:"64"`
	Gas       string `json:"gas" id:"5" len:"64"`
	Yas       string `json:"yas" id:"6" len:"64"`
	Taj       string `json:"taj" id:"7" len:"64"`
	Mahal     string `json:"mahal" id:"8" len:"64"`
	Ebal      string `json:"ebal" id:"9" len:"64"`
	Sasal     string `json:"sasal" id:"10" len:"64"`
	Sasal2    string `json:"sasal2" id:"11" len:"64"`
}

func main() {
	harr := harray.New()
	harr.Add("a", 1)
	cmhp_print.PrintBytesColored(harr.Memory, 32, []cmhp_print.ColorRange{
		{0, 2, cmhp_print.BgRed},
		{2, 4, cmhp_print.BgGreen},
		{6, 8, cmhp_print.BgRed},
		{6 + 8, 8 * 8, cmhp_print.BgBlue},
		{6 + 8 + 8*8, 1, cmhp_print.BgRed},
		{6 + 8 + 8*8 + 1, 4 * 8, cmhp_print.BgGreen},
	})
}
