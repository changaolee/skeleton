#!/usr/bin/env bash

# Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/changaolee/skeleton.


function skt::util::sourced_variable {
  # Call this function to tell shellcheck that a variable is supposed to
  # be used from other calling context. This helps quiet an "unused
  # variable" warning from shellcheck and also document your code.
  true
}

skt::util::sortable_date() {
  date "+%Y%m%d-%H%M%S"
}

# arguments: target, item1, item2, item3, ...
# returns 0 if target is in the given items, 1 otherwise.
skt::util::array_contains() {
  local search="$1"
  local element
  shift
  for element; do
    if [[ "${element}" == "${search}" ]]; then
      return 0
     fi
  done
  return 1
}

skt::util::wait_for_url() {
  local url=$1
  local prefix=${2:-}
  local wait=${3:-1}
  local times=${4:-30}
  local maxtime=${5:-1}

  command -v curl >/dev/null || {
    skt::log::usage "curl must be installed"
    exit 1
  }

  local i
  for i in $(seq 1 "${times}"); do
    local out
    if out=$(curl --max-time "${maxtime}" -gkfs "${url}" 2>/dev/null); then
      skt::log::status "On try ${i}, ${prefix}: ${out}"
      return 0
    fi
    sleep "${wait}"
  done
  skt::log::error "Timed out waiting for ${prefix} to answer at ${url}; tried ${times} waiting ${wait} between each"
  return 1
}

# Example:  skt::util::wait_for_success 120 5 "sktctl get nodes|grep localhost"
# arguments: wait time, sleep time, shell command
# returns 0 if the shell command get output, 1 otherwise.
skt::util::wait_for_success(){
  local wait_time="$1"
  local sleep_time="$2"
  local cmd="$3"
  while [ "$wait_time" -gt 0 ]; do
    if eval "$cmd"; then
      return 0
    else
      sleep "$sleep_time"
      wait_time=$((wait_time-sleep_time))
    fi
  done
  return 1
}

# Example:  skt::util::trap_add 'echo "in trap DEBUG"' DEBUG
# See: http://stackoverflow.com/questions/3338030/multiple-bash-traps-for-the-same-signal
skt::util::trap_add() {
  local trap_add_cmd
  trap_add_cmd=$1
  shift

  for trap_add_name in "$@"; do
    local existing_cmd
    local new_cmd

    # Grab the currently defined trap commands for this trap
    existing_cmd=$(trap -p "${trap_add_name}" |  awk -F"'" '{print $2}')

    if [[ -z "${existing_cmd}" ]]; then
      new_cmd="${trap_add_cmd}"
    else
      new_cmd="${trap_add_cmd};${existing_cmd}"
    fi

    # Assign the test. Disable the shellcheck warning telling that trap
    # commands should be single quoted to avoid evaluating them at this
    # point instead evaluating them at run time. The logic of adding new
    # commands to a single trap requires them to be evaluated right away.
    # shellcheck disable=SC2064
    trap "${new_cmd}" "${trap_add_name}"
  done
}

# Opposite of skt::util::ensure-temp-dir()
skt::util::cleanup-temp-dir() {
  rm -rf "${SKT_TEMP}"
}

# Create a temp dir that'll be deleted at the end of this bash session.
#
# Vars set:
#   SKT_TEMP
skt::util::ensure-temp-dir() {
  if [[ -z ${SKT_TEMP-} ]]; then
    SKT_TEMP=$(mktemp -d 2>/dev/null || mktemp -d -t iamrnetes.XXXXXX)
    skt::util::trap_add skt::util::cleanup-temp-dir EXIT
  fi
}

skt::util::host_os() {
  local host_os
  case "$(uname -s)" in
    Darwin)
      host_os=darwin
      ;;
    Linux)
      host_os=linux
      ;;
    *)
      skt::log::error "Unsupported host OS.  Must be Linux or Mac OS X."
      exit 1
      ;;
  esac
  echo "${host_os}"
}

