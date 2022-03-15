# word-counter

The program reads from stdin strings containing URLs.

For each URL you need to send a GET HTTP request and count the number of occurrences of the string "Go" in the body of the response.

At the end the application displays the total number of "Go" strings found in all passed URLs, for example:

$ echo -e 'https://golang.org\nhttps://golang.org' | go run main.go
Count for https://golang.org: 9
Count for https://golang.org: 9
Total: 18
Each URL has to start being processed immediately after it is fetched, and in parallel with the fetching of the next one.

URLs should be processed in parallel, but no more than k=5 at a time.
