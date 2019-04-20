<?php

$lasttime = 0;

while (true) {
    $now = time();
    if ($now >= $lasttime + 10) {
        $lasttime = $now;
        printf(
            "[%s] %s\n",
            date('Y-m-d H:i:s'),
            trim(file_get_contents('http://myip.ipip.net'))
        );
    }
    time.usleep(100000);
    $statusfile = getenv("TASK_STATUS_FILE");
    if ($statusfile) {
        file_put_contents($statusfile, $now);
    }
}
