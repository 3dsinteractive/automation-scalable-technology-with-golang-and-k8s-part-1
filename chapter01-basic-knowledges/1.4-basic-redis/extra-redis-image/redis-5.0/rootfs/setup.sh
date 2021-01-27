#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail
#set -o xtrace
# shellcheck disable=SC1091

# Load libraries
. /libos.sh
. /libfs.sh
. /libredis.sh

# Load Redis environment variables
eval "$(redis_env)"

# Ensure Redis environment variables settings are valid
redis_validate
# Ensure Redis is stopped when this script ends
trap "redis_stop" EXIT
am_i_root && ensure_user_exists "$REDIS_DAEMON_USER" "$REDIS_DAEMON_GROUP"
# Ensure Redis is initialized
redis_initialize
