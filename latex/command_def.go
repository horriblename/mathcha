// Definitions for known latex commands (basckslash commands)
// definitions are pulled from mathquill and pie-frameworks' extended version
package latex

import "strings"

// -----------------------------------------------------------------------------
// Definitions for Supported Commands

// Vanilla Symbols: commands that result in one character(rune) and have no
// special grammar rules

// Note on wierd indentation: it's to group them together
type LatexCmd int

const (
	CMD_UNKNOWN LatexCmd = iota

	cmd_1arg_beg // commands that expect 1 arguement
	CMD_sqrt
	cmd_1arg_end

	cmd_2arg_beg // commands that expect 2 arguements
	CMD_frac
	cmd_2arg_end

	vanilla_sym_beg
	CMD_complement
	CMD_nexists
	CMD_sphericalangle
	CMD_iint
	CMD_iiint
	CMD_oiint
	CMD_oiiint
	CMD_backsim
	CMD_backsimeq
	CMD_eqsim
	CMD_ncong
	CMD_approxeq
	CMD_bumpeq
	CMD_Bumpeq
	CMD_doteqdot
	CMD_fallingdotseq
	CMD_risingdotseq
	CMD_eqcirc
	CMD_circeq
	CMD_triangleq
	CMD_leqq
	CMD_geqq
	CMD_lneqq
	CMD_gneqq
	CMD_between
	CMD_nleq
	CMD_ngeq
	CMD_lesssim
	CMD_gtrsim
	CMD_lessgtr
	CMD_gtrless
	CMD_preccurlyeq
	CMD_succcurlyeq
	CMD_precsim
	CMD_succsim
	CMD_nprec
	CMD_nsucc
	CMD_subsetneq
	CMD_supsetneq
	CMD_vDash
	CMD_Vdash
	CMD_Vvdash
	CMD_VDash
	CMD_nvdash
	CMD_nvDash
	CMD_nVdash
	CMD_nVDash
	CMD_vartriangleleft
	CMD_vartriangleright
	CMD_trianglelefteq
	CMD_trianglerighteq
	CMD_multimap
	CMD_Subset
	CMD_Supset
	CMD_Cap
	CMD_Cup
	CMD_pitchfork
	CMD_lessdot
	CMD_gtrdot
	CMD_lll
	CMD_ggg
	CMD_lesseqgtr
	CMD_gtreqless
	CMD_curlyeqprec
	CMD_curlyeqsucc
	CMD_nsim
	CMD_lnsim
	CMD_gnsim
	CMD_precnsim
	CMD_succnsim
	CMD_ntriangleleft
	CMD_ntriangleright
	CMD_ntrianglelefteq
	CMD_ntrianglerighteq
	CMD_blacksquare
	CMD_colon
	CMD_llcorner
	CMD_dotplus
	CMD_intercal
	CMD_veebar
	CMD_barwedge
	CMD_ltimes
	CMD_rtimes
	CMD_leftthreetimes
	CMD_rightthreetimes
	CMD_curlyvee
	CMD_curlywedge
	CMD_circledcirc
	CMD_circledast
	CMD_circleddash
	CMD_boxplus
	CMD_boxminus
	CMD_boxtimes
	CMD_boxdot
	vanilla_sym_end
)

