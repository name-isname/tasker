#!/bin/bash
# Script to manually update homebrew formula after goreleaser release
# This works around the goreleaser v2 directory issue

set -e

VERSION=${1:-$(git describe --tags --abbrev=0)}
REPO_DIR=${2:-/tmp/homebrew-taskctl}

echo "Updating homebrew formula for version $VERSION"

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

# Get commit SHA
COMMIT_SHA=$(git ls-remote https://github.com/name-isname/tasker.git "refs/tags/$VERSION" | awk '{print $1}')

# Create formula file
cat > Formula/taskctl.rb << EOF
# typed: strict
# frozen_string_literal: true

class Taskctl < Formula
  desc "Process-oriented task management tool with CLI, TUI, and Web UI"
  homepage "https://github.com/name-isname/tasker"
  url "https://github.com/name-isname/tasker.git",
      tag:      "$VERSION",
      revision: "$COMMIT_SHA"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/name-isname/tasker/releases/download/$VERSION/tasker_${VERSION}_darwin_arm64.tar.gz"
      sha256 "${CHECKSUMS[darwin_arm64]}"
    end

    on_intel do
      url "https://github.com/name-isname/tasker/releases/download/$VERSION/tasker_${VERSION}_darwin_amd64.tar.gz"
      sha256 "${CHECKSUMS[darwin_amd64]}"
    end
  end

  on_linux do
    on_intel do
      url "https://github.com/name-isname/tasker/releases/download/$VERSION/tasker_${VERSION}_linux_amd64.tar.gz"
      sha256 "${CHECKSUMS[linux_amd64]}"
    end

    on_arm do
      url "https://github.com/name-isname/tasker/releases/download/$VERSION/tasker_${VERSION}_linux_arm64.tar.gz"
      sha256 "${CHECKSUMS[linux_arm64]}"
    end
  end

  def install
    bin.install "taskctl"
  end

  test do
    system bin/"taskctl", "version"
  end
end
EOF

# Show diff
echo "Formula content:"
git diff Formula/taskctl.rb | head -50

# Commit and push
read -p "Commit and push? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git add Formula/taskctl.rb
    git commit -m "Update taskctl to $VERSION"
    git push
    echo "Homebrew formula updated successfully!"
else
    echo "Skipping commit."
fi
