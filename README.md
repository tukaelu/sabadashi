# sabadashi
![Latest GitHub Release](https://img.shields.io/github/release/tukaelu/sabadashi.svg)
[![Go Report Card](https://goreportcard.com/badge/tukaelu/sabadashi)](https://goreportcard.com/report/tukaelu/sabadashi)
![Github Actions Test](https://github.com/tukaelu/sabadashi/actions/workflows/ci.yaml/badge.svg?branch=main)

![](./images/sabadashi-logo.png)

日本語のREADMEは[こちら](README-ja.md)です。

## Description

**sabadashi** is a CLI tool for retrieves metrics posted to Mackerel for a specified period of time.  
Supports retrieving host metrics and service metrics.

I note that this section is very important, but this tool makes many API calls.  
Therefore, please refrain from executing them concurrently.

## Install

### Homebrew tap

```
brew install tukaelu/tap/sabadashi
```

### Binary install

Please download the appropriate Zip archive for your environment from the [releases](https://github.com/tukaelu/sabadashi/releases).

## Usage

### Host metrics

```
NAME:
   sabadashi host - Retrieves host metrics

USAGE:
   sabadashi host [command options] [arguments...]

OPTIONS:
   --id value                       ID of the host from which to retrieve metric
   --from value                     Specify the date to start retrieving metrics in YYYYMMDD format. (e.g. 20230101)
   --to value                       Specify the date to end retrieving metrics in YYYYMMDD format. (e.g. 20231231)
   --granularity value, -g value    Specify the granularity of metric data. Choose from 1m, 5m, 10m, 1h, 2h, 4h or 1d. (default: 1m)
   --with-friendly-date-format, -f  If this flag is enabled, an additional column with a friendly date format is output at the beginning of the CSV line. (default: false)
   --help, -h                       show help
```

### Service metrics
```
NAME:
   sabadashi service - Retrieves service metrics

USAGE:
   sabadashi service [command options] [arguments...]

OPTIONS:
   --name value, -n value           Name of the service from which to retrieve metric
   --from value                     Specify the date to start retrieving metrics in YYYYMMDD format. (e.g. 20230101)
   --to value                       Specify the date to end retrieving metrics in YYYYMMDD format. (e.g. 20231231)
   --granularity value, -g value    Specify the granularity of metric data. Choose from 1m, 5m, 10m, 1h, 2h, 4h or 1d. (default: 1m)
   --with-friendly-date-format, -f  If this flag is enabled, an additional column with a friendly date format is output at the beginning of the CSV line. (default: false)
   --with-external-monitors, -e     If this flag is enabled, it also includes the metric measured in the external monitoring. (default: false)
   --help, -h                       show help
```

Specify the ID of the host from which the metric is to be retrieved, and the start and end dates of the retrieval in the format YYYYYMMDD.  
The tool will retrieve metrics posted from `YYYYY/MM/DD 00:00:00` specified in `from` to `YYYYY/MM/DD 23:59:59` specified in `to`. And the validity period is 460 days.

```
# If the MACKEREL_APIKEY is set in an environment variable
sabadashi host -id <your host id> -from <YYYYMMDD> -to <YYYYMMDD>

# If not, and you explicitly specify
sabadashi host -apkey <your api key> -id <your host id> -from <YYYYMMDD> -to <YYYYMMDD>
```

The API key specified in the environment variable `MACKEREL_APIKEY` or the `-apikey` option must have read permission.

When the command is executed, a directory named by host ID and start/end date will be created under the working directory, and a CSV file for each metric will be output in the directory.

## Notes

- This plugin is unofficial. Please ask questions via Issue or SNS.
- As mentioned earlier, the concurrent act of retrieving metrics for multiple hosts should be avoided, as it puts a load on the service.
- Host metrics are retrieved based on the [List Metric Names API](https://mackerel.io/api-docs/entry/hosts#metric-names) of host metrics.
- The service metrics are retrieved based on the [List Metric Names API](https://mackerel.io/api-docs/entry/services#metric-names) of service metrics and, only if the option to target external monitoring is enabled, additional metrics measured through external monitoring are added from the [List Monitor Configurations API](https://mackerel.io/api-docs/entry/monitors#list).
- If a metric has not been submitted during the specified time period, the data for that time period will not be output as rows in the CSV and may in some cases result in an empty file.

## License

Copyright 2023 tukaelu (Tsukasa NISHIYAMA)

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

```
http://www.apache.org/licenses/LICENSE-2.0
```

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
