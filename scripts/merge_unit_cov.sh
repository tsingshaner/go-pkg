#!/bin/bash

Mods=$1
TotalCoverageOut=$2
ModCoverageOut=$3

SubMods=(${Mods//[\[\]]/})

directoryPath=$(dirname "$TotalCoverageOut")

if [ ! -d "$directoryPath" ]; then
	mkdir -p "$directoryPath"
fi

echo "mode: atomic" >"$TotalCoverageOut"
echo -e "\e[32minfo: \e[0mmerge sub pkgs (${SubMods[*]}) test coverage to $TotalCoverageOut"

for mod in "${SubMods[@]}"; do
	modCovOutFile="$mod/$ModCoverageOut"
	grep -v "mode:" "$modCovOutFile" >>"$TotalCoverageOut"
done

for mod in "${SubMods[@]}"; do
	modCovOutFile="$mod/$ModCoverageOut"
	integrationTestCovOutFile="./internal/coverage/integration/$mod/profile.txt"

	if [ -f "$modCovOutFile" ]; then
		grep -v "mode:" "$modCovOutFile" >>"$TotalCoverageOut"
	fi


done

if [ -f "./internal/coverage/integration/coverage.out" ]; then
	grep -vE "mode:|internal/example" "./internal/coverage/integration/coverage.out" >>"$TotalCoverageOut"
fi
