<script>
  export let selected = [];
  export let options = [];
  export let onchange;
  export let label = "";
  import MultiSelect from "svelte-multiselect";
</script>

<div class="input-div">
  <div class="input-label">{label}</div>
  <div class="input-multiselect_box">
    <MultiSelect
      bind:selected
      bind:options
      on:change={onchange}
      noMatchingOptionsMsg={"ничего не найдено"}
    />
  </div>
</div>

<style>
  :global(div.multiselect > ul.selected > li) {
    width: auto;
    overflow: hidden;
    max-height: 100%;
  }

  .input-multiselect_box {
    font-size: medium;
    margin: auto;
    vertical-align: middle;
    align-items: center;
    flex-wrap: wrap;
    display: flex;
    flex-direction: row;
    padding: 0;
    margin: 0;
  }

  .input-multiselect_box {
    color: var(--tg-theme-text-color);
    padding-right: 10px;
    height: auto;
  }
</style>
