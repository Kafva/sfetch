package lib

const RELEASE = false
// In the release build the _SCRIPT values will
// contain the actual code to be ran and in dev mode
// it will be the path to the script 
const INFO_SCRIPT 	   = `./scripts/info.sh`
const FULL_INFO_SCRIPT = `./scripts/full_info.sh`

const COMMAND_FAILED = "COMMAND_FAILED"
const COMMAND_TIMEOUT = "COMMAND_TIMEOUT"
var HELP_STR = "-----------\nPrint a tree of hosts accessible via ssh"
const HOSTNAME_ANSI_COLOR = "\033[97m"
const LOCALHOST = "localhost"
const EXIT_ERROR = 1

// Command line options
var QUIET *bool 
var SLOW *bool 
var CONNECTION_TIMEOUT *int
var INCLUDE_HOSTNAME *bool
var DEBUG *bool 
var VERBOSE *int 
var SSH_PATH = "" 
var IGNORE_FILE *string
var CONFIG_FILE *string
