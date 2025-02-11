/*
Copyright © 2025 Edgar Piskernik <office@piskernik.com>
*/
package cmd

import (
	"io"
	"log"
	"net/http"
	smtplib "net/smtp"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var debugMode bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "check",
	Short: "Check is a simple uptime monitor",
	Long: `Check is a simple uptime monitor that checks a given URL in a given intervall 
and sends an email notification if the URL is not reachable. The configuration can be 
done via command line flags or a config file. The configuration file must be in YAML format.
The following flags are available:
-h, --help: Help for the check command
-U, --URL: The URL to check
-l, --log: The log file to write to
-a, --author: The author's email address of the email notification
-r, --recipient: The recipient of the email notification
-c, --config: The config file to use
-s, --smtp: The SMTP server to use
-o, --port: The port of the SMTP server
-u, --user: The user for the SMTP server
-p, --password: The password for the SMTP server
-j, --subject: The subject of the email notification
-b, --body: The body of the email notification
Example usage:
check -U https://example.com -l monitor.log`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		debugMode = false
		debugMode, _ = cmd.Flags().GetBool("debug")
		log.SetOutput(os.Stdout)

		// First read the config file (if it exists)
		if debugMode {
			log.Println("Check started with command line flags")
			log.Println("Reading configuration file...")
		}

		viper.SetConfigName(".check")      // name of config file (without extension)
		viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
		viper.AddConfigPath("/etc/check/") // path to look for the config file in
		viper.AddConfigPath("$HOME/")      // call multiple times to add many search paths
		viper.AddConfigPath(".")           // optionally look for config in the working directory
		// Find and read the config file
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if debugMode {
					log.Println("Config file not found; ignoring")
				}
			} else {
				// Config file was found but another error was produced
				log.Printf("Error reading config file: %v", err)
			}
		}
		if debugMode {
			log.Println("Using config file:", viper.ConfigFileUsed())
		}
		url := viper.GetString("URL")
		logFileName := viper.GetString("log")
		author := viper.GetString("author")
		recipient := viper.GetString("recipient")
		config := viper.GetString("config")
		smtp := viper.GetString("smtp")
		port := viper.GetString("port")
		user := viper.GetString("user")
		password := viper.GetString("password")
		subject := viper.GetString("subject")
		body := viper.GetString("body")
		debugMode = viper.GetBool("debug")

		if logFlag, err := cmd.Flags().GetString("log"); err == nil && logFlag != "" {
			logFileName = logFlag
		}

		if logFileName != "" {
			logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
			if err != nil {
				log.Printf("Failed to open log file: %v, no log will be saved", err)
				log.SetOutput(os.Stdout)
			} else {
				mw := io.MultiWriter(os.Stdout, logFile)
				log.SetOutput(mw)
			}
		}

		if debugMode {
			log.Printf("Debug mode is enabled")
			log.Printf("URL: %s\nLog: %s\nRecipient: %s\nConfig: %s\nSMTP: %s\nPort: %s\nUser: %s\nSubject: %s\nBody: %s\n",
				url, logFileName, recipient, config, smtp, port, user, subject, body)
		}

		// Check if newer values exist from command line flags, then overwrite the config values
		if urlFlag, err := cmd.Flags().GetString("URL"); err == nil && urlFlag != "" {
			url = urlFlag
		}
		if recipientFlag, err := cmd.Flags().GetString("recipient"); err == nil && recipientFlag != "" {
			recipient = recipientFlag
		}
		if configFlag, err := cmd.Flags().GetString("config"); err == nil && configFlag != "" {
			config = configFlag
		}
		if smtpFlag, err := cmd.Flags().GetString("smtp"); err == nil && smtpFlag != "" {
			smtp = smtpFlag
		}
		if portFlag, err := cmd.Flags().GetString("port"); err == nil && portFlag != "" {
			port = portFlag
		}
		if userFlag, err := cmd.Flags().GetString("user"); err == nil && userFlag != "" {
			user = userFlag
		}
		if passwordFlag, err := cmd.Flags().GetString("password"); err == nil && passwordFlag != "" {
			password = passwordFlag
		}
		if subjectFlag, err := cmd.Flags().GetString("subject"); err == nil && subjectFlag != "" {
			subject = subjectFlag
		}
		if bodyFlag, err := cmd.Flags().GetString("body"); err == nil && bodyFlag != "" {
			body = bodyFlag
		}

		if url == "" {
			log.Printf("A URL is required to check. Neither a config file entry nor a command line flag was provided.\nRemaining configuration will however be saved to the config file.")
		} else {
			http.DefaultClient.Timeout = 5 * time.Second
			response, err := http.Get(url)
			if err != nil {
				log.Printf("Failed to get website: %v", err)
			} else {
				if debugMode {
					log.Printf("Website response: %s", response.Status)
				}
				if response.StatusCode != http.StatusOK {
					if debugMode {
						log.Printf("Website not reachable: %v", response.Status)
					}
					logFile, err := os.OpenFile(logFileName, os.O_RDONLY, 0644)
					if err != nil {
						log.Printf("Error opening log file: %v", err)
					}
					defer logFile.Close()

					// Read only the last log entry
					logFile.Seek(0, io.SeekEnd)
					buf := make([]byte, 1024)
					stat, _ := logFile.Stat()
					start := stat.Size() - 1024
					if start < 0 {
						start = 0
					}
					logFile.Seek(start, io.SeekStart)
					n, _ := logFile.Read(buf)
					buf = buf[:n]
					lastLine := ""
					lines := strings.Split(string(buf), "\n")
					if len(lines) > 0 {
						lastLine = lines[len(lines)-1]
					}

					if strings.Contains(string(lastLine), subject) {
						if debugMode {
							log.Printf("Subject \"%s\" found in the last log entry. Not sending email notification.", subject)
						}
					} else {
						if debugMode {
							log.Printf("Subject \"%s\" not found in the last log entry. Sending email...", subject)
						}
						if user != "" && password != "" && smtp != "" && port != "" && recipient != "" && subject != "" {
							if author == "" {
								author = user
							}
							sendEmail(user, password, smtp, port, author, recipient, subject, body, url)
							//sendEmail(login, password, smtp, port, sender, recipient, subject, body, url string) {
						} else {
							if debugMode {
								log.Println("Cannot send email. Please provide all required information, like user, password, smtp server, smtp port, recipient email and subject.")
							}
						}
					}
				}
			}
		}
		// Save the configuration to the config file
		if config != "" {
			configPath := filepath.Dir(config)
			configFile := filepath.Base(config)
			if debugMode {
				log.Printf("Path to config: %s\n", configPath)
				log.Printf("Filename of config: %s\n", configFile)
			}
			viper.AddConfigPath(configPath)
			viper.SetConfigName(configFile)
			viper.SetConfigFile(config)
		}

		viper.WriteConfig()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	//rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.check.yaml)")
	rootCmd.PersistentFlags().StringP("URL", "U", "", "URL to check")
	rootCmd.PersistentFlags().StringP("log", "l", "", "Log file to write to")
	rootCmd.PersistentFlags().StringP("author", "a", "", "Author's email address of the email notification")
	rootCmd.PersistentFlags().StringP("recipient", "r", "", "Recipient of the email notification")
	rootCmd.PersistentFlags().StringP("config", "c", "", "Config file to use")
	rootCmd.PersistentFlags().StringP("smtp", "s", "", "SMTP server to use")
	rootCmd.PersistentFlags().StringP("port", "o", "", "Port of the SMTP server")
	rootCmd.PersistentFlags().StringP("user", "u", "", "User login for the SMTP server")
	rootCmd.PersistentFlags().StringP("password", "p", "", "Password for the SMTP server")
	rootCmd.PersistentFlags().StringP("subject", "j", "", "Subject of the email notification")
	rootCmd.PersistentFlags().StringP("body", "b", "", "Body of the email notification")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug mode")

	viper.BindPFlag("URL", rootCmd.PersistentFlags().Lookup("URL"))
	viper.BindPFlag("log", rootCmd.PersistentFlags().Lookup("log"))
	viper.BindPFlag("author", rootCmd.PersistentFlags().Lookup("author"))
	viper.BindPFlag("recipient", rootCmd.PersistentFlags().Lookup("recipient"))
	viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	viper.BindPFlag("smtp", rootCmd.PersistentFlags().Lookup("smtp"))
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user"))
	viper.BindPFlag("password", rootCmd.PersistentFlags().Lookup("password"))
	viper.BindPFlag("subject", rootCmd.PersistentFlags().Lookup("subject"))
	viper.BindPFlag("body", rootCmd.PersistentFlags().Lookup("body"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

func sendEmail(login, password, smtp, port, sender, recipient, subject, body, url string) {
	if debugMode {
		log.Printf("Sending email with %s at server %s:%s from %s to %s with subject: %s about %s: %s", login, smtp, port, sender, recipient, subject, body, url)
	}

	auth := smtplib.PlainAuth("", login, password, smtp)
	to := []string{recipient}
	msg := []byte("From: " + sender + "\r\n" +
		"To: " + recipient + "\r\n" +
		"Subject: " + subject + url + "\r\n" +
		"\r\n" + body + url + "\r\n")

	date := time.Now().Format(time.RFC1123Z)

	err := smtplib.SendMail(smtp+":"+port, auth, login, to, msg)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		log.Printf("%s - %v: %s\n", date, err, url)
	} else {
		log.Printf("%s - Email sent successfully: %s\n", date, subject)
	}
}
