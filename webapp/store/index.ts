import { ActionTree, GetterTree, MutationTree } from 'vuex'
import { formatTs } from '~/plugins/utils'

export const state = () => ({
  startup: 0,
})

export type RootState = ReturnType<typeof state>

export const getters: GetterTree<RootState, RootState> = {
  startup: (state) => state.startup,
  startupString: (state) => formatTs(state.startup),
}

export const mutations: MutationTree<RootState> = {
  START_UP(state) {
    state.startup = new Date().getTime()
  },
}

export const actions: ActionTree<RootState, RootState> = {}
