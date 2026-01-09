# Copilot Repository Instructions

Purpose
- Help contributors and GitHub Copilot produce code that follows our conventions, uses TDD, and maintains high quality.

## Conventions and Formatting
- Adhere to the naming and formatting standards of the language you are using.
- Respect existing repository configurations:
  - Formatters/linters: Prettier, ESLint, Stylelint, Black, Flake8, isort, gofmt, rustfmt, ktlint, Checkstyle, etc.
  - Project configs: .editorconfig, pyproject.toml, package.json, Makefile, tox.ini, mvn/gradle, cargo, dotnet, etc.
- Prefer idiomatic APIs, patterns, and file layouts already present in this repository.
- Naming guidelines (adapt to language norms):
  - Functions/methods and variables: use the language’s conventional casing (e.g., camelCase in JS/TS/Java/Kotlin, snake_case in Python/Ruby, PascalCase for C# methods where applicable).
  - Types/classes: PascalCase.
  - Constants: use the language’s conventional form (e.g., UPPER_SNAKE_CASE in many languages).
- Keep code cohesive and clear; favor readability over cleverness.
- Do not commit generated artifacts, secrets, or environment-specific files. Respect .gitignore and .gitattributes.

## Test-Driven Development (TDD)
- Always use a TDD approach:
  1. Write or update tests that describe the desired behavior; ensure they fail first.
  2. Implement the minimal code to make the tests pass.
  3. Refactor while keeping tests green.
- Coverage expectations:
  - Make sure all new features are covered by tests.
  - Do not forget to test edge cases:
    - null/undefined/None, empty strings/collections
    - zero/negative numbers, boundary and off-by-one values
    - large inputs, performance constraints
    - non-ASCII/Unicode, locales, time zones/DST, leap years
    - precision/overflow, floating-point quirks
    - IO/network errors, retries, timeouts
    - concurrency/race conditions (where applicable)
- Keep tests deterministic, isolated, and fast:
  - Avoid real network/file system where possible; prefer fakes/mocks.
  - Control randomness and time with fixed seeds/clocks.
- Place and name tests according to language/repo conventions:
  - JS/TS: __tests__/ or alongside code with .test.{js,ts}; Jest/Vitest if configured.
  - Python: tests/ with files named test_*.py; pytest if present.
  - Go: *_test.go using testing package.
  - Java: src/test/java with JUnit Jupiter if present.
  - C#: dedicated test project; *Tests.cs with xUnit/NUnit per setup.
  - Ruby: spec/ (RSpec) or test/ (Minitest) per repo.
  - Rust: tests/ integration tests and/or #[cfg(test)] modules.
- When fixing bugs, add regression tests that fail without the fix.

## After Any Code Change
- Run the full test suite and static checks/linters/formatters.
- Fix all failing tests and issues before finishing.
- Use the repository’s standard commands (e.g., package.json scripts, Makefile targets, tox, mvn/gradle, dotnet test, go test, cargo test).
- Update documentation, comments, and changelogs when behavior or public APIs change.

## Dependencies and Configuration
- Reuse existing dependencies when possible; avoid adding new ones without justification.
- If adding/upgrading dependencies:
  - Update lockfiles (e.g., package-lock.json, pnpm-lock.yaml, poetry.lock, Pipfile.lock).
  - Update CI config and docs if needed.
- Follow existing environment, tooling, and CI/CD conventions.

## Commits
- When a feature/task/fix is finished, commit all changes only after all tests pass and quality checks succeed.
- Use a meaningful commit message that clearly describes what was done and the overall goal.
  - Keep the subject concise; add a body explaining context, rationale, and any breaking changes.
  - Reference related issues/PRs when relevant.
  - Conventional Commits style is encouraged:
    - feat: add user profile API to support editing avatars
    - fix: handle empty input in tokenizer to prevent crash
    - refactor: extract date parsing into utility with tests
    - test: add edge cases for DST transitions
    - docs: update README with setup instructions
    - chore: bump eslint and refresh lockfile

## Guidance for GitHub Copilot
- Generate tests alongside any new code (functions, classes, endpoints, scripts).
- Match the repository’s existing file layout, frameworks, and idioms.
- Conform to formatting, linting, and type-checking rules; include necessary imports and stubs.
- Prefer standard libraries and established project utilities over ad-hoc implementations.
- If conventions are unclear, infer from existing code; otherwise, propose a brief comment with the chosen approach.
- Do not introduce breaking changes silently; if unavoidable, document and add migration notes/tests.

## Quick Checklist (per change)
- [ ] Follow language naming and formatting standards.
- [ ] Write failing tests first (cover happy path + edge cases).
- [ ] Implement the minimal code to pass tests.
- [ ] Refactor safely; keep tests green.
- [ ] Run formatters, linters, type checks, and the full test suite.
- [ ] Update docs/comments as needed.
- [ ] Commit with a meaningful message summarizing what and why.