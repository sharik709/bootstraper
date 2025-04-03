#!/usr/bin/env node

const fs = require('fs');
const path = require('path');
const os = require('os');
const https = require('https');
const { execSync } = require('child_process');

const version = process.env.npm_package_version || '0.2.0';
const platform = os.platform();
const arch = os.arch();

// Map Node.js arch to Go arch
const archMap = {
  'x64': 'amd64',
  'arm64': 'arm64',
  'ia32': '386'
};

// Map Node.js platform to Go OS
const platformMap = {
  'darwin': 'darwin',
  'linux': 'linux',
  'win32': 'windows'
};

const goArch = archMap[arch] || 'amd64';
const goOS = platformMap[platform] || 'linux';

const isWindows = platform === 'win32';
const extension = isWindows ? '.exe' : '';
const binName = `bt${extension}`;

const binDir = path.join(__dirname, '..', 'bin');
const binaryPath = path.join(binDir, binName);

// Make sure the bin directory exists
if (!fs.existsSync(binDir)) {
  fs.mkdirSync(binDir, { recursive: true });
}

// Function to download prebuilt binary
function downloadBinary() {
  const binaryFileName = `bt-${goOS}-${goArch}${extension}`;
  const url = `https://github.com/sharik709/bootstraper/releases/download/v${version}/${binaryFileName}`;

  console.log(`Downloading Bootstraper CLI v${version} for ${goOS}-${goArch}...`);

  return new Promise((resolve, reject) => {
    https.get(url, (response) => {
      if (response.statusCode === 302 || response.statusCode === 301) {
        // Handle redirects
        https.get(response.headers.location, (redirectResponse) => {
          if (redirectResponse.statusCode !== 200) {
            reject(new Error(`Failed to download binary: ${redirectResponse.statusCode}`));
            return;
          }

          const file = fs.createWriteStream(binaryPath);
          redirectResponse.pipe(file);

          file.on('finish', () => {
            file.close();
            fs.chmodSync(binaryPath, 0o755); // Make binary executable
            resolve();
          });

          file.on('error', (err) => {
            fs.unlinkSync(binaryPath);
            reject(err);
          });
        }).on('error', reject);
      } else if (response.statusCode === 200) {
        const file = fs.createWriteStream(binaryPath);
        response.pipe(file);

        file.on('finish', () => {
          file.close();
          fs.chmodSync(binaryPath, 0o755); // Make binary executable
          resolve();
        });

        file.on('error', (err) => {
          fs.unlinkSync(binaryPath);
          reject(err);
        });
      } else {
        reject(new Error(`Failed to download binary: ${response.statusCode}`));
      }
    }).on('error', (err) => {
      reject(err);
    });
  });
}

// Alternative: Try to build from source if Go is available
function buildFromSource() {
  console.log('Attempting to build from source...');

  try {
    // Check if Go is installed
    execSync('go version', { stdio: 'ignore' });

    // Get package directory
    const packageDir = path.join(__dirname, '..');

    // Build binary
    execSync('go build -o ' + binaryPath, {
      cwd: packageDir,
      stdio: 'inherit'
    });

    fs.chmodSync(binaryPath, 0o755); // Make binary executable
    console.log('Successfully built Bootstraper CLI from source!');
    return true;
  } catch (err) {
    console.log('Failed to build from source:', err.message);
    return false;
  }
}

// Copy pre-built binary if it exists in dist directory
function copyPreBuiltBinary() {
  const distBinaryPath = path.join(__dirname, '..', 'dist', 'bt');

  if (fs.existsSync(distBinaryPath)) {
    console.log('Using pre-built binary from dist directory...');
    fs.copyFileSync(distBinaryPath, binaryPath);
    fs.chmodSync(binaryPath, 0o755); // Make binary executable
    return true;
  }

  return false;
}

// Main installation process
async function install() {
  try {
    // First check if we already have a pre-built binary in the dist folder
    if (copyPreBuiltBinary()) {
      console.log('Successfully installed Bootstraper CLI!');
      return;
    }

    // Try to download binary from GitHub releases
    try {
      await downloadBinary();
      console.log('Successfully downloaded Bootstraper CLI!');
    } catch (err) {
      console.log('Failed to download binary:', err.message);

      // Try to build from source as a fallback
      if (!buildFromSource()) {
        console.error('Installation failed! Please install manually from https://github.com/sharik709/bootstraper');
        process.exit(1);
      }
    }

    console.log('Bootstraper CLI installed successfully!');
    console.log('You can now use the "bt" command to bootstrap projects.');
  } catch (err) {
    console.error('Error during installation:', err);
    process.exit(1);
  }
}

// Run the installation
install().catch(console.error);
