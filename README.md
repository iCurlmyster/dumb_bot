# dumb_bot
Twitch bot written in Golang that supports plugins

Currently there is only an initial set up.

A simple twitch bot that connects via websockets using the gorilla websocket library

Program loads in a `config.json` and `plugins.json` file at start, located in the `resources` folder.

`config.json` is where your OAuth token, Username, and the Channel name that you want to connect the bot to.

`plugins.json` is where you list plugins that should be loaded into the bot. The json file is key value pairs, where the Key is the command that should fire the plugin, and the Value is the path to the .so plugin file.

#### example

We have a plugin print.go which prints out what the user typed in the command.

example command in twitch chat: `!print Hello, World!`

example output in twitch chat: `your message: Hello, World!`

how to setup up in plugins.json:

```javascript
{
  "print": "./plugins/print/print.so"
}
```

If you add plugins you should add them to the Makefile so that when you run make all plugins can be built without having to do them indivdually.

### Requires:
- Go 1.10+
- GNU Make 4.1+ (This is optional if you want to compile the plugins individually)

## TODOs
- Add more functional plugins to the bot
- Clean up where needed
