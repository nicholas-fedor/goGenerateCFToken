#!/bin/sh
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Install, update, or uninstall gogeneratecftoken from GitHub Releases.
#
# Usage:
#   tmp=$(mktemp)
#   curl -sSfL https://raw.githubusercontent.com/nicholas-fedor/goGenerateCFToken/main/scripts/install.sh -o "$tmp" && sh "$tmp"
#   sh "$tmp" update
#   sh "$tmp" uninstall
#
# Environment:
#   VERSION   Release tag (v1.2.0 or 1.2.0). Default: latest.
#   PREFIX    Directory for archive installs. Default: $HOME/go/bin, or the
#             prefix stored from a previous archive install.
#   METHOD    auto (default), package, or archive.
#             auto installs a native package on Linux when sudo/root is
#             available, otherwise extracts the release archive.

set -eu

REPO="nicholas-fedor/goGenerateCFToken"
BINARY="gogeneratecftoken"
GITHUB="https://github.com/${REPO}"
DEFAULT_PREFIX="${HOME}/go/bin"
STATE_DIR="${XDG_STATE_HOME:-${HOME}/.local/state}/${BINARY}"
STATE_PREFIX_FILE="${STATE_DIR}/prefix"
COSIGN_IDENTITY_REGEXP="^https://github.com/${REPO}/"
COSIGN_OIDC_ISSUER="https://token.actions.githubusercontent.com"

VERSION="${VERSION:-}"
METHOD="${METHOD:-auto}"

PREFIX_SET=0
if [ "${PREFIX+x}" = "x" ]; then
	PREFIX_SET=1
fi

info() {
	printf '%s\n' "$*" >&2
}

die() {
	printf 'error: %s\n' "$*" >&2
	exit 1
}

need_cmd() {
	command -v "$1" >/dev/null 2>&1 || die "required command not found: $1"
}

download() {
	url=$1
	dest=$2

	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "$url" -o "$dest"
	elif command -v wget >/dev/null 2>&1; then
		wget -qO "$dest" "$url"
	else
		die "curl or wget is required"
	fi
}

try_download() {
	url=$1
	dest=$2

	if command -v curl >/dev/null 2>&1; then
		curl -fsSL "$url" -o "$dest" || return 1
	elif command -v wget >/dev/null 2>&1; then
		wget -qO "$dest" "$url" || return 1
	else
		die "curl or wget is required"
	fi
}

resolve_prefix() {
	if [ "$PREFIX_SET" -eq 1 ]; then
		return 0
	fi

	if [ -f "$STATE_PREFIX_FILE" ]; then
		PREFIX=$(tr -d '\r\n' <"$STATE_PREFIX_FILE")
		if [ -n "$PREFIX" ]; then
			return 0
		fi
	fi

	PREFIX="$DEFAULT_PREFIX"
}

save_prefix() {
	mkdir -p "$STATE_DIR"
	printf '%s\n' "$PREFIX" >"$STATE_PREFIX_FILE"
}

clear_prefix_state() {
	rm -f "$STATE_PREFIX_FILE"
}

prefix_on_path() {
	case ":${PATH}:" in
	*":${PREFIX}:"*) return 0 ;;
	*) return 1 ;;
	esac
}

resolve_version() {
	if [ -n "$VERSION" ]; then
		case "$VERSION" in
		v*) printf '%s\n' "$VERSION" ;;
		*) printf 'v%s\n' "$VERSION" ;;
		esac
		return
	fi

	need_cmd curl

	tag=$(curl -fsSLI "${GITHUB}/releases/latest" | tr -d '\r' | awk -F/ 'tolower($0) ~ /^location:/ {print $NF; exit}')
	[ -n "$tag" ] || die "could not determine the latest release tag"
	printf '%s\n' "$tag"
}

os_name() {
	uname_s=$(uname -s)
	case "$uname_s" in
	Linux) printf 'linux\n' ;;
	Darwin) printf 'macOS\n' ;;
	MINGW* | MSYS* | CYGWIN*) die "Windows is not supported by this script; download the .zip from ${GITHUB}/releases" ;;
	*) die "unsupported OS: ${uname_s}" ;;
	esac
}

arch_name() {
	uname_m=$(uname -m)
	case "$uname_m" in
	x86_64 | amd64) printf 'amd64\n' ;;
	i386 | i486 | i586 | i686) printf 'i386\n' ;;
	aarch64 | arm64) printf 'arm64v8\n' ;;
	armv7l | armv6l | armv5l | arm) printf 'armhf\n' ;;
	riscv64) printf 'riscv64\n' ;;
	*) die "unsupported architecture: ${uname_m}" ;;
	esac
}

