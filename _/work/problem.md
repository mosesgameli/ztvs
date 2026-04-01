Axios npm Package Compromised: Supply Chain Attack Delivers Cross-Platform RAT
Written by
Headshot of Liran Tal
Liran Tal

March 30, 2026

13 mins read
On March 31, 2026, two malicious versions of axios, the enormously popular JavaScript HTTP client with over 100 million weekly downloads, were briefly published to npm via a compromised maintainer account. The packages contained a hidden dependency that deployed a cross-platform remote access trojan (RAT) to any machine that ran npm install (or equivalent in other package managers like Bun) during a two-hour window.

The malicious versions (1.14.1 and 0.30.4) were removed from npm by 03:29 UTC. But in the window they were live, anyone whose CI/CD pipeline, developer environment, or build system pulled a fresh install could have been compromised without ever touching a line of Axios code.

TL;DR
Snyk Advisory

SNYK-JS-AXIOS-15850650

Affected versions

axios@1.14.1, axios@0.30.4

Root cause

Hijacked npm maintainer account

Malicious dependency

plain-crypto-js@4.2.1 (SNYK-JS-PLAINCRYPTOJS-15850652)

Payload

Cross-platform RAT (macOS, Windows, Linux)

C2 server

sfrclak[.]com:8000

Published

1.14.1 at 00:21 UTC; 0.30.4 at 01:00 UTC

Removed

03:29 UTC (March 31, 2026)

Safe versions

Any version other than 1.14.1or 0.30.4

Immediate action

Audit lockfiles for affected versions; rotate secrets if exposed

How the attack was constructed
This is not a case of a typosquatted package or a rogue dependency slipping into a build. The attacker had (or gained) direct publishing access to the official axios package on npm, likely by compromising a maintainer's account. According to a collaborator in the official GitHub issue thread, the suspected compromised account belonged to maintainer @jasonsaayman, whose repository permissions were higher than those of other collaborators, complicating rapid remediation.

The attacker did not modify any Axios source files directly. Instead, they added a pre-staged malicious dependency, plain-crypto-js@4.2.1, to the package.json of the new axios releases. The plain-crypto-js package itself was purpose-built for this attack: an earlier "clean" version (4.2.0) had been published 18 hours prior, likely to give it a brief history on the registry. Version 4.2.1 contained the malicious payload.

When a developer or CI system runs npm install axios@1.14.1, npm resolves the dependency tree, pulls plain-crypto-js@4.2.1, and automatically runs its postinstall hook: node setup.js. That single script execution is where the compromise begins.

The dropper: Double-obfuscated and self-erasing
The setup.js postinstall dropper uses two layers of obfuscation to avoid static analysis:

Reversed Base64 encoding with padding character substitution

XOR cipher with the key OrDeR_7077 and a constant value of 333

Once deobfuscated, the script detects the host operating system via os.platform() and reaches out to the C2 server at sfrclak[.]com:8000 (IP: 142.11.206.73) to download a second-stage payload appropriate for the platform.

After execution, the malware erases its own tracks: it deletes setup.js, removes the package.json that contained the postinstall hook, and replaces it with a clean package.md renamed to package.json. If you inspect node_modules/plain-crypto-js after the fact, you would find no obvious signs of a postinstall script ever having been there.

Platform-specific payloads
The second-stage payloads are purpose-built for each platform.

macOS
An AppleScript downloads a binary to /Library/Caches/com.apple.act.mond, deliberately spoofing an Apple background daemon naming convention to blend in. Once established, the RAT:

Generates a 16-character unique victim ID

Fingerprints the system: hostname, username, macOS version, boot/install times, CPU architecture (mac_arm or mac_x64), running processes

Beacons to the C2 every 60 seconds using a fake IE8/Windows XP User-Agent string

Accepts four commands from the attacker:

peinject: receives a Base64-encoded binary from the C2, decodes it, writes it to a hidden temp file (e.g., /private/tmp/.XXXXXX), performs ad-hoc code signing via codesign --force --deep --sign - to bypass Gatekeeper, and executes it

