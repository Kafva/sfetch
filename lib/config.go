package lib

const RELEASE = false
// In the release build the INFO_SCRIPT values will
// contain the actual code to be ran and in dev mode
// it will be the path to the script 
const INFO_SCRIPT 	   = `./scripts/info.sh`
const FULL_INFO_SCRIPT = `./scripts/full_info.sh`

const WINDOWS_FULL_INFO = "wmic baseboard get product ; wmic os get name ; wmic os get version"
const WINDOWS_INFO 		= "wmic os get name ; wmic os get version"
const WINDOWS_PREFIX	= "\033[96m \033[0m "

const COMMAND_FAILED = "COMMAND_FAILED"
const COMMAND_TIMEOUT = "COMMAND_TIMEOUT"
const COMMAND_IN_PROGRESS = "COMMAND_IN_PROGRESS"
var HELP_STR = "-----------\nPrint a tree of hosts accessible via ssh"
const HOSTNAME_ANSI_COLOR = "\033[97m"
const LOCALHOST = "localhost"
const EXIT_ERROR = 1


// Command line options
var SKIP_WINDOWS_CHECK *bool
var QUIET *bool 
var SLOW *bool 
var CONNECTION_TIMEOUT *int
var INCLUDE_HOSTNAME *bool
var DEBUG *bool 
var VERBOSE *int 
var SSH_PATH = "" 
var IGNORE_FILE *string
var CONFIG_FILE *string

var TARGETS *string
var TARGET_MAP = make(map[string]struct{})