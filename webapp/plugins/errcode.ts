export declare interface Result {
  code: number
  info: string
  data: any
}

export const errCode = {
  failed: -1,
  success: 0,
  ok: 1,
  download_ing: 100200,
  res_caches: 100201,
  support_webrtc: 100202,
  play_not_add: 100203,
  file_stream_not_support: 100204,
  download_begin: 100201,
}
