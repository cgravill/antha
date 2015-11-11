// antha/AnthaStandardLibrary/Packages/Parser/RebaseParser.go: Part of the Antha language
// Copyright (C) 2015 The Antha authors. All rights reserved.
//
// This program is free software; you can redistribute it and/or
// modify it under the terms of the GNU General Public License
// as published by the Free Software Foundation; either version 2
// of the License, or (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program; if not, write to the Free Software
// Foundation, Inc., 51 Franklin Street, Fifth Floor, Boston, MA  02110-1301, USA.
//
// For more information relating to the software or licensing issues please
// contact license@antha-lang.org or write to the Antha team c/o
// Synthace Ltd. The London Bioscience Innovation Centre
// 2 Royal College St, London NW1 0NH UK

package parser

import (
	"bufio"
	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/sequences"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"io"
	"math"
	"strconv"
	"strings"
)

func ParseGenbankfile(filename string) (annotatedseq sequences.AnnotatedSeq, err error) {
	line := ""
	genbanklines := make([]string, 0)
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line = fmt.Sprintln(scanner.Text())
		genbanklines = append(genbanklines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	annotated, err := HandleGenbank(genbanklines)

	return

}
func HandleGenbank(lines []string) (annotatedseq sequences.AnnotatedSeq, err error) {

	if lines[0][0:5] == `LOCUS` {
		fmt.Println("in Locus")
		name, _, seqtype, circular, _, err := Locusline(lines[0])
		fmt.Println("foundout this stuff", name, err)
		if err != nil {
			return annotatedseq, err
		}
		if seqtype != "DNA" {
			err = fmt.Errorf("Can't parse genbank files which are not classified as type DNA at present")
			fmt.Println(err.Error())
			return annotatedseq, err
		}
		seq := HandleSequence(lines)
		fmt.Println("foundout this seq", seq)

		features := HandleFeatures(lines, seq, seqtype)
		fmt.Println("found these features", features)
		annotatedseq, err = sequences.MakeAnnotatedSeq(name, seq, circular, features)

	} else {
		err = fmt.Errorf("no LOCUS found on first line")
	}
	return
}
func Locusline(line string) (name string, seqlength int, seqtype string, circular bool, date string, err error) {

	fields := strings.SplitN(line, " ", 2)
	//fmt.Println("length of fields", len(fields))

	restofline := fields[1]

	fields = strings.Split(restofline, " ")
	//fmt.Println("length of fields", len(fields))

	newarray := make([]string, 0)
	for _, s := range fields {
		if s != "" && s != " " {
			newarray = append(newarray, s)
		}
	}
	fields = newarray
	//fmt.Println("length of fields", len(fields))
	//fmt.Println(fields)
	if len(fields) > 1 {
		name = fields[0]
		i, newerr := strconv.Atoi(fields[1])
		if newerr != nil {
			err = newerr
		}
		seqlength = i
		seqtype = fields[3]
		if fields[4] == "circular" {
			circular = true
		} else {
			circular = false
		}
		if len(fields) > 5 {
			date = fields[5]
		} else {
			date = "No date supplied"
		}
		return
	} else {
		err = fmt.Errorf("invalid genbank line: ", line)
	}

	return
}
func Cleanup(line string) (cleanarray []string) {
	fields := strings.Split(line, " ")

	for _, s := range fields {

		if s != "" && s != " " {
			cleanarray = append(cleanarray, s)
		}

	}

	return
}

func Featureline1(line string) (reverse bool, class string, startposition int, endposition int, err error) {

	newarray := Cleanup(line)

	class = newarray[0]

	for _, s := range newarray {

		if strings.Contains(s, `complement`) {
			reverse = true
			s = strings.TrimLeft(s, `(complement)`)
			s = strings.TrimRight(s, ")")
		}
		index := strings.Index(s, "..")
		if index != -1 {

			startposition, err = strconv.Atoi(s[0:index])
			if err != nil {
				fmt.Println(err.Error())
			}
			ss := strings.SplitAfter(s, "..")
			if strings.Contains(ss[1], ")") {
				ss[1] = strings.Replace(ss[1], ")", "", -1)
			}
			endposition, err = strconv.Atoi(strings.TrimRight(ss[1], "\n"))
			if err != nil {
				fmt.Println(err.Error())
			}
		}
	}
	return
}
func Featureline2(line string) (description string, found bool) {

	fields := strings.Split(line, " ")
	//fmt.Println("length of fields", len(fields))

	newarray := make([]string, 0)
	for _, s := range fields {
		if s != "" && s != " " {
			newarray = append(newarray, s)
		}
	}

	for _, line := range newarray {
		if strings.Contains(line, `/label`) {
			parts := strings.SplitAfterN(line, "=", 2)
			if len(parts) == 2 {
				description = strings.TrimSpace(parts[1])
				found = true
				return
			}

		}

	}
	return
}

func HandleFeature(lines []string) (description string, reverse bool, class string, startposition int, endposition int, err error) {

	if len(lines) > 0 {
		reverse, class, startposition, endposition, err := Featureline1(lines[0])
		//	fmt.Println(reverse, class, startposition, endposition, err)

		if err != nil {
			fmt.Errorf("Error with Featureline1 func", lines[0])
			return description, reverse, class, startposition, endposition, err
		}
		for i := 1; i < len(lines); i++ {

			description, found := Featureline2(lines[i])
			if found {
				return description, reverse, class, startposition, endposition, err
			}

		}
	}
	return
}
func DetectFeature(lines []string) (detected bool, startlineindex int, endlineindex int) {
	for i := 0; i < len(lines); i++ {

		if startlineindex != -1 && endlineindex != 0 {
			detected = true
			//		fmt.Println("Yay, detected")
			return
		}
		//	fmt.Println(lines[i])
		if string(lines[i][7]) != " " {
			startlineindex = i
			//		fmt.Println("start:", i, lines[i])
		}

		_, found := Featureline2(lines[i])
		if found {
			endlineindex = i + 1
			//		fmt.Println("end:", i, lines[i])
		}
	}

	return
}
func HandleFeatures(lines []string, seq string, seqtype string) (features []sequences.Feature) {

	features = make([]sequences.Feature, 0)
	var feature sequences.Feature

	for i := 0; i < len(lines); i++ { //, line := range lines {
		//	fmt.Println(lines)
		//	fmt.Println(line)
		if lines[i][0:8] == "FEATURES" {
			fmt.Println(lines[i])
			lines = lines[i+1 : len(lines)]
			fmt.Println("broken")
			fmt.Println(lines)
			//fmt.Println(line)
			//fmt.Println(lines[i])
			break
		}
	}
	fmt.Println("broken again")
	linesatstart := lines

	for i := 0; i < len(linesatstart); i++ {

		//jumpout := false

		if string(lines[0][0]) != " " {

			return
		}

		detected, start, end := DetectFeature(lines)
		fmt.Println("start", start, "end", end)
		if detected {
			fmt.Println("detected!!!!!!!!!!!!!", lines[start:end])

			description, reverse, class, startposition, endposition, err := HandleFeature(lines[start:end])
			fmt.Println("featuredectected: ", description, reverse, class, startposition, endposition, err)
			if err != nil {
				panic(err.Error())
			}
			rev := ""
			if reverse {
				rev = "Reverse"
			}
			feature = sequences.MakeFeature(description, seq[startposition:endposition], seqtype, class, rev)
			if start > end {
				return
			}
			features = append(features, feature)
			lines = lines[end:len(lines)]

		}

	}
	return

}

var (
	illegal string = "1234567890"
)

func HandleSequence(lines []string) (dnaseq string) {
	originallines := len(lines)
	originfound := false
	fmt.Println(originallines)
	if len(lines) > 0 {
		for i := 0; i < originallines; i++ {

			fmt.Println("lines", lines[i])
			if len([]byte(lines[0])) > 0 {
				if originfound == false {
					if lines[i][0:6] == "ORIGIN" {
						originfound = true
					}
				}
				if originfound {

					fmt.Println("i+1", i, len(lines))
					fmt.Println(lines[i+1])
					lines = lines[i+1 : originallines]
					seq := strings.Join(lines, "")
					seq = strings.Replace(seq, " ", "", -1)

					for _, character := range illegal {
						seq = strings.Replace(seq, string(character), "", -1)
					}
					seq = strings.Replace(seq, "\n", "", -1)
					seq = strings.Replace(seq, "//", "", -1)
					dnaseq = seq
					return
				}
			}
		}

	}
	return
}
