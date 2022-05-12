package helpers

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"time"

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

func CreateFile(data []byte, pathFile string) error {
	file, err := os.Create(pathFile)
	if err != nil {
		os.Exit(1)
		return errors.New("invalid file path")
	}

	_, err = file.Write(data)
	if err != nil {
		os.Exit(1)
		return fmt.Errorf("error when write file %s: %s \n", pathFile, err)
	}

	return nil
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

func GenerateRandomString(length int) string {
	var charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789" +
		"!@#$%&*_+"

	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func CreateFolder(folders []string) {
	jsonStr, _ := json.Marshal(folders)

	var data []string
	_ = json.Unmarshal(jsonStr, &data)

	for _, name := range data {
		if _, err := os.Stat(name); err != nil {
			os.Mkdir(name, os.ModePerm)
		}
	}
}
