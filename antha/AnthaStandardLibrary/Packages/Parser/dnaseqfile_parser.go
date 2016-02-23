// antha/AnthaStandardLibrary/Packages/Parser/fasta_parser.go: Part of the Antha language
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
	"fmt"
	//"os"
	"strings"

	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

func DNAFiletoDNASequence(filename string, plasmid bool) (seqs []wtype.DNASequence, err error) {

	seqs = make([]wtype.DNASequence, 0)
	var seq wtype.DNASequence
	//for _, file := range files {
	switch fn := filename; {
	case strings.HasSuffix(fn, ".gdx"):

		seqs, err = GDXtoDNASequence(filename)
	case strings.HasSuffix(fn, ".fasta"):
		if plasmid {
			seqs, err = FASTAtoPlasmidDNASeqs(filename)
		} else {
			seqs, err = FASTAtoLinearDNASeqs(filename)
		}
	case strings.HasSuffix(fn, ".gb"):

		seq, err = GenbanktoDNASequence(filename)

		seqs = append(seqs, seq)
	default:
		err = fmt.Errorf("non valid sequence file")
	}

	if err != nil {
		return seqs, err
	}
	//}
	return
}
