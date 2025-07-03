function Install {

    [CmdletBinding()]
    [OutputType([string])]
    param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Id,
        [Parameter(Mandatory=$true,Position=1)]
        [string]$Version,
        [Parameter(Position=2)]
        [int]$ExpectedCode
    )

    if ( $null -eq $ExpectedCode ) {
        $ExpectedCode = 0;
    }

    $global:LASTEXITCODE = $null;

    $Cmd = ".\output\bin\ezstore.exe install $Id --ver $Version --verbosity d";
    Write-Host $Cmd;
    $Output = Invoke-Expression -Command "$Cmd 2>&1";
    foreach ($Line in ($Output -split "\n")) {
        Write-Host $Line;
    }

    if ( $ExpectedCode -ne $global:LASTEXITCODE ) {
        throw "expected exit code $ExpectedCode, actual $($global:LASTEXITCODE)"
    }

    return $Output;

}

$ColorRegexp = '\x1b\[[0-9;]*m';
