# USB Gadget Go Library

A Go library for creating and managing USB gadget configurations on Linux systems via ConfigFS.

## Features

- Define USB gadgets, configurations, and functions in Go  
- Built-in ECM and ACM functions, plus support for custom functions  
- Easy enable, disable, and teardown of gadgets  

## Installation

With Go modules enabled:

```bash
go install github.com/ardelean-calin/gadgetlib@latest
```

Or add to your project:

```bash
go get github.com/ardelean-calin/gadgetlib
```

## Quick Start

```go
package main

import (
  "log"

  "github.com/ardelean-calin/gadgetlib/functions"
  "github.com/ardelean-calin/gadgetlib/gadget"
)

func main() {
  opts := gadget.GadgetOptions{
    Name:         "g1",
    Manufacturer: "Acme Corp",
    Serial:       "ABC123",
    Controller:   "dummy_udc.0",
    Configs: []gadget.Config{
      {
        Number: 1,
        Functions: []gadget.Function{
          functions.FunctionECM{
            InstanceName: "net0",
            DevAddr:      "02:00:00:00:00:02",
            HostAddr:     "02:00:00:00:00:01",
          },
        },
      },
    },
  }

  g, err := gadget.New(opts)
  if err != nil {
    log.Fatalf("failed to create gadget: %v", err)
  }
  defer g.Teardown()

  if err := g.Enable(); err != nil {
    log.Fatalf("failed to enable gadget: %v", err)
  }

  // Gadget is now active...
}
```

## Testing

```bash
go test ./gadget ./functions
```

## Contributing

Contributions welcome! Please open issues or pull requests for bug fixes and enhancements.

## License

This project is licensed under the MIT License.
