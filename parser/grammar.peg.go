package parser

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

type namespace struct {
	scope     string
	namespace string
}

type typeDef struct {
	name string
	typ  *Type
}

type exception *Struct

type include string

func toIfaceSlice(v interface{}) []interface{} {
	if v == nil {
		return nil
	}
	return v.([]interface{})
}

func ifaceSliceToString(v interface{}) string {
	ifs := toIfaceSlice(v)
	b := make([]byte, len(ifs))
	for i, v := range ifs {
		b[i] = v.([]uint8)[0]
	}
	return string(b)
}

var g = &grammar{
	rules: []*rule{
		{
			name: "Grammer",
			pos:  position{line: 41, col: 1, offset: 514},
			expr: &actionExpr{
				pos: position{line: 41, col: 11, offset: 526},
				run: (*parser).callonGrammer1,
				expr: &seqExpr{
					pos: position{line: 41, col: 11, offset: 526},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 41, col: 11, offset: 526},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 41, col: 14, offset: 529},
							label: "statements",
							expr: &zeroOrMoreExpr{
								pos: position{line: 41, col: 25, offset: 540},
								expr: &seqExpr{
									pos: position{line: 41, col: 27, offset: 542},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 41, col: 27, offset: 542},
											name: "Statement",
										},
										&ruleRefExpr{
											pos:  position{line: 41, col: 37, offset: 552},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 41, col: 44, offset: 559},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 41, col: 44, offset: 559},
									name: "EOF",
								},
								&ruleRefExpr{
									pos:  position{line: 41, col: 50, offset: 565},
									name: "SyntaxError",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "SyntaxError",
			pos:  position{line: 82, col: 1, offset: 1626},
			expr: &actionExpr{
				pos: position{line: 82, col: 15, offset: 1642},
				run: (*parser).callonSyntaxError1,
				expr: &anyMatcher{
					line: 82, col: 15, offset: 1642,
				},
			},
		},
		{
			name: "Include",
			pos:  position{line: 86, col: 1, offset: 1697},
			expr: &actionExpr{
				pos: position{line: 86, col: 11, offset: 1709},
				run: (*parser).callonInclude1,
				expr: &seqExpr{
					pos: position{line: 86, col: 11, offset: 1709},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 86, col: 11, offset: 1709},
							val:        "include",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 21, offset: 1719},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 86, col: 23, offset: 1721},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 86, col: 28, offset: 1726},
								name: "Literal",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 36, offset: 1734},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 90, col: 1, offset: 1779},
			expr: &choiceExpr{
				pos: position{line: 90, col: 13, offset: 1793},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 90, col: 13, offset: 1793},
						name: "Include",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 23, offset: 1803},
						name: "Namespace",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 35, offset: 1815},
						name: "Const",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 43, offset: 1823},
						name: "Enum",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 50, offset: 1830},
						name: "TypeDef",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 60, offset: 1840},
						name: "Struct",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 69, offset: 1849},
						name: "Exception",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 81, offset: 1861},
						name: "Service",
					},
				},
			},
		},
		{
			name: "Namespace",
			pos:  position{line: 92, col: 1, offset: 1870},
			expr: &actionExpr{
				pos: position{line: 92, col: 13, offset: 1884},
				run: (*parser).callonNamespace1,
				expr: &seqExpr{
					pos: position{line: 92, col: 13, offset: 1884},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 92, col: 13, offset: 1884},
							val:        "namespace",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 92, col: 25, offset: 1896},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 92, col: 27, offset: 1898},
							label: "scope",
							expr: &oneOrMoreExpr{
								pos: position{line: 92, col: 33, offset: 1904},
								expr: &charClassMatcher{
									pos:        position{line: 92, col: 33, offset: 1904},
									val:        "[a-z.-]",
									chars:      []rune{'.', '-'},
									ranges:     []rune{'a', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 92, col: 42, offset: 1913},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 92, col: 44, offset: 1915},
							label: "ns",
							expr: &ruleRefExpr{
								pos:  position{line: 92, col: 47, offset: 1918},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 92, col: 58, offset: 1929},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Const",
			pos:  position{line: 99, col: 1, offset: 2040},
			expr: &actionExpr{
				pos: position{line: 99, col: 9, offset: 2050},
				run: (*parser).callonConst1,
				expr: &seqExpr{
					pos: position{line: 99, col: 9, offset: 2050},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 99, col: 9, offset: 2050},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 17, offset: 2058},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 99, col: 19, offset: 2060},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 99, col: 23, offset: 2064},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 33, offset: 2074},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 99, col: 35, offset: 2076},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 99, col: 40, offset: 2081},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 51, offset: 2092},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 99, col: 53, offset: 2094},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 57, offset: 2098},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 99, col: 59, offset: 2100},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 99, col: 65, offset: 2106},
								name: "ConstValue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 76, offset: 2117},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Enum",
			pos:  position{line: 107, col: 1, offset: 2225},
			expr: &actionExpr{
				pos: position{line: 107, col: 8, offset: 2234},
				run: (*parser).callonEnum1,
				expr: &seqExpr{
					pos: position{line: 107, col: 8, offset: 2234},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 107, col: 8, offset: 2234},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 15, offset: 2241},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 107, col: 17, offset: 2243},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 107, col: 22, offset: 2248},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 33, offset: 2259},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 107, col: 35, offset: 2261},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 39, offset: 2265},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 107, col: 42, offset: 2268},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 107, col: 49, offset: 2275},
								expr: &seqExpr{
									pos: position{line: 107, col: 50, offset: 2276},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 107, col: 50, offset: 2276},
											name: "EnumValue",
										},
										&ruleRefExpr{
											pos:  position{line: 107, col: 60, offset: 2286},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 107, col: 65, offset: 2291},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 69, offset: 2295},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 130, col: 1, offset: 2811},
			expr: &actionExpr{
				pos: position{line: 130, col: 13, offset: 2825},
				run: (*parser).callonEnumValue1,
				expr: &seqExpr{
					pos: position{line: 130, col: 13, offset: 2825},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 130, col: 13, offset: 2825},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 130, col: 18, offset: 2830},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 29, offset: 2841},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 130, col: 31, offset: 2843},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 130, col: 37, offset: 2849},
								expr: &seqExpr{
									pos: position{line: 130, col: 38, offset: 2850},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 130, col: 38, offset: 2850},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 130, col: 42, offset: 2854},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 130, col: 44, offset: 2856},
											name: "IntConstant",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 130, col: 58, offset: 2870},
							expr: &ruleRefExpr{
								pos:  position{line: 130, col: 58, offset: 2870},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "TypeDef",
			pos:  position{line: 141, col: 1, offset: 3049},
			expr: &actionExpr{
				pos: position{line: 141, col: 11, offset: 3061},
				run: (*parser).callonTypeDef1,
				expr: &seqExpr{
					pos: position{line: 141, col: 11, offset: 3061},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 141, col: 11, offset: 3061},
							val:        "typedef",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 141, col: 21, offset: 3071},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 141, col: 23, offset: 3073},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 141, col: 27, offset: 3077},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 141, col: 37, offset: 3087},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 141, col: 39, offset: 3089},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 141, col: 44, offset: 3094},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 141, col: 55, offset: 3105},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Struct",
			pos:  position{line: 148, col: 1, offset: 3195},
			expr: &actionExpr{
				pos: position{line: 148, col: 10, offset: 3206},
				run: (*parser).callonStruct1,
				expr: &seqExpr{
					pos: position{line: 148, col: 10, offset: 3206},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 148, col: 10, offset: 3206},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 148, col: 19, offset: 3215},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 148, col: 21, offset: 3217},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 148, col: 24, offset: 3220},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Exception",
			pos:  position{line: 149, col: 1, offset: 3260},
			expr: &actionExpr{
				pos: position{line: 149, col: 13, offset: 3274},
				run: (*parser).callonException1,
				expr: &seqExpr{
					pos: position{line: 149, col: 13, offset: 3274},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 149, col: 13, offset: 3274},
							val:        "exception",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 149, col: 25, offset: 3286},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 149, col: 27, offset: 3288},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 149, col: 30, offset: 3291},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "StructLike",
			pos:  position{line: 150, col: 1, offset: 3342},
			expr: &actionExpr{
				pos: position{line: 150, col: 14, offset: 3357},
				run: (*parser).callonStructLike1,
				expr: &seqExpr{
					pos: position{line: 150, col: 14, offset: 3357},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 150, col: 14, offset: 3357},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 150, col: 19, offset: 3362},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 150, col: 30, offset: 3373},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 150, col: 32, offset: 3375},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 150, col: 36, offset: 3379},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 150, col: 39, offset: 3382},
							label: "fields",
							expr: &ruleRefExpr{
								pos:  position{line: 150, col: 46, offset: 3389},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 150, col: 56, offset: 3399},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 150, col: 60, offset: 3403},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "FieldList",
			pos:  position{line: 160, col: 1, offset: 3537},
			expr: &actionExpr{
				pos: position{line: 160, col: 13, offset: 3551},
				run: (*parser).callonFieldList1,
				expr: &labeledExpr{
					pos:   position{line: 160, col: 13, offset: 3551},
					label: "fields",
					expr: &zeroOrMoreExpr{
						pos: position{line: 160, col: 20, offset: 3558},
						expr: &seqExpr{
							pos: position{line: 160, col: 21, offset: 3559},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 160, col: 21, offset: 3559},
									name: "Field",
								},
								&ruleRefExpr{
									pos:  position{line: 160, col: 27, offset: 3565},
									name: "__",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Field",
			pos:  position{line: 169, col: 1, offset: 3725},
			expr: &actionExpr{
				pos: position{line: 169, col: 9, offset: 3735},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 169, col: 9, offset: 3735},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 169, col: 9, offset: 3735},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 12, offset: 3738},
								name: "IntConstant",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 24, offset: 3750},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 169, col: 26, offset: 3752},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 30, offset: 3756},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 169, col: 32, offset: 3758},
							label: "req",
							expr: &zeroOrOneExpr{
								pos: position{line: 169, col: 36, offset: 3762},
								expr: &ruleRefExpr{
									pos:  position{line: 169, col: 36, offset: 3762},
									name: "FieldReq",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 46, offset: 3772},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 169, col: 48, offset: 3774},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 52, offset: 3778},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 62, offset: 3788},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 169, col: 64, offset: 3790},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 69, offset: 3795},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 80, offset: 3806},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 169, col: 82, offset: 3808},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 169, col: 86, offset: 3812},
								expr: &seqExpr{
									pos: position{line: 169, col: 87, offset: 3813},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 169, col: 87, offset: 3813},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 169, col: 91, offset: 3817},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 169, col: 93, offset: 3819},
											name: "ConstValue",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 169, col: 106, offset: 3832},
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 106, offset: 3832},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldReq",
			pos:  position{line: 184, col: 1, offset: 4092},
			expr: &actionExpr{
				pos: position{line: 184, col: 12, offset: 4105},
				run: (*parser).callonFieldReq1,
				expr: &choiceExpr{
					pos: position{line: 184, col: 13, offset: 4106},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 184, col: 13, offset: 4106},
							val:        "required",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 184, col: 26, offset: 4119},
							val:        "optional",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Service",
			pos:  position{line: 188, col: 1, offset: 4190},
			expr: &actionExpr{
				pos: position{line: 188, col: 11, offset: 4202},
				run: (*parser).callonService1,
				expr: &seqExpr{
					pos: position{line: 188, col: 11, offset: 4202},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 188, col: 11, offset: 4202},
							val:        "service",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 21, offset: 4212},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 23, offset: 4214},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 188, col: 28, offset: 4219},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 39, offset: 4230},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 41, offset: 4232},
							label: "extends",
							expr: &zeroOrOneExpr{
								pos: position{line: 188, col: 49, offset: 4240},
								expr: &seqExpr{
									pos: position{line: 188, col: 50, offset: 4241},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 188, col: 50, offset: 4241},
											val:        "extends",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 60, offset: 4251},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 63, offset: 4254},
											name: "Identifier",
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 74, offset: 4265},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 79, offset: 4270},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 188, col: 82, offset: 4273},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 86, offset: 4277},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 89, offset: 4280},
							label: "methods",
							expr: &zeroOrMoreExpr{
								pos: position{line: 188, col: 97, offset: 4288},
								expr: &seqExpr{
									pos: position{line: 188, col: 98, offset: 4289},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 188, col: 98, offset: 4289},
											name: "Function",
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 107, offset: 4298},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 188, col: 113, offset: 4304},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 188, col: 113, offset: 4304},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 188, col: 119, offset: 4310},
									name: "EndOfServiceError",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 138, offset: 4329},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EndOfServiceError",
			pos:  position{line: 203, col: 1, offset: 4670},
			expr: &actionExpr{
				pos: position{line: 203, col: 21, offset: 4692},
				run: (*parser).callonEndOfServiceError1,
				expr: &anyMatcher{
					line: 203, col: 21, offset: 4692,
				},
			},
		},
		{
			name: "Function",
			pos:  position{line: 207, col: 1, offset: 4758},
			expr: &actionExpr{
				pos: position{line: 207, col: 12, offset: 4771},
				run: (*parser).callonFunction1,
				expr: &seqExpr{
					pos: position{line: 207, col: 12, offset: 4771},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 207, col: 12, offset: 4771},
							label: "oneway",
							expr: &zeroOrOneExpr{
								pos: position{line: 207, col: 19, offset: 4778},
								expr: &seqExpr{
									pos: position{line: 207, col: 20, offset: 4779},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 207, col: 20, offset: 4779},
											val:        "oneway",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 207, col: 29, offset: 4788},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 207, col: 34, offset: 4793},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 38, offset: 4797},
								name: "FunctionType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 51, offset: 4810},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 207, col: 54, offset: 4813},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 59, offset: 4818},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 70, offset: 4829},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 207, col: 72, offset: 4831},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 76, offset: 4835},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 207, col: 79, offset: 4838},
							label: "arguments",
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 89, offset: 4848},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 207, col: 99, offset: 4858},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 103, offset: 4862},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 207, col: 106, offset: 4865},
							label: "exceptions",
							expr: &zeroOrOneExpr{
								pos: position{line: 207, col: 117, offset: 4876},
								expr: &ruleRefExpr{
									pos:  position{line: 207, col: 117, offset: 4876},
									name: "Throws",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 207, col: 125, offset: 4884},
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 125, offset: 4884},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionType",
			pos:  position{line: 230, col: 1, offset: 5265},
			expr: &actionExpr{
				pos: position{line: 230, col: 16, offset: 5282},
				run: (*parser).callonFunctionType1,
				expr: &labeledExpr{
					pos:   position{line: 230, col: 16, offset: 5282},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 230, col: 21, offset: 5287},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 230, col: 21, offset: 5287},
								val:        "void",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 230, col: 30, offset: 5296},
								name: "FieldType",
							},
						},
					},
				},
			},
		},
		{
			name: "Throws",
			pos:  position{line: 237, col: 1, offset: 5403},
			expr: &actionExpr{
				pos: position{line: 237, col: 10, offset: 5414},
				run: (*parser).callonThrows1,
				expr: &seqExpr{
					pos: position{line: 237, col: 10, offset: 5414},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 237, col: 10, offset: 5414},
							val:        "throws",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 237, col: 19, offset: 5423},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 237, col: 22, offset: 5426},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 237, col: 26, offset: 5430},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 237, col: 29, offset: 5433},
							label: "exceptions",
							expr: &ruleRefExpr{
								pos:  position{line: 237, col: 40, offset: 5444},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 237, col: 50, offset: 5454},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "FieldType",
			pos:  position{line: 241, col: 1, offset: 5487},
			expr: &actionExpr{
				pos: position{line: 241, col: 13, offset: 5501},
				run: (*parser).callonFieldType1,
				expr: &labeledExpr{
					pos:   position{line: 241, col: 13, offset: 5501},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 241, col: 18, offset: 5506},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 241, col: 18, offset: 5506},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 241, col: 29, offset: 5517},
								name: "ContainerType",
							},
							&ruleRefExpr{
								pos:  position{line: 241, col: 45, offset: 5533},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "DefinitionType",
			pos:  position{line: 248, col: 1, offset: 5643},
			expr: &actionExpr{
				pos: position{line: 248, col: 18, offset: 5662},
				run: (*parser).callonDefinitionType1,
				expr: &labeledExpr{
					pos:   position{line: 248, col: 18, offset: 5662},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 248, col: 23, offset: 5667},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 248, col: 23, offset: 5667},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 248, col: 34, offset: 5678},
								name: "ContainerType",
							},
						},
					},
				},
			},
		},
		{
			name: "BaseType",
			pos:  position{line: 252, col: 1, offset: 5715},
			expr: &actionExpr{
				pos: position{line: 252, col: 12, offset: 5728},
				run: (*parser).callonBaseType1,
				expr: &choiceExpr{
					pos: position{line: 252, col: 13, offset: 5729},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 252, col: 13, offset: 5729},
							val:        "bool",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 22, offset: 5738},
							val:        "byte",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 31, offset: 5747},
							val:        "i16",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 39, offset: 5755},
							val:        "i32",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 47, offset: 5763},
							val:        "i64",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 55, offset: 5771},
							val:        "double",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 66, offset: 5782},
							val:        "string",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 77, offset: 5793},
							val:        "binary",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContainerType",
			pos:  position{line: 256, col: 1, offset: 5850},
			expr: &actionExpr{
				pos: position{line: 256, col: 17, offset: 5868},
				run: (*parser).callonContainerType1,
				expr: &labeledExpr{
					pos:   position{line: 256, col: 17, offset: 5868},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 256, col: 22, offset: 5873},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 256, col: 22, offset: 5873},
								name: "MapType",
							},
							&ruleRefExpr{
								pos:  position{line: 256, col: 32, offset: 5883},
								name: "SetType",
							},
							&ruleRefExpr{
								pos:  position{line: 256, col: 42, offset: 5893},
								name: "ListType",
							},
						},
					},
				},
			},
		},
		{
			name: "MapType",
			pos:  position{line: 260, col: 1, offset: 5925},
			expr: &actionExpr{
				pos: position{line: 260, col: 11, offset: 5937},
				run: (*parser).callonMapType1,
				expr: &seqExpr{
					pos: position{line: 260, col: 11, offset: 5937},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 260, col: 11, offset: 5937},
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 11, offset: 5937},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 260, col: 20, offset: 5946},
							val:        "map<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 27, offset: 5953},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 260, col: 30, offset: 5956},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 34, offset: 5960},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 44, offset: 5970},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 260, col: 47, offset: 5973},
							val:        ",",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 51, offset: 5977},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 260, col: 54, offset: 5980},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 60, offset: 5986},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 70, offset: 5996},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 260, col: 73, offset: 5999},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SetType",
			pos:  position{line: 268, col: 1, offset: 6098},
			expr: &actionExpr{
				pos: position{line: 268, col: 11, offset: 6110},
				run: (*parser).callonSetType1,
				expr: &seqExpr{
					pos: position{line: 268, col: 11, offset: 6110},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 268, col: 11, offset: 6110},
							expr: &ruleRefExpr{
								pos:  position{line: 268, col: 11, offset: 6110},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 268, col: 20, offset: 6119},
							val:        "set<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 268, col: 27, offset: 6126},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 268, col: 30, offset: 6129},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 268, col: 34, offset: 6133},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 268, col: 44, offset: 6143},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 268, col: 47, offset: 6146},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 275, col: 1, offset: 6219},
			expr: &actionExpr{
				pos: position{line: 275, col: 12, offset: 6232},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 275, col: 12, offset: 6232},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 275, col: 12, offset: 6232},
							val:        "list<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 20, offset: 6240},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 275, col: 23, offset: 6243},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 275, col: 27, offset: 6247},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 37, offset: 6257},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 275, col: 40, offset: 6260},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CppType",
			pos:  position{line: 282, col: 1, offset: 6334},
			expr: &actionExpr{
				pos: position{line: 282, col: 11, offset: 6346},
				run: (*parser).callonCppType1,
				expr: &seqExpr{
					pos: position{line: 282, col: 11, offset: 6346},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 282, col: 11, offset: 6346},
							val:        "cpp_type",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 282, col: 22, offset: 6357},
							label: "cppType",
							expr: &ruleRefExpr{
								pos:  position{line: 282, col: 30, offset: 6365},
								name: "Literal",
							},
						},
					},
				},
			},
		},
		{
			name: "ConstValue",
			pos:  position{line: 286, col: 1, offset: 6399},
			expr: &choiceExpr{
				pos: position{line: 286, col: 14, offset: 6414},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 286, col: 14, offset: 6414},
						name: "Literal",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 24, offset: 6424},
						name: "DoubleConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 41, offset: 6441},
						name: "IntConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 55, offset: 6455},
						name: "ConstMap",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 66, offset: 6466},
						name: "ConstList",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 78, offset: 6478},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "IntConstant",
			pos:  position{line: 288, col: 1, offset: 6490},
			expr: &actionExpr{
				pos: position{line: 288, col: 15, offset: 6506},
				run: (*parser).callonIntConstant1,
				expr: &seqExpr{
					pos: position{line: 288, col: 15, offset: 6506},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 288, col: 15, offset: 6506},
							expr: &charClassMatcher{
								pos:        position{line: 288, col: 15, offset: 6506},
								val:        "[-+]",
								chars:      []rune{'-', '+'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 288, col: 21, offset: 6512},
							expr: &ruleRefExpr{
								pos:  position{line: 288, col: 21, offset: 6512},
								name: "Digit",
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleConstant",
			pos:  position{line: 292, col: 1, offset: 6573},
			expr: &actionExpr{
				pos: position{line: 292, col: 18, offset: 6592},
				run: (*parser).callonDoubleConstant1,
				expr: &seqExpr{
					pos: position{line: 292, col: 18, offset: 6592},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 292, col: 18, offset: 6592},
							expr: &charClassMatcher{
								pos:        position{line: 292, col: 18, offset: 6592},
								val:        "[+-]",
								chars:      []rune{'+', '-'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 292, col: 24, offset: 6598},
							expr: &ruleRefExpr{
								pos:  position{line: 292, col: 24, offset: 6598},
								name: "Digit",
							},
						},
						&litMatcher{
							pos:        position{line: 292, col: 31, offset: 6605},
							val:        ".",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 292, col: 35, offset: 6609},
							expr: &ruleRefExpr{
								pos:  position{line: 292, col: 35, offset: 6609},
								name: "Digit",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 292, col: 42, offset: 6616},
							expr: &seqExpr{
								pos: position{line: 292, col: 44, offset: 6618},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 292, col: 44, offset: 6618},
										val:        "['Ee']",
										chars:      []rune{'\'', 'E', 'e', '\''},
										ignoreCase: false,
										inverted:   false,
									},
									&ruleRefExpr{
										pos:  position{line: 292, col: 51, offset: 6625},
										name: "IntConstant",
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ConstList",
			pos:  position{line: 296, col: 1, offset: 6692},
			expr: &actionExpr{
				pos: position{line: 296, col: 13, offset: 6706},
				run: (*parser).callonConstList1,
				expr: &seqExpr{
					pos: position{line: 296, col: 13, offset: 6706},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 296, col: 13, offset: 6706},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 17, offset: 6710},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 296, col: 20, offset: 6713},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 296, col: 27, offset: 6720},
								expr: &seqExpr{
									pos: position{line: 296, col: 28, offset: 6721},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 296, col: 28, offset: 6721},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 296, col: 39, offset: 6732},
											name: "__",
										},
										&zeroOrOneExpr{
											pos: position{line: 296, col: 42, offset: 6735},
											expr: &ruleRefExpr{
												pos:  position{line: 296, col: 42, offset: 6735},
												name: "ListSeparator",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 296, col: 57, offset: 6750},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 62, offset: 6755},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 296, col: 65, offset: 6758},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ConstMap",
			pos:  position{line: 305, col: 1, offset: 6931},
			expr: &actionExpr{
				pos: position{line: 305, col: 12, offset: 6944},
				run: (*parser).callonConstMap1,
				expr: &seqExpr{
					pos: position{line: 305, col: 12, offset: 6944},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 305, col: 12, offset: 6944},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 305, col: 16, offset: 6948},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 305, col: 19, offset: 6951},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 305, col: 26, offset: 6958},
								expr: &seqExpr{
									pos: position{line: 305, col: 27, offset: 6959},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 305, col: 27, offset: 6959},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 38, offset: 6970},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 305, col: 41, offset: 6973},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 45, offset: 6977},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 48, offset: 6980},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 59, offset: 6991},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 305, col: 63, offset: 6995},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 305, col: 63, offset: 6995},
													val:        ",",
													ignoreCase: false,
												},
												&andExpr{
													pos: position{line: 305, col: 69, offset: 7001},
													expr: &litMatcher{
														pos:        position{line: 305, col: 70, offset: 7002},
														val:        "}",
														ignoreCase: false,
													},
												},
											},
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 75, offset: 7007},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 305, col: 80, offset: 7012},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Literal",
			pos:  position{line: 321, col: 1, offset: 7258},
			expr: &actionExpr{
				pos: position{line: 321, col: 11, offset: 7270},
				run: (*parser).callonLiteral1,
				expr: &choiceExpr{
					pos: position{line: 321, col: 12, offset: 7271},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 321, col: 13, offset: 7272},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 321, col: 13, offset: 7272},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 321, col: 17, offset: 7276},
									expr: &choiceExpr{
										pos: position{line: 321, col: 18, offset: 7277},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 321, col: 18, offset: 7277},
												val:        "\\\"",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 321, col: 25, offset: 7284},
												val:        "[^\"]",
												chars:      []rune{'"'},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 321, col: 32, offset: 7291},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 321, col: 40, offset: 7299},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 321, col: 40, offset: 7299},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 321, col: 45, offset: 7304},
									expr: &choiceExpr{
										pos: position{line: 321, col: 46, offset: 7305},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 321, col: 46, offset: 7305},
												val:        "\\'",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 321, col: 53, offset: 7312},
												val:        "[^']",
												chars:      []rune{'\''},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 321, col: 60, offset: 7319},
									val:        "'",
									ignoreCase: false,
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Identifier",
			pos:  position{line: 328, col: 1, offset: 7520},
			expr: &actionExpr{
				pos: position{line: 328, col: 14, offset: 7535},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 328, col: 14, offset: 7535},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 328, col: 14, offset: 7535},
							expr: &choiceExpr{
								pos: position{line: 328, col: 15, offset: 7536},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 328, col: 15, offset: 7536},
										name: "Letter",
									},
									&litMatcher{
										pos:        position{line: 328, col: 24, offset: 7545},
										val:        "_",
										ignoreCase: false,
									},
								},
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 328, col: 30, offset: 7551},
							expr: &choiceExpr{
								pos: position{line: 328, col: 31, offset: 7552},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 328, col: 31, offset: 7552},
										name: "Letter",
									},
									&ruleRefExpr{
										pos:  position{line: 328, col: 40, offset: 7561},
										name: "Digit",
									},
									&charClassMatcher{
										pos:        position{line: 328, col: 48, offset: 7569},
										val:        "[._]",
										chars:      []rune{'.', '_'},
										ignoreCase: false,
										inverted:   false,
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ListSeparator",
			pos:  position{line: 332, col: 1, offset: 7621},
			expr: &charClassMatcher{
				pos:        position{line: 332, col: 17, offset: 7639},
				val:        "[,;]",
				chars:      []rune{',', ';'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Letter",
			pos:  position{line: 333, col: 1, offset: 7644},
			expr: &charClassMatcher{
				pos:        position{line: 333, col: 10, offset: 7655},
				val:        "[A-Za-z]",
				ranges:     []rune{'A', 'Z', 'a', 'z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Digit",
			pos:  position{line: 334, col: 1, offset: 7664},
			expr: &charClassMatcher{
				pos:        position{line: 334, col: 9, offset: 7674},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 338, col: 1, offset: 7685},
			expr: &anyMatcher{
				line: 338, col: 14, offset: 7700,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 339, col: 1, offset: 7702},
			expr: &choiceExpr{
				pos: position{line: 339, col: 11, offset: 7714},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 339, col: 11, offset: 7714},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 339, col: 30, offset: 7733},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 340, col: 1, offset: 7751},
			expr: &seqExpr{
				pos: position{line: 340, col: 20, offset: 7772},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 340, col: 20, offset: 7772},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 340, col: 25, offset: 7777},
						expr: &seqExpr{
							pos: position{line: 340, col: 27, offset: 7779},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 340, col: 27, offset: 7779},
									expr: &litMatcher{
										pos:        position{line: 340, col: 28, offset: 7780},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 340, col: 33, offset: 7785},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 340, col: 47, offset: 7799},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 341, col: 1, offset: 7804},
			expr: &seqExpr{
				pos: position{line: 341, col: 36, offset: 7841},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 341, col: 36, offset: 7841},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 341, col: 41, offset: 7846},
						expr: &seqExpr{
							pos: position{line: 341, col: 43, offset: 7848},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 341, col: 43, offset: 7848},
									expr: &choiceExpr{
										pos: position{line: 341, col: 46, offset: 7851},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 341, col: 46, offset: 7851},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 341, col: 53, offset: 7858},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 341, col: 59, offset: 7864},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 341, col: 73, offset: 7878},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 342, col: 1, offset: 7883},
			expr: &choiceExpr{
				pos: position{line: 342, col: 21, offset: 7905},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 342, col: 22, offset: 7906},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 342, col: 22, offset: 7906},
								val:        "//",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 342, col: 27, offset: 7911},
								expr: &seqExpr{
									pos: position{line: 342, col: 29, offset: 7913},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 342, col: 29, offset: 7913},
											expr: &ruleRefExpr{
												pos:  position{line: 342, col: 30, offset: 7914},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 342, col: 34, offset: 7918},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 342, col: 52, offset: 7936},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 342, col: 52, offset: 7936},
								val:        "#",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 342, col: 56, offset: 7940},
								expr: &seqExpr{
									pos: position{line: 342, col: 58, offset: 7942},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 342, col: 58, offset: 7942},
											expr: &ruleRefExpr{
												pos:  position{line: 342, col: 59, offset: 7943},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 342, col: 63, offset: 7947},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "__",
			pos:  position{line: 344, col: 1, offset: 7963},
			expr: &zeroOrMoreExpr{
				pos: position{line: 344, col: 6, offset: 7970},
				expr: &choiceExpr{
					pos: position{line: 344, col: 8, offset: 7972},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 344, col: 8, offset: 7972},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 344, col: 21, offset: 7985},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 344, col: 27, offset: 7991},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 345, col: 1, offset: 8002},
			expr: &zeroOrMoreExpr{
				pos: position{line: 345, col: 5, offset: 8008},
				expr: &choiceExpr{
					pos: position{line: 345, col: 7, offset: 8010},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 345, col: 7, offset: 8010},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 345, col: 20, offset: 8023},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "WS",
			pos:  position{line: 346, col: 1, offset: 8059},
			expr: &zeroOrMoreExpr{
				pos: position{line: 346, col: 6, offset: 8066},
				expr: &ruleRefExpr{
					pos:  position{line: 346, col: 6, offset: 8066},
					name: "Whitespace",
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 348, col: 1, offset: 8079},
			expr: &charClassMatcher{
				pos:        position{line: 348, col: 14, offset: 8094},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 349, col: 1, offset: 8102},
			expr: &litMatcher{
				pos:        position{line: 349, col: 7, offset: 8110},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 350, col: 1, offset: 8115},
			expr: &choiceExpr{
				pos: position{line: 350, col: 7, offset: 8123},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 350, col: 7, offset: 8123},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 350, col: 7, offset: 8123},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 350, col: 10, offset: 8126},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 350, col: 16, offset: 8132},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 350, col: 16, offset: 8132},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 350, col: 18, offset: 8134},
								expr: &ruleRefExpr{
									pos:  position{line: 350, col: 18, offset: 8134},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 350, col: 37, offset: 8153},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 350, col: 43, offset: 8159},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 350, col: 43, offset: 8159},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 350, col: 46, offset: 8162},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 352, col: 1, offset: 8167},
			expr: &notExpr{
				pos: position{line: 352, col: 7, offset: 8175},
				expr: &anyMatcher{
					line: 352, col: 8, offset: 8176,
				},
			},
		},
	},
}

