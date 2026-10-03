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
	go test ./...

# The browser suite runs against the static export, so it needs a fresh `out/`
# first. CI shards this; locally it is one command. Nothing else runs it, which
# is how a renamed button and a moved page reached main green on the unit suite.
e2e:
	cd web/next-app && rm -rf out && npx next build
	cd web/next-app && npx playwright test

validate:
	go run scripts/validate_graph.go
	python3 scripts/audit_lessons.py
	python3 scripts/audit_solutions.py
	python3 scripts/counts.py --check
	go test ./internal/generator/ -run TestLearnerDomainsHaveSchemas

# Strict LaTeX check: runs every math span in the lesson corpus through Compute
# Engine and reports the malformed ones. Needs node (for the CE dependency) and
# the Go toolchain (for latexdump). Exits 1 while defects remain, so it is not in
# `validate` yet -- see the CI step for why.
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
