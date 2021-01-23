#!/bin/bash
#
# Bitnami Redis library

# shellcheck disable=SC1091

# Load Generic Libraries
. /libfile.sh
. /liblog.sh
. /libos.sh
. /libservice.sh
. /libvalidations.sh

# Functions

########################
# Retrieve a configuration setting value
# Globals:
#   REDIS_BASEDIR
# Arguments:
#   $1 - key
# Returns:
#   None
#########################
redis_conf_get() {
    local key="${1:?missing key}"

    grep -E "^\s*$key " "${REDIS_BASEDIR}/etc/redis.conf" | awk '{print $2}'
}

########################
# Set a configuration setting value
# Globals:
#   REDIS_BASEDIR
# Arguments:
#   $1 - key
#   $2 - value
# Returns:
#   None
#########################
redis_conf_set() {
    # TODO: improve this. Substitute action?
    local key="${1:?missing key}"
    local value="${2:-}"

    # Sanitize inputs
    value="${value//\\/\\\\}"
    value="${value//&/\\&}"
    value="${value//\?/\\?}"
    [[ "$value" = "" ]] && value="\"$value\""

    sed -i "s?^#*\s*$key .*?$key $value?g" "${REDIS_BASEDIR}/etc/redis.conf"
}

########################
# Unset a configuration setting value
# Globals:
#   REDIS_BASEDIR
# Arguments:
#   $1 - key
# Returns:
#   None
#########################
redis_conf_unset() {
    # TODO: improve this. Substitute action?
    local key="${1:?missing key}"
    sed -i "s?^\s*$key .*??g" "${REDIS_BASEDIR}/etc/redis.conf"
}

########################
# Get Redis version
# Globals:
#   REDIS_BASEDIR
# Arguments:
#   None
# Returns:
#   Redis versoon
#########################
redis_version() {
    "${REDIS_BASEDIR}/bin/redis-cli" --version | grep -E -o "[0-9]+.[0-9]+.[0-9]+"
}

########################
# Get Redis major version
# Globals:
#   REDIS_BASEDIR
# Arguments:
#   None
# Returns:
#   Redis major version
#########################
redis_major_version() {
    redis_version | grep -E -o "^[0-9]+"
}

########################
# Check if redis is running
# Globals:
#   REDIS_BASEDIR
# Arguments:
#   None
# Returns:
#   Boolean
#########################
is_redis_running() {
    local pid
    pid="$(get_pid_from_file "$REDIS_BASEDIR/tmp/redis.pid")"

    if [[ -z "$pid" ]]; then
        false
    else
        is_service_running "$pid"
    fi
}

########################
# Stop Redis
# Globals:
#   REDIS_*
# Arguments:
#   None
# Returns:
#   None
#########################
redis_stop() {
    local pass
    local port
    local args

    ! is_redis_running && return
    pass="$(redis_conf_get "requirepass")"
    port="$(redis_conf_get "port")"

    [[ -n "$pass" ]] && args+=("-a" "\"$pass\"")
    [[ "$port" != "0" ]] && args+=("-p" "$port")

    debug "Stopping Redis..."
    if am_i_root; then
        gosu "$REDIS_DAEMON_USER" "${REDIS_BASEDIR}/bin/redis-cli" "${args[@]}" shutdown
    else
        "${REDIS_BASEDIR}/bin/redis-cli" "${args[@]}" shutdown
    fi
    local counter=5
    while is_redis_running ; do
        if [[ "$counter" -ne 0 ]]; then
            break
        fi
        sleep 1;
        counter=$((counter - 1))
    done
}

########################
# Start redis and wait until it's ready
# Globals:
#   REDIS_*
# Arguments:
#   None
# Returns:
#   None
#########################
redis_start() {
    is_redis_running && return
    debug "Starting Redis..."
    if am_i_root; then
        gosu "$REDIS_DAEMON_USER" "${REDIS_BASEDIR}/bin/redis-server" "${REDIS_BASEDIR}/etc/redis.conf"
    else
        "${REDIS_BASEDIR}/bin/redis-server" "${REDIS_BASEDIR}/etc/redis.conf"
    fi
    local counter=3
    while ! is_redis_running ; do
        if [[ "$counter" -ne 0 ]]; then
            break
        fi
        sleep 1;
        counter=$((counter - 1))
    done
}

