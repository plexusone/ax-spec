# AX Spec Tasks

Tracking experiments and measurements to quantify AX benefits.

## Context

elevenlabs-go is used in production by:

- **omniagent** (`github.com/plexusone/omniagent`) — Voice processing for AI agents
- **videoascode** (`github.com/grokify/videoascode`) — Batch TTS for video generation

Both integrations use elevenlabs-go indirectly via OmniVoice abstraction layer.

### Current Integration Patterns

| Aspect | omniagent | videoascode |
|--------|-----------|-------------|
| Integration | OmniVoice abstraction | OmniVoice abstraction |
| Retry logic | None (delegates to provider) | None |
| Rate limit handling | None | None |
| Error handling | `fmt.Errorf` wrapping | `fmt.Errorf` wrapping |
| Failure recovery | Return error to caller | Resume via manifest checkpoints |
| Pre-flight validation | None | None |
| Batch behavior | N/A | Fail-fast (stops on first error) |

### Opportunities for AX Integration

1. **Rate limit handling** — Neither handles 429 errors with backoff
2. **Error categorization** — Errors are wrapped but not categorized for recovery
3. **Retry logic** — Safe operations (GET) are not retried
4. **Pre-flight validation** — API calls made without validating required fields
5. **Partial batch success** — videoascode stops on first error instead of continuing

---

## Prerequisites

Before running experiments, OmniVoice must be updated to be AX-aware.

### Design Document

