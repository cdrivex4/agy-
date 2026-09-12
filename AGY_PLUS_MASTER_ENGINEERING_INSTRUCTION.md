# AGY++ — MASTER ENGINEERING BOOTSTRAP INSTRUCTION

## Document Purpose

This document is the authoritative bootstrap instruction for an AI coding agent tasked with creating **AGY++**, an independent application derived from the user-facing strengths of Antigravity CLI and the earlier Gemini CLI experience.

The goal is **not** to replace Google's Antigravity infrastructure.

The goal is:

> **AGY++ = Antigravity CLI on steroids, while remaining compatible with Antigravity's normal installation, authentication, configuration, MCP, skills, agents, models, projects, sessions, and future updates wherever technically and legally appropriate.**

AGY++ must be installable on the developer's Windows machine, maintainable as a normal GitHub repository, updateable through GitHub releases, and architected so that future changes to Google's Antigravity CLI can be incorporated without rewriting AGY++.

This is an engineering specification, not a suggestion list. The coding agent must follow it as the initial project contract.

---

# 1. AUTHORITATIVE UPSTREAM PROJECTS

The implementation must continuously track these upstream projects.

## Antigravity CLI

Official repository:

https://github.com/google-antigravity/antigravity-cli

Official releases:

https://github.com/google-antigravity/antigravity-cli/releases

Official changelog:

https://github.com/google-antigravity/antigravity-cli/blob/main/CHANGELOG.md

Official documentation:

https://antigravity.google/docs/cli/overview

Official installation:

https://antigravity.google/cli/install.ps1

The current Antigravity CLI is the primary compatibility target. Its README describes it as the terminal interface to the same core agent engine used by Antigravity 2.0, with shared settings/permissions and session export. Its current authentication flow uses the system keyring and Google Sign-In when needed.

The current repository also demonstrates active releases and ongoing feature development. AGY++ must therefore assume that upstream behaviour will change.

## Gemini CLI

The former Gemini CLI remains the principal reference implementation for features that AGY++ should restore or improve.

Repository:

https://github.com/google-gemini/gemini-cli

CLI reference:

https://github.com/google-gemini/gemini-cli/blob/main/docs/cli/cli-reference.md

Architecture/project instructions:

https://github.com/google-gemini/gemini-cli/blob/main/GEMINI.md

Installation/reference:

https://github.com/google-gemini/gemini-cli/blob/main/docs/get-started/installation.mdx

Roadmap:

https://github.com/google-gemini/gemini-cli/blob/main/ROADMAP.md

Gemini CLI is open source and provides the historical reference for interactive workflows, sessions, checkpointing, authentication, approval modes, tools, MCP, and terminal UX.

IMPORTANT:

Do not blindly copy Gemini CLI code.

Use its architecture and behaviour as a reference. If code is ever copied or adapted, verify its current license, preserve required copyright/license notices, record provenance, and isolate the imported component.

---

# 2. CORE PRODUCT DEFINITION

AGY++ is a new application.

It is not:

- a renamed AGY executable;
- a private fork that blindly tracks Google's source;
- a replacement Gemini API client;
- a project that bypasses Google's authentication or service controls;
- an attempt to disguise traffic or defeat Google's abuse detection;
- a requirement to install Gemini CLI;
- a permanent patch against Google's binaries.

AGY++ is:

- an independent CLI/TUI;
- a compatibility layer around normal AGY functionality where useful;
- a feature layer restoring useful Gemini CLI workflows;
- a persistent user-owned session/workflow layer;
- an extensible architecture capable of adapting when AGY changes;
- a GitHub-maintained application with its own releases.

The first-generation product should continue using Google's AGY plumbing wherever AGY already provides the required capability.

---

# 3. NON-NEGOTIABLE ARCHITECTURAL PRINCIPLE

The application must separate **AGY compatibility** from **AGY++ functionality**.

Use this conceptual structure:

    AGY++
       |
       +-----------------------------+
       |                             |
       v                             v
  AGY++ UX/Features             AGY Compatibility Layer
       |                             |
       |                       auth/config/MCP/
       |                       models/sessions/
       |                       agent capabilities
       |                             |
       +-------------+---------------+
                     |
                AGY backend/
                official interfaces
                or stable supported
                integration points

Never scatter AGY-specific assumptions throughout the application.

All AGY-specific behaviour belongs behind an adapter/interface.

---

# 4. UPSTREAM COMPATIBILITY MUST BE A FIRST-CLASS FEATURE

Create an explicit compatibility subsystem.

Recommended conceptual package:

    internal/upstream/
        agy/
            adapter
            discovery
            capabilities
            version
            config
            auth
            sessions
            mcp
            diagnostics
        gemini/
            reference
            feature_matrix

