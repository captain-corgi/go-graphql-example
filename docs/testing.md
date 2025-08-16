# Testing Strategy

This project adopts a pragmatic testing approach with first-class support for mocks generated from interfaces.

## Goals

- Keep domain logic easily testable without frameworks.
- Use generated mocks for interaction testing around application and interfaces layers.
- Favor table-driven tests and clear expectations.

## Libraries

- testing (stdlib)
- gomock: expectations and verifications (`github.com/golang/mock/gomock`)
- mockgen: code generator for gomock (`github.com/golang/mock/mockgen`)

## Directory & Package Conventions

- For any package that defines interfaces, place generated mocks in a sibling directory named `mocks/` under the same parent directory.
- Generated mocks use package name `mocks`.

Examples:

- Interface file: `internal/domain/user/repository.go`
- Generated file: `internal/domain/user/mocks/mock_repository.go` (package `mocks`)

## go:generate Directives (per interface file)

Add a `//go:generate` directive at the top of each file that declares one or more interfaces. This generates mocks for ALL interfaces in that file into `./mocks` with package name `mocks`.

Recommended (pin version via `go run`):

```go
//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks
```

Alternative (uses installed binary):

```go
//go:generate mockgen -source=$GOFILE -destination=./mocks/mock_$GOFILE -package=mocks
```

Notes:

- `$GOFILE` is expanded by `go generate` to the current file name.
- Ensure the `mocks/` directory exists or let mockgen create it by providing the path.
- Run generation from the module root:

```bash
go generate ./...
```

## Optional: Pin tool dependency

Add a tools file so the module records the mockgen dependency version (no runtime impact):

```go
//go:build tools
// +build tools

package tools

import _ "github.com/golang/mock/mockgen"
```

Place it at `internal/tools/tools.go` (or any location), commit it, and run `go mod tidy` when you first use it.

## Using gomock in tests

Example skeleton showing controller lifecycle and expectations:

```go
package mypkg_test

import (
    "testing"

    "github.com/golang/mock/gomock"
    mymocks "github.com/captain-corgi/go-graphql-example/internal/domain/user/mocks"
)

func TestSomething(t *testing.T) {
    ctrl := gomock.NewController(t)
    t.Cleanup(ctrl.Finish)

    repo := mymocks.NewMockUserRepository(ctrl)
    repo.EXPECT().FindByID(gomock.Any(), "u1").Return(/* domain.User */ nil, nil)

    // invoke system under test using repo
}
```

## Guidance

- Prefer table-driven tests for multiple cases.
- Keep mocks focused on behavior under test; avoid over-specifying expectations.
- Generate or update mocks whenever interfaces change: `go generate ./...`.
- Do not implement production code here; this document only defines conventions.
