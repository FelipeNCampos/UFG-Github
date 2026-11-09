Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

function Invoke-NativeCommand {
    param(
        [Parameter(Mandatory = $true)]
        [string] $Command,

        [Parameter()]
        [string[]] $Arguments = @()
    )

    & $Command @Arguments
    if ($LASTEXITCODE -ne 0) {
        throw "Command failed with exit code ${LASTEXITCODE}: $Command $($Arguments -join ' ')"
    }
}

$projectRoot = Split-Path -Parent $MyInvocation.MyCommand.Path
$resolvedProjectRoot = [System.IO.Path]::GetFullPath($projectRoot)
$buildDir = Join-Path $resolvedProjectRoot "build"

if (Test-Path $buildDir) {
    $resolvedBuildDir = (Resolve-Path $buildDir).Path
    $projectRootPrefix = $resolvedProjectRoot.TrimEnd([System.IO.Path]::DirectorySeparatorChar, [System.IO.Path]::AltDirectorySeparatorChar) + [System.IO.Path]::DirectorySeparatorChar
    if (-not $resolvedBuildDir.StartsWith($projectRootPrefix, [System.StringComparison]::OrdinalIgnoreCase)) {
        throw "Refusing to clean an unexpected directory: $resolvedBuildDir"
    }

    Remove-Item -LiteralPath $resolvedBuildDir -Recurse -Force
}

$classesDir = Join-Path (Join-Path $buildDir "classes") "main"
$distDir = Join-Path $buildDir "dist"
$tmpDir = Join-Path $buildDir "tmp"
$manifestFile = Join-Path $tmpDir "manifest.mf"
$jarFile = Join-Path $distDir "assinador-verificador.jar"
$mainSourceRoot = Join-Path $resolvedProjectRoot "src\main\java"

New-Item -ItemType Directory -Force -Path $classesDir, $distDir, $tmpDir | Out-Null

$mainSources = Get-ChildItem -Path $mainSourceRoot -Recurse -Filter *.java | ForEach-Object { $_.FullName }
if (-not $mainSources) {
    throw "No Java sources were found under $mainSourceRoot"
}

Invoke-NativeCommand -Command "javac" -Arguments (@("-encoding", "UTF-8", "-d", $classesDir) + $mainSources)

[System.IO.File]::WriteAllText(
    $manifestFile,
    "Manifest-Version: 1.0`nMain-Class: br.ufg.runner.assinador.AssinadorApplication`n`n",
    [System.Text.UTF8Encoding]::new($false)
)

Invoke-NativeCommand -Command "jar" -Arguments @("--create", "--file", $jarFile, "--manifest", $manifestFile, "-C", $classesDir, ".")

Write-Host "JAR criado em $jarFile"
