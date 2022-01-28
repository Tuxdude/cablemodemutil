# cablemodemutil

[![PkgGoDev](https://pkg.go.dev/badge/github.com/tuxdude/cablemodemutil)](https://pkg.go.dev/github.com/tuxdude/cablemodemutil) [![Build](https://github.com/Tuxdude/cablemodemutil/actions/workflows/build.yml/badge.svg)](https://github.com/Tuxdude/cablemodemutil/actions/workflows/build.yml) [![Tests](https://github.com/Tuxdude/cablemodemutil/actions/workflows/tests.yml/badge.svg)](https://github.com/Tuxdude/cablemodemutil/actions/workflows/tests.yml) [![Lint](https://github.com/Tuxdude/cablemodemutil/actions/workflows/lint.yml/badge.svg)](https://github.com/Tuxdude/cablemodemutil/actions/workflows/lint.yml) [![CodeQL](https://github.com/Tuxdude/cablemodemutil/actions/workflows/codeql-analysis.yml/badge.svg)](https://github.com/Tuxdude/cablemodemutil/actions/workflows/codeql-analysis.yml) [![Go Report Card](https://goreportcard.com/badge/github.com/tuxdude/cablemodemutil)](https://goreportcard.com/report/github.com/tuxdude/cablemodemutil)

A go library for interfacing with Cable Modems.

This currently only works (and has been tested with) an Arris S33 Cable
Modem. If you would like to add support for other cable modems, please
file an Issue or submit a pull request with details for further discussion.

If you're looking for a command-line interface to use this library, please
see [`cablemodemcli`](https://github.com/Tuxdude/cablemodemcli).

# Example usage

```
package main

import (
  "fmt"
  "os"

  "github.com/tuxdude/cablemodemutil"
)

func main() {
  host := "192.168.100.1"
  protocol := "https"
  user := "admin"
  pass := "password"

  // Use the cable modem information to build the retriever.
  input := cablemodemutil.RetrieverInput{
    Host:           host,
    Protocol:       protocol,
    SkipVerifyCert: true,
    Username:       user,
    ClearPassword:  pass,
  }
  cm := cablemodemutil.NewStatusRetriever(&input)

  // This is a synchronous call to retrieve the status and takes
  // anywhere from two to ten seconds on average.
  st, err := cm.Status()
  if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %s\n", err)
    os.Exit(1)
  }

  // Access the status information.
  // More detailed fields are available in the status, please
  // refer to the documentation.
  fmt.Printf("Model: %s\n", st.Info.Model)
  fmt.Printf("Serial Number: %s\n", st.Info.SerialNumber)
  fmt.Printf("MAC Address: %s\n", st.Info.MACAddress)
  fmt.Printf("Connection Established Timestamp: %s\n", st.Connection.EstablishedAt)
  fmt.Printf("Firmaware version: %s\n", st.Software.FirmwareVersion)
  fmt.Printf("DOCSIS Version: %s\n", st.Software.DOCSISSpecVersion)
  fmt.Printf("Downstream Channel: %d Hz\n", st.Startup.Downstream.FrequencyHZ)

  os.Exit(0)
}
```
