package main

import (
	"flag"
	"fmt"

	"github.com/hdevillers/go-seq/seq"
	"github.com/hdevillers/go-seq/seqio"
)

func main() {
	// Retrieve argument values
	input := flag.String("input", "STDIN", "Input fasta file")
	format := flag.String("format", "fasta", "Input format.")
	gunzip := flag.Bool("d", false, "Decompress the input (gz).")
	mode := flag.String("mode", "overall", "GC content computation mode (overall, onebyony, windows, index).")
	winLen := flag.Int("win-length", 500, "Define window length when mode is 'windows'.")
	winSte := flag.Int("win-step", 250, "Define window step when mode is 'windows'.")
	indNum := flag.Int("index-num", -1, "Define sequence index by number when mode is 'index'.")
	indId := flag.String("index-id", "", "Define sequence index by sequence ID when mode is 'index'.")
	indFrom := flag.Int("index-from", -1, "Define sequence start location when mode is 'index'.")
	indTo := flag.Int("index-to", -1, "Define sequence end location when mode is 'index'.")
	flag.Parse()

	if *input == "" {
		panic("You must provide an input fasta file.")
	}

	// Start reading the input sequence
	seqIn := seqio.NewReader(*input, *format, *gunzip)
	seqIn.CheckPanic()
	defer seqIn.Close()

	// Initiate a GC object
	gc := seq.NewGc()

	switch *mode {
	case "overall":
		// Compute overall GC content
		var ngc int64 = 0
		var ntt int64 = 0
		for seqIn.Next() {
			seqIn.CheckPanic()
			seq := seqIn.Seq()
			ngci, skpi, _ := gc.Count(&seq)
			ngc += ngci
			ntt += int64(seq.Length()) - skpi
		}
		gcc := float64(ngc) / float64(ntt) * 100.0
		fmt.Printf("%s\t%.04f\n", *input, gcc)
	case "onebyone":
		// Compute GC content sequence by sequence
		for seqIn.Next() {
			seqIn.CheckPanic()
			seq := seqIn.Seq()
			ngc, skp, _ := gc.Count(&seq)
			// TODO: control that skp is smaller than seq.Length()
			gcc := float64(ngc) / float64(int64(seq.Length())-skp) * 100.0
			fmt.Printf("%s\t%.04f\n", seq.Id, gcc)
		}
	case "windows":
		// Compute GC content using a sliding window algorithm
		// Check window parameters
		if *winLen < 10 {
			panic("Window length must be at least 10 nucleotides.")
		}
		if *winSte < 1 {
			panic("Window step must be greater or equal to 1.")
		}
		for seqIn.Next() {
			seqIn.CheckPanic()
			seq := seqIn.Seq()
			for i := 0; i < (seq.Length() - *winLen + 1); i += *winSte {
				ngc, skp, err := gc.CountIndex(&seq, i+1, i+*winLen)
				if err != nil {
					panic(err)
				}
				if int(skp) < *winLen {
					gcc := float64(ngc) / float64(int64(*winLen)-skp) * 100.0
					fmt.Printf("%s.%d\t%.04f\n", seq.Id, i+1, gcc)
				} // Just skip reporting to avoid division by zero
			}
			// NOTE: Sequences shorter then window length are just skipped.
		}
	case "index":
		// Compute GC content of a precise region using index
		// Selection by index number
		ok := false
		var seq seq.Seq
		if *indNum != -1 {
			if *indId != "" {
				panic("You cannot use both an index number and a sequence ID when using index mode.")
			}
			if *indNum < 1 {
				panic("Sequence index must be greater or equal to 1.")
			}
			// Check if managed to find the sequence
			ind := 0
		SEARCH_IND:
			for seqIn.Next() {
				seqIn.CheckPanic()
				ind++
				if ind == *indNum {
					seq = seqIn.Seq()
					ok = true
					break SEARCH_IND
				}
			}
		} else if *indId != "" {
		SEARCH_ID:
			for seqIn.Next() {
				seqIn.CheckPanic()
				tmp := seqIn.Seq()
				if tmp.Id == *indId {
					seq = tmp
					ok = true
					break SEARCH_ID
				}
			}
		} else {
			panic("You must provide either a sequence index number or a sequence ID when using index mode.")
		}

		if ok {
			// Check if precise location are required
			if *indFrom == -1 && *indTo == -1 {
				// Compute the GC content for the whole sequence
				ngc, skp, _ := gc.Count(&seq)
				gcc := float64(ngc) / float64(int64(seq.Length())-skp) * 100.0
				fmt.Printf("%s\t%.04f\n", seq.Id, gcc)
			} else {
				if *indFrom < 1 {
					panic("Sequence start location must be at least 1 when using index mode.")
				}
				if *indTo < *indFrom {
					panic("Sequence end location must be higher than sequence start location.")
				}
				if *indTo > seq.Length() {
					panic("Sequence end location is out of bound.")
				}
				ngc, skp, err := gc.CountIndex(&seq, *indFrom, *indTo)
				if err != nil {
					panic(err)
				}
				gcc := float64(ngc) / float64(int64(*indTo-*indFrom+1)-skp) * 100.0
				fmt.Printf("%s.%dto%d\t%.04f\n", seq.Id, *indFrom, *indTo, gcc)
			}
		} else {
			panic("Failed to find the required sequence, please check you argument value(s).")
		}
	default:
		panic(fmt.Sprintf("The provided mode (%s) is not supported, please use 'overall', 'onebyone', 'windows' or 'index'.", *mode))
	}
}
