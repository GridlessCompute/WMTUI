package wmapi

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/GehirnInc/crypt/md5_crypt"
	"github.com/andreburgaud/crypt2go/ecb"
)

// WhatsminerAccessToken represents a reusable token to access and/or control a single Whatsminer ASIC.
type WhatsminerAccessToken struct {
	Created       time.Time
	IPAddress     string
	Port          int
	AdminPassword string
	Cipher        cipher.Block
	Sign          string
}

type ApiError struct {
	Code        int    `json:"Code"`
	Description string `json:"Description"`
	Msg         string `json:"Msg"`
	Status      string `json:"STATUS"`
	When        int    `json:"When"`
}

type SummaryS struct {
	Status []struct {
		Msg    string `json:"Msg"`
		Status string `json:"STATUS"`
	} `json:"STATUS"`
	Summary []struct {
		Accepted        int     `json:"Accepted"`
		BtminerFastBoot string  `json:"Btminer Fast Boot"`
		ChipTempAvg     float64 `json:"Chip Temp Avg"`
		ChipTempMax     float64 `json:"Chip Temp Max"`
		ChipTempMin     float64 `json:"Chip Temp Min"`
		Debug           string  `json:"Debug"`
		Elapsed         int     `json:"Elapsed"`
		EnvTemp         float64 `json:"Env Temp"`
		FactoryGHS      int     `json:"Factory GHS"`
		FanSpeedIn      int     `json:"Fan Speed In"`
		FanSpeedOut     int     `json:"Fan Speed Out"`
		HSRT            float64 `json:"HS RT"`
		HashDeviation   float64 `json:"Hash Deviation%"`
		//HashStable            bool    `json:"Hash Stable"`
		//HashStableCostSeconds int     `json:"Hash Stable Cost Seconds"`
		MHS15M       float64 `json:"MHS 15m"`
		MHS1M        float64 `json:"MHS 1m"`
		MHS5M        float64 `json:"MHS 5m"`
		MHS5S        float64 `json:"MHS 5s"`
		MHSAv        float64 `json:"MHS av"`
		PoolRejected float64 `json:"Pool Rejected%"`
		PoolStale    int     `json:"Pool Stale%"`
		Power        int     `json:"Power"`
		PowerLimit   int     `json:"Power Limit"`
		PowerMode    string  `json:"Power Mode"`
		PowerRate    float64 `json:"Power Rate"`
		Rejected     int     `json:"Rejected"`
		SecurityMode int     `json:"Security Mode"`
		TargetFreq   int     `json:"Target Freq"`
		TargetMHS    int     `json:"Target MHS"`
		Temperature  float64 `json:"Temperature"`
		TotalMH      int64   `json:"Total MH"`
		Uptime       int     `json:"Uptime"`
		FreqAvg      int     `json:"freq_avg"`
	} `json:"SUMMARY"`
	ID int `json:"id"`
}

type WMstruct struct {
	Enc  int    `json:"enc"`
	Data string `json:"data"`
}

//	type ErrorsMessage struct {
//		Code        int    `json:"Code"`
//		Description string `json:"Description"`
//		Msg         Msg    `json:"Msg"`
//		Status      string `json:"STATUS"`
//		When        int    `json:"When"`
//	}
//
//	type Msg struct {
//		ErrorCode []any `json:"error_code"`
//	}
type GetErrorCodeS struct {
	Code        int    `json:"Code"`
	Description string `json:"Description"`
	Msg         struct {
		ErrorCode any `json:"error_code"`
	}
	Status string `json:"STATUS"`
	When   int    `json:"When"`
}

type GetMinerInfoS struct {
	Code        int    `json:"Code"`
	Description string `json:"Description"`
	Msg         struct {
		DNS      string `json:"dns"`
		Gateway  string `json:"gateway"`
		Hostname string `json:"hostname"`
		IP       string `json:"ip"`
		Ledstat  string `json:"ledstat"`
		Mac      string `json:"mac"`
		Minersn  string `json:"minersn"`
		Netmask  string `json:"netmask"`
		Powersn  string `json:"powersn"`
		Proto    string `json:"proto"`
	} `json:"Msg"`
	Status string `json:"STATUS"`
	When   int    `json:"When"`
}

