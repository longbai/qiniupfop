package main

import (
	"flag"
	"log"
	"net/url"
	"os"
)

import (
	"github.com/qiniu/api/resumable/io"
)

func main() {
	file := flag.String("f", "", "file path")
	accessKey := flag.String("ak", "", "access key")
	secretKey := flag.String("sk", "", "secret key")
	pfop := flag.String("pfop", "", "pfop option")
	notify := flag.String("u", "", "notify url")
	custom := flag.String("x", "", "custom args foo=1&bar=2")
	flag.Parse()
	if *token == "" || *file == "" || *key == "" {
		flag.PrintDefaults()
		log.Fatalln("invalid args")
		return
	}

	f, err := os.Open(*file)
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
	if custom != nil && *custom != "" {
		values, err := url.ParseQuery(*custom)
		if err != nil {
			log.Fatalln(err.Error())
			return
		}
		for k, v := range values {
			params["x:"+k] = v[0]
		}
		extra.Params = params
	}
	var ret io.PutRet
	err = io.PutFile(nil, &ret, *token, *key, *file, extra)
	if err != nil {
		log.Fatalln(err)
	}
}
