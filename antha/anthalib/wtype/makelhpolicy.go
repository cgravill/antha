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

package wtype

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	//. "github.com/antha-lang/antha/antha/AnthaStandardLibrary/Packages/doe"
	//"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/ghodss/yaml"
)

//func MakeLysatePolicy() LHPolicy {
//        lysatepolicy := make(LHPolicy, 6)
//        lysatepolicy["ASPSPEED"] = 1.0
//        lysatepolicy["DSPSPEED"] = 1.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["ASP_WAIT"] = 2.0
//        lysatepolicy["DSP_WAIT"] = 2.0
//        lysatepolicy["PRE_MIX"] = 5
//        lysatepolicy["CAN_MSA"]= false
//        return lysatepolicy
//}
//func MakeProteinPolicy() LHPolicy {
//        proteinpolicy := make(LHPolicy, 4)
//        proteinpolicy["DSPREFERENCE"] = 2
//        proteinpolicy["CAN_MULTI"] = true
//        proteinpolicy["PRE_MIX"] = 3
//        proteinpolicy["CAN_MSA"] = false
//        return proteinpolicy
//}

func MakePolicies() map[string]LHPolicy {
	pols := make(map[string]LHPolicy)

	// what policies do we need?
	pols["SmartMix"] = SmartMixPolicy()
	pols["water"] = MakeWaterPolicy()
	pols["multiwater"] = MakeMultiWaterPolicy()
	pols["culture"] = MakeCulturePolicy()
	pols["culturereuse"] = MakeCultureReusePolicy()
	pols["glycerol"] = MakeGlycerolPolicy()
	pols["solvent"] = MakeSolventPolicy()
	pols["default"] = MakeDefaultPolicy()
	pols["dna"] = MakeDNAPolicy()
	pols["DoNotMix"] = MakeDefaultPolicy()
	pols["NeedToMix"] = MakeNeedToMixPolicy()
	pols["PreMix"] = PreMixPolicy()
	pols["PostMix"] = PostMixPolicy()
	pols["MegaMix"] = MegaMixPolicy()
	pols["viscous"] = MakeViscousPolicy()
	pols["Paint"] = MakePaintPolicy()

	// pols["lysate"] = MakeLysatePolicy()
	pols["protein"] = MakeProteinPolicy()
	pols["detergent"] = MakeDetergentPolicy()
	pols["load"] = MakeLoadPolicy()
	pols["loadwater"] = MakeLoadWaterPolicy()
	pols["DispenseAboveLiquid"] = MakeDispenseAboveLiquidPolicy()
	pols["DispenseAboveLiquidMulti"] = MakeDispenseAboveLiquidMultiPolicy()
	pols["PEG"] = MakePEGPolicy()
	pols["Protoplasts"] = MakeProtoplastPolicy()
	pols["dna_mix"] = MakeDNAMixPolicy()
	pols["dna_mix_multi"] = MakeDNAMixMultiPolicy()
	pols["dna_cells_mix"] = MakeDNACELLSMixPolicy()
	pols["dna_cells_mix_multi"] = MakeDNACELLSMixMultiPolicy()
	pols["plateout"] = MakePlateOutPolicy()
	pols["colony"] = MakeColonyPolicy()
	pols["colonymix"] = MakeColonyMixPolicy()
	//      pols["lysate"] = MakeLysatePolicy()
	pols["carbon_source"] = MakeCarbonSourcePolicy()
	pols["nitrogen_source"] = MakeNitrogenSourcePolicy()

	return pols
}

var DefaultPolicies map[string]LHPolicy = MakePolicies()

func GetPolicyByName(policyname PolicyName) (lhpolicy LHPolicy, err error) {
	lhpolicy, policypresent := DefaultPolicies[policyname.String()]

	if !policypresent {
		validPolicies := availablePolicies()
		return LHPolicy{}, fmt.Errorf("policy %s not found in Default list. Valid options: %s", policyname, strings.Join(validPolicies, "\n"))
	}
	return lhpolicy, nil
}

func availablePolicies() (policies []string) {

	for key := range DefaultPolicies {
		policies = append(policies, key)
	}

	sort.Strings(policies)
	return
}

