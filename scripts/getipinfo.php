<?php

date_default_timezone_set('Asia/Shanghai');

printf("[%s] %s\n", date('Y-m-d H:i:s'), trim(file_get_contents('http://myip.ipip.net')));

