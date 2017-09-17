package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"testing"
)

func TestKubectlWrapper(t *testing.T) {

	// Only pass t into top-level Convey calls
	Convey("Given a yaml formatted kubernetes object", t, func() {
		srcYaml := string(mustReadFile("testdata/mergeprunelabel_input.yaml"))

		Convey("When calling mergePrunelabel()", func() {
			mergedYaml, err := mergePrunelabel(srcYaml, "aos9mra4rdf")
			if err != nil {
				t.Fatal(err)
			}
			Convey("Should add the prunelabel, leaving everything else untouched", func() {
				expectedYaml := string(mustReadFile("testdata/mergeprunelabel_expected.yaml"))
				So(mergedYaml, ShouldEqual, expectedYaml)
			})
		})
	})
}

func mustReadFile(path string) []byte {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return content
}
