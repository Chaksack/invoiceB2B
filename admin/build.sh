#!/bin/bash

# Build script for Next.js admin dashboard

echo "Building Invoice B2B Admin Dashboard..."

# Install dependencies
echo "Installing dependencies..."
npm install

# Build the Next.js application
echo "Building Next.js application..."
npm run build

# Export static files
echo "Exporting static files..."
npm run export

echo "Build completed successfully!"
echo "Static files are available in the 'out' directory"