The exact language/package structure may differ after repository inspection, but the architectural separation is mandatory.

Every AGY-specific integration must be identifiable.

For example:

    AgyConfigAdapter
    AgyAuthAdapter
    AgySessionAdapter
    AgyMCPAdapter
    AgyProcessAdapter
    AgyCapabilityDetector

Do not put code such as:

    if agyVersion >= ...

throughout unrelated UI code.

Instead:

    capabilities.Supports("feature")

or:

    backend.Capabilities()

must be used.

---

# 5. UPSTREAM VERSION MANIFEST

Create:

    docs/upstream/UPSTREAM.md

and:

    config/upstream.yaml

The manifest must record:

    agy_repository:
      url: https://github.com/google-antigravity/antigravity-cli
      branch: main
      observed_version:
      observed_commit:
      observed_date:

    gemini_repository:
      url: https://github.com/google-gemini/gemini-cli
      branch: main
      observed_commit:
      observed_date:

The project must never silently claim compatibility with an AGY version that has not been tested.

Maintain a compatibility table:

    AGY version | Tested | Compatible | Notes | Required adapter changes

Example:

    1.1.16 | yes | yes | baseline
    1.1.17 | no  | unknown | pending CI/manual test

---

# 6. UPSTREAM UPDATE TRACKING

The repository must contain tooling that allows an AI developer or human developer to detect upstream changes.

Create:

    scripts/upstream-check.ps1
    scripts/upstream-check.sh
    scripts/upstream-diff.ps1
    scripts/upstream-diff.sh

These scripts should:

1. query the official AGY repository;
2. determine the latest commit/release;
3. compare it with the recorded compatibility baseline;
4. identify changes in relevant files/features;
5. generate a report;
6. never automatically overwrite AGY++ source.

Generate:

    docs/upstream/reports/YYYY-MM-DD-agy-update.md

The report must contain:

    Upstream version
    Previous known version
    Commit range
    Relevant changed files
    Authentication changes
    Permission changes
    Session changes
    MCP changes
    CLI/TUI changes
    Configuration changes
    Model changes
    Headless changes
    Breaking-change assessment
    AGY++ action required

---

# 7. UPSTREAM CHANGE POLICY

When Google changes AGY:

DO NOT immediately merge or copy changes.

First:

    fetch
    inspect
    classify
    test
    adapt
    validate
    document

Classify every change as:

    NONE
    INFORMATIONAL
    COMPATIBILITY_RELEVANT
    FEATURE_OPPORTUNITY
    BREAKING
    SECURITY_RELEVANT

Only then modify AGY++.

AGY++ must remain buildable even if upstream changes.

---

# 8. DO NOT VENDOR THE ENTIRE AGY PROJECT

Do not copy the entire Antigravity repository into AGY++.

Do not create:

    vendor/google-antigravity/antigravity-cli/

unless a future engineering decision explicitly requires a legally permitted isolated dependency.

Prefer:

    compatibility adapter
    external AGY installation
    documented interfaces
    protocol adapters
    configuration readers
    supported invocation

This keeps AGY++ small and prevents an impossible merge problem.

---

# 9. DO NOT MODIFY USER'S AGY INSTALLATION

AGY++ must not silently patch, replace, binary-edit, or inject code into the user's normal AGY installation.

It may:

- discover AGY;
- read supported configuration;
- invoke supported commands;
- use documented configuration mechanisms;
- use supported authentication;
- use compatible session/configuration interfaces where appropriate.

Any operation that changes AGY configuration must be:

- explicit;
- reversible;
- documented;
- backed up when appropriate.

Never silently corrupt or replace AGY configuration.

---

# 10. CONFIGURATION COMPATIBILITY

AGY++ must discover existing AGY configuration.

Build a configuration abstraction:

    EffectiveConfig
       |
       +-- AGY config
       +-- AGY++ config
       +-- environment
       +-- project config
       +-- command-line overrides

Precedence must be documented.

Recommended:

    command line
        >
    project AGY++ configuration
        >
    user AGY++ configuration
        >
    AGY configuration
        >
    environment/defaults

However, authentication and security-sensitive settings must follow the actual AGY security model rather than blindly following this ordering.

AGY++ must not duplicate settings that AGY already owns unless there is a demonstrated reason.

---

# 11. AUTHENTICATION COMPATIBILITY

Authentication is a critical subsystem.

AGY++ must first determine how the installed AGY version authenticates.

Do not scrape credentials.

Do not print credentials.

Do not copy tokens into logs.

Do not create plaintext credential files.

Do not intercept OAuth callbacks unless required and explicitly supported.

Prefer AGY's normal authentication and system keyring mechanisms.

Provide an AGY++ authentication facade:

    /login
    /logout
    /auth
    /account
    /auth status

