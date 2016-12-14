package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

// Line describes if a line is covered
type Line struct {
	Number int   `xml:"number,attr"`
	Hits   int64 `xml:"hits,attr"`
}

// Method combines lines of a given method
type Method struct {
	Name       string  `xml:"name,attr"`
	Signature  string  `xml:"signature,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Lines      []Line  `xml:"lines>line"`
}

// Class covers a given class
type Class struct {
	Name       string   `xml:"name,attr"`
	Filename   string   `xml:"filename,attr"`
	LineRate   float32  `xml:"line-rate,attr"`
	BranchRate float32  `xml:"branch-rate,attr"`
	Complexity float32  `xml:"complexity,attr"`
	Methods    []Method `xml:"methods>method"`
	Lines      []Line   `xml:"lines>line"`
}

// Package holds classes
type Package struct {
	Name       string  `xml:"name,attr"`
	LineRate   float32 `xml:"line-rate,attr"`
	BranchRate float32 `xml:"branch-rate,attr"`
	Complexity float32 `xml:"complexity,attr"`
	Classes    []Class `xml:"classes>class"`
}

// Coverage is the highest level object
type Coverage struct {
	XMLName    xml.Name  `xml:"coverage"`
	LineRate   float32   `xml:"line-rate,attr"`
	BranchRate float32   `xml:"branch-rate,attr"`
	Version    string    `xml:"version,attr"`
	Timestamp  int64     `xml:"timestamp,attr"`
	Packages   []Package `xml:"packages>package"`
}

// GetLines transforms each line into a string
func GetLines(c Coverage) []string {
	var ret []string
	for _, pkg := range c.Packages {
		for _, cl := range pkg.Classes {
			for _, l := range cl.Lines {
				ret = append(ret, fmt.Sprintf("%s:%d.0,%d.50 1 %d", cl.Filename, l.Number, l.Number, l.Hits))
			}
		}
	}
	return ret
}

func main() {
	fileName := flag.String("filename", "coverage.xml", "Coverage xml report")
	flag.Parse()
	c := Coverage{}
	data, err := ioutil.ReadFile(*fileName)
	if err != nil {
		log.Fatal(err)
	}
	err = xml.Unmarshal([]byte(data), &c)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	for _, l := range GetLines(c) {
		fmt.Println(l)
	}
}
