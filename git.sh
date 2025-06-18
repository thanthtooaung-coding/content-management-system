#!/usr/bin/env bash

commit_message="$1"
branch_name="$2"

git add .
git commit -m "${commit_message}"
git pull origin "${branch_name}"
git push orign "${branch_name}"

gh workflow list

gh run list --workflow=workflow.yml