type GetMinerInfo2S struct {
	Status string `json:"STATUS"`
	When   int    `json:"When"`
	Code   int    `json:"Code"`
	Msg    struct {
		IP       string `json:"ip"`
		Proto    string `json:"proto"`
		Netmask  string `json:"netmask"`
		Gateway  string `json:"gateway"`
		DNS      string `json:"dns"`
		Hostname string `json:"hostname"`
		Mac      string `json:"mac"`
		Ledstat  string `json:"ledstat"`
	} `json:"Msg"`
	Description string `json:"Description"`
}

type GetPoolInfoS struct {
	Pools []struct {
		Accepted            int     `json:"Accepted"`
		BadWork             int     `json:"Bad Work"`
		CurrentBlockHeight  int     `json:"Current Block Height"`
		CurrentBlockVersion int     `json:"Current Block Version"`
		Discarded           int     `json:"Discarded"`
		GetFailures         int     `json:"Get Failures"`
		Getworks            int     `json:"Getworks"`
		LastShareTime       int     `json:"Last Share Time"`
		Pool                int     `json:"POOL"`
		PoolRejected        float64 `json:"Pool Rejected%"`
		PoolStale           float64 `json:"Pool Stale%"`
		Priority            int     `json:"Priority"`
		Quota               int     `json:"Quota"`
		Rejected            int     `json:"Rejected"`
		RemoteFailures      int     `json:"Remote Failures"`
		Stale               int     `json:"Stale"`
		Status              string  `json:"Status"`
		StratumActive       bool    `json:"Stratum Active"`
		StratumDifficulty   float64 `json:"Stratum Difficulty"`
		ToRemove            bool    `json:"To Remove"`
		URL                 string  `json:"URL"`
		User                string  `json:"User"`
		Works               int     `json:"Works"`
	} `json:"POOLS"`
	Status []struct {
		Msg    string `json:"Msg"`
		Status string `json:"STATUS"`
	} `json:"STATUS"`
	ID int `json:"id"`
}

// NewWhatsminerAccessToken creates a new instance of WhatsminerAccessToken.
func NewWhatsminerAccessToken(ipAddress string, port int, adminPassword string) (*WhatsminerAccessToken, error) {
	token := &WhatsminerAccessToken{
		Created:   time.Now(),
		IPAddress: ipAddress,
		Port:      port,
	}

	if adminPassword != "" {
		token.EnableWriteAccess(adminPassword)
	}

	return token, nil
}

// initializeWriteAccess initializes write access for the token.
func (t *WhatsminerAccessToken) initializeWriteAccess(adminPassword string) error {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.IPAddress, t.Port))
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(`{"cmd": "get_token"}`))
	if err != nil {
		return err
	}

	buffer := make([]byte, 4000)
	n, err := conn.Read(buffer)
	if err != nil {
		return err
	}

	tokenInfo := make(map[string]interface{})
	err = json.Unmarshal(buffer[:n], &tokenInfo)
	if err != nil {
		return err
	}
	// log.Printf("TokenInfo: %v\n", tokenInfo)
	if msg, ok := tokenInfo["Msg"].(string); ok && msg == "over max connect" {
		return errors.New(msg)
	}

	salt := fmt.Sprintf("$1$%s$", tokenInfo["Msg"].(map[string]interface{})["salt"])
	// log.Printf("Salt: %v\n", salt)
	r := regexp.MustCompile(`\s*\$(\d+)\$([\w\./]*)\$`)
	match := r.FindStringSubmatch(salt)
	if match == nil {
		return errors.New("salt format is not correct")
	}
	log.Printf("Match: %v\n", match)
	m := md5_crypt.New()
	result, err := m.Generate([]byte(adminPassword), []byte(salt))
	if err != nil {
		// log.Printf("Error: %v\n", err)
	}
	// log.Printf("Hash: %v\n", result)

	//pwd := pbkdf2.Key([]byte(adminPassword), []byte(salt), 1000, 16, sha3.New256)
	//key := hex.EncodeToString(pwd)
	key := strings.Split(result, "$")[3]
	// log.Printf("Key: %v\n", key)
	aesKey := sha256.Sum256([]byte(key))
	// log.Printf("AesKey: %x\n", aesKey)

	t.Cipher, _ = aes.NewCipher(aesKey[:])

	newSalt := fmt.Sprintf("$1$%s$", tokenInfo["Msg"].(map[string]interface{})["newsalt"])
	// log.Printf("NewSalt: %v\n", newSalt)
	// log.Printf("SignVal: %v\n", key+fmt.Sprint(tokenInfo["Msg"].(map[string]interface{})["time"]))
	result, err = m.Generate([]byte(key+fmt.Sprint(tokenInfo["Msg"].(map[string]interface{})["time"])), []byte(newSalt))
	if err != nil {
		log.Printf("Error: %v\n", err)
	}
	log.Printf("NewHash: %v\n", result)

	//fmt.Println(tokenInfo["Msg"].(map[string]interface{}))
	//tmp = crypt(pwd[3] + token_info["time"], "$1$" + token_info["newsalt"] + '$')

	//signPwd := pbkdf2.Key([]byte(key+tokenInfo["Msg"].(map[string]interface{})["time"].(string)), []byte(newSalt), 1000, 16, sha3.New256)
	t.Sign = strings.Split(result, "$")[3]
	// log.Printf("Sign: %v\n", t.Sign)

	t.Created = time.Now()

	return nil
}

