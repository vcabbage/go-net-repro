This repo reproduces a bug observed when two goroutines are sending to each other concurrently over UDP.

I have tried to reproduce on OSX (bare metal) and Ubuntu (Parallels VM), but have only observed the failure within Docker.

Easiest way to reproduce is executing `./run.sh`; provided you have Docker installed. This will build repro.go and execute it under strace repeatedly until it fails.

Note that `run.sh` uses the `--privileged` argument to allow strace to function.

#### Example Output

Failure:

Note the `EPERM (Operation not permitted)` from `sendto`. This error is not documented in `sendto(2)`.
Some searches suggest that it's related to iptables. Disabling iptables on the Docker host does not resolve the error.

In some cases the error comes from the "server" connection, sometimes from the "client".

```
socket(PF_INET, SOCK_STREAM, IPPROTO_TCP) = 3
socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP) = 3
setsockopt(3, SOL_IPV6, IPV6_V6ONLY, [1], 4) = 0
bind(3, {sa_family=AF_INET6, sin6_port=htons(0), inet_pton(AF_INET6, "::1", &sin6_addr), sin6_flowinfo=0, sin6_scope_id=0}, 28) = 0
socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP) = 4
setsockopt(4, SOL_IPV6, IPV6_V6ONLY, [0], 4) = 0
bind(4, {sa_family=AF_INET6, sin6_port=htons(0), inet_pton(AF_INET6, "::ffff:127.0.0.1", &sin6_addr), sin6_flowinfo=0, sin6_scope_id=0}, 28) = 0
socket(PF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 3
setsockopt(3, SOL_SOCKET, SO_BROADCAST, [1], 4) = 0
bind(3, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
getsockname(3, {sa_family=AF_INET, sin_port=htons(45366), sin_addr=inet_addr("0.0.0.0")}, [16]) = 0
socket(PF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 5
setsockopt(5, SOL_SOCKET, SO_BROADCAST, [1], 4) = 0
bind(5, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
getsockname(5, {sa_family=AF_INET, sin_port=htons(37761), sin_addr=inet_addr("0.0.0.0")}, [16]) = 0
sendto(5, "", 0, 0, {sa_family=AF_INET, sin_port=htons(45366), sin_addr=inet_addr("127.0.0.1")}, 16) = -1 EPERM (Operation not permitted)
serverConn: write udp4 0.0.0.0:37761->127.0.0.1:45366: sendto: operation not permitted
(*net.UDPConn)(0xc420026048)({
 conn: (net.conn) {
  fd: (*net.netFD)(0xc420054a80)({
   fdmu: (net.fdMutex) {
    state: (uint64) 0,
    rsema: (uint32) 0,
    wsema: (uint32) 0
   },
   sysfd: (int) 5,
   family: (int) 2,
   sotype: (int) 2,
   isConnected: (bool) false,
   net: (string) (len=4) "udp4",
   laddr: (*net.UDPAddr)(0xc42006c840)(0.0.0.0:37761),
   raddr: (net.Addr) <nil>,
   pd: (net.pollDesc) {
    runtimeCtx: (uintptr) 0x7f8d7e3aae40
   }
  })
 }
})
clientConn: <nil>
(*net.UDPConn)(0xc420026040)({
 conn: (net.conn) {
  fd: (*net.netFD)(0xc420054a10)({
   fdmu: (net.fdMutex) {
    state: (uint64) 0,
    rsema: (uint32) 0,
    wsema: (uint32) 0
   },
   sysfd: (int) 3,
   family: (int) 2,
   sotype: (int) 2,
   isConnected: (bool) false,
   net: (string) (len=4) "udp4",
   laddr: (*net.UDPAddr)(0xc42006c7b0)(0.0.0.0:45366),
   raddr: (net.Addr) <nil>,
   pd: (net.pollDesc) {
    runtimeCtx: (uintptr) 0x7f8d7e3aaf00
   }
  })
 }
})
+++ exited with 1 +++
```


Success:
```
socket(PF_INET, SOCK_STREAM, IPPROTO_TCP) = 3
socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP) = 3
setsockopt(3, SOL_IPV6, IPV6_V6ONLY, [1], 4) = 0
bind(3, {sa_family=AF_INET6, sin6_port=htons(0), inet_pton(AF_INET6, "::1", &sin6_addr), sin6_flowinfo=0, sin6_scope_id=0}, 28) = 0
socket(PF_INET6, SOCK_STREAM, IPPROTO_TCP) = 4
setsockopt(4, SOL_IPV6, IPV6_V6ONLY, [0], 4) = 0
bind(4, {sa_family=AF_INET6, sin6_port=htons(0), inet_pton(AF_INET6, "::ffff:127.0.0.1", &sin6_addr), sin6_flowinfo=0, sin6_scope_id=0}, 28) = 0
socket(PF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 3
setsockopt(3, SOL_SOCKET, SO_BROADCAST, [1], 4) = 0
bind(3, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
getsockname(3, {sa_family=AF_INET, sin_port=htons(39642), sin_addr=inet_addr("0.0.0.0")}, [16]) = 0
socket(PF_INET, SOCK_DGRAM|SOCK_CLOEXEC|SOCK_NONBLOCK, IPPROTO_IP) = 5
setsockopt(5, SOL_SOCKET, SO_BROADCAST, [1], 4) = 0
bind(5, {sa_family=AF_INET, sin_port=htons(0), sin_addr=inet_addr("0.0.0.0")}, 16) = 0
getsockname(5, {sa_family=AF_INET, sin_port=htons(50652), sin_addr=inet_addr("0.0.0.0")}, [16]) = 0
sendto(5, "", 0, 0, {sa_family=AF_INET, sin_port=htons(39642), sin_addr=inet_addr("127.0.0.1")}, 16) = 0
+++ exited with 0 +++
```