The implementation should delegate to AGY's supported login mechanism whenever possible.

For example:

    AgyAuthAdapter.Login()
    AgyAuthAdapter.Logout()
    AgyAuthAdapter.Status()

The user experience should be:

    /login

    Opening Google authentication...

    ✓ Authentication completed
    ✓ AGY credentials detected
    ✓ Active account: <redacted identifier>

Never display access tokens, refresh tokens, cookies, API keys, or secrets.

---

# 12. AUTHENTICATION FALLBACK

If AGY introduces another supported authentication method, AGY++ should detect it through the adapter capability system.

Possible states:

    Google OAuth
    system keyring session
    enterprise authentication
    API key
    application/default credentials
    unavailable

Do not assume all methods exist in every version.

The UI should dynamically display only supported options.

---

# 13. MCP COMPATIBILITY

MCP is part of the AGY ecosystem and must be treated as an existing capability, not reinvented unnecessarily.

AGY++ should:

- discover AGY MCP configuration;
- display MCP status;
- avoid duplicating servers;
- provide diagnostics;
- preserve AGY-compatible configuration;
- expose MCP status in the session UI.

Commands may include:

    /mcp
    /mcp status
    /mcp reload

Only implement configuration management independently if AGY's existing mechanism is insufficient.

Never silently modify MCP servers.

---

# 14. SKILLS, AGENTS, PLUGINS AND RULES

AGY++ should discover and preserve the normal AGY ecosystem.

Create capability adapters for:

    skills
    agents
    plugins
    rules
    MCP
    project configuration

The initial implementation should reuse AGY's existing mechanisms.

Do not create a competing skill format merely because it is easier.

If AGY changes its format, isolate the compatibility change inside the adapter.

---

# 15. SESSION COMPATIBILITY

AGY has evolved persistent conversation storage and currently uses SQLite for conversation data.

AGY++ must first understand the current format before touching it.

Create:

    docs/upstream/AGY_SESSION_FORMAT.md

Document:

    location
    schema
    version
    identifiers
    workspace/project relation
    message structure
    tool trajectory
    artifacts
    checkpoints
    migration risks

Do not assume that AGY's internal SQLite database is a stable public API.

Treat it as:

    READ_ONLY COMPATIBILITY DATA

unless official documentation explicitly permits writes.

AGY++ should preferably maintain its own metadata layer rather than corrupting AGY's database.

---

# 16. AGY++ SESSION METADATA

Create a separate Agy++ data store.

Recommended:

    ~/.agy-plus-plus/

with:

    config/
    sessions/
    cache/
    logs/
    diagnostics/
    upstream/
    state/

Use SQLite for Agy++ metadata.

Example:

    sessions
    bookmarks
    checkpoints
    feature_flags
    upstream_versions
    compatibility_results

The Agy++ session must reference an AGY session rather than pretending to own it.

Conceptually:

    AgyPlusSession
        id
        backend = agy
        agy_session_id
        workspace
        project
        title
        created
        updated
        compatibility_version

---

# 17. PRIMARY FEATURE: YOLO MODE

YOLO is a headline feature.

AGY++ must provide:

    --yolo
    -y

and:

    --approval-mode=yolo

where technically compatible.

Interactive:

    /yolo

Optional toggle:

    Ctrl+Y

if the TUI can safely intercept it.

The UI must clearly display:

    MODE: YOLO

YOLO must never be activated because a prompt, project file, skill, MCP server, or agent asks for it.

Only the user can activate it.

---

# 18. YOLO IMPLEMENTATION RULE

Do not assume that terminal keystroke injection is the correct implementation.

First determine the actual AGY permission/control interfaces.

The implementation priority is:

    1. official AGY permission API/interface
    2. supported CLI flag/configuration
    3. stable structured event interface
    4. controlled subprocess integration
    5. PTY interaction only as a last-resort compatibility adapter

If PTY interaction becomes necessary, isolate it completely:

    internal/upstream/agy/pty/

Do not allow PTY parsing to leak into the rest of Agy++.

The adapter must have automated tests against captured fixtures.

---

# 19. YOLO + SANDBOX

If AGY supports sandboxing, YOLO must not automatically disable it.

The desired semantics are:

    YOLO = automatic approval

not:

    YOLO = disable every safety boundary

If AGY provides:

    sandbox + auto approval

Agy++ should preserve both.

If the AGY version makes those modes incompatible, Agy++ must report that explicitly rather than pretending sandboxing remains active.

Example:

    WARNING: AGY 1.x does not support simultaneous YOLO and sandbox mode.
    Continuing would remove sandbox restrictions.
    Operation cancelled.

Never silently weaken a security boundary.

---

# 20. PLAN MODE

