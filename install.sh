#!/bin/sh 
# Adapted from the Deno installer: Copyright 2019 the Deno authors. All rights reserved. MIT license.
# Ref: https://github.com/denoland/deno_install

# Exit immediately if any command exits with a non-zero status
set -e

os=$(uname -s)
arch=$(uname -m)

echo "Detected operating system: $os"
echo "Detected architecture: $arch"

if [ "$arch" = "aarch64" ]; then 	
	arch="arm64"
fi  

if [ $# -eq 0 ]; then 	
	download_uri="https://github.com/okt-limonikas/bublik-log-viewer/releases/latest/download/bublik-log-viewer_${os}_${arch}.tar.gz"
else 	
	download_uri="https://github.com/okt-limonikas/bublik-log-viewer/releases/download/${1}/bublik-log-viewer_${os}_${arch}.tar.gz"
fi  

log_install="${LOG_INSTALL:-${HOME}/.local}"
bin_dir="${log_install}/bin"
exe="${bin_dir}/bublik-log-viewer"

echo "Bin directory is: ${bin_dir}"

if [ ! -d "${bin_dir}" ]; then 	
	mkdir -p "${bin_dir}"
fi  

echo "Downloading Bublik log viewer from: $download_uri"
curl --silent --show-error --location --fail --location --output "${bin_dir}/bublik-log-viewer.tar.gz" "$download_uri"

echo "Extracting Bublik log viewer..."
tar -xzf "${bin_dir}/bublik-log-viewer.tar.gz" -C "${bin_dir}/"

chmod +x "${exe}"

echo "Cleaning up downloaded archive ${bin_dir}/bublik-log-viewer.tar.gz"
rm "${bin_dir}/bublik-log-viewer.tar.gz"

echo "Bublik log viewer was installed successfully to ${exe}"

if command -v bublik-log-viewer >/dev/null; then 	
	echo "Run 'bublik-log-viewer --help' to get started"
fi  
