#require -Module Pester

function Assert-Fail {

    param(
        [Parameter(Mandatory=$true,Position=0)]
        [string] $Message
    )

    return New-Object PSObject -Property @{
        Succeeded      = $false
        FailureMessage = $Message
    }
}

function Get-ErrorMessage {

    param(
        [Parameter(Mandatory=$true)]
        [bool] $Negate,
        [Parameter(Mandatory=$true)]
        [int] $LineNum,
        [Parameter(Mandatory=$true)]
        [string] $Text
    )

    if ( $Negate ) {
        $ComparisonMessage = "is $Text";
    } else {
        $ComparisonMessage = "is not $Text";
    }

    return "Line at index $LineNum $($ComparisonMessage):";

}

function Should-AssertOutput {

    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
        "PSUseApprovedVerbs", "",
        Justification="'Should' is a verb used in Pester module."
    )]
    param(
        $ActualValue,

        [switch] $Negate,

        [int] $LineNum,

        [AllowEmptyString()]
        [string] $ShouldMatch,

        [AllowEmptyString()]
        [string] $ShouldBeExactly,

        [scriptblock] $Script
    )

    if ( $null -ne $ActualValue ) {
        $Output = [string[]]($ActualValue | ForEach-Object { [string]$_ });
    } else {
        $Output = [string[]]@()
    }
    $Count = $Output.Count;

    if ( $Count -eq 0 ) {
        return Assert-Fail "Output is empty.";
    }

    $Line = $Output[$LineNum];

    if ( $null -ne $Script ) {
        $Line = $Line | ForEach-Object $Script;
    }

    if ( "" -ne $ShouldMatch ) {
        $Pass = $Line -match $ShouldMatch;
        $ErrorMessage = Get-ErrorMessage -Negate $Negate -LineNum $LineNum -Text "match";
    } elseif ( "" -ne $ShouldBeExactly ) {
        $Pass = $Line.Equals($ShouldBeExactly);
        $ErrorMessage = Get-ErrorMessage -Negate $Negate -LineNum $LineNum -Text "equal to";
    } else {
        throw "Either -ShouldMatch or -ShouldBeExactly param must be provided"
    }

    if ( $Negate ) {
        $Pass = -not $Pass;
    }

    if ( -not $Pass ) {
        if ( $LineNum -lt 0 ) {
            $DesiredIndex = $Count + $LineNum;
        } else {
            $DesiredIndex = $LineNum;
        }

        $FullErrorMessage = @(
            $ErrorMessage
        )

        for ($i = 0; $i -lt $Count; $i++) {
            if ( $i -eq $DesiredIndex ) {
                $Mark = ">";
            } else {
                $Mark = " ";
            }

            $FullErrorMessage += "$Mark '$($Output[$i])'";
            if ( $i -eq $DesiredIndex ) {
                if ( "" -ne $ShouldMatch ) {
                    $FullErrorMessage += "> '$ShouldMatch'";
                } else {
                    $FullErrorMessage += "> '$ShouldBeExactly'";
                }
            }
        }

        return Assert-Fail ($FullErrorMessage -join [Environment]::NewLine);
    }

    return New-Object PSObject -Property @{
            Succeeded      = $true
            FailureMessage = $null
        };

}

function Should-EqualOutput {

    [System.Diagnostics.CodeAnalysis.SuppressMessageAttribute(
        "PSUseApprovedVerbs", "",
        Justification="'Should' is a verb used in Pester module."
    )]
    param(
        $ActualValue,
        [string[]] $ExpectedValue
    )

    if ( $null -ne $ActualValue ) {
        $Output = [string[]]($ActualValue | ForEach-Object { [string]$_ });
    } else {
        $Output = [string[]]@();
    }
    $ActualCount = $Output.Count;
    $ExpectedCount = $ExpectedValue.Count;

    $ErrorMessage = [string[]]@();

    if ( $ActualCount -ne $ExpectedCount ) {
        $ErrorMessage += "Expected output length $ExpectedCount, but got $ActualCount";

        for ( $i = 0; $i -lt $ActualCount; $i++ ) {
            $ErrorMessage += "$($i + 1)`t$($Output[$i])";
        }

        return Assert-Fail ($ErrorMessage -join [Environment]::NewLine);
    }

    $Pass = $true;
    $ErrorMessage += "Output is invalid:"

    for ( $i = 0; $i -lt $ActualCount; $i++ ) {
        $ActualLine = $Output[$i];
        $ExpectedLine = $ExpectedValue[$i];

        if ( $ActualLine.Equals($ExpectedLine) ) {
            $ErrorMessage += "  '$ActualLine'";
        } else {
            $ErrorMessage += "> '$ActualLine'";
            $ErrorMessage += "> '$ExpectedLine'";
            $Pass = $false;
        }
    }

    if ( $Pass ) {
        return New-Object PSObject -Property @{
               Succeeded      = $true
               FailureMessage = $null
           }
    } else {
        return Assert-Fail ($ErrorMessage -join [Environment]::NewLine);
    }

}

Add-ShouldOperator -Name AssertOutput -InternalName 'Should-AssertOutput' -Test ${function:Should-AssertOutput} `
    -SupportsArrayInput;

Add-ShouldOperator -Name EqualOutput -InternalName 'Should-EqualOutput' -Test ${function:Should-EqualOutput} `
    -SupportsArrayInput;
