<?php

$signalhandler = function ($signo) {
    if ($signo == SIGINT) {
        printf("\n");
    }
    printf(
        "[%s] recv signal: %d\n",
        date('Y-m-d H:i:s'),
        $signo
    );
    exit(0);
};

pcntl_signal(SIGINT, $signalhandler);
pcntl_signal(SIGHUP, $signalhandler);
pcntl_signal(SIGTERM, $signalhandler);

$lasttime = 0;

while (true) {
    pcntl_signal_dispatch();
    $now = time();
    if ($now >= $lasttime + 10) {
        $lasttime = $now;
        printf(
            "[%s] %s\n",
            date('Y-m-d H:i:s'),
            trim(file_get_contents('http://myip.ipip.net'))
        );
    }
    usleep(100000);
    $statusfile = getenv("TASK_STATUS_FILE");
    if ($statusfile) {
        file_put_contents($statusfile, $now);
    }
}
