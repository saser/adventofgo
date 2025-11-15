package day24

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hexops/autogold/v2"
	"go.saser.se/adventofgo/aocdata"
	"go.saser.se/adventofgo/aoctest"
)

func parseSystemT(t *testing.T, input string) *system {
	t.Helper()
	s, err := parseSystem(strings.TrimSpace(input))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

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
	test := func(v uint64) string {
		var bits []uint64
		for _, b := range bitsOf(v) {
			bits = append(bits, b)
		}
		return fmt.Sprintf("%b (MSB first) -> %v (LSB first)", v, bits)
	}
	var got string

	got = test(0)
	autogold.Expect("0 (MSB first) -> [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0] (LSB first)").Equal(t, got)
	got = test(^uint64(0))
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
	if got, want := s.z(), uint64(0b100); got != want {
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
	if got, want := s.z(), uint64(0b0011111101000); got != want {
		t.Errorf("s.z() = %b; want %b", got, want)
	}
}

func allSet(upTo int) uint64 {
	return ^uint64(0) >> (64 - (upTo + 1))
}

func TestAllSet(t *testing.T) {
	got := func(upTo int) string { return fmt.Sprintf("%b", allSet(upTo)) }
	autogold.Expect("1").Equal(t, got(0))
	autogold.Expect("11").Equal(t, got(1))
	autogold.Expect("111").Equal(t, got(2))
	autogold.Expect("1111").Equal(t, got(3))
	autogold.Expect("11111111111").Equal(t, got(10))
}

func TestSweep(t *testing.T) {
	s, err := parseSystem(strings.TrimSpace(aocdata.InputT(t, 2024, 24)))
	if err != nil {
		t.Fatal(err)
	}

	t.Log("before swapping:")
	for i := range 45 {
		var x, y uint64 = 1 << i, 0
		var want uint64 = x + y
		if got := s.run(x, y); got != want {
			t.Errorf("     444443333333333222222222211111111110000000000")
			t.Errorf("     432109876543210987654321098765432109876543210")
			t.Errorf("     %045b", x)
			t.Errorf("   + %045b", y)
			t.Errorf("   = %045b;", got)
			t.Errorf("want %045b", want)
			t.Error()
		}
	}

	s.swap("z09", "z10")
	s.swap("z20", "z21")
	s.swap("z30", "z31")
	s.swap("z34", "z35")

	t.Log("after swapping:")
	for i := range 45 {
		var x, y uint64 = 1 << i, 0
		var want uint64 = x + y
		if got := s.run(x, y); got != want {
			t.Errorf("     444443333333333222222222211111111110000000000")
			t.Errorf("     432109876543210987654321098765432109876543210")
			t.Errorf("     %045b", x)
			t.Errorf("   + %045b", y)
			t.Errorf("   = %045b;", got)
			t.Errorf("want %045b", want)
			t.Error()
		}
	}
}

func TestSwap(t *testing.T) {
	s := parseSystemT(t, `
x00: 0
x01: 1
x02: 0
x03: 1
x04: 0
x05: 1
y00: 0
y01: 0
y02: 1
y03: 1
y04: 0
y05: 1

x00 AND y00 -> z05
x01 AND y01 -> z02
x02 AND y02 -> z01
x03 AND y03 -> z03
x04 AND y04 -> z04
x05 AND y05 -> z00
	`)

	var x, y uint64 = 0b101010, 0b101100
	var want uint64 = x & y
	t.Errorf("before swapping:")
	got := s.run(x, y)
	t.Errorf("     0000000000")
	t.Errorf("     9876543210")
	t.Errorf("     %010b", x)
	t.Errorf("   & %010b", y)
	t.Errorf("   = %010b;", got)
	t.Errorf("want %010b", want)
	t.Error()

	s.swap("z00", "z05")
	s.swap("z01", "z02")

	t.Errorf("after swapping:")
	got = s.run(x, y)
	t.Errorf("     0000000000")
	t.Errorf("     9876543210")
	t.Errorf("     %010b", x)
	t.Errorf("   & %010b", y)
	t.Errorf("   = %010b;", got)
	t.Errorf("want %010b", want)
	t.Error()
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
