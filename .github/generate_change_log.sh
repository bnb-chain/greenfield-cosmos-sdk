#!/usr/bin/env bash
declare change_log_file="./CHANGELOG.md"
declare version="## $@"
declare version_prefix="## v"
declare start=0
declare CHANGE_LOG=""

while read line; do
    if [[ $line == *"$version"* ]]; then
        start=1
        continue
    fi
    if [[ $line == *"$version_prefix"* ]] && [ $start == 1 ]; then
        break;
    fi
    if [ $start == 1 ]; then
        CHANGE_LOG+="$line\n"
    fi
done < "${change_log_file}"

OUTPUT=$(cat <<-END
${CHANGE_LOG}\n
END
)

echo -e "${OUTPUT}"