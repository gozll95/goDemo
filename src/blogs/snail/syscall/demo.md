root@test1:~# strace ./hello
execve("./hello", ["./hello"], [/* 15 vars */]) = 0
arch_prctl(ARCH_SET_FS, 0x4faea8)       = 0
sched_getaffinity(0, 8192, {3, 0, 0, 0, 0, 0, 0, 0}) = 64
mmap(0xc000000000, 65536, PROT_NONE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0xc000000000
munmap(0xc000000000, 65536)             = 0
mmap(NULL, 262144, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0x7f3f61ef3000
mmap(0xc420000000, 1048576, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0xc420000000
mmap(0xc41fff8000, 32768, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0xc41fff8000
mmap(0xc000000000, 4096, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0xc000000000
mmap(NULL, 65536, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0x7f3f61ee3000
clock_gettime(CLOCK_MONOTONIC, {9194246, 124761395}) = 0
mmap(NULL, 65536, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0x7f3f61ed3000
clock_gettime(CLOCK_MONOTONIC, {9194246, 124861886}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 124884210}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 124905667}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 124927847}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 124957580}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 124982904}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125009223}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125031945}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125055211}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125087039}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125129177}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125155672}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125179152}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125216981}) = 0
rt_sigprocmask(SIG_SETMASK, NULL, [], 8) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125386146}) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 125418252}) = 0
sigaltstack(NULL, {ss_sp=0, ss_flags=SS_DISABLE, ss_size=0}) = 0
sigaltstack({ss_sp=0xc420002000, ss_flags=0, ss_size=32672}, NULL) = 0
rt_sigprocmask(SIG_SETMASK, [], NULL, 8) = 0
gettid()                                = 2735
rt_sigaction(SIGHUP, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGHUP, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGINT, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGINT, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGQUIT, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGQUIT, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGILL, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGILL, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGTRAP, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGTRAP, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGABRT, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGABRT, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGBUS, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGBUS, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGFPE, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGFPE, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGUSR1, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGUSR1, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGSEGV, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGSEGV, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGUSR2, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGUSR2, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGPIPE, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGPIPE, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGALRM, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGALRM, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGTERM, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGTERM, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGSTKFLT, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGSTKFLT, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGCHLD, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGCHLD, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGURG, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGURG, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGXCPU, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGXCPU, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGXFSZ, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGXFSZ, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGVTALRM, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGVTALRM, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGPROF, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGPROF, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGWINCH, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGWINCH, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGIO, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGIO, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGPWR, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGPWR, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGSYS, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGSYS, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRTMIN, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_1, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_2, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_2, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_3, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_3, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_4, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_4, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_5, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_5, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_6, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_6, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_7, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_7, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_8, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_8, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_9, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_9, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_10, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_10, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_11, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_11, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_12, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_12, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_13, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_13, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_14, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_14, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_15, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_15, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_16, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_16, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_17, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_17, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_18, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_18, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_19, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_19, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_20, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_20, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_21, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_21, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_22, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_22, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_23, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_23, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_24, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_24, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_25, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_25, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_26, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_26, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_27, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_27, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_28, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_28, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_29, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_29, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_30, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_30, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_31, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_31, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
rt_sigaction(SIGRT_32, NULL, {SIG_DFL, [], 0}, 8) = 0
rt_sigaction(SIGRT_32, {0x44e650, ~[], SA_RESTORER|SA_STACK|SA_RESTART|SA_SIGINFO, 0x44e780}, NULL, 8) = 0
clock_gettime(CLOCK_MONOTONIC, {9194246, 129038851}) = 0
rt_sigprocmask(SIG_SETMASK, ~[], [], 8) = 0
clone(child_stack=0xc420036000, flags=CLONE_VM|CLONE_FS|CLONE_FILES|CLONE_SIGHAND|CLONE_THREAD) = 2736
rt_sigprocmask(SIG_SETMASK, [], NULL, 8) = 0
rt_sigprocmask(SIG_SETMASK, ~[], [], 8) = 0
clone(child_stack=0xc420032000, flags=CLONE_VM|CLONE_FS|CLONE_FILES|CLONE_SIGHAND|CLONE_THREAD) = 2737
rt_sigprocmask(SIG_SETMASK, [], NULL, 8) = 0
rt_sigprocmask(SIG_SETMASK, ~[], [], 8) = 0
clone(child_stack=0xc420034000, flags=CLONE_VM|CLONE_FS|CLONE_FILES|CLONE_SIGHAND|CLONE_THREAD) = 2738
rt_sigprocmask(SIG_SETMASK, [], NULL, 8) = 0
futex(0x4faf50, FUTEX_WAIT, 0, NULL)    = 0
readlinkat(AT_FDCWD, "/proc/self/exe", "/root/hello", 128) = 11
mmap(NULL, 262144, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0x7f3f61e93000
write(1, "Hello, GopherCon!\n", 18Hello, GopherCon!
)     = 18
exit_group(0)                           = ?
root@test1:~#



root@test1:~# strace -c ./hello
Hello, GopherCon!
% time     seconds  usecs/call     calls    errors syscall
------ ----------- ----------- --------- --------- ----------------
  -nan    0.000000           0         1           write
  -nan    0.000000           0         8           mmap
  -nan    0.000000           0         1           munmap
  -nan    0.000000           0       114           rt_sigaction
  -nan    0.000000           0         8           rt_sigprocmask
  -nan    0.000000           0         3           clone
  -nan    0.000000           0         1           execve
  -nan    0.000000           0         2           sigaltstack
  -nan    0.000000           0         1           arch_prctl
  -nan    0.000000           0         1           gettid
  -nan    0.000000           0         4         1 futex
  -nan    0.000000           0         1           sched_getaffinity
  -nan    0.000000           0        18           clock_gettime
  -nan    0.000000           0         1           readlinkat
------ ----------- ----------- --------- --------- ----------------
100.00    0.000000                   164         1 total


ptrace 
strace