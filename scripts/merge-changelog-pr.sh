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

# Must match the branch configured for create-pull-request in the changelog workflow.
EXPECTED_BRANCH=docs/update-changelog

head_ref="$(gh pr view "$PR_NUMBER" --json headRefName --jq .headRefName)"
head_sha="$(gh pr view "$PR_NUMBER" --json headRefOid --jq .headRefOid)"

if [ "$head_ref" != "$EXPECTED_BRANCH" ]; then
  echo "::error::Refusing to merge a pull request from '$head_ref'."
  exit 1
fi

if ! gh pr merge --auto --squash --match-head-commit "$head_sha" "$PR_NUMBER"; then
  echo "Auto-merge is unavailable; merging the clean pull request directly"
  gh pr merge --squash --match-head-commit "$head_sha" "$PR_NUMBER"
fi
