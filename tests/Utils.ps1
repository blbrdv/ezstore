function Invoke-Ezstore {

    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
        "PSReviewUnusedParameter", "",
        Justification="Parameters are used below. PS linter just blind."
    )]
    param(
        [Parameter(Mandatory=$true,Position=0)]
        [string] $Path,
        [Parameter(Mandatory=$true,Position=1)]
        [string[]] $Arguments
    )

    $Path = [IO.Path]::Combine((Get-Location).Path, $Path).Replace('\.','') | Convert-Path;
    $CancelStatuses = @(
        "Stopped",
        "Blocked",
        "Suspended",
        "Disconnected"
    )

    if ( -not (Test-Path -Path $Path) ) {
        throw "Path does not exists: $Path";
    }

    $Job = Start-Job {
        $global:LASTEXITCODE = $null;
        Set-Location $using:Path;

        .\ezstore.exe @using:Arguments 2>&1;
        $global:LASTEXITCODE;
    };

    $Job | Wait-Job -Timeout 600 >$null;
    $Result = Receive-Job -Job $Job -ErrorAction "Stop";

    if ( $CancelStatuses.Contains($Job.State) ) {
        $Output = $Result;
        $ExitCode = 124;
    } else {
        $Count = ($Result | Measure-Object).Count;
        $Output = $Result[0..($Count - 2)];
        $ExitCode = $Result[-1];
    }

    return [string[]]$Output, [int]$ExitCode;

}

function Install-ModuleSafe {

    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
        "PSReviewUnusedParameter", "",
        Justification="Parameters are used below. PS linter just blind."
    )]
    param(
        [Parameter(Mandatory=$true)]
        [string] $Name,
        [AllowEmptyString()]
        [string] $Version
    )

    $Params = @{
        Name = $Name
        SkipPublisherCheck = $true
        Confirm = $false
        Force = $true
        ErrorAction = 'Stop'
    };

    if ( "" -ne $Version ) {
        $Params["MinimumVersion"] = $Version;
        $List = Get-Module -ListAvailable -Name $Name | Where-object Version -ge $Version;
        $Message = "Module $Name ($Version) successfully installed."
    } else {
        $List = Get-Module -ListAvailable -Name $Name;
        $Message = "Module $Name successfully installed."
    }

    if ( $null -eq $List ) {
        Install-Module @Params 3>$null;
        Write-Output $Message;
    } else {
        Write-Output "Module $Name already installed."
    }

}

# See https://github.com/PowerShell/PowerShell/issues/13138#issuecomment-1820195503
function Import-ModuleSafe {

    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
        "PSReviewUnusedParameter", "",
        Justification="Parameters are used below. PS linter just blind."
    )]
    param(
        [Parameter(Mandatory=$true)]
        [string] $Name,
        [AllowEmptyString()]
        [string] $Version,
        [switch] $UseWindowsPowerShell
    )

    $Params = @{
        Name = $Name
        Force = ($null -eq $Env:GITHUB_ACTIONS)
        ErrorAction = 'Stop'
    };

    if ( "" -ne $Version ) {
        $Params["MinimumVersion"] = $Version;
        $Message = "Module $Name ($Version) imported."
    } else {
        $Message = "Module $Name imported."
    }

    if ( ($PSVersionTable.PSVersion.Major -gt 5) -and $UseWindowsPowerShell ) {
        Import-Module @Params -UseWindowsPowerShell 3>$null;
    } else {
        Import-Module @Params 3>$null;
    }

    Write-Output $Message;

}

function Get-PackageFullName {
    return Get-AppxPackage | ForEach-Object { $_.PackageFullName; };
}

function Assert-PackageInstalled {

    param(
        [Parameter(Mandatory=$true)]
        [string]$Name,
        [Parameter(Mandatory=$true)]
        [string]$Version,
        [Parameter(Mandatory=$true)]
        [string]$PackageId
    )

    Get-PackageFullName | Where-Object { $_ -match "^$($Name)_$($Version)_[^_]+_[^_]*_$PackageId$" }

}
