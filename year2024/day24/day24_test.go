package day24

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hexops/autogold/v2"
	"go.saser.se/adventofgo/aoctest"
)

func TestParseWires(t *testing.T) {
	test := func(line string) string {
		w, err := parseWire(line)
		if err != nil {
			return fmt.Sprintf("%q -> error: %v", line, err)
		}
		return fmt.Sprintf("%q -> %+v", line, w)
	}
	var sb strings.Builder
	sb.WriteByte('\n')
	for _, line := range []string{
		"x00: 1",
		"x01: 1",
		"x02: 1",
		"y00: 0",
		"y01: 1",
		"y02: 0",
	} {
		fmt.Fprintln(&sb, test(line))
	}
	got := sb.String()
	autogold.Expect(`
"x00: 1" -> x00(1)
"x01: 1" -> x01(1)
"x02: 1" -> x02(1)
"y00: 0" -> y00(0)
"y01: 1" -> y01(1)
"y02: 0" -> y02(0)
`).Equal(t, got)
}

func TestParseGate(t *testing.T) {
	test := func(lines ...string) string {
		var sb strings.Builder
		sb.WriteByte('\n')
		for _, line := range lines {
			g, err := parseGateDesc(line)
			if err != nil {
				fmt.Fprintf(&sb, "%q -> error: %v\n", line, err)
			} else {
				fmt.Fprintf(&sb, "%q -> %+v\n", line, g)
			}
		}
		return sb.String()
	}
	got := test(
		"x00 AND y00 -> z00",
		"x01 XOR y01 -> z01",
		"x02 OR y02 -> z02",
	)
	autogold.Expect(`
"x00 AND y00 -> z00" -> {W1:x00 Op:AND W2:y00 Out:z00}
"x01 XOR y01 -> z01" -> {W1:x01 Op:XOR W2:y01 Out:z01}
"x02 OR y02 -> z02" -> {W1:x02 Op:OR W2:y02 Out:z02}
`).Equal(t, got)
	got = test(
		"ntg XOR fgs -> mjb",
		"y02 OR x01 -> tnw",
		"kwq OR kpj -> z05",
		"x00 OR x03 -> fst",
		"tgd XOR rvg -> z01",
		"vdt OR tnw -> bfw",
		"bfw AND frj -> z10",
	)
	autogold.Expect(`
"ntg XOR fgs -> mjb" -> {W1:ntg Op:XOR W2:fgs Out:mjb}
"y02 OR x01 -> tnw" -> {W1:y02 Op:OR W2:x01 Out:tnw}
"kwq OR kpj -> z05" -> {W1:kwq Op:OR W2:kpj Out:z05}
"x00 OR x03 -> fst" -> {W1:x00 Op:OR W2:x03 Out:fst}
"tgd XOR rvg -> z01" -> {W1:tgd Op:XOR W2:rvg Out:z01}
"vdt OR tnw -> bfw" -> {W1:vdt Op:OR W2:tnw Out:bfw}
"bfw AND frj -> z10" -> {W1:bfw Op:AND W2:frj Out:z10}
`).Equal(t, got)
}

func TestEvaluate(t *testing.T) {
	s, err := parseSystem(strings.TrimSpace(`
x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0

x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02
	`))
	if err != nil {
		t.Fatalf("parse system: %v", err)
	}
	z00 := s.evaluate(s.wires["z00"])
	z01 := s.evaluate(s.wires["z01"])
	z02 := s.evaluate(s.wires["z02"])
	t.Logf("z00: %+v\n", z00)
	t.Logf("z01: %+v\n", z01)
	t.Logf("z02: %+v\n", z02)
}

func TestBitsOf(t *testing.T) {
	test := func(v int64) string {
		var bits []int64
		for _, b := range bitsOf(v) {
			bits = append(bits, b)
		}
		return fmt.Sprintf("%b (MSB first) -> %v (LSB first)", v, bits)
	}
	var got string

	got = test(0)
	autogold.Expect("0 (MSB first) -> [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] (LSB first)").Equal(t, got)
	got = test(^0)
	autogold.Expect("-1 (MSB first) -> [1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1] (LSB first)").Equal(t, got)
	got = test(123456)
	autogold.Expect("11110001001000000 (MSB first) -> [0 0 0 0 0 0 1 0 0 1 0 0 0 1 1 1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] (LSB first)").Equal(t, got)
}

func TestRun(t *testing.T) {
	s, err := parseSystem(strings.TrimSpace(`
x00: 1
x01: 1
x02: 1
y00: 0
y01: 1
y02: 0

x00 AND y00 -> z00
x01 XOR y01 -> z01
x02 OR y02 -> z02
`))
	if err != nil {
		t.Fatal(err)
	}
	s.evaluateAll()
	if got, want := s.z(), int64(0b100); got != want {
		t.Errorf("s.z() = %b; want %b", got, want)
	}

	s, err = parseSystem(strings.TrimSpace(`
x00: 1
x01: 0
x02: 1
x03: 1
x04: 0
y00: 1
y01: 1
y02: 1
y03: 1
y04: 1

ntg XOR fgs -> mjb
y02 OR x01 -> tnw
kwq OR kpj -> z05
x00 OR x03 -> fst
tgd XOR rvg -> z01
vdt OR tnw -> bfw
bfw AND frj -> z10
ffh OR nrd -> bqk
y00 AND y03 -> djm
y03 OR y00 -> psh
bqk OR frj -> z08
tnw OR fst -> frj
gnj AND tgd -> z11
bfw XOR mjb -> z00
x03 OR x00 -> vdt
gnj AND wpb -> z02
x04 AND y00 -> kjc
djm OR pbm -> qhw
nrd AND vdt -> hwm
kjc AND fst -> rvg
y04 OR y02 -> fgs
y01 AND x02 -> pbm
ntg OR kjc -> kwq
psh XOR fgs -> tgd
qhw XOR tgd -> z09
pbm OR djm -> kpj
x03 XOR y03 -> ffh
x00 XOR y04 -> ntg
bfw OR bqk -> z06
nrd XOR fgs -> wpb
frj XOR qhw -> z04
bqk OR frj -> z07
y03 OR x01 -> nrd
hwm AND bqk -> z03
tgd XOR rvg -> z12
tnw OR pbm -> gnj
	`))
	if err != nil {
		t.Fatal(err)
	}
	s.evaluateAll()
	if got, want := s.z(), int64(0b0011111101000); got != want {
		t.Errorf("s.z() = %b; want %b", got, want)
	}
}

func TestPart1(t *testing.T) {
	aoctest.Test(t, 2024, 24, 1, Part1)
}

func TestPart2(t *testing.T) {
	aoctest.Test(t, 2024, 24, 2, Part2)
}

func BenchmarkPart1(b *testing.B) {
	aoctest.Benchmark(b, 2024, 24, 1, Part1)
}

func BenchmarkPart2(b *testing.B) {
	aoctest.Benchmark(b, 2024, 24, 2, Part2)
}
