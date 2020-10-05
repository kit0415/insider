package supervisors

import (
	"bytes"
	"encoding/json"
	"fmt"
	"insider/util"
	"log"
	"os"
	"strconv"
	"strings"
)

// SourceCodeInfo holds information about the received code to analyze
type SourceCodeInfo struct {
	Path         string
	Tech         string
	SastID       string
	Version      string
	ComponentID  string
	PhysicalPath string

	// Hashing
	MD5Hash    string
	SHA1Hash   string
	SHA256Hash string
}

// reportResult will handle the logic to upload the final report about the source code
// being analyzed to somewhere it can be feed into other tools or used by hand.
// By default, in development environment it will save the report to a file in the
// current directory with the name of report-[here will be the SAST ID].json
func reportResult(bReport []byte, ignoreWarnings bool,targetSaveLoc string) error {
	// Running on debug mode, should avoid communication with the Console.
	log.Println("target log file: %s",targetSaveLoc)
	var reportFilename string
	if !ignoreWarnings {
		//currentTime := time.Now()
		//reportFilename = fmt.Sprintf("report-%v.json", currentTime.Format("20060102150405"))
		log.Println("Checking Dir: %s",targetSaveLoc)
		if _,err :=os.Stat("/mnt/c/git/uploads/"+targetSaveLoc+"/report"); os.IsNotExist(err){
			log.Println("Dir Not Exist")
			os.Mkdir("/mnt/c/git/uploads/"+targetSaveLoc+"/report",0666)
		}
		reportFilename = fmt.Sprintf("report-%s.json", targetSaveLoc)
	} else {
		reportFilename = "report.json"
	}

	pwd, err := os.Getwd()
	log.Println("Dir is ",pwd)
	if err != nil {
		return err
	}

	log.Printf("Json Report %s/%s\n", pwd+"/report", reportFilename)
	//change the export location here
	file, err := os.OpenFile("/mnt/c/git/uploads/"+targetSaveLoc+"/report/"+reportFilename, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return fmt.Errorf("Problems writing the report to the JSON file: %w", err)
	}
	defer file.Close()

	// Makes sure we start to write in the beginning of the file
	// and overwriting anything that was previously inside the file
	if err := file.Truncate(0); err != nil {
		return fmt.Errorf("Problems writing the report to the JSON file: %w", err)
	}

	var outputBuffer bytes.Buffer

	// This will format the output according to the
	// JSON specification to be easier to read by a human
	if err := json.Indent(&outputBuffer, bReport, "", " "); err != nil {
		return err
	}

	formattedContent, err := strconv.Unquote(strings.Replace(strconv.Quote(outputBuffer.String()), ``, ``, -1))
	if err != nil {
		return fmt.Errorf("Problems writing the report to the JSON file: %w", err)
	}

	bytesWritten, err := file.Write([]byte(formattedContent))
	if err != nil {
		return fmt.Errorf("Problems writing the report to the JSON file: %w", err)
	}

	log.Printf("Json Report %v bytes written successfully\n", util.ByteCountSI(int64(bytesWritten)))

	return nil
}
