package helpers

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

//EnvVar function is for read .env file
func EnvVar(key string, defaultVal string) string {
	godotenv.Load()
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultVal
	}
	return value
}

func FileExists(filepath string) bool {
	_, err := os.Stat(filepath)
	return !os.IsNotExist(err)
}

func MergeFilters(newFilter map[string]interface{}, filter interface{}) interface{} {
	iter := reflect.ValueOf(filter).MapRange()
	for iter.Next() {
		k := iter.Key()
		v := iter.Value()

		newFilter[k.String()] = v.Interface()
	}

	return newFilter
}

func ConvertImageToBase64(pathImage string) (string, error) {
	data, err := ioutil.ReadFile(pathImage)
	if err != nil {
		return "", err
	}

	imgBase64Str := base64.StdEncoding.EncodeToString(data)

	return imgBase64Str, nil
}
