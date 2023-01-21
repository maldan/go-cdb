package util_test

import (
	"encoding/json"
	"fmt"
	"github.com/maldan/go-cdb/pack"
	"testing"
	"unsafe"
)

type Test struct {
	Id  uint8  `json:"id"`
	Id2 uint16 `json:"id2"`
	Id3 uint32 `json:"id3"`

	S string `json:"s"`

	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth"`

	Age        uint32 `json:"age"`
	UserId     uint32 `json:"user_id"`
	DocumentId uint32 `json:"document_id"`
	ReviewId   uint32 `json:"review_id"`

	Address               string `json:"address"`
	City                  string `json:"city"`
	State                 string `json:"state"`
	ZipCode               string `json:"zip_code"`
	Email                 string `json:"email"`
	Phone                 string `json:"phone"`
	EmergencyPersonName   string `json:"emergency_person_name"`
	EmergencyPersonPhone  string `json:"emergency_person_phone"`
	ParentOrLegalGuardian string `json:"parent_or_legal_guardian"`
	Date                  string `json:"date"`

	// Arr [4]int `json:"arr"`

	/*
		ReviewEnabled         bool   `json:"review_enabled"`

		Created time.Time `json:"created"`*/
}

/*func TestA(t *testing.T) {
	vv := util.Pack(Test{Id: 2, Id2: 4, Id3: 8, S: "hi"})
	fmt.Printf("%v\n", vv)
	fmt.Printf("%v\n", string(vv))
	util.WriteBin("t.bin", vv)
}
*/

func TestA2(t *testing.T) {
	y := [4]rune{'a', 'с'}
	fmt.Printf("%v\n", y[0])
	fmt.Printf("%v\n", unsafe.Sizeof(Test{}))
	//b1 := []byte{1, 2, 3, 4, 5}
	//b2 := []byte{10, 10, 10, 10, 10}
	//copy(b1[3:], b2)

	// fmt.Printf("%v\n", reflect.TypeOf([4]int{}).Bits())

	vv := pack.Pack(&Test{
		Id:                   1,
		Id2:                  2,
		Id3:                  4,
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
		Date:                 "suka",
	})
	x := pack.Unpack[Test](&vv)
	fmt.Printf("%+v\n", x)

	/*util.WriteBin("sas.bin", vv)
	util.WriteJson("sas.json", vv)*/
}

/*func Benchmark1(b *testing.B) {
	v := Test{Id: 2, Id2: 4, Id3: 8, S: "hi"}
	prt := unsafe.Pointer(&v)
	t := reflect.TypeOf(&v).Elem()
	fooID := xunsafe.FieldByName(t, "Id")
	for i := 0; i < b.N; i++ {

		prt = unsafe.Pointer(&v)
		t = reflect.TypeOf(&v).Elem()
		fooID = xunsafe.FieldByName(t, "Id")
	}
	fmt.Printf("%v\n", prt)
	fmt.Printf("%v\n", t)
	fmt.Printf("%v\n", fooID)
}*/

func BenchmarkJson(b *testing.B) {
	a := Test{
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
	}
	for i := 0; i < b.N; i++ {
		v, _ := json.Marshal(&a)
		b.SetBytes(int64(len(v)))
	}
}

func BenchmarkMy2(b *testing.B) {
	a := Test{
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
	}
	for i := 0; i < b.N; i++ {
		v := pack.Pack(&a)
		b.SetBytes(int64(len(v)))
	}
}

func BenchmarkUnJson(b *testing.B) {
	a := Test{
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
	}
	bytes, _ := json.Marshal(a)
	vv := Test{}

	for i := 0; i < b.N; i++ {
		json.Unmarshal(bytes, &vv)
		b.SetBytes(int64(len(bytes)))
	}
}

func BenchmarkUnMy(b *testing.B) {
	a := Test{
		Address:              "1231234gfdfg dfsdhsfghdfj sd sdfh ssj",
		State:                "Сукаа маруляяяя",
		FirstName:            "Сукаа маруляяяя",
		LastName:             "Сукаа маруляяяя",
		DateOfBirth:          "Сукаа маруляяяя",
		Email:                "Сукаа маруляяяя",
		Phone:                "Сукаа маруляяяя",
		EmergencyPersonName:  "Сукаа маруляяяя",
		EmergencyPersonPhone: "Сукаа маруляяяя",
	}
	bytes := pack.Pack(&a)

	for i := 0; i < b.N; i++ {
		pack.Unpack[Test](&bytes)
		b.SetBytes(int64(len(bytes)))
	}
}

/*func BenchmarkR(b *testing.B) {
	ok := false
	var cache2 = map[reflect.Type]*util.TpInfo{}
	var tt = &Test{Id: 1, Id2: 2, Id3: 3, S: "hi"}
	cache2[reflect.TypeOf(tt).Elem()] = &util.TpInfo{}
	for i := 0; i < b.N; i++ {
		_, ok = cache2[reflect.TypeOf(tt).Elem()]
	}
	if !ok {
		fmt.Printf("%v", ok)
	}
}

func BenchmarkMM(b *testing.B) {
	sas := 1
	out := make([]byte, sas)
	for i := 0; i < b.N; i++ {
		if i == 10 {
			sas = 32
		}
		out = make([]byte, sas)
	}
	if sas == -1 {
		fmt.Printf("%v", out)
	}
}*/

/*func BenchmarkG1(b *testing.B) {
	var a = Test{Id: 1, Id2: 4, Id3: 4, S: "hi"}
	for i := 0; i < b.N; i++ {
		reflect.TypeOf(&a).Elem()
	}
}

func BenchmarkG2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		reflect.TypeOf(&Test{Id: 1, Id2: 4, Id3: 4, S: "hi"}).Elem()
	}
}*/

/*func BenchmarkR2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		util.Pack2(&Test{Id: uint8(i), Id2: 4, Id3: uint32(i), S: "hi"})
	}
}
*/
