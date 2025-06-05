#!/bin/bash

set -e  # Exit on any error

echo "Installing PostgreSQL client tools versions 13-17 for Linux (Debian/Ubuntu)..."
echo

# Check if running on supported system
if ! command -v apt-get &> /dev/null; then
    echo "Error: This script requires apt-get (Debian/Ubuntu-like system)"
    exit 1
fi

# Check if running as root or with sudo
if [[ $EUID -eq 0 ]]; then
    SUDO=""
else
    SUDO="sudo"
    echo "This script requires sudo privileges to install packages."
fi

# Create postgresql directory
mkdir -p postgresql

# Get absolute path
POSTGRES_DIR="$(pwd)/postgresql"

echo "Installing PostgreSQL client tools to: $POSTGRES_DIR"
echo

# Add PostgreSQL official APT repository
echo "Adding PostgreSQL official APT repository..."
$SUDO apt-get update -qq
$SUDO apt-get install -y wget ca-certificates

# Add GPG key
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | $SUDO apt-key add -

# Add repository
echo "deb http://apt.postgresql.org/pub/repos/apt/ $(lsb_release -cs)-pgdg main" | $SUDO tee /etc/apt/sources.list.d/pgdg.list

# Update package list
echo "Updating package list..."
$SUDO apt-get update -qq

# Install client tools for each version
versions="13 14 15 16 17"

for version in $versions; do
    echo "Installing PostgreSQL $version client tools..."
    
    # Install client tools only
    $SUDO apt-get install -y postgresql-client-$version
    
    # Create version-specific directory and symlinks
    version_dir="$POSTGRES_DIR/postgresql-$version"
    mkdir -p "$version_dir/bin"
    
    # Create symlinks to the installed binaries
    if [ -f "/usr/bin/pg_dump" ]; then
        # If multiple versions, binaries are usually named with version suffix
        if [ -f "/usr/bin/pg_dump-$version" ]; then
            ln -sf "/usr/bin/pg_dump-$version" "$version_dir/bin/pg_dump"
            ln -sf "/usr/bin/pg_dumpall-$version" "$version_dir/bin/pg_dumpall"
            ln -sf "/usr/bin/psql-$version" "$version_dir/bin/psql"
            ln -sf "/usr/bin/pg_restore-$version" "$version_dir/bin/pg_restore"
            ln -sf "/usr/bin/createdb-$version" "$version_dir/bin/createdb"
            ln -sf "/usr/bin/dropdb-$version" "$version_dir/bin/dropdb"
        else
            # Fallback to non-versioned names (latest version)
            ln -sf "/usr/bin/pg_dump" "$version_dir/bin/pg_dump"
            ln -sf "/usr/bin/pg_dumpall" "$version_dir/bin/pg_dumpall"
            ln -sf "/usr/bin/psql" "$version_dir/bin/psql"
            ln -sf "/usr/bin/pg_restore" "$version_dir/bin/pg_restore"
            ln -sf "/usr/bin/createdb" "$version_dir/bin/createdb"
            ln -sf "/usr/bin/dropdb" "$version_dir/bin/dropdb"
        fi
        
        echo "PostgreSQL $version client tools installed successfully"
    else
        echo "Warning: PostgreSQL $version client tools may not have installed correctly"
    fi
    echo
done

echo "Installation completed!"
echo "PostgreSQL client tools are available in: $POSTGRES_DIR"
echo

# List installed versions
echo "Installed PostgreSQL client versions:"
for version in $versions; do
    version_dir="$POSTGRES_DIR/postgresql-$version"
    if [ -f "$version_dir/bin/pg_dump" ]; then
        echo "  postgresql-$version: $version_dir/bin/"
    fi
done

echo
echo "Usage example:"
echo "  $POSTGRES_DIR/postgresql-15/bin/pg_dump --version" 