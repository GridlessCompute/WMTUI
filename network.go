package main

import (
	api "WMTUI/wmapi"
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

var wapi api.WhatsminerAPI

type Miner struct {
	Ip          string
	Mac         string
	Errcode     string
	UpTime      int
	Hrrt        int
	Wt          int
	W           int
	Limit       int
	Fastboot    string
	Sleep       string
	AcvtivePool string
}

type MinerObj struct {
	Miner   Miner
	Token   api.WhatsminerAccessToken
	Created time.Time
}

func generateAddress(ip string, port int) string {
	return fmt.Sprintf("%s:%d", ip, port)

}

// func generateIp(root string, o4 int) string {
// 	return fmt.Sprintf("%s.%d", root, o4)

// }

func aquireToken(ip string, port int, flags string) (api.WhatsminerAccessToken, error) {
	token, tokenErr := api.NewWhatsminerAccessToken(ip, port, flags)

	if tokenErr != nil {
		// fmt.Printf("unable to generate token for %s", ip)
		return *new(api.WhatsminerAccessToken), tokenErr
	}
	return *token, nil
}

func getFromApi(token api.WhatsminerAccessToken, cmd string) (map[string]interface{}, error) {
	summary, err := wapi.GetReadOnlyInfo(&token, cmd, nil)

	if err != nil {
		return nil, err
	}
	return summary, nil
}

func SendToApi(token api.WhatsminerAccessToken, cmd string, params map[string]interface{}) (map[string]interface{}, error) {
	res, err := wapi.ExecCommand(&token, cmd, params)
	if err != nil {
		fmt.Println(err)
	}
	return res, err
}

func parseSummary(summary map[string]interface{}) (api.SummaryS, api.ApiError, error) {
	var smry api.SummaryS
	var smryErr api.ApiError
	var smryErrJ error

	s, sErr := json.Marshal(summary)

	if sErr != nil {
		// fmt.Println("SOMETHING BAD IN PARSESUMMARY")
	}

	smryj := json.Unmarshal(s, &smry)
	if smryj != nil {
		smryErrJ = json.Unmarshal(s, &smryErr)
	}
	if smryErrJ != nil {
		// fmt.Print("Unable to get info from api")
	}
	//fmt.Print(smry)
	return smry, smryErr, nil
}

func InitScanOne(ip string, hashChan chan MinerObj, wg *sync.WaitGroup) {
	// mnr := new(Miner)
	defer wg.Done()
	mnrO := new(MinerObj)
	mnr := new(Miner)
	address := generateAddress(ip, 4028)
	// Initial connection
	conn, err := net.DialTimeout("tcp", address, 3*time.Second)
	if err2 := recover(); err2 != nil || err != nil {
		// 	neterr, ok := err.(net.Error)
		// 	if ok && neterr.Timeout() {
		// 		logger.Printf("Connection timeout for %s\n", ip)
		// 	} else {
		// 		logger.Printf("Connection Refused: %s\n", ip)
		// 	}
	} else {
		// get token if able to connect
		token, tokenErr := aquireToken(ip, 4028, "admin")
		if tokenErr != nil {
			// fmt.Println(tokenErr.Error())
		}

		//Get Miner Info
		info, err := getFromApi(token, "get_miner_info")
		if err != nil {
			// fmt.Println(err.Error())
		}

		i, err := json.Marshal(info)
		if err != nil {
			// fmt.Println(err.Error())
		}

		var infoStruct api.GetMinerInfoS
		var infoStruct2 api.GetMinerInfo2S
		err = json.Unmarshal(i, &infoStruct)
		if err != nil {
			_ = json.Unmarshal(i, &infoStruct2)
			mnr.Mac = infoStruct2.Msg.Mac
		} else {
			mnr.Mac = infoStruct.Msg.Mac
		}

		mnr.Ip = ip

		mnrO.Created = time.Now()
		mnrO.Token = token
		mnrO.Miner = *mnr

		hashChan <- *mnrO

		conn.Close()

	}

}

func GetMinerData(wg *sync.WaitGroup, mnrO MinerObj, hashChannel chan MinerObj) {
	defer wg.Done()
	//Get Miner Info
	info, err := getFromApi(mnrO.Token, "get_miner_info")
	if err != nil {
		fmt.Println(err.Error())
	}

	i, err := json.Marshal(info)
	if err != nil {
		fmt.Println(err.Error())
	}

	var infoStruct api.GetMinerInfoS
	var infoStruct2 api.GetMinerInfo2S
	err = json.Unmarshal(i, &infoStruct)
	if err != nil {
		err = json.Unmarshal(i, &infoStruct2)
	}

	// // Get Pool Info
	// pools, err := getFromApi(mnrO.Token, "pools")
	// if err != nil {
	// 	// fmt.Println(err.Error())
	// }

	// p, err := json.Marshal(pools)
	// if err != nil {
	// 	// fmt.Println(err.Error())
	// }
	// var poolStruct api.GetPoolInfoS
	// err = json.Unmarshal(p, &poolStruct)
	// if err != nil {
	// 	// fmt.Println(err.Error())
	// }

	// Get Error Code
	er, erErr := getFromApi(mnrO.Token, "get_error_code")
	if erErr != nil {
		fmt.Println(erErr.Error())
	}

	e, eErr := json.Marshal(er)
	if eErr != nil {
		fmt.Println(eErr.Error())
	}

	var err1 api.GetErrorCodeS
	jsonErr := json.Unmarshal(e, &err1)
	if jsonErr != nil {
		fmt.Println(jsonErr.Error())
	}

	re := regexp.MustCompile(`(?:\"|\')(?P<key>[0-9]+)(?:\"|\')`)

	matched := re.FindAllString(string(e), -1)
	matchedString := strings.Join(matched, ", ")

	// Get Summary Info
	summary, summaryErr := getFromApi(mnrO.Token, "summary")
	if summaryErr != nil {
		fmt.Println(summaryErr.Error())
	}

	// var ap string
	// if len(poolStruct.Pools) > 0 {
	// 	ap = poolStruct.Pools[0].URL
	// } else {
	// 	ap = ""
	// }

	res, _, parseErr := parseSummary(summary)
	if parseErr != nil {

	} else {
		if len(res.Summary) > 0 {
			mnrO.Miner.AcvtivePool = "pool"
			mnrO.Miner.Hrrt = int(res.Summary[0].HSRT)
			mnrO.Miner.Limit = res.Summary[0].PowerLimit
			mnrO.Miner.UpTime = res.Summary[0].Uptime
			mnrO.Miner.W = res.Summary[0].Power
			mnrO.Miner.Wt = int(res.Summary[0].PowerRate)
			mnrO.Miner.Errcode = matchedString
			mnrO.Miner.Fastboot = res.Summary[0].BtminerFastBoot

			hashChannel <- mnrO
		}
	}

	//fmt.Println("done updating miner")

}
