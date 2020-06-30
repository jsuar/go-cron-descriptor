# go-cron-descriptor

[![GoDoc](https://img.shields.io/badge/GoDoc-reference-007d9c?style=flat-square)](https://pkg.go.dev/github.com/jsuar/go-cron-descriptor/pkg/envconfig)

A Go library that converts cron expressions into human readable strings. Translated to Go from [cron-expression-descriptor (C#)](https://github.com/bradymholt/cron-expression-descriptor) via [Cron Descriptor (Python)](https://github.com/Salamek/cron-descriptor).

Original Author & Credit: Brady Holt (http://www.geekytidbits.com).

**Note:** I did not write the logic for this package. I've only made minor modifications to better conform to the Go language. If you see incorrect results or have any other recommended changes, please let me know by filing an issue.

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

### Debug Statements

Since I was unfamiliar with the code and logic, I added many debug statements throughout the codebase for troubleshooting. Set the below environment variable to see debug statements.

```
export CRON_DESCRIPTOR_LOG_LEVEL=debug
```

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