Restore/enhance the Gemini-style plan workflow.

Provide:

    --plan
    --approval-mode=plan
    /plan

Plan mode should be visibly distinct.

Where AGY itself provides plan functionality, use it.

Where Agy++ can enforce read-only behaviour independently, it may add an outer policy layer.

Do not claim a tool is read-only unless the implementation actually enforces it.

---

# 21. COMPACT

Implement:

    /compact

as an Agy++ workflow.

The first version should use AGY's supported context/session facilities where available.

If AGY lacks a direct equivalent, Agy++ may maintain its own compacted session metadata.

Never destroy the original transcript merely because context was compacted.

Keep:

    original transcript
    compact summary
    active context

separate.

---

# 22. REWIND / CHECKPOINTS

Implement:

    /rewind
    /checkpoint
    /diff

in phases.

First version:

    conversation checkpoint
    Git diff checkpoint

Do not attempt arbitrary filesystem restoration initially.

If the workspace is a Git repository:

    checkpoint = Git state + Agy++ session state

Later support may include:

    Git worktrees
    Jujutsu
    filesystem snapshots
    Windows VSS

but these are future capabilities.

---

# 23. INTERACTIVE SESSION

Agy++ must support:

    agy++

    agy++ "initial prompt"

    agy++ -i "initial prompt"

    agy++ --resume latest

    agy++ --resume <id>

The interactive UI should retain the strengths of Gemini CLI:

- persistent prompt
- streaming output
- tool activity
- cancellation
- history
- slash commands
- model selection
- session selection
- clear permission state
- visible backend state

---

# 24. COMMAND DESIGN

Initial command set:

    /help
    /login
    /logout
    /auth
    /account
    /yolo
    /plan
    /default
    /auto-edit
    /sessions
    /resume
    /compact
    /rewind
    /checkpoint
    /diff
    /model
    /status
    /usage
    /permissions
    /mcp
    /skills
    /agents
    /config
    /diagnostics
    /upstream
    /quit

Do not implement all commands at once.

Create the command registry so unsupported commands can report:

    Command unavailable with current AGY backend/version.

---

# 25. STATUS PANEL

The TUI should make the underlying system visible.

Display:

    AGY++
    AGY version
    compatibility status
    authenticated account
    model
    project
    workspace
    execution mode
    sandbox status
    MCP status
    session ID

Never display secrets.

---

# 26. AGY++ ARCHITECTURE

Use an architecture approximately equivalent to:

    cmd/
        agy-plus-plus/

    internal/
        app/
        cli/
        tui/
        commands/
        config/
        auth/
        session/
        policy/
        upstream/
            agy/
            gemini-reference/
        backend/
        process/
        mcp/
        skills/
        agents/
        checkpoint/
        compact/
        diagnostics/
        logging/

    tests/
        unit/
        integration/
        compatibility/
        fixtures/
        security/

    docs/
        architecture/
        upstream/
        features/
        security/
        development/

    scripts/
        upstream/
        build/
        release/
        install/

The exact structure may change after the coding agent evaluates the chosen language and UI framework.

The separation of responsibilities must not change.

---

# 27. LANGUAGE

The default recommendation is Go.

Reasons:

- excellent CLI support;
- strong Windows support;
- static binaries;
- process supervision;
- concurrency;
- simple distribution;
- low runtime overhead;
- suitable TUI libraries;
- AGY ecosystem alignment.

Do not select Go merely because AGY is written in Go.

Before committing, document the decision in:

    docs/architecture/ADR-0001-language.md

If Go is selected, use a modern Go version and pin dependencies.

---

# 28. WINDOWS IS A FIRST-CLASS TARGET

The primary developer installation target is Windows.

The application must support:

    Windows 10/11 where dependencies permit

The build must produce:

    agy++.exe

Installation should ultimately support:

    PowerShell installer
    ZIP release
    GitHub release asset

Optional future:

    winget
    Scoop
    Chocolatey

Do not require WSL merely to install Agy++.

---

# 29. INSTALLATION MODEL

The user should eventually be able to install:

    agy++

and execute:

    agy++

without knowing the source tree.

Preferred installation locations:

    %LOCALAPPDATA%\AgyPlusPlus\

or:

    %LOCALAPPDATA%\Programs\AgyPlusPlus\

User data must NOT live inside the installation directory.

Use:

    %APPDATA%\AgyPlusPlus\
    %LOCALAPPDATA%\AgyPlusPlus\

according to Windows conventions.

---

# 30. GITHUB UPDATE MODEL

Agy++ must be a normal GitHub project.

Repository:

    https://github.com/<OWNER>/agy-plus-plus

Do not hard-code the user's eventual GitHub username.

Provide:

    agy++ update

