# route-break

Find where the break in the trace route is.


I have a very weird problem where occassionially my internet will just cut out for a few seconds.  I don't know where in the path its failing, and it doesn't happen that often, so I built this thing to track it and check it.




## How

- run a traceroute to 8.8.8.8 (Google DNS) and record all the IPs in between
- ping each IP separately in a go routine.  The ping is limited by a count.
- if pings start failing, start logging and yelling about it


## To Run

On Linux, need to use sudo as the ping and traceroute library needs permissions (can read more [here](github.com/go-ping/ping)).

``` 
sudo go run .
```