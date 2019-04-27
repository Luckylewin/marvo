package libs

import (
    "bytes"
    "crypto/md5"
    "sort"
    "strconv"
    "encoding/hex"
)

func MakeSignature(srcmap map[string]interface{}, bizkey string) string {
    md5ctx := md5.New()
    keys := make([]string, 0, len(srcmap))

    for k := range srcmap {
        if k == "Sign" {
            continue
        }
        keys = append(keys, k)
    }

    sort.Strings(keys)
    var buf bytes.Buffer
    for _, k := range keys {
        vs := srcmap[k]
        if vs == "" {
            continue
        }
        if buf.Len() > 0 {
            buf.WriteByte('&')
        }

        buf.WriteString(k)
        buf.WriteByte('=')

        switch value := vs.(type) {
        case string:
            buf.WriteString(value)
        case int:
            buf.WriteString(strconv.FormatInt(int64(value), 10))
        default:
            panic("params type not supported")
        }
    }

    buf.WriteString(bizkey)
    md5ctx.Write([]byte(buf.String()))

    return hex.EncodeToString(md5ctx.Sum(nil))
}
