# go-cron-descriptor

[![GoDoc](https://img.shields.io/badge/GoDoc-reference-007d9c?style=flat-square)](https://pkg.go.dev/github.com/jsuar/go-cron-descriptor/pkg/envconfig) [![Go Report Card](https://goreportcard.com/badge/github.com/jsuar/go-cron-descriptor)](https://goreportcard.com/report/github.com/jsuar/go-cron-descriptor)

Go-cron-descriptor translates [cron expressions](https://en.wikipedia.org/wiki/Cron) to English for quicker or easier interpretation of the expression. Only English is supported at the moment (see [Contributing](#contributing) if you would like to help localize to other languages). The ability to convert cron expressions to a human readable format allows for a better user experience. For example, in another project of mine, [the nomad-custodian CLI](https://github.com/jsuar/nomad-custodian#listing-batch-type-jobs) lists batch jobs which includes the associated cron expression and human readable format.

Translated to Go from [cron-expression-descriptor (C#)](https://github.com/bradymholt/cron-expression-descriptor) via [Cron Descriptor (Python)](https://github.com/Salamek/cron-descriptor). Original Author & Credit: Brady Holt (http://www.geekytidbits.com).

**Note:** I did not write the logic for this package. I've only made minor modifications to better conform to the Go language. If you see incorrect results or have any other recommended changes, please let me know by filing an issue.

#### Quick Examples
```
* * * * *         =>    Every minute
0 0 * * FRI       =>    At 00:00 AM, only on Friday
0 1 12 */7 *      =>    At 01:00 AM, on day 12 of the month, every 7 months
0 0 12 LW * ?     =>    At 12:00 PM, on the last weekday of the month
```

## Features

* [x] Supports all cron expression special characters including: `*` `/` `,` `-` `?` `L` `W` `#`
* [x] Supports 5, 6 (w/ seconds or year), or 7 (w/ seconds and year) part cron expressions
* [x] Provides casing options (sentence, title, lower)
* [x]  Support for non-standard non-zero-based week day numbers
* [ ] Supports printing to locale specific human readable format
* [ ] Supports displaying times in specific timezones

## Installation

```
go get https://github.com/jsuar/go-cron-descriptor
```

## Usage

### Options

#### `DayOfWeekIndexZero`

Toggling this option allows for Sunday to be either 0 or 1 in the day of the week field.

```
*/5 15 * * 0-6 => Every 5 minutes, at 03:00 PM, Sunday through Saturday
*/5 15 * * 1-7 => Every 5 minutes, at 03:00 PM, Sunday through Saturday
```

### Debug Statements

Since I was unfamiliar with the code and logic, I added many debug statements throughout the codebase for troubleshooting. Set the below environment variable to see debug statements.

```
export CRON_DESCRIPTOR_LOG_LEVEL=debug
```

### Examples

#### Zero based day of week

    package main

    import (
        "fmt"

        "github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
    )

    func main() {
        cronExpression := "*/5 15 * * 1-5"
        cd, _ := crondescriptor.NewCronDescriptor(cronExpression)
        
        fullDescription, _ := cd.GetDescription(crondescriptor.Full)
        fmt.Printf("%s => %s\n", cronExpression, *fullDescription)

        cronExpression = "0 0/30 8-9 5,20 * ?"
        cd.Parse(cronExpression)

        fullDescription, _ = cd.GetDescription(crondescriptor.Full)
        fmt.Printf("%s => %s\n", cronExpression, *fullDescription)
    }

Output:

    */5 15 * * 1-5 => Every 5 minutes, at 03:00 PM, Monday through Friday
    0 0/30 8-9 5,20 * ? => Every 30 minutes, 08:00 AM through 09:59 AM, on day 5 and 20 of the month

#### Non-zero based day of week

Setting `DayOfWeekIndexZero` to `false` will treat Sunday as the first day of the week instead of Monday.

    cronExpression := "*/5 15 * * 1-5"
	options := crondescriptor.Options{DayOfWeekIndexZero: false}
	cd, _ := crondescriptor.NewCronDescriptorWithOptions(cronExpression, options)

	fullDescription, _ := cd.GetDescription(crondescriptor.Full)
	fmt.Printf("%s => %s\n", cronExpression, *fullDescription)

Output:

    */5 15 * * 1-5 => Every 5 minutes, at 03:00 PM, Sunday through Thursday

### Contributing

1. Fork it
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create new Pull Request

## References

* Testing
  * https://crontab.guru
  * https://crontab.cronhub.io
  * https://www.freeformatter.com/cron-expression-generator-quartz.html
* Repositories
  * https://github.com/Salamek/cron-descriptor
  * https://github.com/golang-standards/project-layout/tree/master/pkg
  * https://github.com/sethvargo/go-envconfig
* Other
  * https://dave.cheney.net/practical-go/presentations/qcon-china.html