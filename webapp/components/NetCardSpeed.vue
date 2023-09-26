<template>
  <div id="myChart" style="width: 600px;height:400px;"></div>
</template>

<script>
import Vue from 'vue';
export default Vue.extend({
  name: 'NetCardSpeed',

  props: {

  },

  data() {
    return {
      myChart: null
    }
  },
  computed: {
    echartsTitle() {
      console.log('--------------------------- echartsTitle')
      return this.$store.state.sysstat.echartsNetCard.title
    },
    echartsNetCardXAxisData() {
      console.log('--------------------------- echartsNetCardXAxisData')
      return this.$store.state.sysstat.echartsNetCard.xAxisData
    },
    echartsNetCardSeries() {
      console.log('--------------------------- echartsNetCardSeries')

      return this.$store.state.sysstat.echartsNetCard.series
    }
  },
  watch: {
    echartsNetCardSeries(value) {
      console.log('--------------------------- value', value)
    },
    "$store.state.sysstat.echartsNetCardChange": {
      handler: function (newValue, oldValue) {
        console.log('--------------------------- newValue', newValue)
        console.log('--------------------------- oldValue', oldValue)
        console.log('--------------------------- this.myChart', this.myChart)
        if (this.myChart === undefined || this.myChart === null) {
          return
        }
        this.myChart.setOption({
          title: { text: this.echartsTitle },
          tooltip: {},
          xAxis: {
            data: this.echartsNetCardXAxisData
          },
          yAxis: {},
          series: this.echartsNetCardSeries
        })
      },
      deep: true,
      immediate: true
    }
  },
  mounted() {
    this.myChart = this.$echarts.init(document.getElementById('myChart'))
    this.myChart.setOption({
      title: { text: this.echartsTitle },
      tooltip: {},
      xAxis: {
        data: this.echartsNetCardXAxisData
      },
      yAxis: {},
      series: this.echartsNetCardSeries
    })
  }
})
</script>
