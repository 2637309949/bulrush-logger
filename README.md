# bulrush-logger
Provides logging capabilities.
- EXAMPLE:   
```go
// Logger Plugin init
var Logger = &logger.Logger{
	Path: "logs",
	Format: func(p *logger.Payload, ctx *gin.Context) string {
		if p.Type == logger.INT {
			startTime := time.Unix(p.StartUnix, 0).Format("2006/01/02 15:04:05")
			return fmt.Sprintf("[%v bulrush] => %s %6s %s", startTime, p.IP, p.Method, p.URL)
		} else if p.Type == logger.OUT {
			endOfTime := time.Unix(p.EndUnix, 0).Format("2006/01/02 15:04:05")
			latency := float64(time.Unix(p.EndUnix, 0).Sub(time.Unix(p.StartUnix, 0)) / time.Millisecond)
			return fmt.Sprintf("[%v bulrush] <= %.2fms %s %6s %s", endOfTime, latency, p.IP, p.Method, p.URL)
		}
		return "FROMAT ERROR"
	},
}
```

## MIT License

Copyright (c) 2018-2020 Double

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.