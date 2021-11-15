package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hdevillers/go-seq/seqio"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	// Retrieve argument values
	in1 := flag.String("in1", "", "Input (#1) fastq file.")
	in2 := flag.String("in2", "", "Input (#2) fastq file (optional).")
	out1 := flag.String("out1", "", "Output (#1) fastq file.")
	out2 := flag.String("out2", "", "Output (#2) fastq file (required if -in2 provided).")
	nsam := flag.Int("n", 0, "Number of reads to keep (n-first).")
	psam := flag.Float64("p", 0, "Proportion of reads to keep [0; 1].")
	gunzip := flag.Bool("gz", false, "Input and output are compressed (gz).")
	flag.Parse()

	// Prepare seqio objects
	if *in1 == "" {
		if *in2 == "" {
			*in1 = "STDIN"
		} else {
			panic("Can use SDTIN when working with two files. You must provide an input fastq file #1 (-in1).")
		}
	}
	if *out1 == "" {
		if *out2 == "" {
			*out2 = "STDOUT"
		} else {
			panic("Can print outputs to STDOUT when working with two files. You must provide an output file #1 (-out2).")
		}
	}
	if *in2 != "" {
		if *out2 == "" {
			panic("You must provide an output file name for the second file (-out2).")
		}
	}
	if *nsam == 0 && *psam == 0.0 {
		panic("You must choose one selection strategy (-n: n first reads or -p: proportion of reads).")
	}
	if *nsam != 0 && *psam != 0.0 {
		panic("You cannot choose both selection strategies (-p and -n).")
	}

	// Select reads according the strategy
	if *in2 == "" {
		seqin1 := seqio.NewReader(*in1, "fastq", *gunzip)
		seqin1.CheckPanic()
		defer seqin1.Close()

		seqout1 := seqio.NewWriter(*out1, "fastq", *gunzip)
		seqout1.CheckPanic()
		defer seqout1.Close()

		if *nsam != 0 {
			nstored := 0
		SEQ:
			for seqin1.Next() {
				seqin1.CheckPanic()
				nstored++
				seqout1.Write(seqin1.Seq())
				if nstored == *nsam {
					break SEQ
				}
			}
			if *nsam > nstored {
				fmt.Fprintf(os.Stderr, "Required %d reads but input file contains only %d reads.", *nsam, nstored)
			}
		}
	}

}
