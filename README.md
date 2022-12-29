# plailyst
PlaiLyst applies custom filtering on given youtube channels and creates custom Youtube playlist

# Note
This is a super beta version as of now. It has no UI and no deployment script as of now. I use this tool for my personal usage and I am not planning to create a user interface for it. Feel free to add one, if you'd like to contribute.

# How to configure
Create an OAuth 2.0 Client ID for a web app at: https://console.cloud.google.com/apis/credentials
Make sure to add your own email to the allowed lists, if you're using it for testing.

Using the same Google account, create an empty play list on Youtube.

Save the secret file as `client_secret.json` at the top of your directory.

Create a custom config file with custom teams, filters and Youtube channels and name it `configs.yaml`. Something similar to what I use:
```
---
channels:
  # Get channel IDs from https://commentpicker.com/youtube-channel-id.php
  - beinsportsusa:
    id: UC0YatYmg5JRYzXJPxIdRd8g
  - cbssportsgolazo:
    id: UCET00YnetHT7tOpu12v8jxg
  - NBCSports:
    id: UCqZQlzSHbVJrwrn5XvzrzcA
  - ESPNFC:
    id: UC6c1z7bA__85CIWZ_jpCK-Q
# Make sure above channels show matches from below teams :)
teams:
  - AC Milan
  - PSG
  - Manchester United
  - Man United
  - Manchester City
  - Man City
  - Arsenal
  - Liverpool
  - Real Madrid
  - Barcelona
# These are must terms in each video title. I personally only want to watch highlights. Feel free to make it more strict for yourself.
terms:
  - highlights
  - highlight
# The playlist ID that you created in an earlier step.
playlist: PLppTncK5Ux-fVwFzMrbzd3qpzweW6Lxyx
```

Install go-task if you don't have it already: https://taskfile.dev/installation/

Run: `go-task run`. 

On your browser go to: http://localhost:8765/login

Go through the usual OAuth 2.0 authentication flow with Google and after a while you should be redirected to the main page with a message saying your playlist was updated.
