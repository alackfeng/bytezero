/* eslint-disable promise/always-return */
import { ActionTree, GetterTree, MutationTree } from 'vuex'
import { formatTs } from '~/plugins/utils'
import { RootState } from '~/store'

export declare interface SeriesData {
  name: string
  data: number[]
  type: 'line'
  stack: 'x'
}

export declare interface NetCardSpeedEcharts {
  title: string
  xAxisData: string[]
  series: SeriesData[]
}

export declare interface NetCard {
  Name: string
  Mac: string
  IP: string[]
  Speed: number
}

declare interface NetCardList {
  loading: Boolean
  finished: Boolean
  error: Boolean
  errorCode: Number
  errorText: String
  list: NetCard[]
}

export enum ListOperType {
  LOTNetCard = 1,
  LOTNetCardSubscribe = 2,
  LOTNetCard3 = 3,
  LOTNetCard4 = 4,
}

export const state = () => {
  return {
    listNetCard: {} as NetCardList,
    echartsNetCard: {} as NetCardSpeedEcharts,
    echartsNetCardChange: false as Boolean,
    echartsNetCardName: '' as String,
  }
}
export type SysStatState = ReturnType<typeof state>

export const getters: GetterTree<SysStatState, RootState> = {
  netCardList: (state) => state.listNetCard.list,
  netCardListLength: (state) => state.listNetCard.list.length,
}

export const mutations: MutationTree<SysStatState> = {
  SET_NETCARD_LIST(state, { data }) {
    if (data.Info && data.Info.length > 0) {
      if (state.listNetCard.list === undefined) {
        state.listNetCard.list = [] as NetCard[]
      }
      state.listNetCard.list = [] as NetCard[]
      state.listNetCard.list.push(...data.Info)
    }
    // state.listNetCard.finished = true
  },
  SET_LIST_ERROR(state, { type, res }) {
    if (type === ListOperType.LOTNetCard) {
      state.listNetCard.errorText = `[${res.code}]${res.info}`
      state.listNetCard.errorCode = res.code
      state.listNetCard.error = !(res.code === 0 || res.code === 1)
    }
  },
  SET_LIST_LOADING: (state, { type, value }) => {
    if (type === ListOperType.LOTNetCard) {
      state.listNetCard.loading = value
      if (value) {
        state.listNetCard.finished = false
      }
    }
  },
  SET_NETCARD_SUBSCRIBE(state, { data }) {
    if (data.Info && data.Info.length > 0) {
      state.echartsNetCard.title = '网卡实时速率'
      for (const info of data.Info) {
        // console.log('----------------info: ', info)
        if (state.echartsNetCard.xAxisData === undefined)
          state.echartsNetCard.xAxisData = []

        state.echartsNetCard.xAxisData.push(formatTs(info.NowMs))
        for (let i = 0; i < info.Stat.length; i++) {
          const stat = info.Stat[i]
          if (state.echartsNetCardName !== stat.name) {
            continue
          }
          console.log('----------------stat: ', stat)
          if (state.echartsNetCard.series === undefined)
            state.echartsNetCard.series = []

          const serie = state.echartsNetCard.series.at(i)
          if (serie === undefined) {
            const data: SeriesData = {
              name: stat.name,
              type: 'line',
              stack: 'x',
              data: [stat.bytesRecv],
            }
            state.echartsNetCard.series.push(data)
          } else {
            serie.data.push(stat.bytesRecv)
          }
        }
        // state.echartsNetCard.series.push
      }
      console.log('---------  echartsNetCard', state.echartsNetCard)
      state.echartsNetCardChange = !state.echartsNetCardChange
    }
  },
  SET_SUBSCRIBE_ERROR(_state, { type, res }) {
    if (type === ListOperType.LOTNetCardSubscribe) {
      console.log(
        '-------------------- SET_SUBSCRIBE_ERROR ',
        res.code,
        res.info
      )
    }
  },
  SET_SUBSCRIBE_NAME(state, { Name }) {
    state.echartsNetCardName = Name
  },
}

export const actions: ActionTree<SysStatState, RootState> = {
  // 获取全部网卡列表.
  async getNetCardList({ commit }) {
    const res = await this.$api.getNetCardList()
    console.log('---------------getNetCardList res : ', res)
    if (this.$api.success(res)) {
      commit('SET_NETCARD_LIST', { data: res.data })
    } else {
      commit('SET_LIST_ERROR', { type: ListOperType.LOTNetCard, res })
    }
  },
  async getNetCardSubscribe({ commit }, { Name }) {
    const res = await this.$api.getNetCardSubscribe({ Name })
    console.log('--------------- getNetCardSubscriberes : ', res)
    commit('SET_SUBSCRIBE_NAME', { Name })
    if (this.$api.success(res)) {
      commit('SET_NETCARD_SUBSCRIBE', { data: res.data })
    } else {
      commit('SET_SUBSCRIBE_ERROR', {
        type: ListOperType.LOTNetCardSubscribe,
        res,
      })
    }
  },
}