and:

    agy++ update --check

The update mechanism should:

1. check the GitHub release API;
2. compare current version;
3. display release notes;
4. ask before installing;
5. download the correct signed/hashed artifact;
6. verify integrity;
7. replace the executable safely;
8. preserve configuration and sessions;
9. provide rollback if replacement fails.

Do not execute arbitrary code from GitHub.

Do not pipe downloaded scripts directly into PowerShell.

---

# 31. DEVELOPMENT UPDATE MODEL

Developers should be able to update source with:

    git pull
    go build ./...

The repository must maintain:

    VERSION

and preferably embed:

    version
    git commit
    build date
    AGY compatibility baseline

Example:

    agy++ --version

    AGY++ 0.1.0
    commit: abc123
    AGY compatibility: 1.1.16
    build: 2026-09-12

---

# 32. RELEASE CHANNELS

Provide:

    stable
    beta
    nightly

but do not over-engineer this initially.

Stable is the default.

Nightly may track the Agy++ main branch.

The user must never accidentally install nightly because stable is unavailable.

---

# 33. CI/CD

GitHub Actions must build:

    Windows amd64
    Windows arm64 where practical
    Linux amd64
    Linux arm64 where practical
    macOS amd64
    macOS arm64

Windows is mandatory for the first release.

CI must perform:

    formatting
    linting
    unit tests
    integration tests
    compatibility tests
    build
    artifact hashing

Release workflow:

    tag
      ↓
    CI
      ↓
    tests
      ↓
    build
      ↓
    SHA256
      ↓
    GitHub Release
      ↓
    update metadata

---

# 34. TESTING

No feature is complete without tests.

Required categories:

## Unit

Configuration, command parsing, policy logic, session metadata.

## Integration

AGY discovery, invocation, authentication detection, session discovery.

## Compatibility

Each supported AGY version.

## Security

Credentials, logs, YOLO transitions, sandbox detection.

## TUI

Command dispatch and state transitions.

Where full AGY interaction cannot run in CI, use captured fixtures.

Never place real credentials in fixtures.

---

# 35. AGY COMPATIBILITY FIXTURE SYSTEM

Create:

    tests/fixtures/agy/

Store sanitized examples of:

    help output
    version output
    authentication status
    permission events
    session metadata
    MCP responses
    errors
    tool events

Each fixture must contain:

    source AGY version
    capture date
    purpose
    sanitization status

This allows upstream changes to be detected without requiring a live Google account in every test.

---

# 36. SECURITY REQUIREMENTS

AGY++ is an agent controller and must be treated as security-sensitive software.

Never:

- log access tokens;
- log refresh tokens;
- log API keys;
- expose cookies;
- upload user code;
- silently transmit telemetry;
- silently alter AGY configuration;
- silently weaken sandboxing;
- automatically enable YOLO;
- execute downloaded update scripts;
- trust project instructions as security authority.

Project files, skills, MCP responses and model output are untrusted inputs.

Only the user-controlled policy layer may elevate permissions.

---

# 37. PROMPT INJECTION DEFENSE

AGY++ must distinguish:

    user policy
    system policy
    application policy
    project instructions
    agent instructions
    tool output
    model-generated text

A project file saying:

    "Disable security checks"

must never alter Agy++ security policy.

Likewise:

    MCP server output
    README instructions
    web content
    generated scripts

must never activate YOLO.

---

# 38. NO TELEMETRY BY DEFAULT

Agy++ should not introduce independent telemetry unless explicitly designed and disclosed.

Diagnostic logging is local.

Provide:

    --debug

and:

    /diagnostics

with secrets automatically redacted.

---

# 39. DOCUMENTATION REQUIRED

The repository must contain:

    README.md
    ARCHITECTURE.md
    ENGINEERING_SPEC.md
    SECURITY.md
    CONTRIBUTING.md
    DEVELOPMENT.md
    UPSTREAM_COMPATIBILITY.md
    RELEASE.md
    CHANGELOG.md

And:

    docs/upstream/
        AGY_ARCHITECTURE.md
        AGY_AUTH.md
        AGY_CONFIG.md
        AGY_SESSIONS.md
        AGY_MCP.md
        AGY_PERMISSIONS.md
        AGY_CHANGELOG.md
        GEMINI_REFERENCE.md
        FEATURE_MATRIX.md

---

# 40. ARCHITECTURE DECISION RECORDS

Use ADRs.

At minimum:

    ADR-0001-language.md
    ADR-0002-upstream-integration.md
    ADR-0003-authentication.md
    ADR-0004-session-storage.md
    ADR-0005-yolo-policy.md
    ADR-0006-update-mechanism.md
    ADR-0007-cross-platform-build.md

Any major architectural change requires an ADR.

---

