function Install {

    [CmdletBinding()]
    [OutputType([string])]
    param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Id,
        [Parameter(Mandatory=$true,Position=1)]
        [string]$Version
    )

    $Cmd = ".\output\bin\ezstore.exe install $Id --ver $Version --verbosity m";
    Write-Host $Cmd;
    $Output = Invoke-Expression -Command "$Cmd 2>&1";
    Write-Host $Output;

    return $Output;

}

$ColorRegexp = '\x1b\[[0-9;]*m';

$Data = @{
    Id = "9mvsm3j7zj7c"
    Name = "PeterEtelej.TreeCLI"
    Version = "1.1.0.0"
};
