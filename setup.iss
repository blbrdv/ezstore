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
Source: "output\bin\ezstore.exe"; DestDir: "{app}/bin"; Flags: ignoreversion
Source: "output\README.txt"; DestDir: "{app}"; Flags: ignoreversion
Source: "output\update.ps1"; DestDir: "{app}"; Flags: ignoreversion


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
