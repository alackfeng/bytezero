export const phoneReg =
  /^1(3[0-9]|4[01456879]|5[0-35-9]|6[2567]|7[0-8]|8[0-9]|9[0-35-9])\d{8}$/

export const recaptchaReg = /^\d{6}$/

export const blobToUrl = (data: any) => window.URL.createObjectURL(data)

export const blobToBase64 = (data: Blob) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.readAsDataURL(data)
    reader.onload = () => resolve(reader.result)
    reader.onerror = (error) => reject(error)
  })
}

export const nowMs = () => new Date().getTime()

export const seed = (radix: number = 16) =>
  Math.floor(Math.random() * 2 ** 18).toString(radix)
export const uuid = (l: number = 4): string => {
  let s: string = ''
  for (let i = 0; i < l; i++) {
    s += seed() + (i + 1 === l ? '' : '-')
  }
  return s
}

export const formatDuring = (mss: any): string => {
  if (typeof mss === 'string') {
    mss = parseInt(mss)
  }
  mss = new Date().getTime() - mss
  const days = Math.floor(mss / (1000 * 60 * 60 * 24))
  const hours = Math.floor((mss % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60))
  const minutes = Math.floor((mss % (1000 * 60 * 60)) / (1000 * 60))
  const seconds = Math.floor((mss % (1000 * 60)) / 1000)
  return days + '天' + hours + '时' + minutes + '分钟' + seconds + '秒'
}

export const formatTs = function (
  d: number,
  fmt: string = 'yyyy-MM-dd hh:mm:ss'
): string {
  const date = new Date(d)
  const o: any = {
    'M+': date.getMonth() + 1,
    'd+': date.getDate(),
    'h+': date.getHours(),
    'm+': date.getMinutes(),
    's+': date.getSeconds(),
    'q+': Math.floor((date.getMonth() + 3) / 3),
    S: date.getMilliseconds(),
  }
  if (/(y+)/.test(fmt)) {
    fmt = fmt.replace(
      RegExp.$1,
      (date.getFullYear() + '').substr(4 - RegExp.$1.length)
    )
  }
  for (const k in o) {
    if (new RegExp('(' + k + ')').test(fmt)) {
      fmt = fmt.replace(
        RegExp.$1,
        RegExp.$1.length === 1 ? o[k] : ('00' + o[k]).substr(('' + o[k]).length)
      )
    }
  }
  return fmt
}