runscript: runs arbitrary shell commands via /bin/sh or executes AppleScript files via osascript

rundir: enumerates filesystem metadata from /Applications, ~/Library, and ~/Application Support

kill: terminates the RAT process

Windows
A VBScript downloader copies the PowerShell binary to %PROGRAMDATA%\wt.exe (masquerading as Windows Terminal) and executes a hidden PowerShell RAT with execution policy bypass flags.

Linux
A Python RAT is downloaded to /tmp/ld.py and launched as an orphaned background process via nohup python3, detaching it from the terminal session that spawned it.


Additional compromised packages
Two other packages were observed shipping the malicious plain-crypto-js dependency:

@qqbrowser/openclaw-qbot@0.0.130 — includes a tampered axios@1.14.1 with the injected dependency (SNYK-JS-QQBROWSEROPENCLAWQBOT-15850776)

@shadanai/openclaw (versions 2026.3.31-1, 2026.3.31-2) — vendors plain-crypto-js directly (SNYK-JS-SHADANAIOPENCLAW-15850775)

These secondary packages suggest either coordinated attacker infrastructure or that the malicious plain-crypto-js was being actively used in related campaigns.

Who is actually at risk
The three-hour publication window (00:21 to 03:29 UTC) is the key constraint. Risk is highest for:

CI/CD pipelines that do not pin dependency versions and run npm install on a schedule or on commit — especially those that run overnight or in the early morning UTC.

Developers who ran npm install or npm update in that window and happened to pull the affected versions.

Projects depending on @qqbrowser/openclaw-qbot or @shadanai/openclaw, whose exposure does not depend on the window.

If your lockfile (package-lock.json or yarn.lock) was committed before the malicious versions were published and your install did not update it, you were not affected. Lockfiles are your first line of defense here.

The malicious versions have been removed from the npm registry. However, anyone who installed them during the window should assume a full system compromise: the RAT was live, beaconing, and capable of executing arbitrary follow-on payloads.

Snyk remediation and how to check your exposure
If you are a user or customer of Snyk, then any of the various Snyk integrations will alert you of any projects that vendor the compromised and malicious version of the axios dependency, whether via the Snyk CLI, the Snyk app integration, or otherwise.

Snyk's database includes entries for both SNYK-JS-AXIOS-15850650 and SNYK-JS-PLAINCRYPTOJS-15850652, so snyk test will flag the affected versions and the malicious transitive dependency.

Additionally, if you’re on the Enterprise plan, Snyk you will see a Zero Day report in the application, similar to how you’d find earlier zero day security incidents such as LiteLLM, Shai-Hulud and others, giving you a system-wide view to easily locate and pin-point affected projects and repositories that have the vulnerable axios dependency:


In the case you’re not yet using Snyk, there’s a free tier, and you can easily get started and audit your environment for potential axios compromise or other security issues as follows:

# Install Snyk CLI if you haven't already
npm install -g snyk

# Authenticate
snyk auth

# Test your project
snyk test
For Bun users: Snyk workaround (native bun.lock support is limited in the Snyk CLI at time of writing):

The recommended workaround is to generate a yarn.lock compatible lockfile using Bun's built-in -y flag, which Snyk can parse:

# 1. Regenerate lockfile in yarn.lock format
bun install -y

# 2. Run snyk against the generated yarn.lock
snyk test --file=yarn.lock

Otherwise, you can follow any of the steps below to locate and check if you’re affected by the axios compromise:

Step 1: Check your lockfile for affected versions

# Check for axios 1.14.1 or 0.30.4
grep -E '"axios"' package-lock.json | grep -E '1\.14\.1|0\.30\.4'

# Or with yarn
grep -E 'axios@' yarn.lock | grep -E '1\.14\.1|0\.30\.4'
Step 2: Check for the malicious dependency

# Look for plain-crypto-js in your dependency tree
npm ls plain-crypto-js