/*
Available policy field names and policy types to use:

Here is a list of everything currently implemented in the liquid handling policy framework

ASPENTRYSPEED,                    ,float64,      ,allows slow moves into liquids
ASPSPEED,                                ,float64,     ,aspirate pipetting rate
ASPZOFFSET,                           ,float64,      ,mm above well bottom when aspirating
ASP_WAIT,                                   ,float64,     ,wait time in seconds post aspirate
BLOWOUTOFFSET,                    ,float64,     ,mm above BLOWOUTREFERENCE
BLOWOUTREFERENCE,          ,int,             ,where to be when blowing out: 0 well bottom, 1 well top
BLOWOUTVOLUME,                ,float64,      ,how much to blow out
CAN_MULTI,                              ,bool,         ,is multichannel operation allowed?
DSPENTRYSPEED,                    ,float64,     ,allows slow moves into liquids
DSPREFERENCE,                      ,int,            ,where to be when dispensing: 0 well bottom, 1 well top
DSPSPEED,                              ,float64,       ,dispense pipetting rate
DSPZOFFSET,                         ,float64,          ,mm above DSPREFERENCE
DSP_WAIT,                               ,float64,        ,wait time in seconds post dispense
EXTRA_ASP_VOLUME,            ,wunit.Volume,       ,additional volume to take up when aspirating
EXTRA_DISP_VOLUME,           ,wunit.Volume,       ,additional volume to dispense
JUSTBLOWOUT,                      ,bool,            ,shortcut to get single transfer
POST_MIX,                               ,int,               ,number of mix cycles to do after dispense
POST_MIX_RATE,                    ,float64,          ,pipetting rate when post mixing
POST_MIX_VOL,                      ,float64,          ,volume to post mix (ul)
POST_MIX_X,                          ,float64,           ,x offset from centre of well (mm) when post-mixing
POST_MIX_Y,                          ,float64,           ,y offset from centre of well (mm) when post-mixing
POST_MIX_Z,                          ,float64,           ,z offset from centre of well (mm) when post-mixing
PRE_MIX,                                ,int,               ,number of mix cycles to do before aspirating
PRE_MIX_RATE,                     ,float64,           ,pipetting rate when pre mixing
PRE_MIX_VOL,                       ,float64,           ,volume to pre mix (ul)
PRE_MIX_X,                              ,float64,          ,x offset from centre of well (mm) when pre-mixing
PRE_MIX_Y,                              ,float64,           ,y offset from centre of well (mm) when pre-mixing
PRE_MIX_Z,                              ,float64,           ,z offset from centre of well (mm) when pre-mixing
TIP_REUSE_LIMIT,                    ,int,                ,number of times tips can be reused for asp/dsp cycles
TOUCHOFF,                              ,bool,             ,whether to move to TOUCHOFFSET after dispense
TOUCHOFFSET,                         ,float64,          ,mm above wb to touch off at


*/

func MakePEGPolicy() LHPolicy {
	policy := make(LHPolicy, 10)
	policy["ASPSPEED"] = 1.5
	policy["DSPSPEED"] = 1.5
	policy["ASP_WAIT"] = 2.0
	policy["DSP_WAIT"] = 2.0
	policy["ASPZOFFSET"] = 1.0
	policy["DSPZOFFSET"] = 1.0
	policy["POST_MIX"] = 3
	policy["POST_MIX_Z"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["POST_MIX_VOLUME"] = 190.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = true
	policy["RESET_OVERRIDE"] = true
	policy["DESCRIPTION"] = "Customised for handling Poly Ethylene Glycol solutions. Similar to mixing required for viscous solutions. 3 post-mixes."
	return policy
}

func MakeProtoplastPolicy() LHPolicy {
	policy := make(LHPolicy, 8)
	policy["ASPSPEED"] = 0.5
	policy["DSPSPEED"] = 0.5
	policy["ASPZOFFSET"] = 1.0
	policy["DSPZOFFSET"] = 1.0
	policy["BLOWOUTVOLUME"] = 100.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["TIP_REUSE_LIMIT"] = 5
	policy["CAN_MULTI"] = true
	policy["DESCRIPTION"] = "Customised for handling protoplast solutions. Pipettes very gently. No post-mix."
	return policy
}

func MakePaintPolicy() LHPolicy {
	policy := make(LHPolicy, 14)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.5
	policy["ASPSPEED"] = 1.5
	policy["DSPSPEED"] = 1.5
	policy["ASP_WAIT"] = 1.0
	policy["DSP_WAIT"] = 1.0
	//policy["PRE_MIX"] = 3
	policy["POST_MIX"] = 3
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = true
	policy["CAN_MULTI"] = true
	policy["DESCRIPTION"] = "Customised for handling paint solutions. Similar to mixing required for viscous solutions. 3 post-mixes."
	return policy
}

func MakeDispenseAboveLiquidPolicy() LHPolicy {
	policy := make(LHPolicy, 8)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	policy["ASPSPEED"] = 3.0
	policy["DSPSPEED"] = 3.0
	//policy["ASP_WAIT"] = 1.0
	//policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = false
	policy["DESCRIPTION"] = "Dispense solution above the liquid to facilitate tip reuse but sacrifice pipetting accuracy at low volumes. No post-mix. No multi channel"
	return policy
}
func MakeDispenseAboveLiquidMultiPolicy() LHPolicy {
	policy := make(LHPolicy, 8)
	policy["DSPREFERENCE"] = 1 // 1 indicates dispense at top of well
	policy["ASPSPEED"] = 3.0
	policy["DSPSPEED"] = 3.0
	//policy["ASP_WAIT"] = 1.0
	//policy["DSP_WAIT"] = 1.0
	policy["BLOWOUTVOLUME"] = 50.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = true
	policy["DESCRIPTION"] = "Dispense solution above the liquid to facilitate tip reuse but sacrifice pipetting accuracy at low volumes. No post Mix. Allows multi-channel pipetting."
	return policy
}

func MakeColonyPolicy() LHPolicy {
	policy := make(LHPolicy, 12)
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.0
	policy["ASPSPEED"] = 3.0
	policy["DSPSPEED"] = 3.0
	policy["ASP_WAIT"] = 1.0
	policy["POST_MIX"] = 1
	policy["BLOWOUTVOLUME"] = 0.0
	policy["BLOWOUTVOLUMEUNIT"] = "ul"
	policy["TOUCHOFF"] = false
	policy["CAN_MULTI"] = false
	policy["RESET_OVERRIDE"] = true
	policy["DESCRIPTION"] = "Designed for colony picking. 1 post-mix and no blowout (to avoid potential cross contamination), no multichannel."
	return policy
}

func MakeColonyMixPolicy() LHPolicy {
	policy := MakeColonyPolicy()
	policy["POST_MIX"] = 3
	policy["DESCRIPTION"] = "Designed for colony picking but with added post-mixes. 3 post-mix and no blowout (to avoid potential cross contamination), no multichannel."
	return policy
}

func MakeWaterPolicy() LHPolicy {
	waterpolicy := make(LHPolicy, 6)
	waterpolicy["DSPREFERENCE"] = 0
	waterpolicy["CAN_MSA"] = true
	waterpolicy["CAN_SDD"] = true
	waterpolicy["CAN_MULTI"] = false
	waterpolicy["DSPZOFFSET"] = 1.0
	waterpolicy["BLOWOUTVOLUME"] = 50.0
	waterpolicy["DESCRIPTION"] = "Default policy designed for pipetting water. Includes a blowout step for added accuracy and no post-mixing, no multi channel."
	return waterpolicy
}

func MakeMultiWaterPolicy() LHPolicy {
	pol := MakeWaterPolicy()
	pol["CAN_MULTI"] = true
	pol["DESCRIPTION"] = "Default policy designed for pipetting water but permitting multi-channel use. Includes a blowout step for added accuracy and no post-mixing."
	return pol
}

func MakeCulturePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	culturepolicy.Set("PRE_MIX", 2)
	culturepolicy.Set("PRE_MIX_VOLUME", 19.0)
	culturepolicy.Set("PRE_MIX_RATE", 3.74)
	culturepolicy.Set("ASPSPEED", 2.0)
	culturepolicy.Set("DSPSPEED", 2.0)
	culturepolicy.Set("CAN_MULTI", true)
	culturepolicy.Set("CAN_MSA", false)
	culturepolicy.Set("CAN_SDD", false)
	culturepolicy.Set("DSPREFERENCE", 0)
	culturepolicy.Set("DSPZOFFSET", 0.5)
	culturepolicy.Set("TIP_REUSE_LIMIT", 0)
	culturepolicy.Set("NO_AIR_DISPENSE", true)
	culturepolicy.Set("BLOWOUTVOLUME", 0.0)
	culturepolicy.Set("BLOWOUTVOLUMEUNIT", "ul")
	culturepolicy.Set("TOUCHOFF", false)
	culturepolicy.Set("DESCRIPTION", "Designed for cell cultures. Tips will not be reused to minimise any risk of cross contamination and 2 pre-mixes will be performed prior to aspirating.")
	return culturepolicy
}

func MakePlateOutPolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 17)
	culturepolicy["CAN_MULTI"] = true
	culturepolicy["ASP_WAIT"] = 1.0
	culturepolicy["DSP_WAIT"] = 1.0
	culturepolicy["DSPZOFFSET"] = 0.0
	culturepolicy["TIP_REUSE_LIMIT"] = 7
	culturepolicy["NO_AIR_DISPENSE"] = true
	culturepolicy["TOUCHOFF"] = false
	culturepolicy["RESET_OVERRIDE"] = true
	culturepolicy["DESCRIPTION"] = "Designed for plating out cultures onto agar plates. Dispense will be performed at the well bottom and no blowout will be performed (to minimise risk of cross contamination)"
	return culturepolicy
}

func MakeCultureReusePolicy() LHPolicy {
	culturepolicy := make(LHPolicy, 10)
	culturepolicy.Set("PRE_MIX", 2)
	culturepolicy.Set("PRE_MIX_VOLUME", 19.0)
	culturepolicy.Set("PRE_MIX_RATE", 3.74)
	culturepolicy.Set("ASPSPEED", 2.0)
	culturepolicy.Set("DSPSPEED", 2.0)
	culturepolicy.Set("CAN_MULTI", true)
	culturepolicy.Set("CAN_MSA", true)
	culturepolicy.Set("CAN_SDD", true)
	culturepolicy.Set("DSPREFERENCE", 0)
	culturepolicy.Set("DSPZOFFSET", 0.5)
	culturepolicy.Set("NO_AIR_DISPENSE", true)
	culturepolicy.Set("BLOWOUTVOLUME", 0.0)
	culturepolicy.Set("BLOWOUTVOLUMEUNIT", "ul")
	culturepolicy.Set("TOUCHOFF", false)
	culturepolicy.Set("DESCRIPTION", "Designed for cell cultures but permitting tip reuse when handling the same culture. 2 pre-mixes will be performed prior to aspirating.")
	return culturepolicy
}

func MakeGlycerolPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 9)
	glycerolpolicy["ASPSPEED"] = 1.5
	glycerolpolicy["DSPSPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	glycerolpolicy["TIP_REUSE_LIMIT"] = 0
	glycerolpolicy["CAN_MULTI"] = true
	glycerolpolicy["POST_MIX"] = 3
	glycerolpolicy["POST_MIX_VOLUME"] = 20.0
	glycerolpolicy["POST_MIX_RATE"] = 3.74 // Should this be the same rate as the asp and dsp speeds?
	glycerolpolicy["DESCRIPTION"] = "Designed for viscous samples, in particular enzymes stored in glycerol. 3 gentle post-mixes of 20ul will be performed. Tips will not be reused in order to increase accuracy."
	return glycerolpolicy
}

