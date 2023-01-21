package pack

import (
	"reflect"
	"unsafe"
)

type (
	TpInfo struct {
		HeaderSize int
		Name       [][]byte
		Offset     []uintptr
		Type       []uint8

		MapOffset map[string]uintptr
		MapType   map[string]uint8
	}
)

var cache = map[reflect.Type]*TpInfo{}

func getTypeInfo[T comparable](v *T) *TpInfo {
	vv, ok := cache[reflect.TypeOf(v).Elem()]
	if ok {
		return vv
	}

	typeOf := reflect.TypeOf(v).Elem()
	valueOf := reflect.ValueOf(v).Elem()
	info := TpInfo{}

	info.HeaderSize = 1 // num of fields
	info.MapOffset = map[string]uintptr{}
	info.MapType = map[string]uint8{}

	for i := 0; i < typeOf.NumField(); i++ {
		info.Offset = append(info.Offset, typeOf.Field(i).Offset)
		info.MapOffset[typeOf.Field(i).Name] = typeOf.Field(i).Offset
		info.Name = append(info.Name, []byte(typeOf.Field(i).Name))
		info.HeaderSize += 1                         // field len
		info.HeaderSize += len(typeOf.Field(i).Name) // len
		info.HeaderSize += 1                         // type

		switch valueOf.Field(i).Interface().(type) {
		case uint8:
			info.Type = append(info.Type, _uint8)
			info.MapType[typeOf.Field(i).Name] = _uint8
			info.HeaderSize += 1
		case uint16:
			info.Type = append(info.Type, _uint16)
			info.MapType[typeOf.Field(i).Name] = _uint16
			info.HeaderSize += 2
		case uint32:
			info.Type = append(info.Type, _uint32)
			info.MapType[typeOf.Field(i).Name] = _uint32
			info.HeaderSize += 4
		case string:
			info.Type = append(info.Type, _string)
			info.MapType[typeOf.Field(i).Name] = _string
			info.HeaderSize += 4 // length of string
		default:
			panic("unknown type")
		}
	}

	cache[typeOf] = &info
	return &info
}

func Pack[T comparable](v *T) []byte {
	info := getTypeInfo(v)

	// pos in array
	p := 0

	// Start of struct
	start := uintptr(unsafe.Pointer(v))

	// Calculate string size
	size := info.HeaderSize
	for i := 0; i < len(info.Type); i++ {
		if info.Type[i] == _string {
			bb := *(*[]byte)(unsafe.Pointer(start + info.Offset[i]))
			ln := len(bb)
			size += ln
		}
	}

	// Prepare buffer
	out := make([]byte, size)

	// Header
	p += __write8(&out, p, uint8(len(info.Type))) // Field amount

	// Content
	for i := 0; i < len(info.Offset); i++ {
		// Field name
		p += __copySmallString(&out, p, uintptr(unsafe.Pointer(&info.Name[i])))

		// Type info
		p += __write8(&out, p, info.Type[i])

		switch info.Type[i] {
		case _uint8:
			p += __copy8(&out, p, start, info.Offset[i])
			break
		case _uint16:
			p += __copy16(&out, p, start, info.Offset[i])
			break
		case _uint32:
			p += __copy32(&out, p, start, info.Offset[i])
			break
		case _string:
			p += __copyBigString(&out, p, start+info.Offset[i])
			break
		}
	}

	// 48 nc - base

	return out
}

func Unpack[T comparable](b *[]byte) T {
	out := new(T)
	info := getTypeInfo(out)

	// Start of struct
	startOfStruct := uintptr(unsafe.Pointer(out))
	// startOfBin := uintptr(unsafe.Pointer(b))

	// Read fields number
	p := 0
	numOfFields := (*b)[p]
	p++

	// Read fields
	for i := 0; i < int(numOfFields); i++ {
		// Read field length
		fieldLength := int((*b)[p])
		p += 1

		// Read field name
		fieldName := string((*b)[p : p+fieldLength])
		p += fieldLength

		// Read field type
		vType := (*b)[p]
		p += 1

		// Field offset
		offset, _ := info.MapOffset[fieldName]

		switch vType {
		case _uint8:
			*(*uint8)(unsafe.Pointer(startOfStruct + offset)) = (*b)[p]
			p += 1
			break
		case _uint16:
			*(*uint16)(unsafe.Pointer(startOfStruct + offset)) = __read16(b, p)
			p += 2
			break
		case _uint32:
			*(*uint32)(unsafe.Pointer(startOfStruct + offset)) = __read32(b, p)
			p += 4
			break
		case _string:
			// string length
			strLength := __read32(b, p)
			p += 4

			// Copy
			str := make([]byte, strLength)
			copy(str, (*b)[p:])
			*(*string)(unsafe.Pointer(startOfStruct + offset)) = string(str)
			p += int(strLength)

			break
		}
	}

	return *out
}
