package mcrypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	KEY_SIZE = 16
	IV       = "ZigULaoyNevFdRED"
)

func main() {
	// key := []byte("1602ex_dota_x4Q95T")
	// key = []byte("1602ex_dota")
	// key = []byte("1602")
	// b := make([]byte, 16)
	// copy(b, key)
	// fmt.Println(b)
	// testAes()
	// s, _ := EncryptV1([]byte("5125453"), 1602)
	// // fmt.Println(fmt.Sprintf("%X", s))
	// s, _ = DecryptV1([]byte(s), 1602)
	// fmt.Println(s)
	// s, _ = Decrypt([]byte(s), []byte("1602"))
	// fmt.Println(s)
}

func reverseBytes(data []byte, length uint) {
	var half uint = length / 2
	var i uint
	var a byte
	for i = 0; i < half; i++ {
		a = data[i]
		data[i] = data[length-i-1]
		data[length-i-1] = a
	}
}

// key1 is ek, key2 is length, key3 is i;
func getMyKey(key1 uint, key2 uint, key3 uint) uint {
	var mykey uint = (key1*key2 + key3*key3) % 13

	if mykey < 0 {
		mykey = -mykey
	}

	if mykey >= 8 {
		mykey %= 7
	}

	if mykey == 0 { //  防止 mykey == 0不移位
		mykey = 7
	}

	return mykey
}

func EncryptV1(origData []byte, key int) (s string, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	strlen := uint(len(origData))
	result := origData
	reverseBytes(result, strlen)
	var mykey uint
	var i uint
	var a, b, c byte
	for i = 0; i < strlen; i++ {
		mykey = getMyKey(uint(key), strlen, i)

		a = result[i] << mykey
		b = result[i] >> (8 - mykey)
		c = b | a

		if c == 0x03 {
			a = c << mykey
			b = c >> (8 - mykey)
			c = b | a
		}

		result[i] = c
	}

	s = strings.ToUpper(hex.EncodeToString(result))
	return
}

func DecryptV1(origData []byte, key int) (s string, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	origData, _ = hex.DecodeString(string(origData))
	strlen := uint(len(origData))
	result := origData
	var mykey uint
	var i uint
	var a, b, c byte
	for i = 0; i < strlen; i++ {
		mykey = getMyKey(uint(key), strlen, i)
		a = result[i] >> mykey
		b = result[i] << (8 - mykey)
		c = b | a

		if c == 0x03 {
			a = c >> mykey
			b = c << (8 - mykey)
			c = b | a
		}
		result[i] = c
	}

	reverseBytes(result, strlen)
	s = string(result)
	return
}

func testAes() {
	// AES-128。key长度：16, 24, 32 bytes 对应 AES-128, AES-192, AES-256
	key := []byte("1602ex_dota_x4Q95T")
	key = []byte("1602")
	b := make([]byte, 16)
	copy(b, key)
	iv := []byte("ZigULaoyNevFdRED")
	result, err := AesEncrypt([]byte("5125453"), b, iv)
	if err != nil {
		panic(err)
	}
	fmt.Println(base64.StdEncoding.EncodeToString(result))
	origData, err := AesDecrypt(result, b, iv)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(origData))
}

func EncryptV2(origData, key []byte) (s string, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	var t []byte
	if len(key) > KEY_SIZE {
		t, err = AesEncrypt(origData, key[:16], []byte(IV))
		if err != nil {
			return
		}
	} else if len(key) < KEY_SIZE {
		b := make([]byte, 16)
		copy(b, key)
		t, err = AesEncrypt(origData, b, []byte(IV))
		if err != nil {
			return
		}
	} else {
		t, err = AesEncrypt(origData, key, []byte(IV))
		if err != nil {
			return
		}
	}
	r := strings.NewReplacer("+", "-", "/", "_")
	s = base64.StdEncoding.EncodeToString(t)
	s = base64.StdEncoding.EncodeToString([]byte(s))
	s = r.Replace(s)
	s = strings.TrimRight(s, "=")

	return
}

func DecryptV2(origData, key []byte) (s string, err error) {
	defer func() {
		if err := recover(); err != nil {
			return
		}
	}()
	var t []byte

	r := strings.NewReplacer("-", "+", "_", "/")
	s = r.Replace(string(origData))
	t = append([]byte(s), bytes.Repeat([]byte("="), len(s)%4)...)
	t, err = base64.StdEncoding.DecodeString(string(t))
	t, err = base64.StdEncoding.DecodeString(string(t))

	if len(key) > KEY_SIZE {
		t, err = AesDecrypt(t, key[:16], []byte(IV))
		if err != nil {
			return
		}
	} else if len(key) < KEY_SIZE {
		b := make([]byte, 16)
		copy(b, key)
		t, err = AesDecrypt(t, b, []byte(IV))
		if err != nil {
			return
		}
	} else {
		t, err = AesDecrypt(t, key, []byte(IV))
		if err != nil {
			return
		}
	}

	s = string(t)
	return
}

func AesEncrypt(origData, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key[:16])
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, iv)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func AesDecrypt(crypted, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	//blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv)
	origData := make([]byte, len(crypted))
	// origData := crypted
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return origData, nil
}

func ZeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding)
	return append(ciphertext, padtext...)
}

func ZeroUnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	// 去掉最后一个字节 unpadding 次
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
