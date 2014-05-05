package main

import (
	"encoding/base64"
	"flag"
	"log"
	"net/url"
	"os"
)

import (
	"github.com/qiniu/api/conf"
	"github.com/qiniu/api/resumable/io"
	"github.com/qiniu/api/rs"
)

type pfopRet struct {
	PersistentId string `json:"persistentId"`
	Code         int    `json:"code"`
	Error        string `json:"error"`
}

func pfop(pfops, bucket, key, notifyUrl string) (ret pfopRet, err error) {
	client := rs.New(nil)
	param := url.Values{}
	param.Set("bucket", bucket)
	param.Set("key", key)
	param.Set("fops", pfops)
	param.Set("notifyURL", notifyUrl)
	err = client.Conn.CallWithForm(nil, &ret, "http://api.qiniu.com/pfop/", param)
	return
}

func buildOps(convert, saveas string) string {
	pfops := convert
	if saveas != "" {
		pfops = pfops + "|saveas/" + base64.URLEncoding.EncodeToString([]byte(saveas))
	}
	return pfops
}

func genToken(bucket, fops, notifyUrl string) string {
	policy := rs.PutPolicy{
		Scope:               bucket,
		PersistentNotifyUrl: notifyUrl,
		PersistentOps:       fops,
	}
	return policy.Token(nil)
}

func put(token, key, file string) (ret io.PutRet, err error) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalln("file not exist")
		return
	}
	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		log.Fatalln("invalid file")
		return
	}

	blockNotify := func(blkIdx int, blkSize int, ret *io.BlkputRet) {
		log.Println("size", stat.Size(), "block id", blkIdx, "offset", ret.Offset)
	}

	params := map[string]string{}
	extra := &io.PutExtra{
		ChunkSize: 8192,
		Notify:    blockNotify,
		Params:    params,
	}

	err = io.PutFile(nil, &ret, token, key, file, extra)
	return
}

func main() {
	file := flag.String("f", "", "local file")
	bucket := flag.String("b", "", "bucket")
	key := flag.String("k", "", "file key")
	accessKey := flag.String("ak", "", "access key")
	secretKey := flag.String("sk", "", "secret key")
	convert := flag.String("c", "", "pfop single convert option avthumb/m3u8/segtime/10/vcodec/libx264/s/320x240")
	notify := flag.String("n", "", "notify url")
	saveas := flag.String("o", "", "output key, outputBucket:key")
	flag.Parse()
	if *bucket == "" || *key == "" || *accessKey == "" || *accessKey == "" || *convert == "" {
		flag.PrintDefaults()
		log.Fatalln("invalid args")
		return
	}
	conf.ACCESS_KEY = *accessKey
	conf.SECRET_KEY = *secretKey
	ops := buildOps(*convert, *saveas)
	if *file != "" {
		token := genToken(*bucket, ops, *notify)
		ret, err := put(token, *key, *file)
		if err != nil {
			log.Fatalln("put error", err, ret)
			return
		}
		log.Println("put ret", ret.Key)
		return
	}
	ret, err := pfop(ops, *bucket, *key, *notify)
	if err != nil {
		log.Fatalln("pfop error", err, ret)
		return
	}
	log.Println("pfop id", ret.PersistentId)
}