########################
# Load global variables used on Redis configuration.
# Globals:
#   REDIS_*
# Arguments:
#   None
# Returns:
#   Series of exports to be used as 'eval' arguments
#########################
redis_env() {
    cat <<"EOF"
export REDIS_BASEDIR="/opt/bitnami/redis"
export REDIS_EXTRAS_DIR="/opt/bitnami/extra/redis"
export REDIS_VOLUME="/bitnami/redis"
export REDIS_TEMPLATES_DIR="${REDIS_EXTRAS_DIR}/templates"
export REDIS_TMPDIR="${REDIS_BASEDIR}/tmp"
export REDIS_LOGDIR="${REDIS_BASEDIR}/logs"
export PATH="${REDIS_BASEDIR}/bin:$PATH"
export REDIS_DAEMON_USER="redis"
export REDIS_DAEMON_GROUP="redis"
export REDIS_SENTINEL_HOST="${REDIS_SENTINEL_HOST:-}"
export REDIS_SENTINEL_MASTER_NAME="${REDIS_SENTINEL_MASTER_NAME:-}"
export REDIS_SENTINEL_PORT_NUMBER="${REDIS_SENTINEL_PORT_NUMBER:-26379}"
EOF
    if [[ -f "$REDIS_PASSWORD_FILE" ]]; then
        cat <<"EOF"
export REDIS_PASSWORD="$(< "${REDIS_PASSWORD_FILE}")"
EOF
    fi
    if [[ -f "$REDIS_MASTER_PASSWORD_FILE" ]]; then
        cat <<"EOF"
export REDIS_MASTER_PASSWORD="$(< "${REDIS_MASTER_PASSWORD_FILE}")"
EOF
    fi
}

########################
# Validate settings in REDIS_* env vars.
# Globals:
#   REDIS_*
# Arguments:
#   None
# Returns:
#   None
#########################
redis_validate() {
    debug "Validating settings in REDIS_* env vars.."

    # Auxiliary functions
    empty_password_enabled_warn() {
        warn "You set the environment variable ALLOW_EMPTY_PASSWORD=${ALLOW_EMPTY_PASSWORD}. For safety reasons, do not use this flag in a production environment."
    }
    empty_password_error() {
        error "The $1 environment variable is empty or not set. Set the environment variable ALLOW_EMPTY_PASSWORD=yes to allow the container to be started with blank passwords. This is recommended only for development."
        exit 1
    }

    if is_boolean_yes "$ALLOW_EMPTY_PASSWORD"; then
        empty_password_enabled_warn
    else
        [[ -z "$REDIS_PASSWORD" ]] && empty_password_error REDIS_PASSWORD
    fi
    if [[ -n "$REDIS_REPLICATION_MODE" ]]; then
        if [[ "$REDIS_REPLICATION_MODE" =~ ^(slave|replica)$ ]]; then
            if [[ -n "$REDIS_MASTER_PORT_NUMBER" ]]; then
                if ! err=$(validate_port "$REDIS_MASTER_PORT_NUMBER"); then
                    error "An invalid port was specified in the environment variable REDIS_MASTER_PORT_NUMBER: $err"
                    exit 1
                fi
            fi
            if ! is_boolean_yes "$ALLOW_EMPTY_PASSWORD" && [[ -z "$REDIS_MASTER_PASSWORD" ]]; then
                empty_password_error REDIS_MASTER_PASSWORD
            fi
        elif [[ "$REDIS_REPLICATION_MODE" != "master" ]]; then
            error "Invalid replication mode. Available options are 'master/replica'"
            exit 1
        fi
    fi
}



