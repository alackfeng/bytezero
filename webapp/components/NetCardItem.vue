<template>
  <div class="box">
    <van-cell clickable show-right @click="onClick">
      <template #icon>
        <div>
          <van-image block width="28" height="28" :src="thumb" @click="onClickIcon"></van-image>
        </div>
      </template>
      <template #title>
        <div class="box-title progress" :style="progress">{{ title }}</div>
      </template>
      <template #label>
        <div class="box-title">{{ subTitle }}</div>
      </template>
    </van-cell>
  </div>
</template>

<script lang="ts">
import Vue from 'vue'
import { mapActions, mapGetters } from 'vuex'
import { NetCard } from '~/store/sysstat'

export default Vue.extend({
  name: 'NetCardItem',
  props: {
    prop: {
      type: Object,
      default() {
        return {} as NetCard
      },
    },
    index: {
      type: Number,
      default: 0,
    },
  },
  data() {
    return {
      checked: false,
      clicked: true,
      progress: '--progress: 0%',
    }
  },

  computed: {
    thumb() {
      return require('@/static/icon.png')
    },
    title(): string {
      return `${this.prop.Name} - ${this.prop.Mac} `
    },
    subTitle(): string {
      return this.prop.IP
    },
    showRight(): boolean {
      return this.prop?.isFile
    },
  },
  watch: {
    prop(_item) {
      // console.log("ExtLinkItem.watch() called", item)
      // 切换页面时可以更新list数据.
      // eslint-disable-next-line no-return-assign, promise/catch-or-return
      // this.getThumbnail({ item }).then((url) => (this.thumb = url))
      console.log('-----------item: ', _item)
    },
  },
  mounted() {
    console.log('-------------prop: ', this.prop)
  },
  methods: {
    ...mapGetters('global', ['thumbDeaultGet']),
    ...mapActions('global', ['getThumbnail']),

    onChecked() {
      this.clicked = false
      this.$emit('checked', {
        checked: this.checked,
        item: this.prop,
        index: this.index,
        ui: this,
      })
    },
    onClick() {
      this.$emit('click', { index: this.index, item: this.prop, ui: this })
    },
    onClickIcon() {
      // this.clicked = false
      // this.$emit('icon', { index: this.index, item: this.prop, myself: this  })
    },
    setProgress(current: number) {
      // console.log('current',  (current*100/this.prop.fileLen).toFixed(0))
      this.progress =
        '--progress: ' + ((current * 100) / this.prop.fileLen).toFixed(0) + '%'
    },
  },
})
</script>

<style>
.box {
  display: flex;
}

.box-title {
  margin-left: 5px;
  min-width: 375px;
  max-width: 750px;
  width: 80%;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.progress {
  background: linear-gradient(90deg, #0f0, #0ff var(--progress), transparent 0);
}

.box-name {
  background-color: aqua;
}

.checkbox {
  background-color: transparent;
  padding: 16px;
}
</style>
