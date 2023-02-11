# GoSysLat
[![Downloads][1]][2] [![GitHub stars][3]][4]

[1]: https://img.shields.io/github/downloads/spddl/GoSysLat/total.svg
[2]: https://github.com/spddl/GoSysLat/releases "Downloads"

[3]: https://img.shields.io/github/stars/spddl/GoSysLat.svg
[4]: https://github.com/spddl/GoSysLat/stargazers "GitHub stars"

I didn't understand the original SysLat software and how it interacted with the device at first, so I wanted to rebuild it to understand it better. This allowed me to add some functions afterwards and increase the speed of the query.

To use this GoSysLat client you need the device [SysLat](https://syslat.com) and a customized [firmware](https://github.com/spddl/SysLat_Firmware)

[Here are some of my test results.](https://bit.ly/gosyslat) These results are of course only momentary recordings from my system.
It should be clear that the latency measurement is only a polling of different hardware and software. That means you will never get the same results if you compare 2 settings that do not change the latency.

In this repo are also the TestCase's I used if someone wants to recreate this.

Thanks for the support and help to [timecard](https://github.com/djdallmann), [Skew](https://github.com/Skewjo) and [henmill](https://github.com/henmill)

![example](/example.png)