# Or search node_modules directly
find node_modules -name "plain-crypto-js" -type d
Step 3: Check for Bun runtime installs for the malicious axios dependency

If you are using Bun, check your bun.lock (text lockfile, Bun v1.1+):

grep -E 'axios' bun.lock | grep -E '1\.14\.1|0\.30\.4'
Also, check for the malicious transitive dependency:

grep 'plain-crypto-js' bun.lock
Note: Older Bun versions produce a binary bun.lockb. To inspect it, convert first:

> bun bun.lockb  # prints human-readable output to stdout
> bun bun.lockb | grep -E 'axios.*1\.14\.1|axios.*0\.30\.4'
> bun bun.lockb | grep 'plain-crypto-js'
Step 4: Check for IOCs on compromised systems

If you believe a machine ran npm install in the affected window, look for these indicators:

Platform

IOC

macOS

/Library/Caches/com.apple.act.mond binary

Windows

%PROGRAMDATA%\wt.exe (PowerShell masquerading as Windows Terminal)

Linux

/tmp/ld.py Python script

Network

Outbound connections to sfrclak[.]com / 142.11.206.73:8000

Further npm package manager remediation advice
If you are not affected (precautionary):

Pin axios to a known safe version in your package.json. Any version other than 1.14.1 or 0.30.4 is clean.

Commit your lockfile and ensure CI uses npm ci (not npm install) to enforce lockfile integrity.

Add plain-crypto-js to a blocklist in your package manager or security tooling.

Consider enabling --ignore-scripts for npm installs in CI environments where lifecycle hooks are not needed:

npm ci --ignore-scripts
This prevents postinstall scripts from running entirely, which would have blocked this attack vector. Be aware that it can break packages that legitimately need post-install steps (native addons, for example).

Additionally, consider using and rolling to your developers the npq open-source project that introduces security and health signal pre-checks prior to installing dependencies.

Finally, you’d likely want to review and consult these publicly curated npm security best practices.

If you are affected (assume breach):

Contain immediately: Isolate any systems that ran npm install in the affected window.

Rotate all secrets: Treat every credential on the affected machine as compromised — API keys, SSH keys, cloud credentials, npm tokens, GitHub tokens. Do not rotate in place; revoke and reissue.

Review for lateral movement: Check logs for outbound connections to sfrclak[.]com or 142.11.206.73. If the RAT was active, the attacker had arbitrary code execution and may have enumerated or exfiltrated further.

Rebuild environments: Do not attempt to clean compromised systems. Rebuild from a known-clean snapshot or base image.

Audit CI pipelines: Review build logs for the March 31, 2026 UTC window to determine which pipelines installed the affected versions.

The bigger picture: Maintainer account security
This attack follows a now-familiar pattern: compromise a legitimate maintainer account, publish a malicious version of a trusted package, and rely on the ecosystem's implicit trust of registered packages. We've seen this playbook used against ESLint's Prettier plugin, against multiple packages owned by a prolific developer via phishing, and against the Shai-Hulud campaign that compromised over 600 packages.

What makes Axios particularly significant is the scale: 100 million weekly downloads means even a two-hour malicious window represents an enormous potential blast radius. The attacker also showed meaningful operational sophistication, pre-staging the malicious dependency, using a "clean" version history, double-obfuscating the dropper, building platform-specific RATs, and implementing anti-forensic self-deletion. This was not opportunistic

For organizations that depend on open source at scale, the lesson is not to stop using npm or to distrust all dependencies. It's to understand which supply chain controls would have caught this: lockfile enforcement, postinstall script auditing, and runtime monitoring for unexpected process spawns or outbound network connections from build environments. Snyk's guide to preventing npm supply chain attacks and lockfile security considerations are worth revisiting in the context of this incident.

If you want to understand the class of attack at a conceptual level, Snyk Learn has a lesson specifically on compromise of legitimate packages that walks through the attack patterns and defensive controls.