func (c *current) onGrammer1(statements interface{}) (interface{}, error) {
	thrift := &Thrift{
		Includes:   make(map[string]string),
		Namespaces: make(map[string]string),
		Typedefs:   make(map[string]*Type),
		Constants:  make(map[string]*Constant),
		Enums:      make(map[string]*Enum),
		Structs:    make(map[string]*Struct),
		Exceptions: make(map[string]*Struct),
		Services:   make(map[string]*Service),
	}
	stmts := toIfaceSlice(statements)
	for _, st := range stmts {
		switch v := st.([]interface{})[0].(type) {
		case *namespace:
			thrift.Namespaces[v.scope] = v.namespace
		case *Constant:
			thrift.Constants[v.Name] = v
		case *Enum:
			thrift.Enums[v.Name] = v
		case *typeDef:
			thrift.Typedefs[v.name] = v.typ
		case *Struct:
			thrift.Structs[v.Name] = v
		case exception:
			thrift.Exceptions[v.Name] = (*Struct)(v)
		case *Service:
			thrift.Services[v.Name] = v
		case include:
			name := string(v)
			if ix := strings.LastIndex(name, "."); ix > 0 {
				name = name[:ix]
			}
			thrift.Includes[name] = string(v)
		default:
			return nil, fmt.Errorf("parser: unknown value %#v", v)
		}
	}
	return thrift, nil
}

