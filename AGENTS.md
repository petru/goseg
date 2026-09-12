
**Core Formatting & Tooling**

* **Always Format:** Run `gofmt` and `goimports` on all generated code. Code must strictly adhere to standard Go formatting.
* **Linting:** Write code that passes standard `golangci-lint` checks (e.g., `errcheck`, `staticcheck`, `gosec`).

**Naming Conventions**

* **Casing:** Use `MixedCaps` or `mixedCaps` (camelCase) rather than underscores.
* **Acronyms:** Keep initialisms consistently cased (e.g., `ServeHTTP`, `userID`, `parseJSON`).
* **Brevity:** Use short, descriptive variable names (e.g., `c` or `cli` for `Client`, `idx` for index). Context should provide meaning.
* **Packages:** Use short, single-word, lowercase names for packages (e.g., `http`, `log`, `auth`). Avoid `util` or `common`.

**Error Handling**

* **Explicit Checks:** Always handle errors explicitly with `if err != nil { ... }`. Never ignore or swallow them.
* **Contextual Wrapping:** Wrap errors with actionable context using `fmt.Errorf("failed to fetch user %q: %w", id, err)`.
* **No Panics:** Avoid `panic` in library code. Return errors and let the caller decide how to handle exceptional states.

**Interfaces & Types**

* **Small Interfaces:** Keep interfaces small, ideally 1-2 methods (e.g., `io.Reader`).
* **Accept Interfaces, Return Structs:** Functions should accept interfaces (if they need to abstract behavior) but return concrete implementations (structs) to avoid forcing allocation and abstraction on the caller.

**Concurrency & Context**

* **Communication:** "Don't communicate by sharing memory; share memory by communicating." Prefer channels over mutexes when coordinating data flow.
* **Context:** Always pass `context.Context` as the first argument to functions performing I/O, network requests, or long-running operations. Use it for cancellation and timeouts.
* **Goroutine Leaks:** Never start a goroutine without knowing exactly how and when it will stop.

**Testing**

* **Table-Driven Tests:** Structure tests using slices of anonymous structs to handle multiple test cases cleanly.
* **Standard Library:** Rely on the standard `testing` package. Avoid complex third-party assertion libraries unless explicitly requested.
