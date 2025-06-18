#!/usr/bin/env bash

list_users_in_user_pool() {
  local pool_id="$1"
  aws cognito-idp list-users --user-pool-id "$pool_id"  | jq .
}

cognito_user_creation() {
  aws cognito-idp admin-create-user  \
    --user-pool-id "$USER_POOL_ID" \
    --username testuser@example.com \
    --user-attributes Name=email,Value=testuser@example.com \
    --message-action SUPPRESS \
    --temporary-password "$TEMP_PASSWORD"
}

display_pool_client_auth_policy() {
  local explicit_auth_flows=(
    "ALLOW_ADMIN_USER_PASSWORD_AUTH"
    "ALLOW_CUSTOM_AUTH"
    "ALLOW_USER_PASSWORD_AUTH"
    "ALLOW_USER_SRP_AUTH"
    "ALLOW_REFRESH_TOKEN_AUTH"
  )

  for policy in "${explicit_auth_flows[@]}"; do
    echo "$policy"
  done
}

update_user_pool_client_policy() {
  if [ -z "$USER_POOL_ID" ] || [ -z "$USER_POOL_CLIENT_ID" ]; then
    echo "Please set USER_POOL_ID and USER_POOL_CLIENT_ID environment variables in config.sh"
    exit 1
  fi

  local explicit_auth_flows=(
    "ALLOW_ADMIN_USER_PASSWORD_AUTH"
    "ALLOW_CUSTOM_AUTH"
    "ALLOW_USER_PASSWORD_AUTH"
    "ALLOW_USER_SRP_AUTH"
    "ALLOW_REFRESH_TOKEN_AUTH"
  )

  local flows_json
  flows_json=$(printf '"%s",' "${explicit_auth_flows[@]}")
  flows_json="[${flows_json%,}]"

  echo "Updating User Pool Client explicit auth flows to:"
  printf '%s\n' "${explicit_auth_flows[@]}"

  aws cognito-idp update-user-pool-client \
    --user-pool-id "$USER_POOL_ID" \
    --client-id "$USER_POOL_CLIENT_ID" \
    --explicit-auth-flows "$flows_json"

  if [ $? -eq 0 ]; then
    echo "Update successful"
  else
    echo "Update failed"
  fi
}
