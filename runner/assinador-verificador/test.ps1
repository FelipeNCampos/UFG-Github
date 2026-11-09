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
$buildScript = Join-Path $resolvedProjectRoot "build.ps1"

& $buildScript

$mainClassesDir = Join-Path (Join-Path $resolvedProjectRoot "build") "classes"
$mainClassesDir = Join-Path $mainClassesDir "main"
$testClassesDir = Join-Path (Join-Path $resolvedProjectRoot "build") "classes"
$testClassesDir = Join-Path $testClassesDir "test"
$testSourceRoot = Join-Path $resolvedProjectRoot "src\test\java"

New-Item -ItemType Directory -Force -Path $testClassesDir | Out-Null

$testSources = Get-ChildItem -Path $testSourceRoot -Recurse -Filter *.java | ForEach-Object { $_.FullName }
if (-not $testSources) {
    throw "No Java tests were found under $testSourceRoot"
}

$pathSeparator = [System.IO.Path]::PathSeparator
$classPath = "$mainClassesDir$pathSeparator$testClassesDir"

Invoke-NativeCommand -Command "javac" -Arguments (@("-encoding", "UTF-8", "-cp", $mainClassesDir, "-d", $testClassesDir) + $testSources)
Invoke-NativeCommand -Command "java" -Arguments @("-cp", $classPath, "br.ufg.runner.assinador.AssinadorApplicationTest")

Write-Host "Testes concluidos com sucesso."