func (p *parser) callonGrammer1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onGrammer1(stack["statements"])
}

func (c *current) onSyntaxError1() (interface{}, error) {
	return nil, errors.New("parser: syntax error")
}

func (p *parser) callonSyntaxError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSyntaxError1()
}

func (c *current) onInclude1(file interface{}) (interface{}, error) {
	return include(file.(string)), nil
}

func (p *parser) callonInclude1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onInclude1(stack["file"])
}

func (c *current) onNamespace1(scope, ns interface{}) (interface{}, error) {
	return &namespace{
		scope:     ifaceSliceToString(scope),
		namespace: string(ns.(Identifier)),
	}, nil
}

func (p *parser) callonNamespace1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamespace1(stack["scope"], stack["ns"])
}

func (c *current) onConst1(typ, name, value interface{}) (interface{}, error) {
	return &Constant{
		Name:  string(name.(Identifier)),
		Type:  typ.(*Type),
		Value: value,
	}, nil
}

func (p *parser) callonConst1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConst1(stack["typ"], stack["name"], stack["value"])
}

func (c *current) onEnum1(name, values interface{}) (interface{}, error) {
	vs := toIfaceSlice(values)
	en := &Enum{
		Name:   string(name.(Identifier)),
		Values: make(map[string]*EnumValue, len(vs)),
	}
	// Assigns numbers in order. This will behave badly if some values are
	// defined and other are not, but I think that's ok since that's a silly
	// thing to do.
	next := 0
	for _, v := range vs {
		ev := v.([]interface{})[0].(*EnumValue)
		if ev.Value < 0 {
			ev.Value = next
		}
		if ev.Value >= next {
			next = ev.Value + 1
		}
		en.Values[ev.Name] = ev
	}
	return en, nil
}

