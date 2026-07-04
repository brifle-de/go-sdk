# Project

This project is and sdk for the Brifle API that wrapps the openapi specs defined at:
@openapi.yaml

The sdk has the goal to prevent breaking changes, by updating namings in the openapi specs like module name changing. Therefore it wrapps all the features of the Brifle API in separate static modules, that supports versioning.

Always write tests for all the endpoints you use. Test against the actual api breaks without authentication. For Testing use the mock server that contains the expected files: @test/mock for the server endpoints. Mocks are used to test the converting between the openapi specs and the sdk. The mock server is a simple express server that returns the expected responses for the endpoints defined in the openapi specs.

# Architecture

The SDK has three layers:

1. `sdk/api/client.gen.go` — the low-level transport, **auto-generated** by `oapi-codegen` from
   `openapi.yaml`. Do **not** hand-edit it. Regenerate with `make -f Makefile generate` (i.e.
   `go generate ./...`). Config: `cfg.yaml`, trigger: `generate.go`.
2. `sdk/client`, `sdk/middleware`, `sdk/sdk.go` — the client wrapper, automatic auth/token renewal,
   and helpers (`sdk.String`, `sdk.Base64Encode`, `sdk.Base64EncodeString`).
3. `sdk/endpoints/*` — hand-written, stable wrapper packages. These are the insulation layer: they
   absorb generated-type churn so external callers see a stable API.

## Wrapper pattern

Every endpoint is a package-level function:

```go
func Name(client *sdkClient.BrifleClient, ctx context.Context, ...args) (*Result, *api.ResponseStatus, error)
```

Guard nil inputs → map a friendly hand-written struct to the generated `api.*` type (via a
`ToApi...()`/`toApiRequest()` converter) → call `client.ApiClient.WebApiController...` → decode with
`api.ValidateHttpResponse(err, resp, &res)` for JSON or `api.ParseResponseAsString` /
`api.ParseResponseAsBytes` for XML/PDF bodies. Return `(&res, status, nil)`. Wrap JSON responses in a
thin struct embedding the generated type (e.g. `type X struct { *api.Y }`). Non-2xx HTTP is reported
via `status.HttpStatus`, not `error`.

## When adding an endpoint from the spec

1. **Regenerate first** (`make generate`) so the generated method + types exist. Regeneration can
   change existing generated type shapes when the spec's `required` lists change (optional→required
   flips pointer fields to value fields) — fix the affected wrappers, that is exactly their job.
2. Add the wrapper in the matching `sdk/endpoints/*` package (or a new package for a new group).
3. Add a skip-safe test mirroring `content_endpoint_test.go`.
4. Document it in `docs/` (see below).

## Endpoint packages

`accounts`, `address`, `auth`, `content` (incl. cover letters, delivery status, receiver bulk check,
paper-mail preview), `mailbox`, `signatures`, `status`, `tenants`, `wallet`.

# Developer documentation

External-developer docs live in `docs/` (`docs/README.md` is the index, one file per endpoint group,
each with runnable Go examples). Keep them in sync when you add or change an endpoint.