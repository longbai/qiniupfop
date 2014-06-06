package main

import (
	"encoding/base64"
	"flag"
	"fmt"
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
	if notifyUrl != "" {
		param.Set("notifyURL", notifyUrl)
	}
	err = client.Conn.CallWithForm(nil, &ret, "http://api.qiniu.com/pfop", param)
	return
}

func buildOps(convert, saveas string) string {
	pfops := convert
	if saveas != "" {
		pfops = pfops + "|saveas/" + base64.URLEncoding.EncodeToString([]byte(saveas))
	}
	return pfops
}

func genToken(bucket, fops, notifyUrl, pipeline string) string {
	policy := rs.PutPolicy{
		Scope:         bucket,
		PersistentOps: fops,

		ReturnBody: `{"persistentId": $(persistentId)}`,
	}
	if notifyUrl != "" {
		policy.PersistentNotifyUrl = notifyUrl
	}
	if pipeline != "" {
		policy.PersistentPipeline = pipeline
	}

	return policy.Token(nil)
}

func put(token, key, file string) (ret pfopRet, err error) {
	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "file not exist")
		os.Exit(1)
		return
	}
	stat, err := f.Stat()
	if err != nil || stat.IsDir() {
		fmt.Fprintln(os.Stderr, "invalid file")
		os.Exit(1)
		return
	}

	// blockNotify := func(blkIdx int, blkSize int, ret *io.BlkputRet) {
	// 	fmt.Println("size", stat.Size(), "block id", blkIdx, "offset", ret.Offset)
	// }

	params := map[string]string{}
	extra := &io.PutExtra{
		ChunkSize: 8192,
		// Notify:    blockNotify,
		Params: params,
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
	pipeline := flag.String("p", "", "pipline")
	saveas := flag.String("o", "", "output key, outputBucket:key")
	flag.Parse()
	if *bucket == "" || *key == "" || *accessKey == "" || *accessKey == "" || *convert == "" {
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "invalid args")
		os.Exit(1)
		return
	}

	conf.ACCESS_KEY = *accessKey
	conf.SECRET_KEY = *secretKey
	ops := buildOps(*convert, *saveas)
	if *file != "" {
		token := genToken(*bucket, ops, *notify, *pipeline)
		ret, err := put(token, *key, *file)
		if err != nil {
			fmt.Fprintln(os.Stderr, "put error", err, ret)
			os.Exit(1)
			return
		}
		fmt.Fprintln(os.Stdout, ret.PersistentId)
		return
	}
	ret, err := pfop(ops, *bucket, *key, *notify)
	if err != nil {
		fmt.Fprintln(os.Stderr, "pfop error", err, ret)
		os.Exit(1)
		return
	}
	fmt.Fprintln(os.Stdout, ret.PersistentId)
}
