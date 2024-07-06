param(
	[string]$Mods = "[util log]",
	[string]$TotalCoverageOut = "coverage/coverage.out",
	[string]$ModCoverageOut = "coverage/unit_cov.txt"
)

$SubMods = $Mods.TrimStart("[").TrimEnd("]").Split(" ")

$directoryPath = [System.IO.Path]::GetDirectoryName($TotalCoverageOut)
if (-not (Test-Path $directoryPath)) {
	New-Item -ItemType Directory -Path $directoryPath
}

Set-Content -Path $TotalCoverageOut -Value "mode: atomic"

Write-Host -NoNewline -ForegroundColor Green "info: "
Write-Host "merge sub pkgs [$($SubMods -join ", ")] test coverage to $TotalCoverageOut"

for ($i = 0; $i -lt $SubMods.Length; $i++) {
	$modCovOutFile = Join-Path -Path $SubMods[$i] -ChildPath $ModCoverageOut
	if (Test-Path $modCovOutFile) {
		Get-Content $modCovOutFile | Where-Object {
			$_ -NotMatch "mode:"
		} | Add-Content $TotalCoverageOut
	}
}

if (Test-Path "./internal/coverage/integration/coverage.out") {
    Get-Content "./internal/coverage/integration/coverage.out" |
    Select-String -Pattern "mode:|internal/example" -NotMatch |
    Add-Content -Path "$TotalCoverageOut"
}