func (p *parser) callonEnum1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnum1(stack["name"], stack["values"])
}

func (c *current) onEnumValue1(name, value interface{}) (interface{}, error) {
	ev := &EnumValue{
		Name:  string(name.(Identifier)),
		Value: -1,
	}
	if value != nil {
		ev.Value = int(value.([]interface{})[2].(int64))
	}
	return ev, nil
}

func (p *parser) callonEnumValue1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEnumValue1(stack["name"], stack["value"])
}

func (c *current) onTypeDef1(typ, name interface{}) (interface{}, error) {
	return &typeDef{
		name: string(name.(Identifier)),
		typ:  typ.(*Type),
	}, nil
}

func (p *parser) callonTypeDef1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onTypeDef1(stack["typ"], stack["name"])
}

func (c *current) onStruct1(st interface{}) (interface{}, error) {
	return st.(*Struct), nil
}

func (p *parser) callonStruct1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStruct1(stack["st"])
}

func (c *current) onException1(st interface{}) (interface{}, error) {
	return exception(st.(*Struct)), nil
}

func (p *parser) callonException1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onException1(stack["st"])
}

func (c *current) onStructLike1(name, fields interface{}) (interface{}, error) {
	st := &Struct{
		Name: string(name.(Identifier)),
	}
	if fields != nil {
		st.Fields = fields.([]*Field)
	}
	return st, nil
}