var latexCmds = [...]string{
	// functional commands
	CMD_sqrt: "\\sqrt",
	CMD_frac: "\\frac",
	// extended symbols pulled from pie-frameworks's mathquill repo
	CMD_complement:       "\\complement",
	CMD_nexists:          "\\nexists",
	CMD_sphericalangle:   "\\sphericalangle",
	CMD_iint:             "\\iint",
	CMD_iiint:            "\\iiint",
	CMD_oiint:            "\\oiint",
	CMD_oiiint:           "\\oiiint",
	CMD_backsim:          "\\backsim",
	CMD_backsimeq:        "\\backsimeq",
	CMD_eqsim:            "\\eqsim",
	CMD_ncong:            "\\ncong",
	CMD_approxeq:         "\\approxeq",
	CMD_bumpeq:           "\\bumpeq",
	CMD_Bumpeq:           "\\Bumpeq",
	CMD_doteqdot:         "\\doteqdot",
	CMD_fallingdotseq:    "\\fallingdotseq",
	CMD_risingdotseq:     "\\risingdotseq",
	CMD_eqcirc:           "\\eqcirc",
	CMD_circeq:           "\\circeq",
	CMD_triangleq:        "\\triangleq",
	CMD_leqq:             "\\leqq",
	CMD_geqq:             "\\geqq",
	CMD_lneqq:            "\\lneqq",
	CMD_gneqq:            "\\gneqq",
	CMD_between:          "\\between",
	CMD_nleq:             "\\nleq",
	CMD_ngeq:             "\\ngeq",
	CMD_lesssim:          "\\lesssim",
	CMD_gtrsim:           "\\gtrsim",
	CMD_lessgtr:          "\\lessgtr",
	CMD_gtrless:          "\\gtrless",
	CMD_preccurlyeq:      "\\preccurlyeq",
	CMD_succcurlyeq:      "\\succcurlyeq",
	CMD_precsim:          "\\precsim",
	CMD_succsim:          "\\succsim",
	CMD_nprec:            "\\nprec",
	CMD_nsucc:            "\\nsucc",
	CMD_subsetneq:        "\\subsetneq",
	CMD_supsetneq:        "\\supsetneq",
	CMD_vDash:            "\\vDash",
	CMD_Vdash:            "\\Vdash",
	CMD_Vvdash:           "\\Vvdash",
	CMD_VDash:            "\\VDash",
	CMD_nvdash:           "\\nvdash",
	CMD_nvDash:           "\\nvDash",
	CMD_nVdash:           "\\nVdash",
	CMD_nVDash:           "\\nVDash",
	CMD_vartriangleleft:  "\\vartriangleleft",
	CMD_vartriangleright: "\\vartriangleright",
	CMD_trianglelefteq:   "\\trianglelefteq",
	CMD_trianglerighteq:  "\\trianglerighteq",
	CMD_multimap:         "\\multimap",
	CMD_Subset:           "\\Subset",
	CMD_Supset:           "\\Supset",
	CMD_Cap:              "\\Cap",
	CMD_Cup:              "\\Cup",
	CMD_pitchfork:        "\\pitchfork",
	CMD_lessdot:          "\\lessdot",
	CMD_gtrdot:           "\\gtrdot",
	CMD_lll:              "\\lll",
	CMD_ggg:              "\\ggg",
	CMD_lesseqgtr:        "\\lesseqgtr",
	CMD_gtreqless:        "\\gtreqless",
	CMD_curlyeqprec:      "\\curlyeqprec",
	CMD_curlyeqsucc:      "\\curlyeqsucc",
	CMD_nsim:             "\\nsim",
	CMD_lnsim:            "\\lnsim",
	CMD_gnsim:            "\\gnsim",
	CMD_precnsim:         "\\precnsim",
	CMD_succnsim:         "\\succnsim",
	CMD_ntriangleleft:    "\\ntriangleleft",
	CMD_ntriangleright:   "\\ntriangleright",
	CMD_ntrianglelefteq:  "\\ntrianglelefteq",
	CMD_ntrianglerighteq: "\\ntrianglerighteq",
	CMD_blacksquare:      "\\blacksquare",
	CMD_colon:            "\\colon",
	CMD_llcorner:         "\\llcorner",
	CMD_dotplus:          "\\dotplus",
	CMD_intercal:         "\\intercal",
	CMD_veebar:           "\\veebar",
	CMD_barwedge:         "\\barwedge",
	CMD_ltimes:           "\\ltimes",
	CMD_rtimes:           "\\rtimes",
	CMD_leftthreetimes:   "\\leftthreetimes",
	CMD_rightthreetimes:  "\\rightthreetimes",
	CMD_curlyvee:         "\\curlyvee",
	CMD_curlywedge:       "\\curlywedge",
	CMD_circledcirc:      "\\circledcirc",
	CMD_circledast:       "\\circledast",
	CMD_circleddash:      "\\circleddash",
	CMD_boxplus:          "\\boxplus",
	CMD_boxminus:         "\\boxminus",
	CMD_boxtimes:         "\\boxtimes",
	CMD_boxdot:           "\\boxdot",
}

func (cmd LatexCmd) GetCmd() string { return latexCmds[cmd] }

// FIXME if we don't need to map from constant to string command, swap key-values
// in the latexCmds map then change this function
func MatchLatexCmd(cmd string) LatexCmd {
	cmd = strings.TrimSpace(cmd)
	for k, v := range latexCmds {
		if v == cmd {
			return LatexCmd(k)
		}
	}
	return CMD_UNKNOWN
}

func (cmd LatexCmd) IsVanillaSym() bool {
	return vanilla_sym_beg < cmd && cmd < vanilla_sym_end
}

func (cmd LatexCmd) TakesOneArg() bool {
	return cmd_1arg_beg < cmd && cmd < cmd_1arg_end
}

func (cmd LatexCmd) TakesTwoArg() bool {
	return cmd_2arg_beg < cmd && cmd < cmd_2arg_end
}
