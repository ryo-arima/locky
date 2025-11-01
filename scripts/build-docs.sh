#!/bin/bash

# Build documentation for GitHub Pages
# This script generates all documentation and organizes it for GitHub Pages hosting

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
DOCS_ROOT="$PROJECT_ROOT/docs"
DIST_DIR="$DOCS_ROOT/dist"

echo "Building documentation for GitHub Pages..."

# Clean previous build
echo "Cleaning previous build..."
rm -rf "$DIST_DIR"

# Build mdbook
echo "Building mdbook..."
cd "$DOCS_ROOT/books"
mdbook build

# Copy architecture diagrams
echo "Copying architecture diagrams..."
mkdir -p "$DIST_DIR/architecture"
if [ -f "$DOCS_ROOT/architecture/high-level-architecture.svg" ]; then
    cp "$DOCS_ROOT/architecture/high-level-architecture.svg" "$DIST_DIR/architecture/"
fi
if [ -f "$DOCS_ROOT/architecture/high-level-architecture.mmd" ]; then
    cp "$DOCS_ROOT/architecture/high-level-architecture.mmd" "$DIST_DIR/architecture/"
fi

# Copy swagger documentation
echo "Copying Swagger documentation..."
mkdir -p "$DIST_DIR/swagger"
if [ -f "$DOCS_ROOT/swagger/swagger.yaml" ]; then
    cp "$DOCS_ROOT/swagger/swagger.yaml" "$DIST_DIR/swagger/"
fi

# Create Swagger UI viewer
echo "Creating Swagger UI page..."
cat > "$DIST_DIR/swagger/index.html" << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>Locky API - Swagger UI</title>
    <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
    <style>
        html { box-sizing: border-box; overflow: -moz-scrollbars-vertical; overflow-y: scroll; }
        *, *:before, *:after { box-sizing: inherit; }
        body { margin:0; padding:0; }
    </style>
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
    <script>
        window.onload = function() {
            window.ui = SwaggerUIBundle({
                url: "./swagger.yaml",
                dom_id: '#swagger-ui',
                deepLinking: true,
                presets: [
                    SwaggerUIBundle.presets.apis,
                    SwaggerUIStandalonePreset
                ],
                plugins: [
                    SwaggerUIBundle.plugins.DownloadUrl
                ],
                layout: "StandaloneLayout"
            });
        };
    </script>
</body>
</html>
EOF

# Generate godoc HTML
echo "Generating godoc HTML..."
mkdir -p "$DIST_DIR/godoc"

# Create godoc index page
cat > "$DIST_DIR/godoc/index.html" << 'EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Locky - Go Documentation</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
            max-width: 800px;
            margin: 50px auto;
            padding: 20px;
            line-height: 1.6;
        }
        h1 { color: #333; }
        .package-list {
            background: #f5f5f5;
            padding: 20px;
            border-radius: 5px;
            margin: 20px 0;
        }
        .package-list ul {
            list-style: none;
            padding: 0;
        }
        .package-list li {
            margin: 10px 0;
        }
        .package-list a {
            color: #0366d6;
            text-decoration: none;
        }
        .package-list a:hover {
            text-decoration: underline;
        }
        .back-link {
            margin-top: 30px;
        }
    </style>
</head>
<body>
    <h1>Locky - Go Documentation</h1>
    <p>Go package documentation for Locky RBAC service.</p>
    
    <div class="package-list">
        <h2>Available on pkg.go.dev</h2>
        <ul>
            <li><a href="https://pkg.go.dev/github.com/ryo-arima/locky" target="_blank">github.com/ryo-arima/locky</a> - Main package</li>
            <li><a href="https://pkg.go.dev/github.com/ryo-arima/locky/pkg/server" target="_blank">pkg/server</a> - HTTP server implementation</li>
            <li><a href="https://pkg.go.dev/github.com/ryo-arima/locky/pkg/client" target="_blank">pkg/client</a> - CLI client implementations</li>
            <li><a href="https://pkg.go.dev/github.com/ryo-arima/locky/pkg/config" target="_blank">pkg/config</a> - Configuration management</li>
            <li><a href="https://pkg.go.dev/github.com/ryo-arima/locky/pkg/entity" target="_blank">pkg/entity</a> - Data models and DTOs</li>
        </ul>
    </div>

    <h2>Local Documentation</h2>
    <p>To view documentation locally, run:</p>
    <pre><code>godoc -http=:6060
open http://localhost:6060/pkg/github.com/ryo-arima/locky/</code></pre>

    <div class="back-link">
        <a href="../">← Back to Documentation Home</a>
    </div>
</body>
</html>
EOF

# Create .nojekyll file for GitHub Pages
echo "Creating .nojekyll file..."
touch "$DIST_DIR/.nojekyll"

# Create CNAME file if needed (optional)
# echo "yourdomain.com" > "$DIST_DIR/CNAME"

echo ""
echo "✅ Documentation build complete!"
echo "📁 Output directory: $DIST_DIR"
echo ""
echo "To preview locally:"
echo "  cd $DIST_DIR && python3 -m http.server 8000"
echo "  open http://localhost:8000"
echo ""
echo "To deploy to GitHub Pages:"
echo "  git add docs/dist"
echo "  git commit -m 'Update documentation'"
echo "  git push origin main"