skt::util::host_arch() {
  local host_arch
  case "$(uname -m)" in
    x86_64*)
      host_arch=amd64
      ;;
    i?86_64*)
      host_arch=amd64
      ;;
    amd64*)
      host_arch=amd64
      ;;
    aarch64*)
      host_arch=arm64
      ;;
    arm64*)
      host_arch=arm64
      ;;
    arm*)
      host_arch=arm
      ;;
    i?86*)
      host_arch=x86
      ;;
    s390x*)
      host_arch=s390x
      ;;
    ppc64le*)
      host_arch=ppc64le
      ;;
    *)
      skt::log::error "Unsupported host arch. Must be x86_64, 386, arm, arm64, s390x or ppc64le."
      exit 1
      ;;
  esac
  echo "${host_arch}"
}

# This figures out the host platform without relying on golang.  We need this as
# we don't want a golang install to be a prerequisite to building yet we need
# this info to figure out where the final binaries are placed.
skt::util::host_platform() {
  echo "$(skt::util::host_os)/$(skt::util::host_arch)"
}

# looks for $1 in well-known output locations for the platform ($2)
# $SKT_ROOT must be set
skt::util::find-binary-for-platform() {
  local -r lookfor="$1"
  local -r platform="$2"
  local locations=(
    "${SKT_ROOT}/_output/bin/${lookfor}"
    "${SKT_ROOT}/_output/${platform}/${lookfor}"
    "${SKT_ROOT}/_output/local/bin/${platform}/${lookfor}"
    "${SKT_ROOT}/_output/platforms/${platform}/${lookfor}"
  )

  # List most recently-updated location.
  local -r bin=$( (ls -t "${locations[@]}" 2>/dev/null || true) | head -1 )
  echo -n "${bin}"
}

# looks for $1 in well-known output locations for the host platform
# $SKT_ROOT must be set
skt::util::find-binary() {
  skt::util::find-binary-for-platform "$1" "$(skt::util::host_platform)"
}

# Run all known doc generators (today gendocs and genman for sktctl)
# $1 is the directory to put those generated documents
skt::util::gen-docs() {
  local dest="$1"

  # Find binary
  gendocs=$(skt::util::find-binary "gendocs")
  gensktdocs=$(skt::util::find-binary "gensktdocs")
  genman=$(skt::util::find-binary "genman")
  genyaml=$(skt::util::find-binary "genyaml")
  genfeddocs=$(skt::util::find-binary "genfeddocs")

  # TODO: If ${genfeddocs} is not used from anywhere (it isn't used at
  # least from k/k tree), remove it completely.
  skt::util::sourced_variable "${genfeddocs}"

  mkdir -p "${dest}/docs/guide/en-US/cmd/sktctl/"
  "${gendocs}" "${dest}/docs/guide/en-US/cmd/sktctl/"

  mkdir -p "${dest}/docs/guide/en-US/cmd/"
  "${gensktdocs}" "${dest}/docs/guide/en-US/cmd/" "skt-apiserver"
  "${gensktdocs}" "${dest}/docs/guide/en-US/cmd/" "skt-authz-server"
  "${gensktdocs}" "${dest}/docs/guide/en-US/cmd/" "skt-pump"
  "${gensktdocs}" "${dest}/docs/guide/en-US/cmd/" "skt-watcher"
  "${gensktdocs}" "${dest}/docs/guide/en-US/cmd/sktctl" "sktctl"

  mkdir -p "${dest}/docs/man/man1/"
  "${genman}" "${dest}/docs/man/man1/" "skt-apiserver"
  "${genman}" "${dest}/docs/man/man1/" "skt-authz-server"
  "${genman}" "${dest}/docs/man/man1/" "skt-pump"
  "${genman}" "${dest}/docs/man/man1/" "skt-watcher"
  "${genman}" "${dest}/docs/man/man1/" "sktctl"

  mkdir -p "${dest}/docs/guide/en-US/yaml/sktctl/"
  "${genyaml}" "${dest}/docs/guide/en-US/yaml/sktctl/"

  # create the list of generated files
  pushd "${dest}" > /dev/null || return 1
  touch docs/.generated_docs
  find . -type f | cut -sd / -f 2- | LC_ALL=C sort > docs/.generated_docs
  popd > /dev/null || return 1
}

