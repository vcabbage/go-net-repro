This repo reproduces a bug observed when two goroutines are sending to each other concurrently over UDP.

I have tried to reproduce on OSX (bare metal) and Ubuntu (Parallels VM), but have only observed the failure within Docker.

Easiest way to reproduce is executing `./run.sh`; provided you have Docker installed. This will build repro.go and execute it under strace repeatedly until it fails.

Note that `run.sh` uses the `--privileged` argument to allow strace to function.

#### Example Output

Failure:

Note the `EPERM (Operation not permitted)` from `sendto`. This error is not documented in `sendto(2)`.
Some searches suggest that it's related to iptables. Disabling iptables on the Docker host does not resolve the error.

In some cases the error comes from the "server" connection, sometimes from the "client".

The error seems to have some relation to the `sendto` being interruptted. In all observed instances at least one of the `sendto`
calls includes `<unfinished ...>` and `<... sendto resumed> )      = -1 EPERM (Operation not permitted)`. However, the error
does not occur everytime `sendto` is interruptted.

Full logs of a run are in [output.log](output.log).

```
close(3)                                = 0
close(3)                                = 0
close(3)                                = 0
Process 5643 attached
Process 5644 attached
Process 5645 attached
Process 5646 attached
[pid  5642] close(3)                    = 0
[pid  5642] socket(PF_INET, SOCK_STREAM, IPPROTO_TCP) = 3
[pid  5642] close(3)                    = 0
[pid  5642] socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP) = 3
[pid  5642] setsockopt(3, SOL_IPV6, IPV6_V6ONLY, [1], 4) = 0
[pid  5642] bind(3, {sa_family=AF_INET6, sin6_port=htons(0), inet_pton(AF_INET6, "::1", &sin6_addr), sin6_flowinfo=0, sin6_scope_id=0}, 28) = 0
[pid  5642] socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP) = 4
[pid  5642] setsockopt(4, SOL_IPV6, IPV6_V6ONLY, [0], 4) = 0
[pid  5642] bind(4, {sa_family=AF_INET6, sin6_port=htons(0), inet_pton(AF_INET6, "::ffff:127.0.0.1", &sin6_addr), sin6_flowinfo=0, sin6_scope_id=0}, 28) = 0
[pid  5642] close(4)                    = 0
[pid  5642] close(3)                    = 0
IN MAIN
[pid  5642] socket(PF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 3
[pid  5642] setsockopt(3, SOL_SOCKET, SO_BROADCAST, [1], 4) = 0
[pid  5642] bind(3, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
[pid  5642] getsockname(3, {sa_family=AF_INET, sin_port=htons(60315), sin_addr=inet_addr("0.0.0.0")}, [16]) = 0
[pid  5642] socket(PF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 5
[pid  5642] setsockopt(5, SOL_SOCKET, SO_BROADCAST, [1], 4) = 0
[pid  5642] bind(5, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
[pid  5642] getsockname(5, {sa_family=AF_INET, sin_port=htons(36289), sin_addr=inet_addr("0.0.0.0")}, [16]) = 0
[pid  5642] sendto(5, "", 0, 0, {sa_family=AF_INET, sin_port=htons(60315), sin_addr=inet_addr("127.0.0.1")}, 16 <unfinished ...>
[pid  5646] sendto(3, "", 0, 0, {sa_family=AF_INET, sin_port=htons(36289), sin_addr=inet_addr("127.0.0.1")}, 16 <unfinished ...>
[pid  5642] <... sendto resumed> )      = 0
[pid  5646] <... sendto resumed> )      = -1 EPERM (Operation not permitted)
[pid  5646] close(3)                    = 0
[pid  5642] close(5)                    = 0
serverConn: <nil>
clientConn: write udp4 0.0.0.0:60315->127.0.0.1:36289: sendto: operation not permitted
[pid  5645] +++ exited with 1 +++
[pid  5644] +++ exited with 1 +++
[pid  5643] +++ exited with 1 +++
[pid  5646] +++ exited with 1 +++
+++ exited with 1 +++
```


Success:
```
close(3)                                = 0
close(3)                                = 0
close(3)                                = 0
Process 316 attached
Process 317 attached
Process 318 attached
[pid   315] close(3)                    = 0
Process 319 attached
[pid   315] socket(PF_INET, SOCK_STREAM, IPPROTO_TCP) = 3
[pid   315] close(3)                    = 0
[pid   315] socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP) = 3
[pid   315] setsockopt(3, SOL_IPV6, IPV6_V6ONLY, [1], 4) = 0
[pid   315] bind(3, {sa_family=AF_INET6, sin6_port=htons(0), inet_pton(AF_INET6, "::1", &sin6_addr), sin6_flowinfo=0, sin6_scope_id=0}, 28) = 0
[pid   315] socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP) = 4
[pid   315] setsockopt(4, SOL_IPV6, IPV6_V6ONLY, [0], 4) = 0
[pid   315] bind(4, {sa_family=AF_INET6, sin6_port=htons(0), inet_pton(AF_INET6, "::ffff:127.0.0.1", &sin6_addr), sin6_flowinfo=0, sin6_scope_id=0}, 28) = 0
[pid   315] close(4)                    = 0
[pid   315] close(3)                    = 0
IN MAIN
[pid   315] socket(PF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 3
[pid   315] setsockopt(3, SOL_SOCKET, SO_BROADCAST, [1], 4) = 0
[pid   315] bind(3, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
[pid   315] getsockname(3, {sa_family=AF_INET, sin_port=htons(54441), sin_addr=inet_addr("0.0.0.0")}, [16]) = 0
[pid   315] socket(PF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 5
[pid   315] setsockopt(5, SOL_SOCKET, SO_BROADCAST, [1], 4) = 0
[pid   315] bind(5, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
[pid   315] getsockname(5, {sa_family=AF_INET, sin_port=htons(58925), sin_addr=inet_addr("0.0.0.0")}, [16]) = 0
[pid   315] sendto(5, "", 0, 0, {sa_family=AF_INET, sin_port=htons(54441), sin_addr=inet_addr("127.0.0.1")}, 16) = 0
[pid   315] close(5)                    = 0
[pid   315] sendto(3, "", 0, 0, {sa_family=AF_INET, sin_port=htons(58925), sin_addr=inet_addr("127.0.0.1")}, 16) = 0
[pid   315] close(3)                    = 0
[pid   319] <... epoll_wait resumed> )  = ? <unavailable>
[pid   319] +++ exited with 0 +++
[pid   318] +++ exited with 0 +++
[pid   317] +++ exited with 0 +++
[pid   316] +++ exited with 0 +++
+++ exited with 0 +++
```