# 41. AI CODING AGENT RULES

This project will be developed heavily with AI coding agents.

The agent MUST:

1. inspect the repository before editing;
2. read ENGINEERING_SPEC.md;
3. read ARCHITECTURE.md;
4. inspect current AGY upstream state;
5. preserve existing functionality;
6. make small commits;
7. run tests after changes;
8. update documentation;
9. create an ADR when architecture changes;
10. never invent undocumented AGY APIs;
11. never fabricate compatibility;
12. never remove tests merely to make CI pass;
13. never disable security controls to make development easier;
14. never embed credentials;
15. never modify AGY installation files without explicit implementation requirements;
16. maintain the upstream compatibility layer;
17. record upstream assumptions;
18. flag uncertain behaviour rather than guessing.

---

# 42. AI DEVELOPMENT LOOP

Every implementation task must follow:

    READ
      ↓
    UNDERSTAND
      ↓
    INSPECT UPSTREAM
      ↓
    PLAN
      ↓
    IMPLEMENT
      ↓
    TEST
      ↓
    REVIEW
      ↓
    DOCUMENT
      ↓
    COMMIT

For compatibility work:

    detect upstream change
      ↓
    reproduce
      ↓
    identify affected adapter
      ↓
    add regression fixture
      ↓
    modify adapter
      ↓
    test
      ↓
    document

---

# 43. COMMIT STRATEGY

Use meaningful commits.

Examples:

    feat: add AGY installation discovery
    feat: add interactive session shell
    feat: add yolo policy state
    feat: add AGY authentication adapter
    feat: add session metadata store
    fix: adapt AGY 1.1.17 permission event
    test: add AGY 1.1.17 compatibility fixtures
    docs: document AGY authentication compatibility

Do not produce enormous one-commit implementations.

---

# 44. INITIAL IMPLEMENTATION PHASES

## PHASE 0 — Repository Intelligence

Before coding functionality:

- inspect AGY repository;
- inspect latest release;
- inspect README;
- inspect changelog;
- inspect CLI help;
- inspect installation;
- inspect authentication;
- inspect configuration;
- inspect MCP;
- inspect sessions;
- inspect permissions;
- inspect headless mode;
- inspect relevant source/interfaces;
- inspect Gemini CLI reference architecture.

Produce:

    docs/upstream/AGY_ARCHITECTURE.md
    docs/upstream/FEATURE_MATRIX.md
    docs/upstream/AGY_AUTH.md
    docs/upstream/AGY_CONFIG.md
    docs/upstream/AGY_SESSIONS.md
    docs/upstream/AGY_PERMISSIONS.md
    docs/upstream/GEMINI_REFERENCE.md

Do not start large implementation until these exist.

---

# 45. PHASE 1 — Minimal Agy++ Shell

Implement:

    agy++ --version
    agy++ --help
    agy++

Capabilities:

- discover AGY;
- report AGY version;
- report compatibility;
- start normal AGY-compatible interaction;
- cleanly exit;
- preserve normal AGY behaviour.

Success condition:

> A user can install Agy++ and use it as their normal AGY CLI without losing the normal AGY experience.

---

# 46. PHASE 2 — Interactive Command Framework

Implement the slash-command framework.

Start with:

    /help
    /status
    /login
    /logout
    /model
    /sessions
    /quit

Do not implement every command yet.

The framework must allow commands to declare:

    name
    aliases
    description
    arguments
    required capabilities
    execution mode
    backend requirements

---

# 47. PHASE 3 — YOLO

Implement:

    -y
    --yolo
    --approval-mode=yolo
    /yolo

Then:

    Ctrl+Y

if technically safe.

Add tests.

The user must always see:

    YOLO ENABLED

and:

    YOLO DISABLED

when the state changes.

---

# 48. PHASE 4 — LOGIN

Implement:

    /login
    /logout
    /auth status
    /account

Use AGY's supported authentication mechanism.

Test:

    first login
    existing login
    expired session
    logout
    re-login
    SSH/remote case where supported

Never expose credentials.

---

# 49. PHASE 5 — SESSION POWER FEATURES

Implement:

    /sessions
    /resume
    /checkpoint
    /rewind
    /compact

Do not destroy AGY data.

Use Agy++ metadata to augment AGY sessions.

---

# 50. PHASE 6 — POWER-USER TUI

Implement:

- polished terminal layout;
- keyboard shortcuts;
- tool status;
- session information;
- permission state;
- model state;
- backend state;
- compact/rewind UI;
- command palette;
- error recovery.

The TUI must remain usable over SSH.

---

# 51. PHASE 7 — HEADLESS

Support:

    agy++ -p "prompt"

and:

    --output-format text
    --output-format json
    --output-format stream-json

where compatible with the backend.

