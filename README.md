# plailyst
PlaiLyst applies custom filtering on given youtube channels and creates custom Youtube playlist.

I create this tool to add my favorite football clubs match highlights to a custom Youtube's playlist so that I can just sit on the couch and watch all the highlights without searching for them and potentially getting spoilers by finding out about the results. 

In case you share the same interests as me and live in the United States, you can simply add this playlist to your account: https://youtube.com/playlist?list=PLppTncK5Ux-fVwFzMrbzd3qpzweW6Lxyx

The content of this playlist gets updated by me manually. I have to figure out a way to update youtube playlist contents using a service account, so that it can be periodically updated in the background without my intervention.

If you're not interested in the teams that I follow or you cannot watch the contents from channels that I use, because you're not living in the United States, you can follow the rest of instructions to set up your own channel. I personally use this for football, but it can be used for anything, as you can customize list of youtube channels, mandatory keywords and optional ones.

# Big Note
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
