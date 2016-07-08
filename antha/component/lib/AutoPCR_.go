package lib

import (
	"github.com/antha-lang/antha/antha/anthalib/wtype"
	"github.com/antha-lang/antha/antha/anthalib/wunit"
	"github.com/antha-lang/antha/bvendor/golang.org/x/net/context"
	"github.com/antha-lang/antha/execute"
	"github.com/antha-lang/antha/inject"
	"github.com/antha-lang/antha/microArch/factory"
)

// Input parameters for this protocol (data)

// PCRprep parameters
// e.g. ["left homology arm"]:"templatename"
// e.g. ["left homology arm"]:"fwdprimer","revprimer"

// Data which is returned from this protocol, and data types

// Physical Inputs to this protocol with types

// Physical outputs from this protocol with types

func _AutoPCRRequirements() {
}

// Conditions to run on startup
func _AutoPCRSetup(_ctx context.Context, _input *AutoPCRInput) {
}

// The core process for this protocol, with the steps to be performed
// for every input
func _AutoPCRSteps(_ctx context.Context, _input *AutoPCRInput, _output *AutoPCROutput) {
	_output.Reactions = make([]*wtype.LHComponent, 0)
	for reactionname, templatename := range _input.Reactiontotemplate {
		result := PCR_volRunSteps(_ctx, &PCR_volInput{ReactionVolume: wunit.NewVolume(25, "ul"),
			BufferConcinX:         5,
			FwdPrimerName:         _input.Reactiontoprimerpair[reactionname][0],
			RevPrimerName:         _input.Reactiontoprimerpair[reactionname][1],
			TemplateName:          templatename,
			ReactionName:          reactionname,
			FwdPrimerVol:          wunit.NewVolume(1, "ul"),
			RevPrimerVol:          wunit.NewVolume(1, "ul"),
			AdditiveVols:          []wunit.Volume{wunit.NewVolume(5, "ul")},
			Templatevolume:        wunit.NewVolume(1, "ul"),
			PolymeraseVolume:      wunit.NewVolume(1, "ul"),
			DNTPVol:               wunit.NewVolume(1, "ul"),
			Numberofcycles:        30,
			InitDenaturationtime:  wunit.NewTime(30, "s"),
			Denaturationtime:      wunit.NewTime(5, "s"),
			Annealingtime:         wunit.NewTime(10, "s"),
			AnnealingTemp:         wunit.NewTemperature(72, "C"), // Should be calculated from primer and template binding
			Extensiontime:         wunit.NewTime(60, "s"),        // should be calculated from template length and polymerase rate
			Finalextensiontime:    wunit.NewTime(180, "s"),
			Hotstart:              false,
			AddPrimerstoMasterMix: false,

			FwdPrimer:     _input.FwdPrimertype,
			RevPrimer:     _input.RevPrimertype,
			DNTPS:         factory.GetComponentByType("DNTPs"),
			PCRPolymerase: factory.GetComponentByType("Q5Polymerase"),
			Buffer:        factory.GetComponentByType("Q5buffer"),
			Water:         factory.GetComponentByType("water"),
			Template:      _input.Templatetype,
			Additives:     []*wtype.LHComponent{factory.GetComponentByType("GCenhancer")},
			OutPlate:      _input.Plate},
		)

		_output.Reactions = append(_output.Reactions, result.Outputs.Reaction)

	}
}

// Run after controls and a steps block are completed to
// post process any data and provide downstream results
func _AutoPCRAnalysis(_ctx context.Context, _input *AutoPCRInput, _output *AutoPCROutput) {
}

// A block of tests to perform to validate that the sample was processed correctly
// Optionally, destructive tests can be performed to validate results on a
// dipstick basis
func _AutoPCRValidation(_ctx context.Context, _input *AutoPCRInput, _output *AutoPCROutput) {
}
func _AutoPCRRun(_ctx context.Context, input *AutoPCRInput) *AutoPCROutput {
	output := &AutoPCROutput{}
	_AutoPCRSetup(_ctx, input)
	_AutoPCRSteps(_ctx, input, output)
	_AutoPCRAnalysis(_ctx, input, output)
	_AutoPCRValidation(_ctx, input, output)
	return output
}

func AutoPCRRunSteps(_ctx context.Context, input *AutoPCRInput) *AutoPCRSOutput {
	soutput := &AutoPCRSOutput{}
	output := _AutoPCRRun(_ctx, input)
	if err := inject.AssignSome(output, &soutput.Data); err != nil {
		panic(err)
	}
	if err := inject.AssignSome(output, &soutput.Outputs); err != nil {
		panic(err)
	}
	return soutput
}

func AutoPCRNew() interface{} {
	return &AutoPCRElement{
		inject.CheckedRunner{
			RunFunc: func(_ctx context.Context, value inject.Value) (inject.Value, error) {
				input := &AutoPCRInput{}
				if err := inject.Assign(value, input); err != nil {
					return nil, err
				}
				output := _AutoPCRRun(_ctx, input)
				return inject.MakeValue(output), nil
			},
			In:  &AutoPCRInput{},
			Out: &AutoPCROutput{},
		},
	}
}

var (
	_ = execute.MixInto
	_ = wunit.Make_units
)

type AutoPCRElement struct {
	inject.CheckedRunner
}

type AutoPCRInput struct {
	FwdPrimertype        *wtype.LHComponent
	Plate                *wtype.LHPlate
	Reactiontoprimerpair map[string][]string
	Reactiontotemplate   map[string]string
	RevPrimertype        *wtype.LHComponent
	Templatetype         *wtype.LHComponent
}

type AutoPCROutput struct {
	Reactions []*wtype.LHComponent
}

type AutoPCRSOutput struct {
	Data struct {
	}
	Outputs struct {
		Reactions []*wtype.LHComponent
	}
}

func init() {
	if err := addComponent(Component{Name: "AutoPCR",
		Constructor: AutoPCRNew,
		Desc: ComponentDesc{
			Desc: "",
			Path: "antha/component/an/Liquid_handling/PCR/AutoPCR.an",
			Params: []ParamDesc{
				{Name: "FwdPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Plate", Desc: "", Kind: "Inputs"},
				{Name: "Reactiontoprimerpair", Desc: "e.g. [\"left homology arm\"]:\"fwdprimer\",\"revprimer\"\n", Kind: "Parameters"},
				{Name: "Reactiontotemplate", Desc: "PCRprep parameters\n\ne.g. [\"left homology arm\"]:\"templatename\"\n", Kind: "Parameters"},
				{Name: "RevPrimertype", Desc: "", Kind: "Inputs"},
				{Name: "Templatetype", Desc: "", Kind: "Inputs"},
				{Name: "Reactions", Desc: "", Kind: "Outputs"},
			},
		},
	}); err != nil {
		panic(err)
	}
}