func MakeViscousPolicy() LHPolicy {
	glycerolpolicy := make(LHPolicy, 7)
	glycerolpolicy["ASPSPEED"] = 1.5
	glycerolpolicy["DSPSPEED"] = 1.5
	glycerolpolicy["ASP_WAIT"] = 1.0
	glycerolpolicy["DSP_WAIT"] = 1.0
	glycerolpolicy["CAN_MULTI"] = true
	glycerolpolicy["POST_MIX"] = 3
	glycerolpolicy["POST_MIX_RATE"] = 1.5
	glycerolpolicy["DESCRIPTION"] = "Designed for viscous samples. 3 post-mixes of the volume of the sample being transferred will be performed. No tip reuse limit."
	return glycerolpolicy
}
func MakeSolventPolicy() LHPolicy {
	solventpolicy := make(LHPolicy, 5)
	solventpolicy.Set("PRE_MIX", 3)
	solventpolicy.Set("DSPREFERENCE", 0)
	solventpolicy.Set("DSPZOFFSET", 0.5)
	solventpolicy.Set("NO_AIR_DISPENSE", true)
	solventpolicy.Set("CAN_MULTI", true)
	solventpolicy.Set("DESCRIPTION", "Designed for handling solvents. No post-mixes are performed")
	return solventpolicy
}

func MakeDNAPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	dnapolicy["ASPSPEED"] = 2.0
	dnapolicy["DSPSPEED"] = 2.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["RESET_OVERRIDE"] = true
	dnapolicy["TOUCHOFF"] = false
	dnapolicy["DESCRIPTION"] = "Designed for DNA samples. No tip reuse is permitted, no blowout and no post-mixing."
	return dnapolicy
}

func MakeDNAMixPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 10.0
	dnapolicy["POST_MIX"] = 5
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["CAN_MULTI"] = false
	dnapolicy["DESCRIPTION"] = "Designed for DNA samples but with 5 post-mixes of 10ul. No tip reuse is permitted, no blowout, no multichannel."
	return dnapolicy
}

func MakeDNAMixMultiPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 10.0
	dnapolicy["POST_MIX"] = 5
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["DESCRIPTION"] = "Designed for DNA samples but with 5 post-mixes of 10ul. No tip reuse is permitted, no blowout. Allows multi-channel pipetting."
	return dnapolicy
}

func MakeDNACELLSMixPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 20.0
	dnapolicy["POST_MIX"] = 2
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 1.0
	dnapolicy["DESCRIPTION"] = "Designed for mixing DNA with cells. 2 gentle post-mixes are performed. No tip reuse is permitted, no blowout."
	return dnapolicy
}
func MakeDNACELLSMixMultiPolicy() LHPolicy {
	dnapolicy := MakeDNAPolicy()
	dnapolicy["POST_MIX_VOLUME"] = 20.0
	dnapolicy["POST_MIX"] = 2
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 1.0
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["DESCRIPTION"] = "Designed for mixing DNA with cells. 2 gentle post-mixes are performed. No tip reuse is permitted, no blowout. Allows multi-channel pipetting."
	return dnapolicy
}

func MakeDetergentPolicy() LHPolicy {
	detergentpolicy := make(LHPolicy, 9)
	//        detergentpolicy["POST_MIX"] = 3
	detergentpolicy["ASPSPEED"] = 1.0
	detergentpolicy["DSPSPEED"] = 1.0
	detergentpolicy["CAN_MSA"] = false
	detergentpolicy["CAN_SDD"] = false
	detergentpolicy["DSPREFERENCE"] = 0
	detergentpolicy["DSPZOFFSET"] = 0.5
	detergentpolicy["TIP_REUSE_LIMIT"] = 8
	detergentpolicy["NO_AIR_DISPENSE"] = true
	detergentpolicy["DESCRIPTION"] = "Designed for solutions containing detergents. Gentle aspiration and dispense and a tip reuse limit of 8 to reduce problem of foam build up inside the tips."
	return detergentpolicy
}
func MakeProteinPolicy() LHPolicy {
	proteinpolicy := make(LHPolicy, 12)
	proteinpolicy["POST_MIX"] = 5
	proteinpolicy["POST_MIX_VOLUME"] = 50.0
	proteinpolicy["ASPSPEED"] = 2.0
	proteinpolicy["DSPSPEED"] = 2.0
	proteinpolicy["CAN_MSA"] = false
	proteinpolicy["CAN_SDD"] = false
	proteinpolicy["DSPREFERENCE"] = 0
	proteinpolicy["DSPZOFFSET"] = 0.5
	proteinpolicy["TIP_REUSE_LIMIT"] = 0
	proteinpolicy["NO_AIR_DISPENSE"] = true
	proteinpolicy["DESCRIPTION"] = "Designed for protein solutions. Slightly gentler aspiration and dispense and a tip reuse limit of 0 to prevent risk of cross contamination. 5 post-mixes of 50ul will be performed."
	return proteinpolicy
}
func MakeLoadPolicy() LHPolicy {

	loadpolicy := make(LHPolicy, 14)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
	loadpolicy["CAN_MSA"] = false
	loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 0
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["BLOWOUTREFERENCE"] = 1
	loadpolicy["BLOWOUTOFFSET"] = 0.0
	loadpolicy["BLOWOUTVOLUME"] = 0.0
	loadpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	loadpolicy["DESCRIPTION"] = "Designed for loading a sample onto an agarose gel. Very slow dispense rate, no tip reuse and no blowout."
	return loadpolicy
}