pkg_kind() {
	[ -r /etc/os-release ] || return 1
	# shellcheck disable=SC1091
	. /etc/os-release

	id=$(printf '%s' "${ID:-}" | tr '[:upper:]' '[:lower:]')
	like=$(printf '%s' "${ID_LIKE:-}" | tr '[:upper:]' '[:lower:]')

	case "$id" in
	debian | ubuntu | linuxmint | pop | raspbian | elementary) printf 'deb\n' && return 0 ;;
	fedora | rhel | centos | rocky | almalinux | ol | amzn | sles | opensuse*) printf 'rpm\n' && return 0 ;;
	alpine) printf 'apk\n' && return 0 ;;
	arch | manjaro | endeavouros | garuda) printf 'arch\n' && return 0 ;;
	esac

	case "$like" in
	*debian*) printf 'deb\n' && return 0 ;;
	*rhel* | *fedora* | *suse*) printf 'rpm\n' && return 0 ;;
	*arch*) printf 'arch\n' && return 0 ;;
	esac

	return 1
}

can_elevate() {
	[ "$(id -u)" -eq 0 ] && return 0
	command -v sudo >/dev/null 2>&1
}

run_root() {
	if [ "$(id -u)" -eq 0 ]; then
		"$@"
	else
		sudo "$@"
	fi
}

verify_checksum() {
	file=$1
	sums=$2
	dir=$(dirname "$sums")

	grep -E "[[:space:]]${file}\$" "$sums" >/dev/null || die "no checksum found for ${file}"

	if command -v sha256sum >/dev/null 2>&1; then
		(cd "$dir" && grep -E "[[:space:]]${file}\$" checksums.txt | sha256sum -c -)
	elif command -v shasum >/dev/null 2>&1; then
		(cd "$dir" && grep -E "[[:space:]]${file}\$" checksums.txt | shasum -a 256 -c -)
	else
		die "sha256sum or shasum is required to verify release artifacts"
	fi
}

verify_checksum_signature() {
	dir=$1
	bundle="${dir}/checksums.txt.sig"

	[ -f "$bundle" ] || return 1
	command -v cosign >/dev/null 2>&1 || die "cosign is required to verify signed checksums before a privileged package install"

	cosign verify-blob \
		--bundle "$bundle" \
		--certificate-identity-regexp "$COSIGN_IDENTITY_REGEXP" \
		--certificate-oidc-issuer "$COSIGN_OIDC_ISSUER" \
		"${dir}/checksums.txt" >/dev/null || return 1
}

asset_stem() {
	os=$1
	arch=$2
	ver=$3
	printf '%s_%s_%s_%s\n' "$BINARY" "$os" "$arch" "$ver"
}

install_package() {
	kind=$1
	file=$2

	case "$kind" in
	deb)
		need_cmd dpkg
		run_root dpkg -i "$file"
		;;
	rpm)
		if command -v rpm >/dev/null 2>&1; then
			run_root rpm -Uvh "$file"
		else
			die "rpm is required to install ${file}"
		fi
		;;
	apk)
		need_cmd apk
		run_root apk add --allow-untrusted "$file"
		;;
	arch)
		need_cmd pacman
		run_root pacman -U --noconfirm "$file"
		;;
	*)
		return 1
		;;
	esac
}

install_archive() {
	os=$1
	arch=$2
	ver=$3
	workdir=$4

	ext="tar.gz"
	if [ "$os" = "windows" ]; then
		ext="zip"
	fi

	file="$(asset_stem "$os" "$arch" "$ver").${ext}"
	url="${GITHUB}/releases/download/v${ver}/${file}"

	info "downloading ${file}"
	download "$url" "${workdir}/${file}"
	verify_checksum "$file" "${workdir}/checksums.txt"

	mkdir -p "$PREFIX"

	if [ "$ext" = "zip" ]; then
		need_cmd unzip
		unzip -qo "${workdir}/${file}" -d "$workdir"
	else
		tar -xzf "${workdir}/${file}" -C "$workdir"
	fi

	bin=$(find "$workdir" -type f \( -name "$BINARY" -o -name "${BINARY}.exe" \) | awk 'NR==1')
	[ -n "$bin" ] || die "binary ${BINARY} not found in ${file}"

	install -m 755 "$bin" "${PREFIX}/${BINARY}"
	save_prefix
	info "installed ${PREFIX}/${BINARY}"

	if ! prefix_on_path; then
		info "add ${PREFIX} to PATH so ${BINARY} is discoverable"
	fi
}

pkg_installed() {
	kind=$1

	case "$kind" in
	deb)
		command -v dpkg-query >/dev/null 2>&1 || return 1
		dpkg-query -W -f='${Status}' "$BINARY" 2>/dev/null | grep -q 'install ok installed'
		;;
	rpm)
		command -v rpm >/dev/null 2>&1 || return 1
		rpm -q "$BINARY" >/dev/null 2>&1
		;;
	apk)
		command -v apk >/dev/null 2>&1 || return 1
		apk info -e "$BINARY" >/dev/null 2>&1
		;;
	arch)
		command -v pacman >/dev/null 2>&1 || return 1
		pacman -Q "$BINARY" >/dev/null 2>&1
		;;
	*)
		return 1
		;;
	esac
}