########################
# Configure Redis replication
# Globals:
#   REDIS_BASEDIR
# Arguments:
#   $1 - Replication mode
# Returns:
#   None
#########################
redis_configure_replication() {
    info "Configuring replication mode..."

    if [[ "$REDIS_REPLICATION_MODE" = "master" ]]; then
        if [[ -n "$REDIS_PASSWORD" ]]; then
            redis_conf_set masterauth "$REDIS_PASSWORD"
        fi
    elif [[ "$REDIS_REPLICATION_MODE" =~ ^(slave|replica)$ ]]; then
        if [[ -n "$REDIS_SENTINEL_HOST" ]]; then
            REDIS_SENTINEL_INFO=($(redis-cli -h "$REDIS_SENTINEL_HOST" -p "$REDIS_SENTINEL_PORT_NUMBER" sentinel get-master-addr-by-name "$REDIS_SENTINEL_MASTER_NAME"))
            REDIS_MASTER_HOST=${REDIS_SENTINEL_INFO[0]}
            REDIS_MASTER_PORT_NUMBER=${REDIS_SENTINEL_INFO[1]}
        fi
        wait-for-port --host "$REDIS_MASTER_HOST" "$REDIS_MASTER_PORT_NUMBER"
        [[ -n "$REDIS_MASTER_PASSWORD" ]] && redis_conf_set masterauth "$REDIS_MASTER_PASSWORD"
        # Starting with Redis 5, use 'replicaof' instead of 'slaveof'. Maintaining both for backward compatibility
        local parameter="replicaof"
        [[ $(redis_major_version) -lt 5 ]] && parameter="slaveof"
        redis_conf_set "$parameter" "$REDIS_MASTER_HOST $REDIS_MASTER_PORT_NUMBER"
    fi
}

########################
# Disable Redis command(s)
# Globals:
#   REDIS_BASEDIR
# Arguments:
#   $1 - Array of commands to disable
# Returns:
#   None
#########################
redis_disable_unsafe_commands() {
    # The current syntax gets a comma separated list of commands, we split them
    # before passing to redis_disable_unsafe_commands
    read -r -a disabledCommands <<< "$(tr ',' ' ' <<< "$REDIS_DISABLE_COMMANDS")"
    debug "Disabling commands: ${disabledCommands[*]}"
    for cmd in "${disabledCommands[@]}"; do
        if grep -E -q "^\s*rename-command\s+$cmd\s+\"\"\s*$" "${REDIS_BASEDIR}/etc/redis.conf"; then
            debug "$cmd was already disabled"
            continue
        fi
        echo "rename-command $cmd \"\"" >> "$REDIS_BASEDIR/etc/redis.conf"
    done
}

########################
# Ensure Redis is initialized
# Globals:
#   REDIS_*
# Arguments:
#   None
# Returns:
#   None
#########################
redis_initialize() {
    info "Initializing Redis..."

    # User injected custom configuration
    if [[ -e "$REDIS_BASEDIR/etc/redis.conf" ]]; then
        if [[ -e "$REDIS_BASEDIR/etc/redis-default.conf" ]]; then
            rm "${REDIS_BASEDIR}/etc/redis-default.conf"
        fi
    else
        debug "Ensuring expected directories/files exist..."
        for dir in "${REDIS_VOLUME}/data" "${REDIS_BASEDIR}/tmp" "${REDIS_LOGDIR}"; do
            ensure_dir_exists "$dir"
            if am_i_root; then
                chown "$REDIS_DAEMON_USER:$REDIS_DAEMON_GROUP" "$dir"
            fi
        done
        mv "$REDIS_BASEDIR/etc/redis-default.conf" "$REDIS_BASEDIR/etc/redis.conf"

        # Redis config
        debug "Setting Redis config file..."
        redis_conf_set dir "${REDIS_VOLUME}/data"
        redis_conf_set logfile "" # Log to stdout
        redis_conf_set pidfile "${REDIS_BASEDIR}/tmp/redis.pid"
        redis_conf_set daemonize yes
        redis_conf_set bind 0.0.0.0 # Allow remote connections
        # Enable AOF https://redis.io/topics/persistence#append-only-file
        # Leave default fsync (every second)
        redis_conf_set appendonly yes
        # Disable RDB persistence, AOF persistence already enabled.
        # Ref: https://redis.io/topics/persistence#interactions-between-aof-and-rdb-persistence
        redis_conf_set save ""
        if [[ -n "$REDIS_PASSWORD" ]]; then
            redis_conf_set requirepass "$REDIS_PASSWORD"
        else
            redis_conf_unset requirepass
        fi
        if [[ -n "$REDIS_DISABLE_COMMANDS" ]]; then
            redis_disable_unsafe_commands
        fi
        # Configure Replication mode
        if [[ -n "$REDIS_REPLICATION_MODE" ]]; then
            redis_configure_replication
        fi
    fi
}