func MakeLoadWaterPolicy() LHPolicy {
	loadpolicy := make(LHPolicy)
	loadpolicy["ASPSPEED"] = 1.0
	loadpolicy["DSPSPEED"] = 0.1
	loadpolicy["CAN_MSA"] = false
	//loadpolicy["CAN_SDD"] = false
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["NO_AIR_DISPENSE"] = true
	loadpolicy["TOUCHOFF"] = false
	loadpolicy["TIP_REUSE_LIMIT"] = 100
	loadpolicy["BLOWOUTREFERENCE"] = 1
	loadpolicy["BLOWOUTOFFSET"] = 0.0
	loadpolicy["BLOWOUTVOLUME"] = 0.0
	loadpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	loadpolicy["DESCRIPTION"] = "Designed for loading water into agarose gel wells so permits tip reuse. Very slow dispense rate and no blowout."
	return loadpolicy
}

func MakeNeedToMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 16)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["PRE_MIX"] = 3
	dnapolicy["PRE_MIX_VOLUME"] = 20.0
	dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 pre-mixes and 3 post-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy
}

func PreMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	//dnapolicy["POST_MIX"] = 3
	//dnapolicy[""POST_MIX_VOLUME"] = 10.0
	//dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["PRE_MIX"] = 3
	dnapolicy["PRE_MIX_VOLUME"] = 19.0
	dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 pre-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy

}

func PostMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	dnapolicy["POST_MIX"] = 3
	dnapolicy["POST_MIX_RATE"] = 3.74
	//dnapolicy["PRE_MIX"] = 3
	//dnapolicy["PRE_MIX_VOLUME"] = 10
	//dnapolicy["PRE_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "3 post-mixes of the sample being transferred.  No tip reuse permitted."
	return dnapolicy
}

// 3 post mixes of the sample being transferred. Volume is adjusted based upon the volume of liquid in the destination well.
// No tip reuse permitted.
// Rules added to adjust post mix volume based on volume of the destination well.
// volume now capped at max for tip type (MIX_VOLUME_OVERRIDE_TIP_MAX)
func SmartMixPolicy() LHPolicy {
	policy := make(LHPolicy, 12)
	policy["POST_MIX"] = 3
	policy["POST_MIX_RATE"] = 3.74
	policy["ASPSPEED"] = 3.74
	policy["DSPSPEED"] = 3.74
	policy["CAN_MULTI"] = true
	policy["CAN_MSA"] = false
	policy["CAN_SDD"] = false
	policy["DSPREFERENCE"] = 0
	policy["DSPZOFFSET"] = 0.5
	policy["TIP_REUSE_LIMIT"] = 0
	policy["NO_AIR_DISPENSE"] = true
	policy["DESCRIPTION"] = "3 post-mixes of the sample being transferred. Volume is adjusted based upon the volume of liquid in the destination well.  No tip reuse permitted."
	policy["MIX_VOLUME_OVERRIDE_TIP_MAX"] = true
	return policy
}

func MegaMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 12)
	dnapolicy["POST_MIX"] = 10
	dnapolicy["POST_MIX_RATE"] = 3.74
	dnapolicy["ASPSPEED"] = 3.74
	dnapolicy["DSPSPEED"] = 3.74
	dnapolicy["CAN_MULTI"] = true
	dnapolicy["CAN_MSA"] = false
	dnapolicy["CAN_SDD"] = false
	dnapolicy["DSPREFERENCE"] = 0
	dnapolicy["DSPZOFFSET"] = 0.5
	dnapolicy["TIP_REUSE_LIMIT"] = 0
	dnapolicy["NO_AIR_DISPENSE"] = true
	dnapolicy["DESCRIPTION"] = "10 post-mixes of the sample being transferred. No tip reuse permitted."
	return dnapolicy

}

func MakeDefaultPolicy() LHPolicy {
	defaultpolicy := make(LHPolicy, 29)
	defaultpolicy["MIX_VOLUME_OVERRIDE_TIP_MAX"] = false
	defaultpolicy["OFFSETZADJUST"] = 0.0
	defaultpolicy["TOUCHOFF"] = false
	defaultpolicy["TOUCHOFFSET"] = 0.5
	defaultpolicy["ASPREFERENCE"] = 0
	defaultpolicy["ASPZOFFSET"] = 0.5
	defaultpolicy["DSPREFERENCE"] = 0
	defaultpolicy["DSPZOFFSET"] = 0.5
	defaultpolicy["CAN_MSA"] = false
	defaultpolicy["CAN_SDD"] = true
	defaultpolicy["CAN_MULTI"] = true
	defaultpolicy["TIP_REUSE_LIMIT"] = 100
	defaultpolicy["BLOWOUTREFERENCE"] = 1
	defaultpolicy["BLOWOUTVOLUME"] = 50.0
	defaultpolicy["BLOWOUTOFFSET"] = 0.0 //-5.0
	defaultpolicy["BLOWOUTVOLUMEUNIT"] = "ul"
	defaultpolicy["PTZREFERENCE"] = 1
	defaultpolicy["PTZOFFSET"] = -0.5
	defaultpolicy["NO_AIR_DISPENSE"] = true // SERIOUSLY??
	defaultpolicy["DEFAULTPIPETTESPEED"] = 3.0
	defaultpolicy["MANUALPTZ"] = false
	defaultpolicy["JUSTBLOWOUT"] = false
	defaultpolicy["DONT_BE_DIRTY"] = true
	defaultpolicy["POST_MIX_Z"] = 0.5
	defaultpolicy["PRE_MIX_Z"] = 0.5
	defaultpolicy["LLFABOVESURFACE"] = 3.0 //distance above liquid level for dispensing with LiqudLevelFollowing
	defaultpolicy["LLFBELOWSURFACE"] = 3.0 //distance below liquid level for aspirating with LLF
	defaultpolicy["DESCRIPTION"] = "Default mix Policy. Blowout performed, no touch off, no mixing, tip reuse permitted for the same solution."

	return defaultpolicy
}

