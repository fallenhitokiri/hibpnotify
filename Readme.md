# Have I Been Pwned Notify - hibpnotify
This little tool accepts a CSV file with email addresses as input and regularly checks
[https://haveibeenpwned.com](https://haveibeenpwned.com) if they were compromised in a data
breach. If it detects a newly breached account it will send a short info mail to the breached
account and a third, freely configurable email address. Knowing an account was breached allows
to put some special logging for this account in place. Maybe forcing a password change in case
a password was reused. Or simply as an information for people who were victim of a breach in
your organisation so they know and can react.

## Installation and configuration
With the first official release there will be binaries for at least Linux and macOS, right now
you have to run and build the service manually.

```
go get github.com/fallenhitokiri/hibpnotify
```

Now you can initialise a new configuration file.

```
cd $GOPATH/src/github.com/fallenhitokiri/hibpnotify/cmd/
go run hibpnotify.go -init -config path/to/config.json
```

You can no edit the config file and run hibbpnotify

```
cd $GOPATH/src/github.com/fallenhitokiri/hibpnotify/cmd/
go run hibpnotify.go -config path/to/config.json
```

## Current state
### What is working?
- csv input
- email notifications

### Currently WIP
- using Google Apps as data source

### Planned
- Active Directory support
- Slack notifications
- web interface for users to submit their private email addresses

## License
MIT 