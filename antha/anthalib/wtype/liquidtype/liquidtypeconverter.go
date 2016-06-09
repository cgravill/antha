// liquidtypeconverter.go
package liquidtype

import (
	"strconv"
	"strings"
)

// horrible copy and paste of makeLHPOlicyfile code to get around import cycle issue!!!
// this will become out of date so needs to be solved better

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
}

func LiquidTypeFromPolicyDOE(s string) (match bool, number int) {
	for _, polfile := range AvailablePolicyfiles {
		if strings.Contains(s, polfile.Prepend()) {
			fields := strings.SplitAfter(s, polfile.Prepend())

			runnumber, err := strconv.Atoi(fields[1])
			if err != nil {
				panic("for Liguid type " + s + err.Error())
			}
			number = int(polfile.StarterNumber() + runnumber)
			match = true
			return
		}
	}
	return
}

func StringFromLiquidTypeNumber(lt int) (match bool, str string) {

	if lt > 99 {

		var liquidstring string
		var smallestgreaterthanzero int = 1000000000 // set to an absurdly high number to initialise

		for _, polfile := range AvailablePolicyfiles {
			if lt > polfile.StarterNumber() && int(lt)-polfile.StarterNumber() < smallestgreaterthanzero && lt > 0 {
				smallestgreaterthanzero = int(lt) - polfile.StarterNumber()
				liquidstring = polfile.Prepend() + strconv.Itoa(int(lt)-polfile.StarterNumber())

			}
		}
		str = liquidstring
		match = true
		return
	}
	return
}