func MakeJBPolicy() LHPolicy {
	jbp := make(LHPolicy, 1)
	jbp.Set("JUSTBLOWOUT", true)
	jbp.Set("TOUCHOFF", true)
	return jbp
}

func MakeTOPolicy() LHPolicy {
	top := make(LHPolicy, 1)
	top.Set("TOUCHOFF", true)
	return top
}

func MakeLVExtraPolicy() LHPolicy {
	lvep := make(LHPolicy, 2)
	lvep.Set("EXTRA_ASP_VOLUME", wunit.NewVolume(0.5, "ul"))
	lvep.Set("EXTRA_DISP_VOLUME", wunit.NewVolume(0.5, "ul"))
	return lvep
}

func MakeLVDNAMixPolicy() LHPolicy {
	dnapolicy := make(LHPolicy, 4)
	dnapolicy["RESET_OVERRIDE"] = true
	dnapolicy["POST_MIX_VOLUME"] = 5.0
	dnapolicy["POST_MIX"] = 1
	dnapolicy["POST_MIX_Z"] = 0.5
	dnapolicy["POST_MIX_RATE"] = 3.0
	dnapolicy["TOUCHOFF"] = false
	return dnapolicy
}

func TurnOffBlowoutPolicy() LHPolicy {
	loadpolicy := make(LHPolicy, 1)
	loadpolicy["RESET_OVERRIDE"] = true
	return loadpolicy
}

func MakeHVOffsetPolicy() LHPolicy {
	lvop := make(LHPolicy, 6)
	lvop["OFFSETZADJUST"] = 0.75
	lvop["POST_MIX_RATE"] = 37
	lvop["PRE_MIX_RATE"] = 37
	lvop["ASPSPEED"] = 37
	lvop["DSPSPEED"] = 37
	return lvop
}

func AdjustPostMixVolume(mixToVol wunit.Volume) LHPolicy {
	vol := mixToVol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
	policy := make(LHPolicy, 1)
	policy["POST_MIX_VOLUME"] = vol
	return policy
}

func AdjustPreMixVolume(mixToVol wunit.Volume) LHPolicy {
	vol := mixToVol.ConvertTo(wunit.ParsePrefixedUnit("ul"))
	policy := make(LHPolicy, 1)
	policy["PRE_MIX_VOLUME"] = vol
	return policy
}

// deprecated; see above
func MakeHVFlowRatePolicy() LHPolicy {
	policy := make(LHPolicy, 4)
	policy["POST_MIX_RATE"] = 37
	policy["PRE_MIX_RATE"] = 37
	policy["ASPSPEED"] = 37
	policy["DSPSPEED"] = 37
	return policy
}

func MakeCarbonSourcePolicy() LHPolicy {
	cspolicy := make(LHPolicy, 1)
	cspolicy["DSPREFERENCE"] = 1
	cspolicy["DESCRIPTION"] = "Custom policy for carbon source which dispenses above destination solution."
	return cspolicy
}

func MakeNitrogenSourcePolicy() LHPolicy {
	nspolicy := make(LHPolicy, 1)
	nspolicy["DSPREFERENCE"] = 1
	nspolicy["DESCRIPTION"] = "Custom policy for nitrogen source which dispenses above destination solution."
	return nspolicy
}

// newConditionalRule makes a new LHPolicyRule with conditions to apply to an LHPolicy.
//
// An error is returned if an invalid Condition Class or SetPoint is specified.
// The valid Setpoints can be found in MakeInstructionParameters()
func newConditionalRule(ruleName string, conditions ...condition) (LHPolicyRule, error) {
	var errs []string

	rule := NewLHPolicyRule(ruleName)
	for _, condition := range conditions {
		err := condition.AddToRule(rule)
		if err != nil {
			errs = append(errs, err.Error())
		}
	}
	if len(errs) > 0 {
		return rule, fmt.Errorf(strings.Join(errs, ".\n"))
	}
	return rule, nil
}

type condition interface {
	AddToRule(LHPolicyRule) error
}

type categoricCondition struct {
	Class    string
	SetPoint string
}

func (c categoricCondition) AddToRule(rule LHPolicyRule) error {
	return rule.AddCategoryConditionOn(c.Class, c.SetPoint)
}

type numericCondition struct {
	Class string
	Range conditionRange
}

type conditionRange struct {
	Lower float64
	Upper float64
}

func (c numericCondition) AddToRule(rule LHPolicyRule) error {
	return rule.AddNumericConditionOn(c.Class, c.Range.Lower, c.Range.Upper)
}

