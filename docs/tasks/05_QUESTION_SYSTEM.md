# Question System Tasks

## QUES-01: Create Question Catalog and Immutable Version Model
**Priority:** Critical · **Phase:** 3 Question/content foundation · **Component:** Database / Backend · **Effort:** M · **Parallelism:** Central content model owner

**Description:** Model categories, difficulty, types, localized questions and immutable published versions. **Business Requirements:** PRD §5,8; GDD §3; DB §5. **Technical Requirements:** Unique version per question/locale, status lifecycle, media and accessibility text, answer key restricted to server. **Dependencies:** BACK-04, FOUND-04.

**Implementation Steps:** 1. Add category/question/version migrations. 2. Add constraints/indexes. 3. Define content schema for supported question types. 4. Add migration fixtures. **Expected Files / Modules:** `internal/question`, SQL migrations, contract schemas. **API / Database Impact:** Question tables and admin routes.

**Acceptance Criteria:** Published content is immutable; localized versions are independently identified; invalid answer/options fail validation. **Testing Requirements:** Migration, schema, uniqueness and permission tests. **Completion Definition:** Data model accepted by admin, duel and localization owners.

## QUES-02: Implement Question Draft, Review and Publication
**Priority:** High · **Phase:** 3 / P1 admin · **Component:** Admin / Backend · **Effort:** L · **Parallelism:** Admin UX can parallelize after QUES-01

**Description:** Provide backend content lifecycle use cases for create/edit/preview/review/publish/archive and bulk import. **Business Requirements:** PRD §5 P1; Roadmap §3; UI/UX §5. **Technical Requirements:** Role-gated publish; version history; duplicate detection; media validation; publication audit. **Dependencies:** QUES-01, ADMIN-01.

**Implementation Steps:** 1. Implement draft CRUD and validation. 2. Add preview and version diff. 3. Add reviewer/publisher permissions. 4. Validate bulk import and report row errors. 5. Audit publish/archive. **Expected Files / Modules:** Admin question feature and question application service. **API / Database Impact:** Admin question routes, version and audit records.

**Acceptance Criteria:** Draft edits do not mutate a published version or active duel; unauthorized role cannot publish; bulk import is all-or-explicitly-partial with row report. **Testing Requirements:** Permission, version immutability, duplicate and import validation. **Completion Definition:** Content editor workflow accepted by content owner.

## QUES-03: Implement Balanced Question Selection and Match Snapshot
**Priority:** Critical · **Phase:** 3-4 · **Component:** Backend · **Effort:** L · **Parallelism:** Coordinated with BATTLE

**Description:** Select deterministic, eligible question sequence for a duel before it starts. **Business Requirements:** GDD §3; Roadmap Epic 3. **Technical Requirements:** Locale/category/difficulty balance, no immediate player repeats where possible, exposure controls, immutable selected version and rules snapshot. **Dependencies:** QUES-01, BATTLE-01, FOUND-02.

**Implementation Steps:** 1. Define selection constraints/config. 2. Select/persist sequence at match initialization. 3. Record exposure metadata safely. 4. Define fallback when content pool is insufficient. 5. Add repeat/quality telemetry. **Expected Files / Modules:** selector, repository queries and config. **API / Database Impact:** duel_questions and selected version references.

**Acceptance Criteria:** Both participants receive same ordered IDs/versions; sequence is fixed for match lifetime; insufficient pool fails explicitly or follows approved fallback. **Testing Requirements:** Determinism, distribution, locale fallback and concurrent selection tests. **Completion Definition:** Match fixture proves stable sequence across restart.

## QUES-04: Secure Question Delivery and Answer-Key Separation
**Priority:** Critical · **Phase:** 4 · **Component:** Backend / Security · **Effort:** M · **Parallelism:** Central review with BATTLE

**Description:** Deliver localized question content without exposing correctness before round resolution. **Business Requirements:** PRD §8; GDD §3; API §12. **Technical Requirements:** Projection allowlist; answer key excluded from WebSocket, REST, logs, analytics and client bundles; reveal only after valid resolve. **Dependencies:** QUES-01, QUES-03, BACK-08.

