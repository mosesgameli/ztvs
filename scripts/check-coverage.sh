#!/bin/bash
set -e

# Run tests and generate coverage profile
echo "Running unit tests and generating coverage profile..."
# Exclude /test/ directory and use all packages
go test -count=1 -coverprofile=coverage.out ./...

# Extract total coverage percentage
COVERAGE=$(go tool cover -func=coverage.out | grep total | grep -Eo '[0-9]+\.[0-9]+')
echo "Current coverage: $COVERAGE%"

# Minimum coverage required
THRESHOLD=85.0

# Compare coverage using bc for decimal comparison
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
  echo "❌ Error: Code coverage is below $THRESHOLD% ($COVERAGE%)"
  exit 1
fi

echo "✅ Code coverage is sufficient ($COVERAGE%)"
exit 0
