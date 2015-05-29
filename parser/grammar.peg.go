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
	language  string
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
			pos:  position{line: 41, col: 1, offset: 517},
			expr: &actionExpr{
				pos: position{line: 41, col: 11, offset: 529},
				run: (*parser).callonGrammer1,
				expr: &seqExpr{
					pos: position{line: 41, col: 11, offset: 529},
					exprs: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 41, col: 11, offset: 529},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 41, col: 14, offset: 532},
							label: "statements",
							expr: &zeroOrMoreExpr{
								pos: position{line: 41, col: 25, offset: 543},
								expr: &seqExpr{
									pos: position{line: 41, col: 27, offset: 545},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 41, col: 27, offset: 545},
											name: "Statement",
										},
										&ruleRefExpr{
											pos:  position{line: 41, col: 37, offset: 555},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 41, col: 44, offset: 562},
							alternatives: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 41, col: 44, offset: 562},
									name: "EOF",
								},
								&ruleRefExpr{
									pos:  position{line: 41, col: 50, offset: 568},
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
			pos:  position{line: 82, col: 1, offset: 1632},
			expr: &actionExpr{
				pos: position{line: 82, col: 15, offset: 1648},
				run: (*parser).callonSyntaxError1,
				expr: &anyMatcher{
					line: 82, col: 15, offset: 1648,
				},
			},
		},
		{
			name: "Include",
			pos:  position{line: 86, col: 1, offset: 1703},
			expr: &actionExpr{
				pos: position{line: 86, col: 11, offset: 1715},
				run: (*parser).callonInclude1,
				expr: &seqExpr{
					pos: position{line: 86, col: 11, offset: 1715},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 86, col: 11, offset: 1715},
							val:        "include",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 21, offset: 1725},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 86, col: 23, offset: 1727},
							label: "file",
							expr: &ruleRefExpr{
								pos:  position{line: 86, col: 28, offset: 1732},
								name: "Literal",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 86, col: 36, offset: 1740},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Statement",
			pos:  position{line: 90, col: 1, offset: 1785},
			expr: &choiceExpr{
				pos: position{line: 90, col: 13, offset: 1799},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 90, col: 13, offset: 1799},
						name: "Include",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 23, offset: 1809},
						name: "Namespace",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 35, offset: 1821},
						name: "Const",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 43, offset: 1829},
						name: "Enum",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 50, offset: 1836},
						name: "TypeDef",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 60, offset: 1846},
						name: "Struct",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 69, offset: 1855},
						name: "Exception",
					},
					&ruleRefExpr{
						pos:  position{line: 90, col: 81, offset: 1867},
						name: "Service",
					},
				},
			},
		},
		{
			name: "Namespace",
			pos:  position{line: 92, col: 1, offset: 1876},
			expr: &actionExpr{
				pos: position{line: 92, col: 13, offset: 1890},
				run: (*parser).callonNamespace1,
				expr: &seqExpr{
					pos: position{line: 92, col: 13, offset: 1890},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 92, col: 13, offset: 1890},
							val:        "namespace",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 92, col: 25, offset: 1902},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 92, col: 27, offset: 1904},
							label: "language",
							expr: &oneOrMoreExpr{
								pos: position{line: 92, col: 36, offset: 1913},
								expr: &charClassMatcher{
									pos:        position{line: 92, col: 36, offset: 1913},
									val:        "[a-z]",
									ranges:     []rune{'a', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 92, col: 43, offset: 1920},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 92, col: 45, offset: 1922},
							label: "ns",
							expr: &oneOrMoreExpr{
								pos: position{line: 92, col: 48, offset: 1925},
								expr: &charClassMatcher{
									pos:        position{line: 92, col: 48, offset: 1925},
									val:        "[a-zA-z.]",
									chars:      []rune{'.'},
									ranges:     []rune{'a', 'z', 'A', 'z'},
									ignoreCase: false,
									inverted:   false,
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 92, col: 59, offset: 1936},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Const",
			pos:  position{line: 99, col: 1, offset: 2052},
			expr: &actionExpr{
				pos: position{line: 99, col: 9, offset: 2062},
				run: (*parser).callonConst1,
				expr: &seqExpr{
					pos: position{line: 99, col: 9, offset: 2062},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 99, col: 9, offset: 2062},
							val:        "const",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 17, offset: 2070},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 99, col: 19, offset: 2072},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 99, col: 23, offset: 2076},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 33, offset: 2086},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 99, col: 35, offset: 2088},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 99, col: 40, offset: 2093},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 51, offset: 2104},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 99, col: 53, offset: 2106},
							val:        "=",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 57, offset: 2110},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 99, col: 59, offset: 2112},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 99, col: 65, offset: 2118},
								name: "ConstValue",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 99, col: 76, offset: 2129},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Enum",
			pos:  position{line: 107, col: 1, offset: 2237},
			expr: &actionExpr{
				pos: position{line: 107, col: 8, offset: 2246},
				run: (*parser).callonEnum1,
				expr: &seqExpr{
					pos: position{line: 107, col: 8, offset: 2246},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 107, col: 8, offset: 2246},
							val:        "enum",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 15, offset: 2253},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 107, col: 17, offset: 2255},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 107, col: 22, offset: 2260},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 33, offset: 2271},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 107, col: 35, offset: 2273},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 39, offset: 2277},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 107, col: 42, offset: 2280},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 107, col: 49, offset: 2287},
								expr: &seqExpr{
									pos: position{line: 107, col: 50, offset: 2288},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 107, col: 50, offset: 2288},
											name: "EnumValue",
										},
										&ruleRefExpr{
											pos:  position{line: 107, col: 60, offset: 2298},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 107, col: 65, offset: 2303},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 107, col: 69, offset: 2307},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EnumValue",
			pos:  position{line: 130, col: 1, offset: 2823},
			expr: &actionExpr{
				pos: position{line: 130, col: 13, offset: 2837},
				run: (*parser).callonEnumValue1,
				expr: &seqExpr{
					pos: position{line: 130, col: 13, offset: 2837},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 130, col: 13, offset: 2837},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 130, col: 18, offset: 2842},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 130, col: 29, offset: 2853},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 130, col: 31, offset: 2855},
							label: "value",
							expr: &zeroOrOneExpr{
								pos: position{line: 130, col: 37, offset: 2861},
								expr: &seqExpr{
									pos: position{line: 130, col: 38, offset: 2862},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 130, col: 38, offset: 2862},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 130, col: 42, offset: 2866},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 130, col: 44, offset: 2868},
											name: "IntConstant",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 130, col: 58, offset: 2882},
							expr: &ruleRefExpr{
								pos:  position{line: 130, col: 58, offset: 2882},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "TypeDef",
			pos:  position{line: 141, col: 1, offset: 3061},
			expr: &actionExpr{
				pos: position{line: 141, col: 11, offset: 3073},
				run: (*parser).callonTypeDef1,
				expr: &seqExpr{
					pos: position{line: 141, col: 11, offset: 3073},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 141, col: 11, offset: 3073},
							val:        "typedef",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 141, col: 21, offset: 3083},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 141, col: 23, offset: 3085},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 141, col: 27, offset: 3089},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 141, col: 37, offset: 3099},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 141, col: 39, offset: 3101},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 141, col: 44, offset: 3106},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 141, col: 55, offset: 3117},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "Struct",
			pos:  position{line: 148, col: 1, offset: 3207},
			expr: &actionExpr{
				pos: position{line: 148, col: 10, offset: 3218},
				run: (*parser).callonStruct1,
				expr: &seqExpr{
					pos: position{line: 148, col: 10, offset: 3218},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 148, col: 10, offset: 3218},
							val:        "struct",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 148, col: 19, offset: 3227},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 148, col: 21, offset: 3229},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 148, col: 24, offset: 3232},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "Exception",
			pos:  position{line: 149, col: 1, offset: 3272},
			expr: &actionExpr{
				pos: position{line: 149, col: 13, offset: 3286},
				run: (*parser).callonException1,
				expr: &seqExpr{
					pos: position{line: 149, col: 13, offset: 3286},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 149, col: 13, offset: 3286},
							val:        "exception",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 149, col: 25, offset: 3298},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 149, col: 27, offset: 3300},
							label: "st",
							expr: &ruleRefExpr{
								pos:  position{line: 149, col: 30, offset: 3303},
								name: "StructLike",
							},
						},
					},
				},
			},
		},
		{
			name: "StructLike",
			pos:  position{line: 150, col: 1, offset: 3354},
			expr: &actionExpr{
				pos: position{line: 150, col: 14, offset: 3369},
				run: (*parser).callonStructLike1,
				expr: &seqExpr{
					pos: position{line: 150, col: 14, offset: 3369},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 150, col: 14, offset: 3369},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 150, col: 19, offset: 3374},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 150, col: 30, offset: 3385},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 150, col: 32, offset: 3387},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 150, col: 36, offset: 3391},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 150, col: 39, offset: 3394},
							label: "fields",
							expr: &ruleRefExpr{
								pos:  position{line: 150, col: 46, offset: 3401},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 150, col: 56, offset: 3411},
							val:        "}",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 150, col: 60, offset: 3415},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "FieldList",
			pos:  position{line: 160, col: 1, offset: 3549},
			expr: &actionExpr{
				pos: position{line: 160, col: 13, offset: 3563},
				run: (*parser).callonFieldList1,
				expr: &labeledExpr{
					pos:   position{line: 160, col: 13, offset: 3563},
					label: "fields",
					expr: &zeroOrMoreExpr{
						pos: position{line: 160, col: 20, offset: 3570},
						expr: &seqExpr{
							pos: position{line: 160, col: 21, offset: 3571},
							exprs: []interface{}{
								&ruleRefExpr{
									pos:  position{line: 160, col: 21, offset: 3571},
									name: "Field",
								},
								&ruleRefExpr{
									pos:  position{line: 160, col: 27, offset: 3577},
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
			pos:  position{line: 169, col: 1, offset: 3737},
			expr: &actionExpr{
				pos: position{line: 169, col: 9, offset: 3747},
				run: (*parser).callonField1,
				expr: &seqExpr{
					pos: position{line: 169, col: 9, offset: 3747},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 169, col: 9, offset: 3747},
							label: "id",
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 12, offset: 3750},
								name: "IntConstant",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 24, offset: 3762},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 169, col: 26, offset: 3764},
							val:        ":",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 30, offset: 3768},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 169, col: 32, offset: 3770},
							label: "req",
							expr: &zeroOrOneExpr{
								pos: position{line: 169, col: 36, offset: 3774},
								expr: &ruleRefExpr{
									pos:  position{line: 169, col: 36, offset: 3774},
									name: "FieldReq",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 46, offset: 3784},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 169, col: 48, offset: 3786},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 52, offset: 3790},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 62, offset: 3800},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 169, col: 64, offset: 3802},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 69, offset: 3807},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 169, col: 80, offset: 3818},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 169, col: 82, offset: 3820},
							label: "def",
							expr: &zeroOrOneExpr{
								pos: position{line: 169, col: 86, offset: 3824},
								expr: &seqExpr{
									pos: position{line: 169, col: 87, offset: 3825},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 169, col: 87, offset: 3825},
											val:        "=",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 169, col: 91, offset: 3829},
											name: "_",
										},
										&ruleRefExpr{
											pos:  position{line: 169, col: 93, offset: 3831},
											name: "ConstValue",
										},
									},
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 169, col: 106, offset: 3844},
							expr: &ruleRefExpr{
								pos:  position{line: 169, col: 106, offset: 3844},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FieldReq",
			pos:  position{line: 184, col: 1, offset: 4104},
			expr: &actionExpr{
				pos: position{line: 184, col: 12, offset: 4117},
				run: (*parser).callonFieldReq1,
				expr: &choiceExpr{
					pos: position{line: 184, col: 13, offset: 4118},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 184, col: 13, offset: 4118},
							val:        "required",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 184, col: 26, offset: 4131},
							val:        "optional",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Service",
			pos:  position{line: 188, col: 1, offset: 4202},
			expr: &actionExpr{
				pos: position{line: 188, col: 11, offset: 4214},
				run: (*parser).callonService1,
				expr: &seqExpr{
					pos: position{line: 188, col: 11, offset: 4214},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 188, col: 11, offset: 4214},
							val:        "service",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 21, offset: 4224},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 23, offset: 4226},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 188, col: 28, offset: 4231},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 39, offset: 4242},
							name: "_",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 41, offset: 4244},
							label: "extends",
							expr: &zeroOrOneExpr{
								pos: position{line: 188, col: 49, offset: 4252},
								expr: &seqExpr{
									pos: position{line: 188, col: 50, offset: 4253},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 188, col: 50, offset: 4253},
											val:        "extends",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 60, offset: 4263},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 63, offset: 4266},
											name: "Identifier",
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 74, offset: 4277},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 79, offset: 4282},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 188, col: 82, offset: 4285},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 86, offset: 4289},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 188, col: 89, offset: 4292},
							label: "methods",
							expr: &zeroOrMoreExpr{
								pos: position{line: 188, col: 97, offset: 4300},
								expr: &seqExpr{
									pos: position{line: 188, col: 98, offset: 4301},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 188, col: 98, offset: 4301},
											name: "Function",
										},
										&ruleRefExpr{
											pos:  position{line: 188, col: 107, offset: 4310},
											name: "__",
										},
									},
								},
							},
						},
						&choiceExpr{
							pos: position{line: 188, col: 113, offset: 4316},
							alternatives: []interface{}{
								&litMatcher{
									pos:        position{line: 188, col: 113, offset: 4316},
									val:        "}",
									ignoreCase: false,
								},
								&ruleRefExpr{
									pos:  position{line: 188, col: 119, offset: 4322},
									name: "EndOfServiceError",
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 188, col: 138, offset: 4341},
							name: "EOS",
						},
					},
				},
			},
		},
		{
			name: "EndOfServiceError",
			pos:  position{line: 203, col: 1, offset: 4682},
			expr: &actionExpr{
				pos: position{line: 203, col: 21, offset: 4704},
				run: (*parser).callonEndOfServiceError1,
				expr: &anyMatcher{
					line: 203, col: 21, offset: 4704,
				},
			},
		},
		{
			name: "Function",
			pos:  position{line: 207, col: 1, offset: 4770},
			expr: &actionExpr{
				pos: position{line: 207, col: 12, offset: 4783},
				run: (*parser).callonFunction1,
				expr: &seqExpr{
					pos: position{line: 207, col: 12, offset: 4783},
					exprs: []interface{}{
						&labeledExpr{
							pos:   position{line: 207, col: 12, offset: 4783},
							label: "oneway",
							expr: &zeroOrOneExpr{
								pos: position{line: 207, col: 19, offset: 4790},
								expr: &seqExpr{
									pos: position{line: 207, col: 20, offset: 4791},
									exprs: []interface{}{
										&litMatcher{
											pos:        position{line: 207, col: 20, offset: 4791},
											val:        "oneway",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 207, col: 29, offset: 4800},
											name: "__",
										},
									},
								},
							},
						},
						&labeledExpr{
							pos:   position{line: 207, col: 34, offset: 4805},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 38, offset: 4809},
								name: "FunctionType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 51, offset: 4822},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 207, col: 54, offset: 4825},
							label: "name",
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 59, offset: 4830},
								name: "Identifier",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 70, offset: 4841},
							name: "_",
						},
						&litMatcher{
							pos:        position{line: 207, col: 72, offset: 4843},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 76, offset: 4847},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 207, col: 79, offset: 4850},
							label: "arguments",
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 89, offset: 4860},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 207, col: 99, offset: 4870},
							val:        ")",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 207, col: 103, offset: 4874},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 207, col: 106, offset: 4877},
							label: "exceptions",
							expr: &zeroOrOneExpr{
								pos: position{line: 207, col: 117, offset: 4888},
								expr: &ruleRefExpr{
									pos:  position{line: 207, col: 117, offset: 4888},
									name: "Throws",
								},
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 207, col: 125, offset: 4896},
							expr: &ruleRefExpr{
								pos:  position{line: 207, col: 125, offset: 4896},
								name: "ListSeparator",
							},
						},
					},
				},
			},
		},
		{
			name: "FunctionType",
			pos:  position{line: 230, col: 1, offset: 5277},
			expr: &actionExpr{
				pos: position{line: 230, col: 16, offset: 5294},
				run: (*parser).callonFunctionType1,
				expr: &labeledExpr{
					pos:   position{line: 230, col: 16, offset: 5294},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 230, col: 21, offset: 5299},
						alternatives: []interface{}{
							&litMatcher{
								pos:        position{line: 230, col: 21, offset: 5299},
								val:        "void",
								ignoreCase: false,
							},
							&ruleRefExpr{
								pos:  position{line: 230, col: 30, offset: 5308},
								name: "FieldType",
							},
						},
					},
				},
			},
		},
		{
			name: "Throws",
			pos:  position{line: 237, col: 1, offset: 5415},
			expr: &actionExpr{
				pos: position{line: 237, col: 10, offset: 5426},
				run: (*parser).callonThrows1,
				expr: &seqExpr{
					pos: position{line: 237, col: 10, offset: 5426},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 237, col: 10, offset: 5426},
							val:        "throws",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 237, col: 19, offset: 5435},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 237, col: 22, offset: 5438},
							val:        "(",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 237, col: 26, offset: 5442},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 237, col: 29, offset: 5445},
							label: "exceptions",
							expr: &ruleRefExpr{
								pos:  position{line: 237, col: 40, offset: 5456},
								name: "FieldList",
							},
						},
						&litMatcher{
							pos:        position{line: 237, col: 50, offset: 5466},
							val:        ")",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "FieldType",
			pos:  position{line: 241, col: 1, offset: 5499},
			expr: &actionExpr{
				pos: position{line: 241, col: 13, offset: 5513},
				run: (*parser).callonFieldType1,
				expr: &labeledExpr{
					pos:   position{line: 241, col: 13, offset: 5513},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 241, col: 18, offset: 5518},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 241, col: 18, offset: 5518},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 241, col: 29, offset: 5529},
								name: "ContainerType",
							},
							&ruleRefExpr{
								pos:  position{line: 241, col: 45, offset: 5545},
								name: "Identifier",
							},
						},
					},
				},
			},
		},
		{
			name: "DefinitionType",
			pos:  position{line: 248, col: 1, offset: 5655},
			expr: &actionExpr{
				pos: position{line: 248, col: 18, offset: 5674},
				run: (*parser).callonDefinitionType1,
				expr: &labeledExpr{
					pos:   position{line: 248, col: 18, offset: 5674},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 248, col: 23, offset: 5679},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 248, col: 23, offset: 5679},
								name: "BaseType",
							},
							&ruleRefExpr{
								pos:  position{line: 248, col: 34, offset: 5690},
								name: "ContainerType",
							},
						},
					},
				},
			},
		},
		{
			name: "BaseType",
			pos:  position{line: 252, col: 1, offset: 5727},
			expr: &actionExpr{
				pos: position{line: 252, col: 12, offset: 5740},
				run: (*parser).callonBaseType1,
				expr: &choiceExpr{
					pos: position{line: 252, col: 13, offset: 5741},
					alternatives: []interface{}{
						&litMatcher{
							pos:        position{line: 252, col: 13, offset: 5741},
							val:        "bool",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 22, offset: 5750},
							val:        "byte",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 31, offset: 5759},
							val:        "i16",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 39, offset: 5767},
							val:        "i32",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 47, offset: 5775},
							val:        "i64",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 55, offset: 5783},
							val:        "double",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 66, offset: 5794},
							val:        "string",
							ignoreCase: false,
						},
						&litMatcher{
							pos:        position{line: 252, col: 77, offset: 5805},
							val:        "binary",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ContainerType",
			pos:  position{line: 256, col: 1, offset: 5862},
			expr: &actionExpr{
				pos: position{line: 256, col: 17, offset: 5880},
				run: (*parser).callonContainerType1,
				expr: &labeledExpr{
					pos:   position{line: 256, col: 17, offset: 5880},
					label: "typ",
					expr: &choiceExpr{
						pos: position{line: 256, col: 22, offset: 5885},
						alternatives: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 256, col: 22, offset: 5885},
								name: "MapType",
							},
							&ruleRefExpr{
								pos:  position{line: 256, col: 32, offset: 5895},
								name: "SetType",
							},
							&ruleRefExpr{
								pos:  position{line: 256, col: 42, offset: 5905},
								name: "ListType",
							},
						},
					},
				},
			},
		},
		{
			name: "MapType",
			pos:  position{line: 260, col: 1, offset: 5937},
			expr: &actionExpr{
				pos: position{line: 260, col: 11, offset: 5949},
				run: (*parser).callonMapType1,
				expr: &seqExpr{
					pos: position{line: 260, col: 11, offset: 5949},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 260, col: 11, offset: 5949},
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 11, offset: 5949},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 260, col: 20, offset: 5958},
							val:        "map<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 27, offset: 5965},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 260, col: 30, offset: 5968},
							label: "key",
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 34, offset: 5972},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 44, offset: 5982},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 260, col: 47, offset: 5985},
							val:        ",",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 51, offset: 5989},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 260, col: 54, offset: 5992},
							label: "value",
							expr: &ruleRefExpr{
								pos:  position{line: 260, col: 60, offset: 5998},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 260, col: 70, offset: 6008},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 260, col: 73, offset: 6011},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "SetType",
			pos:  position{line: 268, col: 1, offset: 6110},
			expr: &actionExpr{
				pos: position{line: 268, col: 11, offset: 6122},
				run: (*parser).callonSetType1,
				expr: &seqExpr{
					pos: position{line: 268, col: 11, offset: 6122},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 268, col: 11, offset: 6122},
							expr: &ruleRefExpr{
								pos:  position{line: 268, col: 11, offset: 6122},
								name: "CppType",
							},
						},
						&litMatcher{
							pos:        position{line: 268, col: 20, offset: 6131},
							val:        "set<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 268, col: 27, offset: 6138},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 268, col: 30, offset: 6141},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 268, col: 34, offset: 6145},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 268, col: 44, offset: 6155},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 268, col: 47, offset: 6158},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ListType",
			pos:  position{line: 275, col: 1, offset: 6231},
			expr: &actionExpr{
				pos: position{line: 275, col: 12, offset: 6244},
				run: (*parser).callonListType1,
				expr: &seqExpr{
					pos: position{line: 275, col: 12, offset: 6244},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 275, col: 12, offset: 6244},
							val:        "list<",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 20, offset: 6252},
							name: "WS",
						},
						&labeledExpr{
							pos:   position{line: 275, col: 23, offset: 6255},
							label: "typ",
							expr: &ruleRefExpr{
								pos:  position{line: 275, col: 27, offset: 6259},
								name: "FieldType",
							},
						},
						&ruleRefExpr{
							pos:  position{line: 275, col: 37, offset: 6269},
							name: "WS",
						},
						&litMatcher{
							pos:        position{line: 275, col: 40, offset: 6272},
							val:        ">",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "CppType",
			pos:  position{line: 282, col: 1, offset: 6346},
			expr: &actionExpr{
				pos: position{line: 282, col: 11, offset: 6358},
				run: (*parser).callonCppType1,
				expr: &seqExpr{
					pos: position{line: 282, col: 11, offset: 6358},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 282, col: 11, offset: 6358},
							val:        "cpp_type",
							ignoreCase: false,
						},
						&labeledExpr{
							pos:   position{line: 282, col: 22, offset: 6369},
							label: "cppType",
							expr: &ruleRefExpr{
								pos:  position{line: 282, col: 30, offset: 6377},
								name: "Literal",
							},
						},
					},
				},
			},
		},
		{
			name: "ConstValue",
			pos:  position{line: 286, col: 1, offset: 6411},
			expr: &choiceExpr{
				pos: position{line: 286, col: 14, offset: 6426},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 286, col: 14, offset: 6426},
						name: "Literal",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 24, offset: 6436},
						name: "DoubleConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 41, offset: 6453},
						name: "IntConstant",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 55, offset: 6467},
						name: "ConstMap",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 66, offset: 6478},
						name: "ConstList",
					},
					&ruleRefExpr{
						pos:  position{line: 286, col: 78, offset: 6490},
						name: "Identifier",
					},
				},
			},
		},
		{
			name: "IntConstant",
			pos:  position{line: 288, col: 1, offset: 6502},
			expr: &actionExpr{
				pos: position{line: 288, col: 15, offset: 6518},
				run: (*parser).callonIntConstant1,
				expr: &seqExpr{
					pos: position{line: 288, col: 15, offset: 6518},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 288, col: 15, offset: 6518},
							expr: &charClassMatcher{
								pos:        position{line: 288, col: 15, offset: 6518},
								val:        "[-+]",
								chars:      []rune{'-', '+'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&oneOrMoreExpr{
							pos: position{line: 288, col: 21, offset: 6524},
							expr: &ruleRefExpr{
								pos:  position{line: 288, col: 21, offset: 6524},
								name: "Digit",
							},
						},
					},
				},
			},
		},
		{
			name: "DoubleConstant",
			pos:  position{line: 292, col: 1, offset: 6585},
			expr: &actionExpr{
				pos: position{line: 292, col: 18, offset: 6604},
				run: (*parser).callonDoubleConstant1,
				expr: &seqExpr{
					pos: position{line: 292, col: 18, offset: 6604},
					exprs: []interface{}{
						&zeroOrOneExpr{
							pos: position{line: 292, col: 18, offset: 6604},
							expr: &charClassMatcher{
								pos:        position{line: 292, col: 18, offset: 6604},
								val:        "[+-]",
								chars:      []rune{'+', '-'},
								ignoreCase: false,
								inverted:   false,
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 292, col: 24, offset: 6610},
							expr: &ruleRefExpr{
								pos:  position{line: 292, col: 24, offset: 6610},
								name: "Digit",
							},
						},
						&litMatcher{
							pos:        position{line: 292, col: 31, offset: 6617},
							val:        ".",
							ignoreCase: false,
						},
						&zeroOrMoreExpr{
							pos: position{line: 292, col: 35, offset: 6621},
							expr: &ruleRefExpr{
								pos:  position{line: 292, col: 35, offset: 6621},
								name: "Digit",
							},
						},
						&zeroOrOneExpr{
							pos: position{line: 292, col: 42, offset: 6628},
							expr: &seqExpr{
								pos: position{line: 292, col: 44, offset: 6630},
								exprs: []interface{}{
									&charClassMatcher{
										pos:        position{line: 292, col: 44, offset: 6630},
										val:        "['Ee']",
										chars:      []rune{'\'', 'E', 'e', '\''},
										ignoreCase: false,
										inverted:   false,
									},
									&ruleRefExpr{
										pos:  position{line: 292, col: 51, offset: 6637},
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
			pos:  position{line: 296, col: 1, offset: 6704},
			expr: &actionExpr{
				pos: position{line: 296, col: 13, offset: 6718},
				run: (*parser).callonConstList1,
				expr: &seqExpr{
					pos: position{line: 296, col: 13, offset: 6718},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 296, col: 13, offset: 6718},
							val:        "[",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 17, offset: 6722},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 296, col: 20, offset: 6725},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 296, col: 27, offset: 6732},
								expr: &seqExpr{
									pos: position{line: 296, col: 28, offset: 6733},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 296, col: 28, offset: 6733},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 296, col: 39, offset: 6744},
											name: "__",
										},
										&zeroOrOneExpr{
											pos: position{line: 296, col: 42, offset: 6747},
											expr: &ruleRefExpr{
												pos:  position{line: 296, col: 42, offset: 6747},
												name: "ListSeparator",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 296, col: 57, offset: 6762},
											name: "__",
										},
									},
								},
							},
						},
						&ruleRefExpr{
							pos:  position{line: 296, col: 62, offset: 6767},
							name: "__",
						},
						&litMatcher{
							pos:        position{line: 296, col: 65, offset: 6770},
							val:        "]",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "ConstMap",
			pos:  position{line: 305, col: 1, offset: 6943},
			expr: &actionExpr{
				pos: position{line: 305, col: 12, offset: 6956},
				run: (*parser).callonConstMap1,
				expr: &seqExpr{
					pos: position{line: 305, col: 12, offset: 6956},
					exprs: []interface{}{
						&litMatcher{
							pos:        position{line: 305, col: 12, offset: 6956},
							val:        "{",
							ignoreCase: false,
						},
						&ruleRefExpr{
							pos:  position{line: 305, col: 16, offset: 6960},
							name: "__",
						},
						&labeledExpr{
							pos:   position{line: 305, col: 19, offset: 6963},
							label: "values",
							expr: &zeroOrMoreExpr{
								pos: position{line: 305, col: 26, offset: 6970},
								expr: &seqExpr{
									pos: position{line: 305, col: 27, offset: 6971},
									exprs: []interface{}{
										&ruleRefExpr{
											pos:  position{line: 305, col: 27, offset: 6971},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 38, offset: 6982},
											name: "__",
										},
										&litMatcher{
											pos:        position{line: 305, col: 41, offset: 6985},
											val:        ":",
											ignoreCase: false,
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 45, offset: 6989},
											name: "__",
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 48, offset: 6992},
											name: "ConstValue",
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 59, offset: 7003},
											name: "__",
										},
										&choiceExpr{
											pos: position{line: 305, col: 63, offset: 7007},
											alternatives: []interface{}{
												&litMatcher{
													pos:        position{line: 305, col: 63, offset: 7007},
													val:        ",",
													ignoreCase: false,
												},
												&andExpr{
													pos: position{line: 305, col: 69, offset: 7013},
													expr: &litMatcher{
														pos:        position{line: 305, col: 70, offset: 7014},
														val:        "}",
														ignoreCase: false,
													},
												},
											},
										},
										&ruleRefExpr{
											pos:  position{line: 305, col: 75, offset: 7019},
											name: "__",
										},
									},
								},
							},
						},
						&litMatcher{
							pos:        position{line: 305, col: 80, offset: 7024},
							val:        "}",
							ignoreCase: false,
						},
					},
				},
			},
		},
		{
			name: "Literal",
			pos:  position{line: 321, col: 1, offset: 7270},
			expr: &actionExpr{
				pos: position{line: 321, col: 11, offset: 7282},
				run: (*parser).callonLiteral1,
				expr: &choiceExpr{
					pos: position{line: 321, col: 12, offset: 7283},
					alternatives: []interface{}{
						&seqExpr{
							pos: position{line: 321, col: 13, offset: 7284},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 321, col: 13, offset: 7284},
									val:        "\"",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 321, col: 17, offset: 7288},
									expr: &choiceExpr{
										pos: position{line: 321, col: 18, offset: 7289},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 321, col: 18, offset: 7289},
												val:        "\\\"",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 321, col: 25, offset: 7296},
												val:        "[^\"]",
												chars:      []rune{'"'},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 321, col: 32, offset: 7303},
									val:        "\"",
									ignoreCase: false,
								},
							},
						},
						&seqExpr{
							pos: position{line: 321, col: 40, offset: 7311},
							exprs: []interface{}{
								&litMatcher{
									pos:        position{line: 321, col: 40, offset: 7311},
									val:        "'",
									ignoreCase: false,
								},
								&zeroOrMoreExpr{
									pos: position{line: 321, col: 45, offset: 7316},
									expr: &choiceExpr{
										pos: position{line: 321, col: 46, offset: 7317},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 321, col: 46, offset: 7317},
												val:        "\\'",
												ignoreCase: false,
											},
											&charClassMatcher{
												pos:        position{line: 321, col: 53, offset: 7324},
												val:        "[^']",
												chars:      []rune{'\''},
												ignoreCase: false,
												inverted:   true,
											},
										},
									},
								},
								&litMatcher{
									pos:        position{line: 321, col: 60, offset: 7331},
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
			pos:  position{line: 328, col: 1, offset: 7532},
			expr: &actionExpr{
				pos: position{line: 328, col: 14, offset: 7547},
				run: (*parser).callonIdentifier1,
				expr: &seqExpr{
					pos: position{line: 328, col: 14, offset: 7547},
					exprs: []interface{}{
						&oneOrMoreExpr{
							pos: position{line: 328, col: 14, offset: 7547},
							expr: &choiceExpr{
								pos: position{line: 328, col: 15, offset: 7548},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 328, col: 15, offset: 7548},
										name: "Letter",
									},
									&litMatcher{
										pos:        position{line: 328, col: 24, offset: 7557},
										val:        "_",
										ignoreCase: false,
									},
								},
							},
						},
						&zeroOrMoreExpr{
							pos: position{line: 328, col: 30, offset: 7563},
							expr: &choiceExpr{
								pos: position{line: 328, col: 31, offset: 7564},
								alternatives: []interface{}{
									&ruleRefExpr{
										pos:  position{line: 328, col: 31, offset: 7564},
										name: "Letter",
									},
									&ruleRefExpr{
										pos:  position{line: 328, col: 40, offset: 7573},
										name: "Digit",
									},
									&charClassMatcher{
										pos:        position{line: 328, col: 48, offset: 7581},
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
			pos:  position{line: 332, col: 1, offset: 7633},
			expr: &charClassMatcher{
				pos:        position{line: 332, col: 17, offset: 7651},
				val:        "[,;]",
				chars:      []rune{',', ';'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Letter",
			pos:  position{line: 333, col: 1, offset: 7656},
			expr: &charClassMatcher{
				pos:        position{line: 333, col: 10, offset: 7667},
				val:        "[A-Za-z]",
				ranges:     []rune{'A', 'Z', 'a', 'z'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "Digit",
			pos:  position{line: 334, col: 1, offset: 7676},
			expr: &charClassMatcher{
				pos:        position{line: 334, col: 9, offset: 7686},
				val:        "[0-9]",
				ranges:     []rune{'0', '9'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "SourceChar",
			pos:  position{line: 338, col: 1, offset: 7697},
			expr: &anyMatcher{
				line: 338, col: 14, offset: 7712,
			},
		},
		{
			name: "Comment",
			pos:  position{line: 339, col: 1, offset: 7714},
			expr: &choiceExpr{
				pos: position{line: 339, col: 11, offset: 7726},
				alternatives: []interface{}{
					&ruleRefExpr{
						pos:  position{line: 339, col: 11, offset: 7726},
						name: "MultiLineComment",
					},
					&ruleRefExpr{
						pos:  position{line: 339, col: 30, offset: 7745},
						name: "SingleLineComment",
					},
				},
			},
		},
		{
			name: "MultiLineComment",
			pos:  position{line: 340, col: 1, offset: 7763},
			expr: &seqExpr{
				pos: position{line: 340, col: 20, offset: 7784},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 340, col: 20, offset: 7784},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 340, col: 25, offset: 7789},
						expr: &seqExpr{
							pos: position{line: 340, col: 27, offset: 7791},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 340, col: 27, offset: 7791},
									expr: &litMatcher{
										pos:        position{line: 340, col: 28, offset: 7792},
										val:        "*/",
										ignoreCase: false,
									},
								},
								&ruleRefExpr{
									pos:  position{line: 340, col: 33, offset: 7797},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 340, col: 47, offset: 7811},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "MultiLineCommentNoLineTerminator",
			pos:  position{line: 341, col: 1, offset: 7816},
			expr: &seqExpr{
				pos: position{line: 341, col: 36, offset: 7853},
				exprs: []interface{}{
					&litMatcher{
						pos:        position{line: 341, col: 36, offset: 7853},
						val:        "/*",
						ignoreCase: false,
					},
					&zeroOrMoreExpr{
						pos: position{line: 341, col: 41, offset: 7858},
						expr: &seqExpr{
							pos: position{line: 341, col: 43, offset: 7860},
							exprs: []interface{}{
								&notExpr{
									pos: position{line: 341, col: 43, offset: 7860},
									expr: &choiceExpr{
										pos: position{line: 341, col: 46, offset: 7863},
										alternatives: []interface{}{
											&litMatcher{
												pos:        position{line: 341, col: 46, offset: 7863},
												val:        "*/",
												ignoreCase: false,
											},
											&ruleRefExpr{
												pos:  position{line: 341, col: 53, offset: 7870},
												name: "EOL",
											},
										},
									},
								},
								&ruleRefExpr{
									pos:  position{line: 341, col: 59, offset: 7876},
									name: "SourceChar",
								},
							},
						},
					},
					&litMatcher{
						pos:        position{line: 341, col: 73, offset: 7890},
						val:        "*/",
						ignoreCase: false,
					},
				},
			},
		},
		{
			name: "SingleLineComment",
			pos:  position{line: 342, col: 1, offset: 7895},
			expr: &choiceExpr{
				pos: position{line: 342, col: 21, offset: 7917},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 342, col: 22, offset: 7918},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 342, col: 22, offset: 7918},
								val:        "//",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 342, col: 27, offset: 7923},
								expr: &seqExpr{
									pos: position{line: 342, col: 29, offset: 7925},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 342, col: 29, offset: 7925},
											expr: &ruleRefExpr{
												pos:  position{line: 342, col: 30, offset: 7926},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 342, col: 34, offset: 7930},
											name: "SourceChar",
										},
									},
								},
							},
						},
					},
					&seqExpr{
						pos: position{line: 342, col: 52, offset: 7948},
						exprs: []interface{}{
							&litMatcher{
								pos:        position{line: 342, col: 52, offset: 7948},
								val:        "#",
								ignoreCase: false,
							},
							&zeroOrMoreExpr{
								pos: position{line: 342, col: 56, offset: 7952},
								expr: &seqExpr{
									pos: position{line: 342, col: 58, offset: 7954},
									exprs: []interface{}{
										&notExpr{
											pos: position{line: 342, col: 58, offset: 7954},
											expr: &ruleRefExpr{
												pos:  position{line: 342, col: 59, offset: 7955},
												name: "EOL",
											},
										},
										&ruleRefExpr{
											pos:  position{line: 342, col: 63, offset: 7959},
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
			pos:  position{line: 344, col: 1, offset: 7975},
			expr: &zeroOrMoreExpr{
				pos: position{line: 344, col: 6, offset: 7982},
				expr: &choiceExpr{
					pos: position{line: 344, col: 8, offset: 7984},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 344, col: 8, offset: 7984},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 344, col: 21, offset: 7997},
							name: "EOL",
						},
						&ruleRefExpr{
							pos:  position{line: 344, col: 27, offset: 8003},
							name: "Comment",
						},
					},
				},
			},
		},
		{
			name: "_",
			pos:  position{line: 345, col: 1, offset: 8014},
			expr: &zeroOrMoreExpr{
				pos: position{line: 345, col: 5, offset: 8020},
				expr: &choiceExpr{
					pos: position{line: 345, col: 7, offset: 8022},
					alternatives: []interface{}{
						&ruleRefExpr{
							pos:  position{line: 345, col: 7, offset: 8022},
							name: "Whitespace",
						},
						&ruleRefExpr{
							pos:  position{line: 345, col: 20, offset: 8035},
							name: "MultiLineCommentNoLineTerminator",
						},
					},
				},
			},
		},
		{
			name: "WS",
			pos:  position{line: 346, col: 1, offset: 8071},
			expr: &zeroOrMoreExpr{
				pos: position{line: 346, col: 6, offset: 8078},
				expr: &ruleRefExpr{
					pos:  position{line: 346, col: 6, offset: 8078},
					name: "Whitespace",
				},
			},
		},
		{
			name: "Whitespace",
			pos:  position{line: 348, col: 1, offset: 8091},
			expr: &charClassMatcher{
				pos:        position{line: 348, col: 14, offset: 8106},
				val:        "[ \\t\\r]",
				chars:      []rune{' ', '\t', '\r'},
				ignoreCase: false,
				inverted:   false,
			},
		},
		{
			name: "EOL",
			pos:  position{line: 349, col: 1, offset: 8114},
			expr: &litMatcher{
				pos:        position{line: 349, col: 7, offset: 8122},
				val:        "\n",
				ignoreCase: false,
			},
		},
		{
			name: "EOS",
			pos:  position{line: 350, col: 1, offset: 8127},
			expr: &choiceExpr{
				pos: position{line: 350, col: 7, offset: 8135},
				alternatives: []interface{}{
					&seqExpr{
						pos: position{line: 350, col: 7, offset: 8135},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 350, col: 7, offset: 8135},
								name: "__",
							},
							&litMatcher{
								pos:        position{line: 350, col: 10, offset: 8138},
								val:        ";",
								ignoreCase: false,
							},
						},
					},
					&seqExpr{
						pos: position{line: 350, col: 16, offset: 8144},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 350, col: 16, offset: 8144},
								name: "_",
							},
							&zeroOrOneExpr{
								pos: position{line: 350, col: 18, offset: 8146},
								expr: &ruleRefExpr{
									pos:  position{line: 350, col: 18, offset: 8146},
									name: "SingleLineComment",
								},
							},
							&ruleRefExpr{
								pos:  position{line: 350, col: 37, offset: 8165},
								name: "EOL",
							},
						},
					},
					&seqExpr{
						pos: position{line: 350, col: 43, offset: 8171},
						exprs: []interface{}{
							&ruleRefExpr{
								pos:  position{line: 350, col: 43, offset: 8171},
								name: "__",
							},
							&ruleRefExpr{
								pos:  position{line: 350, col: 46, offset: 8174},
								name: "EOF",
							},
						},
					},
				},
			},
		},
		{
			name: "EOF",
			pos:  position{line: 352, col: 1, offset: 8179},
			expr: &notExpr{
				pos: position{line: 352, col: 7, offset: 8187},
				expr: &anyMatcher{
					line: 352, col: 8, offset: 8188,
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
			thrift.Namespaces[v.language] = v.namespace
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

func (c *current) onNamespace1(language, ns interface{}) (interface{}, error) {
	return &namespace{
		language:  ifaceSliceToString(language),
		namespace: ifaceSliceToString(ns),
	}, nil
}

func (p *parser) callonNamespace1() (interface{}, error) {
	stack := p.vstack[len(p.vstack)-1]
	_ = stack
	return p.cur.onNamespace1(stack["language"], stack["ns"])
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
