package main

import (
	"fmt"
	"log"

	"github.com/jsuar/go-cron-descriptor/pkg/crondescriptor"
)

func main() {
	cronExpression := "*/5 15 * * 0-5"
	cd, err := crondescriptor.NewCronDescriptor(cronExpression)
	if err != nil {
		cd.Logger.Panic(err.Error())
	}
	fullDescription, err := cd.GetDescription(crondescriptor.Full)
	if err != nil {
		cd.Logger.Panic(err.Error())
	}
	fmt.Printf("%s => %s\n", cronExpression, *fullDescription)

	cd.Options.Verbose = false
	cd.Options.DayOfWeekIndexZero = false
	cronExpression = "*/5 15 * * 1-6"
	if err = cd.Parse(cronExpression); err != nil {
		cd.Logger.Panic(err.Error())
	}

	fullDescription, err = cd.GetDescription(crondescriptor.Full)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s => %s\n", cronExpression, *fullDescription)

}
