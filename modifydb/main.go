// modifydb is a command which can perform various
// modifications on the grades database
package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var (
	flags = pflag.NewFlagSet("", pflag.ContinueOnError)

	dbPathFlag  string
	commandFlag string
	problemFlag string
	studentFlag string
	pointsFlag  float64
	commentFlag string
	resetFlag   bool
	verboseFlag bool
	filesFlag   bool
)

const (
	exitUsage = 1 + iota // usage errors
	exitLogic            // logic errors (such as nonexistent students)
	exitError            // filesystem errors and so forth
)

var usage = `Usage: modifydb [flags...] [[--files | -f] [<file> [...]]]
Optional flags:
                -c, --command [grade | view] - The command to execute
                -p, --problem <problem>      - The problem to grade
                -s, --student <username>     - The student to grade
                    --points <points>        - The number of points to assign on the problem
                    --comment <comment>      - An optional comment on the grade
                -f, --files [<file> [...]]   - Names of files that were included in this handin
                -v                           - Be verbose
                    --reset                  - Reset the database to its initial state
                    --db-path <path>         - An optional custom database path (overrides default)

Commands:
				   grade - Alter a grade in the database
                   view  - Print the contents of the database
`

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(exitUsage)
	}

	if verboseFlag {
		fmt.Println("Verbose mode enabled...")
	}

	db, err := ReadDB(dbPathFlag)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not read db: %v\n", err)
		os.Exit(exitError)
	}

	// reset database instead of performing other operations
	if resetFlag {
		if verboseFlag {
			fmt.Printf("Resetting database at %s\n", dbPathFlag)
		}
		err := DestroyDB(dbPathFlag)
		if err != nil {
			fmt.Printf("could not reset db: %v\n", err)
		}
		os.Exit(exitError)
	}

	switch commandFlag {
	case "grade":
		if verboseFlag {
			fmt.Printf("Command: grade\n")
			fmt.Printf("[%s] Problem: %s, Student: %s, Points: %d, Comment: %s\n", commandFlag, problemFlag, studentFlag, pointsFlag, commentFlag)
		}

		if pointsFlag < 0 {
			fmt.Fprintln(os.Stderr, "Error: expected --points flag with grade command.")
			os.Exit(exitUsage)
		} else if studentFlag == "" {
			fmt.Fprintln(os.Stderr, "Error: expected --student flag with grade command.")
			os.Exit(exitUsage)
		} else if problemFlag == "" {
			fmt.Fprintln(os.Stderr, "Error: expected --problem flag with grade command.")
			os.Exit(exitUsage)
		}

		var files []string
		if filesFlag {
			// these are all of the non-flag arguments
			files = flags.Args()
		}

		// attempt to add the given grade
		if !grade(problemFlag, studentFlag, pointsFlag, commentFlag, files[:], db) {
			fmt.Fprintf(os.Stderr, "could not grade: unknown assignment %s\n", problemFlag)
			os.Exit(7)
		}

		// write the modified database back to the file
		err = WriteDB(db, dbPathFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "could not write changes to db: %v\n", err)
			os.Exit(exitError)
		} else if verboseFlag {
			fmt.Println("Database write successful!")
		}
	case "view":
		// if the --student flag is provided, filter
		// by the given student
		if flags.Lookup("student").Changed {
			db = db.FilterUser(studentFlag)
		}
		b, err := json.MarshalIndent(db, "", "\t")
		if err == nil {
			fmt.Println(string(b))
		} else {
			fmt.Fprintf(os.Stderr, "could not encode db as json: %v\n", err)
		}
	default:
		if flags.Lookup("command").Changed {
			fmt.Fprintf(os.Stderr, "unrecognized command: %q\n", commandFlag)
		} else {
			fmt.Fprintln(os.Stderr, "must specify command (use --command flag)")
		}
		os.Exit(exitUsage)
	}
}

func init() {

	flags.StringVarP(&commandFlag, "command", "c", "", "The command to execute [grade | view]")         // command
	flags.StringVarP(&problemFlag, "problem", "p", "", "The problem to grade")                          // problem
	flags.StringVarP(&studentFlag, "student", "s", "", "The student to grade")                          // student
	flags.Float64Var(&pointsFlag, "points", -1, "The number of points to assign on the problem")        // points
	flags.StringVar(&commentFlag, "comment", "", "An optional comment on the grade")                    // comment
	flags.BoolVarP(&filesFlag, "files", "f", false, "Names of files that were included in this handin") // files
	flags.BoolVarP(&verboseFlag, "verbose", "v", false, "Be verbose.")                                  // verbose
	flags.BoolVar(&resetFlag, "reset", false, "Reset the database to its initial state.")               // reset
	flags.StringVar(&dbPathFlag, "db-path", "/course/cs666/.db/db", "The path to the database file.")   // database path

}
