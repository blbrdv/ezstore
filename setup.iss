#define Name        "ezstore"
#define Publisher   "blbrdv"
#define URL         "https://github.com/blbrdv/ezstore"
#define ExeName     "ezstore.exe"


;------------------------------------------------------------------------------
[Setup]
AppId={{27F6680A-E22B-4611-A5A0-089A828C5F96}
AppName={#Name}
AppVersion={#PV}
VersionInfoVersion={#FV}
AppPublisher={#Publisher}
AppPublisherURL={#URL}
AppSupportURL={#URL}
AppUpdatesURL={#URL}

DefaultDirName={commonpf}\{#Name}
DefaultGroupName={#Name}

OutputDir=.\release
OutputBaseFileName=ezsetup

LicenseFile=LICENSE
SetupIconFile=icons\icon.ico

Compression=lzma
SolidCompression=yes
ChangesEnvironment=true


;------------------------------------------------------------------------------
[Files]
Source: "output\386\bin\ezstore.exe"; DestDir: "{app}/bin"; Check: IsI386; Flags: "ignoreversion solidbreak"
Source: "output\amd64\bin\ezstore.exe"; DestDir: "{app}/bin"; Check: IsAmd64; Flags: "ignoreversion solidbreak"
Source: "output\arm\bin\ezstore.exe"; DestDir: "{app}/bin"; Check: IsArm; Flags: "ignoreversion solidbreak"
Source: "output\arm64\bin\ezstore.exe"; DestDir: "{app}/bin"; Check: IsArm64; Flags: "ignoreversion solidbreak"
Source: "cmd\README.txt"; DestDir: "{app}"; Flags: "ignoreversion solidbreak"
Source: "cmd\update.ps1"; DestDir: "{app}"; Flags: "ignoreversion solidbreak"


;------------------------------------------------------------------------------
[Icons]
Name: "{group}\{#Name}"; Filename: "{app}\{#ExeName}"


;------------------------------------------------------------------------------
[Tasks]
Name: envPath; Description: "Add to PATH variable"


;------------------------------------------------------------------------------
#include "innosetup/environment.iss"
[Code]
procedure CurStepChanged(CurStep: TSetupStep);
begin
  if CurStep = ssPostInstall
  then EnvAddPath(ExpandConstant('{app}') + '\bin');
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  if CurUninstallStep = usPostUninstall
  then EnvRemovePath(ExpandConstant('{app}') + '\bin');
end;

function IsAmd64: Boolean;
begin
  Result := (ProcessorArchitecture = paX64);
end;

function IsArm64: Boolean;
begin
  Result := (ProcessorArchitecture = paARM64);
end;

function IsI386: Boolean;
begin
  Result := (ProcessorArchitecture = paX86);
end;

function IsArm: Boolean;
begin
  Result := (not IsAmd64) and (not IsArm64) and (not IsI386);
end;
