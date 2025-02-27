package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/jpeg"
	"io"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
)

// dataは変換元, box(ポインタ)は変換後
func MarshalAndInsert(data any, box any) {
	marshaledData, _ := json.Marshal(data)
	json.Unmarshal(marshaledData, box)
}

func GenId() string {
	uuidWithHyphen, _ := uuid.NewRandom()
	return uuidWithHyphen.String()
}

// Pointerにしたい値を入力する
func Ptr[T any](x T) *T {
	return &x
}

func Aes256Encode(plaintext string, key string) *string {
	bPlaintext := pkcs5Padding(plaintext)
	block, _ := aes.NewCipher([]byte(getMD5Hash(key)))
	ciphertext := make([]byte, len(bPlaintext))
	mode := cipher.NewCBCEncrypter(block, []byte(os.Getenv("IV")))
	mode.CryptBlocks(ciphertext, bPlaintext)
	return Ptr(hex.EncodeToString(ciphertext))
}

func pkcs5Padding(plaintext string) []byte {
	ciphertext := []byte(plaintext)
	padding := (aes.BlockSize - len(ciphertext)%aes.BlockSize)
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func getMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func unpad(b []byte) string {
	return string(b[:len(b)-int(b[len(b)-1])])
}

func Aes256Decode(cipherText string, key string) (decryptedString string) {
	cipherTextDecoded, _ := hex.DecodeString(cipherText)
	block, _ := aes.NewCipher([]byte(getMD5Hash(key)))
	mode := cipher.NewCBCDecrypter(block, []byte(os.Getenv("IV")))
	mode.CryptBlocks([]byte(cipherTextDecoded), []byte(cipherTextDecoded))
	return unpad(cipherTextDecoded)
}

func MakeHMAC(msg, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(msg))
	return hex.EncodeToString(mac.Sum(nil))
}

func Base64Enc(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func NewDTO[T1 any, T2 any](model T1) (dto *T2) {
	MarshalAndInsert(model, &dto)
	return
}

func NewDTOs[T1 any, T2 any](models []T1) (dtos []T2) {
	MarshalAndInsert(models, &dtos)
	return
}

func EnhanceResponseWriter(w *http.ResponseWriter) {
	(*w).Header().Set("Content-Type", "application/json")
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

func GetIdFromPath(p string) (id string) {
	p = path.Clean("/" + p)
	slc := strings.Split(p, "/")
	return slc[len(slc)-1]
}

func ImageToReader(img image.Image) (io.Reader, error) {
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, img, nil) // または jpeg.Encode(buf, img, nil)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf.Bytes()), nil
}
