#!/bin/bash

# Script to check test coverage for all packages

set -e

echo "========================================="
echo "Running Unit Tests with Coverage"
echo "========================================="
echo ""

# Get all pkg subdirectories that have tests in test/unit/pkg
packages=()
for dir in test/unit/pkg/*/; do
    if [ -d "$dir" ]; then
        pkg_name=$(basename "$dir")
        # Check if there are test files
        if ls "$dir"*_test.go >/dev/null 2>&1; then
            packages+=("$pkg_name")
        fi
    fi
done

# Also check subdirectories (e.g., entity/model, entity/request, etc.)
for dir in test/unit/pkg/*/*/; do
    if [ -d "$dir" ]; then
        parent=$(basename $(dirname "$dir"))
        pkg_name="$parent/$(basename "$dir")"
        # Check if there are test files
        if ls "$dir"*_test.go >/dev/null 2>&1; then
            packages+=("$pkg_name")
        fi
    fi
done

total_coverage=0
package_count=0

for pkg in "${packages[@]}"; do
    echo "Testing pkg/$pkg..."
    result=$(go test -cover -coverpkg=./pkg/$pkg/... ./test/unit/pkg/$pkg/... 2>&1 | grep "coverage:" || echo "coverage: 0.0%")
    coverage=$(echo "$result" | grep -oE '[0-9]+\.[0-9]+%' | head -1)
    
    if [ -n "$coverage" ]; then
        echo "  ✓ Coverage: $coverage"
        # Extract numeric value for averaging
        numeric_coverage=$(echo "$coverage" | grep -oE '[0-9]+\.[0-9]+')
        total_coverage=$(echo "$total_coverage + $numeric_coverage" | bc)
        package_count=$((package_count + 1))
    else
        echo "  ✗ No coverage data"
    fi
    echo ""
done

echo "========================================="
echo "Summary"
echo "========================================="

if [ $package_count -gt 0 ]; then
    avg_coverage=$(echo "scale=1; $total_coverage / $package_count" | bc)
    echo "Packages tested: $package_count"
    echo "Average coverage: ${avg_coverage}%"
    
    if (( $(echo "$avg_coverage >= 100.0" | bc -l) )); then
        echo "Status: ✓ Target achieved! (100%)"
    elif (( $(echo "$avg_coverage >= 80.0" | bc -l) )); then
        echo "Status: ⚠ Good progress (>80%)"
    else
        echo "Status: ✗ Needs improvement (<80%)"
    fi
else
    echo "No packages tested"
fi

echo ""
echo "To generate detailed HTML coverage report:"
echo "  go test -coverpkg=./pkg/... ./test/unit/pkg/... -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html"
