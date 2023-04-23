type record = Record<string | number, string>
type records = Array<record>

export const useToCamelCase = (arr: records): records => {
    const newArr: records = []

    for (let i = 0; i < arr.length; i++) {
        const newObj: record = {}
        for (let key in arr[i]) {
            let newKey = key.replace(/\_(\w)/g, function (all, letter) {
                return letter.toUpperCase();
            })
            newObj[newKey] = arr[i][key]
        }
        newArr.push(newObj)
    }

    return newArr
}