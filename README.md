# certcal

Provide an iCal web feed for certificate expiry.


## Usage

There are a couple of ways to run this

### Standalone
You can download one of the binaries and run it as a standalone:

```shell
% PORT=3000 CERTCAL_INTERVAL=5h CERTCAL_HOSTS="unwiredcouch.com" ./certcal
```


### As part of an existing http mux

You can include the handler in your existing mux, something along the lines
of:

```go
import (
  "github.com/mrtazz/certcal/handler"
  "github.com/mrtazz/certcal/hosts"
  "net/http"
  "time"
)

func Run() {

  ...

  hosts.AddHosts([]string{"unwiredcouch.com"})
  hosts.UpdateEvery(5 * time.Hour)

  http.HandleFunc("/hosts", handler.Handler)
  http.ListenAndServe(":3000", nil)

}
```


### Via Docker
There is a docker image as well that you can use


## FAQ

### Shouldn't certs renew automatically?
Probably. But sometimes they aren't.

### Shouldn't this be an alert somewhere?
Maybe, up to you.

### Old expiry events are disappearing!
That is by design. The UID of the `VEVENT` is the `sha256sum` of the summary
of the event. Because generally if the cert got renewed the old event is just
cruft.


## Inspiration
[genuinetools/certok](https://github.com/genuinetools/certok) inspired this
