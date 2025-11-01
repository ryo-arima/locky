#!/bin/bash

# Script to publish package to pkg.go.dev
# This creates a version tag and pushes it to trigger pkg.go.dev indexing

set -e

VERSION=${1:-v0.1.0}

echo "Publishing Locky to pkg.go.dev..."
echo "Version: $VERSION"
echo ""

# Check if tag already exists
if git tag | grep -q "^${VERSION}$"; then
    echo "Error: Tag $VERSION already exists"
    echo "Available tags:"
    git tag
    exit 1
fi

# Ensure we're on main branch
CURRENT_BRANCH=$(git branch --show-current)
if [ "$CURRENT_BRANCH" != "main" ]; then
    echo "Warning: You are on branch '$CURRENT_BRANCH', not 'main'"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Check for uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo "Error: You have uncommitted changes"
    echo "Please commit or stash your changes first"
    exit 1
fi

# Create annotated tag
echo "Creating tag $VERSION..."
git tag -a "$VERSION" -m "Release $VERSION

- Initial public release
- Complete RBAC implementation with Casbin
- JWT authentication
- User, Group, Member, Role management
- Multi-tier API (Public, Internal, Private)
- Redis caching
- MySQL/TiDB support
- CLI clients (Admin, App, Anonymous)
- Comprehensive documentation
"

echo ""
echo "Tag created successfully!"
echo ""
echo "Pushing tag to GitHub..."
git push origin "$VERSION"

echo ""
echo "✅ Tag pushed successfully!"
echo ""
echo "pkg.go.dev will index your package automatically within a few minutes."
echo ""
echo "You can check the status at:"
echo "  https://pkg.go.dev/github.com/ryo-arima/locky@$VERSION"
echo ""
echo "Or request immediate indexing by visiting:"
echo "  https://proxy.golang.org/github.com/ryo-arima/locky/@v/$VERSION.info"
echo ""
echo "Once indexed, your package will be available at:"
echo "  https://pkg.go.dev/github.com/ryo-arima/locky"
echo "  https://pkg.go.dev/github.com/ryo-arima/locky/pkg/server"
echo "  https://pkg.go.dev/github.com/ryo-arima/locky/pkg/client"
echo "  https://pkg.go.dev/github.com/ryo-arima/locky/pkg/config"
echo "  https://pkg.go.dev/github.com/ryo-arima/locky/pkg/entity"
echo ""
