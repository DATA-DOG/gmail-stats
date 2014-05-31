# Gmail unread mail stats

A command line utility which uses gmail atom feed to get unread email statistics.

## Installation

    go get github.com/codegangsta/cli
    go get github.com/DATA-DOG/gmail-stats

The binary **gmail-stats** will be installed in **$GOPATH/bin** which should be in your $PATH.

## Usage

    gmail-stats --help

### Showing only count of unread messages

    gmail-stats -u account@gmail.com -p secretpassword -c

### Showing a list of unread messages

    gmail-stats -u account@gmail.com -p secretpassword

Prints all unread messages separated by new lines:

    Subject;From;Summary
    Subject;From;Summary
    Subject;From;Summary

Where

- Subject - is an email message subject
- From - is senders email address
- Summary - is a short summary of message body

### Running as daemon

Gmail stats can be run as daemon. It will use **send-notify** standard linux notification system
to send a notification if there are new unread emails. It does not bug you if unread email count
is not changing, but will remind you on increase.

    gmail-stats -u account@gmail.com -p secretpassword -d

By default it will check mailbox every **60** seconds, you can change this interval to 20 seconds for instance:

    gmail-stats -u account@gmail.com -p secretpassword -d -i 20

You may use **notify-osd** or other utility set which integrates together with **libnotify** on all linux platforms.