// Conditions to apply to LHpolicyRules based on liquid policy used
var (
	OnSmartMix  = categoricCondition{"LIQUIDCLASS", "SmartMix"}
	OnPostMix   = categoricCondition{"LIQUIDCLASS", "PostMix"}
	OnPreMix    = categoricCondition{"LIQUIDCLASS", "PreMix"}
	OnNeedToMix = categoricCondition{"LIQUIDCLASS", "NeedToMix"}
)

// Conditions to apply to LHpolicyRules based on volume of liquid that a sample is being pipetted into at the destination well
var (
	IntoLessThan20ul          = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{Lower: 0.0, Upper: 20.0}}
	IntoBetween20ulAnd50ul    = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{20.0, 50.0}}
	IntoBetween50ulAnd100ul   = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{50.0, 100.0}}
	IntoBetween100ulAnd200ul  = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{100.0, 200.0}}
	IntoBetween200ulAnd1000ul = numericCondition{Class: "TOWELLVOLUME", Range: conditionRange{200.0, 1000.0}}
)

// Conditions to apply to LHpolicyRules based on volume of liquid being transferred
var (
	LessThan20ul = numericCondition{Class: "VOLUME", Range: conditionRange{0.0, 20.0}}
)

// Conditions to apply to LHpolicyRules based on volume of liquid in source well from which a sample is taken
var (
	FromBetween100ulAnd200ul  = numericCondition{Class: "WELLFROMVOLUME", Range: conditionRange{100.0, 200.0}}
	FromBetween200ulAnd1000ul = numericCondition{Class: "WELLFROMVOLUME", Range: conditionRange{200.0, 1000.0}}
)

func AddUniversalRules(originalRuleSet *LHPolicyRuleSet, policies map[string]LHPolicy) (lhpr *LHPolicyRuleSet, err error) {

	lhpr = originalRuleSet

	for name, policy := range policies {
		rule := NewLHPolicyRule(name)
		err := rule.AddCategoryConditionOn("LIQUIDCLASS", name)

		if err != nil {
			return nil, err
		}
		lhpr.AddRule(rule, policy)
	}

	// hack to fix plate type problems
	// this really should be removed asap
	rule := NewLHPolicyRule("HVOffsetFix")
	//rule.AddNumericConditionOn("VOLUME", 20.1, 300.0) // what about higher? // set specifically for openPlant configuration

	rule.AddCategoryConditionOn("TIPTYPE", "Gilson200")
	rule.AddCategoryConditionOn("PLATFORM", "GilsonPipetmax")
	// don't get overridden
	rule.Priority = 100
	pol := MakeHVOffsetPolicy()
	lhpr.AddRule(rule, pol)

	// merged the below and the above
	/*
		rule = NewLHPolicyRule("HVFlowRate")
		rule.AddNumericConditionOn("VOLUME", 20.1, 300.0) // what about higher? // set specifically for openPlant configuration
		//rule.AddCategoryConditionOn("FROMPLATETYPE", "pcrplate_skirted_riser")
		pol = MakeHVFlowRatePolicy()
		lhpr.AddRule(rule, pol)
	*/

	rule = NewLHPolicyRule("DNALV")
	rule.AddNumericConditionOn("VOLUME", 0.0, 1.99)
	rule.AddCategoryConditionOn("LIQUIDCLASS", "dna")
	pol = MakeLVDNAMixPolicy()
	lhpr.AddRule(rule, pol)
	return lhpr, nil
}

