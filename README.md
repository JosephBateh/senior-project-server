# senior-project-server
Server for Spotify Playlist Manager

# Environment Variables
An environment variable file needs to be stored at the root of the directory and maintain these variables:

```
SPOTIFY_ID=<spotify_client_id>
SPOTIFY_SECRET=<spotify_secret>
MLAB_LOGIN=<db-driver-url>
MLAB_DB=<db-name>
PRODUCTION=<boolean>
REDIRECT_URI=<location of web app>
```

To run the server, navigate to the directory and run:

```
go get
go install
go run main.go
```