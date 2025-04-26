[CmdletBinding()]
param(
    [ValidateRange(1, 1024)]
    [int]     $SizeGB    = 10,
    [string[]]$Filenames = @('disk1.img','disk2.img','disk3.img','disk4.img')
)

$sizeBytes = $SizeGB * 1GB    # 1 073 741 824 × GB

foreach ($file in $Filenames) {
    try {
        $fs = [System.IO.File]::Open(
            $file,
            [System.IO.FileMode]::Create,
            [System.IO.FileAccess]::Write,
            [System.IO.FileShare]::None
        )

        $fs.SetLength($sizeBytes)
        $fs.Close()
        Write-Output "$file を $SizeGB GB で作成しました。"
    }
    catch {
        Write-Error "$file の作成に失敗しました: $($_.Exception.Message)"
        throw
    }
}