# Removes previously generated docs-- we don't want to check them in. $SKT_ROOT
# must be set.
skt::util::remove-gen-docs() {
  if [ -e "${SKT_ROOT}/docs/.generated_docs" ]; then
    # remove all of the old docs; we don't want to check them in.
    while read -r file; do
      rm "${SKT_ROOT}/${file}" 2>/dev/null || true
    done <"${SKT_ROOT}/docs/.generated_docs"
    # The docs/.generated_docs file lists itself, so we don't need to explicitly
    # delete it.
  fi
}

# Returns the name of the upstream remote repository name for the local git
# repo, e.g. "upstream" or "origin".
skt::util::git_upstream_remote_name() {
  git remote -v | grep fetch |\
    grep -E 'github.com[/:]changaolee/skeleton' |\
    head -n 1 | awk '{print $1}'
}

# Exits script if working directory is dirty. If it's run interactively in the terminal
# the user can commit changes in a second terminal. This script will wait.
skt::util::ensure_clean_working_dir() {
  while ! git diff HEAD --exit-code &>/dev/null; do
    echo -e "\nUnexpected dirty working directory:\n"
    if tty -s; then
        git status -s
    else
        git diff -a # be more verbose in log files without tty
        exit 1
    fi | sed 's/^/  /'
    echo -e "\nCommit your changes in another terminal and then continue here by pressing enter."
    read -r
  done 1>&2
}

# Find the base commit using:
# $PULL_BASE_SHA if set (from Prow)
# current ref from the remote upstream branch
skt::util::base_ref() {
  local -r git_branch=$1

  if [[ -n ${PULL_BASE_SHA:-} ]]; then
    echo "${PULL_BASE_SHA}"
    return
  fi

  full_branch="$(skt::util::git_upstream_remote_name)/${git_branch}"

  # make sure the branch is valid, otherwise the check will pass erroneously.
  if ! git describe "${full_branch}" >/dev/null; then
    # abort!
    exit 1
  fi

  echo "${full_branch}"
}

# Checks whether there are any files matching pattern $2 changed between the
# current branch and upstream branch named by $1.
# Returns 1 (false) if there are no changes
#         0 (true) if there are changes detected.
skt::util::has_changes() {
  local -r git_branch=$1
  local -r pattern=$2
  local -r not_pattern=${3:-totallyimpossiblepattern}

  local base_ref
  base_ref=$(skt::util::base_ref "${git_branch}")
  echo "Checking for '${pattern}' changes against '${base_ref}'"

  # notice this uses ... to find the first shared ancestor
  if git diff --name-only "${base_ref}...HEAD" | grep -v -E "${not_pattern}" | grep "${pattern}" > /dev/null; then
    return 0
  fi
  # also check for pending changes
  if git status --porcelain | grep -v -E "${not_pattern}" | grep "${pattern}" > /dev/null; then
    echo "Detected '${pattern}' uncommitted changes."
    return 0
  fi
  echo "No '${pattern}' changes detected."
  return 1
}

skt::util::download_file() {
  local -r url=$1
  local -r destination_file=$2

  rm "${destination_file}" 2&> /dev/null || true

  for i in $(seq 5)
  do
    if ! curl -fsSL --retry 3 --keepalive-time 2 "${url}" -o "${destination_file}"; then
      echo "Downloading ${url} failed. $((5-i)) retries left."
      sleep 1
    else
      echo "Downloading ${url} succeed"
      return 0
    fi
  done
  return 1
}

# Test whether openssl is installed.
# Sets:
#  OPENSSL_BIN: The path to the openssl binary to use
function skt::util::test_openssl_installed {
    if ! openssl version >& /dev/null; then
      echo "Failed to run openssl. Please ensure openssl is installed"
      exit 1
    fi

    OPENSSL_BIN=$(command -v openssl)
}

# creates a client CA, args are sudo, dest-dir, ca-id, purpose
# purpose is dropped in after "key encipherment", you usually want
# '"client auth"'
# '"server auth"'
# '"client auth","server auth"'
function skt::util::create_signing_certkey {
    local sudo=$1
    local dest_dir=$2
    local id=$3
    local purpose=$4
    # Create client ca
    ${sudo} /usr/bin/env bash -e <<EOF
    rm -f "${dest_dir}/${id}-ca.crt" "${dest_dir}/${id}-ca.key"
    ${OPENSSL_BIN} req -x509 -sha256 -new -nodes -days 365 -newkey rsa:2048 -keyout "${dest_dir}/${id}-ca.key" -out "${dest_dir}/${id}-ca.crt" -subj "/C=xx/ST=x/L=x/O=x/OU=x/CN=ca/emailAddress=x/"
    echo '{"signing":{"default":{"expiry":"43800h","usages":["signing","key encipherment",${purpose}]}}}' > "${dest_dir}/${id}-ca-config.json"
EOF
}