**Location:** [omnivoice/docs/design/ax-integration.md](https://github.com/plexusone/omnivoice/blob/main/docs/design/ax-integration.md)

### Implementation Goals

| ID | Goal | Target | How to Measure |
|----|------|--------|----------------|
| G1 | Auto-recover from transient errors | >70% recovery rate | `RetrySuccesses / RetriedCalls` |
| G2 | Reduce wasted API calls | >80% validation catches | `ValidationCatches / (ValidationCatches + API400s)` |
| G3 | Intelligent retry with backoff | >60% retry success | `RetrySuccesses / TotalRetries` |
| G4 | Smart fallback decisions | >90% appropriate | `FallbacksAvoided / TotalRetryableErrors` |
| G5 | Preserve error context | 100% have category | Audit error returns |
| G6 | Minimize latency overhead | <5ms added | Benchmark without retries |
| G7 | Backward compatible | 0 breaking changes | API compatibility check |

### Implementation Phases

| Phase | Scope | Duration | Status |
|-------|-------|----------|--------|
| 1. Core Resilience | omnivoice-core: `resilience/` package | 1 week | Not started |
| 2. AX Bridge | omnivoice: ElevenLabs error classifier | 1 week | Not started |
| 3. Smart Fallback | omnivoice-core: TTS/STT client updates | 0.5 week | Not started |
| 4. elevenlabs-go | ax/: OmniVoice helper functions | 0.5 week | Not started |
| 5. Integration Testing | End-to-end tests, benchmarks | 1 week | Not started |
| 6. Documentation | Docs, changelog, release | 0.5 week | Not started |

### Changes by Repository

**omnivoice-core** (new `resilience` package):

- `resilience/error.go` — Error types, categories, `ProviderError`
- `resilience/retry.go` — Retry logic with generics
- `resilience/backoff.go` — Exponential backoff with jitter
- `tts/tts.go` — Smart fallback (only on permanent errors)
- `stt/stt.go` — Smart fallback (only on permanent errors)

**omnivoice** (AX-aware ElevenLabs provider):

- `providers/elevenlabs/ax_bridge.go` — Error classifier using elevenlabs-go/ax
- `providers/elevenlabs/tts.go` — Retry logic, pre-flight validation, metrics
- `providers/elevenlabs/stt.go` — Retry logic, pre-flight validation, metrics
- `providers/elevenlabs/metrics.go` — Metrics for goal measurement

**elevenlabs-go** (minor additions):

- `ax/omnivoice.go` — Helper functions for OmniVoice integration

### Blocking Tasks

Experiments cannot run until:

- [ ] **P1:** `resilience` package implemented in omnivoice-core
- [ ] **P2:** AX bridge implemented in omnivoice
- [ ] **P3:** Metrics collection added to ElevenLabs provider
- [ ] **P4:** Integration tests passing

---

## Outcome Metrics to Measure

### 1. Error Auto-Recovery Rate

**Definition:** Percentage of errors that can be handled programmatically without human intervention.

**Measurement:**
```go
type ErrorOutcome struct {
    Total           int
    AutoRecovered   int  // Handled by AX-aware logic
    ManualRequired  int  // Required human intervention
    Unrecoverable   int  // Permanent failures (auth, validation)
}

AutoRecoveryRate = AutoRecovered / (Total - Unrecoverable) * 100
```

**Experiment:**
1. Instrument elevenlabs-go calls to log all errors
2. Run videoascode batch TTS on 100+ slides
3. Categorize errors: rate_limit, not_found, auth, validation, server
4. Compare: How many could AX auto-recover vs current fail-fast?

**Target:** "X% of transient errors can be auto-recovered with AX"

---

### 2. Retry Success Rate

**Definition:** Percentage of retryable operations that succeed after retry.

**Measurement:**
```go
type RetryOutcome struct {
    TotalRetryable    int  // Errors on retryable operations
    SucceededOnRetry  int  // Succeeded after 1+ retries
    ExhaustedRetries  int  // Failed after max retries
}

RetrySuccessRate = SucceededOnRetry / TotalRetryable * 100
```

**Experiment:**
1. Add retry logic with AX `IsRetryable()` checks
2. Configure exponential backoff (1s, 2s, 4s, max 3 retries)
3. Run batch operations during peak hours (more rate limits)
4. Measure: How many rate-limited requests succeed on retry?

**Target:** "X% of rate-limited requests succeed with AX retry logic"

---

### 3. Pre-flight Validation Savings

**Definition:** API calls avoided by validating required fields before calling.

**Measurement:**
```go
type ValidationOutcome struct {
    TotalAttempts       int  // All attempted API calls
    CaughtByValidation  int  // Blocked by pre-flight check
    FailedAtAPI         int  // 400 errors from API
    Succeeded           int
}

ValidationSavings = CaughtByValidation / (CaughtByValidation + FailedAtAPI) * 100
```

**Experiment:**
1. Add `ax.ValidateFields()` before all mutation operations
2. Log when validation catches missing fields
3. Run with intentionally incomplete requests
4. Measure: How many 400 errors would be prevented?

**Target:** "X% of validation errors caught before API call"

---

### 4. Batch Completion Rate

**Definition:** Percentage of batch items that complete successfully.

**Measurement:**
```go
type BatchOutcome struct {
    TotalItems      int
    Succeeded       int
    FailedRetryable int  // Could have succeeded with retry
    FailedPermanent int  // Would fail regardless
}

// Current (fail-fast)
CurrentCompletion = Succeeded / TotalItems * 100

// With AX (continue + retry)
AXCompletion = (Succeeded + RecoveredWithRetry) / TotalItems * 100
```

**Experiment (videoascode):**
1. Run batch TTS on 50 slides
2. Inject failures: rate limits at slide 10, 25, 40
3. Current behavior: Stops at slide 10 (20% complete)
4. AX behavior: Retries rate limits, continues on permanent failures
5. Measure: Final completion rate

**Target:** "Batch completion improved from X% to Y% with AX"

---

### 5. Time to Error Resolution

**Definition:** Time from error occurrence to successful resolution.

**Measurement:**
```go
type ResolutionTiming struct {
    ErrorTime      time.Time
    ResolutionTime time.Time
    Method         string  // "auto_retry", "manual_resume", "code_fix"
}

AvgResolutionTime = sum(ResolutionTime - ErrorTime) / count
```

**Experiment:**
1. Time how long videoascode takes to complete with interruptions
2. Current: Error → manual re-run → resume from checkpoint
3. AX: Error → auto-retry → continue
4. Measure: End-to-end time for same workload

**Target:** "Error resolution time reduced from X minutes to Y seconds"

---

### 6. API Cost Efficiency

**Definition:** Ratio of successful API calls to total API calls (including retries).

**Measurement:**
```go
type CostOutcome struct {
    TotalAPICalls    int
    SuccessfulCalls  int
    FailedCalls      int
    RetryAttempts    int
    CharactersBilled int  // ElevenLabs bills by character
}

CostEfficiency = SuccessfulCalls / TotalAPICalls * 100
WastedCalls = FailedCalls / TotalAPICalls * 100
```

**Experiment:**
1. Track all API calls and their outcomes
2. Calculate: calls that succeeded vs wasted calls
3. With AX: Pre-flight validation prevents some wasted calls
4. With AX: Retry logic reduces need for full re-runs

**Target:** "X% reduction in wasted API calls"

---

## Experiments

### Experiment 1: Rate Limit Recovery (videoascode)

**Hypothesis:** AX retry logic will recover from rate limits without manual intervention.

**Setup:**
1. Configure videoascode for 100-slide batch
2. Set ElevenLabs to trigger rate limits (aggressive request rate)
3. Run with current code (baseline)
4. Run with AX-aware retry logic (treatment)

**Metrics:**
- Completion rate (% slides generated)
- Total time to completion
- Number of manual re-runs required
- Total API calls made

**Implementation:**

```go
// Current (videoascode/pkg/omnivoice/tts/provider.go)
result, err := p.provider.Synthesize(ctx, text, config)
if err != nil {
    return nil, fmt.Errorf("%s tts failed: %w", p.name, err)
}

// With AX
result, err := p.provider.Synthesize(ctx, text, config)
if err != nil {
    if code, ok := elevenlabs.GetAXErrorCode(err); ok {
        info := ax.GetErrorInfo(code)
        if info.Retryable {
            // Exponential backoff retry
            for attempt := 1; attempt <= 3; attempt++ {
                time.Sleep(time.Duration(attempt*attempt) * time.Second)
                result, err = p.provider.Synthesize(ctx, text, config)
                if err == nil {
                    break
                }
            }
        }
    }
    if err != nil {
        return nil, fmt.Errorf("%s tts failed: %w", p.name, err)
    }
}
```

---

### Experiment 2: Error Categorization (omniagent)

**Hypothesis:** AX error categories enable smarter recovery strategies.

**Setup:**
1. Instrument omniagent voice processor with error logging
2. Run agent workflows that trigger various errors
3. Categorize errors and measure recovery options

**Metrics:**
- Error distribution by category
- % of errors with clear recovery path
- % of errors requiring human intervention

**Error Categories to Track:**

| Category | Example | Recovery Strategy |
|----------|---------|-------------------|
| rate_limit | 429 Too Many Requests | Exponential backoff |
| not_found | Voice ID doesn't exist | Use default voice |
| auth | Invalid API key | Alert user |
| validation | Text too long | Split into chunks |
| server | 500 Internal Error | Retry with backoff |
| quota | Monthly limit exceeded | Alert user |

---

### Experiment 3: Validation Savings (both)

**Hypothesis:** Pre-flight validation prevents wasted API calls.

**Setup:**
1. Add logging for all API call attempts
2. Add `ax.ValidateFields()` checks before calls
3. Run workflows with various input completeness

**Test Cases:**
- Missing voice_id (required for TTS)
- Missing text (required for TTS)
- Invalid model_id
- Empty audio for STT

**Metrics:**
- API calls prevented by validation
- Time saved (no round-trip for invalid requests)
- Error message quality (specific vs generic)

---

### Experiment 4: End-to-End Agent Success (omniagent)

**Hypothesis:** AX integration improves agent task completion rate.

**Setup:**
1. Define a standard voice task (transcribe → process → synthesize)
2. Run 100 iterations with current code
3. Run 100 iterations with AX-enhanced code
4. Inject transient failures (rate limits, network errors)

**Metrics:**
- Task success rate (%)
- Average task completion time
- Human intervention rate
- Error recovery rate

---

## Implementation Plan

### Phase 0: Prerequisites (Weeks 1-3)

See [Prerequisites](#prerequisites) section above. Full design: [omnivoice/docs/design/ax-integration.md](https://github.com/plexusone/omnivoice/blob/main/docs/design/ax-integration.md)

- [ ] **Week 1:** Implement `resilience` package in omnivoice-core
- [ ] **Week 2:** Implement AX bridge and update ElevenLabs provider in omnivoice
- [ ] **Week 3:** Integration testing and metrics validation

### Phase 1: Baseline Measurement (Week 4)

- [ ] Run videoascode batch TTS (100 slides) with current code
- [ ] Run omniagent voice workflows with current code
- [ ] Record baseline metrics:
  - Completion rate
  - Error distribution by type
  - Manual intervention count
  - Total time to completion

### Phase 2: Experiments with AX (Week 5)

- [ ] Run Experiment 1: Rate Limit Recovery (videoascode)
- [ ] Run Experiment 2: Error Categorization (omniagent)
- [ ] Run Experiment 3: Validation Savings (both)
- [ ] Run Experiment 4: End-to-End Agent Success (omniagent)

### Phase 3: Analysis & Documentation (Week 6)

- [ ] Compare baseline vs AX-enabled metrics
- [ ] Calculate improvement percentages
- [ ] Update case studies with quantitative results
- [ ] Create before/after comparison charts
- [ ] Document best practices from experiments
- [ ] Publish findings

---

## Success Criteria

| Metric | Current (estimated) | Target |
|--------|---------------------|--------|
| Error auto-recovery rate | 0% | >50% |
| Retry success rate | N/A | >70% |
| Batch completion rate | ~60% (with failures) | >90% |
| Pre-flight validation catches | 0% | >80% of validation errors |
| Manual intervention rate | High | <10% of errors |

---

## Notes

### Why OmniVoice is the Integration Point

Both omniagent and videoascode use OmniVoice as their voice abstraction layer. This means:

1. **Single integration point** — Adding AX to OmniVoice benefits both projects
2. **Provider-agnostic patterns** — AX patterns can apply to other providers too
3. **Clean separation** — elevenlabs-go provides AX data, OmniVoice uses it

### Integration Architecture

```
omniagent/videoascode
        │
        ▼
    OmniVoice (abstraction)  ← Add AX-aware error handling here
        │
        ▼
    elevenlabs-go + ax/      ← AX metadata source
        │
        ▼
    ElevenLabs API
```

### Related Repositories

- `github.com/plexusone/elevenlabs-go` — SDK with AX package
- `github.com/plexusone/omnivoice` — Voice abstraction layer
- `github.com/plexusone/omniagent` — AI agent framework
- `github.com/grokify/videoascode` — Video generation tool
- `github.com/plexusone/ax-spec` — AX specification and tools