func GetLHPolicyForTest() (*LHPolicyRuleSet, error) {
	// make some policies

	policies := MakePolicies()

	// now make rules

	lhpr := NewLHPolicyRuleSet()

	lhpr, err := AddUniversalRules(lhpr, policies)

	if err != nil {
		return nil, err
	}

	for name, policy := range policies {
		rule := NewLHPolicyRule(name)
		err := rule.AddCategoryConditionOn("LIQUIDCLASS", name)

		if err != nil {
			return nil, err
		}
		lhpr.AddRule(rule, policy)
	}

	adjustPostMix, err := newConditionalRule("mixInto20ul", OnSmartMix, IntoBetween20ulAnd50ul)

	if err != nil {
		return lhpr, err
	}

	adjustVol20 := AdjustPostMixVolume(wunit.NewVolume(20, "ul"))
	adjustVol50 := AdjustPostMixVolume(wunit.NewVolume(50, "ul"))
	adjustVol100 := AdjustPostMixVolume(wunit.NewVolume(100, "ul"))
	adjustVol200 := AdjustPostMixVolume(wunit.NewVolume(200, "ul"))

	lhpr.AddRule(adjustPostMix, adjustVol20)

	adjustPostMix50, err := newConditionalRule("mixInto50ul", OnSmartMix, IntoBetween50ulAnd100ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPostMix50, adjustVol50)

	adjustPostMix100, err := newConditionalRule("mixInto100ul", OnSmartMix, IntoBetween100ulAnd200ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPostMix100, adjustVol100)

	adjustPostMix200, err := newConditionalRule("mixInto200ul", OnSmartMix, IntoBetween200ulAnd1000ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPostMix200, adjustVol200)

	// adjust original PostMix and NeedToMix policy to only set post mix volume if low volume.
	postmix20ul, err := newConditionalRule("PostMix20ul", OnPostMix, LessThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(postmix20ul, adjustVol20)

	needToMix20ul, err := newConditionalRule("NeedToMix20ul", OnNeedToMix, LessThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(needToMix20ul, adjustVol20)

	// now pre mix values for PreMix and NeedToMix
	adjustPreMixVol20 := AdjustPreMixVolume(wunit.NewVolume(20, "ul"))
	adjustPreMixVol100 := AdjustPreMixVolume(wunit.NewVolume(100, "ul"))
	adjustPreMixVol200 := AdjustPreMixVolume(wunit.NewVolume(200, "ul"))

	// PreMix
	adjustPreMix, err := newConditionalRule("preMix20ul", OnPreMix, LessThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPreMix, adjustPreMixVol20)

	adjustPreMix100ul, err := newConditionalRule("PreMixFrom100ul", OnPreMix, FromBetween100ulAnd200ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPreMix100ul, adjustPreMixVol100)

	adjustPreMix200ul, err := newConditionalRule("PreMixFrom200ul", OnPreMix, FromBetween200ulAnd1000ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustPreMix200ul, adjustPreMixVol200)

	// NeedToMix
	adjustNeedToMix, err := newConditionalRule("NeedToPreMix20ul", OnNeedToMix, LessThan20ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustNeedToMix, adjustPreMixVol20)

	adjustNeedToMix100ul, err := newConditionalRule("NeedToPreMixFrom100ul", OnNeedToMix, FromBetween100ulAnd200ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustNeedToMix100ul, adjustPreMixVol100)

	adjustNeedToMix200ul, err := newConditionalRule("NeedToPreMixFrom200ul", OnNeedToMix, FromBetween200ulAnd1000ul)

	if err != nil {
		return lhpr, err
	}

	lhpr.AddRule(adjustNeedToMix200ul, adjustPreMixVol200)

	//fix for removing blowout in DNA only if EGEL 48 plate type is used
	rule := NewLHPolicyRule("EPAGE48Load")
	rule.AddCategoryConditionOn("TOPLATETYPE", "EPAGE48")
	pol := TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 48 plate type is used
	rule = NewLHPolicyRule("EGEL48Load")
	rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL48")
	pol = TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 96_1 plate type is used
	rule = NewLHPolicyRule("EGEL961Load")
	rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL96_1")
	pol = TurnOffBlowoutPolicy()
	lhpr.AddRule(rule, pol)

	//fix for removing blowout in DNA only if EGEL 96_2 plate type is used
	rule = NewLHPolicyRule("EGEL962Load")
	rule.AddCategoryConditionOn("TOPLATETYPE", "EGEL96_2")
	pol = TurnOffBlowoutPolicy()

	lhpr.AddRule(rule, pol)

	return lhpr, nil

}

func LoadLHPoliciesFromFile() (*LHPolicyRuleSet, error) {
	lhPoliciesFileName := os.Getenv("ANTHA_LHPOLICIES_FILE")
	if lhPoliciesFileName == "" {
		return nil, fmt.Errorf("Env variable ANTHA_LHPOLICIES_FILE not set")
	}
	contents, err := ioutil.ReadFile(lhPoliciesFileName)
	if err != nil {
		return nil, err
	}
	lhprs := NewLHPolicyRuleSet()
	lhprs.Policies = make(map[string]LHPolicy)
	lhprs.Rules = make(map[string]LHPolicyRule)
	//	err = readYAML(contents, lhprs)
	err = readJSON(contents, lhprs)
	if err != nil {
		return nil, err
	}
	return lhprs, nil
}

func readYAML(fileContents []byte, ruleSet *LHPolicyRuleSet) error {
	if err := yaml.Unmarshal(fileContents, ruleSet); err != nil {
		return err
	}
	return nil
}

func readJSON(fileContents []byte, ruleSet *LHPolicyRuleSet) error {
	if err := json.Unmarshal(fileContents, ruleSet); err != nil {
		return err
	}
	return nil
}

/*
	LTNIL = LiquidType("nil")
	LTWater = LiquidType("water")
	LTDefault = LiquidType("default")
	LTCulture = LiquidType("culture")
	LTProtoplasts = LiquidType("protoplasts")
	LTDNA = LiquidType("dna")
	LTDNAMIX = LiquidType("dna_mix")
	LTProtein = LiquidType("protein")
	LTMultiWater = LiquidType("multiwater")
	LTLoad = LiquidType("load")
	LTVISCOUS = LiquidType("viscous")
	LTPEG = LiquidType("peg")
	LTPAINT = LiquidType("paint")
	LTNeedToMix = LiquidType("NeedToMix")
	LTPostMix = LiquidType("PostMix")
	LTload = LiquidType("load")
	LTGlycerol = LiquidType("glycerol")
	LTPLATEOUT = LiquidType("plateout")
	LTDetergent = LiquidType("detergent")
	LTCOLONY = LiquidType("colony")
	LTNSrc = LiquidType("nitrogen_source")
	InvalidPolicyName = LiquidType("InvalidPolicyName")
	LTSmartMix = LiquidType("SmartMix")
	LTPreMix = LiquidType("PreMix")
	LTDISPENSEABOVE = LiquidType("DispenseAboveLiquid")
	LTMegaMix = LiquidType("MegaMix")
	LTDoNotMix = LiquidType("DoNotMix")
	LTDNACELLSMIX = LiquidType("dna_cells_mix")
	LTloadwater = LiquidType("loadwater")
*/
