import Vue from 'vue'
import Vuex from 'vuex'
import { GetterTree, MutationTree, ActionTree } from 'vuex'
import { User } from '../types'

Vue.use(Vuex)

class State {
  user: User | null = null
  token: any | null = null
}

const getters = <GetterTree<State, any>>{
  user: (state: State) => state.user,
  token: (state: State) => state.token,
}

const mutations = <MutationTree<State>>{
  setUser(state, payload) {
    state.user = payload
  },
  setToken(state, payload) {
    state.token = payload
  },
}

const actions = <ActionTree<State, any>>{}

export const store = new Vuex.Store({
  state: new State(),
  mutations,
  actions,
  getters,
})
