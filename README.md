# Flydav

A webdav application written in Go.

## Get started in 30 seconds

1. Start by downloading FlyDav from their website at [release page](https://github.com/pluveto/flydav/releases).
2. Run `./flydav -H 0.0.0.0` to start the server. Then you'll input the password for the default user `flydav`.
3. Open `http://YOUR_IP:7086/webdav` in your webdav client such as RaiDrive.

## Configuring FlyDav

1. Start by downloading FlyDav from their website at [release page](https://github.com/pluveto/flydav/releases).

2. Now that you have the software, you need to create a configuration file for it. Start by creating a new file called `flydav.toml`.

3. Inside the configuration file, you will need to add the following information:

- `[server]`: This section will define the host, port, and path of the webdav server.
  - `host`: The IP address of the host. This should be set to “0.0.0.0” if you want to make the server accessible from any IP address.
  - `port`: The port number to use for the webdav server.
  - `path`: The path of the webdav server.
- `fs_dir`: The directory on the server where the webdav files will be stored.
- `[auth]`: This section will define the authentication settings for the webdav server.
  - `[[auth.user]]`: This subsection will define the username and credentials for each user that has access to the webdav server.
    - `username`: The username of the user.
    - `sub_fs_dir`: The subdirectory of the fs_dir to which the user will have access.
    - `sub_path`: The path that the user will access the webdav server from.
    - `password_hash`: The hashed password of the user.
    - `password_crypt`: The type of hashing algorithm used to hash the password. This should be set to “bcrypt”.
- `[log]`: This section will define the logging settings for the webdav server.
  - `level`: The log level of the server.
  - `[[log.file]]`: This subsection will define the settings for the log file. Ignore this subsection if you do not want to log to a file.
    - `format`: The format of the log file.
    - `path`: The path of the log file.
    - `max_size`: The maximum size of the log file in megabytes.
    - `max_age`: The maximum age of the log file in days.
  - `[[log.stdout]]`: This subsection will define the settings for the log output to the console. Ignore this subsection if you do not want to log to the console.
    - `format`: The format of the log output.
    - `output`: The output stream for the log output.

4. Save the configuration file and run the FlyDav server. You should now be able to access the webdav server with the configured settings.

## Features

- [x] Basic authentication
- [x] Multiple users
- [x] Different root directory for each user
- [x] Different path prefix for each user
- [x] Logging
- [ ] SSL - You can use a reverse proxy like Nginx to enable SSL.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
