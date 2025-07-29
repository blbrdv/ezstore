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
        [Parameter(Mandatory=$true,Position=0)]
        [bool] $Negate,
        [Parameter(Mandatory=$true,Position=1)]
        [int] $LineNum,
        [Parameter(Mandatory=$true,Position=2)]
        [string] $Text,
        [Parameter(Mandatory=$true,Position=3)]
        [string] $Value
    )

    if ( $Negate ) {
        $ComparisonMessage = $Text;
    } else {
        $ComparisonMessage = "didn't $Text";
    }

    return "Line at index $LineNum $ComparisonMessage '$Value':";

}

function Should-AssertOutput {

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

    if ( $ActualValue.Count -lt 1 ) {
        return Assert-Fail "Output is empty.";
    }

    $Line = $ActualValue[$LineNum];

    if ( $null -ne $Script ) {
        $Line = $Line | % $Script;
    }

    if ( "" -ne $ShouldMatch ) {
        $Pass = $Line -match $ShouldMatch;
        $ErrorMessage = Get-ErrorMessage $Negate $LineNum "match" $ShouldMatch;
    } elseif ( "" -ne $ShouldBeExactly ) {
        $Pass = $Line.Equals($ShouldBeExactly);
        $ErrorMessage = Get-ErrorMessage $Negate $LineNum "equal to" $ShouldBeExactly;
    } else {
        throw "Either -ShouldMatch or -ShouldBeExactly param must be provided"
    }

    if ( $Negate ) {
        $Pass = -not $Pass;
    }

    if ( -not $Pass ) {
        if ( $LineNum -lt 0 ) {
            $DesiredIndex = $ActualValue.Count + $LineNum;
        } else {
            $DesiredIndex = $LineNum;
        }

        $FullErrorMessage = @(
            $ErrorMessage;
        )

        for ($i = 0; $i -lt $ActualValue.Count; $i++) {
            if ( $i -eq $DesiredIndex ) {
                $Mark = ">";
            } else {
                $Mark = " ";
            }

            $FullErrorMessage += "$Mark '$($ActualValue[$i])'";
        }

        return Assert-Fail ($FullErrorMessage -join [Environment]::NewLine);
    }

    return New-Object PSObject -Property @{
            Succeeded      = $true
            FailureMessage = $null
        };

}

Add-ShouldOperator -Name AssertOutput -InternalName 'Should-AssertOutput' -Test ${function:Should-AssertOutput} -SupportsArrayInput;
