run:
	go run ./cmd/mathua

build:
	go build -o bin/mathua ./cmd/mathua

build-all:
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build -o dist/mathua-windows-amd64.exe ./cmd/mathua
	GOOS=linux GOARCH=amd64 go build -o dist/mathua-linux-amd64 ./cmd/mathua
	GOOS=darwin GOARCH=amd64 go build -o dist/mathua-darwin-amd64 ./cmd/mathua
	GOOS=darwin GOARCH=arm64 go build -o dist/mathua-darwin-arm64 ./cmd/mathua

test:
	python3 scripts/check_duplicate_tests.py
	@unformatted="$$(gofmt -l $$(git ls-files '*.go'))"; \
	if [ -n "$$unformatted" ]; then echo "gofmt needed on:"; echo "$$unformatted"; exit 1; fi
	go vet ./...
	go test ./...

# The browser suite runs against the static export, so it needs a fresh `out/`
# first. CI shards this; locally it is one command. Nothing else runs it, which
# is how a renamed button and a moved page reached main green on the unit suite.
e2e:
	cd web/next-app && rm -rf out && npx next build
	cd web/next-app && npx playwright test

# Prove the browser suite does not depend on a real API. The bundle is normally
# built against the deployed API (web/next-app/.env.production), so an endpoint no
# stub claims is answered by production, and a page that waits for that answer can
# pass on a fast reply and hang on a slow one. This build points every request at
# 192.0.2.1, a documentation address that never answers, so a stub that misses its
# URL is a failure here rather than a flake in CI. See CONTRIBUTING.md.
e2e-hermetic:
	cd web/next-app && rm -rf out && NEXT_PUBLIC_API_URL=http://192.0.2.1:9099 npx next build
	cd web/next-app && npx playwright test

# The diagnosis suite measures the answer-to-submit timeline and prints what it
# finds instead of asserting, so a scenario that breaks still reports green. It is
# not part of `e2e` and not part of CI: it cost 273.6s of e2e test time across 27
# executions on 2026-10-23, and one run printed four FAILURE lines while the shard
# passed. This target is how it runs on all three profiles when someone is
# investigating a phone that will not submit.
# The harness runs under the same network guard as the suite, so it needs the hermetic build:
# a production URL would have every unstubbed call aborted and reported as an escape.
e2e-diagnosis:
	cd web/next-app && rm -rf out && NEXT_PUBLIC_API_URL=http://192.0.2.1:9099 npx next build
	cd web/next-app && PW_DIAGNOSIS=1 npx playwright test e2e/answer-diagnosis.spec.ts

validate:
	python3 scripts/check_duplicate_tests.py
	go run scripts/validate_graph.go
	go run ./cmd/validate-mastery-units
	python3 scripts/audit_lessons.py
	python3 scripts/audit_solutions.py
	python3 scripts/counts.py --check
	python3 scripts/check_docs.py
	python3 scripts/check_ste.py
	node web/next-app/scripts/build-docs.mjs --check
	go test ./internal/generator/ -run TestLearnerDomainsHaveSchemas

# Strict LaTeX check: runs every math span in the lesson corpus through Compute
# Engine and reports the malformed ones. Needs node (for the CE dependency) and
# the Go toolchain (for latexdump). Exits 1 while defects remain, so it is not in
# `validate` yet -- see the CI step for why.
# The mastery-unit partition is a generated draft (ADR-046 D3). The band width is a parameter,
# not policy: 2, 3 and 4 are all candidates and the width decides how often a learner is
# assessed. Regenerating at a different width needs no code change.
mastery-units:
	go run ./cmd/masteryunits -band 3

latex-normalize:
	cd web/next-app && node scripts/normalize-latex.mjs

lint-go:
	go run github.com/curtbushko/go-ai-lint/cmd/go-ai-lint@v1.0.1-0.20260620203811-c6ce4ee5624f ./...

lint-slop:
	npx --yes aislop@0.16.1 scan

fuzz:
	go test ./internal/generator/... -run TestFuzz -count 1000

tidy:
	go mod tidy

serve:
	go run ./cmd/mathua --serve --port 8080
