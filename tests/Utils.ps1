function Invoke-Ezstore {

    param(
        [Parameter(Mandatory=$true,Position=0)]
        [string] $Path,
        [Parameter(Mandatory=$true,Position=1)]
        [string[]] $Arguments
    )

    $global:LASTEXITCODE = $null;

    $Path = [IO.Path]::Combine((Get-Location).Path, $Path).Replace('\.','');

    try {
        Push-Location;
        Set-Location $Path;

        $Output = & .\ezstore.exe @Arguments 2>&1;

        $ExitCode = $global:LASTEXITCODE;
        $global:LASTEXITCODE = $null;

        $Output = $Output -split [Environment]::NewLine;

        Write-Verbose "Exit code: $ExitCode";
        Write-Verbose "Output:";
        foreach ( $Line in $Output ) {
            Write-Verbose $Line;
        }

        return $Output, $ExitCode;
    } finally {
        Pop-Location;
    }

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

    & {
        $ProgressPreference = 'Ignore';
        $Params = @{
            Name = $Name
        };

        if ( "" -eq $Version ) {
            $Params["MinimumVersion"] = $Version;
            $List = Get-Module -ListAvailable -Name $Name | Where-object Version -ge $Version;
        } else {
            $List = Get-Module -ListAvailable -Name $Name;
        }

        if ( $null -eq $List ) {
            Install-Module @Params -SkipPublisherCheck 3>$null;
        }
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

    & {
        $ProgressPreference = 'Ignore';
        $Params = @{
            Name = $Name
        };

        if ( "" -eq $Version ) {
            $Params["MinimumVersion"] = $Version;
        }

        if ( $null -ne $Env:GITHUB_ACTIONS ) {
            Import-Module @Params -Force 3>$null;
        } elseif ( $UseWindowsPowerShell ) {
            Import-Module @Params -UseWindowsPowerShell 3>$null;
        } else {
            Import-Module @Params 3>$null;
        }
    }

}
