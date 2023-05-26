#!/usr/bin/env bash

# Copyright 2023 lichangao(李长傲) <changao.li.work@outlook.com>. All rights reserved.
# Use of this source code is governed by a MIT style
# license that can be found in the LICENSE file. The original repo for
# this file is https://github.com/changaolee/skeleton.

# -----------------------------------------------------------------------------
# Version management helpers.  These functions help to set, save and load the
# following variables:
#
#    SKT_GIT_COMMIT - The git commit id corresponding to this
#          source code.
#    SKT_GIT_TREE_STATE - "clean" indicates no changes since the git commit id
#        "dirty" indicates source code changes after the git commit id
#        "archive" indicates the tree was produced by 'git archive'
#    SKT_GIT_VERSION - "vX.Y" used to indicate the last release version.
#    SKT_GIT_MAJOR - The major part of the version
#    SKT_GIT_MINOR - The minor component of the version

# Grovels through git to set a set of env variables.
#
# If SKT_GIT_VERSION_FILE, this function will load from that file instead of
# querying git.
skt::version::get_version_vars() {
  if [[ -n ${SKT_GIT_VERSION_FILE-} ]]; then
    skt::version::load_version_vars "${SKT_GIT_VERSION_FILE}"
    return
  fi

  # shellcheck disable=SC2016,SC2050
  # Disabled as we're not expanding these at runtime, but rather expecting
  # that another tool may have expanded these and rewritten the source (!)
  if [[ '$Format:%%$' == "%" ]]; then
    SKT_GIT_COMMIT='$Format:%H$'
    SKT_GIT_TREE_STATE="archive"
    # When a 'git archive' is exported, the '$Format:%D$' below will look
    # something like 'HEAD -> release-1.8, tag: v1.8.3' where then 'tag: '
    # can be extracted from it.
    if [[ '$Format:%D$' =~ tag:\ (v[^ ,]+) ]]; then
      SKT_GIT_VERSION="${BASH_REMATCH[1]}"
    fi
  fi

  local git=(git --work-tree "${SKT_ROOT}")

  if [[ -n ${SKT_GIT_COMMIT-} ]] || SKT_GIT_COMMIT=$("${git[@]}" rev-parse "HEAD^{commit}" 2>/dev/null); then
    if [[ -z ${SKT_GIT_TREE_STATE-} ]]; then
      # Check if the tree is dirty.  default to dirty
      if git_status=$("${git[@]}" status --porcelain 2>/dev/null) && [[ -z ${git_status} ]]; then
        SKT_GIT_TREE_STATE="clean"
      else
        SKT_GIT_TREE_STATE="dirty"
      fi
    fi

    # Use git describe to find the version based on tags.
    if [[ -n ${SKT_GIT_VERSION-} ]] || SKT_GIT_VERSION=$("${git[@]}" describe --tags --always --match='v*' 2>/dev/null); then
      # This translates the "git describe" to an actual semver.org
      # compatible semantic version that looks something like this:
      #   v1.1.0-alpha.0.6+84c76d1142ea4d
      #
      # TODO: We continue calling this "git version" because so many
      # downstream consumers are expecting it there.
      #
      # These regexes are painful enough in sed...
      # We don't want to do them in pure shell, so disable SC2001
      # shellcheck disable=SC2001
      DASHES_IN_VERSION=$(echo "${SKT_GIT_VERSION}" | sed "s/[^-]//g")
      if [[ "${DASHES_IN_VERSION}" == "---" ]]; then
        # shellcheck disable=SC2001
        # We have distance to subversion (v1.1.0-subversion-1-gCommitHash)
        SKT_GIT_VERSION=$(echo "${SKT_GIT_VERSION}" | sed "s/-\([0-9]\{1,\}\)-g\([0-9a-f]\{14\}\)$/.\1\+\2/")
      elif [[ "${DASHES_IN_VERSION}" == "--" ]]; then
        # shellcheck disable=SC2001
        # We have distance to base tag (v1.1.0-1-gCommitHash)
        SKT_GIT_VERSION=$(echo "${SKT_GIT_VERSION}" | sed "s/-g\([0-9a-f]\{14\}\)$/+\1/")
      fi
      if [[ "${SKT_GIT_TREE_STATE}" == "dirty" ]]; then
        # git describe --dirty only considers changes to existing files, but
        # that is problematic since new untracked .go files affect the build,
        # so use our idea of "dirty" from git status instead.
        # TODO?
        #SKT_GIT_VERSION+="-dirty"
        :
      fi

      # Try to match the "git describe" output to a regex to try to extract
      # the "major" and "minor" versions and whether this is the exact tagged
      # version or whether the tree is between two tagged versions.
      if [[ "${SKT_GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?([-].*)?([+].*)?$ ]]; then
        SKT_GIT_MAJOR=${BASH_REMATCH[1]}
        SKT_GIT_MINOR=${BASH_REMATCH[2]}
        if [[ -n "${BASH_REMATCH[4]}" ]]; then
          SKT_GIT_MINOR+="+"
        fi
      fi

      # If SKT_GIT_VERSION is not a valid Semantic Version, then refuse to build.
      if ! [[ "${SKT_GIT_VERSION}" =~ ^v([0-9]+)\.([0-9]+)(\.[0-9]+)?(-[0-9A-Za-z.-]+)?(\+[0-9A-Za-z.-]+)?$ ]]; then
        echo "SKT_GIT_VERSION should be a valid Semantic Version. Current value: ${SKT_GIT_VERSION}"
        echo "Please see more details here: https://semver.org"
        exit 1
      fi
    fi
  fi
}

# Saves the environment flags to $1
skt::version::save_version_vars() {
  local version_file=${1-}
  [[ -n ${version_file} ]] || {
    echo "!!! Internal error.  No file specified in skt::version::save_version_vars"
    return 1
  }

  cat <<EOF >"${version_file}"
SKT_GIT_COMMIT='${SKT_GIT_COMMIT-}'
SKT_GIT_TREE_STATE='${SKT_GIT_TREE_STATE-}'
SKT_GIT_VERSION='${SKT_GIT_VERSION-}'
SKT_GIT_MAJOR='${SKT_GIT_MAJOR-}'
SKT_GIT_MINOR='${SKT_GIT_MINOR-}'
EOF
}

# Loads up the version variables from file $1
skt::version::load_version_vars() {
  local version_file=${1-}
  [[ -n ${version_file} ]] || {
    echo "!!! Internal error.  No file specified in skt::version::load_version_vars"
    return 1
  }

  source "${version_file}"
}
