```
go build -ldflags "-s -w" cabaymau.go && ./cabaymau --file /home/chiro/gits/playground-go/data/test.txt --chunk [4 times of your number of CPU threads]
```

Example result:

`FreeBSD-13.0-RELEASE-amd64-dvd1.iso` 4,5 GBs

```
chiro@Luna:~/gits/examples/md4-rainbow$ ./cabaymau --file /home/chiro/Downloads/FreeBSD-13.0-RELEASE-amd64-dvd1.iso --chunk 64
2021/06/14 13:09:13 Calculation took 148.364859 sec
```