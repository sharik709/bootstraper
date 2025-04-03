#!/usr/bin/env node

const path = require('path');
const { spawn } = require('child_process');
const fs = require('fs');

// Determine the path to the binary
const localBinary = path.join(__dirname, 'bt');
let binaryPath = localBinary;

// Check if the binary exists
if (!fs.existsSync(binaryPath)) {
  // Try dist directory as fallback
  const distBinary = path.join(__dirname, '..', 'dist', 'bt');
  if (fs.existsSync(distBinary)) {
    binaryPath = distBinary;
  } else {
    console.error('Error: Bootstraper CLI binary not found. Please reinstall the package.');
    process.exit(1);
  }
}

// Execute the binary with the provided arguments
const child = spawn(binaryPath, process.argv.slice(2), { stdio: 'inherit' });

// Handle binary execution
child.on('error', (err) => {
  console.error('Error executing Bootstraper CLI:', err.message);
  process.exit(1);
});

// Forward the exit code
child.on('close', (code) => {
  process.exit(code);
});
