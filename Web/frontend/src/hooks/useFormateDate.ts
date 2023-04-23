type stringKey = Record<string, string>

export function useFormateDate(dateTime: Date | null, fmt = 'yyyy-mm-dd hh:MM:ss') {
    // 如果为null,则格式化当前时间
    if (!dateTime) dateTime = new Date()
    let ret
    let opt: stringKey = {
        y: dateTime.getFullYear().toString(), // 年
        m: (dateTime.getMonth() + 1).toString(), // 月
        d: dateTime.getDate().toString(), // 日
        h: dateTime.getHours().toString(), // 时
        M: dateTime.getMinutes().toString(), // 分
        s: dateTime.getSeconds().toString() // 秒
    }
    for (let k in opt) {
        ret = new RegExp("(" + k + "+)").exec(fmt)
        if (ret) {
            fmt = fmt.replace(ret[1], (ret[1].length == 1) ? opt[k] : (opt[k].padStart(ret[1].length, "0")))
        }
    }
    return fmt
}