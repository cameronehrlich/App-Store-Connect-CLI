# ASC CLI - Project Plan

## Vision

A fast, AI-agent-friendly CLI for App Store Connect that enables developers to ship iOS apps with zero friction.

**Problem:** Manual App Store Connect workflows are slow, and there’s no AI‑agent‑friendly CLI with clean JSON output.

**Solution:** A lightweight Go-based CLI focused on speed, simplicity, and AI-first design.

---

## Current Reality (v0.1 - Implemented, Validated Locally)

**Last Updated:** 2026-01-20

### What Works

- Project structure and Go module setup
- CLI skeleton with ffcli framework
- Commands register and show help: `asc --help`, `asc feedback --help`, `asc auth --help`
- ECDSA JWT signing wired to `.p8` keys
- Feedback/crash/review endpoints aligned to ASC OpenAPI spec
- Code compiles and unit tests run
- Live API validation: feedback/crashes return data; reviews may be empty if no reviews exist

### What Doesn't Work Yet

- **Pagination** - No paging support yet for large datasets
- **Filtering** - Reviews support rating/territory filters; feedback/crashes are unfiltered
- **Keychain** - Keychain support implemented; local config fallback remains
- **Tests** - Basic unit coverage (config, auth, query, output)

### Files Status

```
main.go              ✓ Compiles, entry point
cmd/commands.go      ✓ Compiles, commands defined
cmd/auth.go          ✓ Compiles, login/logout/status defined
internal/asc/        ✓ JWT signing with ECDSA (ES256)
internal/auth/       ✓ Keychain support with config fallback
internal/config/     ✓ Basic config file handling
Makefile             ✓ Build targets exist
```

---

## What Was Actually Delivered

### For TestFlight Feedback (Screenshot Submissions)

**Expected Data:**
- Submission ID
- Created date
- Tester email
- Comment text
- Screenshot submission metadata

**Current Status:** ✅ Implemented (requires live credentials to verify)

```bash
asc feedback --app "123456789" --json
# Requires valid credentials
```

**Required in ASC API:**
- `GET /v1/apps/{id}/betaFeedbackScreenshotSubmissions`

### For Crash Reports (Crash Submissions)

**Expected Data:**
- Submission ID
- Created date
- Tester email
- Crash log (when available)
- Device metadata (model, OS version)

**Current Status:** ✅ Implemented (requires live credentials to verify)

```bash
asc crashes --app "123456789" --json
# Requires valid credentials
```

**Required in ASC API:**
- `GET /v1/apps/{id}/betaFeedbackCrashSubmissions`

### For App Store Reviews

**Expected Data:**
- Review ID
- Created date
- Rating
- Title
- Body
- Territory

**Current Status:** ✅ Implemented (requires live credentials to verify)

```bash
asc reviews --app "123456789" --json
# Requires valid credentials
```

**Required in ASC API:**
- `GET /v1/apps/{id}/customerReviews`

---

## Roadmap

### Phase 1: Foundation - REVISED (Current)

**Goal:** Ensure API calls work and add basic filters where supported

#### What Needs Fixing

- [x] Test actual API authentication with App Store Connect
- [ ] Add pagination support
- [ ] Add integration tests for auth flow

#### Deliverables (Actual)

```
✓ go.mod
✓ main.go with ffcli
✓ cmd/commands.go (feedback, crashes, reviews, auth)
✓ cmd/auth.go (login, logout, status)
✓ internal/asc/client.go - JWT signing with ECDSA and ASC endpoints
✓ internal/auth/keychain.go - Keychain support with config fallback
✓ internal/config/config.go
✓ Makefile
✓ CLAUDE.md
```

---

### Phase 2: Core Features (v0.1) - REVISED

**Goal:** Validate feedback, crashes, and reviews against live API

#### Features (Actual Implementation)

