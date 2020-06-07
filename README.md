# Roll Of Dice

This application was created to meet the test of a company using the Go programming language, where to display the form using the default <template> go. At first I thought of using WebSocket because of the problem with the multiplayer provisions. Because if I use web socket I don't need to create a database and the data is realtime. But because I have never used WebSocket, so I decided to use sqlite database. I know this application is far from perfect because of the shortcomings that I have, with the experience of using the Go programming language under 1 year. But I am very happy to be able to work on projects with the Go programming language even if it's just a test.

## Prerequisites

**Install Go v 1.11+**

Please check the [Official Golang Documentation](https://golang.org/doc/install) for installation.

## Installation

**Clone this repository**

```bash
git clone git@github.com:humamalamin/dice.git
# Switch to the repository folder
cd dice
```
**Run Dice Application**

```bash
go run main.go
```
Open your browser and input http://localhost:9001

## References
* [HTML in Golang](https://medium.com/@thedevsaddam/easy-way-to-render-html-in-go-34575f858026)
* [Disable Button in Javascript](https://flaviocopes.com/how-to-disable-button-javascript/)
* [Websocket](https://gowebexamples.com/websockets/)
* [Using Sqlite in Golang](https://www.thepolyglotdeveloper.com/2017/04/using-sqlite-database-golang-application/)

## Contributing

When contributing to this repository, please note we have a code standards, please follow it in all your interactions with the project.

#### Steps to contribute

1. Clone this repository.
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Submit pull request.