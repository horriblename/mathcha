// Definitions for known latex commands (basckslash commands)
// definitions are pulled from mathquill and pie-frameworks' extended version
package latex

// -----------------------------------------------------------------------------
// Definitions for Supported Commands

// Vanilla Symbols: commands that result in one character(rune) and have no
// special grammar rules

// Note on wierd indentation: it's to group them together
type LatexCmd int

const (
	CMD_UNKNOWN  LatexCmd = iota
	cmd_text_beg          // temporary group name for commands that take a 'raw' string as parameter
	// text formatting
	CMD_text
	// CMD_textnormal
	// CMD_textrm
	// CMD_textup
	// CMD_textmd
	// CMD_emph
	// CMD_italic
	// CMD_textit
	// CMD_textbf
	// CMD_textsf
	// CMD_texttt
	// CMD_textsc
	// CMD_uppercase
	cmd_text_end

	cmd_1arg_beg // commands that expect 1 arguement
	// accents
	CMD_underline
	CMD_overline
	CMD_subscript
	CMD_superscript
	// CMD_lowercase
	CMD_sqrt
	cmd_1arg_end

	cmd_2arg_beg // commands that expect 2 arguements
	CMD_binom
	CMD_frac
	cmd_2arg_end

	cmd_enclosing_beg
	CMD_left
	CMD_right
	cmd_enclosing_end

	vanilla_sym_beg
	CMD_alpha
	CMD_beta
	CMD_gamma
	CMD_delta
	CMD_zeta
	CMD_eta
	CMD_theta
	CMD_iota
	CMD_kappa
	CMD_mu
	CMD_nu
	CMD_xi
	CMD_rho
	CMD_sigma
	CMD_tau
	CMD_chi
	CMD_psi

	CMD_phi
	CMD_phiv
	CMD_varphi

	CMD_epsilon
	CMD_epsiv
	CMD_varepsilon

	CMD_piv
	CMD_varpi

	CMD_sigmaf
	CMD_sigmav
	CMD_varsigma

	CMD_thetav
	CMD_vartheta
	CMD_thetasym

	CMD_upsilon
	CMD_upsi

	CMD_gammad
	CMD_Gammad
	CMD_digamma

	CMD_kappav
	CMD_varkappa

	CMD_rhov
	CMD_varrho

	CMD_pi
	CMD_lambda

	CMD_Upsilon
	CMD_Upsi
	CMD_upsih
	CMD_Upsih

	CMD_Gamma
	CMD_Delta
	CMD_Theta
	CMD_Lambda
	CMD_Xi
	CMD_Pi
	CMD_Sigma
	CMD_Phi
	CMD_Psi
	CMD_Omega

	CMD_cdot
	CMD_sim
	CMD_cong
	CMD_equiv
	CMD_oplus
	CMD_otimes
	CMD_times
	CMD_div
	CMD_ne
	CMD_ast
	CMD_therefor
	CMD_cuz

	CMD_prop
	CMD_asymp

	CMD_lt
	CMD_gt
	CMD_le
	CMD_ge
	CMD_isin
	CMD_notin
	CMD_ni
	CMD_notni

	CMD_sub
	CMD_sup
	CMD_nsub
	CMD_nsup

	CMD_sube
	CMD_supe

	CMD_nsube
	CMD_nsupe

	CMD_sum
	CMD_prod
	CMD_coprod
	CMD_int

	CMD_N
	CMD_P
	CMD_Z
	CMD_Q
	CMD_R
	CMD_C
	CMD_H

	// spacing
	CMD_SPACE
	CMD_quad
	CMD_emsp
	CMD_qquad

	CMD_diamond
	CMD_bigtriangleup
	CMD_ominus
	CMD_uplus
	CMD_bigtriangledown
	CMD_sqcap
	CMD_triangleleft
	CMD_sqcup
	CMD_triangleright
	CMD_odot
	CMD_bigcirc
	CMD_dagger
	CMD_ddagger
	CMD_wr
	CMD_amalg

	CMD_models
	CMD_prec
	CMD_succ
	CMD_preceq
	CMD_succeq
	CMD_simeq
	CMD_mid
	CMD_ll
	CMD_gg
	CMD_parallel
	CMD_bowtie
	CMD_sqsubset
	CMD_sqsupset
	CMD_smile
	CMD_sqsubseteq
	CMD_sqsupseteq
	CMD_doteq
	CMD_frown
	CMD_vdash
	CMD_dashv

	CMD_longleftarrow
	CMD_longrightarrow
	CMD_Longleftarrow
	CMD_Longrightarrow
	CMD_longleftrightarrow
	CMD_updownarrow
	CMD_Longleftrightarrow
	CMD_Updownarrow
	CMD_mapsto
	CMD_nearrow
	CMD_hookleftarrow
	CMD_hookrightarrow
	CMD_searrow
	CMD_leftharpoonup
	CMD_rightharpoonup
	CMD_swarrow
	CMD_leftharpoondown
	CMD_rightharpoondown
	CMD_nwarrow

	CMD_ldots
	CMD_cdots
	CMD_vdots
	CMD_ddots
	CMD_surd
	CMD_triangle
	CMD_ell
	CMD_top
	CMD_flat
	CMD_natural
	CMD_sharp
	CMD_wp
	CMD_bot
	CMD_clubsuit
	CMD_diamondsuit
	CMD_heartsuit
	CMD_spadesuit

	CMD_oint
	CMD_bigcap
	CMD_bigcup
	CMD_bigsqcup
	CMD_bigvee
	CMD_bigwedge
	CMD_bigodot
	CMD_bigotimes
	CMD_bigoplus
	CMD_biguplus

	CMD_lfloor
	CMD_rfloor
	CMD_lceil
	CMD_rceil
	CMD_slash
	CMD_opencurlybrace
	CMD_closecurlybrace

	CMD_caret
	CMD_underscore
	CMD_backslash
	CMD_vert
	CMD_perp
	CMD_nabla
	CMD_hbar
	CMD_AA
	CMD_ring
	CMD_bull
	CMD_setminus
	CMD_not
	CMD_dots

	CMD_converges
	CMD_dArr
	CMD_diverges
	CMD_uArr
	CMD_to
	CMD_implies
	CMD_gets
	CMD_impliedby
	CMD_harr
	CMD_iff

	CMD_Re
	CMD_Im
	CMD_part

	CMD_inf
	CMD_alef

	CMD_forall
	CMD_xist
	CMD_and
	CMD_or

	CMD_o
	CMD_cup
	CMD_cap

	CMD_deg
	CMD_ang

	CMD_ln
	CMD_lg
	CMD_log
	CMD_span
	CMD_proj
	CMD_det
	CMD_dim
	CMD_min
	CMD_max
	CMD_mod
	CMD_lcm
	CMD_gcd
	CMD_gcf
	CMD_hcf
	CMD_lim

	CMD_sin
	CMD_cos
	CMD_tan
	CMD_sec
	CMD_cosec
	CMD_csc
	CMD_cotan
	CMD_cot

	CMD_sinh
	CMD_cosh
	CMD_tanh
	CMD_sech
	CMD_cosech
	CMD_csch
	CMD_cotanh
	CMD_coth

	CMD_asin
	CMD_acos
	CMD_atan
	CMD_asec
	CMD_acosec
	CMD_acsc
	CMD_acotan
	CMD_acot

	CMD_asinh
	CMD_acosh
	CMD_atanh
	CMD_asech
	CMD_acosech
	CMD_acsch
	CMD_acotanh
	CMD_acoth

	CMD_arcsin
	CMD_arccos
	CMD_arctan
	CMD_arcsec
	CMD_arccosec
	CMD_arccsc
	CMD_arccotan
	CMD_arccot

	CMD_arcsinh
	CMD_arccosh
	CMD_arctanh
	CMD_arcsech
	CMD_arccosech
	CMD_arccsch
	CMD_arccotanh
	CMD_arccoth

	// extended symbols by pie framework
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
	CMD_nmid
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

var latexCmds = map[string]LatexCmd{
	// functional commands
	`\text`: CMD_text,
	// accents
	`\underline`:   CMD_underline,
	`\overline`:    CMD_overline,
	`\subscript`:   CMD_subscript,
	`\superscript`: CMD_superscript,
	`_`:            CMD_subscript,
	`^`:            CMD_superscript,
	// text formatting

	`\left`:  CMD_left,
	`\right`: CMD_right,

	`\sqrt`: CMD_sqrt,

	// 2 parameter commands
	`\binom`:    CMD_binom,
	`\frac`:     CMD_frac,
	`\dfrac`:    CMD_frac,
	`\cfrac`:    CMD_frac,
	`\fraction`: CMD_frac,

	"\\alpha":      CMD_alpha,
	"\\beta":       CMD_beta,
	"\\gamma":      CMD_gamma,
	"\\delta":      CMD_delta,
	"\\zeta":       CMD_zeta,
	"\\eta":        CMD_eta,
	"\\theta":      CMD_theta,
	"\\iota":       CMD_iota,
	"\\kappa":      CMD_kappa,
	"\\mu":         CMD_mu,
	"\\nu":         CMD_nu,
	"\\xi":         CMD_xi,
	"\\rho":        CMD_rho,
	"\\sigma":      CMD_sigma,
	"\\tau":        CMD_tau,
	"\\chi":        CMD_chi,
	"\\psi":        CMD_psi,
	"\\phi":        CMD_phi,
	"\\phiv":       CMD_phiv,
	"\\varphi":     CMD_varphi,
	"\\epsilon":    CMD_epsilon,
	"\\epsiv":      CMD_epsiv,
	"\\varepsilon": CMD_varepsilon,
	"\\piv":        CMD_piv,
	"\\varpi":      CMD_varpi,
	"\\sigmaf":     CMD_sigmaf,
	"\\sigmav":     CMD_sigmav,
	"\\varsigma":   CMD_varsigma,
	"\\thetav":     CMD_thetav,
	"\\vartheta":   CMD_vartheta,
	"\\thetasym":   CMD_thetasym,
	"\\upsilon":    CMD_upsilon,
	"\\upsi":       CMD_upsi,
	"\\gammad":     CMD_gammad,
	"\\Gammad":     CMD_Gammad,
	"\\digamma":    CMD_digamma,
	"\\kappav":     CMD_kappav,
	"\\varkappa":   CMD_varkappa,
	"\\rhov":       CMD_rhov,
	"\\varrho":     CMD_varrho,
	"\\pi":         CMD_pi,
	"\\lambda":     CMD_lambda,
	"\\Upsilon":    CMD_Upsilon,
	"\\Upsi":       CMD_Upsi,
	"\\upsih":      CMD_upsih,
	"\\Upsih":      CMD_Upsih,
	"\\Gamma":      CMD_Gamma,
	"\\Delta":      CMD_Delta,
	"\\Theta":      CMD_Theta,
	"\\Lambda":     CMD_Lambda,
	"\\Xi":         CMD_Xi,
	"\\Pi":         CMD_Pi,
	"\\Sigma":      CMD_Sigma,
	"\\Phi":        CMD_Phi,
	"\\Psi":        CMD_Psi,
	"\\Omega":      CMD_Omega,

	"\\cdot ":              CMD_cdot,
	"\\sdot":               CMD_cdot,
	"\\sim":                CMD_sim,
	"\\cong":               CMD_cong,
	"\\equiv":              CMD_equiv,
	"\\oplus":              CMD_oplus,
	"\\otimes":             CMD_otimes,
	"\\times":              CMD_times,
	"\\div":                CMD_div,
	"\\divide":             CMD_div,
	"\\divides":            CMD_div,
	"\\ne":                 CMD_ne,
	"\\neq":                CMD_ne,
	"\\ast":                CMD_ast,
	"\\star":               CMD_ast,
	"\\loast":              CMD_ast,
	"\\lowast":             CMD_ast,
	"\\therefor":           CMD_therefor,
	"\\therefore":          CMD_therefor,
	"\\cuz":                CMD_cuz,
	"\\because":            CMD_cuz,
	"\\prop":               CMD_prop,
	"\\propto":             CMD_prop,
	"\\asymp":              CMD_asymp,
	"\\approx":             CMD_asymp,
	"\\lt":                 CMD_lt,
	"\\gt":                 CMD_gt,
	"\\le":                 CMD_le,
	"\\leq":                CMD_le,
	"\\ge":                 CMD_ge,
	"\\geq":                CMD_ge,
	"\\isin":               CMD_isin,
	"\\in":                 CMD_isin,
	"\\notin":              CMD_notin,
	"\\ni":                 CMD_ni,
	"\\contains":           CMD_ni,
	"\\notni":              CMD_notni,
	"\\niton":              CMD_notni,
	"\\notcontains":        CMD_notni,
	"\\doesnotcontain":     CMD_notni,
	"\\sub":                CMD_sub,
	"\\subset":             CMD_sub,
	"\\sup":                CMD_sup,
	"\\supset":             CMD_sup,
	"\\superset":           CMD_sup,
	"\\nsub":               CMD_nsub,
	"\\notsub":             CMD_nsub,
	"\\nsubset":            CMD_nsub,
	"\\notsubset":          CMD_nsub,
	"\\nsup":               CMD_nsup,
	"\\notsup":             CMD_nsup,
	"\\nsupset":            CMD_nsup,
	"\\notsupset":          CMD_nsup,
	"\\nsuperset":          CMD_nsup,
	"\\notsuperset":        CMD_nsup,
	"\\sube":               CMD_sube,
	"\\subeq":              CMD_sube,
	"\\subsete":            CMD_sube,
	"\\subseteq":           CMD_sube,
	"\\supe":               CMD_supe,
	"\\supeq":              CMD_supe,
	"\\supsete":            CMD_supe,
	"\\supseteq":           CMD_supe,
	"\\supersete":          CMD_supe,
	"\\superseteq":         CMD_supe,
	"\\nsube":              CMD_nsube,
	"\\nsubeq":             CMD_nsube,
	"\\notsube":            CMD_nsube,
	"\\notsubeq":           CMD_nsube,
	"\\nsubsete":           CMD_nsube,
	"\\nsubseteq":          CMD_nsube,
	"\\notsubsete":         CMD_nsube,
	"\\notsubseteq":        CMD_nsube,
	"\\nsupe":              CMD_nsupe,
	"\\nsupeq":             CMD_nsupe,
	"\\notsupe":            CMD_nsupe,
	"\\notsupeq":           CMD_nsupe,
	"\\nsupsete":           CMD_nsupe,
	"\\nsupseteq":          CMD_nsupe,
	"\\notsupsete":         CMD_nsupe,
	"\\notsupseteq":        CMD_nsupe,
	"\\nsupersete":         CMD_nsupe,
	"\\nsuperseteq":        CMD_nsupe,
	"\\notsupersete":       CMD_nsupe,
	"\\notsuperseteq":      CMD_nsupe,
	"\\sum":                CMD_sum,
	"\\summation":          CMD_sum,
	"\\prod":               CMD_prod,
	"\\product":            CMD_prod,
	"\\coprod":             CMD_coprod,
	"\\coproduct":          CMD_coprod,
	"\\int":                CMD_int,
	"\\integral":           CMD_int,
	"\\N":                  CMD_N,
	"\\naturals":           CMD_N,
	"\\Naturals":           CMD_N,
	"\\P":                  CMD_P,
	"\\primes":             CMD_P,
	"\\Primes":             CMD_P,
	"\\projective":         CMD_P,
	"\\Projective":         CMD_P,
	"\\probability":        CMD_P,
	"\\Probability":        CMD_P,
	"\\Z":                  CMD_Z,
	"\\integers":           CMD_Z,
	"\\Integers":           CMD_Z,
	"\\Q":                  CMD_Q,
	"\\rationals":          CMD_Q,
	"\\Rationals":          CMD_Q,
	"\\R":                  CMD_R,
	"\\reals":              CMD_R,
	"\\Reals":              CMD_R,
	"\\C":                  CMD_C,
	"\\complex":            CMD_C,
	"\\Complex":            CMD_C,
	"\\complexes":          CMD_C,
	"\\Complexes":          CMD_C,
	"\\complexplane":       CMD_C,
	"\\Complexplane":       CMD_C,
	"\\ComplexPlane":       CMD_C,
	"\\H":                  CMD_H,
	"\\Hamiltonian":        CMD_H,
	"\\quaternions":        CMD_H,
	"\\Quaternions":        CMD_H,
	"\\ ":                  CMD_SPACE,
	"\\quad":               CMD_quad,
	"\\emsp":               CMD_emsp,
	"\\qquad":              CMD_qquad,
	"\\diamond":            CMD_diamond,
	"\\bigtriangleup":      CMD_bigtriangleup,
	"\\ominus":             CMD_ominus,
	"\\uplus":              CMD_uplus,
	"\\bigtriangledown":    CMD_bigtriangledown,
	"\\sqcap":              CMD_sqcap,
	"\\triangleleft":       CMD_triangleleft,
	"\\sqcup":              CMD_sqcup,
	"\\triangleright":      CMD_triangleright,
	"\\odot":               CMD_odot,
	"\\bigcirc":            CMD_bigcirc,
	"\\dagger":             CMD_dagger,
	"\\ddagger":            CMD_ddagger,
	"\\wr":                 CMD_wr,
	"\\amalg":              CMD_amalg,
	"\\models":             CMD_models,
	"\\prec":               CMD_prec,
	"\\succ":               CMD_succ,
	"\\preceq":             CMD_preceq,
	"\\succeq":             CMD_succeq,
	"\\simeq":              CMD_simeq,
	"\\mid":                CMD_mid,
	"\\ll":                 CMD_ll,
	"\\gg":                 CMD_gg,
	"\\parallel":           CMD_parallel,
	"\\bowtie":             CMD_bowtie,
	"\\sqsubset":           CMD_sqsubset,
	"\\sqsupset":           CMD_sqsupset,
	"\\smile":              CMD_smile,
	"\\sqsubseteq":         CMD_sqsubseteq,
	"\\sqsupseteq":         CMD_sqsupseteq,
	"\\doteq":              CMD_doteq,
	"\\frown":              CMD_frown,
	"\\vdash":              CMD_vdash,
	"\\dashv":              CMD_dashv,
	"\\longleftarrow":      CMD_longleftarrow,
	"\\longrightarrow":     CMD_longrightarrow,
	"\\Longleftarrow":      CMD_Longleftarrow,
	"\\Longrightarrow":     CMD_Longrightarrow,
	"\\longleftrightarrow": CMD_longleftrightarrow,
	"\\updownarrow":        CMD_updownarrow,
	"\\Longleftrightarrow": CMD_Longleftrightarrow,
	"\\Updownarrow":        CMD_Updownarrow,
	"\\mapsto":             CMD_mapsto,
	"\\nearrow":            CMD_nearrow,
	"\\hookleftarrow":      CMD_hookleftarrow,
	"\\hookrightarrow":     CMD_hookrightarrow,
	"\\searrow":            CMD_searrow,
	"\\leftharpoonup":      CMD_leftharpoonup,
	"\\rightharpoonup":     CMD_rightharpoonup,
	"\\swarrow":            CMD_swarrow,
	"\\leftharpoondown":    CMD_leftharpoondown,
	"\\rightharpoondown":   CMD_rightharpoondown,
	"\\nwarrow":            CMD_nwarrow,
	"\\ldots":              CMD_ldots,
	"\\cdots":              CMD_cdots,
	"\\vdots":              CMD_vdots,
	"\\ddots":              CMD_ddots,
	"\\surd":               CMD_surd,
	"\\triangle":           CMD_triangle,
	"\\ell":                CMD_ell,
	"\\top":                CMD_top,
	"\\flat":               CMD_flat,
	"\\natural":            CMD_natural,
	"\\sharp":              CMD_sharp,
	"\\wp":                 CMD_wp,
	"\\bot":                CMD_bot,
	"\\clubsuit":           CMD_clubsuit,
	"\\diamondsuit":        CMD_diamondsuit,
	"\\heartsuit":          CMD_heartsuit,
	"\\spadesuit":          CMD_spadesuit,
	"\\oint":               CMD_oint,
	"\\bigcap":             CMD_bigcap,
	"\\bigcup":             CMD_bigcup,
	"\\bigsqcup":           CMD_bigsqcup,
	"\\bigvee":             CMD_bigvee,
	"\\bigwedge":           CMD_bigwedge,
	"\\bigodot":            CMD_bigodot,
	"\\bigotimes":          CMD_bigotimes,
	"\\bigoplus":           CMD_bigoplus,
	"\\biguplus":           CMD_biguplus,
	"\\lfloor":             CMD_lfloor,
	"\\rfloor":             CMD_rfloor,
	"\\lceil":              CMD_lceil,
	"\\rceil":              CMD_rceil,
	"\\slash":              CMD_slash,
	"\\opencurlybrace":     CMD_opencurlybrace,
	"\\closecurlybrace":    CMD_closecurlybrace,
	"\\caret":              CMD_caret,
	"\\underscore":         CMD_underscore,
	"\\backslash":          CMD_backslash,
	"\\vert":               CMD_vert,
	"\\perp":               CMD_perp,
	"\\perpendicular":      CMD_perp,
	"\\nabla":              CMD_nabla,
	"\\del":                CMD_nabla,
	"\\hbar":               CMD_hbar,
	"\\AA":                 CMD_AA,
	"\\Angstrom":           CMD_AA,
	"\\angstrom":           CMD_AA,
	"\\ring":               CMD_ring,
	"\\circ":               CMD_ring,
	"\\circle":             CMD_ring,
	"\\bull":               CMD_bull,
	"\\bullet":             CMD_bull,
	"\\setminus":           CMD_setminus,
	"\\smallsetminus":      CMD_setminus,
	"\\not":                CMD_not,
	"\\neg":                CMD_not,
	"\\dots":               CMD_dots,
	"\\ellip":              CMD_dots,
	"\\hellip":             CMD_dots,
	"\\ellipsis":           CMD_dots,
	"\\hellipsis":          CMD_dots,
	"\\converges":          CMD_converges,
	"\\darr":               CMD_converges,
	"\\dnarr":              CMD_converges,
	"\\dnarrow":            CMD_converges,
	"\\downarrow":          CMD_converges,
	"\\dArr":               CMD_dArr,
	"\\dnArr":              CMD_dArr,
	"\\dnArrow":            CMD_dArr,
	"\\Downarrow":          CMD_dArr,
	"\\diverges":           CMD_diverges,
	"\\uarr":               CMD_diverges,
	"\\uparrow":            CMD_diverges,
	"\\uArr":               CMD_uArr,
	"\\Uparrow":            CMD_uArr,
	"\\to":                 CMD_to,
	"\\rarr":               CMD_to,
	"\\rightarrow":         CMD_to,
	"\\implies":            CMD_implies,
	"\\rArr":               CMD_implies,
	"\\Rightarrow":         CMD_implies,
	"\\gets":               CMD_gets,
	"\\larr":               CMD_gets,
	"\\leftarrow":          CMD_gets,
	"\\impliedby":          CMD_impliedby,
	"\\lArr":               CMD_impliedby,
	"\\Leftarrow":          CMD_impliedby,
	"\\harr":               CMD_harr,
	"\\lrarr":              CMD_harr,
	"\\leftrightarrow":     CMD_harr,
	"\\iff":                CMD_iff,
	"\\hArr":               CMD_iff,
	"\\lrArr":              CMD_iff,
	"\\Leftrightarrow":     CMD_iff,
	"\\Re":                 CMD_Re,
	"\\Real":               CMD_Re,
	"\\real":               CMD_Re,
	"\\Im":                 CMD_Im,
	"\\imag":               CMD_Im,
	"\\image":              CMD_Im,
	"\\imagin":             CMD_Im,
	"\\imaginary":          CMD_Im,
	"\\Imaginary":          CMD_Im,
	"\\part":               CMD_part,
	"\\partial":            CMD_part,
	"\\inf":                CMD_inf,
	"\\infin":              CMD_inf,
	"\\infty":              CMD_inf,
	"\\infinity":           CMD_inf,
	"\\alef":               CMD_alef,
	"\\alefsym":            CMD_alef,
	"\\aleph":              CMD_alef,
	"\\alephsym":           CMD_alef,
	"\\forall":             CMD_forall,
	"\\xist":               CMD_xist,
	"\\xists":              CMD_xist,
	"\\exist":              CMD_xist,
	"\\exists":             CMD_xist,
	"\\and":                CMD_and,
	"\\land":               CMD_and,
	"\\wedge":              CMD_and,
	"\\or":                 CMD_or,
	"\\lor":                CMD_or,
	"\\vee":                CMD_or,
	"\\o":                  CMD_o,
	"\\O":                  CMD_o,
	"\\empty":              CMD_o,
	"\\emptyset":           CMD_o,
	"\\oslash":             CMD_o,
	"\\Oslash":             CMD_o,
	"\\nothing":            CMD_o,
	"\\varnothing":         CMD_o,
	"\\cup":                CMD_cup,
	"\\union":              CMD_cup,
	"\\cap":                CMD_cap,
	"\\intersect":          CMD_cap,
	"\\intersection":       CMD_cap,
	"\\deg":                CMD_deg,
	"\\degree":             CMD_deg,
	"\\ang":                CMD_ang,
	"\\angle":              CMD_ang,
	"\\ln":                 CMD_ln,
	"\\lg":                 CMD_lg,
	"\\log":                CMD_log,
	"\\span":               CMD_span,
	"\\proj":               CMD_proj,
	"\\det":                CMD_det,
	"\\dim":                CMD_dim,
	"\\min":                CMD_min,
	"\\max":                CMD_max,
	"\\mod":                CMD_mod,
	"\\lcm":                CMD_lcm,
	"\\gcd":                CMD_gcd,
	"\\gcf":                CMD_gcf,
	"\\hcf":                CMD_hcf,
	"\\lim":                CMD_lim,
	"\\sin":                CMD_sin,
	"\\cos":                CMD_cos,
	"\\tan":                CMD_tan,
	"\\sec":                CMD_sec,
	"\\cosec":              CMD_cosec,
	"\\csc":                CMD_csc,
	"\\cotan":              CMD_cotan,
	"\\cot":                CMD_cot,
	"\\sinh":               CMD_sinh,
	"\\cosh":               CMD_cosh,
	"\\tanh":               CMD_tanh,
	"\\sech":               CMD_sech,
	"\\cosech":             CMD_cosech,
	"\\csch":               CMD_csch,
	"\\cotanh":             CMD_cotanh,
	"\\coth":               CMD_coth,
	"\\asin":               CMD_asin,
	"\\acos":               CMD_acos,
	"\\atan":               CMD_atan,
	"\\asec":               CMD_asec,
	"\\acosec":             CMD_acosec,
	"\\acsc":               CMD_acsc,
	"\\acotan":             CMD_acotan,
	"\\acot":               CMD_acot,
	"\\asinh":              CMD_asinh,
	"\\acosh":              CMD_acosh,
	"\\atanh":              CMD_atanh,
	"\\asech":              CMD_asech,
	"\\acosech":            CMD_acosech,
	"\\acsch":              CMD_acsch,
	"\\acotanh":            CMD_acotanh,
	"\\acoth":              CMD_acoth,
	"\\arcsin":             CMD_arcsin,
	"\\arccos":             CMD_arccos,
	"\\arctan":             CMD_arctan,
	"\\arcsec":             CMD_arcsec,
	"\\arccosec":           CMD_arccosec,
	"\\arccsc":             CMD_arccsc,
	"\\arccotan":           CMD_arccotan,
	"\\arccot":             CMD_arccot,
	"\\arcsinh":            CMD_arcsinh,
	"\\arccosh":            CMD_arccosh,
	"\\arctanh":            CMD_arctanh,
	"\\arcsech":            CMD_arcsech,
	"\\arccosech":          CMD_arccosech,
	"\\arccsch":            CMD_arccsch,
	"\\arccotanh":          CMD_arccotanh,
	"\\arccoth":            CMD_arccoth,
	// extended symbols from pie-frameworks's mathquill repo
	"\\complement":       CMD_complement,
	"\\nexists":          CMD_nexists,
	"\\sphericalangle":   CMD_sphericalangle,
	"\\iint":             CMD_iint,
	"\\iiint":            CMD_iiint,
	"\\oiint":            CMD_oiint,
	"\\oiiint":           CMD_oiiint,
	"\\backsim":          CMD_backsim,
	"\\backsimeq":        CMD_backsimeq,
	"\\eqsim":            CMD_eqsim,
	"\\ncong":            CMD_ncong,
	"\\approxeq":         CMD_approxeq,
	"\\bumpeq":           CMD_bumpeq,
	"\\Bumpeq":           CMD_Bumpeq,
	"\\doteqdot":         CMD_doteqdot,
	"\\fallingdotseq":    CMD_fallingdotseq,
	"\\risingdotseq":     CMD_risingdotseq,
	"\\eqcirc":           CMD_eqcirc,
	"\\circeq":           CMD_circeq,
	"\\triangleq":        CMD_triangleq,
	"\\leqq":             CMD_leqq,
	"\\geqq":             CMD_geqq,
	"\\lneqq":            CMD_lneqq,
	"\\gneqq":            CMD_gneqq,
	"\\between":          CMD_between,
	"\\nleq":             CMD_nleq,
	"\\ngeq":             CMD_ngeq,
	"\\lesssim":          CMD_lesssim,
	"\\gtrsim":           CMD_gtrsim,
	"\\lessgtr":          CMD_lessgtr,
	"\\gtrless":          CMD_gtrless,
	"\\preccurlyeq":      CMD_preccurlyeq,
	"\\succcurlyeq":      CMD_succcurlyeq,
	"\\precsim":          CMD_precsim,
	"\\succsim":          CMD_succsim,
	"\\nprec":            CMD_nprec,
	"\\nsucc":            CMD_nsucc,
	"\\subsetneq":        CMD_subsetneq,
	"\\supsetneq":        CMD_supsetneq,
	"\\vDash":            CMD_vDash,
	"\\Vdash":            CMD_Vdash,
	"\\Vvdash":           CMD_Vvdash,
	"\\VDash":            CMD_VDash,
	"\\nvdash":           CMD_nvdash,
	"\\nvDash":           CMD_nvDash,
	"\\nVdash":           CMD_nVdash,
	"\\nVDash":           CMD_nVDash,
	"\\vartriangleleft":  CMD_vartriangleleft,
	"\\vartriangleright": CMD_vartriangleright,
	"\\trianglelefteq":   CMD_trianglelefteq,
	"\\trianglerighteq":  CMD_trianglerighteq,
	"\\multimap":         CMD_multimap,
	"\\Subset":           CMD_Subset,
	"\\Supset":           CMD_Supset,
	"\\Cap":              CMD_Cap,
	"\\Cup":              CMD_Cup,
	"\\pitchfork":        CMD_pitchfork,
	"\\lessdot":          CMD_lessdot,
	"\\gtrdot":           CMD_gtrdot,
	"\\lll":              CMD_lll,
	"\\ggg":              CMD_ggg,
	"\\lesseqgtr":        CMD_lesseqgtr,
	"\\gtreqless":        CMD_gtreqless,
	"\\curlyeqprec":      CMD_curlyeqprec,
	"\\curlyeqsucc":      CMD_curlyeqsucc,
	"\\nsim":             CMD_nsim,
	"\\lnsim":            CMD_lnsim,
	"\\gnsim":            CMD_gnsim,
	"\\precnsim":         CMD_precnsim,
	"\\succnsim":         CMD_succnsim,
	"\\ntriangleleft":    CMD_ntriangleleft,
	"\\ntriangleright":   CMD_ntriangleright,
	"\\ntrianglelefteq":  CMD_ntrianglelefteq,
	"\\ntrianglerighteq": CMD_ntrianglerighteq,
	"\\blacksquare":      CMD_blacksquare,
	"\\colon":            CMD_colon,
	"\\llcorner":         CMD_llcorner,
	"\\dotplus":          CMD_dotplus,
	"\\nmid":             CMD_nmid,
	"\\intercal":         CMD_intercal,
	"\\veebar":           CMD_veebar,
	"\\barwedge":         CMD_barwedge,
	"\\ltimes":           CMD_ltimes,
	"\\rtimes":           CMD_rtimes,
	"\\leftthreetimes":   CMD_leftthreetimes,
	"\\rightthreetimes":  CMD_rightthreetimes,
	"\\curlyvee":         CMD_curlyvee,
	"\\curlywedge":       CMD_curlywedge,
	"\\circledcirc":      CMD_circledcirc,
	"\\circledast":       CMD_circledast,
	"\\circleddash":      CMD_circleddash,
	"\\boxplus":          CMD_boxplus,
	"\\boxminus":         CMD_boxminus,
	"\\boxtimes":         CMD_boxtimes,
	"\\boxdot":           CMD_boxdot,
}

// BUG map variables are unordered, this returns a different string everytime if multiple commands are available
func (cmd LatexCmd) GetCmd() string {
	for k, v := range latexCmds {
		if v == cmd {
			return k
		}
	}
	return "unmapped command"
}

func MatchLatexCmd(cmd string) LatexCmd {
	return latexCmds[cmd]
}

func (cmd LatexCmd) TakesRawStrArg() bool {
	return cmd_text_beg < cmd && cmd < cmd_text_end
}

func (cmd LatexCmd) IsVanillaSym() bool {
	return vanilla_sym_beg < cmd && cmd < vanilla_sym_end
}

func (cmd LatexCmd) IsTextCmd() bool {
	return cmd_text_beg < cmd && cmd < cmd_text_end
}

func (cmd LatexCmd) TakesOneArg() bool {
	return cmd_1arg_beg < cmd && cmd < cmd_1arg_end
}

func (cmd LatexCmd) TakesTwoArg() bool {
	return cmd_2arg_beg < cmd && cmd < cmd_2arg_end
}

func (cmd LatexCmd) IsEnclosing() bool {
	return cmd_enclosing_beg < cmd && cmd < cmd_enclosing_end
}