func (p *parser) callonStructLike1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onStructLike1(stack["name"], stack["fields"])
}

func (c *current) onFieldList1(fields interface{}) (interface{}, error) {
	fs := fields.([]interface{})
	flds := make([]*Field, len(fs))
	for i, f := range fs {
		flds[i] = f.([]interface{})[0].(*Field)
	}
	return flds, nil
}

func (p *parser) callonFieldList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldList1(stack["fields"])
}

func (c *current) onField1(id, req, typ, name, def interface{}) (interface{}, error) {
	f := &Field{
		ID:   int(id.(int64)),
		Name: string(name.(Identifier)),
		Type: typ.(*Type),
	}
	if req != nil && !req.(bool) {
		f.Optional = true
	}
	if def != nil {
		f.Default = def.([]interface{})[2]
	}
	return f, nil
}

func (p *parser) callonField1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onField1(stack["id"], stack["req"], stack["typ"], stack["name"], stack["def"])
}

func (c *current) onFieldReq1() (interface{}, error) {
	return !bytes.Equal(c.text, []byte("optional")), nil
}

func (p *parser) callonFieldReq1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldReq1()
}

func (c *current) onService1(name, extends, methods interface{}) (interface{}, error) {
	ms := methods.([]interface{})
	svc := &Service{
		Name:    string(name.(Identifier)),
		Methods: make(map[string]*Method, len(ms)),
	}
	if extends != nil {
		svc.Extends = string(extends.([]interface{})[2].(Identifier))
	}
	for _, m := range ms {
		mt := m.([]interface{})[0].(*Method)
		svc.Methods[mt.Name] = mt
	}
	return svc, nil
}

