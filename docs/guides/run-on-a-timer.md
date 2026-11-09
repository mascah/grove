# Run on a timer

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

Grove does not schedule itself. On a Mac, launchd can run the line from
[Running unattended](../running.md) every so often, in a checkout you keep
for one project:

```sh
K=$(grove claim) && timeout 2h claude -p "/grove-implement $K --unattended"
```

A cycle you mark `unattended` is then started without anyone watching.
Each firing either starts one session or starts nothing. It starts nothing
while the project is paused after an attempt, from this checkout or any
other, or when no marked cycle can be claimed. Marking a cycle is what
puts work in front of the timer; merging a change alone does not.

## Before the timer

Set up the checkout as "Set up a checkout" in
[Running unattended](../running.md) says, and run the line in it once by
hand to see a session through. The timer's session has the same
permissions, from the checkout's `.claude/settings.local.json`, and nobody
to answer it. So a tool those settings do not allow fails, and the agent
works around it or stops. Claude Code uses the login you made with
`claude` in Terminal. A job in your user's launchd domain, as below, can
read it.

## The job

Save this as `~/Library/LaunchAgents/com.you.grove-claim.ledger.plist`,
with your user's home in place of `/Users/you`, one label and file per
project:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.you.grove-claim.ledger</string>
  <key>ProgramArguments</key>
  <array>
    <string>/bin/sh</string>
    <string>-c</string>
    <string>K=$(grove claim) &amp;&amp; timeout 2h claude -p "/grove-implement $K --unattended"</string>
  </array>
  <key>WorkingDirectory</key>
  <string>/Users/you/grove-runs/ledger</string>
  <key>EnvironmentVariables</key>
  <dict>
    <key>PATH</key>
    <string>/Users/you/.local/bin:/opt/homebrew/bin:/usr/bin:/bin</string>
  </dict>
  <key>StartInterval</key>
  <integer>900</integer>
  <key>StandardOutPath</key>
  <string>/Users/you/grove-runs/ledger.log</string>
  <key>StandardErrorPath</key>
  <string>/Users/you/grove-runs/ledger.log</string>
</dict>
</plist>
```

- `&amp;&amp;` is `&&` written as XML needs it.
- `PATH` is all the session gets: launchd does not read your shell's
  profile. It needs the folders holding `grove`, `claude`, `git`, `gh`,
  `timeout` and your project's test tools. `which grove claude gh timeout`
  in Terminal shows them; `timeout` comes with Homebrew's `coreutils`.
- `StartInterval` is in seconds. launchd starts no second copy of a job
  that is still running, so a firing during a session is skipped. A
  `grove claim` you run by hand meanwhile finds the project paused.
- `timeout 2h` stops the session after two hours. Change it in the line,
  such as `timeout 45m`.

Load it, and it fires every interval, the first one interval from now:

```sh
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.you.grove-claim.ledger.plist
```

## Where to look

- `/Users/you/grove-runs/ledger.log` gets one line per firing from Grove.
  The line says the project is paused, there is nothing to claim, or what
  went wrong. A firing that claims a cycle adds its ID, then whatever
  `claude -p` prints, which is its final message.
- Each attempt's record is on `grove/runs`, and
  [Watching the project](../watching.md) lists them:
  `git log -1 --format=%B origin/grove/runs` shows the latest.
- `launchctl print gui/$(id -u)/com.you.grove-claim.ledger` shows whether
  the job is running, how many times it has run, and its last exit code.
  The code is 0 for a session that exited 0. It is non-zero for a firing
  that had nothing to do, a refusal, a session that failed, or one that
  `timeout` stopped (124).

When the project is paused, deal with the attempt. Then either run
`grove claim --retry` with the agent in the checkout yourself, or merge
its pull request. The timer picks up from there.

## Stopping it

```sh
launchctl bootout gui/$(id -u)/com.you.grove-claim.ledger
```

If a session is running, launchd stops the job's shell and the session
with it. The attempt's branch and files stay in the checkout for you.
