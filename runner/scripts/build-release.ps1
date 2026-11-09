param(
    [Parameter()]
    [string] $Version = "0.0.0-dev",

    [Parameter()]
    [string] $DistDir = "dist"
)

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Invoke-NativeCommand {
    param(
        [Parameter(Mandatory = $true)]
        [string] $Command,

        [Parameter()]
        [string[]] $Arguments = @(),

        [Parameter()]
        [string] $WorkingDirectory = $PWD.Path
    )

    Push-Location $WorkingDirectory
    try {
        & $Command @Arguments
        if ($LASTEXITCODE -ne 0) {
            throw "Command failed with exit code ${LASTEXITCODE}: $Command $($Arguments -join ' ')"
        }
    }
    finally {
        Pop-Location
    }
}

function Build-GoBinary {
    param(
        [Parameter(Mandatory = $true)]
        [string] $ModuleDir,

        [Parameter(Mandatory = $true)]
        [string] $OutputPath,

        [Parameter(Mandatory = $true)]
        [string] $GOOS,

        [Parameter(Mandatory = $true)]
        [string] $GOARCH
    )

    $previousGOOS = $env:GOOS
    $previousGOARCH = $env:GOARCH
    try {
        $env:GOOS = $GOOS
        $env:GOARCH = $GOARCH
        Invoke-NativeCommand -Command "go" -Arguments @("build", "-trimpath", "-ldflags", "-s -w", "-o", $OutputPath, ".") -WorkingDirectory $ModuleDir
    }
    finally {
        $env:GOOS = $previousGOOS
        $env:GOARCH = $previousGOARCH
    }
}

function Write-Checksums {
    param(
        [Parameter(Mandatory = $true)]
        [string] $Directory
    )

    $checksumFile = Join-Path $Directory "SHA256SUMS.txt"
    $lines = Get-ChildItem -Path $Directory -File |
        Where-Object { $_.Name -ne "SHA256SUMS.txt" } |
        Sort-Object Name |
        ForEach-Object {
            $hash = Get-FileHash -Algorithm SHA256 -LiteralPath $_.FullName
            "$($hash.Hash.ToLowerInvariant())  $($_.Name)"
        }

    [System.IO.File]::WriteAllLines($checksumFile, $lines, [System.Text.UTF8Encoding]::new($false))
}

$repoRoot = Split-Path -Parent $PSScriptRoot
$resolvedRepoRoot = [System.IO.Path]::GetFullPath($repoRoot)
$resolvedDistDir = [System.IO.Path]::GetFullPath((Join-Path $resolvedRepoRoot $DistDir))
$repoRootPrefix = $resolvedRepoRoot.TrimEnd([System.IO.Path]::DirectorySeparatorChar, [System.IO.Path]::AltDirectorySeparatorChar) + [System.IO.Path]::DirectorySeparatorChar

if (-not $resolvedDistDir.StartsWith($repoRootPrefix, [System.StringComparison]::OrdinalIgnoreCase)) {
    throw "Refusing to write release artifacts outside the repository: $resolvedDistDir"
}

if (Test-Path $resolvedDistDir) {
    Remove-Item -LiteralPath $resolvedDistDir -Recurse -Force
}
New-Item -ItemType Directory -Force -Path $resolvedDistDir | Out-Null

$assinadorModule = Join-Path $resolvedRepoRoot "assinador-cli"
$simuladorModule = Join-Path $resolvedRepoRoot "simulador-cli"
$javaModule = Join-Path $resolvedRepoRoot "assinador-verificador"

& (Join-Path $javaModule "build.ps1")
$javaDistDir = Join-Path (Join-Path $javaModule "build") "dist"
Copy-Item -LiteralPath (Join-Path $javaDistDir "assinador-verificador.jar") -Destination (Join-Path $resolvedDistDir "assinador-verificador-$Version.jar")

$targets = @(
    @{ GOOS = "windows"; GOARCH = "amd64"; Suffix = "windows-amd64.exe" },
    @{ GOOS = "linux"; GOARCH = "amd64"; Suffix = "linux-amd64" },
    @{ GOOS = "darwin"; GOARCH = "amd64"; Suffix = "macos-amd64" }
)

foreach ($target in $targets) {
    Build-GoBinary `
        -ModuleDir $assinadorModule `
        -OutputPath (Join-Path $resolvedDistDir "assinatura-$Version-$($target.Suffix)") `
        -GOOS $target.GOOS `
        -GOARCH $target.GOARCH

    Build-GoBinary `
        -ModuleDir $simuladorModule `
        -OutputPath (Join-Path $resolvedDistDir "simulador-$Version-$($target.Suffix)") `
        -GOOS $target.GOOS `
        -GOARCH $target.GOARCH
}

Write-Checksums -Directory $resolvedDistDir

Write-Host "Artefatos de release gerados em $resolvedDistDir"
