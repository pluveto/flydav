# Go CLI Boilerplate

A powerful Golang CLI application scaffold integrated with logrus, go-arg, config, testify and Github Action.

## Usage

### 1. Use as template

Click `Use this template` button to create a new repository. Then clone the new repository to your local machine.

### 2. Do some changes

1. Change the package name
    + Replace `example.com/m` with your own package name.
    + The default is `example.com/m` You can replace by `sed`:

    ```bash
    sed -i 's/example.com\/m/your.package.name/g' $(find . -type f)
    ```

    Or use global search and replace in your IDE.

2. Change the module name
    + The default application name is `greet`
    + Rename `cmd/greet` to `cmd/YOUR_APP_NAME`
    + Rename inside `Makefile` for `BIN_NAMES`
    + Rename inside `conf/*.toml` for logger file name
    + Rename insdie `.vscode` for debugging

3. Add more application if needed
    + Add more application in `cmd` folder just like `greet` app.
    + Add more application name in `Makefile` for `BIN_NAMES`, for example:

    ```makefile
    BIN_NAMES=app1 app2
    ```

4. Enable write permission for workflow
    + Click `Settings` tab of your repository.
    + Select `Action -> General` in left sidebar.
    + Locate `Workflow permissions` section.
    + Check `Read and write permission`

### 3. Release

1. Edit `README.md`. Switch `LICENSE` to your own license.
2. Coding for your application.
3. Run `make dev` to build for your local machine.
4. Run `make test` to run unit tests.
5. Run `make package` to cross compile for different platforms.
6. Click `Action` tab of your repository to see the Github Action workflow. In release workflow, click `Run workflow` button to release your application.

## Run Example

```bash
go run cmd/greet/main.go Jack -v -c ./dist/config.example.toml
```

## Features

+ [x] Logrus logger. JSON format
+ [x] Multiple logging output: file, stdout, stderr
+ [x] Command line args, such as `--verbose` and `--config`
+ [x] Config file, such as `./dist/config.example.toml`
+ [x] Log rotation

Use `| jq` to pretty print JSON log.

Generate structs for your config file: [Toml to Go](https://xuri.me/toml-to-go/)
