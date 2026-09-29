#!/bin/sh
# SPDX-License-Identifier: AGPL-3.0-or-later
#
# Merge the generated changelog pull request.
#
# Requests auto-merge, and falls back to merging directly when the
# repository does not have auto-merge enabled, so the changelog stays
# fully automated.
#
# Usage:
#   sh ./scripts/merge-changelog-pr.sh <pull-request-number>

set -eu

PR_NUMBER="${1:?usage: merge-changelog-pr.sh <pull-request-number>}"

head_ref="$(gh pr view "$PR_NUMBER" --json headRefName --jq .headRefName)"

case "$head_ref" in
  docs/*) ;;
  *)
    echo "::error::Refusing to merge a changelog pull request from '$head_ref'."
    exit 1
    ;;
esac

if ! gh pr merge --auto --squash "$PR_NUMBER"; then
  echo "Auto-merge is unavailable; merging the clean pull request directly"
  gh pr merge --squash "$PR_NUMBER"
fi
