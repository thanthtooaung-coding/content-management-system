#!/usr/bin/env bash

source "$(dirname "$0")/../config.sh"
source "$(dirname "$0")/../lib/user.sh"
source "$(dirname "$0")/../lib/pool.sh"

main() {
  case "$1" in
    list-pool)
      list_user_pools
      ;;

    update-user-pool-client-policy)
      update_user_pool_client_policy
      ;;

    describe-pool)
      if [ -z "$2" ] || [ -z "$3" ]; then
        echo "Usage: $0 describe-pool <USER_POOL_ID> <PARAMETER>"
        exit 1
      fi
      describe_user_pool "$2" "$3"
      ;;
    update-user-pool-client-policy)
        update_user_pool_client_policy
        ;;
    dis-pool-auth_policy)
      display_pool_client_auth_policy
      ;;

    list-users)
      if [ -z "$2" ]; then
        echo "Usage: $0 list-users <USER_POOL_ID>"
        exit 1
      fi
      list_users_in_user_pool "$2"
      ;;

    create-user)
      cognito_user_creation
      ;;

   *)
     echo -e "\033[1;33mUsage:\033[0m $0 {"
     echo -e "  list-pool |"
     echo -e "  dis-pool-auth_policy |"
     echo -e "  update-user-pool-client-policy |"
     echo -e "  describe-pool <USER_POOL_ID> <PARAMETER> |"
     echo -e "  list-users <USER_POOL_ID> |"
     echo -e "  create-user"
     echo -e "}"
     exit 1
     ;;

  esac
}

main "$@"