# signs a client certificate: args are sudo, dest-dir, CA, filename (roughly), username, groups...
function skt::util::create_client_certkey {
    local sudo=$1
    local dest_dir=$2
    local ca=$3
    local id=$4
    local cn=${5:-$4}
    local groups=""
    local SEP=""
    shift 5
    while [ -n "${1:-}" ]; do
        groups+="${SEP}{\"O\":\"$1\"}"
        SEP=","
        shift 1
    done
    ${sudo} /usr/bin/env bash -e <<EOF
    cd ${dest_dir}
    echo '{"CN":"${cn}","names":[${groups}],"hosts":[""],"key":{"algo":"rsa","size":2048}}' | ${CFSSL_BIN} gencert -ca=${ca}.crt -ca-key=${ca}.key -config=${ca}-config.json - | ${CFSSLJSON_BIN} -bare client-${id}
    mv "client-${id}-key.pem" "client-${id}.key"
    mv "client-${id}.pem" "client-${id}.crt"
    rm -f "client-${id}.csr"
EOF
}

# signs a serving certificate: args are sudo, dest-dir, ca, filename (roughly), subject, hosts...
function skt::util::create_serving_certkey {
    local sudo=$1
    local dest_dir=$2
    local ca=$3
    local id=$4
    local cn=${5:-$4}
    local hosts=""
    local SEP=""
    shift 5
    while [ -n "${1:-}" ]; do
        hosts+="${SEP}\"$1\""
        SEP=","
        shift 1
    done
    ${sudo} /usr/bin/env bash -e <<EOF
    cd ${dest_dir}
    echo '{"CN":"${cn}","hosts":[${hosts}],"key":{"algo":"rsa","size":2048}}' | ${CFSSL_BIN} gencert -ca=${ca}.crt -ca-key=${ca}.key -config=${ca}-config.json - | ${CFSSLJSON_BIN} -bare serving-${id}
    mv "serving-${id}-key.pem" "serving-${id}.key"
    mv "serving-${id}.pem" "serving-${id}.crt"
    rm -f "serving-${id}.csr"
EOF
}

# creates a self-contained sktconfig: args are sudo, dest-dir, ca file, host, port, client id, token(optional)
function skt::util::write_client_sktconfig {
    local sudo=$1
    local dest_dir=$2
    local ca_file=$3
    local api_host=$4
    local api_port=$5
    local client_id=$6
    local token=${7:-}
    cat <<EOF | ${sudo} tee "${dest_dir}"/"${client_id}".sktconfig > /dev/null
apiVersion: v1
kind: Config
clusters:
  - cluster:
      certificate-authority: ${ca_file}
      server: https://${api_host}:${api_port}/
    name: local-up-cluster
users:
  - user:
      token: ${token}
      client-certificate: ${dest_dir}/client-${client_id}.crt
      client-key: ${dest_dir}/client-${client_id}.key
    name: local-up-cluster
contexts:
  - context:
      cluster: local-up-cluster
      user: local-up-cluster
    name: local-up-cluster
current-context: local-up-cluster
EOF

    # flatten the sktconfig files to make them self contained
    username=$(whoami)
    ${sudo} /usr/bin/env bash -e <<EOF
    $(skt::util::find-binary sktctl) --sktconfig="${dest_dir}/${client_id}.sktconfig" config view --minify --flatten > "/tmp/${client_id}.sktconfig"
    mv -f "/tmp/${client_id}.sktconfig" "${dest_dir}/${client_id}.sktconfig"
    chown ${username} "${dest_dir}/${client_id}.sktconfig"
EOF
}

