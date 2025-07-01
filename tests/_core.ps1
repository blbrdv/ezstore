function Install {

    [CmdletBinding()]
    [OutputType([string])]
    param (
        [Parameter(Mandatory=$true,Position=0)]
        [string]$Id,
        [Parameter(Mandatory=$true,Position=1)]
        [string]$Version
    )

    $Cmd = ".\output\bin\ezstore.exe install $Id --ver $Version --verbosity d";
    Write-Output $Cmd;
    $Output = Invoke-Expression -Command "$Cmd 2>&1";
    Write-Output $Output;

    return $Output;

}

$ColorRegexp = '\x1b\[[0-9;]*m';

$Targets = @{

    # Tree CLI app
    # no dependencies
    "9mvsm3j7zj7c" = @{
        Name = "PeterEtelej.TreeCLI"
        Version = "1.1.0.0"
    }

    # Wikipedia app
    # one dependency
    "9wzdncrfhwm4" = @{
        Name = "WikimediaFoundation.Wikipedia"
        Version = "1.0.0.0"
    }

    # WhatsApp app
    # 5 dependencies
    "9nksqgp7f2nh" = @{
        Name = "5319275A.WhatsAppDesktop"
        Version = "2.2524.4.0"
    }

};
