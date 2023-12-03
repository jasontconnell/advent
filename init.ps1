param (
    [string]$year = "",
    [string]$day = ""
)

if ($year -eq "") {
    $year = Get-Date -Format "yyyy"
}

if ($day -eq "") {
    $day = Get-Date -Format "dd"
}

initaoc -y $year -d $day