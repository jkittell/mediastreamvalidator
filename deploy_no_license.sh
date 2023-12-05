#!/usr/bin/env bash
# -*- coding:utf-8 -*-
#
# Copyright (C) 2023 Apple, Inc.
#

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" >/dev/null && pwd)"

# Script commands
COMMAND_INSTALL="INSTALL"
COMMAND_UNINSTALL="UNINSTALL"

#MARK: Defaults
LIST_ROW_LEN=3
DEPENDENCY_DIR="Dependencies"

# Install by default. (Without -i option)
COMMAND=${COMMAND_INSTALL}

LICENSE_FILE="License.txt"

# MARK: Logging
exec 3>&2
verbosity=3
silent_lvl=0
crt_lvl=1
err_lvl=2
wrn_lvl=3
inf_lvl=4
dbg_lvl=5
trc_lvl=6
BUILD_CONFIG="Deploy"
notify() { log $silent_lvl "$1"; } # Always prints
critical() { log $crt_lvl "CRITICAL: $1"; }
error() { log $err_lvl "ERROR: $1"; }
warn() { log $wrn_lvl "WARNING: $1"; }
inf() { log $inf_lvl "INFO: $1"; } # "info" is already a command
debug() { log $dbg_lvl "DEBUG: $1"; }
trace() { log $trc_lvl "TRACE: $1"; }
log ()
{
	if [ $verbosity -ge $1 ]; then
		echo "$2" >&3
	fi
}

# MARK: xtrace
shopt -s expand_aliases
_xtrace() {
	if [ $verbosity -ge $trc_lvl ]; then
		case $1 in
			on)
				set -x
				;;
			off)
				set +x
				;;
		esac
	fi
}

alias xtrace='{ _xtrace $(cat); } 2>/dev/null <<<'

# MARK: usage()
usage()
{
cat << EOF
usage: $0 [options]

options:
   -e          : erase (uninstall) packages
   -h          : help
   -i          : install packages
   -v          : verbose
EOF
}

function dislpay_license()
{
  if [ -f ${LICENSE_FILE} ]; then
    cat ${LICENSE_FILE} | more
    read -p "Do you accept this license? [Y/N]: " -n 1 -r
    echo    # (optional) move to a new line
    if [[ ! $REPLY =~ ^[Yy]$ ]]
    then
      notify "Exiting installation as license is not accepted"
      exit 1
    fi
  fi
}

function install_packages()
{
	local length=${#DEPENDENCY_LIST[*]}
	for (( i=0; i<length; i += LIST_ROW_LEN ));	do
		local package=${DEPENDENCY_LIST[$i]}
		local install_options="-Uvh ${DEPENDENCY_LIST[$((i + 1))]}"

		debug "index: $i name: ${package} install_options: ${install_options}"

		local matching_package_name="$(ls -t1 | grep "$package-.*\.rpm" | head -n 1)"; shift
		notify ""
		notify "Installing ${package}..."

		debug "matching_package_name ${matching_package_name}"

		if [[ -n "${matching_package_name}" ]]; then
			debug "rpm ${install_options} ${matching_package_name}"
			rpm ${install_options} "${matching_package_name}"
		fi
	done
}

function remove_packages()
{
    local length=${#DEPENDENCY_LIST[*]}
    # Remove in reverse order
    # First index from the back
    local first_index=$((length - LIST_ROW_LEN))
    local prev_package=""
    for ((i=first_index ; i>=0 ; i -= LIST_ROW_LEN)); do
        local installed_packages=( $(rpm -qa) )
        local package=${DEPENDENCY_LIST[$i]}
        # e.g. ifs4l-release --> ifs4l
        local sanitized_package=${package%-*}
        local remove_options="-ev ${DEPENDENCY_LIST[$((i + 2))]}"

        if [[ ${prev_package} == ${sanitized_package} ]]; then
            continue
        fi

        debug "index: $i name: ${package} remove_options: ${remove_options}"

        local matching_package_name="$(printf '%s\n' ${installed_packages[@]} | grep -i "${sanitized_package,,}-.*" | head -n 1)"; shift
        notify ""
        notify "Removing ${package}..."

        if [[ -n "${matching_package_name}" ]]; then
            debug "rpm ${remove_options} `rpm -qa | grep -i ${sanitized_package}-.*`"
            rpm ${remove_options} `rpm -qa | grep -i ${sanitized_package}-.*`
            prev_package=${sanitized_package}
        else
            warn "No ${package} package installed. Skipping."
        fi
    done
}

# Listed in installation order. The order is important.
# List is used for removal in reverse order.
# List format: 'package_name' '<install_options>' '<remove_options>'
#
# install_options - additional options used during package installation. Ex: --force --nodeps
# remove_options - additional options used during package removal.

DEPENDENCY_LIST=(
	'libdispatch' '' ''

	'foundation' '' ''

	'icucore' '' ''

	'ifs4l-release' '' ''

	'CoreVideo-release' '' ''

	'calinuxbase-release' '' ''

	'caulk-release' '' ''

	'audiotoolboxcore-release' '' ''

	'libblocksruntime' '' ''

	'libbsd' '' ''

	'audiocodecs-hls-public-release' '' ''
)

pushd "${SCRIPT_DIR}" > /dev/null

# MARK: - Collect Options
while getopts dehirv option
do
	case $option in
		e)
			COMMAND=${COMMAND_UNINSTALL}
			;;
		h)
			usage
			exit
			;;
		i)
			COMMAND=${COMMAND_INSTALL}
			;;
		v)
			if [ $verbosity -lt $dbg_lvl ];then
				verbosity=$dbg_lvl
			fi
			;;
		?)
			usage
			exit
			;;
	esac
done

debug "DEPENDENCY_LIST: ${DEPENDENCY_LIST[*]}"
# Check dependency list
length=${#DEPENDENCY_LIST[*]}
if [ $((length % $LIST_ROW_LEN)) -ne 0 -o  $length -eq 0 ]; then
	usage "Config error: incomplete dependency list"
fi

if [ -z ${COMMAND} ]; then
	echo "No command specified."
	usage
	exit
fi

case "${COMMAND}" in
	${COMMAND_INSTALL} )
        # display license
        # dislpay_license

		notify "Installing packages..."
		if [ ! -d "${DEPENDENCY_DIR}/${BUILD_CONFIG}" ]
		then
			echo "${DEPENDENCY_DIR}/${BUILD_CONFIG} does not exist..."
			exit
		fi

		# Install packages
		pushd "${DEPENDENCY_DIR}/${BUILD_CONFIG}" > /dev/null
		install_packages
		popd > /dev/null

		# Install tools
		notify ""
		notify "Installing hlstools..."
		if [ -f hlstools-*.rpm ]; then
			debug "hlstools package found."
			rpm -Uvh hlstools-*.rpm
		else
			warn "No hlstools package found. Skipping."
		fi

		notify ""
		notify "Done"
		;;

	${COMMAND_UNINSTALL} )
		notify "Removing packages..."

		hlstools_package_name="$(rpm -qa | grep "hlstools-.*" | head -n 1)"

		# Remove tools
		notify ""
		notify "Removing hlstools..."
		if [ -z ${hlstools_package_name} ]; then
			warn "No hlstools package found. Skipping."
		else
			rpm -ev "hlstools"
		fi

		# Remove packages
		pushd "${DEPENDENCY_DIR}" > /dev/null
		remove_packages
		popd > /dev/null

		notify ""
		notify "Done"
		;;
	?)
		error "Undefined command"
		exit -1 ;;
esac

popd > /dev/null