1. **TestFlight Feedback Command**
   ```bash
   asc feedback --app "APP_ID" --json
   # Must return: id, createdDate, email, comment, screenshot metadata
   ```

2. **TestFlight Crashes Command**
   ```bash
   asc crashes --app "APP_ID" --json
   # Must return: id, createdDate, email, crash metadata
   ```

3. **App Store Reviews Command**
   ```bash
   asc reviews --app "APP_ID" --json
   # Must return: id, createdDate, rating, title, body, territory
   ```

#### Technical Tasks

- [x] Verify API endpoint paths against ASC OpenAPI spec:
  - `GET /v1/apps/{id}/betaFeedbackScreenshotSubmissions`
  - `GET /v1/apps/{id}/betaFeedbackCrashSubmissions`
  - `GET /v1/apps/{id}/customerReviews`
- [x] Add review filters (`--stars`, `--territory`) via query params
- [ ] Add pagination support for all list endpoints
- [ ] Add feedback/crash filters where supported (device model, OS version, etc.)
- [ ] Write integration tests (mock or real API)

---

### Phase 3: App Management (v0.2)

**Goal:** Add commands for managing apps and builds

#### Features

1. **List Apps**
   ```bash
   asc apps list
   asc apps list --json
   ```

2. **List Builds**
   ```bash
   asc builds list --app "APP_ID"
   asc builds list --app "APP_ID" --json
   ```

3. **Build Details**
   ```bash
   asc builds info --build "BUILD_ID"
   ```

4. **Expire Build**
   ```bash
   asc builds expire --build "BUILD_ID" --app "APP_ID"
   ```

#### Technical Tasks

- [ ] Implement `GET /v1/apps`
- [ ] Implement `GET /v1/apps/{id}/builds`
- [ ] Implement `PATCH /v1/builds/{id}`
- [ ] Add build expiration workflow
- [ ] Add pagination support

---

### Phase 4: Beta Management (v0.3)

**Goal:** Add commands for managing beta testers and groups

#### Features

1. **Beta Groups**
   ```bash
   asc beta-groups list --app "APP_ID"
   asc beta-groups create --app "APP_ID" --name "Beta Testers"
   ```

2. **Beta Testers**
   ```bash
   asc beta-testers list --app "APP_ID"
   asc beta-testers add --app "APP_ID" --email "tester@example.com" --group "Beta"
   asc beta-testers remove --app "APP_ID" --email "tester@example.com"
   ```

3. **Tester Invitations**
   ```bash
   asc beta-testers invite --app "APP_ID" --email "tester@example.com"
   ```

#### Technical Tasks

- [ ] Implement `GET /v1/apps/{id}/betaGroups`
- [ ] Implement `POST /v1/betaGroups`
- [ ] Implement `GET /v1/apps/{id}/betaTesters`
- [ ] Implement `POST /v1/betaTesters`
- [ ] Implement `DELETE /v1/betaTesters/{id}`

---

### Phase 5: Localization (v0.4)

**Goal:** Add commands for managing app localizations

#### Features

1. **Download Localizations**
   ```bash
   asc localizations download --app "APP_ID" --output ./locales/
   ```

2. **Upload Localizations**
   ```bash
   asc localizations upload --app "APP_ID" --path ./locales/
   ```

3. **List Localizations**
   ```bash
   asc localizations list --app "APP_ID"
   ```

#### Technical Tasks

- [ ] Implement `GET /v1/apps/{id}/appStoreVersions`
- [ ] Implement `GET /v1/apps/{id}/appStoreVersionLocalizations`
- [ ] Implement file upload/download for localization files
- [ ] Add validation for localization files

---

### Phase 6: Submission (v0.5)

**Goal:** Add commands for submitting apps

#### Features

1. **Submit for Review**
   ```bash
   asc submit --app "APP_ID" --build "BUILD_ID" --submit-type "APP_STORE"
   ```

2. **Check Submission Status**
   ```bash
   asc submit status --app "APP_ID"
   ```

