package main

import (
	"github.com/khulnasoft-lab/vulmap/pkg/testutils"
)

var dnsTestCases = []TestCaseInfo{
	{Path: "protocols/dns/basic.yaml", TestCase: &dnsBasic{}},
	{Path: "protocols/dns/ptr.yaml", TestCase: &dnsPtr{}},
	{Path: "protocols/dns/caa.yaml", TestCase: &dnsCAA{}},
	{Path: "protocols/dns/tlsa.yaml", TestCase: &dnsTLSA{}},
	{Path: "protocols/dns/variables.yaml", TestCase: &dnsVariables{}},
	{Path: "protocols/dns/payload.yaml", TestCase: &dnsPayload{}},
	{Path: "protocols/dns/dsl-matcher-variable.yaml", TestCase: &dnsDSLMatcherVariable{}},
}

type dnsBasic struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsBasic) Execute(filePath string) error {
	results, err := testutils.RunVulmapTemplateAndGetResults(filePath, "one.one.one.one", debug)
	if err != nil {
		return err
	}
	return expectResultsCount(results, 1)
}

type dnsPtr struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsPtr) Execute(filePath string) error {
	results, err := testutils.RunVulmapTemplateAndGetResults(filePath, "1.1.1.1", debug)
	if err != nil {
		return err
	}
	return expectResultsCount(results, 1)
}

type dnsCAA struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsCAA) Execute(filePath string) error {
	results, err := testutils.RunVulmapTemplateAndGetResults(filePath, "google.com", debug)
	if err != nil {
		return err
	}
	return expectResultsCount(results, 1)
}

type dnsTLSA struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsTLSA) Execute(filePath string) error {
	results, err := testutils.RunVulmapTemplateAndGetResults(filePath, "scanme.sh", debug)
	if err != nil {
		return err
	}
	return expectResultsCount(results, 0)
}

type dnsVariables struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsVariables) Execute(filePath string) error {
	results, err := testutils.RunVulmapTemplateAndGetResults(filePath, "one.one.one.one", debug)
	if err != nil {
		return err
	}
	return expectResultsCount(results, 1)
}

type dnsPayload struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsPayload) Execute(filePath string) error {
	results, err := testutils.RunVulmapTemplateAndGetResults(filePath, "google.com", debug)
	if err != nil {
		return err
	}
	if err := expectResultsCount(results, 3); err != nil {
		return err
	}

	// override payload from CLI
	results, err = testutils.RunVulmapTemplateAndGetResults(filePath, "google.com", debug, "-var", "subdomain_wordlist=subdomains.txt")
	if err != nil {
		return err
	}
	return expectResultsCount(results, 4)
}

type dnsDSLMatcherVariable struct{}

// Execute executes a test case and returns an error if occurred
func (h *dnsDSLMatcherVariable) Execute(filePath string) error {
	results, err := testutils.RunVulmapTemplateAndGetResults(filePath, "one.one.one.one", debug)
	if err != nil {
		return err
	}
	return expectResultsCount(results, 1)
}
