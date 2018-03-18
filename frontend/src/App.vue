<template>
  <div class="app">
    <div class="legend">Converter Romani et Arabic numbers</div>
    <div class="controls">
      <input type="text" class="input" v-model="input" @keydown="filter" @keyup.enter="convert" @keyup="check" autofocus="true">
      <div class="convert" :class="{ disabled: convertDisabled }" @click="convert"></div>
      <input type="text" class="output" v-model="output" readonly="true">
    </div>
  </div>
</template>

<script>
export default {
  name: 'App',
  data () {
    return {
      input: '',
      output: '',
      convertDisabled: true
    }
  },
  methods: {
    filter: function (event) {
      if (event.key === 'Backspace' || event.key === 'Enter' || event.key === 'Tab') {
        return false
      }

      let key = event.key.toUpperCase()
      if (!/[0-9MDCLXVI]/.test(key)) {
        event.preventDefault()
        return false
      }
    },
    check () {
      let re = /^((\d{0,4})|(M{0,3})(D?C{0,3}|C[DM])(L?X{0,3}|X[LC])(V?I{0,3}|I[VX]))$/
      this.convertDisabled = !re.test(this.input.toUpperCase())

      let n = this.input
      if (!isNaN(n) && (+n < 1 || +n > 3999)) {
        this.convertDisabled = true
      }
    },
    convert () {
      if (this.convertDisabled) {
        return false
      }

      let number = this.input.toUpperCase()
      this.axios.get(`http://localhost:8080/convert/${number}`).then(response => {
        this.output = response.data.Res
      }).catch(err => {
        console.log(err)
      })
    }
  }
}
</script>

<style lang="stylus" scoped>
  @import "assets/styles/variables.styl"

  .app
    margin-top 60px
    padding 0 20px

  .legend
    font-size 40px
    font-weight bold
    text-align center
    margin 20px 0

    @media screen and (max-width 991px)
      font-size 30px

  .controls
    display flex
    justify-content center
    align-items center

    @media screen and (max-width 767px)
      flex-direction column

  input[type=text]
    width 40%
    flex-grow: 1

    @media screen and (max-width 767px)
      width 80%

  .convert
    height 0
    width 0
    margin 0 20px
    border 24px solid transparent
    cursor pointer

    @media screen and (min-width 768px)
      border-left-color $btn-color
      border-right 0

    @media screen and (max-width 767px)
      border-top-color $btn-color
      border-bottom 0

  .disabled
    @media screen and (min-width 768px)
      border-left-color $gray-color
    @media screen and (max-width 767px)
      border-top-color $gray-color
    &:hover
      cursor not-allowed

  .output
    background $main-color
    &:hover
      cursor not-allowed
</style>
