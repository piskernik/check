# Project Title

Check is a simple command line application to check the accessability to a certain website.

## Description

If the website is not reachable, you can send an automated e-mail and/or log it in a custom log file. Everything is controlled by commandline flags or by a simple YAML configuration file `.check`, see the `.check.example`for further details. Anytime you call the programm with flags, the flags are stored in the configuration file, so you don't need to repeat anything, which is already set-up.


## Getting Started

### Dependencies

`check` uses the following Go modules
* cobra for command line interaction
* viper for configuration management
* smtplib for sending e-mails

### Installing

* Save the appropriate version according to your system in your home directory.
```
wget www.github.com/piskernik/check/releases/check_arm_mac_0.1
mv check_arm_mac_0.1 check
chmod u+x check
```

### Executing program

`check -U https://www.example.com -l check.log`

## Help

```
# Use
	Use:   "check",
	Short: "Check is a simple uptime monitor",
	Long: `Check is a simple uptime monitor that checks a given URL in a given intervall 
and sends an email notification if the URL is not reachable. The configuration can be 
done via command line flags or a config file. The configuration file must be in YAML format.
The following flags are available:
- URL: The URL to check
- log: The log file to write to
- recipient: The recipient of the email notification
- config: The config file to use
- SMTP: The SMTP server to use
- port: The port of the SMTP server
- user: The user for the SMTP server
- password: The password for the SMTP server
- intervall: The intervall in minutes to check the URL
- subject: The subject of the email notification
- body: The body of the email notification
Example usage:
check -U https://example.com -l monitor.log -r`,
```

## Authors
Contributors
Edgar Piskernik

## Version History

* 0.1
    * Initial Release

## License

This project is licensed under the MIT License - see the LICENSE.md file for details

## Acknowledgments

Inspiration, code snippets, etc.
* [Uptime Monitor in Python](https://github.com/Manu-Abuya/Website-Uptime-Monitor/blob/master/WebsiteUptimeMonitor/website_monitor.py)
* [cobra](https://github.com/spf13/cobra)
* [viper](https://github.com/spf13/viper)