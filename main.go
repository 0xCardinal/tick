package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	_ "github.com/mattn/go-sqlite3"
)

// Setup - Create a database if not exist (tt -s)
func setup() {

	// Header
	fmt.Print(`
      _____             _____          
     /\    \           /\    \         
    /::\    \         /::\    \        
    \:::\    \        \:::\    \       
     \:::\    \        \:::\    \      
      \:::\    \        \:::\    \     
       \:::\    \        \:::\    \    
       /::::\    \       /::::\    \   
      /::::::\    \     /::::::\    \  
     /:::/\:::\    \   /:::/\:::\    \ 
    /:::/  \:::\____\ /:::/  \:::\____\
   /:::/    \::/    //:::/    \::/    /
  /:::/    / \/____//:::/    / \/____/ 
 /:::/    /        /:::/    /          
/:::/    /        /:::/    /           
\::/    /         \::/    /            
 \/____/           \/____/             
                                       
									  
Welcome to Tick!
	
As name suggests this is a todo list that can be accessed on a terminal.
	
Type 'tt -h' for more info.`)

	fmt.Print("\nDo you want to setup tick? (y/n): ")
	reader := bufio.NewReader(os.Stdin)
	consent, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	database, err := sql.Open("sqlite3", "./tick.db")
	consent = strings.Replace(consent, "\n", "", -1)
	if strings.Compare(consent, "y") == 0 {
		color.Set(color.FgCyan, color.Bold)
		if err != nil {
			panic(err)
		} else {
			fmt.Print("Initilizing Database...")
		}
		stmt, stmterr := database.Prepare("CREATE TABLE IF NOT EXISTS AllTasks (ID INTEGER PRIMARY KEY, TASKS TEXT, DONE INTEGER DEFAULT 0, URGENCY INTEGER DEFAULT 0)") // 0 = not
		if stmterr != nil {
			panic(stmterr)
		} else {
			stmt.Exec()
		}

		fmt.Println("Done.\n")
		color.Unset()
		color.Set(color.FgGreen, color.Bold)

		fmt.Println("Setup Complete...\nType 'tt -h' for help. \nWe hope Tick makes you more productive!")

		color.Unset()
	} else if strings.Compare(consent, "n") == 0 {
		color.Red("Hope to see you soon!")
	}
}

// Add Task (tt -a)
// Open a shell, that lets user to add tasks to db
// For Urgent - !urgent
// For Exit - !exit
func addTask() {
	database, err := sql.Open("sqlite3", "./tick.db")
	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(os.Stdin)
	red := color.New(color.FgRed).SprintFunc()
	fmt.Printf("Welcome to Tick!\n================\nType %s at the start, for urgent tasks.\nType %s when you are done adding tasks.\nAdd Tasks,\n", red("!urgent"), red("!exit"))
	i := 1
	for {
		fmt.Print(">")
		task, _ := reader.ReadString('\n')
		task = strings.Replace(task, "\n", "", -1)

		if strings.Compare(task, "!exit") == 0 {
			fmt.Println("Tasks added to list! \nExiting...")
			break
		} else {
			if strings.HasPrefix(task, "!urgent ") {
				task = strings.Trim(task, "!urgent ")
				stmt, _ := database.Prepare("INSERT INTO AllTasks (TASKS, URGENCY) VALUES (?,?)")
				stmt.Exec(task, 1)
			} else {
				stmt, _ := database.Prepare("INSERT INTO AllTasks (TASKS) VALUES (?)")
				stmt.Exec(task)
			}
		}
		i++
	}
}

// Show Urgent Tasks
func showUrgent() {
	database, err := sql.Open("sqlite3", "./tick.db")
	if err != nil {
		panic(err)
	}
	// Colors
	heading := color.New(color.FgWhite, color.BgRed, color.Bold)
	subHead := color.New(color.BgBlue, color.FgWhite, color.Bold)
	urgentList := color.New(color.FgGreen, color.Bold)

	heading.Print(" Welcome to Tick! ")
	fmt.Println("\n------------------")

	// Count Check
	rows, err := database.Query("SELECT COUNT(*) as count FROM AllTasks WHERE URGENCY=1 and DONE=0")
	if checkCount(rows) > 0 {
		subHead.Print(" URGENT ")
		fmt.Println("\n")
		rows, err = database.Query("SELECT ID, TASKS FROM AllTasks WHERE URGENCY=1 and DONE=0")
		// Data Retrival
		var id int
		var task string

		for rows.Next() {
			err = rows.Scan(&id, &task)
			if err != nil {
				panic(err)
			}
			urgentList.Printf("[ %d ]  %s\n", id, task)
		}

		rows.Close()
	} else {
		subHead.Print("  No URGENT Tasks Present! ")
		fmt.Println()
		rows.Close()
	}
}

func checkCount(rows *sql.Rows) (count int) {
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			panic(err)
		}
	}
	return count
}