func (p *parser) callonService1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onService1(stack["name"], stack["extends"], stack["methods"])
}

func (c *current) onEndOfServiceError1() (interface{}, error) {
	return nil, errors.New("parser: expected end of service")
}

func (p *parser) callonEndOfServiceError1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onEndOfServiceError1()
}

func (c *current) onFunction1(oneway, typ, name, arguments, exceptions interface{}) (interface{}, error) {
	m := &Method{
		Name: string(name.(Identifier)),
	}
	t := typ.(*Type)
	if t.Name != "void" {
		m.ReturnType = t
	}
	if oneway != nil {
		m.Oneway = true
	}
	if arguments != nil {
		m.Arguments = arguments.([]*Field)
	}
	if exceptions != nil {
		m.Exceptions = exceptions.([]*Field)
		for _, e := range m.Exceptions {
			e.Optional = true
		}
	}
	return m, nil
}

func (p *parser) callonFunction1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunction1(stack["oneway"], stack["typ"], stack["name"], stack["arguments"], stack["exceptions"])
}

func (c *current) onFunctionType1(typ interface{}) (interface{}, error) {
	if t, ok := typ.(*Type); ok {
		return t, nil
	}
	return &Type{Name: string(c.text)}, nil
}

func (p *parser) callonFunctionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFunctionType1(stack["typ"])
}

func (c *current) onThrows1(exceptions interface{}) (interface{}, error) {
	return exceptions, nil
}

func (p *parser) callonThrows1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onThrows1(stack["exceptions"])
}

func (c *current) onFieldType1(typ interface{}) (interface{}, error) {
	if t, ok := typ.(Identifier); ok {
		return &Type{Name: string(t)}, nil
	}
	return typ, nil
}

func (p *parser) callonFieldType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onFieldType1(stack["typ"])
}

func (c *current) onDefinitionType1(typ interface{}) (interface{}, error) {
	return typ, nil
}

func (p *parser) callonDefinitionType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDefinitionType1(stack["typ"])
}

func (c *current) onBaseType1() (interface{}, error) {
	return &Type{Name: string(c.text)}, nil
}

func (p *parser) callonBaseType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onBaseType1()
}

func (c *current) onContainerType1(typ interface{}) (interface{}, error) {
	return typ, nil
}

func (p *parser) callonContainerType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onContainerType1(stack["typ"])
}

func (c *current) onMapType1(key, value interface{}) (interface{}, error) {
	return &Type{
		Name:      "map",
		KeyType:   key.(*Type),
		ValueType: value.(*Type),
	}, nil
}

func (p *parser) callonMapType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onMapType1(stack["key"], stack["value"])
}

func (c *current) onSetType1(typ interface{}) (interface{}, error) {
	return &Type{
		Name:      "set",
		ValueType: typ.(*Type),
	}, nil
}

func (p *parser) callonSetType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onSetType1(stack["typ"])
}

func (c *current) onListType1(typ interface{}) (interface{}, error) {
	return &Type{
		Name:      "list",
		ValueType: typ.(*Type),
	}, nil
}

func (p *parser) callonListType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onListType1(stack["typ"])
}

func (c *current) onCppType1(cppType interface{}) (interface{}, error) {
	return cppType, nil
}

func (p *parser) callonCppType1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onCppType1(stack["cppType"])
}

func (c *current) onIntConstant1() (interface{}, error) {
	return strconv.ParseInt(string(c.text), 10, 64)
}

func (p *parser) callonIntConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIntConstant1()
}

func (c *current) onDoubleConstant1() (interface{}, error) {
	return strconv.ParseFloat(string(c.text), 64)
}

func (p *parser) callonDoubleConstant1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onDoubleConstant1()
}

func (c *current) onConstList1(values interface{}) (interface{}, error) {
	valueSlice := values.([]interface{})
	vs := make([]interface{}, len(valueSlice))
	for i, v := range valueSlice {
		vs[i] = v.([]interface{})[0]
	}
	return vs, nil
}

func (p *parser) callonConstList1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstList1(stack["values"])
}

func (c *current) onConstMap1(values interface{}) (interface{}, error) {
	if values == nil {
		return nil, nil
	}
	vals := values.([]interface{})
	kvs := make([]KeyValue, len(vals))
	for i, kv := range vals {
		v := kv.([]interface{})
		kvs[i] = KeyValue{
			Key:   v[0],
			Value: v[4],
		}
	}
	return kvs, nil
}

func (p *parser) callonConstMap1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onConstMap1(stack["values"])
}

func (c *current) onLiteral1() (interface{}, error) {
	if len(c.text) != 0 && c.text[0] == '\'' {
		return strconv.Unquote(`"` + strings.Replace(string(c.text[1:len(c.text)-1]), `\'`, `'`, -1) + `"`)
	}
	return strconv.Unquote(string(c.text))
}

func (p *parser) callonLiteral1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onLiteral1()
}

func (c *current) onIdentifier1() (interface{}, error) {
	return Identifier(string(c.text)), nil
}

func (p *parser) callonIdentifier1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onIdentifier1()
}

var (
	// errNoRule is returned when the grammar to parse has no rule.
	errNoRule = errors.New("grammar has no rule")

	// errInvalidEncoding is returned when the source is not properly
	// utf8-encoded.
	errInvalidEncoding = errors.New("invalid encoding")

	// errNoMatch is returned if no match could be found.
	errNoMatch = errors.New("no match found")
)

// Option is a function that can set an option on the parser. It returns
// the previous setting as an Option.
type Option func(*parser) Option

// Debug creates an Option to set the debug flag to b. When set to true,
// debugging information is printed to stdout while parsing.
//
// The default is false.
func Debug(b bool) Option {
	return func(p *parser) Option {
		old := p.debug
		p.debug = b
		return Debug(old)
	}
}

// Memoize creates an Option to set the memoize flag to b. When set to true,
// the parser will cache all results so each expression is evaluated only
// once. This guarantees linear parsing time even for pathological cases,
// at the expense of more memory and slower times for typical cases.
//
// The default is false.
func Memoize(b bool) Option {
	return func(p *parser) Option {
		old := p.memoize
		p.memoize = b
		return Memoize(old)
	}
}

// Recover creates an Option to set the recover flag to b. When set to
// true, this causes the parser to recover from panics and convert it
// to an error. Setting it to false can be useful while debugging to
// access the full stack trace.
//
// The default is true.
func Recover(b bool) Option {
	return func(p *parser) Option {
		old := p.recover
		p.recover = b
		return Recover(old)
	}
}

// ParseFile parses the file identified by filename.
func ParseFile(filename string, opts ...Option) (interface{}, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseReader(filename, f, opts...)
}

// ParseReader parses the data from r using filename as information in the
// error messages.
func ParseReader(filename string, r io.Reader, opts ...Option) (interface{}, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	return Parse(filename, b, opts...)
}

// Parse parses the data from b using filename as information in the
// error messages.
func Parse(filename string, b []byte, opts ...Option) (interface{}, error) {
	return newParser(filename, b, opts...).parse(g)
}

// position records a position in the text.
type position struct {
	line, col, offset int
}

func (p position) String() string {
	return fmt.Sprintf("%d:%d [%d]", p.line, p.col, p.offset)
}

// savepoint stores all state required to go back to this point in the
// parser.
type savepoint struct {
	position
	rn rune
	w  int
}

type current struct {
	pos  position // start position of the match
	text []byte   // raw text of the match
}

// the AST types...

type grammar struct {
	pos   position
	rules []*rule
}

type rule struct {
	pos         position
	name        string
	displayName string
	expr        interface{}
}

type choiceExpr struct {
	pos          position
	alternatives []interface{}
}

type actionExpr struct {
	pos  position
	expr interface{}
	run  func(*parser) (interface{}, error)
}

type seqExpr struct {
	pos   position
	exprs []interface{}
}

