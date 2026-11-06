// Command goversion is a CLI for parsing, comparing, and checking
// software versions using github.com/khulnasoft/go-version-lib.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/khulnasoft/go-version-lib/pkg/reference"
	"github.com/khulnasoft/go-version-lib/pkg/version"
)

// stdout and stderr are variables so tests can capture output.
var (
	stdout io.Writer = os.Stdout
	stderr io.Writer = os.Stderr
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) < 1 {
		printUsage()
		return 2
	}

	switch args[0] {
	case "parse":
		return runParse(args[1:])
	case "compare":
		return runCompare(args[1:])
	case "check":
		return runCheck(args[1:])
	case "detect":
		return runDetect(args[1:])
	case "help", "-h", "--help":
		printUsage()
		return 0
	default:
		fmt.Fprintf(stderr, "unknown command: %s\n\n", args[0])
		printUsage()
		return 2
	}
}

func printUsage() {
	fmt.Fprintln(stdout, "Usage: goversion <command> [arguments]")
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Commands:")
	fmt.Fprintln(stdout, "  parse <version>    Parse and print a version")
	fmt.Fprintln(stdout, "  compare <v1> <v2>  Compare two versions")
	fmt.Fprintln(stdout, "  check <v> <c>      Check if a version satisfies a constraint")
	fmt.Fprintln(stdout, "  detect <version>   Detect the version type of a version string")
	fmt.Fprintln(stdout)
	fmt.Fprintln(stdout, "Examples:")
	fmt.Fprintln(stdout, "  goversion parse 1.2.3")
	fmt.Fprintln(stdout, "  goversion compare 1.2.3 1.2.4")
	fmt.Fprintln(stdout, "  goversion check 1.2.3 '>=1.2.0'")
	fmt.Fprintln(stdout, "  goversion detect 1.2.3")
}

func runParse(args []string) int {
	fs := flag.NewFlagSet("parse", flag.ContinueOnError)
	fs.Parse(args)

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: goversion parse <version>")
		return 2
	}

	ref, err := reference.NewReference(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %q: %v\n", fs.Arg(0), err)
		return 1
	}

	fmt.Fprintln(stdout, ref.String())
	return 0
}

func runCompare(args []string) int {
	fs := flag.NewFlagSet("compare", flag.ContinueOnError)
	fs.Parse(args)

	if fs.NArg() != 2 {
		fmt.Fprintln(stderr, "usage: goversion compare <v1> <v2>")
		return 2
	}

	r1, err := reference.NewReference(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %q: %v\n", fs.Arg(0), err)
		return 1
	}

	r2, err := reference.NewReference(fs.Arg(1))
	if err != nil {
		fmt.Fprintf(stderr, "error parsing %q: %v\n", fs.Arg(1), err)
		return 1
	}

	switch result := r1.Compare(r2); {
	case result < 0:
		fmt.Fprintf(stdout, "%s < %s\n", fs.Arg(0), fs.Arg(1))
	case result > 0:
		fmt.Fprintf(stdout, "%s > %s\n", fs.Arg(0), fs.Arg(1))
	default:
		fmt.Fprintf(stdout, "%s == %s\n", fs.Arg(0), fs.Arg(1))
	}
	return 0
}

func runCheck(args []string) int {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	fs.Parse(args)

	if fs.NArg() != 2 {
		fmt.Fprintln(stderr, "usage: goversion check <version> <constraint>")
		return 2
	}

	v, c := fs.Arg(0), fs.Arg(1)

	ver, err := version.Parse(v)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing version %q: %v\n", v, err)
		return 1
	}

	constraints, err := version.NewConstraints(c)
	if err != nil {
		fmt.Fprintf(stderr, "error parsing constraint %q: %v\n", c, err)
		return 1
	}

	if constraints.Check(ver) {
		fmt.Fprintf(stdout, "%s satisfies %s\n", v, c)
		return 0
	}

	fmt.Fprintf(stdout, "%s does not satisfy %s\n", v, c)
	return 1
}

func runDetect(args []string) int {
	fs := flag.NewFlagSet("detect", flag.ContinueOnError)
	fs.Parse(args)

	if fs.NArg() != 1 {
		fmt.Fprintln(stderr, "usage: goversion detect <version>")
		return 2
	}

	typ, err := reference.Detect(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(stderr, "error: %v\n", err)
		return 1
	}

	fmt.Fprintln(stdout, typ)
	return 0
}
