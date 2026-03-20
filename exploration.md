## Exploration: Release v0.2.0 Preparation

### Current State
The project has completed the main P1 features for v0.2.0:
- A1: Number keybindings (1-4) for direct tab switching
- A2: Item counts in tab labels
- A3: Alphabetical list sorting

However, several other items planned for v0.2.0 in `ROADMAP.md` remain unimplemented:
- A4: Hidden agent toggle (`h` key) [P2]
- A5: Copy detail to clipboard (`y` key) [P2]
- A6: Configurable list/detail split ratio [P3]
- A7: Full-screen detail view (`Enter` key) [P2]
- A8: Scroll position indicator in status bar [P3]
- A9: Config file path in status bar [P3]
- E1: Unit tests for api and merge packages [P1]
- E2: GitHub Actions CI pipeline [P1]
- E3: golangci-lint integration [P2]

The `README.md` was recently updated to document A1, A2, and A3. `CHANGELOG.md` only has the `[0.1.0]` release so far.

### Affected Areas
- `ROADMAP.md` — Needs to be updated to push unimplemented features (A4-A9, E1-E3) from v0.2.0 to v0.3.0 or a later release.
- `CHANGELOG.md` — Needs a new `[0.2.0]` block with the completed features (A1, A2, A3).
- `README.md` — Might need version references bumped (if any exist, though currently mostly feature descriptions).

### Approaches
1. **Delay Release until all v0.2.0 items are done**
   - Pros: Delivers exactly what was promised.
   - Cons: Holds back the useful features already implemented (A1, A2, A3).
   - Effort: High (requires implementing CI, tests, and more UI features).

2. **Release what we have and shift the rest**
   - Pros: Ships value immediately. Accurately reflects reality.
   - Cons: Roadmap promise is "broken" but updated dynamically.
   - Effort: Low (just documentation updates).

### Recommendation
**Approach 2 (Release what we have and shift the rest).** Since the user explicitly asked to "sacar release v0.2.0", we should adjust the roadmap to reflect reality. Move the uncompleted items to v0.3.0, document A1, A2, and A3 in `CHANGELOG.md` under `[0.2.0]`, and prepare the commit/push.

### Risks
- CI/CD (E2) and tests (E1) are being delayed, which could reduce code quality confidence in the short term.
- Users expecting the hidden agent toggle (A4) or copy features (A5) will have to wait for v0.3.0.

### Ready for Proposal
Yes. The orchestrator can proceed to the `sdd-propose` phase or directly to `sdd-apply` to update the documentation files (`ROADMAP.md`, `CHANGELOG.md`), commit, tag, and push the release.