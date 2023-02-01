package dson

/*func S(s *[]byte, offset int, v any) int {
	valueOf := reflect.ValueOf(v)
	typeOf := reflect.TypeOf(v)

	if typeOf.Kind() == reflect.Struct {
		if typeOf.Name() == "Time" {
			vx := valueOf.Interface().(time.Time).Format("2006-05-04T15:04:05.999-07:00")
			offset = WriteString(s, offset, vx)
		} else {
			for i := 0; i < typeOf.NumField(); i++ {
				offset = WriteField(s, offset, typeOf.Field(i).Type.Kind(), typeOf.Field(i).Name)
				offset = S(s, offset, valueOf.Field(i).Interface())
			}
		}
	}

	if typeOf.Kind() == reflect.Slice {
		for i := 0; i < valueOf.Len(); i++ {
			offset = S(s, offset, valueOf.Index(i).Interface())
		}
	}

	if typeOf.Kind() == reflect.String {
		offset = WriteString(s, offset, valueOf.String())
	}

	if typeOf.Kind() == reflect.Int {
		*s = append(*s, 0, 0, 0, 0)
		binary.LittleEndian.PutUint32((*s)[offset:], uint32(valueOf.Int()))
		offset += 4
	}
	return offset
}

func Traverser(s []byte) {
	scope := make([]string, 16)
	scopeIndex := 0
	currentContent := ""
	isCaptureMode := false
	isKeyMode := true

	for i := 0; i < len(s); i++ {
		if isCaptureMode {
			if s[i] == '"' {
				if isKeyMode {
					fmt.Printf("K: %v\n", currentContent)
				} else {
					fmt.Printf("V: %v\n", currentContent)
				}

				isCaptureMode = false
				continue
			}
			currentContent += string(s[i])
		}
		if s[i] == ',' {
			isKeyMode = true
		}
		if s[i] == ':' {
			isKeyMode = false
		}
		if s[i] == '"' {
			currentContent = ""
			isCaptureMode = true
		}
		if s[i] == '{' {
			scope[scopeIndex] = ""
			scopeIndex += 1
		}
		if s[i] == '}' {
			scopeIndex -= 1
		}
		//fmt.Printf("%v\n", s[i])
	}
}

func Expand(s *[]byte, amt int) {
	for i := 0; i < amt; i++ {
		*s = append(*s, 0)
	}
}

func WriteField(s *[]byte, offset int, kind reflect.Kind, name string) int {
	Expand(s, 1) // type
	switch kind {
	case reflect.String:
		(*s)[offset] = TypeString
		break
	case reflect.Int:
		(*s)[offset] = TypeI32
		break
	case reflect.Slice:
		(*s)[offset] = TypeSlice
		break
	case reflect.Struct:
		(*s)[offset] = TypeStruct
		break
	default:
		panic("Gas gas gaaaas I'm feeling something in my ass")
	}

	offset += 1

	Expand(s, 1) // len
	(*s)[offset] = uint8(len(name))
	offset += 1

	Expand(s, len(name)) // name
	copy((*s)[offset:], name)
	offset += len(name)

	return offset
}

func WriteString(s *[]byte, offset int, value string) int {
	Expand(s, 1) // len
	(*s)[offset] = uint8(len(value))
	offset += 1

	Expand(s, len(value)) // name
	copy((*s)[offset:], value)
	offset += len(value)

	return offset
}*/
