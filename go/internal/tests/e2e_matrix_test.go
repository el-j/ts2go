package tests

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/el-j/ts2go/internal/transpiler"
)

type e2eTestCase struct {
	name   string
	tsCode string
}

func TestE2EMatrix(t *testing.T) {
	nodePath, err := exec.LookPath("node")
	if err != nil {
		t.Skip("node is not installed, skipping Node-comparison E2E matrix tests")
	}

	testCases := []e2eTestCase{
		{
			name: "arithmetic_and_logical_operators",
			tsCode: `function compute(a: number, b: number): number {
  const sum = a + b;
  const product = a * b;
  const diff = b - a;
  return (sum > 0 && product > 0) ? (diff + 10) : 0;
}
console.log(compute(3, 7));
`,
		},
		{
			name: "loops_and_accumulation",
			tsCode: `function sumNumbers(): number {
  let total = 0;
  for (let i = 0; i < 5; i++) {
    total += i;
  }
  return total;
}
console.log(sumNumbers());
`,
		},
		{
			name: "switch_control_flow",
			tsCode: `function checkVal(x: number): string {
  switch (x) {
    case 1:
      return "one";
    case 2:
      return "two";
    default:
      return "other";
  }
}
console.log(checkVal(2));
console.log(checkVal(99));
`,
		},
		{
			name: "classes_and_constructors",
			tsCode: `class Calculator {
  base: number;
  constructor(initial: number) {
    this.base = initial;
  }
  add(n: number): number {
    return this.base + n;
  }
}
const calc = new Calculator(100);
console.log(calc.add(25));
`,
		},
		{
			name: "arrow_functions_and_closures",
			tsCode: `function execute(n: number): number {
  const multiplier = (x: number): number => x * 3;
  return multiplier(n);
}
console.log(execute(14));
`,
		},
		{
			name: "interfaces_and_struct_instantiation",
			tsCode: `interface Point {
  x: number;
  y: number;
}
function makePoint(px: number, py: number): Point {
  return { x: px, y: py };
}
const pt = makePoint(10, 20);
console.log(pt.x + pt.y);
`,
		},
		{
			name: "nullish_coalescing",
			tsCode: `function fallback(val: string | null): string {
  const result: string = val ?? "default-val";
  return result;
}
console.log(fallback(null));
console.log(fallback("specified"));
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tsFile := filepath.Join(tmpDir, "input.ts")
			goFile := filepath.Join(tmpDir, "main.go")
			binFile := filepath.Join(tmpDir, "program")

			if err := os.WriteFile(tsFile, []byte(tc.tsCode), 0644); err != nil {
				t.Fatalf("failed to write ts file: %v", err)
			}

			// Step 1: Run with Node.js to get ground-truth output
			nodeCmd := exec.Command(nodePath, "--experimental-strip-types", tsFile)
			var nodeOut bytes.Buffer
			var nodeErr bytes.Buffer
			nodeCmd.Stdout = &nodeOut
			nodeCmd.Stderr = &nodeErr
			if err := nodeCmd.Run(); err != nil {
				t.Fatalf("Node execution failed: %v\nStderr: %s", err, nodeErr.String())
			}
			expectedOutput := strings.TrimSpace(nodeOut.String())

			// Step 2: Transpile TS to Go
			if err := transpiler.Transpile(tsFile, goFile); err != nil {
				t.Fatalf("Transpile failed: %v", err)
			}

			// Step 3: Build Go binary
			buildCmd := exec.Command("go", "build", "-o", binFile, goFile)
			var buildOut bytes.Buffer
			buildCmd.Stdout = &buildOut
			buildCmd.Stderr = &buildOut
			if err := buildCmd.Run(); err != nil {
				content, _ := os.ReadFile(goFile)
				t.Fatalf("Go build failed: %v\nOutput: %s\nGenerated Go Code:\n%s", err, buildOut.String(), string(content))
			}

			// Step 4: Execute Go binary
			runCmd := exec.Command(binFile)
			var goOut bytes.Buffer
			var goErr bytes.Buffer
			runCmd.Stdout = &goOut
			runCmd.Stderr = &goErr
			if err := runCmd.Run(); err != nil {
				t.Fatalf("Go binary execution failed: %v\nStderr: %s", err, goErr.String())
			}
			actualOutput := strings.TrimSpace(goOut.String())

			// Step 5: Verify identical outputs
			if actualOutput != expectedOutput {
				t.Errorf("Output mismatch:\nExpected (Node):\n%s\nGot (Go):\n%s", expectedOutput, actualOutput)
			}
		})
	}
}
