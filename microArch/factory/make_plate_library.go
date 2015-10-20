// anthalib/factory/make_plate_library.go: Part of the Antha language
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

package factory

import "github.com/antha-lang/antha/antha/anthalib/wtype"

func makePlateLibrary() map[string]*wtype.LHPlate {
	plates := make(map[string]*wtype.LHPlate)

	swshp := wtype.NewShape("box", "mm", 8.2, 8.2, 41.3)
	welltype := wtype.NewLHWell("DSW96", "", "", "ul", 2000, 25, swshp, 3, 8.2, 8.2, 41.3, 4.7, "mm")
	plate := wtype.NewLHPlate("DSW96", "Unknown", 8, 12, 44.1, "mm", welltype, 9, 9, 0.0, 0.0, 0.0)
	plates[plate.Type] = plate

	rwshp := wtype.NewShape("cylinder", "mm", 8.2, 8.2, 11)
	welltype = wtype.NewLHWell("SRWFB96", "", "", "ul", 500, 10, rwshp, 0, 8.2, 8.2, 11, 1.0, "mm")
	plate = wtype.NewLHPlate("SRWFB96", "Unknown", 8, 12, 15, "mm", welltype, 9, 9, 0.0, 0.0, 0.0)
	plates[plate.Type] = plate

	stshp := wtype.NewShape("box", "mm", 8.2, 72, 41.3)
	welltype = wtype.NewLHWell("DWST12", "", "", "ul", 15000, 1000, stshp, 3, 8.2, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST12", "Unknown", 1, 12, 44.1, "mm", welltype, 9, 9, 0.0, 31.5, 0.0)
	plates[plate.Type] = plate

	stshp = wtype.NewShape("box", "mm", 115.0, 8.2, 41.3)
	welltype = wtype.NewLHWell("DWST8", "", "", "ul", 24000, 1000, stshp, 3, 115, 8.2, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWST8", "Unknown", 8, 1, 44.1, "mm", welltype, 9, 9, 49.5, 0.0, 0.0)
	plates[plate.Type] = plate

	rshp := wtype.NewShape("box", "mm", 115.0, 72.0, 41.3)
	welltype = wtype.NewLHWell("DWR1", "", "", "ul", 300000, 20000, rshp, 3, 115, 72, 41.3, 4.7, "mm")
	plate = wtype.NewLHPlate("DWR1", "Unknown", 1, 1, 44.1, "mm", welltype, 9, 9, 49.5, 31.5, 0.0)
	plates[plate.Type] = plate

	rwshp = wtype.NewShape("cylinder", "mm", 5.5, 5.5, 20.4)
	welltype = wtype.NewLHWell("pcrplate", "", "", "ul", 250, 5, rwshp, 0, 5.5, 5.5, 20.4, 1.4, "mm")
	//plate = wtype.NewLHPlate("pcrplate", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 6.5)
	//plates[plate.Type] = plate
	plate = wtype.NewLHPlate("pcrplate_with_cooler", "Unknown", 8, 12, 25.7, "mm", welltype, 9, 9, 0.0, 0.0, 15.5)
	plates[plate.Type] = plate
	return plates
}

func GetPlateByType(typ string) *wtype.LHPlate {
	plates := makePlateLibrary()
	p := plates[typ]
	return p.Dup()
}

func GetPlateList() []string {
	plates := makePlateLibrary()

	kz := make([]string, len(plates))
	x := 0
	for name, _ := range plates {
		kz[x] = name
		x += 1
	}
	return kz
}