package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/cespare/xxhash"
	"github.com/shenwei356/bio/seq"
	"github.com/shenwei356/bio/seqio/fastx"
	"github.com/shenwei356/xopen"
	"io"
	"os"
)

func main(){
	var alphabet *seq.Alphabet
	idRegexp := "^(\\S+)\\s?"

	LineWidth := 0

	bySeq := true
	quiet := false



	//ForcelyOutputFastq := true
	//files := make([]string,0,100)
	//files = append(files,"R1.fq.gz")
	var infileR1,infileR2,outfileR1,outfileR2 string
	var h,ignoreCase bool
	flag.StringVar(&infileR1,"in1","raw_R1.pe.fastq","input R1 file,support gz or fq")
	flag.StringVar(&infileR2,"in2","raw_R2.pe.fastq","input R2 file,support gz or fq")
	flag.BoolVar(&h, "h", false, "this script is change from seqkit rmdup -by_seq, same algorithm but input R1 R2,output R1 R2\n")
	flag.BoolVar(&ignoreCase, "ignoreCase", false, "ignore case")
	flag.StringVar(&outfileR1,"out1","uniq_R1.fastq","output R1 file,support gz or fq")
	flag.StringVar(&outfileR2,"out2","uniq_R2.fastq","output R2 file,support gz or fq")

	flag.Parse()
	if h{flag.Usage()
	os.Exit(1)
	}
	outfhR1, err := xopen.Wopen(outfileR1)
	if err != nil{
		panic("创建输出文件R1失败")
	}

	outfhR2, err := xopen.Wopen(outfileR2)
	if err != nil{
		panic("创建输出文件R2失败")
	}
	defer outfhR1.Close()
	defer outfhR2.Close()
	//var outfhDup *xopen.Writer
	counterR1 := make(map[uint64]int)
	counterR2 := make(map[uint64]int)
	//names := make(map[uint64][]string)
	var subjectR1 uint64
	var subjectR2 uint64
	var removed int
	var recordR1 *fastx.Record
	var recordR2 *fastx.Record
	var fastxReaderR1 *fastx.Reader
	var fastxReaderR2 *fastx.Reader


		fastxReaderR1, err = fastx.NewReader(alphabet, infileR1, idRegexp)
		fastxReaderR2, err = fastx.NewReader(alphabet, infileR2, idRegexp)

		if err != nil{panic(err)}
		for {
			recordR1, err = fastxReaderR1.Read()
			if err != nil {
				if err == io.EOF {
					break
				}

				break
			}
			if fastxReaderR1.IsFastq {
				LineWidth = 0
				fastx.ForcelyOutputFastq = true
			}
			//R1
			recordR2, err = fastxReaderR2.Read()
			if err != nil {
				if err == io.EOF {
					break
				}

				break
			}

			//R2
			if bySeq {
				if ignoreCase {
					subjectR1 = xxhash.Sum64(bytes.ToLower(recordR1.Seq.Seq))
					subjectR2 = xxhash.Sum64(bytes.ToLower(recordR2.Seq.Seq))
				} else {
					subjectR1 = xxhash.Sum64(recordR1.Seq.Seq)
					subjectR2 = xxhash.Sum64(recordR2.Seq.Seq)
				}
			}

			_, ok := counterR1[subjectR1]
			_,ok2 := counterR2[subjectR2]
			if ok || ok2{ // duplicated
				counterR1[subjectR1]++
				counterR2[subjectR2]++
				removed++


			} else { // new one
				recordR1.FormatToWriter(outfhR1, LineWidth)
				recordR2.FormatToWriter(outfhR2, LineWidth)
				counterR1[subjectR1]++
				counterR2[subjectR2]++
			}
		}



	if !quiet {
		fmt.Printf("%d duplicated records(PE) removed\n", removed)
	}

}
