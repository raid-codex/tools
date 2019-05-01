<?php

function url_get($url) {
    set_error_handler(
        function ($severity, $message, $file, $line) {
            throw new ErrorException($message, $severity, $severity, $file, $line);
        }
    );
    try {
        $ret = array(file_get_contents($url), null);
    }
    catch (Exception $e) {
        $ret = array(null, $e);
    }
    restore_error_handler();
    return $ret;
}

?>