**Implementation Steps:** 1. Define separate public question and private answer-key types. 2. Implement delivery projection. 3. Implement post-resolution answer disclosure. 4. Scan logs/traces for key leakage. **Expected Files / Modules:** question serializer, duel event contracts and redaction tests. **API / Database Impact:** question and resolution events.

**Acceptance Criteria:** Pre-resolution payload never contains answer key or explanation if it reveals answer; post-resolution payload reveals only after both submissions/deadline. **Testing Requirements:** Contract snapshot and serialization security tests. **Completion Definition:** Security review verifies no answer leakage across all transports.

## QUES-05: Implement Server-Side Answer Validation
**Priority:** Critical · **Phase:** 4 · **Component:** Backend · **Effort:** M · **Parallelism:** One owner with BATTLE-04

**Description:** Validate answer shape/correctness using selected question version and server receive time. **Business Requirements:** PRD §8; GDD §4; API §7. **Technical Requirements:** Ignore client correctness/time, enforce membership/round/deadline/question ID/duplicate answer, supported question types and stable error codes. **Dependencies:** QUES-01, QUES-03, BACK-06.

**Implementation Steps:** 1. Implement typed answer validators. 2. Compare server receipt against authoritative deadline. 3. Insert unique answer idempotently. 4. Handle timeout and invalid payload consistently. 5. Persist scoring inputs for audit. **Expected Files / Modules:** answer validator/use case and DB constraints. **API / Database Impact:** `/duels/{id}/answers`, `duel.answer_submit`; duel_answers.

**Acceptance Criteria:** Client-supplied response time/correctness is ignored; one answer per round; late/invalid/foreign answers cannot score. **Testing Requirements:** Boundary timestamp, duplicate, malformed and forged-answer tests. **Completion Definition:** REST and WebSocket paths pass the same validator suite.

## QUES-06: Add Question Quality, Localization and Exposure QA
**Priority:** Medium · **Phase:** 3-8 · **Component:** Admin / Testing · **Effort:** M · **Parallelism:** Content QA can run parallel after QUES-02

**Description:** Establish content validation and operational quality controls. **Business Requirements:** GDD §3; UI/UX §7; Roadmap §3. **Technical Requirements:** Correctness review, ambiguity, locale quality, accessible alt text, rights, duplicate/repetition detection and difficulty audit. **Dependencies:** QUES-02, QUES-03.

**Implementation Steps:** 1. Define editorial checklist and reviewer roles. 2. Add validation reports for imports. 3. Track question exposure/performance. 4. Add deactivate/replace workflow without mutating active snapshots. **Expected Files / Modules:** admin QA panel, reports, operations checklist. **API / Database Impact:** exposure metadata and content audit.

**Acceptance Criteria:** Publish gate checks required fields and reviewer approval; content issues can be removed from future selection without altering active matches. **Testing Requirements:** Missing alt text, duplicate, localization, invalid media and publication tests. **Completion Definition:** Approved content QA checklist and initial reviewed bank.

## QUES-07: Add Question Selection and Answer Contract Tests
**Priority:** High · **Phase:** 4 · **Component:** Testing · **Effort:** M · **Parallelism:** Parallel after schema freeze

**Description:** Verify question and answer contracts independently of UI. **Business Requirements:** Roadmap §4-5. **Technical Requirements:** Contract fixtures for each question type and locale, answer-key security test, deterministic sequence. **Dependencies:** QUES-03, QUES-04, QUES-05.

**Implementation Steps:** 1. Create fixtures for published/localized questions. 2. Validate API/event schemas. 3. Test active-match snapshot through edit/publish changes. 4. Add regression cases for exposure/fallback. **Expected Files / Modules:** question test fixtures and contract CI. **API / Database Impact:** None.

**Acceptance Criteria:** Contract and security tests fail on answer-key leakage or mutation of selected version. **Testing Requirements:** Unit, integration and schema tests. **Completion Definition:** CI protects all supported question types and event variants.
