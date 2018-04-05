// /anthalib/driver/liquidhandling/makelhpolicy.go: Part of the Antha language
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

package liquidtype

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/AnthaPath"
	. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	"github.com/antha-lang/antha/antha/anthalib/wtype"
)

type PolicyFile struct {
	Filename                string
	DXORJMP                 string
	FactorColumns           *[]int
	LiquidTypeStarterNumber int
}

func (polfile PolicyFile) Prepend() (prepend string) {
	nameparts := strings.Split(polfile.Filename, ".")
	prepend = nameparts[0]
	return
}

func (polfile PolicyFile) StarterNumber() (starternumber int) {
	starternumber = polfile.LiquidTypeStarterNumber
	return
}

func MakePolicyFile(filename string, dxorjmp string, factorcolumns *[]int, liquidtypestartnumber int) (policyfile PolicyFile) {
	policyfile.Filename = filename
	policyfile.DXORJMP = dxorjmp
	policyfile.FactorColumns = factorcolumns
	policyfile.LiquidTypeStarterNumber = liquidtypestartnumber
	return
}

// policy files to put in ./antha
var AvailablePolicyfiles []PolicyFile = []PolicyFile{
	MakePolicyFile("170516CCFDesign_noTouchoff_noBlowout.xlsx", "DX", nil, 100),
	MakePolicyFile("2700516AssemblyCCF.xlsx", "DX", nil, 1000),
	MakePolicyFile("newdesign2factorsonly.xlsx", "JMP", &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, 2000),
	MakePolicyFile("190516OnePolicy.xlsx", "JMP", &[]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}, 3000),
	MakePolicyFile("AssemblycategoricScreen.xlsx", "JMP", &[]int{1, 2, 3, 4, 5}, 4000),
	MakePolicyFile("090816dispenseerrordiagnosis.xlsx", "JMP", &[]int{2}, 5000),
	MakePolicyFile("090816combineddesign.xlsx", "JMP", &[]int{1}, 6000),
}

// change to range through several files
//var DOEliquidhandlingFile = "170516CCFDesign_noTouchoff_noBlowout.xlsx" // "2700516AssemblyCCF.xlsx" //"newdesign2factorsonly.xlsx" // "170516CCFDesign_noTouchoff_noBlowout.xlsx" // "170516CFF.xlsx" //"newdesign2factorsonly.xlsx" "170516CCFDesign_noTouchoff_noBlowout.xlsx" // //"newdesign2factorsonly.xlsx" //"8run4cpFactorial.xlsx" //"FullFactorial.xlsx" // "Screenwtype.LHPolicyDOE2.xlsx"
//var DXORJMP = "DX"                                                      //"JMP"
var BASEPolicy = "default" //"dna"

func MakePolicies() map[string]wtype.LHPolicy {
	pols := make(map[string]wtype.LHPolicy)

	pols = wtype.MakePolicies()

	// TODO: Remove this hack
	for _, DOEliquidhandlingFile := range AvailablePolicyfiles {
		if _, err := os.Stat(filepath.Join(anthapath.Path(), DOEliquidhandlingFile.Filename)); err == nil {
			//if antha.Anthafileexists(DOEliquidhandlingFile) {
			//fmt.Println("found lhpolicy doe file", DOEliquidhandlingFile)

			filenameparts := strings.Split(DOEliquidhandlingFile.Filename, ".")

			policies, names, _, err := PolicyMakerfromDesign(BASEPolicy, DOEliquidhandlingFile.DXORJMP, DOEliquidhandlingFile.Filename, filenameparts[0])
			//policies, names, _, err := PolicyMakerfromDesign(BASEPolicy, DXORJMP, DOEliquidhandlingFile, "DOE_run")
			for i, policy := range policies {
				pols[names[i]] = policy
			}
			if err != nil {
				panic(err)
			}
		} else {
			//	fmt.Println("no lhpolicy doe file found named: ", DOEliquidhandlingFile)
		}
	}
	return pols

}

func PolicyFilefromName(filename string) (pol PolicyFile, found bool) {
	for _, policy := range AvailablePolicyfiles {
		if policy.Filename == filename {
			pol = policy
			found = true
			return
		}
	}
	return
}

