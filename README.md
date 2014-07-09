qiniupfop
=========

pfop go example

## pfop with exist key
```
bin/qiniupfop -b "<bucket>" -k "<key>" -ak "<access key>" -sk "<secret key>" -n "<notify url>" -c "<convert>" -o "[saveas]" -p "[pipeline]"

example:
bin/qiniupfop -b shun -k sintel_trailer.mp4 -ak "****" -sk "****" -n "http://uoccz3glhjg7.runscope.net/" -c "avthumb/m3u8/segtime/10/vcodec/libx264/s/320x240" -o "shun:test1.m3u8" -f ~/sintel_trailer.mp4

```
##pfop after upload
```
bin/qiniupfop -b "<bucket>" -k "<key>" -ak "<access key>" -sk "<secret key>" -n "<notify url>" -c "<convert>" -o "[saveas]" -f "[file]" -p "[pipeline]"

example:
bin/qiniupfop -b shun -k sintel_trailer.mp4 -ak "****" -sk "****" -n "http://uoccz3glhjg7.runscope.net/" -c "avthumb/m3u8/segtime/10/vcodec/libx264/s/320x240" -o "shun:test1.m3u8"
```