type labeledExpr struct {
	pos   position
	label string
	expr  interface{}
}

type expr struct {
	pos  position
	expr interface{}
}

type andExpr expr
type notExpr expr
type zeroOrOneExpr expr
type zeroOrMoreExpr expr
type oneOrMoreExpr expr

type ruleRefExpr struct {
	pos  position
	name string
}

type andCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type notCodeExpr struct {
	pos position
	run func(*parser) (bool, error)
}

type litMatcher struct {
	pos        position
	val        string
	ignoreCase bool
}

type charClassMatcher struct {
	pos        position
	val        string
	chars      []rune
	ranges     []rune
	classes    []*unicode.RangeTable
	ignoreCase bool
	inverted   bool
}

type anyMatcher position

// errList cumulates the errors found by the parser.
type errList []error

func (e *errList) add(err error) {
	*e = append(*e, err)
}

func (e errList) err() error {
	if len(e) == 0 {
		return nil
	}
	e.dedupe()
	return e
}

func (e *errList) dedupe() {
	var cleaned []error
	set := make(map[string]bool)
	for _, err := range *e {
		if msg := err.Error(); !set[msg] {
			set[msg] = true
			cleaned = append(cleaned, err)
		}
	}
	*e = cleaned
}

func (e errList) Error() string {
	switch len(e) {
	case 0:
		return ""
	case 1:
		return e[0].Error()
	default:
		var buf bytes.Buffer

		for i, err := range e {
			if i > 0 {
				buf.WriteRune('\n')
			}
			buf.WriteString(err.Error())
		}
		return buf.String()
	}
}

// parserError wraps an error with a prefix indicating the rule in which
// the error occurred. The original error is stored in the Inner field.
type parserError struct {
	Inner  error
	pos    position
	prefix string
}

// Error returns the error message.
func (p *parserError) Error() string {
	return p.prefix + ": " + p.Inner.Error()
}

// newParser creates a parser with the specified input source and options.
func newParser(filename string, b []byte, opts ...Option) *parser {
	p := &parser{
		filename: filename,
		errs:     new(errList),
		data:     b,
		pt:       savepoint{position: position{line: 1}},
		recover:  true,
	}
	p.setOptions(opts)
	return p
}

// setOptions applies the options to the parser.
func (p *parser) setOptions(opts []Option) {
	for _, opt := range opts {
		opt(p)
	}
}

type resultTuple struct {
	v   interface{}
	b   bool
	end savepoint
}

type parser struct {
	filename string
	pt       savepoint
	cur      current

	data []byte
	errs *errList

	recover bool
	debug   bool
	depth   int

	memoize bool
	// memoization table for the packrat algorithm:
	// map[offset in source] map[expression or rule] {value, match}
	memo map[int]map[interface{}]resultTuple

	// rules table, maps the rule identifier to the rule node
	rules map[string]*rule
	// variables stack, map of label to value
	vstack []map[string]interface{}
	// rule stack, allows identification of the current rule in errors
	rstack []*rule

	// stats
	exprCnt int
}

// push a variable set on the vstack.
func (p *parser) pushV() {
	if cap(p.vstack) == len(p.vstack) {
		// create new empty slot in the stack
		p.vstack = append(p.vstack, nil)
	} else {
		// slice to 1 more
		p.vstack = p.vstack[:len(p.vstack)+1]
	}

	// get the last args set
	m := p.vstack[len(p.vstack)-1]
	if m != nil && len(m) == 0 {
		// empty map, all good
		return
	}

	m = make(map[string]interface{})
	p.vstack[len(p.vstack)-1] = m
}

// pop a variable set from the vstack.
func (p *parser) popV() {
	// if the map is not empty, clear it
	m := p.vstack[len(p.vstack)-1]
	if len(m) > 0 {
		// GC that map
		p.vstack[len(p.vstack)-1] = nil
	}
	p.vstack = p.vstack[:len(p.vstack)-1]
}

func (p *parser) print(prefix, s string) string {
	if !p.debug {
		return s
	}

	fmt.Printf("%s %d:%d:%d: %s [%#U]\n",
		prefix, p.pt.line, p.pt.col, p.pt.offset, s, p.pt.rn)
	return s
}

func (p *parser) in(s string) string {
	p.depth++
	return p.print(strings.Repeat(" ", p.depth)+">", s)
}

func (p *parser) out(s string) string {
	p.depth--
	return p.print(strings.Repeat(" ", p.depth)+"<", s)
}

func (p *parser) addErr(err error) {
	p.addErrAt(err, p.pt.position)
}

func (p *parser) addErrAt(err error, pos position) {
	var buf bytes.Buffer
	if p.filename != "" {
		buf.WriteString(p.filename)
	}
	if buf.Len() > 0 {
		buf.WriteString(":")
	}
	buf.WriteString(fmt.Sprintf("%d:%d (%d)", pos.line, pos.col, pos.offset))
	if len(p.rstack) > 0 {
		if buf.Len() > 0 {
			buf.WriteString(": ")
		}
		rule := p.rstack[len(p.rstack)-1]
		if rule.displayName != "" {
			buf.WriteString("rule " + rule.displayName)
		} else {
			buf.WriteString("rule " + rule.name)
		}
	}
	pe := &parserError{Inner: err, prefix: buf.String()}
	p.errs.add(pe)
}

// read advances the parser to the next rune.
func (p *parser) read() {
	p.pt.offset += p.pt.w
	rn, n := utf8.DecodeRune(p.data[p.pt.offset:])
	p.pt.rn = rn
	p.pt.w = n
	p.pt.col++
	if rn == '\n' {
		p.pt.line++
		p.pt.col = 0
	}

	if rn == utf8.RuneError {
		if n > 0 {
			p.addErr(errInvalidEncoding)
		}
	}
}

// restore parser position to the savepoint pt.
func (p *parser) restore(pt savepoint) {
	if p.debug {
		defer p.out(p.in("restore"))
	}
	if pt.offset == p.pt.offset {
		return
	}
	p.pt = pt
}

// get the slice of bytes from the savepoint start to the current position.
func (p *parser) sliceFrom(start savepoint) []byte {
	return p.data[start.position.offset:p.pt.position.offset]
}

func (p *parser) getMemoized(node interface{}) (resultTuple, bool) {
	if len(p.memo) == 0 {
		return resultTuple{}, false
	}
	m := p.memo[p.pt.offset]
	if len(m) == 0 {
		return resultTuple{}, false
	}
	res, ok := m[node]
	return res, ok
}

func (p *parser) setMemoized(pt savepoint, node interface{}, tuple resultTuple) {
	if p.memo == nil {
		p.memo = make(map[int]map[interface{}]resultTuple)
	}
	m := p.memo[pt.offset]
	if m == nil {
		m = make(map[interface{}]resultTuple)
		p.memo[pt.offset] = m
	}
	m[node] = tuple
}

func (p *parser) buildRulesTable(g *grammar) {
	p.rules = make(map[string]*rule, len(g.rules))
	for _, r := range g.rules {
		p.rules[r.name] = r
	}
}

func (p *parser) parse(g *grammar) (val interface{}, err error) {
	if len(g.rules) == 0 {
		p.addErr(errNoRule)
		return nil, p.errs.err()
	}

	// TODO : not super critical but this could be generated
	p.buildRulesTable(g)

	if p.recover {
		// panic can be used in action code to stop parsing immediately
		// and return the panic as an error.
		defer func() {
			if e := recover(); e != nil {
				if p.debug {
					defer p.out(p.in("panic handler"))
				}
				val = nil
				switch e := e.(type) {
				case error:
					p.addErr(e)
				default:
					p.addErr(fmt.Errorf("%v", e))
				}
				err = p.errs.err()
			}
		}()
	}

	// start rule is rule [0]
	p.read() // advance to first rune
	val, ok := p.parseRule(g.rules[0])
	if !ok {
		if len(*p.errs) == 0 {
			// make sure this doesn't go out silently
			p.addErr(errNoMatch)
		}
		return nil, p.errs.err()
	}
	return val, p.errs.err()
}

