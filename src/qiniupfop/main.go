package main

import (
	"flag"
	"log"
	"net/url"
)

import (
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/rs"
)

type pfopRet struct {
	PersistentId string `json:"persistentId"`
	Code         int    `json:"code"`
	Error        string `json:"error"`
}

func main() {
	bucket := flag.String("b", "", "bucket")
	key := flag.String("k", "", "file key")
	accessKey := flag.String("ak", "", "access key")
	secretKey := flag.String("sk", "", "secret key")
	pfop := flag.String("pfop", "", "pfop option")
	notify := flag.String("n", "", "notify url")
	flag.Parse()
	if *bucket == "" || *key == "" || *accessKey == "" || *accessKey == "" || *pfop == "" {
		flag.PrintDefaults()
		log.Fatalln("invalid args")
		return
	}
	conf.ACCESS_KEY = *accessKey
	conf.SECRET_KEY = *secretKey

	client := rs.New(nil)
	param := url.Values{}
	param.Set("bucket", *bucket)
	param.Set("key", url.QueryEscape(*key))
	param.Set("fops", url.QueryEscape(*pfop))
	param.Set("notifyURL", url.QueryEscape(*notify))
	var ret pfopRet
	err := client.Conn.CallWithForm(nil, &ret, "http://api.qiniu.com/pfop/", param)
	if err != nil {
		log.Fatalln("error", err, ret)
	}
	log.Println("id", ret.PersistentId)
}