// EnableWriteAccess enables write access for the token.
func (t *WhatsminerAccessToken) EnableWriteAccess(adminPassword string) {
	t.AdminPassword = adminPassword
	t.initializeWriteAccess(adminPassword)
}

// HasWriteAccess checks write access and refreshes the token if necessary.
func (t *WhatsminerAccessToken) HasWriteAccess() bool {
	if t.AdminPassword == "" {
		return false
	}

	if time.Since(t.Created).Seconds() > 30*60 {
		// Writeable token has expired; reinitialize
		t.initializeWriteAccess(t.AdminPassword)
	}

	return true
}

// WhatsminerAPI represents a stateless class with only class methods for read/write API calls.
type WhatsminerAPI struct{}

// GetReadOnlyInfo sends a READ-ONLY API command.
func (w *WhatsminerAPI) GetReadOnlyInfo(accessToken *WhatsminerAccessToken, cmd string, additionalParams map[string]interface{}) (map[string]interface{}, error) {
	jsonCmd := map[string]interface{}{"cmd": cmd}
	if additionalParams != nil {
		for key, value := range additionalParams {
			jsonCmd[key] = value
		}
	}

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", accessToken.IPAddress, accessToken.Port), 5*time.Second)
	if err2 := recover(); err2 != nil || err != nil {
		return nil, err
	}
	defer conn.Close()

	err = json.NewEncoder(conn).Encode(jsonCmd)
	if err != nil {
		return nil, err
	}

	resp, err := readAll(conn)
	if err != nil {
		return nil, err
	}

	// Modify response to handle weird M31 issues
	resp = strings.ReplaceAll(resp, "inf", "999")
	resp = strings.ReplaceAll(resp, "nan", "0")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		log.Println("Error calling read-only endpoint")
		log.Println(resp)
		return nil, err
	}

	return result, nil
}

// ExecCommand sends a WRITEABLE API command.

