package testDataMaker

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func CreateMountPointAndClient(mpn int, n int) map[string][]string {
	mountpoints := make(map[string][]string, mpn)
	for i := 0; i < mpn; i++ {
		clientList := []string{}
		// fmt.Println(clientList)
		for j := 0; j < n; j++ {
			clientName := fmt.Sprintf("mp%04dclient%04d", i+1, j+1)
			clientList = append(clientList, clientName)
		}
		mpName := fmt.Sprintf("mp%04d", i+1)
		mountpoints[mpName] = clientList
		// fmt.Println(mountpoints)
	}
	return mountpoints
}

func KeyMaker(n int) string {
	d1md5 := md5.Sum([]byte(fmt.Sprintf("%d", n)))
	d1md5str1 := fmt.Sprintf("%x%x%x", d1md5, d1md5, d1md5)
	// fmt.Println(d1md5str1)
	d1Base64 := base64.StdEncoding.EncodeToString([]byte(d1md5str1))
	d2Base64 := base64.StdEncoding.EncodeToString([]byte(d1Base64))
	d3Base64 := base64.StdEncoding.EncodeToString([]byte(d2Base64))
	return d3Base64
}

func DataMaker(n int) []byte {
	a := make(map[int]string)
	a[n] = KeyMaker(n)
	data, _ := json.Marshal(a)
	return data
}
