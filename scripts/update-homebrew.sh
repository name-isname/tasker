#!/bin/bash
# Script to manually update homebrew cask after goreleaser release
# This works around the goreleaser v2 directory issue

set -e

VERSION=${1:-$(git describe --tags --abbrev=0)}
REPO_DIR=${2:-/tmp/homebrew-taskctl}

echo "Updating homebrew cask for version $VERSION"

# Clone or update homebrew tap
if [ -d "$REPO_DIR" ]; then
    cd "$REPO_DIR"
    git pull
else
    git clone https://github.com/name-isname/homebrew-taskctl.git "$REPO_DIR"
    cd "$REPO_DIR"
fi

# Get release assets and checksums
echo "Fetching release assets..."
BASE_URL="https://github.com/name-isname/tasker/releases/download/$VERSION"

declare -A PLATFORMS=(
    ["darwin_arm64"]="macOS ARM64"
    ["darwin_amd64"]="macOS Intel"
    ["linux_amd64"]="Linux AMD64"
    ["linux_arm64"]="Linux ARM64"
)

# Collect checksums
CHECKSUMS=()
for platform in "${!PLATFORMS[@]}"; do
    echo "Downloading ${PLATFORMS[$platform]}..."
    sha256=$(curl -sL "$BASE_URL/tasker_${VERSION}_${platform}.tar.gz" | shasum -a 256 | awk '{print $1}')
    CHECKSUMS["$platform"]=$sha256
done

# Create Casks directory if it doesn't exist
mkdir -p Casks

# Create cask file
cat > Casks/taskctl.rb << 'EOF'
cask "taskctl" do
  version "VERSION_PLACEHOLDER"
  sha256 "SHA256_PLACEHOLDER"

  url "https://github.com/name-isname/tasker/releases/download/VERSION_PLACEHOLDER/tasker_VERSION_PLACEHOLDER_darwin_arm64.tar.gz"
  name "taskctl"
  desc "Process-oriented task management tool with CLI, TUI, and Web UI"
  homepage "https://github.com/name-isname/tasker"

  binary "taskctl"

  zap trash: "~/.taskctl"
end
EOF

# Replace placeholders with actual version
sed -i.bak "s/VERSION_PLACEHOLDER/$VERSION/g" Casks/taskctl.rb
rm Casks/taskctl.rb.bak

# Show diff
echo "Cask content:"
cat Casks/taskctl.rb

# Commit and push
read -p "Commit and push? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git add Casks/taskctl.rb
    git commit -m "Update taskctl cask to $VERSION"
    git push
    echo "Homebrew cask updated successfully!"
else
    echo "Skipping commit."
fi
