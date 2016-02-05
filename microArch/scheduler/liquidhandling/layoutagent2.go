// liquidhandling/layoutagent2.go: Part of the Antha language
// Copyright (C) 2014 the Antha authors. All rights reserved.
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

package liquidhandling

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	//"github.com/antha-lang/antha/antha/anthalib/wutil"
	"github.com/antha-lang/antha/microArch/driver/liquidhandling"
	"github.com/antha-lang/antha/microArch/factory"
	"github.com/antha-lang/antha/microArch/logger"
	//"github.com/antha-lang/antha/microArch/logger"
)

func ImprovedLayoutAgent(request *LHRequest, params *liquidhandling.LHProperties) *LHRequest {
	// we have three kinds of solution
	// 1- ones going to a specific plate
	// 2- ones going to a specific plate type
	// 3- ones going to a plate of our choosing

	// find existing assignments

	plate_choices := get_and_complete_assignments(request)

	// now we know what remains unassigned, we assign it

	plate_choices = choose_plates(request, plate_choices)

	// now we have plates of type 1 & 2

	// make specific plates... this may mean splitting stuff out into multiple plates

	make_plates(request)

	// now we have solutions of type 1 only -- we just need to
	// say where on each plate they will go
	// this needs to set Output_assignments
	make_layouts(request, plate_choices)

	return request
}

type PlateChoice struct {
	Platetype string
	Assigned  []string
	ID        string
	Wells     []string
}

func get_and_complete_assignments(request *LHRequest) []PlateChoice {
	s := make([]PlateChoice, 0, 3)
	m := make(map[int]string)

	// inconsistent plate types will be assigned randomly!
	//	for k, v := range request.LHInstructions {
	for _, k := range request.Order_instructions_added {
		v := request.LHInstructions[k]
		if v.PlateID != "" {
			i := defined(v.PlateID, s)

			if i == -1 {
				s = append(s, PlateChoice{v.Platetype, []string{v.ID}, v.PlateID, []string{}})
			} else {
				s[i].Assigned = append(s[i].Assigned, v.ID)
				s[i].Wells = append(s[i].Wells, v.Welladdress)
			}
		} else if v.Majorlayoutgroup != -1 {
			id, ok := m[v.Majorlayoutgroup]
			if !ok {
				id := wtype.NewUUID()
				m[v.Majorlayoutgroup] = id
				//  fix the plate id to this temporary one
				request.LHInstructions[k].PlateID = id
			}

			i := defined(id, s)

			if i == -1 {
				s = append(s, PlateChoice{v.Platetype, []string{v.ID}, v.PlateID, []string{}})
			} else {
				s[i].Assigned = append(s[i].Assigned, v.ID)
				s[i].Wells = append(s[i].Wells, v.Welladdress)
			}
		}
	}

	// make sure the plate choices all have defined types

	for i, _ := range s {
		if s[i].Platetype == "" {
			s[i].Platetype = request.Output_platetypes[0].Type
		}
	}

	return s
}

func defined(s string, pc []PlateChoice) int {
	r := -1

	for i, v := range pc {
		if v.ID == s {
			r = i
			break
		}
	}
	return r
}

func choose_plates(request *LHRequest, pc []PlateChoice) []PlateChoice {
	for _, v := range request.LHInstructions {
		// this id may be temporary, only things without it still are not assigned to a
		// plate, even a virtual one
		if v.PlateID == "" {
			pt := v.Platetype

			ass := assignmentWithType(pt, pc)

			if ass == -1 {
				// make a new plate
				ass = len(pc)
				pc = append(pc, PlateChoice{chooseAPlate(request, v), []string{v.ID}, wtype.GetUUID(), []string{}})
			}

			pc[ass].Assigned = append(pc[ass].Assigned, v.ID)
		}
	}

	// make sure the plate isn't too full

	pc2 := make([]PlateChoice, 0, len(pc))

	for _, v := range pc {
		plate := factory.GetPlateByType(v.Platetype)

		// chop the assignments up

		pc2 = append(pc2, modpc(v, plate.Nwells)...)
	}

	// copy the choices in

	for _, c := range pc2 {
		for _, i := range c.Assigned {
			request.LHInstructions[i].PlateID = c.ID
			request.LHInstructions[i].Platetype = c.Platetype
		}
	}
	return pc2
}

// chop the assignments up modulo plate size
func modpc(choice PlateChoice, nwell int) []PlateChoice {
	r := make([]PlateChoice, 0, 1)

	for s := 0; s < len(choice.Assigned); s += nwell {
		e := s + nwell
		if e > len(choice.Assigned) {
			e = len(choice.Assigned)
		}
		r = append(r, PlateChoice{choice.Platetype, choice.Assigned[s:e], wtype.GetUUID(), choice.Wells[s:e]})
	}
	return r
}

func assignmentWithType(pt string, pc []PlateChoice) int {
	r := -1

	if pt == "" {
		if len(pc) != 0 {
			r = 0
		}
		return r
	}

	for i, v := range pc {
		if pt == v.Platetype {
			r = i
			break
		}
	}

	return r
}

func chooseAPlate(request *LHRequest, ins *wtype.LHInstruction) string {
	// for now we ignore ins and just choose the First Output Platetype
	return request.Output_platetypes[0].Type
}
func stringinarray(s string, array []string) int {
	r := -1

	for i, k := range array {
		if k == s {
			r = i
			break
		}
	}

	return r
}

func plateidarray(arr []*wtype.LHPlate) []string {
	ret := make([]string, 0, 3)

	for _, v := range arr {
		ret = append(ret, v.ID)
	}
	return ret
}

// we have potentially added extra theoretical plates above
// now we make real plates and swap them in

func make_plates(request *LHRequest) {
	remap := make(map[string]string)
	for k, v := range request.LHInstructions {
		_, skip := remap[v.PlateID]

		if skip {
			request.LHInstructions[k].PlateID = remap[v.PlateID]
			continue
		}
		_, ok := request.Output_plates[v.PlateID]

		if !ok {
			plate := factory.GetPlateByType(v.Platetype)
			request.Output_plates[plate.ID] = plate
			remap[v.PlateID] = plate.ID
			request.LHInstructions[k].PlateID = remap[v.PlateID]
		}

	}
}

func make_layouts(request *LHRequest, pc []PlateChoice) {
	// we need to fill in the platechoice structure then
	// transfer the info across to the solutions

	opa := request.Output_assignments

	for _, c := range pc {
		// make a temporary plate to hold info

		plat := factory.GetPlateByType(c.Platetype)

		// make an iterator for it

		it := request.OutputIteratorFactory(plat)

		//seed in the existing assignments

		for _, w := range c.Wells {
			if w != "" {
				wc := wtype.MakeWellCoords(w)
				plat.Cols[wc.X][wc.Y].Currvol += 100.0
			}
		}

		for i, _ := range c.Assigned {
			sID := c.Assigned[i]
			well := c.Wells[i]

			var assignment string

			if well == "" {
				wc := plat.NextEmptyWell(it)

				if wc.IsZero() {
					// something very bad has happened
					logger.Fatal("DIRE WARNING: The unthinkable has happened... output plate has too many assignments!")
				}

				plat.Cols[wc.X][wc.Y].Currvol += 100.0
				request.LHInstructions[sID].Welladdress = wc.FormatA1()
				assignment = c.ID + ":" + wc.FormatA1()
			} else {
				assignment = c.ID + ":" + well
			}

			opa[assignment] = append(opa[assignment], sID)
		}
	}

	request.Output_assignments = opa
}