# Determines if docker can be run, failures may simply require that the user be added to the docker group.
function skt::util::ensure_docker_daemon_connectivity {
  IFS=" " read -ra DOCKER <<< "${DOCKER_OPTS}"
  # Expand ${DOCKER[@]} only if it's not unset. This is to work around
  # Bash 3 issue with unbound variable.
  DOCKER=(docker ${DOCKER[@]:+"${DOCKER[@]}"})
  if ! "${DOCKER[@]}" info > /dev/null 2>&1 ; then
    cat <<'EOF' >&2
Can't connect to 'docker' daemon.  please fix and retry.

Possible causes:
  - Docker Daemon not started
    - Linux: confirm via your init system
    - macOS w/ docker-machine: run `docker-machine ls` and `docker-machine start <name>`
    - macOS w/ Docker for Mac: Check the menu bar and start the Docker application
  - DOCKER_HOST hasn't been set or is set incorrectly
    - Linux: domain socket is used, DOCKER_* should be unset. In Bash run `unset ${!DOCKER_*}`
    - macOS w/ docker-machine: run `eval "$(docker-machine env <name>)"`
    - macOS w/ Docker for Mac: domain socket is used, DOCKER_* should be unset. In Bash run `unset ${!DOCKER_*}`
  - Other things to check:
    - Linux: User isn't in 'docker' group.  Add and relogin.
      - Something like 'sudo usermod -a -G docker ${USER}'
      - RHEL7 bug and workaround: https://bugzilla.redhat.com/show_bug.cgi?id=1119282#c8
EOF
    return 1
  fi
}

# Wait for background jobs to finish. Return with
# an error status if any of the jobs failed.
skt::util::wait-for-jobs() {
  local fail=0
  local job
  for job in $(jobs -p); do
    wait "${job}" || fail=$((fail + 1))
  done
  return ${fail}
}

# skt::util::join <delim> <list...>
# Concatenates the list elements with the delimiter passed as first parameter
#
# Ex: skt::util::join , a b c
#  -> a,b,c
function skt::util::join {
  local IFS="$1"
  shift
  echo "$*"
}

# Downloads cfssl/cfssljson/cfssl-certinfo into $1 directory if they do not already exist in PATH
#
# Assumed vars:
#   $1 (cfssl directory) (optional)
#
# Sets:
#  CFSSL_BIN: The path of the installed cfssl binary
#  CFSSLJSON_BIN: The path of the installed cfssljson binary
#  CFSSLCERTINFO_BIN: The path of the installed cfssl-certinfo binary
#
function skt::util::ensure-cfssl {
  if command -v cfssl &>/dev/null && command -v cfssljson &>/dev/null && command -v cfssl-certinfo &>/dev/null; then
    CFSSL_BIN=$(command -v cfssl)
    CFSSLJSON_BIN=$(command -v cfssljson)
    CFSSLCERTINFO_BIN=$(command -v cfssl-certinfo)
    return 0
  fi

  host_arch=$(skt::util::host_arch)

  if [[ "${host_arch}" != "amd64" ]]; then
    echo "Cannot download cfssl on non-amd64 hosts and cfssl does not appear to be installed."
    echo "Please install cfssl, cfssljson and cfssl-certinfo and verify they are in \$PATH."
    echo "Hint: export PATH=\$PATH:\$GOPATH/bin; go get -u github.com/cloudflare/cfssl/cmd/..."
    exit 1
  fi

  # Create a temp dir for cfssl if no directory was given
  local cfssldir=${1:-}
  if [[ -z "${cfssldir}" ]]; then
    cfssldir="$HOME/bin"
  fi

  mkdir -p "${cfssldir}"
  pushd "${cfssldir}" > /dev/null || return 1

  echo "Unable to successfully run 'cfssl' from ${PATH}; downloading instead..."
  kernel=$(uname -s)
  case "${kernel}" in
    Linux)
      curl --retry 10 -L -o cfssl https://pkg.cfssl.org/R1.2/cfssl_linux-amd64
      curl --retry 10 -L -o cfssljson https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64
      curl --retry 10 -L -o cfssl-certinfo https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64
      ;;
    Darwin)
      curl --retry 10 -L -o cfssl https://pkg.cfssl.org/R1.2/cfssl_darwin-amd64
      curl --retry 10 -L -o cfssljson https://pkg.cfssl.org/R1.2/cfssljson_darwin-amd64
      curl --retry 10 -L -o cfssl-certinfo https://pkg.cfssl.org/R1.2/cfssl-certinfo_darwin-amd64
      ;;
    *)
      echo "Unknown, unsupported platform: ${kernel}." >&2
      echo "Supported platforms: Linux, Darwin." >&2
      exit 2
  esac

  chmod +x cfssl || true
  chmod +x cfssljson || true
  chmod +x cfssl-certinfo || true

  CFSSL_BIN="${cfssldir}/cfssl"
  CFSSLJSON_BIN="${cfssldir}/cfssljson"
  CFSSLCERTINFO_BIN="${cfssldir}/cfssl-certinfo"
  if [[ ! -x ${CFSSL_BIN} || ! -x ${CFSSLJSON_BIN} || ! -x ${CFSSLCERTINFO_BIN} ]]; then
    echo "Failed to download 'cfssl'."
    echo "Please install cfssl, cfssljson and cfssl-certinfo and verify they are in \$PATH."
    echo "Hint: export PATH=\$PATH:\$GOPATH/bin; go get -u github.com/cloudflare/cfssl/cmd/..."
    exit 1
  fi
  popd > /dev/null || return 1
}

