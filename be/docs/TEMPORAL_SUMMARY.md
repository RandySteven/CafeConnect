# Temporal Integration – Summary (CafeConnect BE)

Summary of the Temporal workflow integration work and findings.

---

## 1. What We Built

- **Transaction checkout workflow (v3)**  
  Checkout runs as a Temporal workflow that runs activities in order:  
  `checkUser` + `checkCafe` (parallel) → `checkFranchise` → `saveTransactionHeader` → `publishTransaction` → return receipt.
- **Local activities**  
  All steps run as **local activities** (in the same worker process) to avoid round-trips to the server and keep latency low.
- **Worker lifecycle**  
  One worker polls the `cafe_connect` task queue; it is started in `main` and stopped on shutdown.

---

## 2. Fixes and Changes Made

| Issue | Cause | Fix |
|-------|--------|-----|
| **Missing task queue name** | `taskQueue` on `temporalClient` was never set. | In `NewTemporalClient`, set `taskQueue` from config (default `"cafe_connect"`) and pass it to both the worker and the struct. |
| **Context type in workflow** | Workflow used `context.Context` and `context.WithValue`. | Use `workflow.Context` for the workflow and `workflow.WithValue` where needed. |
| **Context deadline exceeded in workflow** | Several causes (see below). | Activity timeouts, correct context types, data via return values, and `GetWorkflowResult` with `context.Background()`. |
| **Activities not running / wrong types** | Activities took `workflow.Context`; data was “passed” via context. | Activities now take `context.Context` and return results; workflow passes those results to the next activity. |
| **Workflow logic bug** | `workflow.Go` + `select` with `default` didn’t wait for the goroutine. | Removed that pattern; workflow runs activities in a clear sequence (parallel where intended, then sequential). |
| **No workers polling** | Worker was never started; workflow/activities never registered. | Call `registerWorkflowAndActivities()` in `NewTransactionWorkflow`; in `main`, call `app.Workflow.Start()` and `defer app.Workflow.Stop()`. |
| **Slow workflow** | Each `ExecuteActivity` did a full round-trip to the server. | Switched to `ExecuteLocalActivity` and `LocalActivityOptions` so work runs in-process. |
| **“Failed reaching server: context deadline exceeded”** | Namespace `cafe_connect` didn’t exist; first gRPC call (GetSystemInfo) times out after 5s. | Start server with `-n cafe_connect`; increased `GetSystemInfoTimeout` to 15s and added hints in config/logs. |

---

## 3. Where “Context Deadline Exceeded” Comes From

- **When:** During **`client.NewClient(...)`** in `pkg/temporal/temporal.go`.
- **What happens:** The SDK immediately does one gRPC call to the server: **GetSystemInfo** (to load capabilities).
- **Timeout:** That call uses a **5 second** timeout by default (`getSystemInfoTimeout` in the SDK).
- **If the server is down or the namespace is wrong:** The call never succeeds; after 5s the context hits its deadline → `context.DeadlineExceeded` → SDK wraps it as **“failed reaching server: context deadline exceeded”**.
- **Our change:** We set `ConnectionOptions.GetSystemInfoTimeout = 15 * time.Second` and added a comment in code so it’s clear what triggers the error.

---

## 4. Running Temporal Locally

- **Start the server** (creates namespace `cafe_connect`):
  ```bash
  temporal server start-dev --ui-port 8080 --db-filename cafe_connect.db -n cafe_connect
  ```
- **Config** (`files/yml/cofeConnect.local.yaml`):  
  `host: localhost`, `port: "7233"`, `namespace: "cafe_connect"`, `taskQueue: "cafe_connect"`.
- **Without `-n cafe_connect`** only the `default` namespace exists; the app will fail to connect (or hit the same timeout) if it uses `cafe_connect`.

---

## 5. Important Files

| File | Role |
|------|------|
| `pkg/temporal/temporal.go` | Temporal client and worker; task queue; `GetSystemInfoTimeout`; connection error hint. |
| `usecases/transactions/workflow.go` | Registers workflow + activities; `CheckoutTransactionV3` starts workflow and waits for result. |
| `usecases/transactions/transaction.go` | Workflow function: runs local activities in order, passes data via return values. |
| `usecases/transactions/check_*.go`, `save_transaction_header.go`, `publish_transaction.go` | Activity implementations: take `context.Context`, return `(result, error)`. |
| `cmd/main/http/main.go` | Starts Temporal worker (`app.Workflow.Start()`) and defers `Stop()`. |
| `apps/app.go` | Builds app; creates Temporal client (app fails to start if Temporal is unreachable). |

---

## 6. Takeaways

1. **Workflow vs activity context**  
   Workflow code uses `workflow.Context`; activities use `context.Context`. Don’t mix them.
2. **Data between activities**  
   Use return values and workflow-level variables. Context values don’t propagate across activity boundaries.
3. **Local vs remote activities**  
   For short, in-process work (DB, publish), local activities avoid scheduling latency; use regular activities for long or external work.
4. **Namespace must exist**  
   For `temporal server start-dev`, create the app’s namespace with `-n <namespace>`.
5. **First gRPC call**  
   “Context deadline exceeded” at client creation is from the initial **GetSystemInfo** call and its timeout; fixing server/namespace or increasing `GetSystemInfoTimeout` addresses it.

---

*Summary of work done on Temporal integration for CafeConnect backend.*
