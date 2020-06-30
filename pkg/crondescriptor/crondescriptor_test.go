package crondescriptor

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	type args struct {
		expression string
	}
	tests := []struct {
		name           string
		args           args
		wantParsedExpr [7]string
		wantErr        bool
	}{
		{
			"Case 0",
			args{"5 4 * * *"},
			[7]string{"", "5", "4", "*", "*", "*", ""},
			false,
		},
		{
			"With - and /",
			args{"23 0-20/2 * * *"}, // At minute 23 past every 2nd hour from 0 through 20.
			[7]string{"", "23", "0-20/2", "*", "*", "*", ""},
			false,
		},
		{
			"At 04:05 on Sunday",
			args{"5 4 * * sun"}, // At 04:05 on Sunday.
			[7]string{"", "5", "4", "*", "*", "0", ""},
			false,
		},
		{
			"With * , /",
			args{"0 0,12 1 */2 *"}, // At minute 0 past hour 0 and 12 on day-of-month 1 in every 2nd month.
			[7]string{"", "0", "0,12", "1", "*/2", "*", ""},
			false,
		},
	}
	for _, tt := range tests {
		cd, err := NewCronDescriptor(tt.args.expression)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Run(tt.name, func(t *testing.T) {
			err := cd.Parse(tt.args.expression)
			if err != nil {
				t.Errorf(err.Error())
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(cd.expressionArray, tt.wantParsedExpr) {
				t.Errorf("Parse() = %v, want %v", cd.expressionArray, tt.wantParsedExpr)
			}
		})
	}
}

func Test_decreaseDaysOfWeek(t *testing.T) {
	type fields struct {
		in0 string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			"Decrement from 1 to 0",
			fields{"1"},
			"0",
			false,
		},
		{
			"Throw error while decrementing from 0 to -1",
			fields{"0"},
			"0",
			true,
		},
	}
	cd, err := NewCronDescriptor("* * * ? * *")
	if err != nil {
		t.Errorf(err.Error())
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := cd.decreaseDaysOfWeek(tt.fields.in0)
			// fmt.Printf("%t - %s \n", err != nil, err)
			if (err != nil) != tt.wantErr {
				t.Errorf("decreaseDaysOfWeek() error = %v, wantErr %v\n", err, tt.wantErr)
				return
			} else if (err != nil) && tt.wantErr {
				fmt.Printf("Error correctly thrown: %s\n", err)
			} else if *got != tt.want {
				t.Errorf("decreaseDaysOfWeek() = %v, want %v\n", *got, tt.want)
			}
		})
	}
}

