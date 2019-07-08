package util

import (
	"os"
	"strings"
	"bytes"
	"errors"
	"io/ioutil"
	"path"
	"encoding/json"
	"regexp"
	"fmt"
	"reflect"
)

const extendTag = "@extend:"
const sep = string(os.PathSeparator)

func IsEmpty (key string) bool {
	return strings.Count(key, "") - 1 <= 0
}

func IsNagative (key int32) bool {
	return key <= 0
}

func Contains (seq,target string) bool {
	return strings.Contains(seq,target)
}

/*
  Splicing all strings in a more efficient way
*/
func Join(keys ...string) string {
	var buffer bytes.Buffer
	for _,key := range keys {
		buffer.WriteString(key)
	}
	return buffer.String()
}

func ExtendFile(filePath string) (string, error) {
	fi, err := os.Stat(filePath)
	if err != nil {
		return "", err
	} else if fi.IsDir() {
		return "", errors.New(filePath + " is not a file.")
	}
	var b []byte
	b, err = ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return ExtendFileContent(path.Dir(filePath), b)
}

func ExtendFileContent(dir string, content []byte) (string, error) {
	//检查是不是规范的json
	test := new(interface{})
	err := json.Unmarshal(content, &test)
	if err != nil {
		return "", err
	}

	//替换子json文件
	reg := regexp.MustCompile(`"` + extendTag + `.*?"`)
	ret := reg.ReplaceAllStringFunc(string(content), func(match string) string {
		match = match[len(extendTag)+1 : len(match)-1]
		var p = match
		if !strings.HasPrefix(match, sep) {
			p = dir + sep + match
		}
		sb, err2 := ExtendFile(p)
		if err2 != nil {
			err = fmt.Errorf("替换json配置[%s]失败：%s\n", match, err2.Error())
		}
		return string(sb)
	})
	return ret, err
}

func Contain(obj interface{}, target interface{}) (bool, error) {
    targetValue := reflect.ValueOf(target)
    switch reflect.TypeOf(target).Kind() {
    case reflect.Slice, reflect.Array:
        for i := 0; i < targetValue.Len(); i++ {
            if targetValue.Index(i).Interface() == obj {
                return true, nil
            }
        }
    case reflect.Map:
        if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
            return true, nil
        }
    }
    return false, errors.New("not in array")
}