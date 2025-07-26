function Invoke-Ezstore {

    param(
        [Parameter(Mandatory=$true,Position=0)]
        [string] $Path,
        [Parameter(Mandatory=$true,Position=1)]
        [string] $Arguments
    )

    $global:LASTEXITCODE = $null;

    $Path = [IO.Path]::Combine((Get-Location).Path, $Path).Replace('\.','');

    try {
        Push-Location;
        Set-Location $Path;

        $Output = Invoke-Expression -Command ".\ezstore.exe $Arguments 2>&1";

        $ExitCode = $global:LASTEXITCODE;
        $global:LASTEXITCODE = $null;

        return [string[]]($Output -split [Environment]::NewLine), [int]$ExitCode;
    } finally {
        Pop-Location;
    }

}

function Install-ModuleSafe {

    param(
        [Parameter(Mandatory=$true)]
        [string] $Name,
        [AllowEmptyString()]
        [string] $Version
    )

    & {
        $ProgressPreference = 'Ignore';

        if ( "" -eq $Version ) {
            Install-Module -Name $Name -SkipPublisherCheck -Force 3>$null;
        } else {
            if ( $null -eq (Get-Module -ListAvailable -Name $Name | Where-object Version -ge $Version) ) {
                Install-Module -Name $Name -MinimumVersion $Version -SkipPublisherCheck -Force 3>$null;
            }
        }
    }

}

# See https://github.com/PowerShell/PowerShell/issues/13138#issuecomment-1820195503
function Import-ModuleSafe {

    param(
        [Parameter(Mandatory=$true)]
        [string] $Name,
        [AllowEmptyString()]
        [string] $Version
    )

    & {
        $ProgressPreference = 'Ignore';

        if ( "" -eq $Version ) {
            Import-Module -Name $Name -UseWindowsPowerShell 3>$null;
        } else {
            if ( $null -eq (Get-Module -ListAvailable -Name $Name | Where-object Version -ge $Version) ) {
                Import-Module -Name $Name -MinimumVersion $Version -UseWindowsPowerShell 3>$null;
            }
        }
    }

}
