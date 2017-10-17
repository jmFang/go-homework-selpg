/**
/* anthor :jiamoufang
/* this program is achieved with reference of 'selpg.c'
/* something special : use flag package to parse cmi instruction
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

/*
* a flag definition of selpg
 */

var (
	selpg = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	s     = selpg.Int("s", 0, "the start page required in command selpg")
	e     = selpg.Int("e", 0, "the end page required in command selpg")
	l     = selpg.Int("l", 72, "the initial number of line per page")
	d     = selpg.String("d", "", "the destination of printer")
	f     = selpg.Bool("f", false, "the flag of new page noted as '\f'")
)

type structSelpg struct {
	startPage  int
	endPage    int
	inFilename string
	pageLen    int
	pageType   int
	printDest  string
}

// the instruction of "selpg"
var progname string

func usage() {
	fmt.Printf("\nUSAGE: %s -sstart_page -eend_page [-f] [-llines_per_page ] [-ddest ] [ in_filename ]\n", progname)
}

func processArgs(ac int, av []string, psa *structSelpg) {

	selpg.Parsed()

	var tmp = []rune(av[0])
	var argno int

	if ac < 3 {
		fmt.Fprintf(os.Stderr, "%s: not enough arguments\n", progname)
		usage()
		os.Exit(1)
	}

	if av[1][0:2] != "-s" {
		fmt.Fprintf(os.Stderr, "%s: 1st arg should be -sstartPage\n", progname)
		usage()
		os.Exit(2)
	}

	if av[2][0:2] != "-e" {
		fmt.Printf("%s: 2nd arg should -eendPage\n", progname)
		usage()
		os.Exit(4)
	}

	start, _ := strconv.Atoi(string(av[1][2:]))
	if start < 1 {
		fmt.Printf("%s : invalid start page %d\n", progname, start)
		usage()
		os.Exit(3)
	}

	psa.startPage = start

	end, _ := strconv.Atoi(string(av[2][2:]))
	if end < 1 || end < psa.startPage {
		fmt.Printf("%s: invalid end page %d\n", progname, end)
		usage()
		os.Exit(5)
	}

	psa.endPage = end

	argno = 3

	for argno <= ac-1 && []rune(av[argno])[0] == '-' {
		tmp = []rune(av[argno])

		switch tmp[1] {

		case 'l':

			k, err := strconv.Atoi(string(tmp[2:]))
			if k < 1 || err != nil {
				fmt.Fprintf(os.Stderr, "%s: invalid page length %d\n", progname, k)
				usage()
				os.Exit(6)
			}

			psa.pageLen = k
			argno++

		case 'f':
			if strings.Compare(string(tmp), "-f") != 0 {
				fmt.Fprintf(os.Stderr, "%s: option should be \"-f\"\n", progname)
				usage()
				os.Exit(7)
			}

			psa.pageType = 'f'
			argno++

		case 'd':
			t := tmp[2:]
			if len(t) < 1 {
				fmt.Fprintf(os.Stderr, "%s: -d option requires a printer destination", progname)
				usage()
				os.Exit(8)
			}

			psa.printDest = string(tmp[2:])
			argno++

		default:
			fmt.Fprintf(os.Stderr, "%s: unknown option %s\n", progname, tmp)
			usage()
			os.Exit(9)

		} /*end switch*/
	} // end while

	if argno <= ac-1 {
		psa.inFilename = av[argno]
		//try to open the file
		if f, err := os.Open(psa.inFilename); err != nil {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" does not exist\n", progname, psa.inFilename)
			f.Close()
			os.Exit(10)
		}
		// weather the file can be read
		if f, err := os.OpenFile(psa.inFilename, os.O_RDWR|os.O_CREATE, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "%s: input file \"%s\" cannot be read\n", progname, psa.inFilename)
			f.Close()
			os.Exit(11)
		}

	} // end if
}

func processInput(psa structSelpg) {
	var fin *os.File
	var fout *os.File
	var err error
	var cmd *exec.Cmd

	if psa.inFilename == "" {
		fin = os.Stdin
	} else {
		fin, err = os.Open(psa.inFilename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open input file \"%s\"\n", progname, psa.inFilename)
			os.Exit(12)
		}
	} // end else

	/* set the output destination */
	if psa.printDest == "" {
		fout = os.Stdout
	} else {
		str := fmt.Sprintf("-d%s", psa.printDest)
		cmd = exec.Command("lp", str)
		_, err := cmd.Output()

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s: could not open pipe to \"%s\"\n", progname, str)
			os.Exit(13)
		}

	} //end else

	// handles for each page type

	var line string
	var lineCtr = 0 //counter of lines
	var pageCtr = 1 //counter of page
	var res string  // as for the string buffer

	rd := bufio.NewReader(fin)
	if psa.pageType == 'l' {

		for true {
			line, err = rd.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}

			lineCtr++

			if lineCtr > psa.pageLen {
				pageCtr++   //for another page
				lineCtr = 1 // start from begin of a new page

			}

			if pageCtr >= psa.startPage && pageCtr <= psa.endPage {
				// not for printer but to stdout
				if psa.printDest == "" {
					fmt.Fprintf(fout, "%s", line)
				} else {
					res += line
				}
			}
		}
	} else { // page type is '\f'
		for true {
			// read the stdin rune by rune
			char, _, err1 := rd.ReadRune()
			if err1 != nil || io.EOF == err {
				break
			}

			if char == '\f' {
				pageCtr++
			}

			if pageCtr >= psa.startPage && pageCtr <= psa.endPage {
				// output to stdout
				if psa.printDest == "" {
					fmt.Fprintf(fout, "%c", char)
				} else {
					res += string(char)
				}
			}
		}
	}
	// print destination is not empty
	if psa.printDest != "" {
		cmd.Stdin = strings.NewReader(res)
		cmd.Stdout = os.Stdout
		err = cmd.Run()
		if err != nil {
			fmt.Printf("printing to %s occurs some errors", psa.printDest)
			os.Exit(1)
		}
	}
	/* end page type handles */
	if pageCtr < psa.startPage {
		fmt.Fprintf(os.Stderr, "%s: start_page (%d) greater than total pages (%d),no output written\n", progname, psa.endPage, pageCtr)
		os.Exit(17)
	} else if pageCtr < psa.endPage {
		fmt.Fprintf(os.Stderr, "%s: end_page (%d) greater than total pages (%d), less output than expected\n", progname, psa.endPage, pageCtr)
	}
}

func main() {
	var psa structSelpg

	av := os.Args
	ac := len(os.Args)

	/*for _, str := range av {
		fmt.Printf("%s \n", str)
	}*/
	//argCount := len(os.Args[1:])

	progname = os.Args[0]

	psa.startPage = -1
	psa.endPage = -1
	psa.pageLen = 3
	psa.inFilename = ""
	psa.printDest = ""
	psa.pageType = 'l'

	processArgs(ac, av, &psa)

	processInput(psa)

	fmt.Println("done\n")

}
