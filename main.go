package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
)


func passCmd(args []string) ([]string, error) {

    if len(args) <= 2 {
        return nil, errors.New("Insufficient args provided, usage: ftrav <command> <path/key>")
    }
    return args[1:], nil

}


func readMap(jsonPath string) map[string]string {
    file, err := os.Open(jsonPath)
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }                                
    
    defer file.Close()
    
    var pathMap map[string]string
    decoder := json.NewDecoder(file)
    if err := decoder.Decode(&pathMap); err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }

    return pathMap

}


func changeDirectory(cmd []string, allPaths map[string]string, jsonPath string) {
    
    if len(allPaths) == 0 {
        fmt.Printf("No fast travel locations set, set locations by navigating to desired destination directory and using 'ftrav set <key>' ")
        os.Exit(1)
    }
    path := allPaths[cmd[1]]
    
    err := clipboard.WriteAll("cd "+"'"+path+"'")
    if err != nil {
        fmt.Printf("Fast travel failed! %v", err)
        os.Exit(1)
    }

    fmt.Printf("cd %v copied to clipboard, paste to fast travel there", path)

}


func ensureJSON(filepath string) {
    
    _, err := os.Stat(filepath)
    if err == nil {
        return
    }

    if !os.IsNotExist(err) {
        fmt.Println("Error: ", err)
        os.Exit(1)
    }

    newFile, err := os.Create(filepath)
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(1)
    }
    
    defer newFile.Close()

    _, err = newFile.WriteString("{}")
    if err != nil {
        fmt.Println("Error: ", err)
        os.Exit(1)
    }

}




func setDirectoryVar(cmd []string, allPaths map[string]string, jsonPath string) {
     
    path, err := os.Getwd()
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
    allPaths[cmd[1]] = path
    
    jsonData, err := json.MarshalIndent(allPaths, "", "  ")
    if err != nil {
        fmt.Println("Error marshalling JSON:", err)
        os.Exit(1)
    }
 
    file, err := os.Create(jsonPath)
    if err != nil {
        fmt.Println("Error creating file:", err)
        os.Exit(1)
    }

    defer file.Close()
   
    _, err = file.Write(jsonData)
    if err != nil {
        fmt.Println("Error writing JSON to file:", err)
        os.Exit(1)                                                                                                        
    }
}

// func displayAllPaths(cmd []string, allPaths map[string]string, jsonPath string) {}
// type cmdArgs struct {
//     cmd []string 
//     allPaths map[string]string  
//     jsonPath string 
// }



// map of available ftrav commands 
var availCmds = map[string]func(cmd []string, allPaths map[string]string, jsonPath string) {
    "to": changeDirectory,
    "set": setDirectoryVar,
    // "ls": displayAllPaths
}




func main() {
    
    // read in json file
    exePath, err := os.Executable()
    if err != nil {
        fmt.Println("Error:", err)
        os.Exit(1)
    }
 
    jsonDirPath := filepath.Dir(exePath)
    jsonPath := jsonDirPath + "\\fastTravel.json"
    ensureJSON(jsonPath)
    allPaths := readMap(jsonPath)
 
    // sanitize input
    inputCommand, err := passCmd(os.Args)
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    action := inputCommand[0]


    // execute user provided action
    exeCmd, ok := availCmds[action]
    if !ok {
        fmt.Println("Invalid command, use 'help' for available commands.")
        os.Exit(1)
    }

    exeCmd(inputCommand, allPaths, jsonPath)

}



