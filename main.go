package main
 
import (
    "fmt"
    "net/http"
    "io/ioutil"
    "os"
    "encoding/xml"
    "math/rand"
    "time"
    "bufio"
    "io"
    "strings"
    )
 
type Titles struct {
  Title string `xml:"title"`
  ViewItemURL string `xml:"viewItemURL"`
}

type Items struct {
  Item []Titles `xml:"item"`
}

type Response struct {
    Ack string `xml:"ack"`
    SearchResult Items `xml:"searchResult"`
}

var lineCount int 	//How many lines in file

func main() {
    
		file, err:= os.Open("words.txt")
		if err != nil { panic(err) }
	    // close on exit and check for its returned error
	    defer func() {
	        if err := file.Close(); err != nil {
	            panic(err)
	        }
	    }()

	    // Get file line length
	    lineCount, err = fileLength(*file)

	    // Handle error
	    if err != nil{
	    	fmt.Println("Error!", err)
	    }

	    //fmt.Println(randomLine(*file, lineCount))
        http.HandleFunc("/", handler)
    	http.ListenAndServe(":8000", nil)
}

func randInt(min int , max int) int {
	return min + rand.Intn(max-min)
}

func handler(w http.ResponseWriter, r *http.Request) {

	file, err:= os.Open("words.txt")
	if err != nil { panic(err) }
    // close on exit and check for its returned error
    defer func() {
        if err := file.Close(); err != nil {
            panic(err)
        }
    }()

	randomWord, err := randomLine(*file, lineCount)

	if err != nil{
		fmt.Println("%s", err)
		os.Exit(1)
	}

	url := "http://svcs.ebay.com/services/search/FindingService/v1?OPERATION-NAME=findItemsByKeywords&SERVICE-VERSION=1.0.0&SECURITY-APPNAME=ordering-85f0-4766-8873-d7ad7c4ae178&RESPONSE-DATA-FORMAT=XML&REST-PAYLOAD&keywords=" + randomWord + "&GLOBAL-ID=EBAY-US&itemFilter(0).name=MaxPrice&itemFilter(0).value=1&itemFilter(0).paramName=Currency&itemFilter(0).paramValue=USD&itemFilter(1).name=FreeShippingOnly&itemFilter(1).value=true&itemFilter(2).name=ListingType&itemFilter(2).value=FixedPrice"

	response, err := http.Get(url)

	    if err != nil {
	        fmt.Printf("%s", err)
	        os.Exit(1)
	    } else {
	        defer response.Body.Close()
	        contents, err := ioutil.ReadAll(response.Body)
	        if err != nil {
	            fmt.Printf("%s", err)
	            os.Exit(1)
	        }

	    	var output Response
	    	err = xml.Unmarshal(contents, &output)

	    	if err != nil{
	    		fmt.Printf("%s", err)
	    		os.Exit(1)
	    	}

	        rand.Seed( time.Now().UTC().UnixNano())
	        // "<a href=%s>%s</a>", 
	        item := output.SearchResult.Item[randInt(1,2)]
	        fmt.Fprintf(w,"<a href=%s>%s</a>", item.ViewItemURL, item.Title)
			}
}

// Function to return a random line from a file
// Can return an empty line if one is present in file
// File must end with an empty line
func randomLine(file os.File, fileLength int)(string, error){
	// Make a read buffer
	r:= bufio.NewReader(&file)

	// Get a random number
	rand.Seed( time.Now().UTC().UnixNano())
	
	// Return a random line
	var err error
	line := ""
	for i := 0; i <= randInt(0,fileLength-1); i++ {
		line, err = r.ReadString('\n')
	}

	if err != nil{
		return "", err
	}

	return strings.TrimSpace(line), nil
}

// Function for determining number of lines in a text file.
// Counts empty lines
func fileLength(file os.File) (int, error){
	// Make a read buffer
	r:= bufio.NewReader(&file)
	counter := 0	//Counter for lines
	var err error

	for err == nil {					//Read until an error is encountered, then break
		_, err = r.ReadString('\n')
		counter++
	}

	// Checks to make sure EOF error was found and not other error
	if err != io.EOF{
		return 0, err
	}

	//Return to beginning of file from last read
	 _, err = file.Seek(0,0)

	 if err != nil{
	 	return 0, err
	 }else{
	 	return counter, nil
	 }	
}