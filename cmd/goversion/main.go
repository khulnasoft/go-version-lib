// Command goversion is a CLI for parsing, comparing, and checking
// software versions using github.com/khulnasoft/go-version-lib.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/khulnasoft/go-version-lib/pkg/reference"
	"github.com/khulnasoft/go-version-lib/pkg/version"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(2)
	}

	switch os.Args[1] {
	case "parse":
		cmdParse(os.Args[2:])
	case "compare":
		cmdCompare(os.Args[2:])
	case "check":
		cmdCheck(os.Args[2:])
	case "detect":
		cmdDetect(os.Args[2:])
	case "help", "-h", "--help":
		printUsage()
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n\n", os.Args[1])
		printUsage()
		os.Exit(2)
	}
}

func printUsage() {
	fmt.Println("Usage: goversion <command> [arguments]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  parse <version>    Parse and print a version")
	fmt.Println("  compare <v1> <v2>  Compare two versions")
	fmt.Println("  check <v> <c>      Check if a version satisfies a constraint")
	fmt.Println("  detect <version>   Detect the version type of a version string")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  goversion parse 1.2.3")
	fmt.Println("  goversion compare 1.2.3 1.2.4")
	fmt.Println("  goversion check 1.2.3 '>=1.2.0'")
	fmt.Println("  goversion detect 1.2.3")
}

func cmdParse(args []string) {
	fs := flag.NewFlagSet("parse", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: goversion parse <version>")
		os.Exit(2)
	}

	ref, err := reference.NewReference(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %q: %v\n", fs.Arg(0), err)
		os.Exit(1)
	}

	fmt.Println(ref.String())
}

func cmdCompare(args []string) {
	fs := flag.NewFlagSet("compare", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: goversion compare <v1> <v2>")
		os.Exit(2)
	}

	r1, err := reference.NewReference(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %q: %v\n", fs.Arg(0), err)
		os.Exit(1)
	}

	r2, err := reference.NewReference(fs.Arg(1))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %q: %v\n", fs.Arg(1), err)
		os.Exit(1)
	}

	switch result := r1.Compare(r2); {
	case result < 0:
		fmt.Printf("%s < %s\n", fs.Arg(0), fs.Arg(1))
	case result > 0:
		fmt.Printf("%s > %s\n", fs.Arg(0), fs.Arg(1))
	default:
		fmt.Printf("%s == %s\n", fs.Arg(0), fs.Arg(1))
	}
}

func cmdCheck(args []string) {
	fs := flag.NewFlagSet("check", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() != 2 {
		fmt.Fprintln(os.Stderr, "usage: goversion check <version> <constraint>")
		os.Exit(2)
	}

	v, c := fs.Arg(0), fs.Arg(1)

	ver, err := version.Parse(v)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing version %q: %v\n", v, err)
		os.Exit(1)
	}

	constraints, err := version.NewConstraints(c)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error parsing constraint %q: %v\n", c, err)
		os.Exit(1)
	}

	if constraints.Check(ver) {
		fmt.Printf("%s satisfies %s\n", v, c)
		return
	}

	fmt.Printf("%s does not satisfy %s\n", v, c)
	os.Exit(1)
}

func cmdDetect(args []string) {
	fs := flag.NewFlagSet("detect", flag.ExitOnError)
	fs.Parse(args)

	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: goversion detect <version>")
		os.Exit(2)
	}

	typ, err := reference.Detect(fs.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(typ)
}
