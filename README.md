# foo_cover_upload
A dead simple, memory based image server and uploader for use with foobar2000 Discord RPC.

## Usage
Download the binary suitable for your operating system in the latest release. A single binary houses the client and the server.
### Client
After you downloaded the binary, put it into somewhere persistent.

In your foobar preferences, head to Discord Rich Presence Integration -> Advanced and fill out the "artwork upload command" field.

Example:

```
"C:/path/to/directory/foo_cover_upload.exe" -url "https://image.your.domain/upload"
```

If the image server requires a password, pass it with the `-auth` flag.

### Server
The server uses Redis as intermediate storage, set up Redis first.

```bash
# Arch
pacman -Syu redis

# Ubuntu
apt install redis

# Fedora/RHEL
dnf install redis

# Enable service
systemctl enable --now redis
```
By default, redis launches on `127.0.0.1:6379`

You can test run the image server with the following
```
foo_cover_upload --serve --host "image.your.domain/upload" --redis-host "127.0.0.1:6379"
```
The image server should now be active on `127.0.0.1:2131`. Visiting the index page in a browser will show you a help page on how to configure your client.

## Help
```
===============================

foo_cover_upload HELP PAGE

===============================

This program houses both the client and the server.
By default, the program will use client mode, which is intended for uploading.
To change mode to server, you must pass the -serve flag to the program along with the necessary parameters.

===============================

AVAILABLE MODES:
-gen
        Generate a key to use for the image server.
-help
        Displays this help message.
-serve
        Launch a server.

===============================

AVAILABLE PARAMETERS:
-addr
        Listening address of the image upload server.
-auth
        Image server password, if required.
-hash
        Sets the MD5 hash of your password on the server.
-host
        The host of your server. Use your domain that can be used by Discord to access images.
-redis-host
        Host of the redis instance.
-redis-password
        Password of the redis instance.
-url
        Address of the server to upload to, including port.
```
## License
MIT