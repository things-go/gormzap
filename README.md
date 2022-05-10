# zap logging driver for gorm v2
zap logging driver for gorm v2

[![GoDoc](https://godoc.org/github.com/things-go/gormzap?status.svg)](https://godoc.org/github.com/things-go/gormzap)
[![Go.Dev reference](https://img.shields.io/badge/go.dev-reference-blue?logo=go&logoColor=white)](https://pkg.go.dev/github.com/things-go/gormzap?tab=doc)
[![codecov](https://codecov.io/gh/things-go/gormzap/branch/main/graph/badge.svg)](https://codecov.io/gh/things-go/gormzap)
![Action Status](https://github.com/things-go/gormzap/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/things-go/gormzap)](https://goreportcard.com/report/github.com/things-go/gormzap)
[![Licence](https://img.shields.io/github/license/things-go/gormzap)](https://raw.githubusercontent.com/things-go/gormzap/main/LICENSE)
[![Tag](https://img.shields.io/github/v/tag/things-go/gormzap)](https://github.com/things-go/gormzap/tags)

## Features


## Usage

### Installation

Use go get.
```bash
    go get github.com/things-go/gormzap
```

Then import the gormzap package into your own code.
```bash
    import gormzap "github.com/things-go/gormzap"
```

### Example

[embedmd]:# (_example/main.go go)
```go
package main

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github/things-go/gormzap"
)

func main() {
	zapL, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	log := gormzap.New(zapL,
		gormzap.WithCustomFields(
			gormzap.Immutable("service", "test"),
			func(ctx context.Context) zap.Field {
				v := ctx.Value("requestId")
				if v == nil {
					return zap.Skip()
				}
				if vv, ok := v.(string); ok {
					return zap.String("requestId", vv)
				}
				return zap.Skip()
			},
		),
		gormzap.WithConfig(logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			Colorful:                  false,
			IgnoreRecordNotFoundError: false,
			LogLevel:                  logger.Info,
		}),
	)
	// your dialector
	db, _ := gorm.Open(nil, &gorm.Config{Logger: log})
	// do your things
	_ = db
}
```

## License

This project is under MIT License. See the [LICENSE](LICENSE) file for the full license text.
