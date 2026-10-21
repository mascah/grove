# Run apply on a timer

A guide is a recipe that uses Grove's existing features and has been tried;
it is not part of what Grove does. Use it if it helps.

## What it is for

Grove does not schedule itself. On a Mac, launchd can run `grove run` every
so often in a checkout you keep for one project, so approved work is started
without anyone watching. Each firing either starts one attempt or does
nothing: while a run is going, while the project is paused after an
attempt, or when nothing is approved. [Running apply
unattended](../running.md) says what one run does.

## Before the timer

Set up the checkout as "Set up a checkout" in
[Running apply unattended](../running.md) says, and run `grove run` in it
once by hand to see a session through. The timer's session has the same
permissions, from the checkout's `.claude/settings.local.json`, and nobody
to answer it, so a tool those settings do not allow fails and the agent
works around it or stops. Claude Code uses the login you made with
`claude` in Terminal; a job in your user's launchd domain, as below, can
read it.

## The job

Save this as `~/Library/LaunchAgents/com.you.grove-run.ledger.plist`, with
your user's home in place of `/Users/you`, one label and file per project:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
  <key>Label</key>
  <string>com.you.grove-run.ledger</string>
  <key>ProgramArguments</key>
  <array>
    <string>/Users/you/.local/bin/grove</string>
    <string>run</string>
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

- `PATH` is all the session gets: launchd does not read your shell's
  profile. It needs the folders holding `grove`, `claude`, `git`, `gh` and
  your project's test tools; `which claude gh` in Terminal shows them.
- `StartInterval` is in seconds. launchd starts no second copy of a job
  that is still running, so a firing during a session is skipped, and
  Grove's own lock refuses a `grove run` you start by hand meanwhile.
- `grove run` stops its session after 2 hours. To change that, add
  `<string>--timeout</string>` and, say, `<string>45m</string>` after
  `<string>run</string>`.

Load it, and it fires every interval, the first one interval from now:

```sh
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.you.grove-run.ledger.plist
```

## Where to look

- `/Users/you/grove-runs/ledger.log` gets one line per firing from Grove:
  the attempt and how it ended, that the project is paused or a run is in
  progress, that there is nothing to apply, or what went wrong.
- Each attempt's record and the session's output are in the checkout's
  `.git/grove/`, named by the attempt, as the log line says.
- `launchctl print gui/$(id -u)/com.you.grove-run.ledger` shows whether
  the job is running, how many times it has, and its last exit code: 0 for
  a session that exited 0 or a firing that had nothing to do, 1 for a
  refusal or an attempt that ended any other way.

When the project is paused, deal with the attempt and run
`grove run --retry` in the checkout yourself, or merge its pull request;
the timer then picks up from there.

## Stopping it

```sh
launchctl bootout gui/$(id -u)/com.you.grove-run.ledger
```

If a session is running, launchd sends `grove run` `SIGTERM`; Grove stops
the session's processes and records the attempt as cancelled. Its branch
and files stay in the checkout for you.
