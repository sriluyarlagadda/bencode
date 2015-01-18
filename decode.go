package bencode

import (
	"bufio"
	"errors"
	"io"
	"strconv"
)

func Decode(r io.Reader) (interface{}, error) {

	bufferedReader := bufio.NewReader(r)

	if isString, stringValue, err := decodeString(bufferedReader); err == nil {
		if isString {
			return stringValue, nil
		}
	} else {
		return nil, err
	}

	if isInt, intValue, err := decodeInt(bufferedReader); err == nil {
		if isInt {
			return intValue, nil
		}
	} else {
		return nil, err
	}

	if isList, listValue, err := decodeList(bufferedReader); err == nil {
		if isList {
			return listValue, nil
		}
	} else {
		return nil, err
	}

	if isDict, dictValue, err := decodeDict(bufferedReader); err == nil {
		if isDict {
			return dictValue, nil
		}
	} else {
		return nil, err
	}

	return 0, errors.New("there is error")
}

func readNext(r *bufio.Reader) (rune, error) {
	rune, _, err := r.ReadRune()
	if err != nil {
		return rune, err
	}

	return rune, nil
}

func unreadNext(r *bufio.Reader) error {
	err := r.UnreadRune()
	if err != nil {
		return err
	}
	return nil
}

func decodeString(r *bufio.Reader) (bool, string, error) {
	value, isByteString, err := isByteString(r)

	if err != nil {
		return false, "", err
	}

	if isByteString {
		byteString, err := parseByteString(value, r)
		if err != nil {
			return false, "", err
		}
		return true, byteString, nil
	}

	return false, "", nil
}

func isByteString(r *bufio.Reader) (int, bool, error) {
	nextRune, err := readNext(r)
	if err != nil {
		return 0, false, err
	}

	err = unreadNext(r)
	if err != nil {
		return 0, false, err
	}

	//check for integer
	value, err := strconv.Atoi(string(nextRune))

	if err != nil {
		return 0, false, nil
	}
	return value, true, nil
}

func parseByteString(count int, r *bufio.Reader) (string, error) {
	//ignore count and colon
	readNext(r)
	nextRune, err := readNext(r)

	if err != nil {
		return "", err
	}

	if nextRune != rune(':') {
		return "", errors.New("expected : for string but could not find")
	}

	var runes []rune
	for i := 0; i < count; i++ {

		value, err := readNext(r)
		if err != nil {
			return "", err
		}

		runes = append(runes, value)
	}

	return string(runes), nil
}

func decodeInt(r *bufio.Reader) (bool, int, error) {
	isInt, err := isInt(r)
	if err != nil {
		return false, 0, err
	}

	if isInt {
		intValue, err := parseInt(r)
		if err != nil {
			return false, 0, err
		}
		return true, intValue, nil
	}
	return false, 0, nil
}

func isInt(r *bufio.Reader) (bool, error) {
	nextRune, err := readNext(r)
	if err != nil {
		return false, err
	}

	err = unreadNext(r)
	if err != nil {
		return false, err
	}

	if nextRune != rune('i') {
		return false, nil
	}

	return true, nil
}

func parseInt(r *bufio.Reader) (int, error) {

	//ignore the 'i' at the beginning
	readNext(r)
	var runes []rune
	for value, _ := readNext(r); value != rune('e'); value, _ = readNext(r) {
		runes = append(runes, value)
	}

	return strconv.Atoi(string(runes))
}

func decodeList(r *bufio.Reader) (bool, []interface{}, error) {
	isList, err := isList(r)
	if err != nil {
		return false, nil, err
	}

	if isList {
		listValue, err := parseList(r)
		if err != nil {
			return false, nil, err
		}
		return true, listValue, nil
	}
	return false, nil, nil
}

func isList(r *bufio.Reader) (bool, error) {
	nextRune, err := readNext(r)
	if err != nil {
		return false, err
	}

	err = unreadNext(r)
	if err != nil {
		return false, err
	}

	if nextRune != rune('l') {
		return false, nil
	}
	return true, nil
}

func parseList(r *bufio.Reader) ([]interface{}, error) {
	//ignore the 'l' at the beginning
	readNext(r)
	var bencodeTypes []interface{}

	for value, _ := readNext(r); value != rune('e'); value, _ = readNext(r) {
		unreadNext(r)

		if isInt, intValue, err := decodeInt(r); err == nil {
			if isInt {
				bencodeTypes = append(bencodeTypes, intValue)
			}
		} else {
			return nil, err
		}

		if isString, stringValue, err := decodeString(r); err == nil {
			if isString {
				bencodeTypes = append(bencodeTypes, stringValue)
			}
		} else {
			return nil, err
		}
	}
	return bencodeTypes, nil
}

func decodeDict(r *bufio.Reader) (bool, map[string]interface{}, error) {
	isDict, err := isDict(r)
	if err != nil {
		return false, nil, err
	}

	if isDict {
		distValue, err := parseDict(r)
		if err != nil {
			return false, nil, err
		}
		return true, distValue, nil
	}
	return false, nil, nil
}

func isDict(r *bufio.Reader) (bool, error) {
	nextRune, err := readNext(r)
	if err != nil {
		return false, err
	}

	unreadNext(r)

	if nextRune != rune('d') {
		return false, nil
	}

	return true, nil
}

func parseDict(r *bufio.Reader) (map[string]interface{}, error) {
	//ignore the 'd' at the beginning
	readNext(r)
	dict := make(map[string]interface{})
	for value, _ := readNext(r); value != rune('e'); value, _ = readNext(r) {
		unreadNext(r)
		var key string

		if isString, stringValue, err := decodeString(r); err == nil {
			if isString {
				key = stringValue
			} else {
				return nil, errors.New("key should be strings")
			}
		} else {
			return nil, err
		}

		if isInt, intValue, err := decodeInt(r); err == nil {
			if isInt {
				dict[key] = intValue
				continue
			}
		} else {
			return nil, err
		}

		if isList, listValue, err := decodeList(r); err == nil {
			if isList {
				dict[key] = listValue
				continue
			}
		} else {
			return nil, err
		}

		if isString, stringValue, err := decodeString(r); err == nil {
			if isString {
				dict[key] = stringValue
				continue
			}
		} else {
			return nil, err
		}

	}

	return dict, nil
}
