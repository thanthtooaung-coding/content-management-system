log() {
  echo "[$(date '+%Y-%m-%d %H:%M:%S')] $1"
}

list_user_pools() {
  aws cognito-idp list-user-pools --max-results 60 | jq .
}
update_user_pool_client_policy(){
  local USER_POOL_ID="us-east-1_jeV07U5ru"
  local CLIENT_ID="6nc06fsj3c4mdl71a18gass9p5"

  aws cognito-idp update-user-pool-client \
    --user-pool-id $USER_POOL_ID \
    --client-id $CLIENT_ID \
    --explicit-auth-flows ALLOW_USER_PASSWORD_AUTH ALLOW_REFRESH_TOKEN_AUTH ALLOW_USER_SRP_AUTH

}
describe_user_pool() {
  local pool_id="$1"
  local parameter="$2"
  local output_file="./.json"

  if [ -z "$pool_id" ]; then
    log "❌ Missing user pool ID."
    echo "Usage: $0 describe-pool <USER_POOL_ID> [PARAMETER]"
    return 1
  fi

  if [ ! -f "$output_file" ]; then
    log "⚠️ Output file '$output_file' does not exist. Creating it..."
    touch "$output_file"
  else
    log "✅ Output file '$output_file' already exists."
  fi

  log "📥 Fetching user pool details for: $pool_id"
  aws cognito-idp describe-user-pool --user-pool-id "$pool_id" | jq '.' > "$output_file"
  log "📄 Full output saved to $output_file"

  if [ -n "$parameter" ]; then
    local value
    value=$(jq -r ".UserPool.$parameter" "$output_file")
    if [ "$value" == "null" ]; then
      log "⚠️ Parameter '$parameter' not found in the user pool data."
      echo ""
      return 1
    fi
    echo "$value"
  else
    log "ℹ️ No parameter specified, full output saved."
  fi
}