func (w *WhatsminerAPI) ExecCommand(accessToken *WhatsminerAccessToken, cmd string, additionalParams map[string]interface{}) (map[string]interface{}, error) {
	if !accessToken.HasWriteAccess() {
		return nil, errors.New("access_token must have write access")
	}

	// Assemble the plaintext json
	jsonCmd := map[string]interface{}{"cmd": string(cmd), "token": accessToken.Sign}
	for key, value := range additionalParams {
		jsonCmd[key] = value
	}

	apiCmd, err := json.Marshal(jsonCmd)
	if err != nil {
		return nil, err
	}

	//
	// This Looks Correct
	//
	log.Printf("Json: %v\n", string(apiCmd))

	//mode := ecb.NewECBEncrypter(accessToken.Cipher)
	paddedCmd := addTo16(apiCmd)
	log.Printf("Padded: %v \n", string(paddedCmd))
	dst := make([]byte, len(paddedCmd))
	mode := ecb.NewECBEncrypter(accessToken.Cipher)
	mode.CryptBlocks(dst, paddedCmd)

	//mode.CryptBlocks(dst, addTo16(apiCmd))
	//for i, j := 0, 16; i < len(paddedCmd); i, j = i+16, j+16 {
	//	accessToken.Cipher.Encrypt(dst[i:j], paddedCmd[i:j])
	//}

	//accessToken.Cipher.Encrypt(dst, addTo16(apiCmd))
	encStr := string(
		base64.StdEncoding.EncodeToString(
			dst),
	)

	encStr = strings.ReplaceAll(encStr, "\n", "")

	//dataEnc := &WMstruct{1, encStr}

	dataEnc := map[string]interface{}{
		"enc":  1,
		"data": encStr,
	}

	apiPacketStr, _ := json.Marshal(dataEnc)
	log.Printf("Packet: %v\n", string(apiPacketStr))
	// Encrypt it and assemble the transport JSON
	// I Dont Think This is Correct
	// Makes 32 Long Byte Array
	//encMsg := make([]byte, len(addTo16(apiCmd)))
	//encMsg := make([]byte, aes.BlockSize+len(addTo16(apiCmd)))
	//log.Printf("[]byte: %v\n", encMsg)

	// log.Printf("addTo16(apiCmd): %v\n", string(addTo16(apiCmd)))

	// Or That This is correct
	// Encryptes apiCmd into encMsg -> Don't know if addTo16 is needed or Why Its Full Of "A"
	// DOes it need to be 64Bit -> probably
	//accessToken.Cipher.Encrypt(encMsg, addTo16(apiCmd))
	//accessToken.Cipher.Encrypt(encMsg, apiCmd)
	// PRINT
	//log.Printf("encMSG: %v\n", encMsg)
	// This just returns encMsg as a string -> its already full of "A"
	//encStr := base64.StdEncoding.EncodeToString(encMsg)
	// PRINT
	//log.Printf("encString: %v\n", encStr)
	// data and enc are backwards -> dont know if thats important
	//dataEnc := map[string]interface{}{"enc": 1} // transmit with "enc" to signal that it's encrypted
	//dataEnc["data"] = encStr

	// apiPacketString is the issue too many "AAAAAAAAAAAAA"
	//apiPacketStr, err := json.Marshal(dataEnc)
	if err != nil {
		return nil, err
	}

	//conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", accessToken.IPAddress, accessToken.Port))
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", accessToken.IPAddress, accessToken.Port))
	if err != nil {
		return nil, err
	}

	defer conn.Close()

	err = json.NewEncoder(conn).Encode(dataEnc)

	if err != nil {
		log.Printf("Send Error: %v\n", err)
	}

	//log.Printf("Packet: %v\n", string(apiPacketStr))

	resp, err := readAll(conn)
	if err != nil {
		return nil, err
	}
	log.Printf("Reponse Code: %v\n", string(resp))
	// Modify response to handle weird M31 issues
	resp = strings.ReplaceAll(resp, "inf", "999")
	resp = strings.ReplaceAll(resp, "nan", "0")

	var result map[string]interface{}
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		log.Println("Error decoding encrypted response")
		log.Println(resp)
		return nil, err
	}

	if result["STATUS"] != nil {
		if result["STATUS"].(string) == "E" {
			// Error Raising Here ********************************************************************
			// Returns Code: 23 -> Invalid JSON Format
			return nil, errors.New(result["Msg"].(string))
		}
	}

	respCiphertext, err := base64.StdEncoding.DecodeString(result["enc"].(string))
	if err != nil {
		return nil, err
	}

	respPlaintext := decrypt(string(respCiphertext), accessToken.Cipher)
	resp = strings.Split(respPlaintext, "\x00")[0]
	err = json.Unmarshal([]byte(resp), &result)
	if err != nil {
		return nil, err
	}
	log.Printf("Response: %v\n", resp)
	return result, nil
}

func decrypt(cipherstring string, block cipher.Block) string {
	// Byte array of the string
	ciphertext := []byte(cipherstring)

	// Key
	//key := []byte(keystring)

	// Create the AES cipher
	//block, err := aes.NewCipher(key)
	//if err != nil {
	//	panic(err)
	//}

	// Before even testing the decryption,
	// if the text is too small, then it is incorrect
	if len(ciphertext) < aes.BlockSize {
		panic("Text is too short")
	}

	mode := ecb.NewECBDecrypter(block)
	mode.CryptBlocks(ciphertext, ciphertext)
	ciphertext = PKCS5UnPadding(ciphertext)

	return string(ciphertext)
}

// readAll reads all available data from a connection.
func readAll(conn net.Conn) (string, error) {
	var result strings.Builder
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return "", err
		}
		result.Write(buffer[:n])
		if n < len(buffer) {
			break
		}
	}
	return result.String(), nil
}

// addTo16 pads a string to a multiple of 16 bytes.
func addTo16(b []byte) []byte {
	bytes := make([]byte, len(b)+(16-len(b)%16))

	for i := 0; i < len(bytes); i++ {
		if i < len(b) {
			bytes[i] = b[i]
		} else {
			bytes[i] = 0
		}
	}

	return bytes
}

// PKCS5UnPadding  pads a certain blob of data with necessary data to be used in AES block cipher
func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])

	return src[:(length - unpadding)]
}
