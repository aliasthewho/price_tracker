#!/bin/bash

# Exit on error
set -e

echo "Setting up Git hooks..."

# Create the pre-commit hook
echo '#!/bin/sh

# Run pre-commit
echo "Running pre-commit checks..."
make pre-commit

# Capture the exit code
EXIT_CODE=$?

# If the checks failed, show a message and abort the commit
if [ $EXIT_CODE -ne 0 ]; then
    echo "\n❌ Pre-commit checks failed. Please fix the issues and try again."
    echo "To skip these checks, use: git commit --no-verify"
    exit 1
fi

exit 0
' > .git/hooks/pre-commit

# Make the hook executable
chmod +x .git/hooks/pre-commit

echo "✅ Git hooks set up successfully!"
