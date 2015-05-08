# Go Challenge 1 
See http://golang-challenge.com/go-challenge1/
Scout's honor, I didn't peak at winners' solutions befor completing my solution.

## Lessons learned after viewing winners' solutions
- Use drum.go to document splice file structure; I placed my comments in 
decoder.go
- Do not ignore Big/Little endian issues
- io.LimitReader. Read the Go lib documentation
- I like the way that the winner hid details of header structs within funcs.
- Where possible, avoid declaring variables as []btye if there's a type that
better reflects what you're storing.