# skt::util::ensure-gnu-sed
# Determines which sed binary is gnu-sed on linux/darwin
#
# Sets:
#  SED: The name of the gnu-sed binary
#
function skt::util::ensure-gnu-sed {
  # NOTE: the echo below is a workaround to ensure sed is executed before the grep.
  sed_help="$(LANG=C sed --help 2>&1 || true)"
  if echo "${sed_help}" | grep -q "GNU\|BusyBox"; then
    SED="sed"
  elif command -v gsed &>/dev/null; then
    SED="gsed"
  else
    skt::log::error "Failed to find GNU sed as sed or gsed. If you are on Mac: brew install gnu-sed." >&2
    return 1
  fi
  skt::util::sourced_variable "${SED}"
}

# skt::util::check-file-in-alphabetical-order <file>
# Check that the file is in alphabetical order
#
function skt::util::check-file-in-alphabetical-order {
  local failure_file="$1"
  if ! diff -u "${failure_file}" <(LC_ALL=C sort "${failure_file}"); then
    {
      echo
      echo "${failure_file} is not in alphabetical order. Please sort it:"
      echo
      echo "  LC_ALL=C sort -o ${failure_file} ${failure_file}"
      echo
    } >&2
  false
  fi
}

# skt::util::require-jq
# Checks whether jq is installed.
function skt::util::require-jq {
  if ! command -v jq &>/dev/null; then
    echo "jq not found. Please install." 1>&2
    return 1
  fi
}

# outputs md5 hash of $1, works on macOS and Linux
function skt::util::md5() {
  if which md5 >/dev/null 2>&1; then
    md5 -q "$1"
  else
    md5sum "$1" | awk '{ print $1 }'
  fi
}

# skt::util::read-array
# Reads in stdin and adds it line by line to the array provided. This can be
# used instead of "mapfile -t", and is bash 3 compatible.
#
# Assumed vars:
#   $1 (name of array to create/modify)
#
# Example usage:
# skt::util::read-array files < <(ls -1)
#
function skt::util::read-array {
  local i=0
  unset -v "$1"
  while IFS= read -r "$1[i++]"; do :; done
  eval "[[ \${$1[--i]} ]]" || unset "$1[i]" # ensures last element isn't empty
}

# Some useful colors.
if [[ -z "${color_start-}" ]]; then
  declare -r color_start="\033["
  declare -r color_red="${color_start}0;31m"
  declare -r color_yellow="${color_start}0;33m"
  declare -r color_green="${color_start}0;32m"
  declare -r color_blue="${color_start}1;34m"
  declare -r color_cyan="${color_start}1;36m"
  declare -r color_norm="${color_start}0m"

  skt::util::sourced_variable "${color_start}"
  skt::util::sourced_variable "${color_red}"
  skt::util::sourced_variable "${color_yellow}"
  skt::util::sourced_variable "${color_green}"
  skt::util::sourced_variable "${color_blue}"
  skt::util::sourced_variable "${color_cyan}"
  skt::util::sourced_variable "${color_norm}"
fi

# ex: ts=2 sw=2 et filetype=sh