remove_package() {
	kind=$1

	case "$kind" in
	deb)
		run_root dpkg -r "$BINARY"
		;;
	rpm)
		run_root rpm -e "$BINARY"
		;;
	apk)
		run_root apk del "$BINARY"
		;;
	arch)
		run_root pacman -R --noconfirm "$BINARY"
		;;
	*)
		return 1
		;;
	esac
}

already_installed() {
	[ -x "${PREFIX}/${BINARY}" ] && return 0
	command -v "$BINARY" >/dev/null 2>&1 && return 0
	kind=$(pkg_kind) || return 1
	pkg_installed "$kind"
}

do_uninstall() {
	removed=0

	if kind=$(pkg_kind) && pkg_installed "$kind"; then
		can_elevate || die "uninstalling the native package requires root or sudo"
		info "removing ${BINARY} ${kind} package"
		remove_package "$kind"
		removed=1
	fi

	if [ -e "${PREFIX}/${BINARY}" ]; then
		rm -f "${PREFIX}/${BINARY}"
		info "removed ${PREFIX}/${BINARY}"
		removed=1
	fi

	clear_prefix_state

	if [ "$removed" -eq 0 ]; then
		die "${BINARY} is not installed"
	fi

	info "uninstalled ${BINARY}"
}

usage() {
	cat <<EOF
Usage: install.sh [install|update|uninstall]

  install (default)  Install or replace the latest (or VERSION) release
  update             Same as install
  uninstall          Remove the native package and/or the archive binary

Environment:
  VERSION  Release tag (v1.2.0 or 1.2.0). Default: latest.
  PREFIX   Directory for archive installs. Default: \$HOME/go/bin, or the last archive prefix.
  METHOD   auto (default), package, or archive.
EOF
}

do_install() {
	os=$(os_name)
	arch=$(arch_name)
	tag=$(resolve_version)
	ver=${tag#v}

	if already_installed; then
		info "updating ${BINARY} to ${tag} for ${os}/${arch}"
	else
		info "installing ${BINARY} ${tag} for ${os}/${arch}"
	fi

	workdir=$(mktemp -d)
	trap 'rm -rf "$workdir"' EXIT INT HUP

	info "downloading checksums"
	download "${GITHUB}/releases/download/${tag}/checksums.txt" "${workdir}/checksums.txt"
	try_download "${GITHUB}/releases/download/${tag}/checksums.txt.sig" "${workdir}/checksums.txt.sig" || true

	if [ "$os" = "linux" ] && [ "$METHOD" != "archive" ]; then
		if kind=$(pkg_kind) && { [ "$METHOD" = "package" ] || can_elevate; }; then
			case "$kind" in
			deb) pkg_ext="deb" ;;
			rpm) pkg_ext="rpm" ;;
			apk) pkg_ext="apk" ;;
			arch) pkg_ext="pkg.tar.zst" ;;
			esac

			pkg="$(asset_stem linux "$arch" "$ver").${pkg_ext}"
			pkg_url="${GITHUB}/releases/download/${tag}/${pkg}"

			if try_download "$pkg_url" "${workdir}/${pkg}"; then
				if [ ! -f "${workdir}/checksums.txt.sig" ]; then
					info "signed checksum bundle missing; skipping native package"
				elif ! verify_checksum_signature "$workdir"; then
					die "checksum signature verification failed; refusing privileged package install"
				else
					verify_checksum "$pkg" "${workdir}/checksums.txt"
					info "installing native package ${pkg}"
					if install_package "$kind" "${workdir}/${pkg}"; then
						info "installed ${BINARY} via ${kind} package"
						return 0
					fi
					info "package install failed; falling back to archive"
				fi
			else
				info "native package ${pkg} not in this release; falling back to archive"
			fi
		elif [ "$METHOD" = "package" ]; then
			die "native package install requested but no supported package manager or privilege elevation is available"
		fi
	fi

	if [ "$METHOD" = "package" ]; then
		die "native package install did not complete"
	fi

	install_archive "$os" "$arch" "$ver" "$workdir"
}

main() {
	action=${1:-install}

	case "$action" in
	install | update)
		resolve_prefix
		do_install
		;;
	uninstall | remove)
		resolve_prefix
		do_uninstall
		;;
	-h | --help | help)
		usage
		;;
	*)
		die "unknown command: ${action} (try install, update, or uninstall)"
		;;
	esac
}

main "$@"