func showTasks() {
	// Urgent Tasks
	showUrgent()

	// Normal Tasks
	database, err := sql.Open("sqlite3", "./tick.db")
	if err != nil {
		panic(err)
	}
	// Colors
	subHead := color.New(color.BgBlue, color.FgWhite, color.Bold)
	normalList := color.New(color.Bold)
	fmt.Println("----------------------")
	rows, err := database.Query("SELECT COUNT(*) as count FROM AllTasks where URGENCY=0 and DONE=0")
	if checkCount(rows) > 0 {
		subHead.Print(" NORMAL ")
		fmt.Println("\n")
		rows, err = database.Query("SELECT ID, TASKS FROM AllTasks WHERE URGENCY=0 and DONE=0")
		// Data Retrival
		var id int
		var task string

		for rows.Next() {
			err = rows.Scan(&id, &task)
			if err != nil {
				panic(err)
			}
			normalList.Printf("[ %d ]  %s\n", id, task)
		}

		rows.Close()
	} else {
		subHead.Print("  No Tasks Present! ")
		fmt.Println()
		rows.Close()
	}
	fmt.Println()
}

func showDeleted() {
	database, err := sql.Open("sqlite3", "./tick.db")
	if err != nil {
		panic(err)
	}
	subHead := color.New(color.BgBlue, color.FgWhite, color.Bold)
	deadList := color.New(color.FgRed, color.Bold)
	heading := color.New(color.FgWhite, color.BgRed, color.Bold)

	heading.Print(" Welcome to Tick! ")
	fmt.Println("\n------------------")

	rows, err := database.Query("SELECT COUNT(*) as count FROM AllTasks where DONE=1")
	if checkCount(rows) > 0 {
		subHead.Print("  COMPLETED TASKS  ")
		fmt.Println("\n")
		rows, err = database.Query("SELECT ID, TASKS FROM AllTasks WHERE DONE=1")
		// Data Retrival
		var id int
		var task string

		for rows.Next() {
			err = rows.Scan(&id, &task)
			if err != nil {
				panic(err)
			}
			deadList.Printf("[ %d ]  %s\n", id, task)
		}
		rows.Close()
	} else {
		subHead.Print("  No Completed Tasks!  ")
		fmt.Println("\n")
		os.Exit(0)
	}
}

func delete() {
	database, err := sql.Open("sqlite3", "./tick.db")
	if err != nil {
		panic(err)
	}
	subHead := color.New(color.BgRed, color.FgWhite, color.Bold)
	rows, err := database.Query("SELECT COUNT(*) as count FROM AllTasks WHERE DONE=0")
	if checkCount(rows) > 0 {
		indices := flag.Args()
		for i := 0; i < len(indices); i++ {
			if _, err := strconv.Atoi(indices[i]); err == nil {
				if stmt, err := database.Prepare("UPDATE AllTasks set DONE = 1 WHERE ID = ?"); err == nil {
					stmt.Exec(indices[i])
				}
			}

		}
		subHead.Print("  Tasks Completed!  ")
		fmt.Println()
	} else {
		subHead.Print("  All tasks are completed!  ")
		fmt.Println()
	}
	fmt.Println()
}

func init() {
	flag.Usage = func() {
		h := `
Welcome to Tick!
A handy todo list, accessible from Terminal! (Terminal Todo)

Usage:
  tt [OPTIONS] {-d|--delete} [args...]
  
Options:
					Shows the complete to-do list
  -s, --setup 		Setup Tick for use (initilize database)
  -a, --add			Interactive shell to add tasks.
						> use !urgent at the start of the tasks to specify urgency
						> use !exit at the start to exit the interactive shell
  -d, --delete		Deletes tasks based on tasks' indices
  -u, --urgent		Prints the urgent tasks only
	  --deleted		Prints the deleted tasks only
	  --version		Prints the version information

Example:
  tt -d 1 2 5		// Deletes the tasks with index 1,2 and 5
  
  `
		fmt.Fprintf(os.Stderr, h)
	}
}

func main() {

	var (
		setupFlag   bool
		addFlag     bool
		delFlag     bool
		urgentFlag  bool
		showDelFlag bool
		versionFlag bool
	)

	flag.BoolVar(&setupFlag, "s", false, "")
	flag.BoolVar(&setupFlag, "setup", false, "")
	flag.BoolVar(&addFlag, "a", false, "")
	flag.BoolVar(&addFlag, "add", false, "")
	flag.BoolVar(&delFlag, "d", false, "")
	flag.BoolVar(&delFlag, "delete", false, "")
	flag.BoolVar(&urgentFlag, "u", false, "")
	flag.BoolVar(&urgentFlag, "urgent", false, "")
	flag.BoolVar(&showDelFlag, "deleted", false, "")
	flag.BoolVar(&versionFlag, "v", false, "")
	flag.BoolVar(&versionFlag, "version", false, "")

	flag.Parse()

	TickVersion := " v1.0.0"
	if versionFlag {
		subHead := color.New(color.BgRed, color.FgWhite, color.Bold)
		subHead.Print("  Welcome to Tick", TickVersion, "  ")
		fmt.Println()

		subHead = color.New(color.BgBlue, color.FgWhite, color.Bold)
		subHead.Print("  Developed by @krAshwin  ")
		fmt.Println()
		os.Exit(1)
	}

	if addFlag {
		addTask()
		os.Exit(1)
	}

	if setupFlag {
		setup()
		os.Exit(1)
	}

	if urgentFlag {
		showUrgent()
		os.Exit(1)
	}

	if showDelFlag {
		showDeleted()
		os.Exit(1)
	}

	if delFlag {
		delete()
		os.Exit(1)
	}

	if _, err := os.Stat("./tick.db"); os.IsNotExist(err) {
		setup()
		os.Exit(1)
	} else {
		showTasks()
		os.Exit(1)
	}
}