3. **Cancel Submission**
   ```bash
   asc submit cancel --app "APP_ID"
   ```

#### Technical Tasks

- [ ] Implement `POST /v1/appStoreVersionSubmissions`
- [ ] Implement `GET /v1/appStoreVersionSubmissions/{id}`
- [ ] Implement `DELETE /v1/appStoreVersionSubmissions/{id}`
- [ ] Add workflow for submission process

---

### Phase 7: Sandbox & Testing (v0.6)

**Goal:** Add commands for sandbox testing

#### Features

1. **Create Sandbox Tester**
   ```bash
   asc sandbox create --email "tester@example.com" --territory US
   ```

2. **List Sandbox Testers**
   ```bash
   asc sandbox list
   ```

3. **Delete Sandbox Tester**
   ```bash
   asc sandbox delete --email "tester@example.com"
   ```

#### Technical Tasks

- [ ] Implement `GET /v1/sandboxTesters`
- [ ] Implement `POST /v1/sandboxTesters`
- [ ] Implement `DELETE /v1/sandboxTesters/{id}`

---

### Phase 8: Analytics (v0.7)

**Goal:** Add commands for viewing analytics

#### Features

1. **Sales Report**
   ```bash
   asc analytics sales --app "APP_ID" --from 2024-01-01 --to 2024-01-31
   ```

2. **Usage Report**
   ```bash
   asc analytics usage --app "APP_ID" --last-30-days
   ```

#### Technical Tasks

- [ ] Implement `GET /v1/salesReports`
- [ ] Implement `GET /v1/usageMetrics`
- [ ] Add report formatting

---

### Future Enhancements (v1.0+)

- **Interactive Mode** - TUI for manual exploration
- **Plugins** - Extendable architecture
- **AI Summarization** - Use local LLMs to summarize feedback
- **Auto-Responder** - AI-powered response to reviews
- **Multi-Account** - Manage multiple ASC accounts
- **Web Dashboard** - Optional web UI

---

## Dependencies

### Core Dependencies

```
github.com/peterbourgon/ff/v3     - CLI framework
github.com/golang-jwt/jwt/v5      - JWT signing
```

### Dev Tools (optional)

```
golangci-lint                     - Linting (if installed)
gosec                             - Security scanning (if installed)
github.com/goreleaser/nfpm/v2     - Packaging via `go run` (optional)
```

---

## Release Strategy

1. **Alpha** - Internal testing, core features
2. **Beta** - Public testing, feedback collection
3. **v0.1** - First public release (feedback, crashes, reviews)
4. **v0.x** - Incremental feature releases
5. **v1.0** - Full feature set, stable API

---

## Success Metrics

- Install via Homebrew: `brew install rudrank/tap/asc`
- Average startup time: < 50ms
- All commands support `--json` flag
- 80%+ test coverage
- Zero security vulnerabilities

---

## Current Status

**Phase 1: Foundation - IMPLEMENTED** (validated locally)

Next: Add pagination, feedback/crash filters, and integration tests

## Known Issues

1. **Pagination**
   - No pagination support yet for large datasets

2. **Filtering**
   - Reviews support rating/territory filters
   - Feedback/crashes do not yet expose API filters

3. **Keychain**
   - Keychain supported; local config fallback still exists

---

## Success Criteria for v0.1

- [x] `asc feedback --app "APP_ID" --json` returns screenshot feedback submissions
- [x] Feedback includes: id, createdDate, email, comment
- [x] `asc crashes --app "APP_ID" --json` returns crash submissions
- [x] Crashes include: id, createdDate, email, crash metadata
- [x] `asc reviews --app "APP_ID" --json` returns customer reviews (may be empty)
- [x] Reviews include: id, createdDate, rating, title, body, territory
- [x] All commands work with real App Store Connect API keys
- [ ] Basic integration test that doesn't require real API credentials
