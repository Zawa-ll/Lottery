package conf

import "time"

const SysTimeform = "2006-01-02 15:04:05" // format/parse time
const SysTimeformShort = "2006-01-02"     // Short form

var SysTimeLocation, _ = time.LoadLocation("America/Vancouver")

var SignSecret = []byte("0123456789abcdef")
var CookieSecret = "hellolottery"