Provide:

    timeout
    exit codes
    cancellation
    machine-readable errors

This allows Agy++ to be used from:

    PowerShell
    CI
    automation
    other AI agents
    n8n
    scripts

---

# 52. PHASE 8 — ADVANCED AGY INTEGRATION

Only after the basic application is stable:

- richer MCP integration;
- agent selection;
- skill management;
- model selection;
- subagent support;
- background tasks;
- artifact support;
- remote workflows.

Reuse AGY capabilities whenever possible.

---

# 53. PHASE 9 — UPSTREAM AUTOMATION

Add:

    agy++ upstream check
    agy++ upstream status
    agy++ diagnostics

The command should report:

    Installed AGY
    Tested AGY
    Latest AGY
    Compatibility status
    Adapter warnings
    Recommended action

Example:

    AGY++ 0.4.0
    AGY installed: 1.1.17
    Compatibility baseline: 1.1.16

    STATUS: WARNING

    A newer AGY version is installed than the tested compatibility baseline.

    Changed areas:
      permissions
      authentication

    Run:
      agy++ upstream check

---

# 54. PHASE 10 — RELEASE AND SELF-UPDATE

Create the GitHub release pipeline.

Agy++ should be downloadable from:

    GitHub Releases

The executable must be reproducible from source.

Publish:

    Windows ZIP
    Windows installer
    checksums

Later:

    macOS
    Linux

The updater must preserve:

    configuration
    sessions
    AGY integration settings

---

# 55. INSTALLER REQUIREMENTS

Create a Windows installer.

It must:

- install Agy++;
- add executable to PATH;
- detect AGY;
- not replace AGY;
- not overwrite AGY configuration;
- create Start Menu entry only if useful;
- support uninstall;
- preserve user data during uninstall unless the user explicitly selects deletion.

Provide:

    install.ps1

but do not require users to execute an untrusted remote PowerShell script blindly.

The official release page should provide the installer.

---

# 56. VERSION COMPATIBILITY

Agy++ has its own version:

    0.x.y

AGY has its own version:

    1.x.y

Never pretend these are the same product/version.

Display both.

Example:

    AGY++ 0.2.0
    AGY 1.1.16
    Compatibility: VERIFIED

---

# 57. FAILURE BEHAVIOUR

If AGY changes incompatibly:

Agy++ must fail gracefully.

Example:

    AGY++ detected an unsupported AGY permission protocol.

    AGY version: 1.2.0
    Tested through: 1.1.16

    Normal AGY functionality may still work.

    Agy++ YOLO integration has been disabled to avoid unsafe behaviour.

    Run:
        agy++ upstream check

Never silently guess a new protocol for security-sensitive features.

---

# 58. BACKWARD COMPATIBILITY

Agy++ must not require the latest AGY if an older supported version is already installed.

Maintain a support window.

Example:

    Supported:
      AGY 1.1.x

    Warning:
      newer untested versions

    Unsupported:
      versions below minimum

The exact versions are determined during Phase 0.

---

# 59. DATA MIGRATION

Never perform destructive migration automatically.

For any AGY++ database schema change:

    detect
    backup
    migrate
    verify
    rollback on failure

Migration scripts must be tested.

---

# 60. OBSERVABILITY

Provide local diagnostic information:

    agy++ diagnostics

Output:

    Agy++ version
    OS
    architecture
    AGY version
    AGY executable path
    configuration paths
    authentication state without secrets
    MCP status
    compatibility baseline
    session database status
    feature capability matrix

Redact:

    tokens
    API keys
    cookies
    passwords
    secrets
    sensitive prompt contents

---

# 61. FEATURE FLAGS

Use feature flags for unstable integrations.

Example:

    experimental.yolo
    experimental.rewind
    experimental.session-bridge

Do not hide unstable features behind undocumented environment variables.

Expose them clearly.

---

# 62. LEGAL / LICENSE PROVENANCE

Maintain:

    THIRD_PARTY_NOTICES.md

If any source code is copied from Gemini CLI or another project:

- identify the source;
- identify license;
- preserve required notices;
- isolate it;
- document modifications.

Do not copy proprietary AGY implementation code merely because it can be obtained from a local installation.

Prefer clean-room adapter implementation where necessary.

The project must remain honest about what is:

    original Agy++ code
    adapted open-source code
    external dependency
    AGY runtime capability

---

# 63. NO SECRET BACKDOORS

The application must never contain:

- hidden remote commands;
- hidden telemetry;
- credential exfiltration;
- undocumented network endpoints;
- hard-coded personal accounts;
- developer-only bypasses;
- hidden YOLO activation;
- remote update execution without verification.

---

# 64. DEFINITION OF DONE FOR INITIAL RELEASE

