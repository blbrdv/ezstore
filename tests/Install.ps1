function Invoke-EzstoreInstall {

    param(
        [Parameter(Mandatory=$true,Position=0)]
        [string] $Path,
        [Parameter(Mandatory=$true,Position=1)]
        [string] $Id,
        [Parameter(Mandatory=$true,Position=2)]
        [string] $Version
    )

    return Invoke-Ezstore $Path @("install", "$Id", "--ver", "$Version", "--verbosity", "d");
}
