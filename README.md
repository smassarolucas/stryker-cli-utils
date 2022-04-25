# stryker-cli-utils

The purpose of this project is providing a set of utilities for improving the workflow for developers that want to use [Stryker.NET](https://stryker-mutator.io/docs/stryker-net/introduction) in their projects. You can find Stryker .NET repository [here](https://github.com/stryker-mutator/stryker-net).

## Features

- Running all mutation tests for the solution with one command and condensing the report
```stryker-cli-utils mutate```
- Specifying a location for the report (must finish with *.html*, for now). Default is *mutation-report.html* in the current folder

```stryker-cli-utils mutate --reportLocation="my-report.html"```
```stryker-cli-utils mutate --reportLocation="../my-report.html"```

```stryker-cli-utils mutate --reportLocation="C:/dev/my-report.html"```

## Roadmap 

- Generate stryker-config-json from templates