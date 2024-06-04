# Sandboxer
<p align="center"><img src="resources/icon_transparent.png" width="200"/></p>

## Inspect objects using Trend Micro Vision One or Deep Discovery Analyzer sandbox

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

## Features

### Major Features
- Inspected objects: Files and URLs.
- Supported sandboxes: Vision One and Deep Discovery Analyzer. 
- Supported platforms: Windows and macOS.

###  Minor Features
- Auto notification and download when a new version is released.
- Store analysis results for two months (this is a configurable option).
- Show system notifications when a malicious file is detected
- Show Vision One sandbox quota
- Support HTTP proxy server including basic and NTLM authentication

Sandboxer submissions window:

<p align="center"><img src="resources/captures/submissions.png" width="300px"/></p>

## Installation
Installation and usage video:

[![Installation and usage](https://img.youtube.com/vi/beRX6YXjC4k/0.jpg)](https://www.youtube.com/watch?v=beRX6YXjC4k)

Download [latest release](https://github.com/mpkondrashin/sandboxer/releases/latest) for your platform, unpack the zip file and run setup.exe (for Windows) or SandboxerInstaller (for macOS).

<p align="center"><kbd><img src="resources/captures/page_0.png" width="500px"/></kbd></p>
<p align="center">Choose your sanbox type: <a href="https://docs.trendmicro.com/en-us/documentation/article/trend-vision-one-sandbox-analysis">Vision One</a> or <a href="https://www.trendmicro.com/en_ca/business/products/network/advanced-threat-protection/analyzer.html">Deep Discovery Analyzer</a>. Then press the "Next" button.</p>
<p align="center"><kbd><img src="resources/captures/page_5.png" width="500px"/></kbd></p>
<p align="center">If you selected Vision One on the first step, then enter Token. Learn more about <a href="https://docs.trendmicro.com/en-us/documentation/article/trend-vision-one-api-keys">API Keys</a> and <a href="https://docs.trendmicro.com/en-US/documentation/article/trend-vision-one-configuring-user-rol">Roles</a>. 
If the correct Domain value is not detected automatically, choose it from the dropdown list.</p> 
<p align="center"><kbd><img src="resources/captures/page_6.png" width="500px"/></kbd></p>
<p align="center">If you selected Deep Discovery Analyzer on the first step, provide its IP/DNS address and API Key. If you are using a self-signed certificate, check TLS Errors Ignore.</p>
<p align="center"><kbd><img src="resources/captures/page_7.png" width="500px"/></kbd></p>
<p align="center">Remove the checkbox if there is no need to run Sandboxer automatically. It will be launched automatically upon file submission.</p>
<p align="center"><kbd><img src="resources/captures/page_9.png" width="500px"/></kbd></p>
<p align="center">Wait for the file copy process to finish.</p>
<p align="center"><kbd><img src="resources/captures/page_11.png" width="500px"/></kbd></p>
<p align="center">Press "Quit" button.</p>

## Usage

### To Submit File On macOS

Right-click on the file and choose Quick Actions -> Sandboxer
<p align="center"><img src="resources/captures/quick_actions.png" width="300px"/></p>

### To Submit Files On Windows

Right-click on the file and choose Send To -> Sandboxer. Note that for the latest Windows, you will have to choose first "Show more options".
<p align="center"><img src="resources/captures/send_to.png" width="300px"/></p>

### To Submit URL
Run Sandboxer, if it is not yet running, and pick from its system tray icon menu "Submit URL" item.
<p align="center"><img src="resources/captures/submit_url.png" width="300px"/></p>

### To Get Results
Pick from the Sandboxer system tray icon menu "Submissions" item. Right-click on the menu icon "⋮" and choose "Show Report" or "Investigation Package".

## Bugs

### Notifications
On Windows, if notifications are disabled not by the Sandboxer Options window, but by using Notification Center, then it is not possible to turn them back on.

### New version update
When Sandboxer shows that a new version is available, but upon downloading it shows an error, it means that the user has to wait several minutes and try again.

### Unregistration
If unregistration is not performed from the Options dialog and/or during uninstallation Deep Discovery Analyzer connection is not available, then Sandboxer will be kept on the Submitters list and the only option will be to remove it using Analyzer Web UI. After installing Sandboxer once more, it will just register one more submitter.

### Install error
If during the files copy phase, you encounter some error, try to return to the previous install phase and try again.

### Install crash
To make sure that your installation is completed, check the sandboxer_setup_wizard.log log file generated along the installer executable. The last line should contain the following: "... INFO G0001 Close Logging ...". If not, try to run the installer once more.

### macOS dark theme
If macOS uses a dark theme, it will be hard to see some text on the UI

#### Boxer picture
Icon is taken from <a href="https://www.flaticon.com/free-icons/dog" title="dog icons">Dog icons created by Freepik - Flaticon</a>