Agy++ v0.1.0 is complete only when all of the following work:

    [ ] installs on Windows
    [ ] detects AGY
    [ ] displays AGY version
    [ ] starts normal AGY workflow
    [ ] interactive session works
    [ ] /help works
    [ ] /status works
    [ ] /login works or delegates correctly
    [ ] /logout works or delegates correctly
    [ ] YOLO works
    [ ] YOLO is visibly active
    [ ] YOLO cannot be activated by model/project instructions
    [ ] normal AGY permissions remain intact outside YOLO
    [ ] MCP remains functional
    [ ] sessions remain usable
    [ ] no credentials are exposed
    [ ] tests pass
    [ ] CI builds Windows artifact
    [ ] GitHub release can install it
    [ ] uninstall works
    [ ] README documents installation
    [ ] compatibility documentation exists
    [ ] upstream version is recorded
    [ ] source can rebuild the release

---

# 65. FIRST DEVELOPMENT COMMAND

After creating the repository, the AI developer must NOT immediately generate the entire application.

First perform:

    1. repository initialization
    2. upstream repository inspection
    3. AGY installation detection
    4. current AGY version detection
    5. Gemini architecture inspection
    6. compatibility matrix creation
    7. architecture decision records
    8. minimal scaffold
    9. build/test pipeline
    10. only then Phase 1 implementation

The initial commit should therefore be a **scaffold and research baseline**, not a huge implementation.

---

# 66. REQUIRED INITIAL REPOSITORY TREE

The initial scaffold should approximately contain:

    agy-plus-plus/
    |
    +-- .github/
    |   +-- workflows/
    |       +-- ci.yml
    |       +-- release.yml
    |
    +-- cmd/
    |   +-- agy-plus-plus/
    |
    +-- internal/
    |   +-- app/
    |   +-- cli/
    |   +-- tui/
    |   +-- commands/
    |   +-- config/
    |   +-- auth/
    |   +-- session/
    |   +-- policy/
    |   +-- backend/
    |   +-- upstream/
    |       +-- agy/
    |       +-- gemini-reference/
    |   +-- diagnostics/
    |   +-- logging/
    |
    +-- tests/
    |   +-- unit/
    |   +-- integration/
    |   +-- compatibility/
    |   +-- fixtures/
    |
    +-- docs/
    |   +-- architecture/
    |   +-- upstream/
    |   +-- security/
    |   +-- development/
    |
    +-- scripts/
    |   +-- upstream/
    |   +-- build/
    |   +-- release/
    |
    +-- ARCHITECTURE.md
    +-- ENGINEERING_SPEC.md
    +-- SECURITY.md
    +-- DEVELOPMENT.md
    +-- UPSTREAM_COMPATIBILITY.md
    +-- CONTRIBUTING.md
    +-- CHANGELOG.md
    +-- THIRD_PARTY_NOTICES.md
    +-- LICENSE
    +-- README.md
    +-- VERSION

---

# 67. FINAL DIRECTIVE TO THE AI DEVELOPER

You are building Agy++ as a long-lived independent project.

Do not rush to implement every feature.

First establish a stable compatibility boundary with Antigravity CLI.

The fundamental rule is:

    USE AGY WHERE AGY ALREADY WORKS.
    ADD Agy++ WHERE AGY IS MISSING CAPABILITY.
    ISOLATE EVERY AGY DEPENDENCY.
    TRACK EVERY UPSTREAM CHANGE.
    NEVER BREAK NORMAL AGY OPERATION.
    NEVER SACRIFICE SECURITY FOR CONVENIENCE.
    NEVER GUESS AN UNDOCUMENTED PROTOCOL.
    NEVER MAKE Agy++ DEPEND ON A SINGLE VERSION OF AGY.

The desired end state is:

    Google Antigravity
          |
          | official/stable integration
          v
        AGY
          |
          | compatibility boundary
          v
       AGY++
          |
          +-- Gemini-style YOLO
          +-- interactive workflow
          +-- improved login UX
          +-- sessions
          +-- rewind
          +-- compact
          +-- plan
          +-- enhanced automation
          +-- compatibility management
          +-- future features

Agy++ must remain useful even when AGY changes.

When AGY improves, Agy++ should gain the improvement where appropriate.

When AGY removes or changes a capability, Agy++ should preserve the user's workflow whenever technically possible.

When an upstream change breaks an adapter, the compatibility system must identify it clearly and prevent unsafe silent behaviour.

Build Agy++ as a product, not a temporary hack.

The first objective is a clean, installable, updateable scaffold.

The second objective is a working AGY-compatible CLI.

The third objective is restoring the Gemini-era power-user experience.

The fourth objective is continuous upstream compatibility.

Everything else comes after those foundations.
