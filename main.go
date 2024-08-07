package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	semver "github.com/Masterminds/semver/v3"
	kingpin "github.com/alecthomas/kingpin/v2"
)

var version = "1.0.0"

var (
	app     = kingpin.New("semver", "Command-line semver tools. On error, print to stderr and exit -1.")
	verbose = app.Flag("verbose", "Verbose mode.").Short('v').Bool()

	satisfies            = app.Command("satisfies", "Test if a version satisfies a constraint. Exit 0 if satisfies, 1 if not. If verbose, print an explanation to stdout.")
	satisfiesVersion     = satisfies.Arg("VERSION", "The version to test").Required().String()
	satisfiesConstraints = satisfies.Arg("CONSTRAINTS", "The constraints to test against").Required().String()

	greater  = app.Command("greater", "Compare two versions. Exit 0 if the first is greater, 1 if not. If verbose, print greater to stdout.")
	greaterA = greater.Arg("A", "Left side of A > B").Required().String()
	greaterB = greater.Arg("B", "Right side of A > B").Required().String()

	lesser  = app.Command("lesser", "Compare two versions. Exit 0 if the first is lesser, 1 if not. If verbose, print lesser to stdout.")
	lesserA = lesser.Arg("A", "Left side of A < B").Required().String()
	lesserB = lesser.Arg("B", "Right side of A < B").Required().String()

	equal  = app.Command("equal", "Compare two versions. Exit 0 if they are equal, 1 if not.")
	equalA = equal.Arg("A", "Left side of A = B").Required().String()
	equalB = equal.Arg("B", "Right side of A = B").Required().String()

	inc          = app.Command("inc", "Increment major, minor, or patch component.")
	incComponent = inc.Arg("COMPONENT", "The component to increment. Possible values: [major, minor, patch]").Required().String()
	incVersion   = inc.Arg("VERSION", "The version to increment.").Required().String()

	get          = app.Command("get", "Get major, minor, patch, prerelease or metadata component.")
	getComponent = get.Arg("COMPONENT", "The component to increment. Possible values: [major, minor, patch, prerelease, metadata]").Required().String()
	getVersion   = get.Arg("VERSION", "The version to retreive component from.").Required().String()

	set          = app.Command("set", "Set prerelease or metadata component.")
	setComponent = set.Arg("COMPONENT", "The component to increment. Possible values: [prerelease, metadata]").Required().String()
	setVersion   = set.Arg("VERSION", "The version of which to set a component.").Required().String()
	setValue     = set.Arg("VALUE", "The value to set.").Required().String()

	greatest           = app.Command("greatest", "Find the greatest version in a list.")
	filter_pre_release = greatest.Flag("filte-pre-release", "Ignores all versions with pre-release information before comparison").Short('p').Bool()
	filter_build       = greatest.Flag("filte-build", "Ignores all versions with build information before comparison").Short('b').Bool()
	versions           = greatest.Arg("VERSIONS", "The versions to compare.").Required().Strings()
)

func main() {
	kingpin.Version(version)

	switch kingpin.MustParse(app.Parse(os.Args[1:])) {
	case satisfies.FullCommand():
		v := mustParseVersion(*satisfiesVersion, "VERSION")
		c := mustParseConstraints(*satisfiesConstraints)

		if does, msgs := c.Validate(v); !does {
			if *verbose {
				for _, m := range msgs {
					fmt.Println(m)
				}
			}

			os.Exit(1)
		}

		os.Exit(0)

	case greater.FullCommand():
		a := mustParseVersion(*greaterA, "A")
		b := mustParseVersion(*greaterB, "B")

		if !a.GreaterThan(b) {
			if *verbose {
				fmt.Println(*greaterB)
			}
			os.Exit(1)
		}

		if *verbose {
			fmt.Println(*greaterA)
		}
		os.Exit(0)

	case lesser.FullCommand():
		a := mustParseVersion(*lesserA, "A")
		b := mustParseVersion(*lesserB, "B")

		if !a.LessThan(b) {
			if *verbose {
				fmt.Println(*lesserB)
			}
			os.Exit(1)
		}

		if *verbose {
			fmt.Println(*lesserA)
		}
		os.Exit(0)

	case equal.FullCommand():
		a := mustParseVersion(*equalA, "A")
		b := mustParseVersion(*equalB, "B")

		if !a.Equal(b) {
			os.Exit(1)
		}

		os.Exit(0)

	case inc.FullCommand():
		v := mustParseVersion(*incVersion, "VERSION")
		var v1 semver.Version
		switch *incComponent {
		case "major":
			v1 = v.IncMajor()
		case "minor":
			v1 = v.IncMinor()
		case "patch":
			v1 = v.IncPatch()
		default:
			fmt.Fprintf(os.Stderr, "unknown component name: '%s'\n", *incComponent)
			os.Exit(-1)
		}
		fmt.Println(v1.String())

	case get.FullCommand():
		v := mustParseVersion(*getVersion, "VERSION")
		var component string
		switch *getComponent {
		case "major":
			component = strconv.FormatUint(v.Major(), 10)
		case "minor":
			component = strconv.FormatUint(v.Minor(), 10)
		case "patch":
			component = strconv.FormatUint(v.Patch(), 10)
		case "prerelease":
			component = v.Prerelease()
		case "metadata":
			component = v.Metadata()
		default:
			fmt.Fprintf(os.Stderr, "unknown component name: '%s'\n", *getComponent)
			os.Exit(-1)
		}
		fmt.Println(component)

	case set.FullCommand():
		v := mustParseVersion(*setVersion, "VERSION")
		var v1 semver.Version
		var err error
		switch *setComponent {
		case "prerelease":
			if v1, err = v.SetPrerelease(*setValue); err != nil {
				fmt.Fprintf(os.Stderr, "invalid prerelease; %v\n", err)
				os.Exit(-1)
			}
		case "metadata":
			if v1, err = v.SetMetadata(*setValue); err != nil {
				fmt.Fprintf(os.Stderr, "invalid metadata; %v\n", err)
				os.Exit(-1)
			}
		default:
			fmt.Fprintf(os.Stderr, "unknown component name: '%s'\n", *setComponent)
			os.Exit(-1)
		}
		fmt.Println(v1.String())

	case greatest.FullCommand():
		all_parsed_versions := []semver.Version{}
		for _, v := range *versions {
			all_parsed_versions = append(all_parsed_versions, *mustParseVersion(v, "VERSION"))
		}

		filtered_versions := all_parsed_versions

		if *filter_pre_release {
			filtered_pre_release := []semver.Version{}
			for _, v := range all_parsed_versions {
				if v.Prerelease() == "" {
					filtered_pre_release = append(filtered_pre_release, v)
				}
			}
			filtered_versions = filtered_pre_release
		}

		if *filter_build {
			filtered_build := []semver.Version{}
			for _, v := range filtered_versions {
				if v.Metadata() == "" {
					filtered_build = append(filtered_build, v)
				}
			}
			filtered_versions = filtered_build
		}

		sort.Slice(filtered_versions, func(i, j int) bool {
			return filtered_versions[i].LessThan(&filtered_versions[j])
		})

		fmt.Println(filtered_versions[len(filtered_versions)-1].String())
	}
}

func mustParseVersion(s, ctx string) *semver.Version {
	v, err := semver.NewVersion(s)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse <%s> version; %v: '%s'\n", ctx, err, s)
		os.Exit(-1)
	}

	return v
}

func mustParseConstraints(s string) *semver.Constraints {
	c, err := semver.NewConstraint(s)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse constraints; %v\n", err)
		os.Exit(-1)
	}

	return c
}
