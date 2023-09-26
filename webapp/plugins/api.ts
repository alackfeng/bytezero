import { Plugin } from '@nuxt/types'
import { apiUrl, appId } from './config'
import { errCode, Result } from './errcode'

import { uuid } from './utils'

// declare interface DcsDeviceToken {
//   deviceId: string
//   deviceToken: string
//   userId: string
//   dcsCookie: string
//   shareId: string
// }

// export declare enum ReportType {
//   RTPoliticalSensitive = 1,
//   RTSexualExcitement = 2,
//   RTBloodyViolence = 3,
//   RTOther = 4,
// }

// declare interface DevcieShareReport {
//   shareId: string
//   phone: string
//   type: ReportType
//   reason: string
// }

// declare interface ExtLinkOption {
//   shareId: string
//   extractedCode?: string
// }

// declare interface DeviceOption {
//   deviceId: string
// }

// declare interface GetExterLinkResListReq {
//   shareId: string
//   pageNo: number
//   pageSize: number
// }

// declare interface GetLinkDirPathReq {
//   fileId: string
//   pageNo: number
//   pageSize: number
// }

// declare interface CreateDownloadTaskV2Req {
//   fileId: number
//   ower: string
//   fileStart: number
//   small: number
//   begin(fileLength: number): void
//   end(): void
//   process(data: ArrayBuffer): void
// }

// declare interface PlayerTaskReq {
//   url: string
//   item: any
// }

// declare interface DownloadTaskToFileReq {
//   fileId: number
//   ower: string
//   fileStart: number
//   small: number
//   name: string
//   totalLen: number
//   begin(name: string, fileLength: number): void
//   end(): void
//   process(currentLen: number): void
// }

declare interface GetNetCardSubscribeReq {
  Name: string
}

declare interface BridgeUploadFileReq {
  name: string
  size: string
  file: any
}

declare interface IAPI {
  // baseUrl: string
  // appId: string
  // deviceShareGetInfo(option: ExtLinkOption): Promise<Result>
}

const sysStatNetList = '/api/v1/sysstat/net/list'
const sysStatNetSubscribe = '/api/v1/sysstat/net/subscribe'
const bridgeUploadFile = '/api/v1/bridge/files/upload'

class API {
  readonly dtype: number = 2
  readonly reqType: number = 5
  readonly v: string = '2.5'
  readonly v2: string = '2.6'
  readonly language: string = 'zh'
  readonly timeout: number = 10000

  appId = appId
  dcsId: string = ''

  errCode = errCode

  constructor() {
    console.log('--------------------%d', 1)
  }

  ok(res: Result): boolean {
    return res.code === this.errCode.ok
  }

  success(res: Result): boolean {
    return res.code === this.errCode.success
  }

  failed(res: Result): boolean {
    return res.code !== this.errCode.ok
  }

  cached(res: Result): boolean {
    return res.code === this.errCode.res_caches || res.code === this.errCode.ok
  }

  baseApi(url: string): string {
    return apiUrl + url
  }

  nowMs(): number {
    return new Date().getTime()
  }

  hexToBytes(hex: string): any {
    const bytes = []
    for (let c = 0; c < hex.length; c += 2)
      bytes.push(parseInt(hex.substr(c, 2), 16))
    return bytes
  }

  req(cmd: string, data: any) {
    return {
      uuid: uuid(),
      content: {
        cmdName: cmd,
        dtype: this.dtype,
        reqType: this.reqType,
        language: this.language,
        AppId: this.appId,
        v: this.v2,
        sign: '27BC0743C8D3AF5B54B7FE1B02EA3EDC',
        data: window.btoa(JSON.stringify(data)),
      },
    }
  }

  resp(data: any): any {
    if (data.content) {
      return {
        code: data.content.code,
        info: data.content.info,
        data: JSON.parse(window.atob(data.content.data)),
      }
    } else {
      return data
    }
  }

  randomNumber(n: number = 6): string {
    let num = ''
    for (let i = 0; i < n; i++) num += Math.floor(Math.random() * 10)
    return num
  }

  getRecaptcha(): Promise<Result> {
    return new Promise((resolve) => {
      const recaptcha = this.randomNumber()
      resolve({ code: this.errCode.ok, info: 'ok', data: { recaptcha } })
    })
  }

  validateRecaptcha(): Promise<Result> {
    return new Promise((resolve) => {
      resolve({ code: this.errCode.ok, info: 'ok', data: {} })
    })
  }

  getToken(_option: any) {}

  getNetCardList(): Promise<Result> {
    return fetch(this.baseApi(sysStatNetList), {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json;charset=UTF-8',
      },
    }).then((res) => res.json().then((data) => Promise.resolve(data)))
  }

  getNetCardSubscribe(req: GetNetCardSubscribeReq): Promise<Result> {
    return fetch(this.baseApi(sysStatNetSubscribe), {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json;charset=UTF-8',
      },
      body: JSON.stringify(req),
    }).then((res) => res.json().then((data) => Promise.resolve(data)))
  }

  uploadFile(req: BridgeUploadFileReq): Promise<Result> {
    const formData = new FormData()
    formData.append('total', req.size)
    formData.append('files', req.file)
    return fetch(this.baseApi(bridgeUploadFile), {
      method: 'POST',
      mode: 'cors',
      body: formData,
    }).then((res) => res.json().then((data) => Promise.resolve(data)))
  }
}

declare module 'vue/types/vue' {
  // this.$myInjectedFunction inside Vue components
  interface Vue {
    $api: API & IAPI
    // $deviceShareGetInfo(option: ExtLinkOption): Promise<Result>
  }
}

declare module '@nuxt/types' {
  // nuxtContext.app.$myInjectedFunction inside asyncData, fetch, plugins, middleware, nuxtServerInit
  interface NuxtAppOptions {
    $api: API & IAPI
  }
  // nuxtContext.$myInjectedFunction
  interface Context {
    $api: API
  }
}

declare module 'vuex/types/index' {
  // this.$myInjectedFunction inside Vuex stores
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  interface Store<S> {
    $api: API & IAPI
  }
}

const api: Plugin = (context, inject) => {
  const { app } = context
  const api = new API()
  app.api = api
  inject('api', api)
}
export default api
