#!/usr/bin/env bash

commit_message="$1"
branch_name="$2"

git add .
git commit -m "${commit_message}"
git pull origin "${branch_name}"
git push origin "${branch_name}"

gh workflow list
gh run list --workflow=workflow.yml --limit=3

echo "What would you like to do?"
echo "1) Watch latest workflow run"
echo "2) View workflow run logs"
echo "3) Trigger workflow manually"
echo "4) Exit"

read -p "Choose an option (1-4): " choice

case $choice in
    1)
        latest_run=$(gh run list --workflow=workflow.yml --limit=1 --json databaseId --jq '.[0].databaseId')
        gh run watch "$latest_run"
        ;;
    2)
        latest_run=$(gh run list --workflow=workflow.yml --limit=1 --json databaseId --jq '.[0].databaseId')
        gh run view "$latest_run" --log
        ;;
    3)
        gh workflow run workflow.yml
        ;;
    4)
        exit 0
        ;;
    *)
        echo "Invalid option"
        ;;
esac