# go-version-lib — Improvement Plan

## Current state

- 7 ecosystem packages (`semver`, `npm`, `pep440`, `deb`, `rpm`, `gem`, `mvn`) + generic
  `version` + `part`/`prerelease` primitives + `purl` + `reference` auto-detection + `goversion` CLI.
- ~9k LOC, testify-based tests, `xerrors` for errors, Go 1.22, no CI config in repo.

## Findings / technical debt

1. **`pkg/reference` boilerplate**: 7 near-identical wrapper structs each re-implementing the
   same 8 `Reference` methods via type assertions (`other.(*semverRef)`), which **panic** when
   comparing across ecosystems instead of returning an error or a defined ordering.
2. **`reference.Detect`/`NewReference`**: linear regex scan over all entries on every call;
   `sortedEntries()` appears to re-sort per call (allocation hot path).
3. **No CI**: no `.github/workflows` — tests/lint never run automatically.
4. **Dated deps**: `xerrors` (superseded by stdlib `fmt.Errorf`/`errors`), testify 1.10,
   Go directive 1.22.
5. **No fuzzing / conformance harness**: parsers are regex-heavy and ecosystem rules are subtle
   (PEP 440 normalization, npm partial versions, Maven token ordering).
6. **Constraint fragmentation**: each ecosystem has its own constraint type; there is no unified
   range model, so cross-ecosystem tooling (e.g. dep-parser) must re-implement logic.
7. **No mutation helpers**: parse/compare only — no `Bump`, no canonicalization API, no
   "suggest next version".
8. **Missing ecosystems**: no `cargo`/crates.io, `nuget`, `conda`, `pub` (Dart), `hex` (Elixir),
   Go pseudo-versions (`v0.0.0-yyyymmddhhmmss-abcdefabcdef`).

## Proposed workstreams

### Phase 0 — Hygiene (quick wins)
- Add CI (test + vet + coverage) and a `Makefile` parity pass with dep-parser conventions.
- Replace `xerrors` with stdlib errors; bump `testify`, raise `go` directive.
- Fix cross-type `Compare` panic in `reference`: return a deterministic result or a typed error.

### Phase 1 — Core robustness
- Fuzz targets for every `Parse`/`NewVersion` (round-trip + no-panic invariants).
- Conformance corpora: semver.org vectors, node-semver fixtures, PyPA `packaging`
  test vectors, Maven `ComparableVersion` cases — run as table tests per ecosystem.
- Benchmarks for `Parse`, `Compare`, `Detect`; fix obvious allocation hotspots
  (precompile/sort registry once, avoid re-splitting in hot paths).

### Phase 2 — Unified constraint engine (advanced feature #1/#2)
- Define a common range AST (`=`, `!=`, `>`, `>=`, `<`, `<=`, hyphen ranges, `||` unions,
  wildcard/x-ranges) that all ecosystems lower into.
- Add `Satisfies`, `Subset`, `Intersect`, `Union`, `Difference`, and gap analysis.

### Phase 3 — Smarter parsing & detection (advanced feature #4)
- Coercion/normalization API (npm-style partial `1` → `1.0.0`, PEP 440 normalization).
- Confidence-scored detection from purl/package coords instead of regex guessing.

### Phase 4 — Release automation (advanced feature #5)
- `Bump(major|minor|patch|prerelease)` across ecosystems + `SuggestNext` from conventional
  commits/changelog.

### Phase 5 — Ecosystem expansion
- Add `cargo` (semver-compatible), `nuget`, `conda`, `pub`, `hex`, Go pseudo-version support.

## Suggested sequencing
Phase 0 → 1 first (they de-risk everything else), then 2 (highest user value for consumers
like dep-parser), then 3–5 by demand.