func Test_getYearDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"*",
			fields{"* * * ? * * *"},
			"",
			false,
		},
		{
			"1970-2099",
			fields{"* * * ? * * 1970-2099"},
			// "between 1970 and 2099",
			", 1970 through 2099",
			false,
		},
		{
			"*/6",
			fields{"* * * ? * * */6"},
			", every 6 years",
			false,
		},
		{
			"2020/2",
			fields{"* * * ? * * 2020/2"},
			", every 2 years, starting in 2020",
			false,
		},
		{
			"2020-2025/2",
			fields{"* * * ? * * 2020-2025/2"},
			", every 2 years, 2020 through 2025",
			false,
		},
		{
			"2020,2021",
			fields{"* * * ? * * 2020,2021"},
			", in 2020 and 2021",
			false,
		},
		{
			"2020,2021,2022",
			fields{"* * * ? * * 2020,2021,2022"},
			", in 2020, 2021, and 2022",
			false,
		},
		{
			"2020,2021,2022,2000",
			fields{"* * * ? * * 2020,2021,2022,2000"},
			", in 2020, 2021, 2022, and 2000",
			false,
		},
	}
	for _, tt := range tests {
		cd, err := NewCronDescriptor(tt.fields.expression)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := cd.getYearDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("getYearDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("getYearDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}

func Test_getDayOfMonthDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"*",
			fields{"* * * * *"},
			", every day",
			false,
		},
		{
			"*/7",
			fields{"* * */7 * *"},
			", every 7 days",
			false,
		},
		{
			"1",
			fields{"* * 1 * *"},
			", on day 1 of the month",
			false,
		},
		{
			"L",
			fields{"* * L * *"},
			", on the last day of the month",
			false,
		},
		{
			"L-1",
			fields{"* * L-1 * *"},
			", between day L and 1 of the month",
			false,
		},
		{
			"15W",
			fields{"* * 15W * *"},
			", on the weekday nearest day 15 of the month",
			false,
		},
	}

	for _, tt := range tests {
		cd, err := NewCronDescriptor(tt.fields.expression)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := cd.getDayOfMonthDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("getDayOfMonthDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("getDayOfMonthDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}

func Test_getMonthDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"*",
			fields{"* * * * *"},
			"",
			false,
		},
		{
			"4",
			fields{"* * * 4 *"},
			", only in April",
			false,
		},
		{
			"DEC",
			fields{"* * * DEC *"},
			", only in December",
			false,
		},
		{
			"JAN,JUN",
			fields{"* * * JAN,JUN *"},
			", only in January and June",
			false,
		},
		{
			"JAN,FEB,MAR,APR",
			fields{"* * * JAN,FEB,MAR,APR *"},
			", only in January, February, March, and April",
			false,
		},
		{
			"9-12",
			fields{"* * * 9-12 *"},
			", September through December",
			false,
		},
		{
			"APR-JUN",
			fields{"* * * APR-JUN *"},
			", April through June",
			false,
		},
	}
	for _, tt := range tests {
		cd, err := NewCronDescriptor(tt.fields.expression)
		if err != nil {
			t.Errorf(err.Error())
		}
		t.Run(tt.name, func(t *testing.T) {
			got, err := cd.getMonthDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("getMonthDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("getMonthDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}

func Test_getDayOfTheWeekDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"*",
			fields{"* * * * *"},
			", every day",
			false,
		},
		{
			"1,4,6",
			fields{"* * * * 1,4,6"},
			", only on Monday, Thursday, and Saturday",
			false,
		},
		{
			"5-0",
			fields{"* * * * 5-0"},
			", Friday through Sunday",
			false,
		},
		{
			"5#3",
			fields{"* * * * 5#3"},
			", on the third Friday of the month",
			false,
		},
		{
			"2#2",
			fields{"* * * * 2#2"},
			", on the second Tuesday of the month",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd, err := NewCronDescriptor(tt.fields.expression)
			if err != nil {
				t.Errorf(err.Error())
			}
			got, err := cd.getDayOfTheWeekDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("CronExpression.getDayOfTheWeekDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if *got != tt.want {
				t.Errorf("CronExpression.getDayOfTheWeekDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}

func TestCronExpression_getHoursDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"*",
			fields{"* * * * *"},
			"every hour",
			false,
		},
		{
			"8-12",
			fields{"* 8-12 * * *"},
			"08:00 AM through 12:59 PM",
			false,
		},
		{
			"8,12,3,21",
			fields{"* 8,12,3,21 * * *"},
			"at 08:00 AM, 12:00 PM, 03:00 AM, and 09:00 PM",
			false,
		},
		{
			"8:30-12:45",
			fields{"* 8:30-12:45 * * *"},
			", 8 through 12",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd, err := NewCronDescriptor(tt.fields.expression)
			if err != nil {
				t.Errorf(err.Error())
			}
			got, err := cd.getHoursDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("CronExpression.getHoursDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) && tt.wantErr {
				fmt.Printf("Error correctly thrown: %s\n", err)
			} else if *got != tt.want {
				t.Errorf("CronExpression.getHoursDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}

func TestCronExpression_getMinutesDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"*",
			fields{"* * * * *"},
			"every minute",
			false,
		},
		{
			"0",
			fields{"0 * * * *"},
			"at 0 minutes past the hour",
			false,
		},
		{
			"5-40",
			fields{"5-40 * * * *"},
			", minutes 5 through 40 past the hour",
			false,
		},
		{
			"5,10,15,30",
			fields{"5,10,15,30 * * * *"},
			"at 5, 10, 15, and 30 minutes past the hour",
			false,
		},
		{
			"Every even minute",
			fields{"*/2 * * * *"},
			"every 2 minutes",
			false,
		},
		{
			"Every uneven minute",
			fields{"1/2 * * * *"},
			"every 2 minutes, starting at 1",
			false,
		},
		{
			"Every 10 minutes",
			fields{"*/10 * * * *"},
			"every 10 minutes",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd, err := NewCronDescriptor(tt.fields.expression)
			if err != nil {
				t.Errorf(err.Error())
			}
			got, err := cd.getMinutesDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("CronExpression.getMinutesDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) && tt.wantErr {
				fmt.Printf("Error correctly thrown: %s\n", err)
			} else if *got != tt.want {
				t.Errorf("CronExpression.getMinutesDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}

func TestCronExpression_getSecondsDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"*",
			fields{"* * * * * *"},
			"every second",
			false,
		},
		{
			"10/4",
			fields{"10/4 * * * * *"},
			"every 4 seconds, starting at 10",
			false,
		},
		{
			"5-50",
			fields{"5-50 * * * * *"},
			", seconds 5 through 50 past the minute",
			false,
		},
		{
			"0,4,8,16",
			fields{"0,4,8,16 * * * * *"},
			"at 0, 4, 8, and 16 seconds past the minute",
			false,
		},
		{
			"60",
			fields{"60 * * * * *"},
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd, err := NewCronDescriptor(tt.fields.expression)
			if err != nil {
				t.Errorf(err.Error())
			}
			got, err := cd.getSecondsDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("CronExpression.getSecondsDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) && tt.wantErr {
				fmt.Printf("Error correctly thrown: %s\n", err)
			} else if *got != tt.want {
				t.Errorf("CronExpression.getSecondsDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}

func TestCronExpression_getTimeOfDayDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			"*",
			fields{"* * * * * *"},
			"every second, every minute, every hour",
			false,
		},
		{
			"0 * *",
			fields{"0 * * * * *"},
			"every minute, every hour",
			false,
		},
		{
			"0 0 *",
			fields{"0 0 * * * * *"},
			"at 0 minutes past the hour, every hour",
			false,
		},
		{
			"* */15 *",
			fields{"0 */15 * * * * *"},
			"every 15 minutes, every hour",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cd, err := NewCronDescriptor(tt.fields.expression)
			if err != nil {
				t.Errorf(err.Error())
			}
			got, err := cd.getTimeOfDayDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("CronExpression.getTimeOfDayDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) && tt.wantErr {
				fmt.Printf("Error correctly thrown: %s\n", err)
			} else if *got != tt.want {
				t.Errorf("CronExpression.getTimeOfDayDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}

func TestCronExpression_getFullDescription(t *testing.T) {
	type fields struct {
		expression string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{"", fields{"* * * * *"}, "Every minute", false},
		{"", fields{"0 2 * * *"}, "At 02:00 AM", false},
		{"", fields{"0 5,17 * * *"}, "At 05:00 AM and 05:00 PM", false},
		{"", fields{"0 17 * * sun"}, "At 05:00 PM, only on Sunday", false},
		{"", fields{"*/10 * * * *"}, "Every 10 minutes", false},
		{"", fields{"* * * jan,may,aug *"}, "Every minute, only in January, May, and August", false},
		{"", fields{"0 17 * * sun,fri"}, "At 05:00 PM, only on Sunday and Friday", false},
		{"", fields{"0 2 * * sun"}, "At 02:00 AM, only on Sunday", false},
		{"", fields{"0 */4 * * *"}, "Every 4 hours", false},
		{"", fields{"0 4,17 * * sun,mon"}, "At 04:00 AM and 05:00 PM, only on Sunday and Monday", false},
		// {"", fields{"* * * ? * *"}, "Every second", false},
		// {"", fields{"0 * * ? * *"}, "Every minute", false},
		// {"", fields{"0 */2 * ? * *"}, "Every even minute", false},
		// {"", fields{"0 1/2 * ? * *"}, "Every uneven minute", false},
		// {"", fields{"0 */2 * ? * *"}, "Every 2 minutes", false},
		// {"", fields{"0 */3 * ? * *"}, "Every 3 minutes", false},
		// {"", fields{"0 */4 * ? * *"}, "Every 4 minutes", false},
		// {"", fields{"0 */5 * ? * *"}, "Every 5 minutes", false},
		// {"", fields{"0 */10 * ? * *"}, "Every 10 minutes", false},
		// {"", fields{"0 */15 * ? * *"}, "Every 15 minutes", false},
		// {"", fields{"0 */30 * ? * *"}, "Every 30 minutes", false},
		// {"", fields{"0 15,30,45 * ? * *"}, "Every hour at minutes 15, 30 and 45", false},
		// {"", fields{"0 0 * ? * *"}, "Every hour", false},
		// {"", fields{"0 0 */2 ? * *"}, "Every hour", false},
		// {"", fields{"0 0 0/2 ? * *"}, "Every even hour", false},
		// {"", fields{"0 0 1/2 ? * *"}, "Every uneven hour", false},
		// {"", fields{"0 0 */3 ? * *"}, "Every three hours", false},
		// {"", fields{"0 0 */4 ? * *"}, "Every four hours", false},
		// {"", fields{"0 0 */6 ? * *"}, "Every six hours", false},
		// {"", fields{"0 0 */8 ? * *"}, "Every eight hours", false},
		// {"", fields{"0 0 */12 ? * *"}, "Every twelve hours", false},
		// {"", fields{"0 0 0 * * ?"}, "Every day at midnight - 12am", false},
		// {"", fields{"0 0 1 * * ?"}, "Every day at 1am", false},
		// {"", fields{"0 0 6 * * ?"}, "Every day at 6am", false},
		// {"", fields{"0 0 12 * * ?"}, "Every day at noon - 12pm", false},
		// {"", fields{"0 0 12 * * ?"}, "Every day at noon - 12pm", false},
		// {"", fields{"0 0 12 * * SUN"}, "Every Sunday at noon", false},
		// {"", fields{"0 0 12 * * MON"}, "Every Monday at noon", false},
		// {"", fields{"0 0 12 * * TUE"}, "Every Tuesday at noon", false},
		// {"", fields{"0 0 12 * * WED"}, "Every Wednesday at noon", false},
		// {"", fields{"0 0 12 * * THU"}, "Every Thursday at noon", false},
		// {"", fields{"0 0 12 * * FRI"}, "Every Friday at noon", false},
		// {"", fields{"0 0 12 * * SAT"}, "Every Saturday at noon", false},
		// {"", fields{"0 0 12 * * MON-FRI"}, "Every Weekday at noon", false},
		// {"", fields{"0 0 12 * * SUN,SAT"}, "Every Saturday and Sunday at noon", false},
		// {"", fields{"0 0 12 */7 * ?"}, "Every 7 days at noon", false},
		// {"", fields{"0 0 12 1 * ?"}, "Every month on the 1st, at noon", false},
		// {"", fields{"0 0 12 2 * ?"}, "Every month on the 2nd, at noon", false},
		// {"", fields{"0 0 12 15 * ?"}, "Every month on the 15th, at noon", false},
		// {"", fields{"0 0 12 1/2 * ?"}, "Every 2 days starting on the 1st of the month, at noon", false},
		// {"", fields{"0 0 12 1/4 * ?"}, "Every 4 days staring on the 1st of the month, at noon", false},
		// {"", fields{"0 0 12 L * ?"}, "Every month on the last day of the month, at noon", false},
		// {"", fields{"0 0 12 L-2 * ?"}, "Every month on the second to last day of the month, at noon", false},
		// {"", fields{"0 0 12 LW * ?"}, "Every month on the last weekday, at noon", false},
		// {"", fields{"0 0 12 1L * ?"}, "Every month on the last Sunday, at noon", false},
		// {"", fields{"0 0 12 2L * ?"}, "Every month on the last Monday, at noon", false},
		// {"", fields{"0 0 12 6L * ?"}, "Every month on the last Friday, at noon", false},
		// {"", fields{"0 0 12 1W * ?"}, "Every month on the nearest Weekday to the 1st of the month, at noon", false},
		// {"", fields{"0 0 12 15W * ?"}, "Every month on the nearest Weekday to the 15th of the month, at noon", false},
		// {"", fields{"0 0 12 ? * 2#1"}, "Every month on the first Monday of the Month, at noon", false},
		// {"", fields{"0 0 12 ? * 6#1"}, "Every month on the first Friday of the Month, at noon", false},
		// {"", fields{"0 0 12 ? * 2#2"}, "Every month on the second Monday of the Month, at noon", false},
		// {"", fields{"0 0 12 ? * 5#3"}, "Every month on the third Thursday of the Month, at noon - 12pm", false},
		// {"", fields{"0 0 12 ? JAN *"}, "Every day at noon in January only", false},
		// {"", fields{"0 0 12 ? JUN *"}, "Every day at noon in June only", false},
		// {"", fields{"0 0 12 ? JAN,JUN *"}, "Every day at noon in January and June", false},
		// {"", fields{"0 0 12 ? DEC *"}, "Every day at noon in December only", false},
		// {"", fields{"0 0 12 ? JAN,FEB,MAR,APR *"}, "Every day at noon in January, February, March and April", false},
		// {"", fields{"0 0 12 ? 9-12 *"}, "Every day at noon between September and December", false},
	}
	for _, tt := range tests {
		t.Run(tt.fields.expression, func(t *testing.T) {
			cd, err := NewCronDescriptor(tt.fields.expression)
			if err != nil {
				t.Errorf(err.Error())
			}
			got, err := cd.getFullDescription()
			if (err != nil) != tt.wantErr {
				t.Errorf("CronExpression.getFullDescription() error = %v, wantErr %v", err, tt.wantErr)
				return
			} else if (err != nil) && tt.wantErr {
				fmt.Printf("Error correctly thrown: %s\n", err)
			} else if *got != tt.want {
				t.Errorf("CronExpression.getFullDescription() = %v (str len=%d)", *got, len(*got))
				t.Errorf("want: %v (str len=%d)", tt.want, len(tt.want))
			}
		})
	}
}
