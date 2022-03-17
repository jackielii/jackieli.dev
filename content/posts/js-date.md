+++ 
date = 2020-07-25T23:29:55+01:00
title = "Javascript Date"
description = "Intricacies of Javascript Date"
slug = "javascript-date" 
tags = ["javascript", "date", "datetime"]
categories = ["javascript"]
externalLink = ""
series = []
+++

So I need a way to represent LocalDate that can be found in `java.time`
package. Turns out it's harder than it looks.

Take `LocalDateTime` for example: `LocalDateTime` is a relative date that's the
same no matter which time zone you're in. The easiest thing is to just use
string: `2020-07-25T01:02:03`. Note that there is no zone/offset information in
the string

But if you want to use a `Date` object in javascript to represent it, it gets
very tricky.

My first attempt is to just subtract the users' zone offset like so:

```ts
const offsetInMillis = new Date().getTimezoneOffset() * 60000 // minutes to ms

function toLocalDatetime(date: Date|number): Date {
  let relativeTimestamp;
  if (date instanceof Date) {
	relativeTimestamp = date.getTime() - offsetInMillis
  } else {
    relativeTimestamp = date - offsetInMillis
  }
  return new Date(relativeTimestamp)
}
```

I thought naturally if you create a date from string without timezone, it would create it in the user's local timezone. But actually it's not, see [MDN](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/parse).

```
❯ node
Welcome to Node.js v17.5.0.
Type ".help" for more information.
> Date.parse('01 Jan 1970 00:00:00 GMT')
0
> Date.parse('01 Jan 1970 00:00:00')
-3600000
> new Date().getTimezoneOffset()
0
```



### Conclusion

Bottom line is using just `Date` is not enough to represent LocalDateTime. If
you're in the same situation, just use
[js-jodatime](https://js-joda.github.io/js-joda/)
