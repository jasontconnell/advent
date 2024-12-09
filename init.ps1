param (
    [string]$year = "",
    [string]$day = ""
)

if ($year -eq "") {
    $year = Get-Date -Format "yyyy"
}

if ($day -eq "") {
    $day = Get-Date -Format "dd"
    if ($day[0] -eq "0") {
        $day = $day[1]
    }
}

initaoc -y $year -d $day