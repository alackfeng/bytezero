<template>
  <div class="bg bg-color">
    <van-col>
      <van-row>
        <h2 class="bg-color title">{{ $t('sysstat.monitor.title') }}</h2>
      </van-row>
      <van-row class="bg-box">
        <van-col>
          <h5>网卡信息:</h5>
          <van-list v-model="loading" class="content" :finished="finished" :finished-text="finishedText"
            :error.sync="errorCode" :error-text="errorText" @load="onLoad">
            <van-checkbox-group ref="checkboxGroup" v-model="selectedAll">
              <net-card-item v-for="(item, index) in netCardList" :key="index" :prop="item" :index="index"
                @click="onClickNetCard">
              </net-card-item>
            </van-checkbox-group>
          </van-list>
        </van-col>
        <van-col>
          <h5>实时网速:</h5>
          <net-card-speed></net-card-speed>
        </van-col>
      </van-row>

      <van-row class="bg-color">
        <h2>文件上传测试</h2>
        <net-card-upload></net-card-upload>
      </van-row>
    </van-col>
  </div>
</template>

<script lang="ts">
import Vue from 'vue';
import { ListOperType } from '~/store/sysstat';
export default Vue.extend({
  name: 'SysStatMonitorPage',

  async asyncData({ store }) {
    console.log('------------res : 111')
    const res = await store.dispatch('sysstat/getNetCardList')
    console.log('------------res: ', res)
  },
  data() {
    return {
      finishedText: this.$t('extlink.no_more'),
      selectedAll: [],
    }
  },
  computed: {
    finished() { return this.$store.state.sysstat.finished },
    errorCode() { return this.$store.state.sysstat.error },
    errorText() { return this.$store.state.sysstat.errorText },
    netCardList() {
      console.log('------------- this.$store.state.sysstat.listNetCard.list ', this.$store.state.sysstat.listNetCard.list)
      return this.$store.state.sysstat.listNetCard.list
    },
    listNetCard() {
      return this.$store.state.sysstat.listNetCard
    },
    loading: {
      get(): boolean {
        return this.$store.state.sysstat.listNetCard.loading
      },
      set(value: boolean) {
        this.$store.commit('sysstat/SET_LIST_LOADING', { type: ListOperType.LOTNetCard, vaule: value })
      },
    }
  },
  mounted() {
  },
  methods: {
    onLoad() {
      console.log('--------------------onLoad ')
    },
    async onClickNetCard({ _index, item, _ui }: any) {
      console.log('--------imte: ', item.Name)
      const res = await this.$store.dispatch('sysstat/getNetCardSubscribe', { Name: item.Name })
      console.log('------------onClickNetCard res: ', res)
    },
  }
})
</script>

<style>
.bg {
  width: 960px;
  margin: auto;
}

.title {
  text-align: center;
}

.content {
  background-color: aquamarine;
}

.bg-color {
  background-color: gray;
}

.bg-box {
  background-color: gray;
  display: flex;
  flex-direction: row;
}
</style>
