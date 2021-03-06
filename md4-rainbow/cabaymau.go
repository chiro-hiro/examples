package main

import (
	"bufio"
	"bytes"
	"encoding/hex"
	"flag"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"golang.org/x/crypto/md4"
)

type Job struct {
	Start int64
	End   int64
}

type Jobs []Job

type Computa struct {
	Salt     []byte
	FileName string
	Chunk    int
	fHandle  *os.File
}

func (c *Computa) Clone() Computa {
	return Computa{Salt: c.Salt, FileName: c.FileName, Chunk: c.Chunk}
}

func (c *Computa) Open() {
	f, eErr := os.OpenFile(c.FileName, os.O_RDONLY, 0755)

	if eErr != nil {
		log.Fatalln(eErr)
	}

	c.fHandle = f

}

func (c *Computa) CalcuateJobs() Jobs {
	f, _ := c.fHandle.Stat()
	if c.Chunk <= 1 {
		return Jobs{Job{Start: 0, End: f.Size()}}
	}
	chunkSize := f.Size() / int64(c.Chunk)
	js := make(Jobs, c.Chunk)
	for i, k := int64(0), 0; k < c.Chunk; k += 1 {
		s := int64(i + chunkSize)
		if s >= f.Size() {
			s = f.Size()
		} else {
			// Seeking for next end of line
			buf := make([]byte, chunkSize/10)
			_, e := c.fHandle.ReadAt(buf, s)
			if e != nil {
				s = f.Size()
			} else {
				s += int64(bytes.Index(buf, []byte{'\n'}))
			}
		}
		j := Job{Start: i, End: s}
		js[k] = j
		i = s + 1
	}
	return js
}

func (c *Computa) CalculateMD4(j Job, cH chan string) {
	if c.fHandle == nil {
		log.Fatalln("File wasn't open, you know.")
	}
	h := md4.New()
	r := bufio.NewReader(c.fHandle)
	c.fHandle.Seek(j.Start, io.SeekStart)
	for k := j.Start; k < j.End; {
		h.Reset()
		b, _, e := r.ReadLine()
		if e != nil {
			break
		}
		t := bytes.Trim(b, "\n\r")
		h.Write(t)
		cH <- hex.EncodeToString(h.Sum(nil)) + "," + string(t) + "\n"
		k += int64(len(b))
	}
}

func (c *Computa) Close() {
	if err := c.fHandle.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	sFilename := flag.String("file", "", "Input filename")
	sOut := flag.String("out", "./a.out", "Onput filename")
	sSalt := flag.String("salt", "", "Extra salt in hexstring")
	nChunk := flag.Int("chunk", 2, "Number of chunks, it used to your CPU cores")
	start := time.Now()
	flag.Parse()
	ctx := Computa{}
	if *sFilename == "" {
		flag.Usage()
		println("ERROR: File name can't be empty")
		return
	}

	if len(*sSalt) > 0 {
		bSalt, eErr := hex.DecodeString(*sSalt)
		if eErr != nil {
			log.Fatal(eErr)
		}
		ctx.Salt = bSalt
	}

	ctx.FileName = *sFilename
	ctx.Chunk = *nChunk

	ctx.Open()
	jobs := ctx.CalcuateJobs()
	defer ctx.Close()

	result := make(chan string, 1)

	go func(result chan string) {
		fH, _ := os.OpenFile(*sOut, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		defer fH.Close()
		for {
			c := <-result
			if c == "" {
				break
			}
			fH.WriteString(c)
		}
	}(result)

	var wg sync.WaitGroup

	for i := 0; i < len(jobs); i += 1 {
		k := ctx.Clone()
		k.Open()
		wg.Add(1)
		go func(i int, wg *sync.WaitGroup) {
			defer wg.Done()
			k.CalculateMD4(jobs[i], result)
		}(i, &wg)
		defer k.Close()
	}
	wg.Wait()
	result <- ""
	log.Printf("Calculation took %f sec", time.Since(start).Seconds())
}
