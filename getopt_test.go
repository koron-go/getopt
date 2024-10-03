package getopt_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/koron-go/getopt"
)

type Option struct {
	getopt.Option
	Err string
}

func checkGetopt(t *testing.T, opts string, args []string, want []getopt.Option, restWant []string) {
	var got []getopt.Option
	for opt, err := range getopt.Getopt(args, opts) {
		if err != nil {
			t.Fatalf("failed opt=%+v: %s", opt, err)
		}
		got = append(got, opt)

	}
	t.Helper()
	if d := cmp.Diff(want, got); d != "" {
		t.Errorf("mismatch: -want +got\n%s", d)
	}
	if d := cmp.Diff(restWant, getopt.RestArgs); d != "" {
		t.Errorf("reset mismatch: -want +got\n%s", d)
	}
}

func checkGetopt2(t *testing.T, opts string, args []string, want []Option, restWant []string) {
	var got []Option
	for opt, err := range getopt.Getopt(args, opts) {
		var errstr string
		if err != nil {
			errstr = err.Error()
		}
		got = append(got, Option{Option: opt, Err: errstr})
	}
	t.Helper()
	if d := cmp.Diff(want, got, cmpopts.EquateErrors()); d != "" {
		t.Errorf("mismatch: -want +got\n%s", d)
	}
	if d := cmp.Diff(restWant, getopt.RestArgs); d != "" {
		t.Errorf("reset mismatch: -want +got\n%s", d)
	}
}

func pstr(s string) *string {
	return &s
}

func TestSimple(t *testing.T) {
	checkGetopt(t, "h", []string{"-h"}, []getopt.Option{
		{Name: 'h', Arg: nil},
	}, []string{})
	checkGetopt(t, "hp", []string{"-h", "-p"}, []getopt.Option{
		{Name: 'h', Arg: nil},
		{Name: 'p', Arg: nil},
	}, []string{})
	checkGetopt(t, "hp", []string{"-hp"}, []getopt.Option{
		{Name: 'h', Arg: nil},
		{Name: 'p', Arg: nil},
	}, []string{})
	checkGetopt(t, "abc", []string{"-bac"}, []getopt.Option{
		{Name: 'b', Arg: nil},
		{Name: 'a', Arg: nil},
		{Name: 'c', Arg: nil},
	}, []string{})
}

func TestArgument(t *testing.T) {
	checkGetopt(t, "a:bc", []string{"-afoo"}, []getopt.Option{
		{Name: 'a', Arg: pstr("foo")},
	}, []string{})
	checkGetopt(t, "a:bc", []string{"-afoo", "-bc"}, []getopt.Option{
		{Name: 'a', Arg: pstr("foo")},
		{Name: 'b', Arg: nil},
		{Name: 'c', Arg: nil},
	}, []string{})
	checkGetopt(t, "a:bc", []string{"-a", "foo"}, []getopt.Option{
		{Name: 'a', Arg: pstr("foo")},
	}, []string{})
	checkGetopt(t, "a:bc", []string{"-a", "foo", "-bc"}, []getopt.Option{
		{Name: 'a', Arg: pstr("foo")},
		{Name: 'b', Arg: nil},
		{Name: 'c', Arg: nil},
	}, []string{})
	checkGetopt(t, "a:bc", []string{"-bcafoo"}, []getopt.Option{
		{Name: 'b', Arg: nil},
		{Name: 'c', Arg: nil},
		{Name: 'a', Arg: pstr("foo")},
	}, []string{})
}

func TestRest(t *testing.T) {
	checkGetopt(t, "abc", []string{"-a", "--", "foo", "bar"}, []getopt.Option{
		{Name: 'a', Arg: nil},
	}, []string{"foo", "bar"})
}

func TestNoOptions(t *testing.T) {
	checkGetopt(t, "abc", []string{"foo", "bar"}, nil, []string{"foo", "bar"})
	checkGetopt(t, "abc", []string{"-a", "foo", "bar"}, []getopt.Option{
		{Name: 'a', Arg: nil},
	}, []string{"foo", "bar"})
}

func TestNotOption(t *testing.T) {
	// Unknown option "-f"
	checkGetopt2(t, "abc", []string{"-f", "foo", "bar"}, []Option{
		{Option: getopt.Option{Name: 'f', Arg: nil}, Err: "illegal option: 'f'"},
	}, []string{"foo", "bar"})
	// Terminate with "-" (single hyphen)
	checkGetopt2(t, "a:", []string{"-afoo", "-"}, []Option{
		{Option: getopt.Option{Name: 'a', Arg: pstr("foo")}, Err: ""},
		{Option: getopt.Option{Name: 0, Arg: nil}, Err: "a single \"-\" is not supported"},
	}, []string{"-"})
	// Missing arguments.
	checkGetopt2(t, "a:", []string{"-a"}, []Option{
		{Option: getopt.Option{Name: 'a', Arg: nil}, Err: "no arguments supplied: 'a'"},
	}, nil)
}
