package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
)

type stocksArray []singleStock

type singleStock struct {
	Name   string `json:"t"`
	Amount string `json:"l"`
}

var (
	err      error
	stocks   stocksArray
	response *http.Response
	body     []byte
)

const stocktickerTemplateHTML = ` 
<!DOCTYPE html>
<html lang="en">

<head>

    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <meta name="description" content="Go, stocks, ticker, Google, Apple, Microsoft, sample app">
    <meta name="author" content="Thomas Jaensch">

    <title>Go Stock Ticker</title>

    <!-- Bootstrap Core CSS -->
    <link href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/3.3.4/css/bootstrap.min.css" rel="stylesheet">
    <!-- Custom CSS -->
    <link href="/css/one-page-wonder.css" rel="stylesheet">
    <!-- Bootstrap JS -->
    <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.4/js/bootstrap.min.js"></script>
    <!-- jQuery -->
    <script src="https://ajax.googleapis.com/ajax/libs/jquery/2.1.3/jquery.min.js"></script>

    <!-- HTML5 Shim and Respond.js IE8 support of HTML5 elements and media queries -->
    <!-- WARNING: Respond.js doesn't work if you view the page via file:// -->
    <!--[if lt IE 9]>
        <script src="https://oss.maxcdn.com/libs/html5shiv/3.7.0/html5shiv.js"></script>
        <script src="https://oss.maxcdn.com/libs/respond.js/1.4.2/respond.min.js"></script>
    <![endif]-->

</head>

<body>

    <!-- Full Width Image Header -->
    <header class="header-image">
        <div class="headline">
            <div class="container">
                <h2>{{range .}}
                    {{.Name}}: ${{.Amount}} <br>
                {{end}}</h2>
                <p>Stock ticker* demo app written in Go with Bootstrap, code on <a href="https://github.com/tjaensch/stockticker_google_sdk">Github</a></p>
                <p><small>*Not ticking? Then it's probably after hours on Wall Street.</small></p>
            </div>
        </div>
    </header>

        <!-- Footer -->
        <footer>
            <div class="row">
                <div class="col-lg-12 text-center">
                    <p><small>&copy; Thomas Jaensch 2015</small></p>
                </div>
            </div>
        </footer>

    </div>
    <!-- /.container -->

    <!-- Ajax call to refresh page periodically -->
    <script type="text/javascript">
        function startRefresh() {
        $.get('', function(data) {
            $(document.body).html(data);    
        });
        }
        $(function() {
            setTimeout(startRefresh,1000);
        });
    </script>

</body>

</html>
`

var stockTemplate = template.Must(template.New("stockticker").Parse(stocktickerTemplateHTML))

func stockticker(w http.ResponseWriter, r *http.Request) {
	// Use http://finance.google.com/finance/info?client=ig&q=NASDAQ:GOOG to get a JSON response
	response, err = http.Get("http://finance.google.com/finance/info?client=ig&q=NASDAQ:GOOG,NASDAQ:AAPL,NASDAQ:MSFT")
	if err != nil {
		fmt.Println(err)
	}
	defer response.Body.Close()

	// Read the data into a byte slice
	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	// Remove whitespace from response
	data := bytes.TrimSpace(body)

	// Remove leading slashes and blank space to get byte slice that can be unmarshaled from JSON
	data = bytes.TrimPrefix(data, []byte("// "))

	// Unmarshal the JSON byte slice to a predefined struct
	err = json.Unmarshal(data, &stocks)
	if err != nil {
		fmt.Println(err)
	}

	// Parse struct data to template
	tempErr := stockTemplate.Execute(w, stocks)
	if tempErr != nil {
		http.Error(w, tempErr.Error(), http.StatusInternalServerError)
	}
}

func main() {

	fs := http.FileServer(http.Dir("css"))
	http.Handle("/css/", http.StripPrefix("/css/", fs))
	http.HandleFunc("/", stockticker)
	http.ListenAndServe(":8080", nil)

}
