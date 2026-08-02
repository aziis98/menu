package main

import "github.com/spf13/pflag"

var (
	initial     = pflag.StringP("initial", "i", "", "Initial prompt text")
	placeholder = pflag.StringP("placeholder", "p", "", "Placeholder text")
	selection   = pflag.BoolP("selection", "s", false, "Only return the selected item")
)