func (p *parser) parseRule(rule *rule) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRule " + rule.name))
	}

	if p.memoize {
		res, ok := p.getMemoized(rule)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
	}

	start := p.pt
	p.rstack = append(p.rstack, rule)
	p.pushV()
	val, ok := p.parseExpr(rule.expr)
	p.popV()
	p.rstack = p.rstack[:len(p.rstack)-1]
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}

	if p.memoize {
		p.setMemoized(start, rule, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseExpr(expr interface{}) (interface{}, bool) {
	var pt savepoint
	var ok bool

	if p.memoize {
		res, ok := p.getMemoized(expr)
		if ok {
			p.restore(res.end)
			return res.v, res.b
		}
		pt = p.pt
	}

	p.exprCnt++
	var val interface{}
	switch expr := expr.(type) {
	case *actionExpr:
		val, ok = p.parseActionExpr(expr)
	case *andCodeExpr:
		val, ok = p.parseAndCodeExpr(expr)
	case *andExpr:
		val, ok = p.parseAndExpr(expr)
	case *anyMatcher:
		val, ok = p.parseAnyMatcher(expr)
	case *charClassMatcher:
		val, ok = p.parseCharClassMatcher(expr)
	case *choiceExpr:
		val, ok = p.parseChoiceExpr(expr)
	case *labeledExpr:
		val, ok = p.parseLabeledExpr(expr)
	case *litMatcher:
		val, ok = p.parseLitMatcher(expr)
	case *notCodeExpr:
		val, ok = p.parseNotCodeExpr(expr)
	case *notExpr:
		val, ok = p.parseNotExpr(expr)
	case *oneOrMoreExpr:
		val, ok = p.parseOneOrMoreExpr(expr)
	case *ruleRefExpr:
		val, ok = p.parseRuleRefExpr(expr)
	case *seqExpr:
		val, ok = p.parseSeqExpr(expr)
	case *zeroOrMoreExpr:
		val, ok = p.parseZeroOrMoreExpr(expr)
	case *zeroOrOneExpr:
		val, ok = p.parseZeroOrOneExpr(expr)
	default:
		panic(fmt.Sprintf("unknown expression type %T", expr))
	}
	if p.memoize {
		p.setMemoized(pt, expr, resultTuple{val, ok, p.pt})
	}
	return val, ok
}

func (p *parser) parseActionExpr(act *actionExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseActionExpr"))
	}

	start := p.pt
	val, ok := p.parseExpr(act.expr)
	if ok {
		p.cur.pos = start.position
		p.cur.text = p.sliceFrom(start)
		actVal, err := act.run(p)
		if err != nil {
			p.addErrAt(err, start.position)
		}
		val = actVal
	}
	if ok && p.debug {
		p.print(strings.Repeat(" ", p.depth)+"MATCH", string(p.sliceFrom(start)))
	}
	return val, ok
}

func (p *parser) parseAndCodeExpr(and *andCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndCodeExpr"))
	}

	ok, err := and.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, ok
}

func (p *parser) parseAndExpr(and *andExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAndExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(and.expr)
	p.popV()
	p.restore(pt)
	return nil, ok
}

func (p *parser) parseAnyMatcher(any *anyMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseAnyMatcher"))
	}

	if p.pt.rn != utf8.RuneError {
		start := p.pt
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseCharClassMatcher(chr *charClassMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseCharClassMatcher"))
	}

	cur := p.pt.rn
	start := p.pt
	if chr.ignoreCase {
		cur = unicode.ToLower(cur)
	}

	// try to match in the list of available chars
	for _, rn := range chr.chars {
		if rn == cur {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of ranges
	for i := 0; i < len(chr.ranges); i += 2 {
		if cur >= chr.ranges[i] && cur <= chr.ranges[i+1] {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	// try to match in the list of Unicode classes
	for _, cl := range chr.classes {
		if unicode.Is(cl, cur) {
			if chr.inverted {
				return nil, false
			}
			p.read()
			return p.sliceFrom(start), true
		}
	}

	if chr.inverted {
		p.read()
		return p.sliceFrom(start), true
	}
	return nil, false
}

func (p *parser) parseChoiceExpr(ch *choiceExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseChoiceExpr"))
	}

	for _, alt := range ch.alternatives {
		p.pushV()
		val, ok := p.parseExpr(alt)
		p.popV()
		if ok {
			return val, ok
		}
	}
	return nil, false
}

func (p *parser) parseLabeledExpr(lab *labeledExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLabeledExpr"))
	}

	p.pushV()
	val, ok := p.parseExpr(lab.expr)
	p.popV()
	if ok && lab.label != "" {
		m := p.vstack[len(p.vstack)-1]
		m[lab.label] = val
	}
	return val, ok
}

func (p *parser) parseLitMatcher(lit *litMatcher) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseLitMatcher"))
	}

	start := p.pt
	for _, want := range lit.val {
		cur := p.pt.rn
		if lit.ignoreCase {
			cur = unicode.ToLower(cur)
		}
		if cur != want {
			p.restore(start)
			return nil, false
		}
		p.read()
	}
	return p.sliceFrom(start), true
}

func (p *parser) parseNotCodeExpr(not *notCodeExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotCodeExpr"))
	}

	ok, err := not.run(p)
	if err != nil {
		p.addErr(err)
	}
	return nil, !ok
}

func (p *parser) parseNotExpr(not *notExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseNotExpr"))
	}

	pt := p.pt
	p.pushV()
	_, ok := p.parseExpr(not.expr)
	p.popV()
	p.restore(pt)
	return nil, !ok
}

func (p *parser) parseOneOrMoreExpr(expr *oneOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseOneOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			if len(vals) == 0 {
				// did not match once, no match
				return nil, false
			}
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseRuleRefExpr(ref *ruleRefExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseRuleRefExpr " + ref.name))
	}

	if ref.name == "" {
		panic(fmt.Sprintf("%s: invalid rule: missing name", ref.pos))
	}

	rule := p.rules[ref.name]
	if rule == nil {
		p.addErr(fmt.Errorf("undefined rule: %s", ref.name))
		return nil, false
	}
	return p.parseRule(rule)
}

func (p *parser) parseSeqExpr(seq *seqExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseSeqExpr"))
	}

	var vals []interface{}

	pt := p.pt
	for _, expr := range seq.exprs {
		val, ok := p.parseExpr(expr)
		if !ok {
			p.restore(pt)
			return nil, false
		}
		vals = append(vals, val)
	}
	return vals, true
}

func (p *parser) parseZeroOrMoreExpr(expr *zeroOrMoreExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrMoreExpr"))
	}

	var vals []interface{}

	for {
		p.pushV()
		val, ok := p.parseExpr(expr.expr)
		p.popV()
		if !ok {
			return vals, true
		}
		vals = append(vals, val)
	}
}

func (p *parser) parseZeroOrOneExpr(expr *zeroOrOneExpr) (interface{}, bool) {
	if p.debug {
		defer p.out(p.in("parseZeroOrOneExpr"))
	}

	p.pushV()
	val, _ := p.parseExpr(expr.expr)
	p.popV()
	// whether it matched or not, consider it a match
	return val, true
}

func rangeTable(class string) *unicode.RangeTable {
	if rt, ok := unicode.Categories[class]; ok {
		return rt
	}
	if rt, ok := unicode.Properties[class]; ok {
		return rt
	}
	if rt, ok := unicode.Scripts[class]; ok {
		return rt
	}

	// cannot happen
	panic(fmt.Sprintf("invalid Unicode class: %s", class))
}
