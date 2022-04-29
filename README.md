# stryker-cli-utils

The purpose of this project is providing a set of utilities for improving the workflow for developers that want to use [Stryker.NET](https://stryker-mutator.io/docs/stryker-net/introduction) in their projects. You can find Stryker .NET repository [here](https://github.com/stryker-mutator/stryker-net).

## Features

- Running Stryker for all projects referenced by the application's main test project and condensing the report

```stryker-cli-utils run --a```

- Specifying a location for the report (must finish with ```.html```, for now). Default is ```mutation-report.html``` in the current folder

```stryker-cli-utils run --a --reportLocation="my-report.html"```

```stryker-cli-utils run --a --reportLocation="../my-report.html"```

```stryker-cli-utils run --a --reportLocation="C:/dev/my-report.html"```

- Running Stryker for all config files in the root project

```stryker-cli-utils run --i```

## Pre-requisites

- Stryker .NET installed
- Go installed
- For the ```--i``` flag, the config files must follow the ```-stryker-config.json``` format and be in the root of your application

## Roadmap 

- Write tests for everything
- Make the ```--a``` flag consider multiple ```Test .csproj``` files
- Make possible to choose the report type used
- Output progress of each Stryker execution
- Generate stryker-config-json from templates