func PolicyMakerfromFilename(filename string) (policies []wtype.LHPolicy, names []string, runs []Run, err error) {

	doeliquidhandlingFile, found := PolicyFilefromName(filename)
	if found == false {
		panic("policyfilename" + filename + "not found")
	}
	filenameparts := strings.Split(doeliquidhandlingFile.Filename, ".")

	policies, names, runs, err = PolicyMakerfromDesign(BASEPolicy, doeliquidhandlingFile.DXORJMP, doeliquidhandlingFile.Filename, filenameparts[0])
	return
}

func PolicyMakerfromDesign(basepolicy string, DXORJMP string, dxdesignfilename string, prepend string) (policies []wtype.LHPolicy, names []string, runs []Run, err error) {

	policyitemmap := wtype.MakePolicyItems()
	intfactors := make([]string, 0)

	for key, val := range policyitemmap {

		if val.Type.Name() == "int" {
			intfactors = append(intfactors, key)

		}

	}
	if DXORJMP == "DX" {
		contents, err := ioutil.ReadFile(filepath.Join(anthapath.Path(), dxdesignfilename))

		if err != nil {
			return policies, names, runs, err
		}

		runs, err = RunsFromDXDesignContents(contents, intfactors)

		if err != nil {
			return policies, names, runs, err
		}

	} else if DXORJMP == "JMP" {

		factorcolumns := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13}
		responsecolumns := []int{14, 15, 16, 17}

		contents, err := ioutil.ReadFile(filepath.Join(anthapath.Path(), dxdesignfilename))

		if err != nil {
			return policies, names, runs, err
		}

		runs, err = RunsFromJMPDesignContents(contents, factorcolumns, responsecolumns, intfactors)
		if err != nil {
			return policies, names, runs, err
		}
	} else {
		return policies, names, runs, fmt.Errorf("only JMP or DX allowed as valid inputs for DXORJMP variable")
	}
	policies, names = PolicyMakerfromRuns(basepolicy, runs, prepend, false)
	return policies, names, runs, err
}

func PolicyMaker(basepolicy string, factors []DOEPair, nameprepend string, concatfactorlevelsinname bool) (policies []wtype.LHPolicy, names []string) {

	runs := AllCombinations(factors)

	policies, names = PolicyMakerfromRuns(basepolicy, runs, nameprepend, concatfactorlevelsinname)

	return
}

func PolicyMakerfromRuns(basepolicy string, runs []Run, nameprepend string, concatfactorlevelsinname bool) (policies []wtype.LHPolicy, names []string) {

	policyitemmap := wtype.MakePolicyItems()

	names = make([]string, 0)
	policies = make([]wtype.LHPolicy, 0)

	policy := wtype.MakeDefaultPolicy()
	policy.Set("CAN_MULTI", false)

	/*base, _ := GetPolicyByName(basepolicy)

	for key, value := range base {
		policy[key] = value
	}
	*/
	//fmt.Println("basepolicy:", basepolicy)
	for _, run := range runs {
		for j, desc := range run.Factordescriptors {

			_, ok := policyitemmap[desc]
			if ok {

				/*if val.Type.Name() == "int" {
					aInt, found := run.Setpoints[j].(int)

					var bInt int

					bInt = int(aInt)
					if found {
						run.Setpoints[j] = interface{}(bInt)
					}
				}*/
				policy[desc] = run.Setpoints[j]
			} /* else {
				panic("policyitem " + desc + " " + "not present! " + "These are present: " + policyitemmap.TypeList())
			}*/
		}

		// raising runtime error when using concat == true
		if concatfactorlevelsinname {
			name := nameprepend
			for key, value := range policy {
				name = fmt.Sprint(name, "_", key, ":", value)

			}

		} else {
			names = append(names, nameprepend+strconv.Itoa(run.RunNumber))
		}
		policies = append(policies, policy)

		//policy := GetPolicyByName(basepolicy)
		policy = wtype.MakeDefaultPolicy()
	}

	return
}
